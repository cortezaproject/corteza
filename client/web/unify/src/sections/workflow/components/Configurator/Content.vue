<template>
  <div class="configurator-section">
    <CFormGroup :label="$t('steps.content.configurator.content-label')">
      <div class="content-rte-wrap">
        <c-rich-text-input
          v-model="label"
          :labels="{
            urlPlaceholder: $t('steps.content.configurator.urlPlaceholder'),
            ok: $t('steps.content.configurator.ok'),
          }"
          @input="$emit('update-value', $event)"
        />
      </div>
    </CFormGroup>
  </div>
</template>

<script>
import base from './base.vue'
import { components } from '@planetcrust/human-vue'
const { CRichTextInput } = components

export default {
  components: {
    CRichTextInput,
  },

  extends: base,

  computed: {
    label: {
      get () {
        return this.item.node.value
      },

      set (label) {
        this.item.node.value = label
      },
    },
  },
}
</script>

<style scoped>
.content-rte-wrap :deep(.ProseMirror),
.content-rte-wrap :deep(.ql-editor),
.content-rte-wrap :deep([contenteditable="true"]) {
  min-height: 200px;
}
</style>
