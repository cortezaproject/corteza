<template>
  <b-button
    data-test-id="button-permissions"
    :title="tooltip"
    :variant="buttonVariant"
    :size="size"
    @click="onClick"
  >
    <slot>
      <font-awesome-icon v-if="showButtonIcon" :icon="['fas', 'lock']" />
      <span v-if="buttonLabel">
        {{ buttonLabel }}
      </span>
    </slot>
  </b-button>
</template>
<script lang="js">
import { modalOpenEventName } from './def.ts'
import { library } from '@fortawesome/fontawesome-svg-core'
import { faLock } from '@fortawesome/free-solid-svg-icons'

library.add(faLock)

export default {
  props: {
    size: {
      type: String,
      default: 'md',
    },

    buttonVariant: {
      type: String,
      default: 'light',
    },

    resource: {
      type: String,
      required: true,
    },

    title: {
      type: String,
      default: undefined,
    },

    buttonLabel: {
      type: String,
      default: undefined,
    },

    modalOpenEvent: {
      type: String,
      default: modalOpenEventName,
    },

    target: {
      type: String,
      required: false,
      default: undefined,
    },

    showButtonIcon: {
      type: Boolean,
      default: true,
    },

    // Use this prop if you want the translations to look for all-specific key instead of all/specific
    allSpecific: {
      type: Boolean,
      default: false
    },

    tooltip: {
      type: String,
      default: '',
    },
  },

  methods: {
    onClick () {
      this.$root.$emit(this.modalOpenEvent, {
        target: this.target,
        resource: this.resource,
        title: this.title,
        allSpecific: this.allSpecific
      })
    },
  },
}
</script>
<style lang="scss" scoped>
.pointer {
  cursor: pointer;
}
</style>
