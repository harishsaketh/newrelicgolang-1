package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	nr "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgin/v1"
)

type key int

const (
   keyNrID key = iota
)

var (
    nrapp newrelic.Application
)

func main() {
	initNewRelic()
	r := gin.Default()
	r.Use(nrgin.Middleware(nrapp))
	r.Use(setNewRelicInContext())

	setUpRoutes(r)
	// Listen and Server in 0.0.0.0:8080
	s := &http.Server{
		Addr:         ":8080",
		Handler:      r,
	}
	s.ListenAndServe()
}

//populateNewRelicInContext get the request context populated
func setNewRelicInContext() gin.HandlerFunc {

	return func(c *gin.Context) {
		//Setup context
		ctx := c.Request.Context()

		//Set newrelic context
		var txn nr.Transaction
		//newRelicTransaction is the key populated by nrgin Middleware
		value, exists := c.Get("newRelicTransaction")
		if exists {
			if v, ok := value.(nr.Transaction); ok {
				txn = v
			}
			ctx = context.WithValue(ctx, keyNrID, txn)
		}
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}


func initNewRelic() {
	var err error
	nrConfig := newrelic.NewConfig("test", "somekey")
	nrapp, err = newrelic.NewApplication(nrConfig)
	if err != nil {
		panic("Failed to setup NewRelic: " + err.Error())
	}
}
