<template>
  <div
    :class="{ 'text-center': inConfirmation }"
    class="d-inline-block"
  >
    <template v-if="!inConfirmation">
      <b-button
        data-test-id="button-delete"
        :variant="variant"
        :size="size"
        :disabled="disabled || processing"
        :title="tooltip"
        :class="`${buttonClass} ${borderless ? 'border-0' : ''}`"
        @click.stop.prevent="onPrompt"
      >
        <b-spinner
          v-if="processing"
          data-test-id="spinner"
          class="align-middle"
          small
        />
        <slot v-else>
          <font-awesome-icon
            v-if="showIcon"
            :icon="icon"
            :class="iconClass"
          />
          <span
            v-if="text"
            :class="textClass"
          >
            {{ text }}
          </span>
        </slot>
      </b-button>
    </template>
    <template v-else>
      <b-button
        data-test-id="button-delete-confirm"
        :variant="variantOk"
        :size="sizeConfirm"
        :disabled="okDisabled"
        class="mr-1"
        :class="[ borderless && 'border-0' ]"
        @blur.prevent="onCancel()"
        @click.prevent.stop="onConfirmation()"
      >
        <slot name="yes">
          <font-awesome-icon
            data-test-id="confirm"
            :icon="['fas', 'check']"
          />
        </slot>
      </b-button>
      <b-button
        data-test-id="button-delete-cancel"
        :variant="variantCancel"
        :size="sizeConfirm"
        :disabled="cancelDisabled"
        :class="[ borderless && 'border-0' ]"
        @click.prevent.stop="onCancel()"
      >
        <slot name="no">
          <font-awesome-icon
            data-test-id="reject"
            :icon="['fas', 'times']"
          />
        </slot>
      </b-button>
    </template>
  </div>
</template>
<script lang="js">
export default {
  props: {
    disabled: Boolean,
    okDisabled: Boolean,
    cancelDisabled: Boolean,
    noPrompt: Boolean,
    processing: Boolean,
    showIcon: Boolean,

    icon: {
      type: Array,
      default: () => ['far', 'trash-alt'],
    },

    buttonClass: {
      type: String,
      default: '',
    },

    iconClass: {
      type: String,
      default: '',
    },

    textClass: {
      type: String,
      default: '',
    },

    borderless: {
      type: Boolean,
      default: true,
    },

    variant: {
      type: String,
      default: 'outline-danger',
    },

    size: {
      type: String,
      default: 'sm',
    },

    variantOk: {
      type: String,
      default: 'danger',
    },

    variantCancel: {
      type: String,
      default: 'light',
    },

    sizeConfirm: {
      type: String,
      default: 'sm',
    },

    tooltip: {
      type: String,
      default: '',
    },

    text: {
      type: String,
      default: '',
    },
  },

  data () {
    return {
      inConfirmation: false,
    }
  },

  methods: {
    onPrompt () {
      if (this.noPrompt) {
        this.$emit('confirmed')
      } else {
        this.inConfirmation = true
      }
    },

    onConfirmation () {
      this.inConfirmation = false
      this.$emit('confirmed')
    },

    onCancel () {
      this.inConfirmation = false
      this.$emit('canceled')
    },
  },
}
</script>
