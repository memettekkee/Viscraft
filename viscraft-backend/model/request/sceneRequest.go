package request

type GenerateSceneRequest struct {
	BaseRequest
	ProjectId              string `json:"projectId" binding:"required"`
	Prompt                 string `json:"prompt" binding:"required"`
	GeneratedPrompt        string `json:"generatedPrompt,omitempty"`
	ReferenceSceneId       string `json:"referenceSceneId,omitempty"`
	UploadedReferenceImage string `json:"uploadedReferenceImage,omitempty"`
}

type GetSceneRequest struct {
	BaseRequest
	Id string `json:"id" binding:"required"`
}

type ListScenesRequest struct {
	BaseRequest
	ProjectId string `json:"projectId" binding:"required"`
}

type DeleteSceneRequest struct {
	BaseRequest
	Id string `json:"id" binding:"required"`
}
