<template>
  <b-card
    :header-class="!hideToolbar ? 'p-0 border-bottom' : 'd-none'"
    body-class="p-0"
    class="border border-light rounded"
  >
    <template
      v-if="editor && !hideToolbar"
      #header
    >
      <r-toolbar
        :editor="editor"
        :formats="toolbar"
        :labels="labels"
        :current-value="currentValue"
      />
    </template>

    <editor-content
      :editor="editor"
      :class="bodyClass"
      class="rt-editor-content rt-content p-2"
      :style="{ minHeight: minBodyHeight, maxHeight: maxBodyHeight }"
      @drop.native="onDrop"
      @paste.native="onPaste"
      @dragover.native.prevent
    />
  </b-card>
</template>

<script>
import { Editor, EditorContent } from '@tiptap/vue-2'
import RToolbar from './RToolbar/index.vue'
import { getFormats, getToolbar } from './lib'

export default {
  name: 'CRichTextInput',

  components: {
    EditorContent,
    RToolbar,
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

    minBodyHeight: {
      type: String,
      default: '10rem',
    },

    maxBodyHeight: {
      type: String,
      default: '',
    },

    bodyClass: {
      type: String,
      default: '',
    },

    placeholder: {
      type: String,
      default: '',
    },

    hideToolbar: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      formats: getFormats({
        placeholder: this.placeholder,
      }),
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
          this.editor.commands.setContent(val, false)
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
    if (this.editor) this.editor.destroy()
  },

  methods: {
    /**
     * Initialize the editor, state, ...
     */
    init () {
      this.editor = new Editor({
        extensions: this.formats,
        content: this.value || '',
        parseOptions: {
          preserveWhitespace: 'full',
        },
        onUpdate: this.onUpdate,
        systemAPI: this.$SystemAPI,
      })
    },

    /**
     * Handle on update events. Process current document & update data model
     * @note Currently, build-in toHTML function removes empty lines.
     * Because of this, we are using `view.dom.innerHTML`. This should be improved at a later point
     */
    onUpdate () {
      // Add <br> tags to empty paragraphs
      const editorValue = this.editor.getHTML().replace(/<p><\/p>/g, '<p><br></p>')

      this.currentValue = editorValue === '<p><br></p>' ? '' : editorValue

      this.emittedContent = true
      this.$emit('input', this.currentValue)
    },

    focus () {
      if (this.editor) {
        this.editor.commands.focus()
      }
    },

    onDrop (event) {
      if (event.dataTransfer && event.dataTransfer.files && event.dataTransfer.files.length > 0) {
        event.preventDefault()
        this.$emit('upload', event.dataTransfer.files)
      }
    },

    onPaste (event) {
      if (event.clipboardData && event.clipboardData.files && event.clipboardData.files.length > 0) {
        event.preventDefault()
        this.$emit('upload', event.clipboardData.files)
      }
    },
  },
}
</script>

<style lang="scss">
.rt-editor-content {
  height: 100%;

  .ProseMirror {
    height: 100%;
  }

  /* Make checkboxes only editable inside the editor */
  input[type="checkbox"] {
    pointer-events: auto !important;
    cursor: pointer !important;
  }
}
</style>
