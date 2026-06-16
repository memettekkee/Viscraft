package request

// GenerateImageRequest holds the fields needed to generate a new concept art image.
// All fields are required and validated by the handler layer.
type GenerateImageRequest struct {
	BaseRequest
	ProjectId string `json:"projectId" binding:"required"`
	Prompt    string `json:"prompt" binding:"required"`
	Genre     string `json:"genre" binding:"required"`
	AssetType string `json:"assetType" binding:"required"`
	Mood      string `json:"mood" binding:"required"`
}

// GetImageRequest retrieves a single image by its ID.
type GetImageRequest struct {
	BaseRequest
	Id string `json:"id" binding:"required"`
}

// ListImagesRequest retrieves all images belonging to a project.
type ListImagesRequest struct {
	BaseRequest
	ProjectId string `json:"projectId" binding:"required"`
}

// DeleteImageRequest removes a single image by its ID.
type DeleteImageRequest struct {
	BaseRequest
	Id string `json:"id" binding:"required"`
}
