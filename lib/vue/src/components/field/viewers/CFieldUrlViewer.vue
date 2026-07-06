<template>
  <div>
    <span v-for="(v, index) in formattedValues" :key="index" :class="{ block: isNewlineDelimiter }">
      <span v-if="field.options.outputPlain || disableClick">
        {{ fixUrl(v) }}{{ index !== formattedValues.length - 1 ? delimiter : '' }}
      </span>

      <a
        v-else
        :href="fixUrl(v)"
        target="_blank"
        rel="noopener noreferrer"
        class="text-primary font-medium hover:underline"
        @click.stop
      >
        {{ fixUrl(v) }}{{ index !== formattedValues.length - 1 ? delimiter : '' }}
      </a>
    </span>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { trimUrlFragment, trimUrlQuery, trimUrlPath, onlySecureUrl } from '../url'

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
  record: {
    type: Object,
    required: true,
  },
  namespace: {
    type: Object,
    default: () => ({}),
  },
  valueOnly: {
    type: Boolean,
    default: false,
  },
  extraOptions: {
    type: Object,
    default: () => ({}),
  },
  disableClick: {
    type: Boolean,
    default: false,
  },
})

const value = computed(() => {
  if (props.field.isSystem) {
    return props.record[props.field.name]
  }
  return props.record?.values?.[props.field.name]
})

const formattedValues = computed(() => {
  const v = value.value
  if (props.field.isMulti && Array.isArray(v)) {
    return v.filter(val => val)
  }
  return [v].filter(val => val)
})

const delimiter = computed(() => props.field.options?.multiDelimiter || ', ')
const isNewlineDelimiter = computed(() => delimiter.value === '\n')

function fixUrl(val) {
  let url = String(val)
  const opts = props.field.options || {}

  if (opts.trimFragment) url = trimUrlFragment(url)
  if (opts.trimQuery) url = trimUrlQuery(url)
  if (opts.trimPath) url = trimUrlPath(url)
  if (opts.onlySecure) url = onlySecureUrl(url)

  // Prepend https:// if no protocol present
  if (!/^https?:\/\//i.test(url)) {
    url = 'https://' + url
  }

  return url
}
</script>
