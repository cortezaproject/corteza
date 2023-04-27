<template>
  <b-container
    v-if="queue"
    class="py-3"
  >
    <c-content-header
      :title="title"
    >
      <span
        class="text-nowrap"
      >
        <b-button
          v-if="queueID && canCreate"
          data-test-id="button-add"
          variant="primary"
          class="mr-2"
          :to="{ name: 'system.queue.new' }"
        >
          {{ $t('new') }}
        </b-button>
      </span>
    </c-content-header>

    <c-queue-editor-info
      :queue="queue"
      :processing="info.processing"
      :success="info.success"
      :can-create="canCreate"
      :consumers="consumers"
      @delete="onDelete"
      @submit="onSubmit"
    />
  </b-container>
</template>

<script>
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CQueueEditorInfo from 'corteza-webapp-admin/src/components/Queues/CQueueEditorInfo'
import { mapGetters } from 'vuex'

export default {
  components: {
    CQueueEditorInfo,
  },

  i18nOptions: {
    namespaces: 'system.queues',
    keyPrefix: 'editor',
  },

  mixins: [
    editorHelpers,
  ],

  props: {
    queueID: {
      type: String,
      required: false,
      default: undefined,
    },
  },

  data () {
    return {
      queue: undefined,

      consumers: [],

      info: {
        processing: false,
        success: false,
      },
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canCreate () {
      return this.can('system/', 'queue.create')
    },

    title () {
      return this.queue.queueID ? this.$t('title.edit') : this.$t('title.new')
    },
  },

  watch: {
    queueID: {
      immediate: true,
      handler () {
        this.fetchQueueConsumers()

        if (this.queueID) {
          this.fetchQueue()
        } else {
          this.queue = {
            consumer: 'corteza',
            meta: {
              poll_delay: '',
              dispatch_events: false,
            },
          }
        }
      },
    },
  },

  methods: {
    fetchQueue () {
      this.incLoader()

      this.$SystemAPI.queuesRead({ queueID: this.queueID })
        .then(q => { this.queue = q })
        .catch(this.toastErrorHandler(this.$t('notification:queue.fetch.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    fetchQueueConsumers () {
      this.incLoader()

      this.consumers = [
        { value: 'store', text: 'Store' },
        { value: 'eventbus', text: 'Eventbus' },
        { value: 'corteza', text: 'Corteza' },
        { value: 'redis', text: 'Redis' },
      ]

      this.decLoader()
    },

    onSubmit (queue) {
      this.incLoader()
      if (this.queueID) {
        this.$SystemAPI.queuesUpdate(queue)
          .then(queue => {
            this.queue = queue

            this.animateSuccess('info')
            this.toastSuccess(this.$t('notification:queue.update.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:queue.update.error')))
          .finally(() => {
            this.decLoader()
          })
      } else {
        this.$SystemAPI.queuesCreate(queue)
          .then(({ queueID }) => {
            this.animateSuccess('info')
            this.toastSuccess(this.$t('notification:queue.create.success'))

            this.$router.push({ name: 'system.queue.edit', params: { queueID } })
          })
          .catch(this.toastErrorHandler(this.$t('notification:queue.create.error')))
          .finally(() => {
            this.decLoader()
          })
      }
    },

    onDelete () {
      this.incLoader()
      const { deletedAt = '' } = this.queue
      const method = deletedAt ? 'queuesUndelete' : 'queuesDelete'
      const event = deletedAt ? 'undelete' : 'delete'

      this.$SystemAPI[method]({ queueID: this.queueID })
        .then(() => {
          this.fetchQueue()
          this.toastSuccess(this.$t(`notification:queue.${event}.success`))
          if (!deletedAt) this.$router.push({ name: 'system.queue' })
        })
        .catch(this.toastErrorHandler(this.$t(`notification:queue.${event}.error`)))
        .finally(() => {
          this.decLoader()
        })
    },
  },
}
</script>
