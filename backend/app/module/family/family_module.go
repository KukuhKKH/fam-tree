package family

import (
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/middleware"
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/module/family/controller"
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/module/family/repository"
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/module/family/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type FamilyRouter struct {
	App            *fiber.App
	Controller     *controller.Controller
	AuthMiddleware *middleware.AuthMiddleware
}

var NewFamilyModule = fx.Options(
	fx.Provide(repository.NewFamilyRepository),
	fx.Provide(repository.NewFamilyMemberRepository),
	fx.Provide(service.NewFamilyService),
	fx.Provide(controller.NewController),
	fx.Provide(controller.NewFamilyController),
	fx.Provide(controller.NewPublicController),
	fx.Provide(NewFamilyRouter),
)

func NewFamilyRouter(app *fiber.App, controller *controller.Controller, authMw *middleware.AuthMiddleware) *FamilyRouter {
	return &FamilyRouter{
		App:            app,
		Controller:     controller,
		AuthMiddleware: authMw,
	}
}

func (r *FamilyRouter) RegisterFamilyRoutes() {
	familyCtrl := r.Controller.Family

	_i := r.App.Group("/families")

	// Must be logged in for all family routes
	_i.Use(r.AuthMiddleware.RequireAuth())

	// CRUD Family
	_i.Post("/", familyCtrl.Create)
	_i.Get("/", familyCtrl.GetMyFamilies)

	// slug routes
	slugGrp := _i.Group("/:slug")
	slugGrp.Get("/", familyCtrl.GetFamilyDetails)
	slugGrp.Patch("/", familyCtrl.UpdateFamily)
	slugGrp.Delete("/", familyCtrl.DeleteFamily)

	// Members
	membersGrp := slugGrp.Group("/members")
	membersGrp.Post("/", familyCtrl.InviteMember)
	membersGrp.Get("/", familyCtrl.GetMembers)
	membersGrp.Patch("/:user_id", familyCtrl.UpdateMemberRole)
	membersGrp.Delete("/:user_id", familyCtrl.RemoveMember)

	// ─── Public routes (no auth) ──────────────────────────────────────────────
	publicCtrl := r.Controller.Public

	pub := r.App.Group("/public/families")
	pub.Get("/:slug", publicCtrl.GetPublicFamily)
}
