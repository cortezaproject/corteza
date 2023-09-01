<template>
  <b-card no-body>
    <slot name="title" />

    <fieldset>
      <b-form-group
        :label="$t('metric.editStyle.color')"
        label-class="text-primary"
      >
        <c-input-color-picker
          v-model="options.color"
          :translations="{
            modalTitle: $t('metric.editStyle.colorPicker'),
            cancelBtnLabel: $t('general:label.cancel'),
            saveBtnLabel: $t('general:label.saveAndClose')
          }"
          class="mb-1"
        />
      </b-form-group>

      <b-form-group
        :label="$t('metric.editStyle.backgroundColor')"
        label-class="text-primary"
      >
        <c-input-color-picker
          v-model="options.backgroundColor"
          :translations="{
            modalTitle: $t('geometry.recordFeed.colorPicker'),
            cancelBtnLabel: $t('general:label.cancel'),
            saveBtnLabel: $t('general:label.saveAndClose')
          }"
          class="mb-1"
        />
      </b-form-group>

      <b-form-group
        :label="$t('metric.editStyle.fontSize')"
        label-class="text-primary"
      >
        <b-form-input
          v-model="options.fontSize"
          type="number"
          placeholder="16"
          min="0.1"
          step="0.1"
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
