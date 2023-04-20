<template>
  <b-container
    fluid
    :class="{ 'shadow border-top': !showRecordModal }"
    class="bg-white p-3"
  >
    <b-row
      align-v="stretch"
      no-gutters
      class="wrap-with-vertical-gutters"
    >
      <b-col
        class="d-flex align-items-center justify-content-start"
      >
        <b-button
          v-if="!(hideBack || settings.hideBack)"
          data-test-id="button-back"
          variant="link"
          :disabled="processing"
          class="text-dark back"
          @click.prevent="$emit('back')"
        >
          <font-awesome-icon

            :icon="['fas', showRecordModal && !inEditing ? 'times' : 'chevron-left']"
            class="back-icon"
          />
          {{ backLabel }}
        </b-button>

        <slot name="start-actions" />
      </b-col>

      <b-col
        class="d-flex align-items-center justify-content-center"
      >
        <b-button-group v-if="recordNavigation.prev || recordNavigation.next">
          <b-button
            pill
            size="lg"
            variant="outline-primary"
            class="mr-2"
            :disabled="!record || processing || !recordNavigation.prev"
            :title="$t('recordNavigation.prev')"
            @click="navigateToRecord(recordNavigation.prev)"
          >
            <font-awesome-icon :icon="['fas', 'angle-left']" />
          </b-button>

          <b-button
            size="lg"
            pill
            variant="outline-primary"
            :disabled="!record || processing || !recordNavigation.next"
            :title="$t('recordNavigation.next')"
            @click="navigateToRecord(recordNavigation.next)"
          >
            <font-awesome-icon :icon="['fas', 'angle-right']" />
          </b-button>
        </b-button-group>

        <slot name="center-actions" />
      </b-col>

      <b-col
        class="d-flex align-items-center justify-content-end"
      >
        <template
          v-if="module"
        >
          <c-input-confirm
            v-if="isCreated && !(isDeleted || hideDelete || settings.hideDelete)"
            :disabled="!canDeleteRecord"
            size="lg"
            size-confirm="lg"
            variant="danger"
            :borderless="false"
            @confirmed="$emit('delete')"
          >
            <b-spinner
              v-if="processingDelete"
              small
              type="grow"
            />

            <span v-else>
              {{ labels.delete || $t('label.delete') }}
            </span>
          </c-input-confirm>

          <c-input-confirm
            v-if="isDeleted"
            :disabled="!canUndeleteRecord"
            size="lg"
            size-confirm="lg"
            variant="warning"
            variant-ok="warning"
            :borderless="false"
            @confirmed="$emit('undelete')"
          >
            <b-spinner
              v-if="processingUndelete"
              small
              type="grow"
            />

            <span v-else>
              {{ $t('label.restore') }}
            </span>
          </c-input-confirm>

          <b-button
            v-if="!inEditing && isCreated && module.canCreateRecord && !(hideClone || settings.hideClone)"
            data-test-id="button-clone"
            variant="light"
            size="lg"
            :disabled="!record || processing"
            class="ml-2"
            @click.prevent="$emit('clone')"
          >
            {{ labels.clone || $t('label.saveAsCopy') }}
          </b-button>

          <b-button
            v-if="!inEditing && isCreated && !(hideEdit || settings.hideEdit)"
            data-test-id="button-edit"
            :disabled="!record.canUpdateRecord || processing"
            variant="light"
            size="lg"
            class="ml-2"
            @click.prevent="$emit('edit')"
          >
            {{ labels.edit || $t('label.edit') }}
          </b-button>

          <b-button
            v-if="!inEditing && module.canCreateRecord && !(hideNew || settings.hideNew)"
            data-test-id="button-add-new"
            variant="primary"
            size="lg"
            :disabled="processing"
            class="ml-2"
            @click.prevent="$emit('add')"
          >
            {{ labels.new || $t('label.addNew') }}
          </b-button>

          <b-button
            v-if="inEditing && !(hideSubmit || settings.hideSubmit)"
            data-test-id="button-submit"
            :disabled="!canSaveRecord || processing"
            class="d-flex align-items-center justify-content-center ml-2"
            variant="primary"
            size="lg"
            @click.prevent="$emit('submit')"
          >
            <b-spinner
              v-if="processingSubmit"
              small
              type="grow"
            />

            <span
              v-else
              data-test-id="button-save"
            >
              {{ labels.submit || $t('label.save') }}
            </span>
          </b-button>

          <slot name="end-actions" />
        </template>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import { compose, NoID } from '@cortezaproject/corteza-js'
import { throttle } from 'lodash'

export default {
  i18nOptions: {
    namespaces: 'general',
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
      return this.record && this.record.recordID !== NoID
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
      if (this.inEditing) {
        return this.$t('label.cancel')
      } else if (this.showRecordModal) {
        return this.$t('label.close')
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
<style lang="scss" scoped>
.back {
  &:hover {
    text-decoration: none;

    .back-icon {
      transition: transform 0.3s ease-out;
      transform: translateX(-4px);
    }
  }
}
</style>
