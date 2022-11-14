export const rgbaRegex = /^rgba\((\d+),.*?(\d+),.*?(\d+),.*?(\d*\.?\d*)\)$/

const ln = (n) => Math.round(n < 0 ? 255 + n : (n > 255) ? n - 255 : n)
export const toRGBA = ([r, g, b, a]) =>
  `rgba(${ln(r)}, ${ln(g)}, ${ln(b)}, ${a})`

export const removeDup = (set, key) => {
  const exists = new Set()
  return set.filter(s => {
    const isN = !exists.has(s[key])
    exists.add(s[key])
    return isN
  })
}
