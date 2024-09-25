<template>
  <b-card
    no-body
    class="editor rt-content"
  >
    <template v-if="editor">
      <b-card-header header-class="p-0 border-bottom">
        <editor-menu-bar
          v-slot="{ commands, isActive, getMarkAttrs, getNodeAttrs }"
          :editor="editor"
        >
          <r-toolbar
            :editor="editor"
            :formats="toolbar"
            :commands="commands"
            :is-active="isActive"
            :get-mark-attrs="getMarkAttrs"
            :get-node-attrs="getNodeAttrs"
            :labels="labels"
            :current-value="currentValue"
          />
        </editor-menu-bar>
      </b-card-header>

      <b-card-body>
        <editor-content
          :editor="editor"
          class="editor__content"
        />
      </b-card-body>
    </template>
  </b-card>
</template>

<script>
import RToolbar from './RToolbar/index.vue'
import { EditorMenuBar, Editor, EditorContent } from 'tiptap'
import { getFormats, getToolbar } from './lib'

export default {
  name: 'CRichTextInput',

  components: {
    EditorContent,
    RToolbar,
    EditorMenuBar,
  },

  props: {
    value: {
      type: String,
      required: false,
      default: null,
    },

    labels: {
      type: Object,
      default: () => ({}),
    },
  },

  data () {
    const formats = getFormats()
    return {
      formats,
      toolbar: getToolbar(),
      // Helper to determine if current content differs from prop's content
      emittedContent: false,
      editor: undefined,
      currentValue: '',
    }
  },

  watch: {
    value: {
      handler: function (val) {
        // Update happened due to external content change, not model change
        if (!this.emittedContent) {
          this.editor.setContent(val)
        }

        this.emittedContent = false
      },
      deep: true,
    },
  },

  mounted () {
    this.init()
  },

  beforeDestroy () {
    this.editor.destroy()
  },

  methods: {
    /**
     * Initialize the editor, state, ...
     */
    init () {
      this.editor = new Editor({
        extensions: this.formats,
        // Bypass Editor default empty white space script with an empty space string if there is no value because it's not really valid html
        // also ensuring that the unsaved changes alert detection is not triggered when the Editor does not have any changes
        content: this.value || ' ',
        parseOptions: {
          preserveWhitespace: 'full',
        },
        onUpdate: this.onUpdate,
      })
    },

    /**
     * Handle on update events. Process current document & update data model
     * @note Currently, build-in toHTML function removes empty lines.
     * Because of this, we are using `view.dom.innerHTML`. This should be improved at a later point
     */
    onUpdate () {
      this.currentValue = this.editor.view.dom.innerHTML

      // Makes sure to default to '' as the value if no text is present, for validation purposes
      this.currentValue = this.currentValue !== '<p><br></p>' ? this.currentValue : ''

      this.emittedContent = true
      this.$emit('input', this.currentValue)
    },
  },
}
</script>

<style>
/* Basic editor styles */
.rt-content {
  min-width: 12rem;
  position: static;
}

.rt-content :first-child {
  margin-top: 0;
}

/* Table-specific styling */
.rt-content  table {
  border-collapse: collapse;
  margin: 0;
  overflow: hidden;
  table-layout: fixed;
  width: 100%;
}

.rt-content  td,
.rt-content  th {
  border: 1px solid var(--gray);
  box-sizing: border-box;
  min-width: 1em;
  padding: 6px 8px;
  position: relative;
  vertical-align: top;
}

.rt-content  td > *,
.rt-content  th > * {
  margin-bottom: 0;
}

.rt-content  th {
  background-color: var(--gray-dark);
  font-weight: bold;
  text-align: left;
}

.rt-content  .selectedCell::after {
  background: var(--gray-dark);
  content: "";
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  pointer-events: none;
  position: absolute;
  z-index: 2;
}

.rt-content  .column-resize-handle {
  background-color: var(--purple);
  bottom: -2px;
  pointer-events: none;
  position: absolute;
  right: -2px;
  top: 0;
  width: 4px;
}

.rt-content  .tableWrapper {
  overflow-x: auto;
}

.rt-content .resize-cursor {
  cursor: ew-resize;
  cursor: col-resize;
}
</style>
