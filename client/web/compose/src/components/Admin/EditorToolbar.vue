<template>
  <b-container
    fluid
    data-test-id="editor-toolbar"
    class="bg-white shadow border-top p-3"
  >
    <b-row
      no-gutters
      class="wrap-with-vertical-gutters align-items-center"
    >
      <div
        class="wrap-with-vertical-gutters align-items-center"
      >
        <b-button
          data-test-id="button-back-without-save"
          variant="link"
          :disabled="processing"
          class="text-dark back mr-auto"
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
        class="ml-auto"
      >
        <slot />
      </div>
      <div
        class="d-flex wrap-with-vertical-gutters align-items-center ml-auto"
      >
        <slot name="delete" />
        <c-input-confirm
          v-if="!hideDelete"
          v-b-tooltip.hover
          :disabled="disableDelete || processing"
          size="lg"
          size-confirm="lg"
          variant="danger"
          :title="deleteTooltip"
          :borderless="false"
          @confirmed="$emit('delete')"
        >
          {{ $t('label.delete') }}
        </c-input-confirm>
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
          class="ml-2"
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
          class="ml-2"
          @submit="$emit('saveAndClose')"
        />
        <c-button-submit
          v-if="!hideSave"
          data-test-id="button-save"
          :disabled="disableSave || processing"
          :processing="processingSave"
          :text="$t('label.save')"
          size="lg"
          class="ml-2"
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
