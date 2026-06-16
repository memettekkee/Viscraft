import { describe, it, expect, vi, afterEach, beforeEach } from 'vitest'
import { render, screen, cleanup, fireEvent, waitFor } from '@testing-library/react'
import { ChakraProvider } from '@chakra-ui/react'
import { SWRConfig } from 'swr'
import { GenerateModal } from './GenerateModal'
import { system } from '../../../components/styles/theme'

// Mock generateImage service
const mockGenerateImage = vi.fn()
vi.mock('../../../service/image', () => ({
  generateImage: (...args: unknown[]) => mockGenerateImage(...args),
}))

// Mock workspaceStore
const mockCloseModal = vi.fn()
vi.mock('../../../store/workspaceStore', () => ({
  useWorkspaceStore: (selector: (state: {
    activeProjectId: string | null
    closeModal: typeof mockCloseModal
  }) => unknown) =>
    selector({ activeProjectId: 'project-123', closeModal: mockCloseModal }),
}))

// Mock useGallery hook
const mockMutate = vi.fn().mockResolvedValue(undefined)
vi.mock('../hooks/useGallery', () => ({
  useGallery: () => ({
    images: [],
    isLoading: false,
    error: undefined,
    mutate: mockMutate,
  }),
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

describe('GenerateModal', () => {
  const defaultProps = {
    isOpen: true,
    onClose: vi.fn(),
  }

  beforeEach(() => {
    mockGenerateImage.mockResolvedValue({
      success: true,
      data: { id: 'img-1', status: 'processing', prompt: 'test', genre: 'fantasy', assetType: 'character', mood: 'dark', createdAt: '2024-01-01' },
      message: 'accepted',
      requestId: 'req-1',
    })
  })

  it('renders the modal title when open', () => {
    renderWithProviders(<GenerateModal {...defaultProps} />)
    expect(screen.getByText('Generate Concept Art')).toBeInTheDocument()
  })

  it('renders mode toggle buttons', () => {
    renderWithProviders(<GenerateModal {...defaultProps} />)
    expect(screen.getByRole('button', { name: 'Create' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'From Reference' })).toBeInTheDocument()
  })

  it('renders prompt textarea with placeholder', () => {
    renderWithProviders(<GenerateModal {...defaultProps} />)
    expect(screen.getByPlaceholderText('Describe the concept art you want to generate...')).toBeInTheDocument()
  })

  it('renders character counter showing 0/300 initially', () => {
    renderWithProviders(<GenerateModal {...defaultProps} />)
    expect(screen.getByText('0/300')).toBeInTheDocument()
  })

  it('updates character counter as user types', () => {
    renderWithProviders(<GenerateModal {...defaultProps} />)
    const textarea = screen.getByPlaceholderText('Describe the concept art you want to generate...')
    fireEvent.change(textarea, { target: { value: 'hello' } })
    expect(screen.getByText('5/300')).toBeInTheDocument()
  })

  it('renders genre, asset type, and mood dropdowns', () => {
    renderWithProviders(<GenerateModal {...defaultProps} />)
    expect(screen.getByText('Genre')).toBeInTheDocument()
    expect(screen.getByText('Asset Type')).toBeInTheDocument()
    expect(screen.getByText('Mood')).toBeInTheDocument()
  })

  it('renders the Generate submit button', () => {
    renderWithProviders(<GenerateModal {...defaultProps} />)
    expect(screen.getByRole('button', { name: 'Generate' })).toBeInTheDocument()
  })

  it('shows validation errors when submitting empty form', async () => {
    renderWithProviders(<GenerateModal {...defaultProps} />)
    const btn = screen.getByRole('button', { name: 'Generate' })
    fireEvent.click(btn)

    await waitFor(() => {
      expect(screen.getByText('Description must be at least 3 characters')).toBeInTheDocument()
      expect(screen.getByText('Genre is required')).toBeInTheDocument()
      expect(screen.getByText('Asset type is required')).toBeInTheDocument()
      expect(screen.getByText('Mood is required')).toBeInTheDocument()
    })
  })

  it('calls generateImage with correct payload on valid submit', async () => {
    renderWithProviders(<GenerateModal {...defaultProps} />)

    // Fill prompt
    const textarea = screen.getByPlaceholderText('Describe the concept art you want to generate...')
    fireEvent.change(textarea, { target: { value: 'A dark forest with tall trees' } })

    // Fill genre
    const selects = screen.getAllByRole('combobox') as HTMLSelectElement[]
    fireEvent.change(selects[0], { target: { value: 'fantasy' } })
    fireEvent.change(selects[1], { target: { value: 'location' } })
    fireEvent.change(selects[2], { target: { value: 'dark' } })

    // Submit
    const btn = screen.getByRole('button', { name: 'Generate' })
    fireEvent.click(btn)

    await waitFor(() => {
      expect(mockGenerateImage).toHaveBeenCalledWith({
        projectId: 'project-123',
        prompt: 'A dark forest with tall trees',
        genre: 'fantasy',
        assetType: 'location',
        mood: 'dark',
        referenceImage: undefined,
      })
    })
  })

  it('closes modal and mutates cache on successful generation', async () => {
    renderWithProviders(<GenerateModal {...defaultProps} />)

    const textarea = screen.getByPlaceholderText('Describe the concept art you want to generate...')
    fireEvent.change(textarea, { target: { value: 'A dark forest with tall trees' } })

    const selects = screen.getAllByRole('combobox') as HTMLSelectElement[]
    fireEvent.change(selects[0], { target: { value: 'fantasy' } })
    fireEvent.change(selects[1], { target: { value: 'location' } })
    fireEvent.change(selects[2], { target: { value: 'dark' } })

    const btn = screen.getByRole('button', { name: 'Generate' })
    fireEvent.click(btn)

    await waitFor(() => {
      expect(mockMutate).toHaveBeenCalled()
      expect(mockCloseModal).toHaveBeenCalled()
    })
  })

  it('shows rate limit error inline on ERR_02 and keeps modal open', async () => {
    mockGenerateImage.mockRejectedValueOnce({
      response: {
        data: {
          success: false,
          errorCode: 'ERR_02',
          message: 'rate limited',
          requestId: 'req-2',
        },
      },
    })

    renderWithProviders(<GenerateModal {...defaultProps} />)

    const textarea = screen.getByPlaceholderText('Describe the concept art you want to generate...')
    fireEvent.change(textarea, { target: { value: 'A dark forest with tall trees' } })

    const selects = screen.getAllByRole('combobox') as HTMLSelectElement[]
    fireEvent.change(selects[0], { target: { value: 'fantasy' } })
    fireEvent.change(selects[1], { target: { value: 'location' } })
    fireEvent.change(selects[2], { target: { value: 'dark' } })

    const btn = screen.getByRole('button', { name: 'Generate' })
    fireEvent.click(btn)

    await waitFor(() => {
      expect(screen.getByText('Too many requests, please wait')).toBeInTheDocument()
    })
    // Modal stays open — onClose not called
    expect(mockCloseModal).not.toHaveBeenCalled()
  })

  it('shows network error when request fails without response', async () => {
    mockGenerateImage.mockRejectedValueOnce(new Error('Network Error'))

    renderWithProviders(<GenerateModal {...defaultProps} />)

    const textarea = screen.getByPlaceholderText('Describe the concept art you want to generate...')
    fireEvent.change(textarea, { target: { value: 'A dark forest with tall trees' } })

    const selects = screen.getAllByRole('combobox') as HTMLSelectElement[]
    fireEvent.change(selects[0], { target: { value: 'fantasy' } })
    fireEvent.change(selects[1], { target: { value: 'location' } })
    fireEvent.change(selects[2], { target: { value: 'dark' } })

    const btn = screen.getByRole('button', { name: 'Generate' })
    fireEvent.click(btn)

    await waitFor(() => {
      expect(screen.getByText('Unable to connect to server')).toBeInTheDocument()
    })
  })

  it('does not render content when closed', () => {
    renderWithProviders(<GenerateModal isOpen={false} onClose={vi.fn()} />)
    expect(screen.queryByText('Generate Concept Art')).not.toBeInTheDocument()
  })

  it('shows blocked word validation error', async () => {
    renderWithProviders(<GenerateModal {...defaultProps} />)

    const textarea = screen.getByPlaceholderText('Describe the concept art you want to generate...')
    fireEvent.change(textarea, { target: { value: 'Something explicit here' } })

    const selects = screen.getAllByRole('combobox') as HTMLSelectElement[]
    fireEvent.change(selects[0], { target: { value: 'fantasy' } })
    fireEvent.change(selects[1], { target: { value: 'location' } })
    fireEvent.change(selects[2], { target: { value: 'dark' } })

    const btn = screen.getByRole('button', { name: 'Generate' })
    fireEvent.click(btn)

    await waitFor(() => {
      expect(screen.getByText('Description contains a blocked word')).toBeInTheDocument()
    })
    // Should not call the API
    expect(mockGenerateImage).not.toHaveBeenCalled()
  })
})
