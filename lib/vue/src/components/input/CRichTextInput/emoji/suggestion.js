import { VueRenderer } from '@tiptap/vue-3'
import EmojiList from './EmojiList.vue'

const EMOJI_BLACKLIST = new Set(['relaxed', 'frowning_face'])

export default {
  items: ({ query, editor }) => {
    const emojis = editor.storage.emoji?.emojis || []

    if (!query) {
      // Show face emojis by default when user just types ':'
      return emojis
        .filter(
          e =>
            (e.group === '' || e.group === 'people & body') &&
            !e.name.startsWith('regional_indicator') &&
            !EMOJI_BLACKLIST.has(e.name),
        )
        .slice(0, 20)
    }

    const q = query.toLowerCase()

    return emojis
      .filter(emoji => {
        // Skip regional indicators
        if (emoji.name.startsWith('regional_indicator')) return false
        if (EMOJI_BLACKLIST.has(emoji.name)) return false

        // Match against name
        if (emoji.name.toLowerCase().includes(q)) return true

        // Match against shortcodes
        if (emoji.shortcodes && emoji.shortcodes.some(s => s.toLowerCase().includes(q))) return true

        // Match against tags
        if (emoji.tags && emoji.tags.some(t => t.toLowerCase().includes(q))) return true

        return false
      })
      .slice(0, 20)
  },

  render: () => {
    let renderer

    return {
      onStart: props => {
        renderer = new VueRenderer(EmojiList, {
          props,
          editor: props.editor,
        })

        if (!props.clientRect) {
          return
        }

        const popup = renderer.element
        if (popup) {
          popup.style.position = 'absolute'
          popup.style.zIndex = '1100'
          popup.style.minWidth = '200px'
          popup.style.pointerEvents = 'auto'

          // Position using clientRect
          const rect = props.clientRect()
          if (rect) {
            popup.style.left = `${rect.left}px`
            popup.style.top = `${rect.bottom + 4}px`
          }

          document.body.appendChild(popup)
        }
      },

      onUpdate(props) {
        renderer.updateProps(props)

        if (!props.clientRect) {
          return
        }

        const popup = renderer.element
        if (popup) {
          const rect = props.clientRect()
          if (rect) {
            popup.style.left = `${rect.left}px`
            popup.style.top = `${rect.bottom + 4}px`
          }
        }
      },

      onKeyDown(props) {
        if (props.event.key === 'Escape') {
          return true
        }

        return renderer.ref?.onKeyDown(props)
      },

      onExit() {
        if (renderer?.element) {
          renderer.element.remove()
        }
        renderer?.destroy()
      },
    }
  },
}
