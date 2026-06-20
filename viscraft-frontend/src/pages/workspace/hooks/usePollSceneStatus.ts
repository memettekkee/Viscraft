import useSWR from 'swr'
import { postFetcher } from '../../../helper/fetcher'
import type { ApiResponse, Scene } from '../../../types'

/**
 * SWR hook that polls scene status while processing.
 * Polls every 3s when status is "processing", stops when completed/failed.
 * Disables fetching when sceneId is null.
 *
 * Validates: Requirements 7.1, 7.2, 7.3
 */
export function usePollSceneStatus(sceneId: string | null) {
  const { data, error, isLoading, mutate } = useSWR<ApiResponse<Scene>>(
    sceneId ? ['/scenes/get', { id: sceneId }] : null,
    postFetcher,
    {
      refreshInterval: (data) =>
        data?.data?.status === 'processing' ? 3000 : 0,
      revalidateOnFocus: false,
    }
  )

  return {
    scene: data?.data ?? null,
    isLoading,
    error,
    mutate,
  }
}
