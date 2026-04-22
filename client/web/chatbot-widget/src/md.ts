// Tiny safe markdown renderer for agent messages. Intentionally minimal:
// bold, italic, inline + fenced code, safe links, bullet lists, newlines.
// Everything is escaped first — no user/agent input reaches innerHTML raw.

const SAFE_URL = /^(https?:|mailto:)/i

function escapeHtml(s: string): string {
  return s
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

export function renderMarkdown(src: string): string {
  if (!src) return ''
  let out = escapeHtml(src)

  // fenced code — captured across lines; language tag ignored
  out = out.replace(/```([a-zA-Z0-9_-]*)\n([\s\S]*?)```/g, (_m, _lang, body) => {
    return `<pre><code>${body.replace(/\n$/, '')}</code></pre>`
  })

  // inline code
  out = out.replace(/`([^`\n]+)`/g, '<code>$1</code>')

  // bold, then italic (bold first so ** doesn't get eaten by *)
  out = out.replace(/\*\*([^*\n]+)\*\*/g, '<strong>$1</strong>')
  out = out.replace(/(^|[^*])\*([^*\n]+)\*/g, '$1<em>$2</em>')
  out = out.replace(/(^|[^_])_([^_\n]+)_/g, '$1<em>$2</em>')

  // links — allowlist http/https/mailto only
  out = out.replace(/\[([^\]]+)\]\(([^)]+)\)/g, (m, label, url) => {
    if (!SAFE_URL.test(url)) return m
    return `<a href="${url}" target="_blank" rel="noopener noreferrer">${label}</a>`
  })

  // bullet lists — group consecutive "- " lines
  out = out.replace(/(?:^|\n)((?:- [^\n]+\n?)+)/g, block => {
    const items = block
      .trim()
      .split('\n')
      .map(l => l.replace(/^- /, ''))
      .map(l => `<li>${l}</li>`)
      .join('')
    return `\n<ul>${items}</ul>\n`
  })

  // paragraph breaks: double newline → block break; single → <br>
  // but don't add <br> inside <pre>/<ul>
  out = out
    .replace(/\n{2,}/g, '\n\n')
    .split('\n\n')
    .map(block => {
      if (/^\s*(<pre|<ul|<ol)/.test(block)) return block
      return block.replace(/\n/g, '<br>')
    })
    .join('\n\n')

  return out
}
