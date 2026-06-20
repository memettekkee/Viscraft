import api from '../../lib/api'
import type { ApiResponse, User } from '../../types'
import type { CreateUserRequest, LoginRequest } from './types'

export type { CreateUserRequest, LoginRequest }

export async function createUser(payload: CreateUserRequest): Promise<ApiResponse<User>> {
  const response = await api.post<ApiResponse<User>>('/users/create', payload)
  return response.data
}

export async function login(payload: LoginRequest): Promise<ApiResponse<User>> {
  const response = await api.post<ApiResponse<User>>('/users/login', payload)
  return response.data
}

export async function getCurrentUser(): Promise<ApiResponse<User>> {
  const response = await api.post<ApiResponse<User>>('/users/get', {})
  return response.data
}

export async function deleteUser(): Promise<ApiResponse<void>> {
  const response = await api.post<ApiResponse<void>>('/users/delete', {})
  return response.data
}

export async function completeTour(): Promise<ApiResponse<void>> {
  const response = await api.post<ApiResponse<void>>('/users/complete-tour', {})
  return response.data
}
