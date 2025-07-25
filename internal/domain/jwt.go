package domain

import (
	"github.com/golang-jwt/jwt"
)

type TokenType string

const (
	TokenTypeAccess        TokenType = "accessToken"
	TokenTypeRefresh       TokenType = "refreshToken"
	TokenTypeResetPassword TokenType = "resetPassword"
)

type UserPermissions struct {
	ProjectPermissions map[ProjectID]ProjectPermission `json:"project_permissions,omitempty"`
	TeamRoles          map[TeamID]Role                 `json:"team_roles,omitempty"`
	CanCreateProjects  bool                            `json:"can_create_projects"`
	CanCreateTeams     bool                            `json:"can_create_teams"`
	CanManageUsers     bool                            `json:"can_manage_users"`
}

type ProjectPermission struct {
	CanRead   bool `json:"can_read"`
	CanWrite  bool `json:"can_write"`
	CanDelete bool `json:"can_delete"`
	CanManage bool `json:"can_manage"`
	TeamRole  Role `json:"team_role,omitempty"`
}

type TokenClaims struct {
	jwt.StandardClaims
	TokenType   TokenType       `json:"type"`
	UserID      uint            `json:"userId"`
	Username    string          `json:"username"`
	IsSuperuser bool            `json:"isSuperuser"`
	Permissions UserPermissions `json:"permissions,omitempty"`
}
