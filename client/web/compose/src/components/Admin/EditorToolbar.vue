<template>
  <b-container
    fluid
    data-test-id="editor-toolbar"
    class="bg-white shadow border-top p-3"
  >
    <b-row
      no-gutters
      class="gap-1 gap-col-3 flex-column flex-md-row align-items-md-center"
    >
      <div
        class="d-flex gap-1 align-items-center"
      >
        <b-button
          data-test-id="button-back-without-save"
          variant="link"
          :disabled="processing"
          class="text-dark back mr-auto p-1"
          @click="$emit('back')"
        >
          <font-awesome-icon
            :icon="['fas', 'chevron-left']"
            class="back-icon"
          />
          {{ $t('label.backWithoutSave') }}
        </b-button>
      </div>

      <div
        class="ml-md-auto d-flex flex-column d-inline-block-md"
      >
        <slot />
      </div>

      <div
        class="d-flex gap-1 flex-md-row flex-column align-items-md-center ml-md-auto"
      >
        <slot name="delete" />

        <c-input-confirm
          v-if="!hideDelete"
          v-b-tooltip.hover
          :disabled="disableDelete || processing"
          :processing="processingDelete"
          :text="$t('label.delete')"
          size="lg"
          size-confirm="lg"
          variant="danger"
          :title="deleteTooltip"
          :borderless="false"
          class="d-flex flex-column"
          @confirmed="$emit('delete')"
        />

        <slot name="saveAsCopy" />

        <c-button-submit
          v-if="!hideClone"
          data-test-id="button-clone"
          :disabled="disableClone || processing"
          :title="cloneTooltip"
          :processing="processingClone"
          variant="light"
          :text="$t('label.saveAsCopy')"
          size="lg"
          class="ml-md-2"
          @submit="$emit('clone')"
        />

        <c-button-submit
          v-if="!hideSave"
          data-test-id="button-save-and-close"
          :disabled="disableSave || processing"
          :processing="processingSaveAndClose"
          variant="light"
          :text="$t('label.saveAndClose')"
          size="lg"
          class="ml-md-2"
          @submit="$emit('saveAndClose')"
        />

        <c-button-submit
          v-if="!hideSave"
          data-test-id="button-save"
          :disabled="disableSave || processing"
          :processing="processingSave"
          :text="$t('label.save')"
          size="lg"
          class="ml-md-2"
          @submit="$emit('save')"
        />
      </div>
    </b-row>
  </b-container>
</template>
<script>

export default {
  i18nOptions: {
    namespaces: 'general',
  },

  inheritAttrs: true,

  props: {
    processing: {
      type: Boolean,
      default: false,
    },

    processingSave: {
      type: Boolean,
      default: false,
    },

    processingSaveAndClose: {
      type: Boolean,
      default: false,
    },

    processingClone: {
      type: Boolean,
      default: false,
    },

    processingDelete: {
      type: Boolean,
    },

    backLink: {
      type: Object,
      required: false,
      default: undefined,
    },

    hideDelete: {
      type: Boolean,
      required: false,
    },

    hideSave: {
      type: Boolean,
      required: false,
    },

    hideClone: {
      type: Boolean,
      required: false,
    },

    disableDelete: {
      type: Boolean,
      required: false,
      default: false,
    },

    disableSave: {
      type: Boolean,
      required: false,
      default: false,
    },

    disableClone: {
      type: Boolean,
      default: false,
    },

    deleteTooltip: {
      type: String,
      required: false,
      default: '',
    },

    cloneTooltip: {
      type: String,
      default: '',
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

[dir="rtl"] {
  .back {
    .back-icon {
      display: none;
    }
  }
}
</style>
