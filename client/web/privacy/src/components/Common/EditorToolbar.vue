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
            :borderless="false"
            variant="danger"
            size="lg"
            size-confirm="lg"
            class="ml-1"
            @confirmed="$emit('delete')"
          >
            {{ deleteLabel }}
          </c-input-confirm>

          <b-button
            v-else
            :disabled="deleteDisabled || processing"
            variant="danger"
            size="lg"
            class="ml-1"
            @click="$emit('delete')"
          >
            {{ deleteLabel }}
          </b-button>
        </template>

        <b-button
          v-if="submitShow"
          :disabled="submitDisabled || processing"
          variant="primary"
          size="lg"
          class="ml-1"
          @click="$emit('submit')"
        >
          {{ submitLabel }}
        </b-button>

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
      required: false,
    },

    backLink: {
      type: Object,
      default: () => ({ name: 'root' }),
    },

    deleteShow: {
      type: Boolean,
      required: false,
    },

    deleteDisabled: {
      type: Boolean,
      required: false,
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
      required: false,
    },

    submitDisabled: {
      type: Boolean,
      required: false,
    },

    submitLabel: {
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
</style>
