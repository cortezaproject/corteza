<template>
  <b-button
    data-test-id="button-submit"
    type="submit"
    :variant="variant"
    :disabled="disabled || processing || success"
    :size="size"
    :block="block"
    :title="title"
    :class="buttonClass"
    @click.prevent="$emit('submit')"
  >
    <template v-if="processing">
      <span
        v-if="loadingText"
        data-test-id="button-loading-text"
        class="loading-text mx-2"
      >
        {{ loadingText }}
      </span>
      <b-spinner
        v-else
        data-test-id="spinner"
        class="align-middle"
        small
      />
    </template>
    <template v-else-if="success">
      <font-awesome-icon
        data-test-id="icon-success"
        :icon="['fas', 'check']"
        :class="iconVariant"
        class="text-white"
      />
    </template>
    <template v-else>
      <span
        data-test-id="button-text"
      >
        {{ text }}
      </span>
    </template>
  </b-button>
</template>

<script>
export default {
  name: 'CButtonSubmit',

  props: {
    processing: {
      type: Boolean,
    },

    success: {
      type: Boolean,
    },

    disabled: {
      type: Boolean,
    },

    title: {
      type: String,
      default: '',
    },

    buttonClass: {
      type: String,
      default: '',
    },

    text: {
      type: String,
      default: '',
    },

    loadingText: {
      type: String,
      default: '',
    },

    size: {
      type: String,
      default: 'md',
    },

    block: {
      type: Boolean,
    },

    variant: {
      type: String,
      default: 'primary',
    },

    iconVariant: {
      type: String,
      default: 'text-white',
    },
  },
}
</script>

<style scoped>
.loading-text::after {
  display: inline-block;
  animation: saving steps(1, end) 1s infinite;
  content: '';
}
</style>
