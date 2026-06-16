import { describe, it, expect, vi, afterEach } from 'vitest'
import { render, screen, cleanup, fireEvent } from '@testing-library/react'
import { ChakraProvider } from '@chakra-ui/react'
import { EmptyState } from './EmptyState'
import { system } from '../styles/theme'

afterEach(cleanup)

function renderWithChakra(ui: React.ReactElement) {
  return render(<ChakraProvider value={system}>{ui}</ChakraProvider>)
}

describe('EmptyState', () => {
  it('renders default title and description', () => {
    renderWithChakra(<EmptyState onAction={vi.fn()} />)
    expect(screen.getByText('No maps charted yet')).toBeInTheDocument()
    expect(screen.getByText('Start generating concept art to fill your collection.')).toBeInTheDocument()
  })

  it('renders custom title and description', () => {
    renderWithChakra(
      <EmptyState
        onAction={vi.fn()}
        title="No images"
        description="Click below to create one."
      />
    )
    expect(screen.getByText('No images')).toBeInTheDocument()
    expect(screen.getByText('Click below to create one.')).toBeInTheDocument()
  })

  it('renders call-to-action button', () => {
    renderWithChakra(<EmptyState onAction={vi.fn()} />)
    expect(screen.getByRole('button', { name: 'Generate your first image' })).toBeInTheDocument()
  })

  it('calls onAction when button is clicked', () => {
    const onAction = vi.fn()
    renderWithChakra(<EmptyState onAction={onAction} />)
    fireEvent.click(screen.getByRole('button', { name: 'Generate your first image' }))
    expect(onAction).toHaveBeenCalledOnce()
  })
})
