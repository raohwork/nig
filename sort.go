// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nig

import (
	"errors"
	"slices"

	"github.com/ninedraft/tsort"
)

type simpleFifo interface {
	add(deps ...anotherDep)
	get() (anotherDep, bool)
}

type fifoImpl []anotherDep

func (f *fifoImpl) add(deps ...anotherDep) {
	*f = append(*f, deps...)
}
func (f *fifoImpl) get() (ret anotherDep, ok bool) {
	if len(*f) < 1 {
		return
	}

	ret, ok = (*f)[0], true
	*f = (*f)[1:]
	return
}

func newFifo() simpleFifo {
	return &fifoImpl{}
}

func sortDeps(deps []anotherDep) ([]anotherDep, error) {
	if len(deps) <= 1 {
		return deps, nil
	}
	// create graph
	edges := map[anotherDep][]anotherDep{}
	fifo := newFifo()
	fifo.add(deps...)
	for cur, ok := fifo.get(); ok; cur, ok = fifo.get() {
		if _, ok := edges[cur]; !ok {
			edges[cur] = nil
		}
		x := cur.dependsOn()
		for _, v := range x {
			if !slices.Contains(edges[v], cur) {
				edges[v] = append(edges[v], cur)
			}
		}
		fifo.add(x...)
	}
	points := make([]anotherDep, 0, len(edges))
	for k := range edges {
		points = append(points, k)
	}

	sorted, cyclic := tsort.Sort(points, func(i anotherDep) []anotherDep {
		return edges[i]
	})
	if cyclic {
		return nil, errors.New("cyclic dependency detected")
	}

	return sorted, nil
}
