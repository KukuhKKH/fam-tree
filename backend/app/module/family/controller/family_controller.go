package controller

import (
	"strconv"

	"git.dev.siap.id/kukuhkkh/app-silsilah/app/middleware"
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/module/family/request"
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/module/family/service"
	"git.dev.siap.id/kukuhkkh/app-silsilah/utils/response"
	"github.com/gofiber/fiber/v2"
)

type familyController struct {
	svc service.FamilyService
}

type FamilyController interface {
	Create(c *fiber.Ctx) error
	GetMyFamilies(c *fiber.Ctx) error
	GetFamilyDetails(c *fiber.Ctx) error
	UpdateFamily(c *fiber.Ctx) error
	DeleteFamily(c *fiber.Ctx) error
	InviteMember(c *fiber.Ctx) error
	GetMembers(c *fiber.Ctx) error
	UpdateMemberRole(c *fiber.Ctx) error
	RemoveMember(c *fiber.Ctx) error
}

func NewFamilyController(svc service.FamilyService) FamilyController {
	return &familyController{svc: svc}
}

func (ctrl *familyController) Create(c *fiber.Ctx) error {
	var req request.CreateFamilyRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusBadRequest, Messages: response.Messages{err.Error()}})
	}

	res, err := ctrl.svc.CreateFamily(middleware.GetUserID(c), req)
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusInternalServerError, Messages: response.Messages{err.Error()}})
	}

	return response.Resp(c, response.Response{Code: fiber.StatusCreated, Data: res, Messages: response.Messages{"Family created successfully"}})
}

func (ctrl *familyController) GetMyFamilies(c *fiber.Ctx) error {
	res, err := ctrl.svc.GetMyFamilies(middleware.GetUserID(c))
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusInternalServerError, Messages: response.Messages{err.Error()}})
	}

	return response.Resp(c, response.Response{Code: fiber.StatusOK, Data: res, Messages: response.Messages{"My families retrieved"}})
}

func (ctrl *familyController) GetFamilyDetails(c *fiber.Ctx) error {
	slug := c.Params("slug")
	res, err := ctrl.svc.GetFamilyDetails(slug, middleware.GetUserID(c))
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusForbidden, Messages: response.Messages{err.Error()}})
	}

	return response.Resp(c, response.Response{Code: fiber.StatusOK, Data: res, Messages: response.Messages{"Family details retrieved"}})
}

func (ctrl *familyController) UpdateFamily(c *fiber.Ctx) error {
	var req request.UpdateFamilyRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusBadRequest, Messages: response.Messages{err.Error()}})
	}

	slug := c.Params("slug")
	res, err := ctrl.svc.UpdateFamily(slug, middleware.GetUserID(c), req)
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusInternalServerError, Messages: response.Messages{err.Error()}})
	}

	return response.Resp(c, response.Response{Code: fiber.StatusOK, Data: res, Messages: response.Messages{"Family updated"}})
}

func (ctrl *familyController) DeleteFamily(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if err := ctrl.svc.DeleteFamily(slug, middleware.GetUserID(c)); err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusForbidden, Messages: response.Messages{err.Error()}})
	}

	return response.Resp(c, response.Response{Code: fiber.StatusOK, Messages: response.Messages{"Family deleted"}})
}

func (ctrl *familyController) InviteMember(c *fiber.Ctx) error {
	var req request.InviteMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusBadRequest, Messages: response.Messages{err.Error()}})
	}

	slug := c.Params("slug")
	res, err := ctrl.svc.InviteMember(slug, middleware.GetUserID(c), req)
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusBadRequest, Messages: response.Messages{err.Error()}})
	}

	return response.Resp(c, response.Response{Code: fiber.StatusCreated, Data: res, Messages: response.Messages{"Member invited successfully"}})
}

func (ctrl *familyController) GetMembers(c *fiber.Ctx) error {
	slug := c.Params("slug")
	res, err := ctrl.svc.GetMembers(slug, middleware.GetUserID(c))
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusForbidden, Messages: response.Messages{err.Error()}})
	}

	return response.Resp(c, response.Response{Code: fiber.StatusOK, Data: res, Messages: response.Messages{"Members retrieved"}})
}

func (ctrl *familyController) UpdateMemberRole(c *fiber.Ctx) error {
	var req request.UpdateMemberRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusBadRequest, Messages: response.Messages{err.Error()}})
	}

	slug := c.Params("slug")
	targetID, err := strconv.ParseUint(c.Params("user_id"), 10, 64)
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusBadRequest, Messages: response.Messages{"Invalid user ID"}})
	}

	err = ctrl.svc.UpdateMemberRole(slug, middleware.GetUserID(c), targetID, req)
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusForbidden, Messages: response.Messages{err.Error()}})
	}

	return response.Resp(c, response.Response{Code: fiber.StatusOK, Messages: response.Messages{"Member role updated"}})
}

func (ctrl *familyController) RemoveMember(c *fiber.Ctx) error {
	slug := c.Params("slug")
	targetID, err := strconv.ParseUint(c.Params("user_id"), 10, 64)
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusBadRequest, Messages: response.Messages{"Invalid user ID"}})
	}

	err = ctrl.svc.RemoveMember(slug, middleware.GetUserID(c), targetID)
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusForbidden, Messages: response.Messages{err.Error()}})
	}

	return response.Resp(c, response.Response{Code: fiber.StatusOK, Messages: response.Messages{"Member removed"}})
}
