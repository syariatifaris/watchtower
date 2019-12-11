package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/syariatifaris/watchtower"
)

type (
	redisDummy struct{}
	dbDummy    struct{}
)

var (
	rds *redisDummy
	db  *dbDummy
)

func init() {
	rds = new(redisDummy)
	db = new(dbDummy)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	done := make(chan bool)

	tw := watchtower.New(true)

	redisFixable := watchtower.Fixable{
		Name: "redis object health",
		Err:  "redis object is nil",
		Healthy: func() bool {
			return rds != nil
		},
		Fix: func() error {
			fmt.Println("fixing redis")
			rds = new(redisDummy)
			return nil
		},
	}

	dbFixable := watchtower.Fixable{
		Name: "database object health",
		Err:  "database nil",
		Healthy: func() bool {
			return db != nil
		},
		Fix: func() error {
			fmt.Println("fixing the db")
			db = new(dbDummy)
			return nil
		},
	}

	tw.AddWatchObject(redisFixable, dbFixable)
	go tw.Run(done)

	go func() {
		time.Sleep(time.Second * 2)
		fmt.Println("set nil")
		rds = nil
	}()

	go func() {
		time.Sleep(time.Second * 10)
		fmt.Println("set nil 2nd")
		db = nil
	}()

	select {
	case <-ctx.Done():
		done <- true
		os.Exit(1)
	}
}
