<template>
  <tr>
    <td />
    <td>
      <b-form-group :label="$t('navigation.fieldLabel')">
        <b-form-input
          v-model="options.item.label"
          type="text"
        />
      </b-form-group>
    </td>
    <td>
      <b-form-group :label="$t('navigation.url')">
        <b-form-input
          v-model="options.item.url"
          type="text"
        />
      </b-form-group>
    </td>
    <td>
      <b-form-group :label="$t('navigation.openIn')">
        <b-form-select
          v-model="options.item.target"
          :options="openInType"
        />
      </b-form-group>
    </td>
  </tr>
</template>

<script>
import base from './base'

export default {
  extends: base,

  data () {
    return {
      value: {},
      openInType: [
        { value: 'sameTab', text: this.$t('navigation.sameTab') },
        { value: 'newTab', text: this.$t('navigation.newTab') },
      ],
    }
  },
}
</script>
