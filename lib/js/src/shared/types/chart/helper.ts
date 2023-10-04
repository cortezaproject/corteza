import { get } from 'lodash'
import colorschemes from './colorschemes'

export const getColorschemeColors = (colorscheme?: string, customColorSchemes?: any[]): string[] => {
  if (!colorscheme) {
    return []
  }

  if (colorscheme.includes('custom') && customColorSchemes) {
    return customColorSchemes.find(({ id }) => id === colorscheme)?.colors || []
  }

  return get(colorschemes, colorscheme)
}
