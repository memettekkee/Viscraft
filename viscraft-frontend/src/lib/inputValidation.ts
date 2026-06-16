import { BLOCKED_WORDS } from '../constants'
import type { Genre, AssetType, Mood } from '../types'

// --- Constants ---
const MIN_PROMPT_LENGTH = 3
const MAX_PROMPT_LENGTH = 300
const MAX_REFERENCE_SIZE_BYTES = 5 * 1024 * 1024 // 5MB

// --- Interfaces ---

export interface GenerateFormData {
  prompt: string
  genre: Genre | ''
  assetType: AssetType | ''
  mood: Mood | ''
  referenceImage?: string // base64
}

export interface ValidationResult {
  valid: boolean
  errors: Record<string, string>
}

// --- Helpers ---

/**
 * Estimates the decoded byte size of a base64 string.
 * Base64 encoding ratio: 4 chars encode 3 bytes, minus padding bytes.
 */
export function estimateBase64DecodedSize(base64: string): number {
  const padding = (base64.match(/=+$/) || [''])[0].length
  return Math.floor((base64.length * 3) / 4) - padding
}

// --- Main validation function ---

/**
 * Validates the generate form data, returning all field errors simultaneously.
 * Never throws — always returns a ValidationResult.
 */
export function validateGenerateForm(form: GenerateFormData): ValidationResult {
  const errors: Record<string, string> = {}

  // 1. Prompt length check
  const trimmed = form.prompt.trim()
  if (trimmed.length < MIN_PROMPT_LENGTH) {
    errors.prompt = `Description must be at least ${MIN_PROMPT_LENGTH} characters`
  } else if (trimmed.length > MAX_PROMPT_LENGTH) {
    errors.prompt = `Description must not exceed ${MAX_PROMPT_LENGTH} characters`
  }

  // 2. Blocked word check (case-insensitive) — only if no length error
  if (!errors.prompt) {
    const lower = trimmed.toLowerCase()
    const found = BLOCKED_WORDS.find((word) => lower.includes(word))
    if (found) {
      errors.prompt = 'Description contains a blocked word'
    }
  }

  // 3. Required dropdown fields
  if (!form.genre) errors.genre = 'Genre is required'
  if (!form.assetType) errors.assetType = 'Asset type is required'
  if (!form.mood) errors.mood = 'Mood is required'

  // 4. Reference image validation (if present)
  if (form.referenceImage) {
    const sizeBytes = estimateBase64DecodedSize(form.referenceImage)
    if (sizeBytes > MAX_REFERENCE_SIZE_BYTES) {
      errors.referenceImage = 'Reference image exceeds 5MB'
    }
  }

  return { valid: Object.keys(errors).length === 0, errors }
}
