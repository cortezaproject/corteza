<template>
  <form @submit.prevent="handleSubmit">
    <div class="input-group">
      <textarea
        ref="textarea"
        v-model="input"
        :placeholder="inputPlaceholder"
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
        <font-awesome-icon
          :icon="['fas', 'paper-plane']"
        />
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
    inputPlaceholder: {
      type: String,
      default: '',
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
