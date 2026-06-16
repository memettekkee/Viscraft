export interface CreateProjectRequest {
  name: string
  description?: string
}

export interface GetProjectRequest {
  id: string
}

export interface DeleteProjectRequest {
  id: string
}
