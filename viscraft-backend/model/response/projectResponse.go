package response

type ProjectData struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	ProductCategory string `json:"productCategory"`
	VisualStyle     string `json:"visualStyle"`
	CreatedAt       string `json:"createdAt"`
}

type CreateProjectResponse struct {
	BaseResponse
	Data *ProjectData `json:"data,omitempty"`
}

type GetProjectResponse struct {
	BaseResponse
	Data *ProjectData `json:"data,omitempty"`
}

type ListProjectsResponse struct {
	BaseResponse
	Data []ProjectData `json:"data"`
}

type DeleteProjectResponse struct {
	BaseResponse
}
