<template>
  <b-button
    v-b-tooltip.hover="{ title: tooltip, container: '#body' }"
    data-test-id="button-translation"
    :variant="buttonVariant"
    :class="buttonClass"
    :disabled="disabled"
    :size="size"
    @click="onClick"
  >
    <slot>
      <font-awesome-icon :icon="['fas', 'language']" />
    </slot>
  </b-button>
</template>
<script lang="js">
import { library } from '@fortawesome/fontawesome-svg-core'
import { faLanguage } from '@fortawesome/free-solid-svg-icons'

library.add(faLanguage)

export default {
  props: {
    buttonVariant: {
      type: String,
      default: () => { return 'light' },
    },

    buttonClass: {
      type: String,
      default: () => { return '' },
    },

    size: {
      type: String,
      default: 'md',
    },

    disabled: {
      type: Boolean,
      default: () => false,
    },

    resource: {
      type: String,
      required: true,
    },

    /**
     * See CTranslatorForm for description
     */
    highlightKey: {
      type: String,
      default: '',
    },

    tooltip: {
      type: String,
      default: '',
    },

    titles: {
      type: Object,
      default: () => ({}),
    },

    fetcher: {
      type: Function,
      default: undefined,
    },

    updater: {
      type: Function,
      default: undefined,
    },

    keyPrettyfier: {
      type: Function,
      default: undefined,
    },
  },

  methods: {
    onClick () {
      this.$root.$emit('c-translator', {
        resource: this.resource,
        titles: this.titles,
        highlightKey: this.highlightKey,
        fetcher: this.fetcher,
        updater: this.updater,
        keyPrettyfier: this.keyPrettyfier,
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
