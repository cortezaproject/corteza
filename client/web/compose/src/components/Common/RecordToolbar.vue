<template>
  <c-toolbar :class="{ 'shadow border-top': !showRecordModal }">
    <template #start>
      <b-button
        v-if="!(hideBack || settings.hideBack)"
        data-test-id="button-back"
        variant="link"
        :disabled="processing"
        class="text-dark back text-left text-nowrap p-1"
        @click.prevent="$emit('back')"
      >
        <font-awesome-icon
          :icon="['fas', hasBack ? 'chevron-left' : 'times']"
          class="back-icon"
        />
        {{ backLabel }}
      </b-button>

      <slot name="start-actions" />
    </template>

    <template #center>
      <div
        v-if="recordNavigation.prev || recordNavigation.next"
        class="d-flex align-items-center fill-width gap-1"
      >
        <b-button
          v-b-tooltip.hover="{ title: $t('recordNavigation.prev'), container: '#body' }"
          pill
          size="lg"
          variant="outline-primary"
          :disabled="!record || processing || !recordNavigation.prev"
          @click="navigateToRecord(recordNavigation.prev)"
        >
          <font-awesome-icon :icon="['fas', 'angle-left']" />
        </b-button>

        <b-button
          v-b-tooltip.hover="{ title: $t('recordNavigation.next'), container: '#body' }"
          size="lg"
          pill
          variant="outline-primary"
          :disabled="!record || processing || !recordNavigation.next"
          @click="navigateToRecord(recordNavigation.next)"
        >
          <font-awesome-icon :icon="['fas', 'angle-right']" />
        </b-button>
      </div>

      <slot name="center-actions" />
    </template>

    <template
      v-if="module"
      #end
    >
      <slot name="end-actions" />

      <c-input-confirm
        v-if="isCreated && !(isDeleted || hideDelete || settings.hideDelete)"
        :disabled="!record || !canDeleteRecord"
        :processing="processingDelete"
        :text="labels.delete || $t('label.delete')"
        size="lg"
        size-confirm="lg"
        variant="danger"
        @confirmed="$emit('delete')"
      />

      <c-input-confirm
        v-if="isDeleted && !(hideDelete || settings.hideDelete)"
        :disabled="!record || !canUndeleteRecord"
        :processing="processingUndelete"
        :text="$t('label.restore')"
        size="lg"
        size-confirm="lg"
        variant="warning"
        variant-ok="warning"
        @confirmed="$emit('undelete')"
      />

      <b-button
        v-if="isCreated && module.canCreateRecord && !(hideClone || settings.hideClone)"
        data-test-id="button-clone"
        variant="light"
        size="lg"
        :disabled="!record || processing"
        class="text-nowrap"
        @click.prevent="$emit('clone')"
      >
        {{ labels.clone || $t('label.saveAsCopy') }}
      </b-button>

      <b-button
        v-if="!inEditing && isCreated && !(hideEdit || settings.hideEdit)"
        data-test-id="button-edit"
        :disabled="!record || !record.canUpdateRecord || processing"
        variant="light"
        size="lg"
        @click.prevent="$emit('edit')"
      >
        {{ labels.edit || $t('label.edit') }}
      </b-button>

      <b-button
        v-else-if="inEditing && isCreated && !(hideEdit || settings.hideEdit)"
        data-test-id="button-view"
        :disabled="!record || !record.canUpdateRecord || processing"
        variant="light"
        size="lg"
        @click.prevent="$emit('view')"
      >
        {{ labels.edit || $t('label.view') }}
      </b-button>

      <b-button
        v-if="!inEditing && module.canCreateRecord && !(hideNew || settings.hideNew)"
        data-test-id="button-add-new"
        variant="primary"
        size="lg"
        :disabled="!record || processing"
        class="text-nowrap"
        @click.prevent="$emit('add')"
      >
        {{ labels.new || $t('label.addNew') }}
      </b-button>

      <c-button-submit
        v-if="inEditing && !(hideSubmit || settings.hideSubmit)"
        data-test-id="button-save"
        :disabled="!record || !canSaveRecord || processingSubmit"
        :processing="processingSubmit"
        :text="labels.submit || $t('label.save')"
        size="lg"
        @submit="$emit('submit')"
      />
    </template>
  </c-toolbar>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'
import { compose, NoID } from '@cortezaproject/corteza-js'
import { throttle } from 'lodash'
const { CToolbar } = components

export default {
  i18nOptions: {
    namespaces: 'general',
  },

  components: {
    CToolbar,
  },

  props: {
    module: {
      type: compose.Module,
      required: false,
      default: undefined,
    },

    record: {
      type: compose.Record,
      required: false,
      default: undefined,
    },

    labels: {
      type: Object,
      default: () => ({}),
    },

    processing: {
      type: Boolean,
      default: false,
    },

    processingDelete: {
      type: Boolean,
      default: false,
    },

    processingUndelete: {
      type: Boolean,
      default: false,
    },

    processingSubmit: {
      type: Boolean,
      default: false,
    },

    inEditing: {
      type: Boolean,
      required: true,
    },

    hideBack: {
      type: Boolean,
      default: () => true,
    },

    hideDelete: {
      type: Boolean,
      default: () => true,
    },

    hideNew: {
      type: Boolean,
      default: () => true,
    },

    hideClone: {
      type: Boolean,
      default: () => true,
    },

    hideEdit: {
      type: Boolean,
      default: () => true,
    },

    hideSubmit: {
      type: Boolean,
      default: () => true,
    },

    showRecordModal: {
      type: Boolean,
      required: false,
    },

    recordNavigation: {
      type: Object,
      required: false,
      default: () => ({}),
    },

    hasBack: {
      type: Boolean,
      default: true,
    },
  },

  computed: {
    isCreated () {
      // The !this.record is intentional, to keep the button visible even when loading a record
      return !this.record || this.record.recordID !== NoID
    },

    isDeleted () {
      return this.record && this.record.deletedAt
    },

    settings () {
      return this.$Settings.get('compose.ui.record-toolbar', {})
    },

    canSaveRecord () {
      if (!this.module || !this.record) {
        return false
      }

      return this.record.recordID === NoID
        ? this.module.canCreateRecord
        : this.record.canUpdateRecord
    },

    canDeleteRecord () {
      if (!this.module || !this.record) {
        return false
      }

      return !this.isDeleted && this.record.canDeleteRecord && !this.processing && this.record.recordID !== NoID
    },

    canUndeleteRecord () {
      if (!this.module || !this.record) {
        return false
      }

      return this.isDeleted && this.record.canUndeleteRecord && !this.processing && this.record.recordID !== NoID
    },

    backLabel () {
      if (this.showRecordModal) {
        return this.hasBack ? this.$t('label.back') : this.$t('label.close')
      }

      return this.hasBack ? this.labels.back || this.$t('label.back') : this.$t('label.home')
    },
  },

  methods: {
    navigateToRecord: throttle(function (recordID) {
      this.$emit('update-navigation', recordID)
    }, 500),
  },
}
</script>
