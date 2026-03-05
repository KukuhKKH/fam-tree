package controller

import (
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/module/family/service"
	"git.dev.siap.id/kukuhkkh/app-silsilah/utils/response"
	"github.com/gofiber/fiber/v2"
)

type publicController struct {
	svc service.FamilyService
}

type PublicController interface {
	GetPublicFamily(c *fiber.Ctx) error
}

func NewPublicController(svc service.FamilyService) PublicController {
	return &publicController{svc: svc}
}

func (ctrl *publicController) GetPublicFamily(c *fiber.Ctx) error {
	slug := c.Params("slug")
	res, err := ctrl.svc.GetPublicFamilyDetails(slug)
	if err != nil {
		return response.Resp(c, response.Response{Code: fiber.StatusNotFound, Messages: response.Messages{err.Error()}})
	}

	return response.Resp(c, response.Response{Code: fiber.StatusOK, Data: res, Messages: response.Messages{"Public family retrieved"}})
}
