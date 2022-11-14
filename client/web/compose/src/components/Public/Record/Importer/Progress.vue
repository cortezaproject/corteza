<template>
  <b-card
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-form-group class="my-4 mx-4">
      <b-progress
        :max="progress.entryCount"
        show-value
        show-progress
        variant="primary"
        height="80px"
        class="bg-light"
      >
        <b-progress-bar
          :value="progress.completed"
          class="progress-label"
          variant="primary"
        >
          <span class="font-weight-bold">{{ $t('recordList.import.progressRatio', progress) }}</span>
        </b-progress-bar>
      </b-progress>
    </b-form-group>

    <b-form-group class="mx-4 mb-0">
      <span
        v-if="progress.finishedAt && !progress.failed"
        class="text-success"
      >

        {{ $t('recordList.import.success') }}
      </span>
    </b-form-group>
  </b-card>
</template>

<script>
let toHandle = null

export default {
  i18nOptions: {
    namespaces: 'block',
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
      handler: function ({ finishedAt, failed }) {
        if (finishedAt && failed) {
          this.clearTimeout()
          this.$emit('importFailed', this.progress)
        } else if (finishedAt) {
          this.clearTimeout()
          this.$root.$emit('recordList.refresh', this.session)
        }
      },
      deep: true,
      immediate: true,
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
