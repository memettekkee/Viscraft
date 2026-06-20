export type TabType = 'login' | 'register'

export interface AuthModalProps {
  isOpen: boolean
  onClose: () => void
  defaultTab?: TabType
}
