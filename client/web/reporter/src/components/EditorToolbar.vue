<template>
  <c-toolbar
    class="bg-white shadow border-top"
  >
    <template #start>
      <b-button
        v-if="backLink"
        data-test-id="button-back"
        variant="link"
        :to="backLink"
        class="d-flex align-items-center text-dark back text-nowrap gap-1 text-decoration-none"
      >
        <font-awesome-icon
          :icon="['fas', 'chevron-left']"
          class="back-icon"
        />
        {{ $t('general:label.back') }}
      </b-button>
    </template>

    <template #center>
      <slot />
    </template>

    <template #end>
      <c-input-confirm
        v-if="!hideDelete"
        size="lg"
        size-confirm="lg"
        variant="danger"
        :disabled="deleteDisabled || processingDelete || processing"
        :processing="processingDelete"
        :text="$t('general:label.delete')"
        :borderless="false"
        @confirmed="$emit('delete')"
      />

      <c-button-submit
        data-test-id="button-clone"
        :disabled="cloneDisabled || processingClone || processing"
        :processing="processingClone"
        variant="light"
        :text="$t('general:label.clone')"
        class="text-nowrap"
        size="lg"
        @submit="$emit('clone')"
      />

      <c-button-submit
        data-test-id="button-save"
        :disabled="saveDisabled || processingSave || processing"
        :processing="processingSave"
        :text="$t('general:label.save')"
        size="lg"
        @submit="$emit('save')"
      />
    </template>
  </c-toolbar>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'
const { CToolbar } = components

export default {
  components: {
    CToolbar,
  },

  props: {
    backLink: {
      type: Object,
      required: false,
      default: () => ({ name: 'root' }),
    },

    hideDelete: {
      type: Boolean,
    },

    deleteDisabled: {
      type: Boolean,
    },

    hideSave: {
      type: Boolean,
    },

    saveDisabled: {
      type: Boolean,
    },

    cloneDisabled: {
      type: Boolean,
    },

    processing: {
      type: Boolean,
    },

    processingDelete: {
      type: Boolean,
    },

    processingSave: {
      type: Boolean,
      required: false,
    },

    processingClone: {
      type: Boolean,
    },
  },
}
</script>
