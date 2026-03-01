package main

import (
	"go.uber.org/fx"

	"git.dev.siap.id/kukuhkkh/app-silsilah/app/middleware"
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/module/auth"
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/router"
	"git.dev.siap.id/kukuhkkh/app-silsilah/internal/bootstrap"
	"git.dev.siap.id/kukuhkkh/app-silsilah/internal/bootstrap/database"
	"git.dev.siap.id/kukuhkkh/app-silsilah/utils/config"
	"git.dev.siap.id/kukuhkkh/app-silsilah/utils/session"
	fxzerolog "github.com/efectn/fx-zerolog"
	_ "go.uber.org/automaxprocs"
)

// @title                       Aplikasi Silsilah API
// @version                     1.0
// @description                 API untuk aplikasi silsilah keluarga.
// @contact.name                Kukuh Rahmadani
// @contact.email               krahmadani1@gmail.com
// @license.name                MIT
// @license.url                 https://opensource.org/license/mit
// @host                        localhost:8080
// @schemes                     http https
// @BasePath                    /
func main() {
	fx.New(
		// config
		fx.Provide(config.NewConfig),
		// logging
		fx.Provide(bootstrap.NewLogger),
		// fiber
		fx.Provide(bootstrap.NewFiber),
		// database
		fx.Provide(database.NewDatabase),
		// session
		fx.Provide(session.NewStore),
		// middleware
		fx.Provide(middleware.NewMiddleware),
		fx.Provide(middleware.NewAuthMiddleware),
		// router
		fx.Provide(router.NewRouter),

		// provide modules
		auth.NewAuthModule,

		// start application
		fx.Invoke(bootstrap.Start),

		// define logger
		fx.WithLogger(fxzerolog.Init()),
	).Run()
}
