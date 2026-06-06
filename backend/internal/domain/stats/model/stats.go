package model

// ProjectDashboardStats holds aggregate metrics for a municipality dashboard.
type ProjectDashboardStats struct {
	TotalProjects          int            `json:"total_projects"`
	OngoingProjects        int            `json:"ongoing_projects"`
	CompletedProjects      int            `json:"completed_projects"`
	ProjectsByStatus       map[string]int `json:"projects_by_status"`
	PendingApprovals       int            `json:"pending_approvals"`
	MissingDocuments       int            `json:"missing_documents"`
	PendingDocumentReviews int            `json:"pending_document_reviews"`
	TotalBuildings         int            `json:"total_buildings"`
	TotalPropertyOwners    int            `json:"total_property_owners"`
}

// AdminDashboardStats holds aggregate metrics for the system admin dashboard.
type AdminDashboardStats struct {
	TotalOrganizations  int            `json:"total_organizations"`
	OrganizationsByStatus map[string]int `json:"organizations_by_status"`
	TotalUsers          int            `json:"total_users"`
	ActiveUsers         int            `json:"active_users"`
	UsersByStatus       map[string]int `json:"users_by_status"`
	RoleDistribution    map[string]int `json:"role_distribution"`
}
