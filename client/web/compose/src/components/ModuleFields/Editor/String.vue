<template>
  <b-form-group
    label-class="text-primary"
    :state="state"
    :class="formGroupStyleClasses"
  >
    <template
      v-if="!valueOnly"
      #label
    >
      <div
        class="d-flex align-items-top"
      >
        <label
          class="mb-0"
        >
          {{ label }}
        </label>

        <hint
          :id="field.fieldID"
          :text="hint"
        />
      </div>
      <small
        class="form-text font-weight-light text-muted"
      >
        {{ description }}
      </small>
    </template>

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
