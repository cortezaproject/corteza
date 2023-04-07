<template>
  <b-card no-body>
    <slot name="title" />

    <fieldset>
      <c-input-color-picker
        v-model="options.color"
        :label="$t('metric.editStyle.color')"
      />

      <c-input-color-picker
        v-model="options.backgroundColor"
        :label="$t('metric.editStyle.backgroundColor.label')"
        :disabled="isBackgroundTransparent"
        class="mb-1"
      />
      <small>
        {{ $t('metric.editStyle.backgroundColor.description') }}
      </small>

      <b-form-group
        :label="$t('metric.editStyle.fontSize')"
        label-cols="4"
        class="mt-1"
      >
        <b-form-input
          v-model="options.fontSize"
          type="number"
          placeholder="16"
          class="mb-1"
        />
      </b-form-group>
    </fieldset>
  </b-card>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'
const { CInputColorPicker } = components

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    CInputColorPicker,
  },

  props: {
    options: {
      type: Object,
      required: true,
      default: () => ({}),
    },

    isBackgroundTransparent: {
      type: Boolean,
    },
  },
}
</script>
