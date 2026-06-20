import { useEffect } from 'react'
import { useSceneList } from '../hooks/useSceneList'
import { useSceneActions } from '../hooks/useSceneActions'
import { usePollSceneStatus } from '../hooks/usePollSceneStatus'
import { ImageCardSkeleton } from '../../../components/skeleton/ImageCardSkeleton'
import { SceneCard } from './SceneCard'
import type { Scene } from '../../../types'

interface PollingSceneCardProps {
  scene: Scene
  scenes: Scene[]
  projectId: string
  onDelete: (sceneId: string) => void
}

export function PollingSceneCard({
  scene,
  scenes,
  projectId,
  onDelete,
}: PollingSceneCardProps) {
  const { mutate: mutateSceneList } = useSceneList(projectId)
  const { handleRegenerate } = useSceneActions()
  const { scene: polledScene } = usePollSceneStatus(scene.id)

  useEffect(() => {
    if (polledScene && polledScene.status !== 'processing') {
      mutateSceneList()
    }
  }, [polledScene?.status]) 

  const resolved = polledScene && polledScene.status !== 'processing' ? polledScene : null

  if (!resolved) {
    return <ImageCardSkeleton label="Generating..." />
  }

  return (
    <SceneCard
      scene={resolved}
      scenes={scenes}
      onDelete={onDelete}
      onRegenerate={handleRegenerate}
    />
  )
}
