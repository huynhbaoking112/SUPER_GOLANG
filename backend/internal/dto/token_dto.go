package dto

type UserTokenData struct {
	GlobalRole           string                         `json:"global_role"`
	WorkspaceMemberships []WorkspaceMembershipTokenData `json:"workspace_memberships,omitempty"`
}

type WorkspaceMembershipTokenData struct {
	WorkspaceID string   `json:"workspace_id"`
	RoleName    string   `json:"role_name"`
	Permissions []string `json:"permissions"`
	Status      string   `json:"status"`
}
