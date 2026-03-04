package service

import (
	"database/sql"
	"errors"
	"time"

	"git.dev.siap.id/kukuhkkh/app-silsilah/app/database/schema"
	family_repo "git.dev.siap.id/kukuhkkh/app-silsilah/app/module/family/repository"
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/module/person/repository"
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/module/person/request"
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/module/person/response"
)

type PersonService interface {
	CreatePerson(familySlug string, userID uint64, req request.CreatePersonRequest) (response.PersonResponse, error)
	GetPerson(familySlug string, userID uint64, personID uint64) (response.PersonResponse, error)
	ListPersons(familySlug string, userID uint64) ([]response.PersonResponse, error)
	UpdatePerson(familySlug string, userID uint64, personID uint64, req request.UpdatePersonRequest) (response.PersonResponse, error)
	DeletePerson(familySlug string, userID uint64, personID uint64) error

	CreateRelationship(familySlug string, userID uint64, req request.CreateRelationshipRequest) (response.RelationshipResponse, error)
	ListRelationships(familySlug string, userID uint64) ([]response.RelationshipResponse, error)
	DeleteRelationship(familySlug string, userID uint64, relID uint64) error

	GetTree(familySlug string, userID uint64) (response.TreeResponse, error)
	GetAncestors(familySlug string, userID uint64, personID uint64) ([]response.PersonResponse, error)
	GetDescendants(familySlug string, userID uint64, personID uint64) ([]response.PersonResponse, error)
}

type personService struct {
	personRepo       repository.PersonRepository
	relationshipRepo repository.RelationshipRepository
	familyRepo       family_repo.FamilyRepository
	familyMemberRepo family_repo.FamilyMemberRepository
}

func NewPersonService(
	personRepo repository.PersonRepository,
	relationshipRepo repository.RelationshipRepository,
	familyRepo family_repo.FamilyRepository,
	familyMemberRepo family_repo.FamilyMemberRepository,
) PersonService {
	return &personService{
		personRepo:       personRepo,
		relationshipRepo: relationshipRepo,
		familyRepo:       familyRepo,
		familyMemberRepo: familyMemberRepo,
	}
}

// resolveFamilyAccess resolves the family by slug and returns the family + user's role.
// Returns error if family not found or user cannot access a private family.
func (s *personService) resolveFamilyAccess(familySlug string, userID uint64) (*schema.Family, string, error) {
	family, err := s.familyRepo.FindBySlug(familySlug)
	if err != nil {
		return nil, "", errors.New("family not found")
	}

	role, err := s.familyMemberRepo.CheckUserRoleInFamily(family.ID, userID)
	if err != nil {
		if family.Visibility == "private" {
			return nil, "", errors.New("unauthorized")
		}
		role = "none"
	}

	return family, role, nil
}

func canEdit(role string) bool {
	return role == schema.FamilyRoleOwner || role == schema.FamilyRoleAdmin || role == schema.FamilyRoleEditor
}

func parseDate(raw string) sql.NullTime {
	if raw == "" {
		return sql.NullTime{}
	}

	for _, layout := range []string{time.DateOnly, time.RFC3339} {
		if t, err := time.Parse(layout, raw); err == nil {
			return sql.NullTime{Time: t, Valid: true}
		}
	}

	return sql.NullTime{}
}

// ─── Person CRUD ──────────────────────────────────────────────────────────────

func (s *personService) CreatePerson(familySlug string, userID uint64, req request.CreatePersonRequest) (response.PersonResponse, error) {
	family, role, err := s.resolveFamilyAccess(familySlug, userID)
	if err != nil {
		return response.PersonResponse{}, err
	}

	if !canEdit(role) {
		return response.PersonResponse{}, errors.New("insufficient permission to add person")
	}

	isAlive := true
	if req.IsAlive != nil {
		isAlive = *req.IsAlive
	}

	person := &schema.Person{
		FamilyID:   family.ID,
		FullName:   req.FullName,
		Nickname:   req.Nickname,
		Gender:     req.Gender,
		BirthDate:  parseDate(req.BirthDate),
		BirthPlace: req.BirthPlace,
		IsAlive:    isAlive,
		DeathDate:  parseDate(req.DeathDate),
		DeathPlace: req.DeathPlace,
		PhotoURL:   req.PhotoURL,
		Bio:        req.Bio,
		UserID:     req.UserID,
	}

	created, err := s.personRepo.Create(person)
	if err != nil {
		return response.PersonResponse{}, err
	}

	return response.FromPersonSchema(created), nil
}

func (s *personService) GetPerson(familySlug string, userID uint64, personID uint64) (response.PersonResponse, error) {
	family, _, err := s.resolveFamilyAccess(familySlug, userID)
	if err != nil {
		return response.PersonResponse{}, err
	}

	person, err := s.personRepo.FindByID(personID, family.ID)
	if err != nil {
		return response.PersonResponse{}, errors.New("person not found")
	}

	return response.FromPersonSchema(person), nil
}

func (s *personService) ListPersons(familySlug string, userID uint64) ([]response.PersonResponse, error) {
	family, _, err := s.resolveFamilyAccess(familySlug, userID)
	if err != nil {
		return nil, err
	}

	persons, err := s.personRepo.ListByFamilyID(family.ID)
	if err != nil {
		return nil, err
	}

	results := make([]response.PersonResponse, 0, len(persons))
	for i := range persons {
		results = append(results, response.FromPersonSchema(&persons[i]))
	}

	return results, nil
}

func (s *personService) UpdatePerson(familySlug string, userID uint64, personID uint64, req request.UpdatePersonRequest) (response.PersonResponse, error) {
	family, role, err := s.resolveFamilyAccess(familySlug, userID)
	if err != nil {
		return response.PersonResponse{}, err
	}

	if !canEdit(role) {
		return response.PersonResponse{}, errors.New("insufficient permission to update person")
	}

	person, err := s.personRepo.FindByID(personID, family.ID)
	if err != nil {
		return response.PersonResponse{}, errors.New("person not found")
	}

	if req.FullName != "" {
		person.FullName = req.FullName
	}

	if req.Nickname != "" {
		person.Nickname = req.Nickname
	}

	if req.Gender != "" {
		person.Gender = req.Gender
	}
	if req.BirthDate != "" {
		person.BirthDate = parseDate(req.BirthDate)
	}

	if req.BirthPlace != "" {
		person.BirthPlace = req.BirthPlace
	}

	if req.IsAlive != nil {
		person.IsAlive = *req.IsAlive
	}

	if req.DeathDate != "" {
		person.DeathDate = parseDate(req.DeathDate)
	}

	if req.DeathPlace != "" {
		person.DeathPlace = req.DeathPlace
	}

	if req.PhotoURL != "" {
		person.PhotoURL = req.PhotoURL
	}

	if req.Bio != "" {
		person.Bio = req.Bio
	}

	if req.UserID != nil {
		person.UserID = req.UserID
	}

	if err := s.personRepo.Update(person); err != nil {
		return response.PersonResponse{}, err
	}

	return response.FromPersonSchema(person), nil
}

func (s *personService) DeletePerson(familySlug string, userID uint64, personID uint64) error {
	family, role, err := s.resolveFamilyAccess(familySlug, userID)
	if err != nil {
		return err
	}

	if !canEdit(role) {
		return errors.New("insufficient permission to delete person")
	}

	if _, err := s.personRepo.FindByID(personID, family.ID); err != nil {
		return errors.New("person not found")
	}

	return s.personRepo.Delete(personID, family.ID)
}

// ─── Relationship CRUD ────────────────────────────────────────────────────────

func (s *personService) CreateRelationship(familySlug string, userID uint64, req request.CreateRelationshipRequest) (response.RelationshipResponse, error) {
	family, role, err := s.resolveFamilyAccess(familySlug, userID)
	if err != nil {
		return response.RelationshipResponse{}, err
	}

	if !canEdit(role) {
		return response.RelationshipResponse{}, errors.New("insufficient permission")
	}

	if req.PersonAID == req.PersonBID {
		return response.RelationshipResponse{}, errors.New("person_a_id and person_b_id cannot be the same")
	}

	// Validate both persons belong to this family
	if _, err := s.personRepo.FindByID(req.PersonAID, family.ID); err != nil {
		return response.RelationshipResponse{}, errors.New("person_a not found in this family")
	}

	if _, err := s.personRepo.FindByID(req.PersonBID, family.ID); err != nil {
		return response.RelationshipResponse{}, errors.New("person_b not found in this family")
	}

	if s.relationshipRepo.ExistsDuplicate(family.ID, req.PersonAID, req.PersonBID, req.RelationshipType) {
		return response.RelationshipResponse{}, errors.New("relationship already exists")
	}

	metadata := req.Metadata
	if len(metadata) == 0 {
		metadata = []byte("{}")
	}

	rel := &schema.Relationship{
		FamilyID:         family.ID,
		PersonAID:        req.PersonAID,
		PersonBID:        req.PersonBID,
		RelationshipType: req.RelationshipType,
		Metadata:         metadata,
	}

	created, err := s.relationshipRepo.Create(rel)
	if err != nil {
		return response.RelationshipResponse{}, err
	}

	// Auto-create parent_child for spouse:
	// If we just created parent_child(A→child), find if A has a spouse B,
	// then auto-create parent_child(B→child) if it doesn't already exist.
	if req.RelationshipType == schema.RelationshipParentChild {
		parentID := req.PersonAID
		childID := req.PersonBID

		rels, _ := s.relationshipRepo.ListByPersonID(parentID, family.ID)
		for _, r := range rels {
			if r.RelationshipType != schema.RelationshipSpouse {
				continue
			}

			spouseID := r.PersonBID
			if r.PersonBID == parentID {
				spouseID = r.PersonAID
			}

			// Auto-create if not duplicate
			if !s.relationshipRepo.ExistsDuplicate(family.ID, spouseID, childID, schema.RelationshipParentChild) {
				spouseRel := &schema.Relationship{
					FamilyID:         family.ID,
					PersonAID:        spouseID,
					PersonBID:        childID,
					RelationshipType: schema.RelationshipParentChild,
					Metadata:         []byte("{}"),
				}
				_, _ = s.relationshipRepo.Create(spouseRel)
			}
		}
	}

	return response.FromRelationshipSchema(created), nil
}

func (s *personService) ListRelationships(familySlug string, userID uint64) ([]response.RelationshipResponse, error) {
	family, _, err := s.resolveFamilyAccess(familySlug, userID)
	if err != nil {
		return nil, err
	}

	rels, err := s.relationshipRepo.ListByFamilyID(family.ID)
	if err != nil {
		return nil, err
	}

	results := make([]response.RelationshipResponse, 0, len(rels))
	for i := range rels {
		results = append(results, response.FromRelationshipSchema(&rels[i]))
	}

	return results, nil
}

func (s *personService) DeleteRelationship(familySlug string, userID uint64, relID uint64) error {
	family, role, err := s.resolveFamilyAccess(familySlug, userID)
	if err != nil {
		return err
	}

	if !canEdit(role) {
		return errors.New("insufficient permission")
	}

	if _, err := s.relationshipRepo.FindByID(relID, family.ID); err != nil {
		return errors.New("relationship not found")
	}

	return s.relationshipRepo.Delete(relID, family.ID)
}

// ─── Tree & Traversal ─────────────────────────────────────────────────────────

// GetTree returns the full family tree as nodes + edges for frontend rendering.
func (s *personService) GetTree(familySlug string, userID uint64) (response.TreeResponse, error) {
	family, _, err := s.resolveFamilyAccess(familySlug, userID)
	if err != nil {
		return response.TreeResponse{}, err
	}

	persons, err := s.personRepo.ListByFamilyID(family.ID)
	if err != nil {
		return response.TreeResponse{}, err
	}

	rels, err := s.relationshipRepo.ListByFamilyID(family.ID)
	if err != nil {
		return response.TreeResponse{}, err
	}

	nodes := make([]response.TreeNode, 0, len(persons))
	for _, p := range persons {
		node := response.TreeNode{
			ID:         p.ID,
			FullName:   p.FullName,
			Nickname:   p.Nickname,
			Gender:     p.Gender,
			BirthPlace: p.BirthPlace,
			IsAlive:    p.IsAlive,
			DeathPlace: p.DeathPlace,
			PhotoURL:   p.PhotoURL,
			Bio:        p.Bio,
			UserID:     p.UserID,
		}

		if p.BirthDate.Valid {
			node.BirthDate = &p.BirthDate.Time
		}

		if p.DeathDate.Valid {
			node.DeathDate = &p.DeathDate.Time
		}

		nodes = append(nodes, node)
	}

	edges := make([]response.TreeEdge, 0, len(rels))
	for _, r := range rels {
		edges = append(edges, response.TreeEdge{
			From:             r.PersonAID,
			To:               r.PersonBID,
			RelationshipType: r.RelationshipType,
			Metadata:         r.Metadata,
		})
	}

	// Calculate the Root Node (Patriarch/Matriarch)
	// A person is a root if they have no parents in this family.
	// We prioritize those who HAVE children (patriarchs/matriarchs).
	hasParent := make(map[uint64]bool)
	hasChildren := make(map[uint64]bool)
	for _, r := range rels {
		if r.RelationshipType == "parent_child" {
			hasParent[r.PersonBID] = true
			hasChildren[r.PersonAID] = true
		}
	}

	var rootID uint64
	// Strategy: find someone who has children but no parent.
	// If multiple, pick the oldest (earliest birth date) or the one with most descendants.
	for _, p := range persons {
		if !hasParent[p.ID] && hasChildren[p.ID] {
			if rootID == 0 {
				rootID = p.ID
			} else {
				// Potential tie-breaker: older person stays as root
				// For now, first found is usually the oldest due to database order/creation
			}
		}
	}

	// Fallback: if no clear patriarch, pick the first person
	if rootID == 0 && len(persons) > 0 {
		rootID = persons[0].ID
	}

	return response.TreeResponse{
		Nodes:  nodes,
		Edges:  edges,
		RootID: rootID,
	}, nil
}

func (s *personService) GetAncestors(familySlug string, userID uint64, personID uint64) ([]response.PersonResponse, error) {
	family, _, err := s.resolveFamilyAccess(familySlug, userID)
	if err != nil {
		return nil, err
	}

	persons, err := s.personRepo.GetAncestors(personID, family.ID)
	if err != nil {
		return nil, err
	}

	results := make([]response.PersonResponse, 0, len(persons))
	for i := range persons {
		results = append(results, response.FromPersonSchema(&persons[i]))
	}

	return results, nil
}

func (s *personService) GetDescendants(familySlug string, userID uint64, personID uint64) ([]response.PersonResponse, error) {
	family, _, err := s.resolveFamilyAccess(familySlug, userID)
	if err != nil {
		return nil, err
	}

	persons, err := s.personRepo.GetDescendants(personID, family.ID)
	if err != nil {
		return nil, err
	}

	results := make([]response.PersonResponse, 0, len(persons))
	for i := range persons {
		results = append(results, response.FromPersonSchema(&persons[i]))
	}

	return results, nil
}
