import lodash from 'lodash-es'
const { get } = lodash

import colorschemes from './colorschemes'

// What a chart paints when it names no scheme.
const defaultColors = [
  '#37A2DA',
  '#32C5E9',
  '#67E0E3',
  '#9FE6B8',
  '#FFDB5C',
  '#ff9f7f',
  '#fb7293',
  '#E062AE',
  '#E690D1',
  '#e7bcf3',
  '#9d96f5',
  '#8378EA',
  '#96BFFF',
]

export const getColorschemeColors = (
  colorscheme?: string,
  customColorSchemes?: any[],
): string[] => {
  if (!colorscheme) {
    return [...defaultColors]
  }

  const colors = colorscheme.includes('custom')
    ? customColorSchemes?.find(({ id }) => id === colorscheme)?.colors
    : get(colorschemes, colorscheme)

  // An id that resolves to nothing — a custom scheme since deleted, a name that
  // never existed — reaches echarts as an undefined palette, so the chart draws
  // its legend and none of its series. The default palette says "no scheme"
  // where a blank panel says "no data".
  return colors?.length ? colors : [...defaultColors]
}
