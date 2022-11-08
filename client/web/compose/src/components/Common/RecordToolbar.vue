<template>
  <b-container
    fluid
    :class="{'shadow border-top': !showRecordModal}"
    class="bg-white p-3"
  >
    <b-row
      no-gutters
      class="wrap-with-vertical-gutters align-items-center"
    >
      <div
        class="wrap-with-vertical-gutters align-items-center"
      >
        <b-button
          v-if="!settings.hideBack"
          data-test-id="button-back"
          variant="link"
          class="text-dark back"
          :disabled="processing"
          @click.prevent="$emit('back')"
        >
          <font-awesome-icon
            :icon="['fas', showRecordModal ? 'times' : 'chevron-left']"
            class="back-icon"
          />
          {{ showRecordModal ? $t('label.close') : labels.back || $t('label.back') }}
        </b-button>
      </div>

      <div
        v-if="module"
        class="d-flex wrap-with-vertical-gutters align-items-center ml-auto"
      >
        <c-input-confirm
          v-if="isCreated && !settings.hideDelete"
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

        <b-button
          v-if="!inEditing && module.canCreateRecord && !hideClone && isCreated && !settings.hideClone"
          data-test-id="button-clone"
          variant="light"
          size="lg"
          :disabled="!record || processing"
          class="ml-2"
          @click.prevent="$emit('clone')"
        >
          {{ labels.clone || $t('label.clone') }}
        </b-button>

        <b-button
          v-if="!inEditing && !settings.hideEdit && isCreated"
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
          v-if="module.canCreateRecord && !hideAdd && !inEditing && !settings.hideNew"
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
          v-if="inEditing && !settings.hideSubmit"
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
      </div>
    </b-row>
  </b-container>
</template>

<script>
import { compose, NoID } from '@cortezaproject/corteza-js'

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

    processingSubmit: {
      type: Boolean,
      default: false,
    },

    inEditing: {
      type: Boolean,
      required: true,
    },

    hideClone: {
      type: Boolean,
      default: () => false,
    },

    hideAdd: {
      type: Boolean,
      default: () => false,
    },

    isDeleted: {
      type: Boolean,
      default: true,
    },

    showRecordModal: {
      type: Boolean,
      required: false,
    },
  },

  computed: {
    isCreated () {
      return this.record && this.record.recordID !== NoID
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
