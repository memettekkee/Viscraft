package request

type CreateProjectRequest struct {
	BaseRequest
	Name            string `json:"name" binding:"required,min=1,max=255"`
	Description     string `json:"description"`
	ProductCategory string `json:"productCategory"`
	VisualStyle     string `json:"visualStyle"`
}

type GetProjectRequest struct {
	BaseRequest
	Id string `json:"id" binding:"required"`
}

type ListProjectsRequest struct {
	BaseRequest
}

type DeleteProjectRequest struct {
	BaseRequest
	Id string `json:"id" binding:"required"`
}
