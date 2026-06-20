import { create } from 'zustand'
import type { Scene } from '../types'

interface WorkspaceState {
  activeProjectId: string | null
  generateModalOpen: boolean
  prefillPrompt: string | null
  regenerateSceneId: string | null
  regenerateFileUrl: string | null
  selectedScene: Scene | null
}

interface WorkspaceActions {
  setActiveProject: (projectId: string) => void
  clearActiveProject: () => void
  openGenerateModal: (prefillPrompt?: string, regenerateScene?: Scene) => void
  closeModal: () => void
  openSceneDetail: (scene: Scene) => void
  closeSceneDetail: () => void
}

export const useWorkspaceStore = create<WorkspaceState & WorkspaceActions>((set) => ({
  activeProjectId: null,
  generateModalOpen: false,
  prefillPrompt: null,
  regenerateSceneId: null,
  regenerateFileUrl: null,
  selectedScene: null,

  setActiveProject: (projectId) => set({ activeProjectId: projectId }),
  clearActiveProject: () => set({ activeProjectId: null }),
  openGenerateModal: (prefillPrompt, regenerateScene) => set({
    generateModalOpen: true,
    prefillPrompt: prefillPrompt ?? null,
    regenerateSceneId: regenerateScene?.id ?? null,
    regenerateFileUrl: regenerateScene?.fileUrl ?? null,
  }),
  closeModal: () => set({ generateModalOpen: false, prefillPrompt: null, regenerateSceneId: null, regenerateFileUrl: null }),
  openSceneDetail: (scene) => set({ selectedScene: scene }),
  closeSceneDetail: () => set({ selectedScene: null }),
}))
