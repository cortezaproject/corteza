import MarkdownIt from 'markdown-it'
import DOMPurify from 'dompurify'
// @ts-expect-error — no type defs published
import { full as emoji } from 'markdown-it-emoji'
// @ts-expect-error — no type defs published
import footnote from 'markdown-it-footnote'
// @ts-expect-error — no type defs published
import mark from 'markdown-it-mark'
// @ts-expect-error — no type defs published
import deflist from 'markdown-it-deflist'
// @ts-expect-error — no type defs published
import sub from 'markdown-it-sub'
// @ts-expect-error — no type defs published
import sup from 'markdown-it-sup'
// @ts-expect-error — no type defs published
import abbr from 'markdown-it-abbr'
// @ts-expect-error — no type defs published
import taskLists from 'markdown-it-task-lists'

const md = new MarkdownIt({
  html: true,
  breaks: true,
  linkify: true,
  typographer: false,
})
  .use(emoji)
  .use(footnote)
  .use(mark)
  .use(deflist)
  .use(sub)
  .use(sup)
  .use(abbr)
  .use(taskLists, { enabled: true, label: false })

const defaultLinkRender =
  md.renderer.rules.link_open ||
  function (tokens, idx, options, _env, self) {
    return self.renderToken(tokens, idx, options)
  }

md.renderer.rules.link_open = function (tokens, idx, options, env, self) {
  const token = tokens[idx]
  const targetIdx = token.attrIndex('target')
  if (targetIdx < 0) token.attrPush(['target', '_blank'])
  else token.attrs![targetIdx][1] = '_blank'

  const relIdx = token.attrIndex('rel')
  if (relIdx < 0) token.attrPush(['rel', 'noopener noreferrer'])
  else token.attrs![relIdx][1] = 'noopener noreferrer'

  return defaultLinkRender(tokens, idx, options, env, self)
}

export function renderMarkdown(source: string): string {
  if (!source) return ''
  const html = md.render(source)
  return DOMPurify.sanitize(html, {
    ADD_ATTR: ['target', 'rel'],
  })
}
