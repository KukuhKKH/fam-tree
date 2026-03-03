package controller

import (
	"strconv"

	"git.dev.siap.id/kukuhkkh/app-silsilah/app/middleware"
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/module/person/request"
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/module/person/service"
	"git.dev.siap.id/kukuhkkh/app-silsilah/utils/response"
	"github.com/gofiber/fiber/v2"
)

type PersonController interface {
	CreatePerson(c *fiber.Ctx) error
	GetPerson(c *fiber.Ctx) error
	ListPersons(c *fiber.Ctx) error
	UpdatePerson(c *fiber.Ctx) error
	DeletePerson(c *fiber.Ctx) error

	CreateRelationship(c *fiber.Ctx) error
	ListRelationships(c *fiber.Ctx) error
	DeleteRelationship(c *fiber.Ctx) error

	GetTree(c *fiber.Ctx) error
	GetAncestors(c *fiber.Ctx) error
	GetDescendants(c *fiber.Ctx) error
}

type personController struct {
	svc service.PersonService
}

func NewPersonController(svc service.PersonService) PersonController {
	return &personController{svc: svc}
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func parsePersonID(c *fiber.Ctx) (uint64, error) {
	return strconv.ParseUint(c.Params("person_id"), 10, 64)
}

// ─── Person endpoints ─────────────────────────────────────────────────────────

func (ctrl *personController) CreatePerson(c *fiber.Ctx) error {
	var req request.CreatePersonRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusBadRequest, Messages: response.Messages{err.Error()}})
	}

	res, err := ctrl.svc.CreatePerson(c.Params("slug"), middleware.GetUserID(c), req)
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusUnprocessableEntity, Messages: response.Messages{err.Error()}})
	}

	return response.Resp(c, response.Response{Code: fiber.StatusCreated, Data: res, Messages: response.Messages{"Person created"}})
}

func (ctrl *personController) ListPersons(c *fiber.Ctx) error {
	res, err := ctrl.svc.ListPersons(c.Params("slug"), middleware.GetUserID(c))
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusForbidden, Messages: response.Messages{err.Error()}})
	}

	return response.Resp(c, response.Response{Code: fiber.StatusOK, Data: res, Messages: response.Messages{"Persons retrieved"}})
}

func (ctrl *personController) GetPerson(c *fiber.Ctx) error {
	personID, err := parsePersonID(c)
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusBadRequest, Messages: response.Messages{"Invalid person ID"}})
	}

	res, err := ctrl.svc.GetPerson(c.Params("slug"), middleware.GetUserID(c), personID)
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusNotFound, Messages: response.Messages{err.Error()}})
	}

	return response.Resp(c, response.Response{Code: fiber.StatusOK, Data: res, Messages: response.Messages{"Person retrieved"}})
}

func (ctrl *personController) UpdatePerson(c *fiber.Ctx) error {
	var req request.UpdatePersonRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusBadRequest, Messages: response.Messages{err.Error()}})
	}

	personID, err := parsePersonID(c)
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusBadRequest, Messages: response.Messages{"Invalid person ID"}})
	}

	res, err := ctrl.svc.UpdatePerson(c.Params("slug"), middleware.GetUserID(c), personID, req)
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusUnprocessableEntity, Messages: response.Messages{err.Error()}})
	}

	return response.Resp(c, response.Response{Code: fiber.StatusOK, Data: res, Messages: response.Messages{"Person updated"}})
}

func (ctrl *personController) DeletePerson(c *fiber.Ctx) error {
	personID, err := parsePersonID(c)
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusBadRequest, Messages: response.Messages{"Invalid person ID"}})
	}

	if err := ctrl.svc.DeletePerson(c.Params("slug"), middleware.GetUserID(c), personID); err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusForbidden, Messages: response.Messages{err.Error()}})
	}

	return response.Resp(c, response.Response{Code: fiber.StatusOK, Messages: response.Messages{"Person deleted"}})
}

// ─── Relationship endpoints ───────────────────────────────────────────────────

func (ctrl *personController) CreateRelationship(c *fiber.Ctx) error {
	var req request.CreateRelationshipRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusBadRequest, Messages: response.Messages{err.Error()}})
	}

	res, err := ctrl.svc.CreateRelationship(c.Params("slug"), middleware.GetUserID(c), req)
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusUnprocessableEntity, Messages: response.Messages{err.Error()}})
	}

	return response.Resp(c, response.Response{Code: fiber.StatusCreated, Data: res, Messages: response.Messages{"Relationship created"}})
}

func (ctrl *personController) ListRelationships(c *fiber.Ctx) error {
	res, err := ctrl.svc.ListRelationships(c.Params("slug"), middleware.GetUserID(c))
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusForbidden, Messages: response.Messages{err.Error()}})
	}

	return response.Resp(c, response.Response{Code: fiber.StatusOK, Data: res, Messages: response.Messages{"Relationships retrieved"}})
}

func (ctrl *personController) DeleteRelationship(c *fiber.Ctx) error {
	relID, err := strconv.ParseUint(c.Params("rel_id"), 10, 64)
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusBadRequest, Messages: response.Messages{"Invalid relationship ID"}})
	}

	if err := ctrl.svc.DeleteRelationship(c.Params("slug"), middleware.GetUserID(c), relID); err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusForbidden, Messages: response.Messages{err.Error()}})
	}

	return response.Resp(c, response.Response{Code: fiber.StatusOK, Messages: response.Messages{"Relationship deleted"}})
}

// ─── Tree endpoints ───────────────────────────────────────────────────────────

func (ctrl *personController) GetTree(c *fiber.Ctx) error {
	res, err := ctrl.svc.GetTree(c.Params("slug"), middleware.GetUserID(c))
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusForbidden, Messages: response.Messages{err.Error()}})
	}

	return response.Resp(c, response.Response{Code: fiber.StatusOK, Data: res, Messages: response.Messages{"Tree retrieved"}})
}

func (ctrl *personController) GetAncestors(c *fiber.Ctx) error {
	personID, err := parsePersonID(c)
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusBadRequest, Messages: response.Messages{"Invalid person ID"}})
	}

	res, err := ctrl.svc.GetAncestors(c.Params("slug"), middleware.GetUserID(c), personID)
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusForbidden, Messages: response.Messages{err.Error()}})
	}

	return response.Resp(c, response.Response{Code: fiber.StatusOK, Data: res, Messages: response.Messages{"Ancestors retrieved"}})
}

func (ctrl *personController) GetDescendants(c *fiber.Ctx) error {
	personID, err := parsePersonID(c)
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusBadRequest, Messages: response.Messages{"Invalid person ID"}})
	}

	res, err := ctrl.svc.GetDescendants(c.Params("slug"), middleware.GetUserID(c), personID)
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusForbidden, Messages: response.Messages{err.Error()}})
	}

	return response.Resp(c, response.Response{Code: fiber.StatusOK, Data: res, Messages: response.Messages{"Descendants retrieved"}})
}
