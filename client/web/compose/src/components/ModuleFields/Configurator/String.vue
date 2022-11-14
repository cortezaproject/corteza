<template>
  <div>
    <b-form-checkbox v-model="f.options.multiLine">
      {{ $t('kind.string.multiLine') }}
    </b-form-checkbox>
    <b-form-checkbox v-model="f.options.useRichTextEditor">
      {{ $t('kind.string.richText') }}
    </b-form-checkbox>
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
