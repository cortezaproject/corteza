<template>
  <ace-editor
    :font-size="14"
    :show-print-margin="true"
    :show-gutter="true"
    :highlight-active-line="true"
    width="100%"
    mode="plain_text"
    theme="chrome"
    name="editor/plain_text"
    :on-change="(v) => template.template = v"
    :value="template.template"
    :editor-props="{
      $blockScrolling: false,
    }"
  />
</template>

<script>
import { Ace as AceEditor } from 'vue2-brace-editor'

import 'brace/mode/plain_text'
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
