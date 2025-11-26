// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nig

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
)

// handlerInfo contains the parsed information about a handler function
type handlerInfo struct {
	handlerFunc reflect.Value
	argStructs  []reflect.Type
}

// fieldInfo contains information about a struct field and its dependency
type fieldInfo struct {
	fieldIndex int
	fieldType  reflect.Type
	depKey     string
}

// validateHandlerType checks if the handler has a valid signature:
// func(*gin.Context, struct1, struct2, ...)
// Returns the handler info and error if validation fails.
func validateHandlerType(handler any) (handlerInfo, error) {
	handlerType := reflect.TypeOf(handler)

	// Must be a function
	if handlerType.Kind() != reflect.Func {
		return handlerInfo{}, fmt.Errorf("handler must be a function")
	}

	// Must have at least 2 parameters (context + at least one struct)
	if handlerType.NumIn() < 2 {
		return handlerInfo{}, fmt.Errorf("handler must accept *gin.Context and at least one struct parameter")
	}

	// First parameter must be *gin.Context
	firstParam := handlerType.In(0)
	contextType := reflect.TypeOf((*gin.Context)(nil))
	if firstParam != contextType {
		return handlerInfo{}, fmt.Errorf("handler's first parameter must be *gin.Context")
	}

	// All remaining parameters must be structs
	argStructs := make([]reflect.Type, 0, handlerType.NumIn()-1)
	for i := 1; i < handlerType.NumIn(); i++ {
		paramType := handlerType.In(i)
		if paramType.Kind() != reflect.Struct {
			return handlerInfo{}, fmt.Errorf("handler's parameter %d must be a struct, got %v", i, paramType.Kind())
		}
		argStructs = append(argStructs, paramType)
	}

	// Handler should not return anything (or we ignore return values)
	if handlerType.NumOut() > 0 {
		return handlerInfo{}, fmt.Errorf("handler must not return anything")
	}

	return handlerInfo{
		handlerFunc: reflect.ValueOf(handler),
		argStructs:  argStructs,
	}, nil
}

// parseStructTags parses the struct tags for all fields in the given struct types.
// Returns a map: struct index -> []fieldInfo
// Returns error if any field doesn't have a "nig" tag.
func parseStructTags(structs []reflect.Type) (map[int][]fieldInfo, error) {
	result := make(map[int][]fieldInfo)

	for structIdx, structType := range structs {
		numFields := structType.NumField()
		fields := make([]fieldInfo, 0, numFields)

		for fieldIdx := range numFields {
			field := structType.Field(fieldIdx)

			// Get the "nig" tag
			depKey, ok := field.Tag.Lookup("nig")
			if !ok || depKey == "" {
				return nil, fmt.Errorf("field %s.%s must have a non-empty 'nig' struct tag",
					structType.Name(), field.Name)
			}

			fields = append(fields, fieldInfo{
				fieldIndex: fieldIdx,
				fieldType:  field.Type,
				depKey:     depKey,
			})
		}

		result[structIdx] = fields
	}

	return result, nil
}

// validateDependencies checks that all required dependencies are registered
// and their types match the struct fields.
// Returns a map: depKey -> anotherDep for middleware collection.
// Returns error if validation fails.
func validateDependencies(deps map[string]any, fieldsMap map[int][]fieldInfo) (map[string]anotherDep, error) {
	depInstances := make(map[string]anotherDep)

	for _, fields := range fieldsMap {
		for _, field := range fields {
			// Check if dependency is registered
			dep, exists := deps[field.depKey]
			if !exists {
				return nil, fmt.Errorf("dependency '%s' is not registered", field.depKey)
			}

			// Check if it implements anotherDep (i.e., it's a *Dep[T])
			depInstance, ok := dep.(anotherDep)
			if !ok {
				return nil, fmt.Errorf("dependency '%s' is not a valid Dep instance", field.depKey)
			}

			// Verify type compatibility by attempting to get the value
			// We need to check if the Dep's type parameter matches the field type
			// This is tricky with reflection and generics...
			// For now, we'll store the dep and do runtime type checking during execution

			depInstances[field.depKey] = depInstance
		}
	}

	return depInstances, nil
}

// collectAndSortMiddlewares collects all middlewares from dependencies
// and sorts them using topological sort to handle dependency order.
// Returns the sorted middleware list.
// Panics if there's a cyclic dependency.
func collectAndSortMiddlewares(depInstances map[string]anotherDep) ([]gin.HandlerFunc, error) {
	// gather deps
	deps := make([]anotherDep, 0, len(depInstances))
	for _, d := range depInstances {
		deps = append(deps, d)
	}

	// sort dependencies
	deps, err := sortDeps(deps)
	if err != nil {
		return nil, err
	}

	// calculate number of middlewares
	size := 0
	for _, d := range deps {
		size += len(d.setup())
	}
	ret := make([]gin.HandlerFunc, 0, size)

	for _, d := range deps {
		ret = append(ret, d.setup()...)
	}

	return ret, nil
}

// wrapHandler creates a gin.HandlerFunc that extracts dependencies from context,
// constructs the argument structs, and calls the original handler.
func wrapHandler(info handlerInfo, fieldsMap map[int][]fieldInfo) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prepare arguments for the handler
		args := make([]reflect.Value, len(info.argStructs)+1)
		args[0] = reflect.ValueOf(c) // First arg is always *gin.Context

		// Construct each struct argument
		for structIdx, structType := range info.argStructs {
			structValue := reflect.New(structType).Elem()
			fields := fieldsMap[structIdx]

			for _, field := range fields {
				// Get the dependency value from context
				depValue, exists := c.Get(field.depKey)
				if !exists {
					// This shouldn't happen if middlewares are set up correctly
					panic(fmt.Sprintf("dependency '%s' not found in context", field.depKey))
				}

				// Set the field value
				fieldValue := structValue.Field(field.fieldIndex)
				depReflectValue := reflect.ValueOf(depValue)

				// Type check
				if !depReflectValue.Type().AssignableTo(fieldValue.Type()) {
					panic(fmt.Sprintf("dependency '%s' type mismatch: cannot assign %v to %v",
						field.depKey, depReflectValue.Type(), fieldValue.Type()))
				}

				fieldValue.Set(depReflectValue)
			}

			args[structIdx+1] = structValue
		}

		// Call the original handler
		info.handlerFunc.Call(args)
	}
}
