package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"git.dev.siap.id/kukuhkkh/app-silsilah/app/database/schema"
	roleRepo "git.dev.siap.id/kukuhkkh/app-silsilah/app/module/user/repository"
	"git.dev.siap.id/kukuhkkh/app-silsilah/utils/config"
)

// LogtoRoleSyncService syncs Logto roles to the local database after login.
type LogtoRoleSyncService interface {
	SyncRolesFromClaims(userID uint64, roleCodes []string) error
	SyncRolesFromManagementAPI(userID uint64, logtoUserID string, accessToken string) error
}

type logtoRoleSyncService struct {
	roleRepo     roleRepo.RoleRepository
	userRoleRepo roleRepo.UserRoleRepository
	cfg          *config.Config
	httpClient   *http.Client
}

type logtoRole struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewLogtoRoleSyncService(
	roleRepo roleRepo.RoleRepository,
	userRoleRepo roleRepo.UserRoleRepository,
	cfg *config.Config,
) LogtoRoleSyncService {
	return &logtoRoleSyncService{
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
		cfg:          cfg,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

// SyncRolesFromClaims is the preferred approach.
func (s *logtoRoleSyncService) SyncRolesFromClaims(userID uint64, roleCodes []string) error {
	if len(roleCodes) == 0 {
		return s.userRoleRepo.SyncUserRoles(userID, []schema.Role{})
	}

	for _, code := range roleCodes {
		_, _ = s.roleRepo.CreateIfNotExists(&schema.Role{
			Name: code,
			Code: code,
		})
	}

	roles, err := s.roleRepo.FindByCodes(roleCodes)
	if err != nil {
		return fmt.Errorf("role sync: find by codes: %w", err)
	}

	return s.userRoleRepo.SyncUserRoles(userID, roles)
}

func (s *logtoRoleSyncService) SyncRolesFromManagementAPI(userID uint64, logtoUserID string, accessToken string) error {
	endpoint := strings.TrimSuffix(s.cfg.Sso.Logto.Endpoint, "/")
	url := fmt.Sprintf("%s/api/users/%s/roles", endpoint, logtoUserID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("role sync: build request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("role sync: http request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("role sync: logto returned %d: %s", resp.StatusCode, string(body))
	}

	var logtoRoles []logtoRole
	if err := json.NewDecoder(resp.Body).Decode(&logtoRoles); err != nil {
		return fmt.Errorf("role sync: decode response: %w", err)
	}

	roleCodes := make([]string, 0, len(logtoRoles))
	for _, r := range logtoRoles {
		roleCodes = append(roleCodes, r.Name)
	}

	return s.SyncRolesFromClaims(userID, roleCodes)
}
