<template>
  <a
    v-if="link"
    data-test-id="link-permissions"
    class="pointer"
    :title="tooltip"
    @click="onClick"
  >
    <font-awesome-icon :icon="['fas', 'lock']" />
  </a>
  <b-button
    v-else
    data-test-id="button-permissions"
    :title="tooltip"
    :variant="buttonVariant"
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
    link: {
      type: Boolean,
    },

    buttonVariant: {
      type: String,
      default: undefined,
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
