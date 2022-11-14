<template>
  <div>
    <div>
      <b-form-checkbox v-model="f.options.outputPlain">
        {{ $t('kind.email.preventToLink') }}
      </b-form-checkbox>
    </div>
  </div>
</template>

<script>
import base from './base'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  extends: base,
}
</script>
