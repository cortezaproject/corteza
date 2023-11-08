<template>
  <b-container
    fluid
    class="bg-white shadow-sm border-top p-3"
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
          class="d-flex align-items-center text-dark back"
        >
          <font-awesome-icon
            :icon="['fas', 'chevron-left']"
            class="back-icon mr-1"
          />
          {{ $t('general:label.back') }}
        </b-button>

        <slot name="left" />
      </b-col>

      <b-col
        class="d-flex justify-content-center"
      >
        <slot name="middle" />
      </b-col>

      <b-col
        class="d-flex justify-content-end"
      >
        <template v-if="deleteShow">
          <c-input-confirm
            v-if="deleteConfirm"
            :disabled="deleteDisabled || processing"
            :processing="processingDelete"
            :text="deleteLabel"
            variant="danger"
            size="lg"
            size-confirm="lg"
            class="ml-1"
            @confirmed="$emit('delete')"
          />

          <b-button
            v-else
            :data-test-id="buttonLabelCypressId(deleteLabel)"
            :disabled="deleteDisabled || processing"
            variant="danger"
            size="lg"
            class="ml-1"
            @click="$emit('delete')"
          >
            {{ deleteLabel }}
          </b-button>
        </template>

        <c-button-submit
          v-if="submitShow"
          :data-test-id="buttonLabelCypressId(submitLabel)"
          :disabled="submitDisabled"
          :processing="processing"
          :text="submitLabel"
          size="lg"
          class="ml-1"
          @submit="$emit('submit')"
        />

        <slot />
      </b-col>
    </b-row>
  </b-container>
</template>

<script>

export default {
  props: {
    processing: {
      type: Boolean,
    },

    processingDelete: {
      type: Boolean,
    },

    backLink: {
      type: Object,
      default: () => ({ name: 'root' }),
    },

    deleteShow: {
      type: Boolean,
    },

    deleteDisabled: {
      type: Boolean,
    },

    deleteConfirm: {
      type: Boolean,
      default: true,
    },

    deleteLabel: {
      type: String,
      default: '',
    },

    submitShow: {
      type: Boolean,
    },

    submitDisabled: {
      type: Boolean,
    },

    submitLabel: {
      type: String,
      default: '',
    },
  },

  methods: {
    buttonLabelCypressId (label) {
      return `button-${label.toLowerCase().split(' ').join('-')}`
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
