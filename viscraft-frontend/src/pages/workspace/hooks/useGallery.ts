import useSWR from 'swr'
import { postFetcher } from '../../../helper/fetcher'
import type { ApiResponse, Image } from '../../../types'

/**
 * SWR hook for fetching images belonging to a project.
 * Disables fetching when projectId is null/undefined.
 *
 * Validates: Requirements 5.1, 13.3
 */
export function useGallery(projectId: string | null | undefined) {
  const { data, error, isLoading, mutate } = useSWR<ApiResponse<Image[]>>(
    projectId ? ['/images/list', { projectId }] : null,
    postFetcher
  )

  return {
    images: data?.data ?? [],
    isLoading,
    error,
    mutate,
  }
}
