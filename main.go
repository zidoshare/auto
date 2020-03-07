package main

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"gitlab.scustartup.com/hnqc/auto/config"
	"time"
)

func main() {
	fmt.Printf("%++v\n", config.Config())
	panic("error")
	app := iris.New()
	app.Logger().SetLevel("info")
	app.Use(recover.New())
	app.Use(logger.New())
	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		app.Shutdown(ctx)
	})
	app.Run(iris.Addr(":3001"), iris.WithoutServerError(iris.ErrServerClosed))
}
