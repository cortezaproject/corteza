<template>
  <b-card
    class="shadow-sm"
    :title="$t('title')"
  >
    <b-row
      v-for="prop in list"
      :key="prop"
    >
      <b-col cols="12">
        <b-form-checkbox
          v-model="properties[prop].enabled"
          :disabled="disabled"
          class="mb-1"
        >
          {{ $t('form.' + kebabCase(prop) + '.checkbox.label') }}
        </b-form-checkbox>

        <b-form-group
          :label="$t('form.' + kebabCase(prop) + '.notes.label')"
          :description="$t('form.' + kebabCase(prop) + '.notes.description')"
          class="ml-4"
        >
          <b-form-textarea
            v-model="properties[prop].notes"
            :disabled="disabled"
          />
        </b-form-group>
      </b-col>
    </b-row>
  </b-card>
</template>
<script>
import { kebabCase } from 'lodash'

export default {
  i18nOptions: {
    namespaces: 'system.connections',
    keyPrefix: 'editor.properties',
  },

  props: {
    disabled: { type: Boolean, default: false },

    properties: {
      type: Object,
      required: true,
    },
  },

  data () {
    return {
      list: [
        'dataAtRestEncryption',
        'dataAtRestProtection',
        'dataAtTransitEncryption',
        'dataRestoration',
      ],
    }
  },

  methods: {
    kebabCase (str) {
      return kebabCase(str)
    },
  },
}
</script>
