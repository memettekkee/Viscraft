import { api } from '../../lib/api'

export interface PromptOption {
  id: string
  category: string
  label: string
  promptValue: string
  sortOrder: number
}

interface PromptOptionsResponse {
  success: boolean
  data: PromptOption[]
}

/**
 * Fetches ALL prompt options in a single request.
 * Frontend maps by category client-side.
 */
export async function fetchAllPromptOptions(): Promise<PromptOption[]> {
  const response = await api.post<PromptOptionsResponse>('/prompt-options', {})
  return response.data.data ?? []
}

/** Groups prompt options by category */
export function groupByCategory(options: PromptOption[]): Record<string, PromptOption[]> {
  const grouped: Record<string, PromptOption[]> = {}
  for (const opt of options) {
    if (!grouped[opt.category]) {
      grouped[opt.category] = []
    }
    grouped[opt.category].push(opt)
  }
  return grouped
}
