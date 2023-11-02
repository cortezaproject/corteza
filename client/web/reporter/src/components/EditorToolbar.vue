<template>
  <c-toolbar
    class="bg-white shadow border-top"
  >
    <template #start>
      <b-button
        v-if="backLink"
        data-test-id="button-back"
        variant="link"
        size="lg"
        :to="backLink"
        class="text-dark back text-left text-nowrap p-1"
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
        :disabled="deleteDisabled || processingDelete"
        :processing="processingDelete"
        :text="$t('general:label.delete')"
        :borderless="false"
        @confirmed="$emit('delete')"
      />

      <c-button-submit
        data-test-id="button-save"
        :disabled="saveDisabled || processing"
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
