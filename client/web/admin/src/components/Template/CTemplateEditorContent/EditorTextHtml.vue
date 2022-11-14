<template>
  <ace-editor
    data-test-id="template-html-editor"
    :font-size="14"
    :show-print-margin="true"
    :show-gutter="true"
    :highlight-active-line="true"
    width="100%"
    mode="html"
    theme="chrome"
    name="editor/html"
    :on-change="(v) => template.template = v"
    :value="template.template"
    :editor-props="{
      $blockScrolling: false,
    }"
  />
</template>

<script>
import { Ace as AceEditor } from 'vue2-brace-editor'

import 'brace/mode/html'
import 'brace/theme/chrome'

export default {
  components: {
    AceEditor,
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
