package main

import (
	"auto/cmd/server/config"
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		logrus.Fatalln("main: invalid configuration")
	}
	app := iris.New()
	app.Logger().SetLevel("info")
	app.Use(recover.New())
	app.Use(logger.New())
	routes(app, cfg)
	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		app.Shutdown(ctx)
	})
	app.Run(iris.Addr(cfg.Server.Addr), iris.WithoutServerError(iris.ErrServerClosed))
}
