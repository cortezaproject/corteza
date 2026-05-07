const ICONS_WITHOUT_DARK_VARIANT = new Set([
  'play',
  'stop',
  'issue',
  'cog',
  'connection-point',
])

export function getIcon (name, theme = 'light') {
  if (!name) return ''
  const basePath = `${import.meta.env.BASE_URL}icons`
  const useDark = theme === 'dark' && !ICONS_WITHOUT_DARK_VARIANT.has(name)
  return `${basePath}/${useDark ? 'dark/' : ''}${name}.svg`
}
