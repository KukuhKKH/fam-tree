package controller

type Controller struct {
	Family FamilyController
	Public PublicController
}

func NewController(family FamilyController, public PublicController) *Controller {
	return &Controller{
		Family: family,
		Public: public,
	}
}
