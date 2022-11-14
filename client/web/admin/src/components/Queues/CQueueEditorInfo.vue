<template>
  <b-card
    class="shadow-sm"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-form @submit="$emit('submit', queue)">
      <b-form-group
        :label="$t('name')"
        label-cols="2"
      >
        <b-form-input
          v-model="queue.queue"
          :state="handleState"
        />
        <b-form-invalid-feedback :state="handleState">
          {{ $t('invalid-handle-characters') }}
        </b-form-invalid-feedback>
      </b-form-group>

      <b-form-group
        :label="$t('consumer')"
        label-cols="2"
      >
        <b-form-select
          v-model="queue.consumer"
          :options="consumers"
        />
      </b-form-group>

      <b-form-group
        label-cols="2"
        :label="$t('poll_delay')"
        :description="metaPollDelayDescription()"
      >
        <b-form-input
          v-model="(queue.meta || {}).poll_delay"
          class="col-xs-2 col-lg-2"
          :state="durationState"
        />
      </b-form-group>

      <b-form-group
        v-if="isMetaDispatchEvents"
        :label="$t('dispatch_events')"
        :description="$t('dispatch_events_desc')"
        label-cols="2"
      >
        <b-form-checkbox
          v-model="queue.meta.dispatch_events"
          name="checkbox-1"
        >
          {{ $t("dispatch_events") }}
        </b-form-checkbox>
      </b-form-group>

      <b-form-group
        v-if="queue.createdAt"
        :label="$t('createdAt')"
        label-cols="2"
      >
        {{ queue.createdAt }}
      </b-form-group>

      <b-form-group
        v-if="queue.updatedAt"
        :label="$t('updatedAt')"
        label-cols="2"
      >
        {{ queue.updatedAt }}
      </b-form-group>

      <b-form-group
        v-if="queue.deletedAt"
        :label="$t('deletedAt')"
        label-cols="2"
      >
        {{ queue.deletedAt }}
      </b-form-group>
    </b-form>

    <template #header>
      <h3 class="m-0">
        {{ $t("title") }}
      </h3>
    </template>

    <template #footer>
      <c-submit-button
        class="float-right"
        :processing="processing"
        :success="success"
        :disabled="saveDisabled"
        @submit="$emit('submit', queue)"
      />

      <confirmation-toggle
        v-if="queue && queue.queueID"
        @confirmed="$emit('delete')"
      >
        {{ getDeleteStatus }}
      </confirmation-toggle>
    </template>
  </b-card>
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'
import { handleState } from 'corteza-webapp-admin/src/lib/handle'
import ConfirmationToggle from 'corteza-webapp-admin/src/components/ConfirmationToggle'
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'

export default {
  name: 'CQueueEditorInfo',

  i18nOptions: {
    namespaces: 'system.queues',
    keyPrefix: 'editor.info',
  },

  components: {
    ConfirmationToggle,
    CSubmitButton,
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
      return queue ? handleState(queue) : false
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
