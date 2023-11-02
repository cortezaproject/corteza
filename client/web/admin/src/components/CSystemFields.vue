<template>
  <b-row
    class="py-4"
  >
    <b-col
      v-if="id"
      cols="12"
    >
      <b-form-group
        :label="$t('id')"
        label-class="text-primary"
      >
        {{ id }}
      </b-form-group>
    </b-col>

    <b-col
      v-for="(f, i) in systemFields"
      :key="i"
      cols="12"
    >
      <b-form-group
        :label="$t(f) || $t(label)"
        label-class="text-primary"
        :data-test-id="`input-${generateTestID(f)}`"
      >
        {{ getFieldValue(f) }}
      </b-form-group>
    </b-col>
    <slot name="custom-field" />
  </b-row>
</template>

<script>
import { getSystemFields, kebabize } from 'corteza-webapp-admin/src/lib/sysFields'
export default {
  name: 'CSystemFields',

  props: {
    resource: {
      type: Object,
      required: true,
    },

    label: {
      type: String,
      default: '',
    },

    id: {
      type: String,
      default: '',
    },
  },

  computed: {
    systemFields () {
      return getSystemFields(this.resource)
    },
  },

  methods: {
    generateTestID (field) {
      return kebabize(field)
    },

    getFieldValue (field) {
      const isTimeValue = field.substring(field.length - 2) === 'At'
      let value = isTimeValue ? this.$options.filters.locFullDateTime(this.resource[field]) : this.resource[field]
      return value
    },
  },
}
</script>

<style>

</style>
