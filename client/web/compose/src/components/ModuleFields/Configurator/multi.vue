<template>
  <div>
    <b-form-group
      v-if="field.isMulti"
      :label="$t('options.multiDelimiter.label')"
      label-class="text-primary"
    >
      <b-form-radio-group
        v-model="field.options.multiDelimiter"
        :options="selectOptions"
        stacked
      />

      <b-form-group
        :label="$t('options.multiDelimiter.customLabel')"
        class="mt-2"
        label-class="text-primary"
      >
        <b-form-input
          v-model="field.options.multiDelimiter"
          :disabled="isNotConfigurable"
          :placeholder="$t('options.multiDelimiter.customPlaceholder')"
        />
      </b-form-group>
    </b-form-group>
  </div>
</template>

<script>
export default {
  i18nOptions: {
    namespaces: 'field',
  },

  props: {
    field: {
      type: Object,
      required: true,
    },
  },

  computed: {
    selectOptions () {
      return [
        { text: this.$t('options.multiDelimiter.newline'), value: '\n' },
        { text: this.$t('options.multiDelimiter.comma'), value: ', ', disabled: this.isNotConfigurable },
      ]
    },

    isNotConfigurable () {
      return ['File'].includes(this.field.kind)
    },
  },
}
</script>
