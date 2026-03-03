package person

import (
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/middleware"
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/module/person/controller"
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/module/person/repository"
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/module/person/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type PersonRouter struct {
	App            *fiber.App
	Controller     PersonController
	AuthMiddleware *middleware.AuthMiddleware
}

// PersonController bundles all person-related controllers.
type PersonController struct {
	Person controller.PersonController
}

var NewPersonModule = fx.Options(
	fx.Provide(repository.NewPersonRepository),
	fx.Provide(repository.NewRelationshipRepository),
	fx.Provide(service.NewPersonService),
	fx.Provide(newPersonController),
	fx.Provide(NewPersonRouter),
)

func newPersonController(svc service.PersonService) *PersonController {
	return &PersonController{
		Person: controller.NewPersonController(svc),
	}
}

func NewPersonRouter(app *fiber.App, ctrl *PersonController, authMw *middleware.AuthMiddleware) *PersonRouter {
	return &PersonRouter{
		App:            app,
		Controller:     *ctrl,
		AuthMiddleware: authMw,
	}
}

func (r *PersonRouter) RegisterPersonRoutes() {
	personCtrl := r.Controller.Person

	// Scoped under /families/:slug — so all person data is always tied to a family
	families := r.App.Group("/families/:slug")
	families.Use(r.AuthMiddleware.OptionalAuth())

	// Persons
	persons := families.Group("/persons")
	persons.Get("/", personCtrl.ListPersons)
	persons.Post("/", personCtrl.CreatePerson)
	persons.Get("/:person_id", personCtrl.GetPerson)
	persons.Patch("/:person_id", personCtrl.UpdatePerson)
	persons.Delete("/:person_id", personCtrl.DeletePerson)

	// Traversal — under a specific person
	persons.Get("/:person_id/ancestors", personCtrl.GetAncestors)
	persons.Get("/:person_id/descendants", personCtrl.GetDescendants)

	// Relationships (scoped per family, not per person)
	rels := families.Group("/relationships")
	rels.Get("/", personCtrl.ListRelationships)
	rels.Post("/", personCtrl.CreateRelationship)
	rels.Delete("/:rel_id", personCtrl.DeleteRelationship)

	// Full tree
	families.Get("/tree", personCtrl.GetTree)
}
