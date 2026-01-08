<template>
  <div :class="['d-flex', isUser ? 'justify-content-end' : 'justify-content-start']">
    <div
      class="w-100"
      style="max-width: 48rem;"
    >
      <div
        :class="[
          'card',
          isUser ? 'bg-primary text-white' : ''
        ]"
      >
        <div class="card-body">
          <p
            class="card-text small mb-0"
            style="white-space: pre-wrap;"
          >
            {{ message.content }}
          </p>

          <div
            v-if="message.sources && message.sources.length > 0"
            class="mt-3 pt-3 border-top"
            :class="isUser ? 'border-white border-opacity-25' : ''"
          >
            <p
              class="small fw-semibold mb-2"
              :class="isUser ? 'text-white-50' : 'text-muted'"
            >
              Sources:
            </p>
            <div class="d-flex flex-wrap gap-2">
              <span
                v-for="(source, idx) in message.sources"
                :key="idx"
                class="badge bg-secondary"
              >
                <a
                  v-if="source.url"
                  :href="source.url"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="text-white text-decoration-none"
                >
                  {{ source.title }}
                </a>
                <span v-else>{{ source.title }}</span>
              </span>
            </div>
          </div>
        </div>
      </div>

      <p class="small text-muted mt-1 mb-0">
        {{ formatTime(message.timestamp) }}
      </p>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ChatMessage',

  props: {
    message: {
      type: Object,
      required: true,
    },
  },

  computed: {
    isUser () {
      return this.message.role === 'user'
    },
  },

  methods: {
    formatTime (timestamp) {
      return new Date(timestamp).toLocaleTimeString([], {
        hour: '2-digit',
        minute: '2-digit',
      })
    },
  },
}
</script>

<style scoped>
.gap-2 {
  gap: 0.5rem;
}
</style>
