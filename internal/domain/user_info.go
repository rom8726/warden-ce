package domain

// UserInfo represents detailed information about a user, включая их команды.
type UserInfo struct {
	User  User
	Teams []UserTeamInfo
}
