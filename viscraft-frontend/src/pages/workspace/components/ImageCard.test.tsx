import { describe, it, expect, vi, afterEach, beforeEach } from 'vitest'
import { render, screen, cleanup, fireEvent } from '@testing-library/react'
import { ChakraProvider } from '@chakra-ui/react'
import { SWRConfig } from 'swr'
import { ImageCard } from './ImageCard'
import { system } from '../../../components/styles/theme'
import type { Image } from '../../../types'

// Mock the image service
vi.mock('../../../service/image', () => ({
  generateImage: vi.fn().mockResolvedValue({ success: true, data: { id: 'new-1', status: 'processing' } }),
  deleteImage: vi.fn().mockResolvedValue({ success: true }),
}))

// Mock workspaceStore
vi.mock('../../../store/workspaceStore', () => ({
  useWorkspaceStore: (selector: (state: Record<string, unknown>) => unknown) =>
    selector({
      activeProjectId: 'project-1',
      openRegenerateModal: vi.fn(),
    }),
}))

afterEach(cleanup)

function renderWithProviders(ui: React.ReactElement) {
  return render(
    <SWRConfig value={{ provider: () => new Map() }}>
      <ChakraProvider value={system}>{ui}</ChakraProvider>
    </SWRConfig>
  )
}

const completedImage: Image = {
  id: 'img-1',
  status: 'completed',
  fileUrl: '/storage/images/img-1.png',
  prompt: 'A dark fantasy castle surrounded by mist and ancient trees in a forgotten land of sorrow',
  genre: 'fantasy',
  assetType: 'location',
  mood: 'dark',
  createdAt: '2024-01-01T00:00:00Z',
}

const failedImageTimeout: Image = {
  id: 'img-2',
  status: 'failed',
  prompt: 'A sci-fi spaceship',
  genre: 'sci-fi',
  assetType: 'item',
  mood: 'epic',
  errorCode: 'ERR_03',
  createdAt: '2024-01-01T00:00:00Z',
}

const failedImageInvalidAI: Image = {
  id: 'img-3',
  status: 'failed',
  prompt: 'A horror creature',
  genre: 'horror',
  assetType: 'creature',
  mood: 'mysterious',
  errorCode: 'ERR_05',
  createdAt: '2024-01-01T00:00:00Z',
}

const failedImageGeneric: Image = {
  id: 'img-4',
  status: 'failed',
  prompt: 'A steampunk item',
  genre: 'steampunk',
  assetType: 'item',
  mood: 'whimsical',
  errorCode: 'ERR_06',
  createdAt: '2024-01-01T00:00:00Z',
}

const processingImage: Image = {
  id: 'img-5',
  status: 'processing',
  prompt: 'A post-apocalyptic scene',
  genre: 'post-apocalyptic',
  assetType: 'location',
  mood: 'dark',
  createdAt: '2024-01-01T00:00:00Z',
}

describe('ImageCard', () => {
  const onRegenerate = vi.fn()
  const onDelete = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Processing state', () => {
    it('renders ImageCardSkeleton with "Mapping..." label', () => {
      renderWithProviders(
        <ImageCard image={processingImage} onRegenerate={onRegenerate} onDelete={onDelete} />
      )
      expect(screen.getByText('Mapping...')).toBeInTheDocument()
    })
  })

  describe('Completed state', () => {
    it('renders the image element', () => {
      renderWithProviders(
        <ImageCard image={completedImage} onRegenerate={onRegenerate} onDelete={onDelete} />
      )
      const img = screen.getByRole('img')
      expect(img).toBeInTheDocument()
      expect(img).toHaveAttribute('src', expect.stringContaining('/storage/images/img-1.png'))
    })

    it('displays truncated prompt (60 chars max)', () => {
      renderWithProviders(
        <ImageCard image={completedImage} onRegenerate={onRegenerate} onDelete={onDelete} />
      )
      // The full prompt is 87 chars; truncated should be 60 chars + ellipsis
      const promptText = screen.getByText(/A dark fantasy castle/)
      expect(promptText.textContent!.length).toBeLessThanOrEqual(61) // 60 + ellipsis char
    })

    it('renders genre·mood stamp badge in uppercase', () => {
      renderWithProviders(
        <ImageCard image={completedImage} onRegenerate={onRegenerate} onDelete={onDelete} />
      )
      const badge = screen.getByTestId('stamp-badge-img-1')
      expect(badge).toBeInTheDocument()
      expect(badge.textContent).toBe('FANTASY · DARK')
    })

    it('renders Regenerate and Delete action buttons', () => {
      renderWithProviders(
        <ImageCard image={completedImage} onRegenerate={onRegenerate} onDelete={onDelete} />
      )
      expect(screen.getByRole('button', { name: 'Regenerate' })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: 'Delete' })).toBeInTheDocument()
    })

    it('calls onRegenerate when Regenerate is clicked', () => {
      renderWithProviders(
        <ImageCard image={completedImage} onRegenerate={onRegenerate} onDelete={onDelete} />
      )
      fireEvent.click(screen.getByRole('button', { name: 'Regenerate' }))
      expect(onRegenerate).toHaveBeenCalledWith(completedImage)
    })

    it('calls onDelete when Delete is clicked', () => {
      renderWithProviders(
        <ImageCard image={completedImage} onRegenerate={onRegenerate} onDelete={onDelete} />
      )
      fireEvent.click(screen.getByRole('button', { name: 'Delete' }))
      expect(onDelete).toHaveBeenCalledWith('img-1')
    })
  })

  describe('Failed state', () => {
    it('displays "Request timed out" for ERR_03', () => {
      renderWithProviders(
        <ImageCard image={failedImageTimeout} onRegenerate={onRegenerate} onDelete={onDelete} />
      )
      expect(screen.getByText('Request timed out')).toBeInTheDocument()
    })

    it('displays "Invalid AI response" for ERR_05', () => {
      renderWithProviders(
        <ImageCard image={failedImageInvalidAI} onRegenerate={onRegenerate} onDelete={onDelete} />
      )
      expect(screen.getByText('Invalid AI response')).toBeInTheDocument()
    })

    it('displays "Image generation failed" for ERR_06', () => {
      renderWithProviders(
        <ImageCard image={failedImageGeneric} onRegenerate={onRegenerate} onDelete={onDelete} />
      )
      expect(screen.getByText('Image generation failed')).toBeInTheDocument()
    })

    it('renders broken-map icon', () => {
      renderWithProviders(
        <ImageCard image={failedImageTimeout} onRegenerate={onRegenerate} onDelete={onDelete} />
      )
      expect(screen.getByText('🗺️✕')).toBeInTheDocument()
    })

    it('renders Retry button', () => {
      renderWithProviders(
        <ImageCard image={failedImageTimeout} onRegenerate={onRegenerate} onDelete={onDelete} />
      )
      expect(screen.getByRole('button', { name: 'Retry' })).toBeInTheDocument()
    })
  })
})
