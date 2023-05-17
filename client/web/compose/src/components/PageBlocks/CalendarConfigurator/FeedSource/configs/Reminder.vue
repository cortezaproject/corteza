<template>
  <div>
    <b-form-group
      horizontal
      :label-cols="3"
      breakpoint="md"
      :label="$t('calendar.reminderFeed.colorLabel')"
    >
      <c-input-color-picker
        v-model="feed.options.color"
      />
    </b-form-group>
  </div>
</template>

<script>
import base from './base'
import CInputColorPicker from '/Users/katrinyordanova/Desktop/repos/corteza/client/web/compose/src/components/ModuleFields/Configurator/CInputColorPicker.vue'

// import { components } from '@cortezaproject/corteza-vue'
// const { CInputColorPicker } = components

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
