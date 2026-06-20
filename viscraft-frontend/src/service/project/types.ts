export interface CreateProjectRequest {
  name: string
  description?: string
  productCategory?: string
  visualStyle?: string
}

export interface GetProjectRequest {
  id: string
}

export interface DeleteProjectRequest {
  id: string
}
