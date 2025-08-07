<template>
  <div class="d-flex flex-wrap">
    <component
      :is="getItem(f)"
      v-for="(f, i) of formats"
      :key="`${f.name}${i}`"
      :format="f"
      v-bind="$props"
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
import { removeMark } from 'tiptap-commands'

export default {
  inheritAttrs: true,

  props: {
    editor: {
      type: Object,
      required: true,
    },
    commands: {
      type: Object,
      required: true,
    },
    isActive: {
      type: Object,
      required: true,
    },
    getMarkAttrs: {
      type: Function,
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
      removeMark(null)(this.editor.view.state, this.editor.view.dispatch)
    },

    triggerCommand (v) {
      this.commands[v.type](v.attrs)
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