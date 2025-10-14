<template>
  <div>
    <b-form-checkbox v-model="f.options.multiLine">
      {{ $t('kind.string.multiLine') }}
    </b-form-checkbox>

    <b-form-checkbox v-model="f.options.useRichTextEditor">
      {{ $t('kind.string.richText') }}
    </b-form-checkbox>

    <b-form-group
      label-class="d-flex align-items-center text-primary"
      class="mt-3"
    >
      <template #label>
        {{ $t('kind.string.sanitizeXSS.label') }}
        <c-hint
          :tooltip="$t('kind.string.sanitizeXSS.warning')"
          icon-class="text-warning"
        />
      </template>

      <c-input-checkbox
        v-model="f.options.sanitizeXSS"
        :switch="true"
        :labels="checkboxLabel"
      />
    </b-form-group>
  </div>
</template>

<script>
import base from './base'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  extends: base,

  data () {
    return {
      checkboxLabel: {
        on: this.$t('general:label.yes'),
        off: this.$t('general:label.no'),
      },
    }
  },
}
</script>
