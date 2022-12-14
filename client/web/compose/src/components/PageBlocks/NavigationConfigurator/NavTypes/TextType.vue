<template>
  <b-row>
    <b-col
      cols="12"
      sm="6"
    >
      <b-form-group :label="$t('navigation.text')">
        <b-form-input
          v-model="column.options.itemOption.text"
          type="text"
          size="sm"
        />
      </b-form-group>
    </b-col>
  </b-row>
</template>

<script>
export default {
  i18nOptions: {
    namespaces: 'block',
  },

  props: {
    column: {
      type: Object,
      required: true,
    },
  },
}
</script>
