import { get } from 'lodash'
import colorschemes from './colorschemes'

export const getColorschemeColors = (colorscheme?: string, customColorSchemes?: any[]): string[] => {
  if (!colorscheme) {
    return ['#37A2DA', '#32C5E9', '#67E0E3', '#9FE6B8', '#FFDB5C', '#ff9f7f', '#fb7293', '#E062AE', '#E690D1', '#e7bcf3', '#9d96f5', '#8378EA', '#96BFFF']
  }

  if (colorscheme.includes('custom') && customColorSchemes) {
    return customColorSchemes.find(({ id }) => id === colorscheme)?.colors || []
  }

  return get(colorschemes, colorscheme)
}
