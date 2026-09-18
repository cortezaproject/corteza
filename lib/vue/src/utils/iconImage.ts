const LOAD_TIMEOUT_MS = 5_000

function loadImage(src: string): Promise<HTMLImageElement> {
  return new Promise((resolve, reject) => {
    const img = new Image()
    const timer = setTimeout(() => reject(new Error(`icon took too long: ${src}`)), LOAD_TIMEOUT_MS)
    img.crossOrigin = 'anonymous'
    img.onload = () => {
      clearTimeout(timer)
      resolve(img)
    }
    img.onerror = () => {
      clearTimeout(timer)
      reject(new Error(`icon did not load: ${src}`))
    }
    img.src = src
  })
}

// Draws the image at `src` into a square PNG data URL, contained and centred,
// with an optional red dot in the top-right corner. Rejects when the image does
// not load in time or its origin does not allow reading it back.
export async function renderIcon(
  src: string,
  { size = 64, dot = false, dotColor = '#ef4444' } = {},
): Promise<string> {
  const img = await loadImage(src)

  const canvas = document.createElement('canvas')
  canvas.width = size
  canvas.height = size
  const ctx = canvas.getContext('2d')
  if (!ctx) throw new Error('no 2d canvas context')

  // An SVG without width/height reports no natural size; it fills the square.
  const w = img.naturalWidth || size
  const h = img.naturalHeight || size
  const scale = Math.min(size / w, size / h)
  const dw = w * scale
  const dh = h * scale
  ctx.drawImage(img, (size - dw) / 2, (size - dh) / 2, dw, dh)

  if (dot) {
    const r = size * 0.22
    const cx = size - r
    const cy = r

    // A transparent ring keeps the dot legible over any icon.
    ctx.globalCompositeOperation = 'destination-out'
    ctx.beginPath()
    ctx.arc(cx, cy, r + size * 0.06, 0, Math.PI * 2)
    ctx.fill()

    ctx.globalCompositeOperation = 'source-over'
    ctx.fillStyle = dotColor
    ctx.beginPath()
    ctx.arc(cx, cy, r, 0, Math.PI * 2)
    ctx.fill()
  }

  return canvas.toDataURL('image/png')
}
