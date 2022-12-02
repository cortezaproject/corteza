<template>
  <b-button
    :data-test-id="cypressID"
    type="submit"
    :variant="variant"
    :disabled="disabled || processing || success"
    class="submit"
    :class="buttonClass"
    @click.prevent="$emit('submit')"
  >
    <b-spinner
      v-if="processing"
      small
      :variant="iconVariant"
      type="grow"
    />
    <font-awesome-icon
      v-else-if="success"
      :icon="['fas', 'check']"
      class="text-white h3 mb-0"
      :class="iconVariant"
    />
    <span
      v-else
    >
      <slot>{{ $t('general.label.submit') }}</slot>
    </span>
  </b-button>
</template>

<script>
export default {
  name: 'CSubmitButton',

  i18nOptions: {
    namespaces: 'admin',
  },

  props: {
    processing: {
      type: Boolean,
      value: false,
    },

    success: {
      type: Boolean,
      value: false,
    },

    disabled: {
      type: Boolean,
      value: false,
    },

    buttonClass: {
      type: String,
      default: '',
    },

    variant: {
      type: String,
      default: 'primary',
    },

    iconVariant: {
      type: String,
      default: 'text-white',
    },

    cypressID: {
      type: String,
      default: 'button-submit',
    },
  },
}
</script>

<style>
.submit {
  min-width: 75px;
  min-height: 35px;
}
</style>
