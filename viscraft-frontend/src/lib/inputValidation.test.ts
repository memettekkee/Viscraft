import { describe, it, expect } from 'vitest'
import {
  validateGenerateForm,
  estimateBase64DecodedSize,
} from './inputValidation'
import type { GenerateFormData } from './inputValidation'

function validForm(overrides: Partial<GenerateFormData> = {}): GenerateFormData {
  return {
    prompt: 'A dark fantasy castle on a cliff',
    genre: 'fantasy',
    assetType: 'location',
    mood: 'dark',
    ...overrides,
  }
}

describe('validateGenerateForm', () => {
  describe('prompt length validation', () => {
    it('returns error when prompt is fewer than 3 characters (trimmed)', () => {
      const result = validateGenerateForm(validForm({ prompt: 'ab' }))
      expect(result.valid).toBe(false)
      expect(result.errors.prompt).toBe('Description must be at least 3 characters')
    })

    it('returns error when prompt is only whitespace below 3 chars', () => {
      const result = validateGenerateForm(validForm({ prompt: '   a  ' }))
      expect(result.valid).toBe(false)
      expect(result.errors.prompt).toBe('Description must be at least 3 characters')
    })

    it('returns error when prompt exceeds 300 characters (trimmed)', () => {
      const longPrompt = 'x'.repeat(301)
      const result = validateGenerateForm(validForm({ prompt: longPrompt }))
      expect(result.valid).toBe(false)
      expect(result.errors.prompt).toBe('Description must not exceed 300 characters')
    })

    it('accepts prompt at exactly 3 characters', () => {
      const result = validateGenerateForm(validForm({ prompt: 'abc' }))
      expect(result.valid).toBe(true)
      expect(result.errors.prompt).toBeUndefined()
    })

    it('accepts prompt at exactly 300 characters', () => {
      const result = validateGenerateForm(validForm({ prompt: 'x'.repeat(300) }))
      expect(result.valid).toBe(true)
      expect(result.errors.prompt).toBeUndefined()
    })
  })

  describe('blocked word detection', () => {
    it('returns error when prompt contains a blocked word (lowercase)', () => {
      const result = validateGenerateForm(validForm({ prompt: 'draw a nude character' }))
      expect(result.valid).toBe(false)
      expect(result.errors.prompt).toBe('Description contains a blocked word')
    })

    it('returns error when prompt contains a blocked word (mixed case)', () => {
      const result = validateGenerateForm(validForm({ prompt: 'An EXPLICIT scene' }))
      expect(result.valid).toBe(false)
      expect(result.errors.prompt).toBe('Description contains a blocked word')
    })

    it('detects nsfw case-insensitively', () => {
      const result = validateGenerateForm(validForm({ prompt: 'something NSFW here' }))
      expect(result.valid).toBe(false)
      expect(result.errors.prompt).toBe('Description contains a blocked word')
    })

    it('detects gore case-insensitively', () => {
      const result = validateGenerateForm(validForm({ prompt: 'a Gore filled scene' }))
      expect(result.valid).toBe(false)
      expect(result.errors.prompt).toBe('Description contains a blocked word')
    })

    it('does not flag clean prompts', () => {
      const result = validateGenerateForm(validForm({ prompt: 'A peaceful forest clearing' }))
      expect(result.valid).toBe(true)
      expect(result.errors.prompt).toBeUndefined()
    })
  })

  describe('required dropdown validation', () => {
    it('returns error when genre is empty', () => {
      const result = validateGenerateForm(validForm({ genre: '' }))
      expect(result.valid).toBe(false)
      expect(result.errors.genre).toBe('Genre is required')
    })

    it('returns error when assetType is empty', () => {
      const result = validateGenerateForm(validForm({ assetType: '' }))
      expect(result.valid).toBe(false)
      expect(result.errors.assetType).toBe('Asset type is required')
    })

    it('returns error when mood is empty', () => {
      const result = validateGenerateForm(validForm({ mood: '' }))
      expect(result.valid).toBe(false)
      expect(result.errors.mood).toBe('Mood is required')
    })
  })

  describe('reference image size validation', () => {
    it('returns error when reference image exceeds 5MB', () => {
      // 5MB = 5,242,880 bytes. Base64 of that is ~6,990,507 chars + padding.
      // Create a base64 string that decodes to > 5MB
      const oversizeBase64 = 'A'.repeat(7_000_000)
      const result = validateGenerateForm(validForm({ referenceImage: oversizeBase64 }))
      expect(result.valid).toBe(false)
      expect(result.errors.referenceImage).toBe('Reference image exceeds 5MB')
    })

    it('accepts reference image under 5MB', () => {
      const smallBase64 = 'A'.repeat(1000)
      const result = validateGenerateForm(validForm({ referenceImage: smallBase64 }))
      expect(result.valid).toBe(true)
      expect(result.errors.referenceImage).toBeUndefined()
    })

    it('does not validate when referenceImage is absent', () => {
      const result = validateGenerateForm(validForm({ referenceImage: undefined }))
      expect(result.valid).toBe(true)
      expect(result.errors.referenceImage).toBeUndefined()
    })
  })

  describe('simultaneous error reporting', () => {
    it('returns all errors at once when multiple fields are invalid', () => {
      const result = validateGenerateForm({
        prompt: 'ab',
        genre: '',
        assetType: '',
        mood: '',
      })
      expect(result.valid).toBe(false)
      expect(result.errors.prompt).toBeDefined()
      expect(result.errors.genre).toBeDefined()
      expect(result.errors.assetType).toBeDefined()
      expect(result.errors.mood).toBeDefined()
    })
  })

  describe('valid form', () => {
    it('returns valid: true with empty errors for a fully valid form', () => {
      const result = validateGenerateForm(validForm())
      expect(result.valid).toBe(true)
      expect(result.errors).toEqual({})
    })
  })
})

describe('estimateBase64DecodedSize', () => {
  it('correctly estimates size without padding', () => {
    // 4 base64 chars = 3 bytes
    expect(estimateBase64DecodedSize('AAAA')).toBe(3)
  })

  it('correctly estimates size with single padding char', () => {
    // "AAA=" → 4 chars, 1 padding → (4*3/4) - 1 = 2
    expect(estimateBase64DecodedSize('AAA=')).toBe(2)
  })

  it('correctly estimates size with double padding', () => {
    // "AA==" → 4 chars, 2 padding → (4*3/4) - 2 = 1
    expect(estimateBase64DecodedSize('AA==')).toBe(1)
  })

  it('estimates larger strings correctly', () => {
    // 100 chars, no padding → floor(100*3/4) - 0 = 75
    const str = 'A'.repeat(100)
    expect(estimateBase64DecodedSize(str)).toBe(75)
  })
})
