package controller

type Controller struct {
	Family FamilyController
}

func NewController(family FamilyController) *Controller {
	return &Controller{
		Family: family,
	}
}
