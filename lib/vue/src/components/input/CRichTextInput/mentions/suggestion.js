import { createPopper } from '@popperjs/core'
import { VueRenderer } from '@tiptap/vue-2'
import MentionList from './MentionList.vue'

export default {
  items: ({ query, editor }) => {
    if (!query) {
      return []
    }

    let systemAPI = editor.options?.systemAPI

    if (!systemAPI) {
      return []
    }

    return systemAPI.userList({
      query,
      limit: 10,
    }).then(({ set }) => {
      return set.map(user => Object.freeze(user))
    }).catch(() => {
      return []
    })
  },

  render: () => {
    let component
    let popup
    let popperInstance

    return {
      onStart: (props) => {
        component = new VueRenderer(MentionList, {
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

        const handleSelection = (event, eventType) => {
          const options = popup.querySelectorAll('.mention-option')
          
          let optionClicked = false
          options.forEach((option, index) => {
            if (option.contains(event.target) || option === event.target) {
              optionClicked = true
              
              if (props.items[index]) {
                const item = props.items[index]
                
                props.command({
                  id: item.userID,
                  label: item.name || item.username || item.email || item.userID,
                })
                
                event.preventDefault()
                event.stopPropagation()
              }
            }
          })
        }

        popup.addEventListener('click', (event) => handleSelection(event, 'clicked'))
        popup.addEventListener('mouseup', (event) => handleSelection(event, 'mouseup'))

        // Track mousedown but don't prevent default to allow click events
        popup.addEventListener('mousedown', (event) => {
        })

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

        // Update popper position
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
