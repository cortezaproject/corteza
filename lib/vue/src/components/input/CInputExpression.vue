<template>
  <div
    class="position-relative"
  >
    <c-ace-editor
      v-model="editorValue"
      auto-complete
      init-expressions
      :auto-complete-suggestions="autoCompleteSuggestions"
      v-bind="{ ...$attrs, ...$props }"
      v-on="$listeners"
    />
  </div>
</template>

<script>
import CAceEditor from './CAceEditor.vue'

export default {
  components: {
    CAceEditor,
  },

  props: {
    value: {
      type: String,
      default: '',
    },

    lang: {
      type: String,
      default: 'text',
    },

    height: {
      type: String,
      default: '80px',
    },

    showLineNumbers: {
      type: Boolean,
      default: false,
    },

    fontSize: {
      type: String,
      default: '14px',
    },

    border: {
      type: Boolean,
      default: true,
    },

    showPopout: {
      type: Boolean,
      default: false,
    },

    readOnly: {
      type: Boolean,
      default: false,
    },

    highlightActiveLine: {
      type: Boolean,
      default: false,
    },

    showPrintMargin: {
      type: Boolean,
      default: false,
    },

    suggestionParams: {
      type: Array,
      default: () => [],
    },

    fontFamily: {
      type: String,
      default: '',
    },

    placeholder: {
      type: String,
      default: '',
    },
  },

  computed: {
    editorValue: {
      get () {
        return this.value
      },

      set (value = '') {
        this.$emit('update:value', value)
      },
    },

    autoCompleteSuggestions () {
      return this.getRecordBasedSuggestions(this.suggestionParams)
    },
  },

  methods: {
    getRecordBasedSuggestions (params = []) {
      const result = {}

      function addSuggestion (key, caption, value) {
        if (!result[key]) result[key] = []
        result[key].push({ caption: caption, value: value })
      }

      function processProperties (prefix, properties, interpolate) {
        (properties || []).forEach((prop) => {
          if (typeof prop === 'string') {
            const value = prefix + '.' + prop + (interpolate ? '}' : '')
            addSuggestion(prefix, prop, value)
          } else {
            const nestedPrefix = prefix + '.' + prop.value + '.'
            addSuggestion(prefix, prop.value, nestedPrefix)

            if (prop.properties) {
              (prop.properties || []).forEach((nestedProp) => {
                const nestedValue = nestedPrefix + nestedProp + (interpolate ? '}' : '')
                addSuggestion(prefix + '.' + prop.value, nestedProp, nestedValue)
              })
            }
          }
        })
      }

      (params || []).forEach((p) => {
        if (typeof p === 'string') {
          addSuggestion('', '', p)
        } else {
          const { interpolate = false, properties = [], value, root = true } = p
          const prefix = interpolate ? '${' : ''
          const suffix = interpolate && !properties.length ? '}' : ''
          const prefixAsValue = prefix + value + suffix + (properties.length > 0 ? '.' : '')

          if (root) {
            addSuggestion('', '', prefixAsValue)
          }

          if (interpolate) {
            addSuggestion('$', prefixAsValue.slice(1), prefixAsValue)
            addSuggestion('${', prefixAsValue.slice(2), prefixAsValue)
          }

          if (properties.length) {
            processProperties(prefix + value, properties, interpolate)
          }
        }
      })

      return result
    },
  },
}
</script>
