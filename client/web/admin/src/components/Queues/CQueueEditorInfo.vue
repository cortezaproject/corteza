<template>
  <b-card
    data-test-id="card-queue-edit"
    header-class="border-bottom"
    footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <template #header>
      <h4 class="m-0">
        {{ $t("title") }}
      </h4>
    </template>

    <b-form @submit="$emit('submit', queue)">
      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('name')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="queue.queue"
              data-test-id="input-name"
              :state="handleState"
            />
            <b-form-invalid-feedback
              :state="handleState"
              data-test-id="feedback-invalid-name"
            >
              {{ $t('invalid-handle-characters') }}
            </b-form-invalid-feedback>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('consumer')"
            label-class="text-primary"
          >
            <b-form-select
              v-model="queue.consumer"
              data-test-id="input-consumer"
              :options="consumers"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('poll_delay')"
            :description="metaPollDelayDescription()"
            label-class="text-primary"
          >
            <b-form-input
              v-model="(queue.meta || {}).poll_delay"
              data-test-id="input-polling"
              :state="durationState"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            v-if="isMetaDispatchEvents"
            :label="$t('dispatch_events')"
            :description="$t('dispatch_events_desc')"
            label-class="text-primary"
          >
            <b-form-checkbox
              v-model="queue.meta.dispatch_events"
              name="checkbox-1"
            >
              {{ $t("dispatch_events") }}
            </b-form-checkbox>
          </b-form-group>
        </b-col>
      </b-row>

      <c-system-fields
        :resource="queue"
      />
    </b-form>

    <template #footer>
      <c-input-confirm
        v-if="queue && queue.queueID && queue.canDeleteQueue"
        :data-test-id="deleteButtonStatusCypressId"
        variant="danger"
        size="md"
        @confirmed="$emit('delete')"
      >
        {{ getDeleteStatus }}
      </c-input-confirm>

      <c-button-submit
        :disabled="saveDisabled"
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="$emit('submit', queue)"
      />
    </template>
  </b-card>
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'
import { handle } from '@cortezaproject/corteza-vue'

export default {
  name: 'CQueueEditorInfo',

  i18nOptions: {
    namespaces: 'system.queues',
    keyPrefix: 'editor.info',
  },

  props: {
    consumers: {
      type: Array,
      required: true,
    },

    queue: {
      type: Object,
      required: true,
    },

    processing: {
      type: Boolean,
      value: false,
    },

    success: {
      type: Boolean,
      value: false,
    },

    canCreate: {
      type: Boolean,
      required: true,
    },
  },

  computed: {
    fresh () {
      return !this.queue.queueID || this.queue.queueID === NoID
    },

    editable () {
      return this.fresh ? this.canCreate : true // this.queue.canUpdateQueue
    },

    saveDisabled () {
      return !this.editable || [this.durationState, this.handleState].includes(false)
    },

    durationState () {
      let pd = (this.queue.meta || {}).poll_delay || ''
      let m = pd.match(/^((\d+h)?(\d+m)?(\d+s)?)|(\s)$/g)

      if (m.length && m[0] === pd) {
        return null
      }

      return false
    },

    handleState () {
      const { queue = '' } = this.queue
      return queue ? handle.handleState(queue) : false
    },

    isMetaPollDelay () {
      if (this.queue.queueID) {
        return ((this.queue.meta || {}).poll_delay || '') === ''
      }

      return true
    },

    isMetaDispatchEvents () {
      return ((this.queue || {}).meta || {}).dispatch_events === null
    },

    getDeleteStatus () {
      return this.queue.deletedAt ? this.$t('undelete') : this.$t('delete')
    },

    deleteButtonStatusCypressId () {
      return `button-${this.getDeleteStatus.toLowerCase()}`
    },
  },

  methods: {
    metaPollDelayDescription () {
      return (((this.queue || {}).meta || {}).poll_delay || null)
        ? this.$t('poll_delay_set')
        : this.$t('poll_delay_empty')
    },

  },
}
</script>
