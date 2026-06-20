export interface GenerateSceneRequest {
  projectId: string
  prompt: string
  generatedPrompt?: string
  referenceSceneId?: string
  uploadedReferenceImage?: string
}

export interface GetSceneRequest {
  id: string
}

export interface ListScenesRequest {
  projectId: string
}

export interface DeleteSceneRequest {
  id: string
}
