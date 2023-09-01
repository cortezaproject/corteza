<template>
  <div>
    <b-form-group
      :label-cols="3"
      :label="$t('calendar.reminderFeed.colorLabel')"
      horizontal
      breakpoint="md"
      label-class="text-primary"
    >
      <c-input-color-picker
        v-model="feed.options.color"
        :translations="{
          modalTitle: $t('calendar.recordFeed.colorPicker'),
          cancelBtnLabel: $t('general:label.cancel'),
          saveBtnLabel: $t('general:label.saveAndClose')
        }"
      />
    </b-form-group>
  </div>
</template>

<script>
import base from './base'
import { components } from '@cortezaproject/corteza-vue'
const { CInputColorPicker } = components

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    CInputColorPicker,
  },

  extends: base,
}
</script>
