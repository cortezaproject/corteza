<template>
  <b-row align-v="center">
    <b-col
      cols="12"
      sm="4"
    >
      <b-form-group :label="$t('navigation.text')">
        <b-form-input
          v-model="column.options.itemOption.text"
          type="text"
          size="sm"
        />
      </b-form-group>
    </b-col>
    <b-col
      cols="12"
      sm="4"
    >
      <b-form-group :label="$t('navigation.url')">
        <b-form-input
          v-model="column.options.itemOption.url"
          type="text"
          size="sm"
        />
      </b-form-group>
    </b-col>
    <b-col
      cols="12"
      sm="4"
    >
      <b-form-group class="mt-4">
        <b-form-checkbox
          v-model="column.options.newWindow"
          size="sm"
          switch
        >
          {{ $t('navigation.newWindow') }}
        </b-form-checkbox>
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
