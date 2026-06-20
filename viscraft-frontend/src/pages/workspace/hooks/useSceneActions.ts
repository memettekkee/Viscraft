import { useSWRConfig } from 'swr'
import { deleteScene } from '../../../service/scene'
import { useWorkspaceStore } from '../../../store/workspaceStore'
import { showToast } from '../../../components/CustomToast'
import { ERROR_MESSAGES } from '../../../constants'
import type { AxiosError } from 'axios'
import type { Scene, ApiResponse } from '../../../types'

/**
 * Extracts scene card action logic: delete and regenerate (open modal with pre-filled prompt).
 */
export function useSceneActions() {
  const { mutate } = useSWRConfig()
  const activeProjectId = useWorkspaceStore((s) => s.activeProjectId)
  const openGenerateModal = useWorkspaceStore((s) => s.openGenerateModal)

  /**
   * Delete: removes a scene and mutates the SWR cache to refresh the list.
   */
  async function handleDelete(sceneId: string) {
    if (!activeProjectId) return

    try {
      await deleteScene({ id: sceneId })

      mutate(
        ['/scenes/list', { projectId: activeProjectId }],
        undefined,
        { revalidate: true }
      )

      showToast({ type: 'success', title: 'Ad shot deleted' })
    } catch (err: unknown) {
      const axiosError = err as AxiosError<ApiResponse>
      const code = axiosError.response?.data?.errorCode
      const backendMessage = axiosError.response?.data?.message
      const message = backendMessage ?? (code ? (ERROR_MESSAGES[code] ?? 'An error occurred') : ERROR_MESSAGES.NETWORK_ERROR)
      showToast({ type: 'error', title: message })
    }
  }

  /**
   * Regenerate: opens the GenerateModal with the scene's prompt pre-filled
   * and the scene's image as reference.
   */
  function handleRegenerate(scene: Scene) {
    openGenerateModal(scene.prompt, scene)
  }

  return { handleDelete, handleRegenerate }
}
