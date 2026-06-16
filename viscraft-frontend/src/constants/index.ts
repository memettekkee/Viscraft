import type { Genre, AssetType, Mood } from '../types'

/**
 * Centralized error message mapping.
 * All UI-facing error strings MUST come from this constant — never display raw backend messages.
 * @see Requirement 12.6
 */
export const ERROR_MESSAGES: Record<string, string> = {
  ERR_01: 'Resource not found',
  ERR_02: 'Too many requests, please wait',
  ERR_03: 'Request timed out',
  ERR_04: 'Validation failed',
  ERR_05: 'Invalid AI response',
  ERR_06: 'Image generation failed',
  ERR_07: 'Image generation failed',
  ERR_08: 'Project not found',
  ERR_09: 'Session expired, please log in again',
  NETWORK_ERROR: 'Unable to connect to server',
}

/** Dropdown option shape for select components */
export interface DropdownOption<T extends string = string> {
  value: T
  label: string
}

/** Genre dropdown options for the Generate Modal */
export const GENRE_OPTIONS: DropdownOption<Genre>[] = [
  { value: 'fantasy', label: 'Fantasy' },
  { value: 'sci-fi', label: 'Sci-Fi' },
  { value: 'post-apocalyptic', label: 'Post-Apocalyptic' },
  { value: 'steampunk', label: 'Steampunk' },
  { value: 'horror', label: 'Horror' },
]

/** Asset type dropdown options for the Generate Modal */
export const ASSET_TYPE_OPTIONS: DropdownOption<AssetType>[] = [
  { value: 'character', label: 'Character' },
  { value: 'location', label: 'Location' },
  { value: 'item', label: 'Item' },
  { value: 'creature', label: 'Creature' },
]

/** Mood dropdown options for the Generate Modal */
export const MOOD_OPTIONS: DropdownOption<Mood>[] = [
  { value: 'dark', label: 'Dark' },
  { value: 'epic', label: 'Epic' },
  { value: 'mysterious', label: 'Mysterious' },
  { value: 'whimsical', label: 'Whimsical' },
]

/**
 * Words that are not allowed in prompts (case-insensitive check).
 * @see Requirement 7.3
 */
export const BLOCKED_WORDS: string[] = ['nude', 'explicit', 'nsfw', 'gore']
