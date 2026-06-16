/**
 * Local types for the landing page.
 */

export type TabType = 'login' | 'register'

export interface AuthModalProps {
  isOpen: boolean
  onClose: () => void
  defaultTab?: TabType
}
