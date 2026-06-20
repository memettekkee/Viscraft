import useSWR from 'swr'
import { postFetcher } from '../../../helper/fetcher'
import type { ApiResponse, Scene } from '../../../types'

/**
 * SWR hook for fetching scenes belonging to a project.
 * Returns scenes sorted by orderIndex in ascending order.
 * Disables fetching when projectId is null/undefined.
 *
 * Validates: Requirements 5.1
 */
export function useSceneList(projectId: string | null | undefined) {
  const { data, error, isLoading, mutate } = useSWR<ApiResponse<Scene[]>>(
    projectId ? ['/scenes/list', { projectId }] : null,
    postFetcher
  )

  const scenes = [...(data?.data ?? [])].sort((a, b) => a.orderIndex - b.orderIndex)

  return {
    scenes,
    isLoading,
    error,
    mutate,
  }
}
