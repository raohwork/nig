// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nig

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// fake zerolog.Logger
type ZerologLogger struct{}

// fake gorm.DB
type GormDB struct{}

// fake github.com/gin-contrib/requestid
type RequestID struct{}

type baseServices struct {
	Logger ZerologLogger
	DB     *GormDB
}

type MyClaim struct {
	// Add fields for your JWT claim
}

var useReqid = CreateDep(
	"request_id",
	func(c *gin.Context) (string, bool) {
		// reqid := requestid.Get(c)
		reqid := "fake-request-id" // replace with actual request ID extraction logic
		if reqid == "" {
			return "", false
		}
		return reqid, true
	},
)

var useLogger = CreateDep(
	"logger",
	func(c *gin.Context) (ZerologLogger, bool) {
		// l := log.With().
		// 	Str("client_ip", c.ClientIP()).
		// 	Str("req_id", useReqid.Get(c)).
		// 	Logger()
		l := ZerologLogger{}
		return l, true
	},
)

func fakeExtractJWT(code string) (*MyClaim, error) {
	return nil, errors.New("fake JWT extraction not implemented")
}

// useJWT extracts JWT claims from the request context. It also depends on the logger
var useJWT = CreateDep(
	"jwt_claim",
	func(c *gin.Context) (ret *MyClaim, ok bool) {
		if !useLogger.Set(c) {
			return
		}

		claim, err := fakeExtractJWT("token code extracted from request")
		if err != nil {
			return nil, false
		}

		// useLogger.Update(
		//     c,
		//     useLogger.Get(c).With().Interface("jwt", claim).Logger,
		// )

		return claim, true
	},
)

func Example() {
	db := &GormDB{} // fake gorm.DB instance

	g := gin.Default()
	// g.Use(cors.Default())
	// g.Use(requestid.New())

	useDB := Use(db, "db")

	baseTool := NewArg(func() gin.HandlersChain {
		return gin.HandlersChain{useLogger.Setup, useDB.Setup}
	}, func(c *gin.Context) baseServices {
		return baseServices{
			Logger: useLogger.Get(c),
			DB:     useDB.Get(c),
		}
	})

	With(g, baseTool).
		GET("/example", func(c *gin.Context, use baseServices) {
			// use.Logger.Debug().Msg("test") // client_ip is set
			// use.DB.Exec("SELECT 1")
		})

	With2Args(g, baseTool, FromDep(useJWT)).
		// user is never nil since useJWT never returns (nil, true)
		GET("/example-with-jwt", func(c *gin.Context, use baseServices, user *MyClaim) {
			// use.Logger.Debug().Msg("test") // client_ip/jwt are set
			// use.DB.Exec("SELECT 1")
			// c.JSON(200, gin.H{"name": user.Name})
		})

	//output:
}
