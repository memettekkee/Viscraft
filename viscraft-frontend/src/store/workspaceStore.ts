import { create } from 'zustand'
import type { Image } from '../types'

interface WorkspaceState {
  activeProjectId: string | null
  generateModalOpen: boolean
  regenerateSource: Image | null
}

interface WorkspaceActions {
  setActiveProject: (projectId: string) => void
  openGenerateModal: () => void
  openRegenerateModal: (source: Image) => void
  closeModal: () => void
}

export const useWorkspaceStore = create<WorkspaceState & WorkspaceActions>((set) => ({
  activeProjectId: null,
  generateModalOpen: false,
  regenerateSource: null,

  setActiveProject: (projectId) => set({ activeProjectId: projectId }),
  openGenerateModal: () => set({ generateModalOpen: true, regenerateSource: null }),
  openRegenerateModal: (source) => set({ generateModalOpen: true, regenerateSource: source }),
  closeModal: () => set({ generateModalOpen: false, regenerateSource: null }),
}))
