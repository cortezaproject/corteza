<template>
  <div class="expression-editor">
    <Textarea
      v-model="expressionValue"
      :auto-resize="true"
      rows="3"
      class="expression-textarea w-full font-mono text-sm"
      spellcheck="false"
    />
  </div>
</template>

<script>
export default {
  props: {
    value: {
      type: String,
      default: '',
    },

    modelValue: {
      type: String,
      default: undefined,
    },

    lang: {
      type: String,
      default: 'text',
    },

    minHeight: {
      type: String,
      default: '6rem',
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
      default: true,
    },

    autoComplete: {
      type: Boolean,
      default: true,
    },
  },

  emits: ['update:value', 'update:modelValue', 'input', 'open'],

  computed: {
    expressionValue: {
      get () {
        // Support both v-model:value and v-model (modelValue)
        return this.modelValue !== undefined ? this.modelValue : this.value
      },

      set (val = '') {
        this.$emit('update:value', val)
        this.$emit('update:modelValue', val)
        this.$emit('input', val)
      },
    },
  },
}
</script>

<style scoped>
.expression-editor {
  width: 100%;
}

.expression-textarea {
  font-family: 'JetBrains Mono', 'Fira Code', 'Consolas', 'Monaco', monospace;
  font-size: 13px;
  line-height: 1.5;
  tab-size: 2;
  resize: vertical;
  min-height: v-bind(minHeight);
}
</style>
