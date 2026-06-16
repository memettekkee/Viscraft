// Enums
export type Genre = 'fantasy' | 'sci-fi' | 'post-apocalyptic' | 'steampunk' | 'horror'
export type AssetType = 'character' | 'location' | 'item' | 'creature'
export type Mood = 'dark' | 'epic' | 'mysterious' | 'whimsical'

// Base response shape — every backend response includes these
export interface ApiResponse<T = unknown> {
  requestId: string
  success: boolean
  message: string
  data?: T
  token?: string
  errorCode?: string
}

// User
export interface User {
  id: string
  email: string
  name?: string
  createdAt: string
}

// Project
export interface Project {
  id: string
  name: string
  description?: string
  createdAt: string
}

// Image
export interface Image {
  id: string
  status: 'processing' | 'completed' | 'failed'
  fileUrl?: string
  prompt: string
  genre: Genre
  assetType: AssetType
  mood: Mood
  errorCode?: string
  createdAt: string
}
