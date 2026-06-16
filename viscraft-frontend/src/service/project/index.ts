import api from '../../lib/api'
import type { ApiResponse, Project } from '../../types'
import type { CreateProjectRequest, GetProjectRequest, DeleteProjectRequest } from './types'

export type { CreateProjectRequest, GetProjectRequest, DeleteProjectRequest }

export async function createProject(payload: CreateProjectRequest): Promise<ApiResponse<Project>> {
  const response = await api.post<ApiResponse<Project>>('/projects/create', payload)
  return response.data
}

export async function getProject(payload: GetProjectRequest): Promise<ApiResponse<Project>> {
  const response = await api.post<ApiResponse<Project>>('/projects/get', payload)
  return response.data
}

export async function listProjects(): Promise<ApiResponse<Project[]>> {
  const response = await api.post<ApiResponse<Project[]>>('/projects/list', {})
  return response.data
}

export async function deleteProject(payload: DeleteProjectRequest): Promise<ApiResponse<void>> {
  const response = await api.post<ApiResponse<void>>('/projects/delete', payload)
  return response.data
}
