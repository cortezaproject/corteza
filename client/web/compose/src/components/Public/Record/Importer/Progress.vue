<template>
  <b-card>
    <c-progress
      :value="progress.completed"
      :max="progress.entryCount"
      labeled
      progress
      :animated="!progress.finishedAt"
      :relative="false"
      variant="success"
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
      v-if="progress.finishedAt && !progress.failed"
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

  data () {
    return {
      progress: this.session.progress || {},
    }
  },

  watch: {
    progress: {
      handler ({ finishedAt, failed }) {
        if (finishedAt && failed) {
          this.clearTimeout()
          this.$emit('importFailed', this.progress)
        } else if (finishedAt) {
          this.clearTimeout()
          this.$root.$emit('recordList.refresh', this.session)
          this.$emit('importSuccessful')
        }
      },
    },
  },

  mounted () {
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
          this.progress = progress
          toHandle = window.setTimeout(this.pool, 2000)
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
