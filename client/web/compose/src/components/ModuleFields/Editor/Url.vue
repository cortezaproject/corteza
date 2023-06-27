<template>
  <b-form-group
    label-class="d-flex align-items-center text-primary p-0"
    :class="formGroupStyleClasses"
  >
    <template
      v-if="!valueOnly"
      #label
    >
      <span class="d-inline-block text-truncate mw-100 py-1">
        {{ label }}
      </span>

      <hint
        :id="field.fieldID"
        :text="hint"
      />
    </template>

    <div
      class="small text-muted"
      :class="{ 'mb-1': description }"
    >
      {{ description }}
    </div>

    <multi
      v-if="field.isMulti"
      v-slot="ctx"
      :value.sync="value"
      :errors="errors"
    >
      <b-form-input
        v-model="value[ctx.index]"
        type="url"
        class="mr-2"
        :placeholder="$t('kind.url.example')"
        :formatter="fixUrl"
        lazy-formatter
      />
    </multi>
    <template
      v-else
    >
      <b-form-input
        v-model="value"
        type="url"
        :placeholder="$t('kind.url.example')"
        :formatter="fixUrl"
        lazy-formatter
      />
      <errors :errors="errors" />
    </template>
  </b-form-group>
</template>
<script>
import base from './base'
import { trimUrlFragment, trimUrlQuery, trimUrlPath, onlySecureUrl } from '../url'

export default {
  i18nOptions: {
    namespaces: 'field',
  },
  extends: base,
  methods: {
    fixUrl (value) {
      // run through all the attributes
      if (this.field.options.trimFragment) {
        value = trimUrlFragment(value)
      }
      if (this.field.options.trimQuery) {
        value = trimUrlQuery(value)
      }
      if (this.field.options.trimPath) {
        value = trimUrlPath(value)
      }
      if (this.field.options.onlySecure) {
        value = onlySecureUrl(value)
      }

      return value
    },
  },
}
</script>
