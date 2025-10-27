<template>
  <form @submit.prevent="handleSubmit">
    <div class="input-group">
      <textarea
        ref="textarea"
        v-model="input"
        placeholder="ask a question... (Shift+Enter for new line)"
        :disabled="isLoading"
        class="form-control"
        rows="1"
        style="resize: none; max-height: 120px;"
        @input="autoResize"
        @keydown="handleKeyDown"
      />
      <button
        type="submit"
        :disabled="!input.trim() || isLoading"
        class="btn btn-primary"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <line
            x1="22"
            y1="2"
            x2="11"
            y2="13"
          />
          <polygon points="22 2 15 22 11 13 2 9 22 2" />
        </svg>
      </button>
    </div>
  </form>
</template>

<script>
export default {
  name: 'ChatInput',
  props: {
    onSendMessage: {
      type: Function,
      required: true,
    },
    isLoading: {
      type: Boolean,
      default: false,
    },
  },
  data () {
    return {
      input: '',
    }
  },
  methods: {
    autoResize () {
      this.$nextTick(() => {
        const textarea = this.$refs.textarea
        if (textarea) {
          textarea.style.height = 'auto'
          textarea.style.height = `${Math.min(textarea.scrollHeight, 120)}px`
        }
      })
    },
    handleSubmit () {
      if (this.input.trim() && !this.isLoading) {
        this.onSendMessage(this.input)
        this.input = ''
        this.$nextTick(() => {
          if (this.$refs.textarea) {
            this.$refs.textarea.style.height = 'auto'
          }
        })
      }
    },
    handleKeyDown (e) {
      if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault()
        this.handleSubmit()
      }
    },
  },
}
</script>
