package main

import (
	"fmt"
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()

	if err := app.Run(iris.Addr("0.0.0.0:3001")); err != nil {
		panic(fmt.Errorf("服务器监听错误：%s", err))
	}
}
