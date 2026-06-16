package response

// ImageData represents the public-facing fields of an image record.
type ImageData struct {
	Id        string `json:"id"`
	Status    string `json:"status"`
	FileUrl   string `json:"fileUrl,omitempty"`
	Prompt    string `json:"prompt"`
	Genre     string `json:"genre"`
	AssetType string `json:"assetType"`
	Mood      string `json:"mood"`
	ErrorCode string `json:"errorCode,omitempty"`
	CreatedAt string `json:"createdAt"`
}

// GenerateImageResponse is returned with HTTP 202 when image generation starts.
type GenerateImageResponse struct {
	BaseResponse
	Data *ImageData `json:"data,omitempty"`
}

// GetImageResponse is returned when fetching a single image by ID.
type GetImageResponse struct {
	BaseResponse
	Data *ImageData `json:"data,omitempty"`
}

// ListImagesResponse is returned when listing all images for a project.
type ListImagesResponse struct {
	BaseResponse
	Data []ImageData `json:"data,omitempty"`
}

// DeleteImageResponse is returned after successfully deleting an image.
type DeleteImageResponse struct {
	BaseResponse
}
