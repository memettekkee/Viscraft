import { api } from '../lib/api'

/**
 * Generic SWR fetcher wrapping axios POST.
 * SWR keys are tuples where key[0] is the endpoint and key[1...] are optional payload fields.
 */
export async function postFetcher<T>([endpoint, ...params]: [string, ...unknown[]]): Promise<T> {
  const body = params.length === 1 && typeof params[0] === 'object' ? params[0] : {}
  const response = await api.post<T>(endpoint, body)
  return response.data
}
