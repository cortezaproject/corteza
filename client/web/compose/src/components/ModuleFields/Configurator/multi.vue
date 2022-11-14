<template>
  <div>
    <b-form-group
      v-if="field.isMulti"
      :label="$t('options.multiDelimiter.label')"
    >
      <b-form-radio-group
        v-model="field.options.multiDelimiter"
        :options="selectOptions"
        stacked
      />
      <b-form-group
        class="mt-2"
        :label="$t('options.multiDelimiter.customLabel')"
      >
        <b-form-input
          v-model="field.options.multiDelimiter"
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

  data () {
    return {
      selectOptions: [
        { text: this.$t('options.multiDelimiter.newline'), value: '\n' },
        { text: this.$t('options.multiDelimiter.comma'), value: ', ' },
      ],
    }
  },
}
</script>
