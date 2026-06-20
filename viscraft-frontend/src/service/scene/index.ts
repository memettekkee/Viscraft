import api from '../../lib/api'
import type { ApiResponse, Scene } from '../../types'
import type { GenerateSceneRequest, GetSceneRequest, ListScenesRequest, DeleteSceneRequest } from './types'

export type { GenerateSceneRequest, GetSceneRequest, ListScenesRequest, DeleteSceneRequest }

export async function generateScene(payload: GenerateSceneRequest): Promise<ApiResponse<Scene>> {
  const response = await api.post<ApiResponse<Scene>>('/scenes/generate', payload)
  return response.data
}

export async function getScene(payload: GetSceneRequest): Promise<ApiResponse<Scene>> {
  const response = await api.post<ApiResponse<Scene>>('/scenes/get', payload)
  return response.data
}

export async function listScenes(payload: ListScenesRequest): Promise<ApiResponse<Scene[]>> {
  const response = await api.post<ApiResponse<Scene[]>>('/scenes/list', payload)
  return response.data
}

export async function deleteScene(payload: DeleteSceneRequest): Promise<ApiResponse<void>> {
  const response = await api.post<ApiResponse<void>>('/scenes/delete', payload)
  return response.data
}
