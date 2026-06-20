// Enums
export type ProductCategory = 'general' | 'food' | 'beverage' | 'cosmetics' | 'fashion' | 'electronics' | 'home'

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
  tourCompleted: boolean
}

// Project
export interface Project {
  id: string
  name: string
  description?: string
  productCategory: string
  visualStyle: string
  createdAt: string
}

// Scene
export interface Scene {
  id: string
  projectId: string
  orderIndex: number
  prompt: string
  generated_prompt: string
  status: 'processing' | 'completed' | 'failed'
  fileUrl?: string
  errorCode?: string
  createdAt: string
}
