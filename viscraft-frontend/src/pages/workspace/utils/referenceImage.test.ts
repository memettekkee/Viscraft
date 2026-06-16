import { describe, it, expect } from 'vitest'
import { validateImageFile } from './referenceImage'

describe('referenceImage utils', () => {
  describe('validateImageFile', () => {
    it('returns null for valid JPEG file under 5MB', () => {
      const file = new File(['x'.repeat(1000)], 'photo.jpg', { type: 'image/jpeg' })
      expect(validateImageFile(file)).toBeNull()
    })

    it('returns null for valid PNG file', () => {
      const file = new File(['x'.repeat(1000)], 'image.png', { type: 'image/png' })
      expect(validateImageFile(file)).toBeNull()
    })

    it('returns null for valid WEBP file', () => {
      const file = new File(['x'.repeat(1000)], 'photo.webp', { type: 'image/webp' })
      expect(validateImageFile(file)).toBeNull()
    })

    it('returns invalid-type error for GIF files', () => {
      const file = new File(['x'.repeat(100)], 'anim.gif', { type: 'image/gif' })
      const result = validateImageFile(file)
      expect(result).not.toBeNull()
      expect(result!.type).toBe('invalid-type')
      expect(result!.message).toBe('Only JPEG, PNG, and WEBP images are accepted')
    })

    it('returns invalid-type error for SVG files', () => {
      const file = new File(['<svg></svg>'], 'icon.svg', { type: 'image/svg+xml' })
      const result = validateImageFile(file)
      expect(result).not.toBeNull()
      expect(result!.type).toBe('invalid-type')
    })

    it('returns invalid-type error for non-image files', () => {
      const file = new File(['hello'], 'doc.pdf', { type: 'application/pdf' })
      const result = validateImageFile(file)
      expect(result).not.toBeNull()
      expect(result!.type).toBe('invalid-type')
    })

    it('returns too-large error for files over 5MB', () => {
      // Create a file just over 5MB
      const size = 5 * 1024 * 1024 + 1
      const file = new File([new ArrayBuffer(size)], 'big.png', { type: 'image/png' })
      // File constructor may not preserve exact size with ArrayBuffer, so use Object.defineProperty
      Object.defineProperty(file, 'size', { value: size })
      const result = validateImageFile(file)
      expect(result).not.toBeNull()
      expect(result!.type).toBe('too-large')
      expect(result!.message).toBe('Image must be smaller than 5MB')
    })

    it('returns null for a file exactly at 5MB', () => {
      const size = 5 * 1024 * 1024
      const file = new File([new ArrayBuffer(10)], 'exact.jpg', { type: 'image/jpeg' })
      Object.defineProperty(file, 'size', { value: size })
      expect(validateImageFile(file)).toBeNull()
    })

    it('checks type before size (invalid type takes priority)', () => {
      const size = 10 * 1024 * 1024 // 10MB
      const file = new File(['x'], 'big.gif', { type: 'image/gif' })
      Object.defineProperty(file, 'size', { value: size })
      const result = validateImageFile(file)
      expect(result).not.toBeNull()
      expect(result!.type).toBe('invalid-type')
    })
  })
})
