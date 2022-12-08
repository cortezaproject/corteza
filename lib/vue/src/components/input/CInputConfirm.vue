<template>
  <span class="text-center">
    <span v-if="!inConfirmation">
      <b-button
        data-test-id="button-delete"
        :variant="variant"
        :size="size"
        :disabled="disabled"
        :title="tooltip"
        :class="`${buttonClass} ${borderless ? 'border-0' : ''}`"
        @click.stop.prevent="onPrompt"
      >
        <slot>
          <font-awesome-icon
            :icon="['far', 'trash-alt']"
          />
        </slot>
      </b-button>

    </span>
    <span v-else>
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
            :icon="['fas', 'times']"
          /></slot>
      </b-button>
    </span>
  </span>
</template>
<script lang="js">
export default {
  props: {
    disabled: Boolean,
    okDisabled: Boolean,
    cancelDisabled: Boolean,
    noPrompt: Boolean,

    buttonClass: {
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
      default: undefined,
    },
    sizeConfirm: {
      type: String,
      default: 'sm',
    },
    tooltip: {
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
