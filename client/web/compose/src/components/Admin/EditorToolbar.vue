<template>
  <c-toolbar
    data-test-id="editor-toolbar"
    class="bg-white shadow border-top"
  >
    <template #start>
      <b-button
        data-test-id="button-back-without-save"
        variant="link"
        :disabled="processing"
        class="text-dark back text-left text-nowrap p-1"
        @click="$emit('back')"
      >
        <font-awesome-icon
          :icon="['fas', 'chevron-left']"
          class="back-icon"
        />
        {{ $t('label.backWithoutSave') }}
      </b-button>
    </template>

    <template #center>
      <slot />
    </template>

    <template #end>
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
        class="text-nowrap"
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
        class="text-nowrap"
        @submit="$emit('saveAndClose')"
      />

      <c-button-submit
        v-if="!hideSave"
        data-test-id="button-save"
        :disabled="disableSave || processing"
        :processing="processingSave"
        :text="$t('label.save')"
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
  i18nOptions: {
    namespaces: 'general',
  },

  components: {
    CToolbar,
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
