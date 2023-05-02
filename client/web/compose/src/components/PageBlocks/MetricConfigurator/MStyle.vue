<template>
  <b-card no-body>
    <slot name="title" />

    <fieldset>
      <b-form-group
        :label="$t('metric.editStyle.color')"
      >
        <c-input-color-picker
          v-model="options.color"
          :translations="{
            modalTitle: $t('metric.editStyle.colorPicker'),
            saveBtnLabel: $t('general:label.saveAndClose')
          }"
          class="mb-1"
        />
      </b-form-group>

      <b-form-group
        :label="$t('metric.editStyle.backgroundColor')"
      >
        <c-input-color-picker
          v-model="options.backgroundColor"
          :translations="{
            modalTitle: $t('geometry.recordFeed.colorPicker'),
            saveBtnLabel: $t('general:label.saveAndClose')
          }"
          class="mb-1"
        />
      </b-form-group>

      <b-form-group
        :label="$t('metric.editStyle.fontSize')"
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
  },
}
</script>
