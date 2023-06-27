<template>
  <b-form-group
    label-class="d-flex align-items-center text-primary p-0"
    :state="state"
    :class="formGroupStyleClasses"
  >
    <template
      v-if="!valueOnly"
      #label
    >
      <span class="d-inline-block text-truncate mw-100 py-1">
        {{ label }}
      </span>

      <hint
        :id="field.fieldID"
        :text="hint"
      />
    </template>

    <div
      class="small text-muted"
      :class="{ 'mb-1': description }"
    >
      {{ description }}
    </div>

    <multi
      v-if="field.isMulti"
      v-slot="ctx"
      :value.sync="value"
      :errors="errors"
    >
      <c-rich-text-input
        v-if="field.options.useRichTextEditor"
        v-model="value[ctx.index]"
        class="mr-2"
      />

      <b-form-textarea
        v-else-if="field.options.multiLine"
        v-model="value[ctx.index]"
        class="mr-2"
      />

      <b-form-input
        v-else
        v-model="value[ctx.index]"
        class="mr-2"
      />
    </multi>

    <template v-else>
      <c-rich-text-input
        v-if="field.options.useRichTextEditor"
        v-model="value"
        class="mr-2"
        :labels="{
          urlPlaceholder: $t('content.urlPlaceholder'),
          ok: $t('content.ok'),
          openLinkInNewTab: $t('content.openLinkInNewTab'),
        }"
      />

      <b-form-textarea
        v-else-if="field.options.multiLine"
        v-model="value"
        class="mr-2"
      />

      <b-form-input
        v-else
        v-model="value"
        class="mr-2"
      />

      <errors :errors="errors" />
    </template>
  </b-form-group>
</template>

<script>
import base from './base'
import { components } from '@cortezaproject/corteza-vue'
const { CRichTextInput } = components

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    CRichTextInput,
  },

  extends: base,
}
</script>
