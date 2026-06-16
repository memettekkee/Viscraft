import useSWR from 'swr'
import { postFetcher } from '../../../helper/fetcher'
import type { ApiResponse, Image } from '../../../types'

/**
 * SWR hook that polls image status while processing.
 * Polls every 3s when status is "processing", stops when completed/failed.
 * Disables fetching when imageId is null.
 *
 * Validates: Requirements 9.1, 9.2, 9.3, 9.4
 */
export function usePollImageStatus(imageId: string | null) {
  const { data, error, isLoading, mutate } = useSWR<ApiResponse<Image>>(
    imageId ? ['/images/get', { id: imageId }] : null,
    postFetcher,
    {
      refreshInterval: (data) =>
        data?.data?.status === 'processing' ? 3000 : 0,
      revalidateOnFocus: false,
    }
  )

  return {
    image: data?.data ?? null,
    isLoading,
    error,
    mutate,
  }
}
