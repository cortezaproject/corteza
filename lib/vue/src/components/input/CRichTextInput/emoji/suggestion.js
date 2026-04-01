import { createPopper } from '@popperjs/core'
import { VueRenderer } from '@tiptap/vue-2'
import EmojiList from './EmojiList.vue'

const EMOJI_BLACKLIST = new Set(['relaxed', 'frowning_face'])

export default {
  items: ({ query, editor }) => {
    const emojis = editor.storage.emoji?.emojis || []

    if (!query) {
      // Show face emojis by default when user just types ':'
      return emojis
        .filter(e => (e.group === '' || e.group === 'people & body') && !e.name.startsWith('regional_indicator') && !EMOJI_BLACKLIST.has(e.name))
        .slice(0, 20)
    }

    const q = query.toLowerCase()

    return emojis
      .filter((emoji) => {
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
    let component
    let popup
    let popperInstance

    return {
      onStart: (props) => {
        component = new VueRenderer(EmojiList, {
          parent: props.editor.view.dom.ownerDocument.defaultView?.Vue || undefined,
          propsData: props,
        })

        if (!props.clientRect) {
          return
        }

        const virtualElement = {
          getBoundingClientRect: props.clientRect,
        }

        popup = document.createElement('div')
        popup.appendChild(component.element)
        popup.style.position = 'absolute'
        popup.style.zIndex = '1100'
        popup.style.minWidth = '200px'
        popup.style.pointerEvents = 'auto'

        document.body.appendChild(popup)

        popperInstance = createPopper(virtualElement, popup, {
          placement: 'bottom-start',
          modifiers: [
            {
              name: 'offset',
              options: {
                offset: [0, 4],
              },
            },
            {
              name: 'preventOverflow',
              options: {
                boundary: 'viewport',
              },
            },
            {
              name: 'flip',
              options: {
                fallbackPlacements: ['top-start', 'bottom-start'],
              },
            },
          ],
        })
      },

      onUpdate(props) {
        component.updateProps(props)

        if (!props.clientRect || !popperInstance) {
          return
        }

        popperInstance.update()
      },

      onKeyDown(props) {
        if (props.event.key === 'Escape') {
          return true
        }

        return component.ref?.onKeyDown(props)
      },

      onExit() {
        if (popperInstance) {
          popperInstance.destroy()
        }
        if (popup) {
          popup.remove()
        }
        component.destroy()
      },
    }
  },
}
