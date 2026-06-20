import type { PromptOption } from '../../../service/promptOptions'

/**
 * Builds the final prompt string sent to the AI from user input + selected options.
 * The generated prompt is what Pollinations receives — not the raw user input.
 */
export function buildPrompt(userPrompt: string, selectedOptions: PromptOption[]): string {
  const safePrompt = userPrompt ?? ''
  const parts = [
    `Product: ${safePrompt.trim()}.`,
    ...selectedOptions.map((o) => o.promptValue + '.'),
    'Commercial product photography, high resolution, sharp focus on product, no watermark, no extra text.',
  ]
  return parts.filter(Boolean).join(' ')
}
