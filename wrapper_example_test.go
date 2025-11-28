// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nig

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Logger2(base zerolog.Logger) *Dep[zerolog.Logger] {
	return NewDep(func(_ *gin.Context) zerolog.Logger { return base })
}

func ReqID2(l *Dep[zerolog.Logger]) *Dep[string] {
	reqid := NewDep(requestid.Get, requestid.New())
	return Updates(l, reqid, func(l zerolog.Logger, reqid string) zerolog.Logger {
		return l.With().Str("request_id", reqid).Logger()
	})
}

type HelloArgs2 struct {
	Log       zerolog.Logger `nig:"logger"`
	RequestID string         `nig:"request_id"`
}

func Example_staticApproach() {
	// discard gin log to prevent output pollution
	gin.DefaultWriter = io.Discard
	g := gin.Default()

	baseLogger := zerolog.New(os.Stdout)
	logger := Logger2(baseLogger)
	reqid := ReqID2(logger)

	helloArgs := NewArg(func(c *gin.Context) HelloArgs2 {
		return HelloArgs2{
			Log:       logger.Get(c),
			RequestID: reqid.Get(c),
		}
	}, logger, reqid)

	Use(helloArgs, g).GET("/", func(c *gin.Context, use HelloArgs2) {
		use.Log.Info().Msg("hello")
		c.JSON(200, gin.H{"data": "hello", "request_id": use.RequestID})
	})

	fmt.Println("Starting server...")
	// actually run the server
	// g.Run(":8080")

	//output: Starting server...
}
