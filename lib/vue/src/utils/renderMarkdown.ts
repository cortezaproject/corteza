export function renderMarkdown(text: string): string {
  if (!text) return ''

  let html = text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')

  // Code blocks
  html = html.replace(/```([\s\S]*?)```/g, (_match, p1) => {
    return `<pre class="bg-surface-ground p-3 rounded-lg overflow-x-auto text-sm font-mono my-2 text-color border border-surface shadow-sm"><code>${p1.trim()}</code></pre>`
  })

  // Inline code
  html = html.replace(/`([^`]+)`/g, '<code class="bg-surface-ground px-1.5 py-0.5 rounded text-sm font-mono text-primary border border-surface shadow-sm">$1</code>')

  // Bold
  html = html.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')

  // Italic
  html = html.replace(/\*([^*]+)\*/g, '<em>$1</em>')

  // Headers
  html = html.replace(/^### (.*$)/gim, '<h3 class="text-base font-bold mt-3 mb-1 text-primary-emphasis">$1</h3>')
  html = html.replace(/^## (.*$)/gim, '<h2 class="text-lg font-bold mt-4 mb-2 pb-1 border-b border-surface text-primary-emphasis">$1</h2>')
  html = html.replace(/^# (.*$)/gim, '<h1 class="text-xl font-bold mt-5 mb-3 pb-1 border-b border-surface text-primary-emphasis">$1</h1>')

  // Unordered lists
  html = html.replace(/^\s*[-*+]\s+(.*)$/gim, '<li class="ml-4 list-disc marker:text-primary">$1</li>')

  // Ordered lists
  html = html.replace(/^\d+\.\s+(.*)$/gim, '<li class="ml-4 list-decimal marker:text-primary">$1</li>')

  // Consolidate adjacent list items
  html = html.replace(/(<li.*<\/li>)\n(<li.*<\/li>)/g, '$1$2')
  html = html.replace(/(<li.*<\/li>)/g, '<ul class="my-2 space-y-1 text-color">$1</ul>')
  html = html.replace(/<\/ul>\n<ul[^>]*>/g, '')

  // Line breaks
  html = html.replace(/\n\n/g, '</p><p class="my-2 text-color">')
  html = html.replace(/\n/g, '<br>')

  return `<div class="prose-sm max-w-none break-words">${html}</div>`
}
