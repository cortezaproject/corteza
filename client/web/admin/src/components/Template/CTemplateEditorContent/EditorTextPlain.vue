<template>
  <c-ace-editor
    v-model="template.template"
    data-test-id="template-plain-text-editor"
    name="editor/plain_text"
    lang="plain_text"
    height="400px"
    font-size="14px"
    show-line-numbers
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
