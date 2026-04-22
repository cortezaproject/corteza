<template>
  <div class="configurator-section">
    <div class="flex flex-col gap-1">
      <label class="font-medium text-primary">
        {{ $t('steps.content.configurator.content-label') }}
      </label>
      <div class="content-rte-wrap border border-surface rounded-border overflow-hidden">
        <c-rich-text-input
          v-model="label"
          :labels="{
            urlPlaceholder: $t('steps.content.configurator.urlPlaceholder'),
            ok: $t('steps.content.configurator.ok'),
          }"
          @input="$emit('update-value', $event)"
        />
      </div>
    </div>
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
