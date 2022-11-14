import { get } from 'lodash'
import colorschemes from './colorschemes'

export const getColorschemeColors = (colorscheme?: string): string[] => {
  if (!colorscheme) {
    return []
  }

  return [...get(colorschemes, colorscheme)]
}
