package main

import (
	"context"

	"go.uber.org/fx"

	"github.com/rfancn/airpress/cache"
	"github.com/rfancn/airpress/config"
	"github.com/rfancn/airpress/dal"
	"github.com/rfancn/airpress/event"
	"github.com/rfancn/airpress/event/listener"
	"github.com/rfancn/airpress/handler"
	"github.com/rfancn/airpress/handler/middleware"
	"github.com/rfancn/airpress/injection"
	"github.com/rfancn/airpress/log"
	"github.com/rfancn/airpress/template"
	"github.com/rfancn/airpress/template/extension"
)

var eventBus event.Bus

// @title           AirPress API
// @version         1.0
// @description     AirPress 博客系统后端 API 文档,通过 swag 注释驱动自动生成。
// @description     包含后台管理(/api/admin)与前台内容(/api/content)两组接口。
// @host            localhost:8080
// @BasePath        /api
// @securityDefinitions.apikey AdminApiKey
// @in header
// @name Admin-Authorization
// @securityDefinitions.apikey ContentApiKey
// @in header
// @name Content-Authorization
func main() {
	app := InitApp()

	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
	eventBus.Publish(context.Background(), &event.StartEvent{})
	<-app.Done()
}

func InitApp() *fx.App {
	options := injection.GetOptions()
	options = append(options,
		fx.NopLogger,
		fx.Provide(
			log.NewLogger,
			log.NewGormLogger,
			event.NewSyncEventBus,
			dal.NewGormDB,
			cache.NewCache,
			config.NewConfig,
			handler.NewServer,
			template.NewTemplate,
			middleware.NewAuthMiddleware,
			middleware.NewGinLoggerMiddleware,
			middleware.NewRecoveryMiddleware,
			middleware.NewInstallRedirectMiddleware,
		),
		fx.Populate(&dal.DB),
		fx.Populate(&eventBus),
		fx.Invoke(
			listener.NewStartListener,
			listener.NewTemplateConfigListener,
			listener.NewLogEventListener,
			listener.NewPostUpdateListener,
			listener.NewCommentListener,
			extension.RegisterCategoryFunc,
			extension.RegisterCommentFunc,
			extension.RegisterTagFunc,
			extension.RegisterMenuFunc,
			extension.RegisterPhotoFunc,
			extension.RegisterLinkFunc,
			extension.RegisterToolFunc,
			extension.RegisterPaginationFunc,
			extension.RegisterPostFunc,
			extension.RegisterStatisticFunc,
			func(s *handler.Server) {
				s.RegisterRouters()
			},
		),
	)
	app := fx.New(
		options...,
	)
	return app
}
