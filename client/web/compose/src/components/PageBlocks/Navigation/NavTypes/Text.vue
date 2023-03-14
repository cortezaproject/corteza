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
  </tr>
</template>

<script>
import base from './base'

export default {
  extends: base,
}
</script>

<style lang="scss" scoped>
th,
td {
  padding-left: 15px;
  padding-right: 15px;
}
</style>
