<template>
  <b-card>
    <c-progress
      :value="progress.completed"
      :max="progress.entryCount"
      labeled
      progress
      :animated="!progress.finishedAt"
      :relative="false"
      :variant="progress.startupError ? 'danger' : 'success'"
      text-style="font-size: 1.5rem;"
      style="height: 4rem;"
      class="mb-4"
    />

    <div
      v-if="!progress.finishedAt"
      class="d-flex"
    >
      <span class="text-secondary">
        <b-spinner
          variant="secondary"
          small
        />
        {{ $t('recordList.import.importing') }}
      </span>

      <b-button
        variant="light"
        class="ml-auto"
        @click="$emit('close')"
      >
        {{ $t('general:label.cancel') }}
      </b-button>
    </div>

    <div
      v-else-if="progress.startupError"
      class="d-flex"
    >
      <span class="text-danger">
        {{ progress.startupError }}
      </span>

      <b-button
        variant="light"
        class="ml-auto"
        @click="$emit('close')"
      >
        {{ $t('general:label.close') }}
      </b-button>
    </div>

    <div
      v-else-if="!progress.failed"
      class="d-flex"
    >
      <span class="text-success">
        {{ $t('recordList.import.success') }}
      </span>

      <b-button
        variant="light"
        class="ml-auto"
        @click="$emit('close')"
      >
        {{ $t('general:label.close') }}
      </b-button>
    </div>
  </b-card>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'
const { CProgress } = components

let toHandle = null

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    CProgress,
  },

  props: {
    session: {
      type: Object,
      required: true,
      default: () => ({}),
    },

    noPool: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    progress () {
      return this.session.progress || {}
    },
  },

  watch: {
    progress: {
      handler ({ finishedAt, failed, startupError }) {
        if (!finishedAt) return

        // Initial run failed; we already render an error state in the
        // template. Routing to ErrorReport here would crash since there's
        // no failLog to display.
        if (startupError) {
          this.clearTimeout()
          return
        }

        if (failed) {
          this.clearTimeout()
          this.$emit('importFailed', this.progress)
          return
        }

        this.clearTimeout()
        this.$root.$emit('recordList.refresh', this.session)
        this.$emit('importSuccessful')
      },
    },
  },

  mounted () {
    // If the initial run already failed before we mounted, don't start polling.
    if (this.progress.startupError) return

    if (!this.noPool) {
      this.pool()
    }
  },

  beforeDestroy () {
    this.clearTimeout()
  },

  methods: {
    clearTimeout () {
      if (toHandle !== null) {
        window.clearTimeout(toHandle)
        toHandle = null
      }
    },

    pool () {
      this.$ComposeAPI.recordImportProgress(this.session)
        .then(({ progress }) => {
          this.$set(this.session, 'progress', progress)
          if (!progress || !progress.finishedAt) {
            toHandle = window.setTimeout(this.pool, 2000)
          }
        })
        .catch((err) => {
          // Progress endpoint blew up (session expired, server down, …).
          // Surface it as a startup error so the spinner stops.
          const message = (err && err.message) || this.$t('recordList.import.startFailed')
          this.$set(this.session, 'progress', {
            ...this.progress,
            finishedAt: new Date().toISOString(),
            startupError: message,
          })
        })
    },
  },
}
</script>

<style lang="scss" scoped>
.progress-label {
  font-size: 15px;
}

</style>
