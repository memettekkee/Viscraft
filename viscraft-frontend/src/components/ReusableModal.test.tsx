import { describe, it, expect, vi, afterEach } from 'vitest'
import { render, screen, cleanup } from '@testing-library/react'
import { ChakraProvider } from '@chakra-ui/react'
import { ReusableModal } from './ReusableModal'
import { system } from './styles/theme'

afterEach(cleanup)

function renderWithChakra(ui: React.ReactElement) {
  return render(<ChakraProvider value={system}>{ui}</ChakraProvider>)
}

describe('ReusableModal', () => {
  it('renders the title when open', () => {
    renderWithChakra(
      <ReusableModal isOpen={true} onClose={vi.fn()} title="Test Modal">
        <p>Modal content</p>
      </ReusableModal>
    )
    expect(screen.getByText('Test Modal')).toBeInTheDocument()
  })

  it('renders children content when open', () => {
    renderWithChakra(
      <ReusableModal isOpen={true} onClose={vi.fn()} title="Title">
        <p>Hello world</p>
      </ReusableModal>
    )
    expect(screen.getByText('Hello world')).toBeInTheDocument()
  })

  it('does not render content when closed', () => {
    renderWithChakra(
      <ReusableModal isOpen={false} onClose={vi.fn()} title="Hidden">
        <p>Should not appear</p>
      </ReusableModal>
    )
    expect(screen.queryByText('Should not appear')).not.toBeInTheDocument()
  })

  it('renders a close button', () => {
    renderWithChakra(
      <ReusableModal isOpen={true} onClose={vi.fn()} title="Close Test">
        <p>Content</p>
      </ReusableModal>
    )
    // DialogCloseTrigger renders a button with an accessible label
    const closeBtn = screen.getByRole('button', { name: /close dialog/i })
    expect(closeBtn).toBeInTheDocument()
  })

  it('accepts custom size prop', () => {
    // This test verifies the component doesn't throw with a custom size
    renderWithChakra(
      <ReusableModal
        isOpen={true}
        onClose={vi.fn()}
        title="Sized Modal"
        size={{ base: 'full', md: 'xl' }}
      >
        <p>Sized content</p>
      </ReusableModal>
    )
    expect(screen.getByText('Sized content')).toBeInTheDocument()
  })
})
