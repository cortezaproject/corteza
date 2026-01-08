<template>
  <b-sidebar
    v-model="showRag"
    header-class="d-flex align-items-center justify-content-between notification-sidebar-header bg-white pl-3 pr-2"
    body-class="d-flex flex-column overflow-hidden bg-white"
    :backdrop="isMobile"
    no-footer
    right
    shadow="sm"
    no-close-on-route-change
    no-close-on-esc
    width="500px"
  >
    <template #header>
      <div>
        <h5 class="text-primary mb-0">
          <b>{{ $t("title") }}</b>
        </h5>
        <p class="small text-muted mb-0">
          <small> {{ $t("description") }} </small>
        </p>
      </div>

      <b-button
        variant="outline-light"
        class="d-flex align-items-center justify-content-center p-2 border-0 text-secondary"
        @click="showRag = false"
      >
        <font-awesome-icon :icon="['fas', 'times']" class="h6 mb-0" />
      </b-button>
    </template>

    <div class="border-bottom bg-white px-4 py-3 shadow-sm">
      <div class="d-flex align-items-center justify-content-between">
        <button
          v-if="messages.length > 0"
          class="btn btn-outline-secondary btn-sm"
          @click="handleClearChat"
        >
          {{ $t("clear-chat") }}
        </button>
      </div>
    </div>

    <div class="flex-fill overflow-auto bg-light">
      <div
        v-if="messages.length === 0"
        class="h-100 d-flex align-items-center justify-content-center"
      >
        <div class="text-center">
          <h2 class="h4 mb-2">{{ $t("welcome-title") }}</h2>
          <p class="text-muted">
            {{ $t("welcome-message") }}
          </p>
        </div>
      </div>

      <div v-else class="container py-4">
        <div class="row">
          <div class="col-12">
            <chat-message
              v-for="message in messages"
              :key="message.id"
              :message="message"
              class="mb-3"
            />

            <!-- Loading indicator -->
            <div v-if="isLoading" class="d-flex justify-content-center py-4">
              <div class="d-flex align-items-center gap-2">
                <div
                  class="spinner-grow spinner-grow-sm text-primary"
                  role="status"
                >
                  <span class="visually-hidden">Loading...</span>
                </div>
                <div
                  class="spinner-grow spinner-grow-sm text-primary"
                  role="status"
                  style="animation-delay: 0.1s"
                >
                  <span class="visually-hidden">Loading...</span>
                </div>
                <div
                  class="spinner-grow spinner-grow-sm text-primary"
                  role="status"
                  style="animation-delay: 0.2s"
                >
                  <span class="visually-hidden">Loading...</span>
                </div>
              </div>
            </div>

            <div ref="messagesEnd" />
          </div>
        </div>
      </div>
    </div>

    <div class="border-top bg-white px-4 py-3 shadow-sm">
      <div class="container">
        <div class="row">
          <div class="col-12">
            <chat-input
              :on-send-message="handleSendMessage"
              :is-loading="isLoading"
              :input-placeholder="$t('input-placeholder')"
            />
          </div>
        </div>
      </div>
    </div>
  </b-sidebar>
</template>

<script>
import ChatMessage from './rag/ChatMessage.vue'
import ChatInput from './rag/ChatInput.vue'

export default {
  name: 'Rag',

  i18nOptions: {
    namespaces: 'rag',
  },

  components: {
    ChatMessage,
    ChatInput,
  },

  props: {
    showRag: {
      type: Boolean,
      required: true,
      default: false,
    },
  },

  data () {
    return {
      messages: [],
      isLoading: false,
    }
  },

  computed: {
    isMobile () {
      return window.innerWidth < 576
    },
  },

  watch: {
    messages: {
      handler () {
        this.$nextTick(() => {
          this.scrollToBottom()
        })
      },
      deep: true,
    },
  },

  methods: {
    scrollToBottom () {
      if (this.$refs.messagesEnd) {
        this.$refs.messagesEnd.scrollIntoView({ behavior: 'smooth' })
      }
    },

    async handleSendMessage (content) {
      const userMessage = {
        id: Date.now().toString(),
        role: 'user',
        content,
        timestamp: new Date(),
      }

      this.messages.push(userMessage)
      this.isLoading = true

      try {
        const data = await this.$DiscoveryAPI.rag(content)

        const assistantMessage = {
          id: (Date.now() + 1).toString(),
          role: 'assistant',
          content: data.response || 'No response received',
          sources: data.sources || [],
          timestamp: new Date(),
        }

        this.messages.push(assistantMessage)
      } catch (error) {
        console.error('Error sending message:', error, content)
        const errorMessage = {
          id: (Date.now() + 1).toString(),
          role: 'assistant',
          content: 'Sorry, I encountered an error processing your request.',
          timestamp: new Date(),
        }
        this.messages.push(errorMessage)
      } finally {
        this.isLoading = false
      }
    },

    handleClearChat () {
      this.messages = []
    },
  },
}
</script>

<style scoped>
.gap-2 {
  gap: 0.5rem;
}
</style>
