import { create } from 'zustand'

interface AuthState {
  token: string | null
  user: unknown | null
  isAuthenticated: boolean
  setAuth: (token: string, user: unknown) => void
  clearAuth: () => void
  updateUser: (user: unknown) => void
}

export const useAuthStore = create<AuthState>((set) => ({
  token: null,
  user: null,
  isAuthenticated: false,
  setAuth: (token, user) => set({ token, user, isAuthenticated: true }),
  clearAuth: () => set({ token: null, user: null, isAuthenticated: false }),
  updateUser: (user) => set({ user }),
}))
