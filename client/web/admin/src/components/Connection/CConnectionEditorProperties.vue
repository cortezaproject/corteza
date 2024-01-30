<template>
  <b-card
    header-class="border-bottom"
    footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <template #header>
      <h4 class="m-0">
        {{ $t('title') }}
      </h4>
    </template>

    <b-row
      v-for="prop in list"
      :key="prop"
    >
      <b-col cols="12">
        <b-form-checkbox
          v-model="properties[prop].enabled"
          class="mb-1"
        >
          {{ $t('form.' + kebabCase(prop) + '.checkbox.label') }}
        </b-form-checkbox>

        <b-form-group
          :label="$t('form.' + kebabCase(prop) + '.notes.label')"
          :description="$t('form.' + kebabCase(prop) + '.notes.description')"
          label-class="text-primary"
          class="ml-4"
        >
          <b-form-textarea
            v-model="properties[prop].notes"
          />
        </b-form-group>
      </b-col>
    </b-row>

    <template #footer>
      <c-button-submit
        :disabled="disabled"
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="$emit('submit')"
      />
    </template>
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
    properties: {
      type: Object,
      required: true,
    },

    disabled: {
      type: Boolean,
      default: false,
    },

    processing: {
      type: Boolean,
      value: false,
    },

    success: {
      type: Boolean,
      value: false,
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
