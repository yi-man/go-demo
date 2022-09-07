package main

import (
	"context"
	"fmt"

	"go.demo/examples/biz"
	"go.demo/examples/conf"
	"go.demo/examples/dal"
	"go.demo/examples/dal/query"
)

func init() {
	dal.DB = dal.ConnectDB(conf.MySQLDSN).Debug()
}

func main() {
	// start your project here
	fmt.Println("hello world")
	defer fmt.Println("bye~")

	query.SetDefault(dal.DB)
	biz.Query(context.Background())
}
