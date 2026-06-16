import type { Genre, AssetType, Mood } from '../../types'

export interface GenerateImageRequest {
  projectId: string
  prompt: string
  genre: Genre
  assetType: AssetType
  mood: Mood
  referenceImage?: string
}

export interface GetImageRequest {
  id: string
}

export interface ListImagesRequest {
  projectId: string
}

export interface DeleteImageRequest {
  id: string
}
