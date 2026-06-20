import { BLOCKED_WORDS } from '../constants'

const MIN_PROMPT_LENGTH = 3
const MAX_PROMPT_LENGTH = 300

export interface GenerateSceneFormData {
  prompt: string
}

export interface CreateProjectFormData {
  name: string
  description?: string
  productCategory?: string
  visualStyle?: string
}

export interface ValidationResult {
  valid: boolean
  errors: Record<string, string>
}

/**
 * Validates the generate form data.
 */
export function validateGenerateSceneForm(form: GenerateSceneFormData): ValidationResult {
  const errors: Record<string, string> = {}

  const trimmed = form.prompt.trim()
  if (trimmed.length < MIN_PROMPT_LENGTH) {
    errors.prompt = `Description must be at least ${MIN_PROMPT_LENGTH} characters`
  } else if (trimmed.length > MAX_PROMPT_LENGTH) {
    errors.prompt = `Description must not exceed ${MAX_PROMPT_LENGTH} characters`
  }

  if (!errors.prompt) {
    const lower = trimmed.toLowerCase()
    const found = BLOCKED_WORDS.find((word) => lower.includes(word))
    if (found) {
      errors.prompt = 'Description contains a blocked word'
    }
  }

  return { valid: Object.keys(errors).length === 0, errors }
}

/**
 * Validates the create project (campaign) form data.
 */
export function validateCreateProjectForm(form: CreateProjectFormData): ValidationResult {
  const errors: Record<string, string> = {}

  const name = (form.name || '').trim()
  if (name.length === 0) {
    errors.name = 'Campaign name is required'
  } else if (name.length > 255) {
    errors.name = 'Campaign name must not exceed 255 characters'
  }

  return { valid: Object.keys(errors).length === 0, errors }
}
