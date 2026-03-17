<template>
  <div :class="classes">
    <span
      v-for="(v, index) of formattedValue"
      :key="index"
      :class="{ 'd-block': field.options.multiDelimiter === '\n' }"
    >
      <span v-if="field.options.outputPlain || disableClick">
        {{ formatValue(v) }}{{ index !== formattedValue.length - 1 ? field.options.multiDelimiter : '' }}
      </span>

      <template v-else>
        <template v-if="isDynamic && detectType(v) === 'email'">
          <a
            :href="'mailto:' + v"
            class="text-decoration-none"
          >{{ formatValue(v) }}{{ index !== formattedValue.length - 1 ? field.options.multiDelimiter : '' }}</a>
        </template>
        <template v-else-if="isDynamic && detectType(v) === 'phone'">
          <a
            :href="'tel:' + v"
            class="text-decoration-none"
          >{{ formatValue(v) }}{{ index !== formattedValue.length - 1 ? field.options.multiDelimiter : '' }}</a>
        </template>
        <template v-else-if="isDynamic && detectType(v) === 'url'">
          <a
            :href="ensureProtocol(formatValue(v))"
            class="text-decoration-none"
            target="_blank"
            rel="noopener"
          >{{ formatValue(v) }}{{ index !== formattedValue.length - 1 ? field.options.multiDelimiter : '' }}</a>
        </template>
        <template v-else>
          <template v-if="field.options.tempProtocol || field.options.customProtocol">
            <a
              :href="getStaticLinkHref(v)"
              class="text-decoration-none"
              target="_blank"
              rel="noopener"
            >{{ formatValue(v) }}{{ index !== formattedValue.length - 1 ? field.options.multiDelimiter : '' }}</a>
          </template>
          <template v-else>
            <a
              v-if="v"
              :href="v"
              class="text-decoration-none"
              target="_blank"
              rel="noopener"
            >{{ v }}{{ index !== formattedValue.length - 1 ? field.options.multiDelimiter : '' }}</a>
            <span v-else>{{ v }}</span>
          </template>
        </template>
      </template>
    </span>
  </div>
</template>

<script>
import base from './base'
import { trimUrlFragment, trimUrlQuery, trimUrlPath, onlySecureUrl } from '../url'

export default {
  extends: base,

  computed: {
    isDynamic () {
      return this.field.options.dynamicMode !== false
    },

    formattedValue () {
      return this.field.isMulti ? this.value : [this.value].filter(v => v)
    },
  },

  methods: {
    getStaticLinkHref (value) {
      let protocol = this.field.options.tempProtocol
      if (protocol === '__CUSTOM__') {
        protocol = this.field.options.customProtocol
      }

      if (!protocol) {
        return value
      }

      if (value && value.startsWith(protocol)) {
        return value
      }

      return protocol + value
    },

    formatValue (value) {
      if (this.isDynamic) {
        const type = this.detectType(value)
        if (type === 'url') {
          return this.fixUrl(value)
        }
      } else {
        const protocol = this.field.options.tempProtocol === '__CUSTOM__'
          ? this.field.options.customProtocol
          : this.field.options.tempProtocol

        if (protocol === 'https://' || protocol === 'http://') {
          return this.fixUrl(value)
        }
      }
      return value
    },

    fixUrl (value) {
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

    ensureProtocol (url) {
      if (!url) return url

      if (/^[a-zA-Z][a-zA-Z0-9+.-]*:/.test(url)) {
        return url
      }

      return 'https://' + url
    },
  },
}
</script>
