<template>
  <c-ace-editor
    v-model="template.template"
    data-test-id="template-html-editor"
    name="editor/html"
    lang="html"
    height="400px"
    show-line-numbers
    highlight-active-line
    show-print-margin
    :border="false"
  />
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'

const { CAceEditor } = components

export default {
  components: {
    CAceEditor,
  },

  props: {
    template: {
      type: Object,
      required: true,
      default: () => ({}),
    },
  },
}
</script>
