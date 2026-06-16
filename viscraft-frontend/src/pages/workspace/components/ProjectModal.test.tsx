import { describe, it, expect, vi, afterEach, beforeEach } from 'vitest'
import { render, screen, cleanup, fireEvent, waitFor } from '@testing-library/react'
import { ChakraProvider } from '@chakra-ui/react'
import { SWRConfig } from 'swr'
import { ProjectModal } from './ProjectModal'
import { system } from '../../../components/styles/theme'

// Mock createProject service
const mockCreateProject = vi.fn()
vi.mock('../../../service/project', () => ({
  createProject: (...args: unknown[]) => mockCreateProject(...args),
}))

// Mock workspaceStore
const mockSetActiveProject = vi.fn()
vi.mock('../../../store/workspaceStore', () => ({
  useWorkspaceStore: (selector: (state: { setActiveProject: typeof mockSetActiveProject }) => unknown) =>
    selector({ setActiveProject: mockSetActiveProject }),
}))

afterEach(() => {
  cleanup()
  vi.clearAllMocks()
})

function renderWithProviders(ui: React.ReactElement) {
  return render(
    <SWRConfig value={{ provider: () => new Map() }}>
      <ChakraProvider value={system}>{ui}</ChakraProvider>
    </SWRConfig>
  )
}

describe('ProjectModal', () => {
  const defaultProps = {
    isOpen: true,
    onClose: vi.fn(),
  }

  beforeEach(() => {
    mockCreateProject.mockResolvedValue({
      success: true,
      data: { id: 'proj-123', name: 'Test Project', createdAt: '2024-01-01' },
      message: 'ok',
      requestId: 'req-1',
    })
  })

  it('renders the modal title when open', () => {
    renderWithProviders(<ProjectModal {...defaultProps} />)
    expect(screen.getByText('New Region')).toBeInTheDocument()
  })

  it('renders name and description fields', () => {
    renderWithProviders(<ProjectModal {...defaultProps} />)
    expect(screen.getByPlaceholderText('Enter project name')).toBeInTheDocument()
    expect(screen.getByPlaceholderText('Describe your project')).toBeInTheDocument()
  })

  it('renders the submit button', () => {
    renderWithProviders(<ProjectModal {...defaultProps} />)
    expect(screen.getByRole('button', { name: 'Create Region' })).toBeInTheDocument()
  })

  it('submit button is disabled when name is empty', () => {
    renderWithProviders(<ProjectModal {...defaultProps} />)
    const btn = screen.getByRole('button', { name: 'Create Region' })
    expect(btn).toBeDisabled()
  })

  it('submit button is enabled when name is valid', () => {
    renderWithProviders(<ProjectModal {...defaultProps} />)
    const input = screen.getByPlaceholderText('Enter project name')
    fireEvent.change(input, { target: { value: 'My Project' } })
    const btn = screen.getByRole('button', { name: 'Create Region' })
    expect(btn).not.toBeDisabled()
  })

  it('shows validation error for whitespace-only name', async () => {
    renderWithProviders(<ProjectModal {...defaultProps} />)
    const input = screen.getByPlaceholderText('Enter project name')
    // Simulate typing spaces then clearing to trigger submit attempt
    fireEvent.change(input, { target: { value: '   ' } })
    // The button should still be disabled because trim results in empty
    const btn = screen.getByRole('button', { name: 'Create Region' })
    expect(btn).toBeDisabled()
  })

  it('calls createProject with correct payload on submit', async () => {
    renderWithProviders(<ProjectModal {...defaultProps} />)
    const nameInput = screen.getByPlaceholderText('Enter project name')
    const descInput = screen.getByPlaceholderText('Describe your project')

    fireEvent.change(nameInput, { target: { value: 'My Project' } })
    fireEvent.change(descInput, { target: { value: 'A description' } })

    const btn = screen.getByRole('button', { name: 'Create Region' })
    fireEvent.click(btn)

    await waitFor(() => {
      expect(mockCreateProject).toHaveBeenCalledWith({
        name: 'My Project',
        description: 'A description',
      })
    })
  })

  it('sets new project as active after successful creation', async () => {
    renderWithProviders(<ProjectModal {...defaultProps} />)
    const nameInput = screen.getByPlaceholderText('Enter project name')
    fireEvent.change(nameInput, { target: { value: 'My Project' } })

    const btn = screen.getByRole('button', { name: 'Create Region' })
    fireEvent.click(btn)

    await waitFor(() => {
      expect(mockSetActiveProject).toHaveBeenCalledWith('proj-123')
    })
  })

  it('calls onClose after successful creation', async () => {
    const onClose = vi.fn()
    renderWithProviders(<ProjectModal isOpen={true} onClose={onClose} />)
    const nameInput = screen.getByPlaceholderText('Enter project name')
    fireEvent.change(nameInput, { target: { value: 'My Project' } })

    const btn = screen.getByRole('button', { name: 'Create Region' })
    fireEvent.click(btn)

    await waitFor(() => {
      expect(onClose).toHaveBeenCalled()
    })
  })

  it('shows API error on failed creation', async () => {
    mockCreateProject.mockRejectedValueOnce(new Error('Network error'))
    renderWithProviders(<ProjectModal {...defaultProps} />)
    const nameInput = screen.getByPlaceholderText('Enter project name')
    fireEvent.change(nameInput, { target: { value: 'My Project' } })

    const btn = screen.getByRole('button', { name: 'Create Region' })
    fireEvent.click(btn)

    await waitFor(() => {
      expect(screen.getByText('Unable to connect to server')).toBeInTheDocument()
    })
  })

  it('does not render content when closed', () => {
    renderWithProviders(<ProjectModal isOpen={false} onClose={vi.fn()} />)
    expect(screen.queryByText('New Region')).not.toBeInTheDocument()
  })

  it('omits description from payload when empty', async () => {
    renderWithProviders(<ProjectModal {...defaultProps} />)
    const nameInput = screen.getByPlaceholderText('Enter project name')
    fireEvent.change(nameInput, { target: { value: 'Name Only' } })

    const btn = screen.getByRole('button', { name: 'Create Region' })
    fireEvent.click(btn)

    await waitFor(() => {
      expect(mockCreateProject).toHaveBeenCalledWith({
        name: 'Name Only',
        description: undefined,
      })
    })
  })
})
