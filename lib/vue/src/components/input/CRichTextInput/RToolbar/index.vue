<template>
  <div class="d-flex flex-wrap">
    <component
      :is="getItem(f)"
      v-for="(f, i) of formats"
      :key="`${f.name}${i}`"
      :format="f"
      v-bind="{ ...$props, ...f.props }"
      :is-active="isActive()"
      :get-mark-attrs="getMarkAttrs"
      :labels="labels"
      :current-value="currentValue"
      class="toolbar-item"
      @click="triggerCommand"
    />

    <!-- Extra button to remove formatting -->
    <b-button
      variant="link"
      class="toolbar-item text-dark font-weight-bold"
      @click="removeMarks"
    >
      <font-awesome-icon icon="remove-format" />
    </b-button>
  </div>
</template>

<script>
import cc from './loader'

export default {
  inheritAttrs: true,

  props: {
    editor: {
      type: Object,
      required: true,
    },
    formats: {
      type: Array,
      required: true,
      default: () => [],
    },
    labels: {
      type: Object,
      default: () => ({}),
    },
    currentValue: {
      type: String,
      default: '',
    },
  },

  methods: {
    isActive () {
      return new Proxy({}, {
        get: (_, prop) => (attrs) => this.editor.isActive(prop, attrs),
      })
    },

    /**
     * Helper method to determine what item type we should display.
     * It can be a simple button (bold, italic, ...) or a dropdown (alignment)
     * @param {Object} f Format in question
     * @returns {Component}
     */
    getItem (f) {
      let b
      if (f.mark) {
        b = cc.mark
      } else if (f.node) {
        b = cc.node
      } else if (f.nodeAttr) {
        b = cc.nodeAttr
      }

      if (!b) {
        throw new Error('invalid node type')
      }

      let comp
      if (f.component) {
        comp = b[f.component]
      } else {
        comp = b.Item
      }

      if (!comp) {
        throw new Error('invalid component type')
      }

      return comp
    },

    /**
     * Helper method for removing marks.
     * It will remove all marks from the current state's range.
     * @returns {Range}
     */
    removeMarks () {
      this.editor.chain().focus().unsetAllMarks().run()
    },

    triggerCommand (v) {
      const e = this.editor.chain().focus()
      const t = v.type
      const a = v.attrs || {}

      switch (t) {
        case 'bold': e.toggleBold().run(); break
        case 'italic': e.toggleItalic().run(); break
        case 'underline': e.toggleUnderline().run(); break
        case 'strike': e.toggleStrike().run(); break
        case 'color': e.setColor(a.color).run(); break
        case 'background': e.setBackgroundColor(a.color).run(); break
        case 'blockquote': e.toggleBlockquote().run(); break
        case 'codeBlock': e.toggleCodeBlock().run(); break
        case 'heading': e.setHeading(a).run(); break
        case 'paragraph': e.setParagraph().run(); break
        case 'orderedList': e.toggleOrderedList().run(); break
        case 'bulletList': e.toggleBulletList().run(); break
        case 'taskList': e.toggleTaskList().run(); break
        case 'horizontalRule': e.setHorizontalRule().run(); break
        case 'alignment': e.setTextAlign(a).run(); break
        case 'link': {
          if (!a.href) {
            e.unsetLink().run()
          } else {
            e.setLink(a).run()
          }
          break
        }

        // table actions
        case 'insertTable': e.insertTable(a).run(); break
        case 'addColumnBefore': e.addColumnBefore().run(); break
        case 'addColumnAfter': e.addColumnAfter().run(); break
        case 'deleteColumn': e.deleteColumn().run(); break
        case 'addRowBefore': e.addRowBefore().run(); break
        case 'addRowAfter': e.addRowAfter().run(); break
        case 'deleteRow': e.deleteRow().run(); break
        case 'mergeCells': e.mergeCells().run(); break
        case 'splitCell': e.splitCell().run(); break
        case 'toggleHeaderRow': e.toggleHeaderRow().run(); break
        case 'toggleHeaderCell': e.toggleHeaderCell().run(); break
        case 'toggleHeaderColumn': e.toggleHeaderColumn().run(); break
        case 'deleteTable': e.deleteTable().run(); break
        case 'emoji': e.setEmoji({ name: a.name }).run(); break
        default:
          break
      }
    },
  },
}
</script>

<style lang="scss">
.toolbar-item {
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.25rem !important;
  width: 2.25rem !important;
  height: 2.25rem !important;
  border-radius: 0.25rem !important;

  &:hover {
    background-color: var(--light) !important;
  }
}
</style>