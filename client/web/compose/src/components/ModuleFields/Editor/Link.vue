<template>
  <b-form-group
    :label-cols-md="horizontal && '5'"
    :label-cols-xl="horizontal && '4'"
    :content-cols-md="horizontal && '7'"
    :content-cols-xl="horizontal && '8'"
    :class="formGroupStyleClasses"
  >
    <template
      #label
    >
      <div
        v-if="!valueOnly"
        class="d-flex align-items-center text-primary px-0"
      >
        <span
          :title="label"
          class="d-inline-block mw-100"
        >
          {{ label }}
        </span>

        <c-hint :tooltip="hint" />

        <slot name="tools" />
      </div>
      <div
        class="small text-muted"
        :class="{ 'mb-1': description }"
      >
        {{ description }}
      </div>
    </template>

    <multi
      v-if="field.isMulti"
      v-slot="ctx"
      :value.sync="value"
      :errors="errors"
    >
      <b-form-input
        :value="value[ctx.index]"
        :placeholder="placeholder"
        :formatter="formatValue"
        lazy-formatter
        @input="setMultiValue($event, ctx.index)"
      />
    </multi>

    <template
      v-else
    >
      <b-form-input
        v-model="value"
        :placeholder="placeholder"
        :formatter="formatValue"
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

  computed: {
    placeholder () {
      return this.$t('kind.link.example')
    },

    effectiveType () {
      if (this.field.options.dynamicMode) {
        return this.detectType(this.value)
      } else {
        const protocol = this.field.options.tempProtocol === '__CUSTOM__'
          ? this.field.options.customProtocol
          : this.field.options.tempProtocol
        return this.getTypeFromProtocol(protocol)
      }
    },
  },

  methods: {
    detectType (input) {
      if (!input) return 'url'

      const trimmed = input.trim()

      if (trimmed.includes('@')) {
        return 'email'
      }

      if (trimmed.startsWith('+') || /^[\d\s\-()]+$/.test(trimmed)) {
        return 'phone'
      }

      return 'url'
    },

    getTypeFromProtocol (protocol) {
      if (!protocol) return 'url'

      const protocolToType = {
        'mailto:': 'email',
        'tel:': 'phone',
        'sms:': 'phone',
        'http://': 'url',
        'https://': 'url',
      }

      const type = protocolToType[protocol]

      if (type) {
        return type
      }

      if (this.field.options.customProtocol && protocol.startsWith(this.field.options.customProtocol)) {
        return 'custom'
      }

      return 'app'
    },

    formatValue (value) {
      if (!value) return value

      const type = this.effectiveType
      let formatted = value

      switch (type) {
        case 'email':
          formatted = formatted.toLowerCase()
          break
        case 'phone':
          formatted = formatted.replace(/[^\d+]/g, '')
          break
        case 'url':
          if (this.field.options.trimFragment) {
            formatted = trimUrlFragment(formatted)
          }
          if (this.field.options.trimQuery) {
            formatted = trimUrlQuery(formatted)
          }
          if (this.field.options.trimPath) {
            formatted = trimUrlPath(formatted)
          }
          if (this.field.options.onlySecure) {
            formatted = onlySecureUrl(formatted)
          }
          break
        case 'app':
          break
        case 'custom':
          if (this.field.options.tempProtocol === '__CUSTOM__') {
            if (this.field.options.customProtocol && !formatted.startsWith(this.field.options.customProtocol)) {
              formatted = this.field.options.customProtocol + formatted
            }
          } else if (this.field.options.tempProtocol && this.field.options.tempProtocol !== '__CUSTOM__') {
            if (!formatted.startsWith(this.field.options.tempProtocol)) {
              formatted = this.field.options.tempProtocol + formatted
            }
          }
          break
      }

      return formatted
    },
  },
}
</script>
