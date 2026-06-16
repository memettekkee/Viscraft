package request

// CreateProjectRequest is used to create a new project.
// Name is required (1-255 characters). Description is optional.
type CreateProjectRequest struct {
	BaseRequest
	Name        string `json:"name" binding:"required,min=1,max=255"`
	Description string `json:"description"`
}

// GetProjectRequest is used to retrieve a single project by ID.
type GetProjectRequest struct {
	BaseRequest
	Id string `json:"id" binding:"required"`
}

// ListProjectsRequest is used to list all projects for the authenticated user.
// The userId is extracted from the JWT context, not from the request body.
type ListProjectsRequest struct {
	BaseRequest
}

// DeleteProjectRequest is used to delete a project by ID.
type DeleteProjectRequest struct {
	BaseRequest
	Id string `json:"id" binding:"required"`
}
