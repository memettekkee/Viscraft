package response

type SceneData struct {
	Id         string `json:"id"`
	ProjectId  string `json:"projectId,omitempty"`
	OrderIndex int    `json:"orderIndex"`
	Prompt     string `json:"prompt"`
	GeneratedPrompt string `json:"generated_prompt"`
	Status     string `json:"status"`
	FileUrl    string `json:"fileUrl,omitempty"`
	ErrorCode  string `json:"errorCode,omitempty"`
	CreatedAt  string `json:"createdAt"`
}

type GenerateSceneResponse struct {
	BaseResponse
	Data *SceneData `json:"data,omitempty"`
}

type GetSceneResponse struct {
	BaseResponse
	Data *SceneData `json:"data,omitempty"`
}

type ListScenesResponse struct {
	BaseResponse
	Data []SceneData `json:"data,omitempty"`
}

type DeleteSceneResponse struct {
	BaseResponse
}
