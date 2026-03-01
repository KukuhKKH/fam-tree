package request

type CreateFamilyRequest struct {
	Name        string `json:"name" validate:"required,min=3"`
	Description string `json:"description"`
	Visibility  string `json:"visibility" validate:"required,oneof=private public link_only"`
}

type UpdateFamilyRequest struct {
	Name        string `json:"name" validate:"omitempty,min=3"`
	Description string `json:"description"`
	Visibility  string `json:"visibility" validate:"omitempty,oneof=private public link_only"`
}

type InviteMemberRequest struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=admin editor viewer"`
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=owner admin editor viewer"`
}
