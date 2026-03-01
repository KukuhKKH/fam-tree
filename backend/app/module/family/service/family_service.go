package service

import (
	"errors"
	"fmt"
	"strings"

	"git.dev.siap.id/kukuhkkh/app-silsilah/app/database/schema"
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/module/family/repository"
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/module/family/request"
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/module/family/response"
	user_repo "git.dev.siap.id/kukuhkkh/app-silsilah/app/module/user/repository"
)

type FamilyService interface {
	CreateFamily(userID uint64, req request.CreateFamilyRequest) (response.FamilyResponse, error)
	GetMyFamilies(userID uint64) ([]response.FamilyResponse, error)
	GetFamilyDetails(slug string, userID uint64) (response.FamilyResponse, error)
	UpdateFamily(slug string, userID uint64, req request.UpdateFamilyRequest) (response.FamilyResponse, error)
	DeleteFamily(slug string, userID uint64) error

	InviteMember(slug string, userID uint64, req request.InviteMemberRequest) (response.FamilyMemberResponse, error)
	GetMembers(slug string, userID uint64) ([]response.FamilyMemberResponse, error)
	UpdateMemberRole(slug string, userID, targetUserID uint64, req request.UpdateMemberRoleRequest) error
	RemoveMember(slug string, userID, targetUserID uint64) error
}

type familyService struct {
	familyRepo       repository.FamilyRepository
	familyMemberRepo repository.FamilyMemberRepository
	userRepo         user_repo.UserRepository
}

func NewFamilyService(
	familyRepo repository.FamilyRepository,
	familyMemberRepo repository.FamilyMemberRepository,
	userRepo user_repo.UserRepository,
) FamilyService {
	return &familyService{
		familyRepo:       familyRepo,
		familyMemberRepo: familyMemberRepo,
		userRepo:         userRepo,
	}
}

func (s *familyService) generateSlug(name string) string {
	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	baseSlug := slug
	counter := 1

	for s.familyRepo.CheckSlugExists(slug) {
		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}

	return slug
}

func (s *familyService) CreateFamily(userID uint64, req request.CreateFamilyRequest) (response.FamilyResponse, error) {
	family := &schema.Family{
		Name:        req.Name,
		Slug:        s.generateSlug(req.Name),
		Description: req.Description,
		Visibility:  req.Visibility,
		CreatedBy:   userID,
	}

	createdFamily, err := s.familyRepo.Create(family)
	if err != nil {
		return response.FamilyResponse{}, err
	}

	// Creator is automatically the owner
	_, err = s.familyMemberRepo.Create(&schema.FamilyMember{
		FamilyID: createdFamily.ID,
		UserID:   userID,
		Role:     schema.FamilyRoleOwner,
	})

	if err != nil {
		return response.FamilyResponse{}, err
	}

	return response.FromFamilySchema(createdFamily, schema.FamilyRoleOwner), nil
}

func (s *familyService) GetMyFamilies(userID uint64) ([]response.FamilyResponse, error) {
	members, err := s.familyRepo.GetFamiliesByUserID(userID)
	if err != nil {
		return nil, err
	}

	var results []response.FamilyResponse
	for _, m := range members {
		results = append(results, response.FromFamilySchema(&m.Family, m.Role))
	}

	return results, nil
}

func (s *familyService) GetFamilyDetails(slug string, userID uint64) (response.FamilyResponse, error) {
	family, err := s.familyRepo.FindBySlug(slug)
	if err != nil {
		return response.FamilyResponse{}, errors.New("family not found")
	}

	role, err := s.familyMemberRepo.CheckUserRoleInFamily(family.ID, userID)
	if err != nil {
		if family.Visibility == "private" {
			return response.FamilyResponse{}, errors.New("unauthorized to view this family")
		}

		role = "none" // Public viewer not in family
	}

	return response.FromFamilySchema(family, role), nil
}

func (s *familyService) UpdateFamily(slug string, userID uint64, req request.UpdateFamilyRequest) (response.FamilyResponse, error) {
	family, err := s.familyRepo.FindBySlug(slug)
	if err != nil {
		return response.FamilyResponse{}, errors.New("family not found")
	}

	role, err := s.familyMemberRepo.CheckUserRoleInFamily(family.ID, userID)
	if err != nil || (role != schema.FamilyRoleOwner && role != schema.FamilyRoleAdmin) {
		return response.FamilyResponse{}, errors.New("insufficient role to update family")
	}

	if req.Name != "" {
		family.Name = req.Name
	}

	if req.Description != "" {
		family.Description = req.Description
	}

	if req.Visibility != "" {
		family.Visibility = req.Visibility
	}

	if err := s.familyRepo.Update(family); err != nil {
		return response.FamilyResponse{}, err
	}

	return response.FromFamilySchema(family, role), nil
}

func (s *familyService) DeleteFamily(slug string, userID uint64) error {
	family, err := s.familyRepo.FindBySlug(slug)
	if err != nil {
		return errors.New("family not found")
	}

	role, err := s.familyMemberRepo.CheckUserRoleInFamily(family.ID, userID)
	if err != nil || role != schema.FamilyRoleOwner {
		return errors.New("only owner can delete family")
	}

	return s.familyRepo.Delete(family.ID)
}

func (s *familyService) InviteMember(slug string, userID uint64, req request.InviteMemberRequest) (response.FamilyMemberResponse, error) {
	family, err := s.familyRepo.FindBySlug(slug)
	if err != nil {
		return response.FamilyMemberResponse{}, errors.New("family not found")
	}

	role, err := s.familyMemberRepo.CheckUserRoleInFamily(family.ID, userID)
	if err != nil || (role != schema.FamilyRoleOwner && role != schema.FamilyRoleAdmin) {
		return response.FamilyMemberResponse{}, errors.New("insufficient permission to invite")
	}

	// Logto-based flow: User must exist in DB (meaning they have logged in at least once)
	targetUser, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil || targetUser == nil {
		return response.FamilyMemberResponse{}, errors.New("user not found. they must register first")
	}

	// Check if already member
	_, err = s.familyMemberRepo.FindByFamilyAndUser(family.ID, targetUser.ID)
	if err == nil {
		return response.FamilyMemberResponse{}, errors.New("user is already a member")
	}

	member := &schema.FamilyMember{
		FamilyID: family.ID,
		UserID:   targetUser.ID,
		Role:     req.Role,
	}

	created, err := s.familyMemberRepo.Create(member)
	if err != nil {
		return response.FamilyMemberResponse{}, err
	}
	created.User = *targetUser

	return response.FromFamilyMemberSchema(created), nil
}

func (s *familyService) GetMembers(slug string, userID uint64) ([]response.FamilyMemberResponse, error) {
	family, err := s.familyRepo.FindBySlug(slug)
	if err != nil {
		return nil, errors.New("family not found")
	}

	role, err := s.familyMemberRepo.CheckUserRoleInFamily(family.ID, userID)
	if err != nil && family.Visibility == "private" {
		return nil, errors.New("unauthorized")
	}
	_ = role // Not strictly used for filtering currently, but validates membership/publicity

	members, err := s.familyMemberRepo.GetByFamilyID(family.ID)
	if err != nil {
		return nil, err
	}

	var results []response.FamilyMemberResponse
	for _, m := range members {
		results = append(results, response.FromFamilyMemberSchema(&m))
	}
	return results, nil
}

func (s *familyService) UpdateMemberRole(slug string, userID, targetUserID uint64, req request.UpdateMemberRoleRequest) error {
	family, err := s.familyRepo.FindBySlug(slug)
	if err != nil {
		return errors.New("family not found")
	}

	myRole, err := s.familyMemberRepo.CheckUserRoleInFamily(family.ID, userID)
	if err != nil || (myRole != schema.FamilyRoleOwner && myRole != schema.FamilyRoleAdmin) {
		return errors.New("insufficient permission")
	}

	targetMember, err := s.familyMemberRepo.FindByFamilyAndUser(family.ID, targetUserID)
	if err != nil {
		return errors.New("target member not found")
	}

	// Only owner can update an owner, or make someone else an owner
	if (targetMember.Role == schema.FamilyRoleOwner || req.Role == schema.FamilyRoleOwner) && myRole != schema.FamilyRoleOwner {
		return errors.New("only owner can grant or modify owner roles")
	}

	return s.familyMemberRepo.UpdateRole(family.ID, targetUserID, req.Role)
}

func (s *familyService) RemoveMember(slug string, userID, targetUserID uint64) error {
	family, err := s.familyRepo.FindBySlug(slug)
	if err != nil {
		return errors.New("family not found")
	}

	myRole, err := s.familyMemberRepo.CheckUserRoleInFamily(family.ID, userID)
	if err != nil {
		return errors.New("not a member")
	}

	targetMember, err := s.familyMemberRepo.FindByFamilyAndUser(family.ID, targetUserID)
	if err != nil {
		return errors.New("target member not found")
	}

	// Rules for removal:
	// 1. Can remove self (leave family) - unless owner
	// 2. Owner can remove anyone
	// 3. Admin can remove editor/viewer

	if userID == targetUserID {
		if myRole == schema.FamilyRoleOwner {
			return errors.New("owner cannot leave family, must transfer ownership or delete family")
		}
	} else {
		if myRole != schema.FamilyRoleOwner && myRole != schema.FamilyRoleAdmin {
			return errors.New("insufficient permission")
		}
		if myRole == schema.FamilyRoleAdmin && (targetMember.Role == schema.FamilyRoleOwner || targetMember.Role == schema.FamilyRoleAdmin) {
			return errors.New("admin cannot remove owners or other admins")
		}
	}

	return s.familyMemberRepo.Delete(family.ID, targetUserID)
}
