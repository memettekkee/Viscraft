/**
 * Utility functions for reference image handling.
 * Converts files and URLs to raw base64 strings (without data URI prefix).
 *
 * Validates: Requirements 6.1, 6.2, 6.3, 6.4
 */

const ACCEPTED_TYPES = ['image/jpeg', 'image/png', 'image/webp']
const MAX_FILE_SIZE = 5 * 1024 * 1024 // 5MB

export interface FileValidationError {
  type: 'invalid-type' | 'too-large'
  message: string
}

/**
 * Validates a file before processing.
 * Returns null if valid, or a FileValidationError if invalid.
 */
export function validateImageFile(file: File): FileValidationError | null {
  if (!ACCEPTED_TYPES.includes(file.type)) {
    return {
      type: 'invalid-type',
      message: 'Only JPEG, PNG, and WEBP images are accepted',
    }
  }
  if (file.size > MAX_FILE_SIZE) {
    return {
      type: 'too-large',
      message: 'Image must be smaller than 5MB',
    }
  }
  return null
}

/**
 * Reads a File and returns raw base64 string (without the data URI prefix).
 *
 * Preconditions:
 *  - file is a valid File object
 *  - file.type is one of: image/jpeg, image/png, image/webp
 *  - file.size ≤ 5MB
 *
 * Postconditions:
 *  - Returns raw base64 string (no "data:image/...;base64," prefix)
 *  - Rejects with error if file cannot be read
 */
export function fileToBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => {
      const result = reader.result as string
      // Strip "data:image/...;base64," prefix, return raw base64
      const base64 = result.split(',')[1]
      resolve(base64)
    }
    reader.onerror = () => reject(new Error('Failed to read file'))
    reader.readAsDataURL(file)
  })
}

/**
 * Fetches an image from a URL and returns raw base64 string (without data URI prefix).
 *
 * Preconditions:
 *  - url is a valid relative path (e.g. /storage/images/uuid.png)
 *
 * Postconditions:
 *  - Returns raw base64 string of the fetched image
 *  - Rejects if fetch fails or image cannot be converted
 */
export async function imageUrlToBase64(url: string): Promise<string> {
  const baseUrl = window.__VISCRAFT_CONFIG__?.API_BASE_URL || 'http://localhost:8080'
  const fullUrl = `${baseUrl}${url}`
  const response = await fetch(fullUrl)

  if (!response.ok) {
    throw new Error('Failed to fetch image')
  }

  const blob = await response.blob()
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => {
      const result = reader.result as string
      const base64 = result.split(',')[1]
      resolve(base64)
    }
    reader.onerror = () => reject(new Error('Failed to convert image'))
    reader.readAsDataURL(blob)
  })
}
