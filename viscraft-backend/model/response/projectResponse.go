package response

// ProjectData holds the project fields returned to the client.
type ProjectData struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}

// CreateProjectResponse is the response for project creation.
type CreateProjectResponse struct {
	BaseResponse
	Data *ProjectData `json:"data,omitempty"`
}

// GetProjectResponse is the response for retrieving a single project.
type GetProjectResponse struct {
	BaseResponse
	Data *ProjectData `json:"data,omitempty"`
}

// ListProjectsResponse is the response for listing all projects.
type ListProjectsResponse struct {
	BaseResponse
	Data []ProjectData `json:"data"`
}

// DeleteProjectResponse is the response for project deletion.
type DeleteProjectResponse struct {
	BaseResponse
}
