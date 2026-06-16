import api from '../../lib/api'
import type { ApiResponse, Image } from '../../types'
import type { GenerateImageRequest, GetImageRequest, ListImagesRequest, DeleteImageRequest } from './types'

export type { GenerateImageRequest, GetImageRequest, ListImagesRequest, DeleteImageRequest }

export async function generateImage(payload: GenerateImageRequest): Promise<ApiResponse<Image>> {
  const response = await api.post<ApiResponse<Image>>('/images/generate', payload)
  return response.data
}

export async function getImage(payload: GetImageRequest): Promise<ApiResponse<Image>> {
  const response = await api.post<ApiResponse<Image>>('/images/get', payload)
  return response.data
}

export async function listImages(payload: ListImagesRequest): Promise<ApiResponse<Image[]>> {
  const response = await api.post<ApiResponse<Image[]>>('/images/list', payload)
  return response.data
}

export async function deleteImage(payload: DeleteImageRequest): Promise<ApiResponse<void>> {
  const response = await api.post<ApiResponse<void>>('/images/delete', payload)
  return response.data
}
