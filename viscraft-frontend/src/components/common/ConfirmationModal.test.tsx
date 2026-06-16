import { describe, it, expect, vi, afterEach } from 'vitest'
import { render, screen, cleanup, fireEvent } from '@testing-library/react'
import { ChakraProvider } from '@chakra-ui/react'
import { ConfirmationModal } from './ConfirmationModal'
import { system } from '../styles/theme'

afterEach(cleanup)

function renderWithChakra(ui: React.ReactElement) {
  return render(<ChakraProvider value={system}>{ui}</ChakraProvider>)
}

describe('ConfirmationModal', () => {
  const defaultProps = {
    isOpen: true,
    onClose: vi.fn(),
    onConfirm: vi.fn(),
    title: 'Delete Project',
    message: 'Are you sure you want to delete this project?',
  }

  it('renders title and message when open', () => {
    renderWithChakra(<ConfirmationModal {...defaultProps} />)
    expect(screen.getByText('Delete Project')).toBeInTheDocument()
    expect(screen.getByText('Are you sure you want to delete this project?')).toBeInTheDocument()
  })

  it('does not render content when closed', () => {
    renderWithChakra(<ConfirmationModal {...defaultProps} isOpen={false} />)
    expect(screen.queryByText('Are you sure you want to delete this project?')).not.toBeInTheDocument()
  })

  it('renders default confirm label as "Delete"', () => {
    renderWithChakra(<ConfirmationModal {...defaultProps} />)
    expect(screen.getByRole('button', { name: 'Delete' })).toBeInTheDocument()
  })

  it('renders custom confirm label', () => {
    renderWithChakra(<ConfirmationModal {...defaultProps} confirmLabel="Remove" />)
    expect(screen.getByRole('button', { name: 'Remove' })).toBeInTheDocument()
  })

  it('renders Cancel button', () => {
    renderWithChakra(<ConfirmationModal {...defaultProps} />)
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument()
  })

  it('calls onConfirm when confirm button is clicked', () => {
    const onConfirm = vi.fn()
    renderWithChakra(<ConfirmationModal {...defaultProps} onConfirm={onConfirm} />)
    fireEvent.click(screen.getByRole('button', { name: 'Delete' }))
    expect(onConfirm).toHaveBeenCalledOnce()
  })

  it('calls onClose when cancel button is clicked', () => {
    const onClose = vi.fn()
    renderWithChakra(<ConfirmationModal {...defaultProps} onClose={onClose} />)
    fireEvent.click(screen.getByRole('button', { name: 'Cancel' }))
    expect(onClose).toHaveBeenCalledOnce()
  })

  it('disables both buttons when isLoading is true', () => {
    renderWithChakra(<ConfirmationModal {...defaultProps} isLoading={true} />)
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeDisabled()
  })
})
