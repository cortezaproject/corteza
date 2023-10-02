<template>
  <b-container
    fluid
    class="bg-white shadow border-top p-3"
  >
    <b-row
      no-gutters
      class="align-items-center"
    >
      <b-col>
        <b-button
          v-if="backLink"
          data-test-id="button-back"
          variant="link"
          size="lg"
          :to="backLink"
          class="d-flex align-items-center p-0 text-dark back"
        >
          <font-awesome-icon
            :icon="['fas', 'chevron-left']"
            class="back-icon mr-1"
          />
          {{ $t('general:label.back') }}
        </b-button>
      </b-col>

      <b-col
        class="d-flex justify-content-center"
      >
        <slot />
      </b-col>

      <b-col
        class="d-flex justify-content-end"
      >
        <c-input-confirm
          v-if="!hideDelete"
          class="mr-1"
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
      </b-col>
    </b-row>
  </b-container>
</template>

<script>

export default {
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
