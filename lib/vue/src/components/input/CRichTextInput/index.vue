<template>
  <b-card
    header-class="p-0 border-bottom"
    body-class="p-0"
    class="border border-light rounded"
  >
    <template
      v-if="editor"
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
        content: this.value || ' ',
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
      this.currentValue = this.editor.getHTML()

      // Makes sure to default to '' as the value if no text is present, for validation purposes
      this.currentValue = this.currentValue !== '<p><br></p>' ? this.currentValue : ''

      this.emittedContent = true
      this.$emit('input', this.currentValue)
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