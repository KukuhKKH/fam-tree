package router

import (
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/module/auth"
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/module/family"
	"git.dev.siap.id/kukuhkkh/app-silsilah/internal/bootstrap/database"
	"git.dev.siap.id/kukuhkkh/app-silsilah/utils/config"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	App          *fiber.App
	Cfg          *config.Config
	Db           *database.Database
	AuthRouter   *auth.AuthRouter
	FamilyRouter *family.FamilyRouter
}

func NewRouter(
	fiber *fiber.App,
	cfg *config.Config,
	db *database.Database,
	authRouter *auth.AuthRouter,
	familyRouter *family.FamilyRouter,
) *Router {
	return &Router{
		App:          fiber,
		Cfg:          cfg,
		Db:           db,
		AuthRouter:   authRouter,
		FamilyRouter: familyRouter,
	}
}

func (r *Router) Register() {
	// Root Redirect to Frontend
	r.App.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect(r.Cfg.App.FrontendUrl)
	})

	// Health Check
	r.App.Get("/health", func(c *fiber.Ctx) error {
		status := fiber.Map{
			"status":   "ok",
			"database": "ok",
		}

		// Check DB connection
		sqlDB, err := r.Db.DB.DB()
		if err != nil {
			status["database"] = "disconnected"
			return c.Status(fiber.StatusInternalServerError).JSON(status)
		}

		if err := sqlDB.Ping(); err != nil {
			status["database"] = "error: " + err.Error()
			return c.Status(fiber.StatusInternalServerError).JSON(status)
		}

		return c.JSON(status)
	})

	// Test Routes
	r.App.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Pong! 👋")
	})

	// routes of modules
	r.AuthRouter.RegisterAuthRoutes()
	r.FamilyRouter.RegisterFamilyRoutes()
}
