package domain

import (
	"time"
)

type TeamID uint

const TeamIDCommon TeamID = 0

// Role represents a user's role in a team.
type Role string

const (
	RoleOwner  Role = "owner"
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

// Team represents a team in the system.
type Team struct {
	ID        TeamID
	Name      string
	CreatedAt time.Time
	Members   []TeamMember
}

type UserTeamInfo struct {
	ID       TeamID
	Name     string
	Role     Role
	CanLeave bool
}

// UserWithTeamsInfo represents a user with information about all their team memberships.
type UserWithTeamsInfo struct {
	User      User
	TeamsInfo []UserTeamInfo
}

// TeamMember represents a user's membership in a team.
type TeamMember struct {
	TeamID TeamID
	UserID UserID
	Role   Role
}

type TeamDTO struct {
	Name string
}
