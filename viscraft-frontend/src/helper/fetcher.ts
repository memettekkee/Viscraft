import { api } from '../lib/api'

export async function postFetcher<T>([endpoint, ...params]: [string, ...unknown[]]): Promise<T> {
  const body = params.length === 1 && typeof params[0] === 'object' ? params[0] : {}
  const response = await api.post<T>(endpoint, body)
  return response.data
}
