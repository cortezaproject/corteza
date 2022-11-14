<template>
  <b-card no-body class="editor rt-content">
    <template v-if="editor">
      <b-card-header header-class="p-0 rounded-sm">
        <editor-menu-bar
          :editor="editor"
          v-slot="{ commands, isActive, getMarkAttrs, getNodeAttrs }">

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
          class="editor__content"
          :editor="editor" />

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
      default: () => ({})
    }
  },

  data () {
    const formats = getFormats()
    return {
      formats,
      toolbar: getToolbar(),
      // Helper to determine if current content differes from prop's content
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
        content: this.value,
        onUpdate: this.onUpdate,
      })

      /**
       * Since we migrated to TipTap, the new content should be emitted
       * after tiptap is done parsing it.
       */
      this.$nextTick(() => {
        this.onUpdate()
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
