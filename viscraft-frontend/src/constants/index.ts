import type { ProductCategory } from '../types'

/**
 * Centralized error message mapping.
 * All UI-facing error strings MUST come from this constant — never display raw backend messages.
 */
export const ERROR_MESSAGES: Record<string, string> = {
  ERR_01: 'Resource not found',
  ERR_02: 'Too many requests, please wait',
  ERR_03: 'Request timed out',
  ERR_04: 'Validation failed',
  ERR_05: 'Invalid AI response',
  ERR_06: 'Generation failed',
  ERR_07: 'Generation failed',
  ERR_08: 'Campaign not found',
  ERR_09: 'Session expired, please log in again',
  ERR_13: 'Content policy violation — prompt rejected',
  NETWORK_ERROR: 'Unable to connect to server',
}

/** Dropdown option shape for select components */
export interface DropdownOption<T extends string = string> {
  value: T
  label: string
}

/** Product category options */
export const PRODUCT_CATEGORY_OPTIONS: DropdownOption<ProductCategory>[] = [
  { value: 'general', label: 'General' },
  { value: 'food', label: 'Food & Snacks' },
  { value: 'beverage', label: 'Beverages' },
  { value: 'cosmetics', label: 'Cosmetics & Skincare' },
  { value: 'fashion', label: 'Fashion & Accessories' },
  { value: 'electronics', label: 'Electronics' },
  { value: 'home', label: 'Home & Living' },
]

/**
 * Words that are not allowed in prompts (case-insensitive check).
 */
export const BLOCKED_WORDS: string[] = ['nude', 'explicit', 'nsfw', 'gore']
