/**
 * Converts a 1-based index to a roman numeral string.
 *
 * Supports values from 1 to 3999.
 * Used by ProjectSidebar to prefix project names.
 *
 * Validates: Requirements 14.4
 */
export function toRomanNumeral(num: number): string {
  if (num < 1 || num > 3999) return String(num)

  const lookup: [number, string][] = [
    [1000, 'M'],
    [900, 'CM'],
    [500, 'D'],
    [400, 'CD'],
    [100, 'C'],
    [90, 'XC'],
    [50, 'L'],
    [40, 'XL'],
    [10, 'X'],
    [9, 'IX'],
    [5, 'V'],
    [4, 'IV'],
    [1, 'I'],
  ]

  let result = ''
  let remaining = num

  for (const [value, numeral] of lookup) {
    while (remaining >= value) {
      result += numeral
      remaining -= value
    }
  }

  return result
}
