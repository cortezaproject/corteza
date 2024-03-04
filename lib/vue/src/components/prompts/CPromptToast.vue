<template>
  <div v-if="!hideToasts">
    <b-toast
      v-for="({ prompt, component, passive }) in toasts"
      :id="'wfPromptToast-'+prompt.stateID"
      :key="'wfPromptToast-'+prompt.stateID"
      :variant="pVal(prompt, 'variant', 'primary')"
      :visible="!isActive"
      solid
      :no-auto-hide="!passive"
      :auto-hide-delay="pVal(prompt, 'timeout', defaultTimeout) * 1000"
      @hide="onToastHide({ prompt, passive })"
    >
      <template #toast-title>
        <strong>{{ pVal(prompt, 'title', 'Workflow prompt') }}</strong>
      </template>

      <component
        v-if="component"
        :is="component"
        :payload="prompt.payload"
        :loading="isLoading"
        @submit="resumeToast({ input: $event, prompt })"
      />
    </b-toast>
  </div>
</template>
<script>
import { mapGetters, mapActions } from 'vuex'
import definitions from './kinds/index.ts'
import { pVal } from './utils.ts'

export default {
  name: 'c-prompt-toast',

  props: {
    hideToasts: {
      type: Boolean,
    },
  },

  data () {
    return {
      passive: new Set(),

      hasFocus: null,
      hasFocusObserver: 0,
    }
  },

  computed: {
    ...mapGetters({
      prompts: 'wfPrompts/all',
      isActive: 'wfPrompts/isActive',
      isLoading: 'wfPrompts/isLoading',
    }),

    /**
     * Prompts with handlers, observed with "watch"
     *
     * Prompts are only returned when document has focus!
     *
     * @returns {*}
     */
    withHandlers () {
      return (this.hasFocus ? this.prompts : [])
        .filter(({ ref }) => !!definitions[ref] && !!definitions[ref].handler)
        .map(prompt => ({ ...definitions[prompt.ref], prompt }))
    },

    /**
     * Prompts with components
     *
     * Prompts are only returned when document has focus!
     *
     * @returns {*}
     */
    withComponents () {
      return (this.hasFocus ? this.prompts : [])
        .filter(({ ref }) => !!definitions[ref] && !!definitions[ref].component)
        .map(prompt => ({ ...definitions[prompt.ref], prompt }))
    },

    /**
     * All non-passive prompts with components
     */
    active() {
      return this.withComponents.filter(({ passive }) => !passive)
    },

    /**
     * Returns list of prompts that we can interpret as toasts: display component is defined
     *
     * Toasts (prompts with components) are displayed in order received but
     * passive (no feedback or input from user required) first and the rest later
     */
    toasts () {
      return this.hideToasts ? [] : [
        ...this.passive.values(),
        ...this.active
      ]
    },

    defaultTimeout () {
      return 7
    }
  },

  watch: {
    // watch prompts with handlers and when a new one arrives
    // shift it from the stack, resume the prompt and handle it
    withHandlers (hh) {
      if (hh.length > 0) {
        const { handler, prompt } = hh.shift()
        this.resume({ input: {}, prompt }).then(() => {
          handler.call(this, prompt.payload)
        })
      }
    },

    /**
     * Make a copy of prompt if it is defined as passive
     *
     * We do this because we do not want it to be removed right away
     * but through a toast component's timeout
     */
    withComponents: {
      immediate: true,
      handler (wc) {
        wc.forEach(p => {
          if (p.passive) {
            this.passive.add(p)
          }
        })
      },
    },
  },

  mounted () {
    this.setDocumentFocusObserver()
  },

  beforeDestroy () {
    this.clearDocumentFocusObserver()
    this.setDefaultValues()
  },

  methods: {
    ...mapActions({
      resume: 'wfPrompts/resume',
      cancel: 'wfPrompts/cancel',
      activate: 'wfPrompts/activate',
    }),

    resumeToast (values) {
      // Only reset input if prompt is kept open
      if (values.input && values.input.keep) {
        values.input = {}
      }

      this.resume(values)
    },

    onToastHide ({ prompt, passive}) {
      if (passive) return

      this.cancel(prompt)
    },

    pVal (prompt, k, def = undefined) {
      return pVal(prompt.payload, k, def)
    },

    clearDocumentFocusObserver() {
      if (this.hasFocusObserver) {
        window.clearInterval(this.hasFocusObserver)
      }
    },

    /**
     * Interval check if window has focus
     */
    setDocumentFocusObserver() {
      this.clearDocumentFocusObserver()

      this.hasFocusObserver = window.setInterval(() => {
        const f = document.hasFocus()
        if (this.hasFocus !== f) {
          this.hasFocus = f
        }
      }, 1000)
    },

    setDefaultValues () {
      this.toasts = []
      this.hasFocus = null
      this.hasFocusObserver = 0
    },
  },
}
</script>

<style lang="scss">
.toast-header {
  align-items: start;
  padding: 0.375rem 0.75rem;

  strong {
    word-break: break-word;
  }

  .close {
    margin-bottom: 0 !important;
  }
}

// .b-toaster-leave-active {
//   width: 100%;
// }

.b-toaster {
  &.b-toaster-top-right,
  &.b-toaster-top-left,
  &.b-toaster-bottom-right,
  &.b-toaster-bottom-left {
    .b-toast {
      &.b-toaster-enter-active,
      &.b-toaster-leave-active,
      &.b-toaster-move {
        transition: transform 0.3s ease-in-out; /* Adjust the timing function for smoother transitions */
        opacity: 1; /* Ensure opacity is set to avoid flickering during transition */
      }

      &.b-toaster-enter {
        transform: translate(0, -100%); /* Start off-screen when entering */
        opacity: 0; /* Start with 0 opacity */
      }

      &.b-toaster-enter-to,
      &.b-toaster-enter-active {
        transform: translate(0, 0); /* Move to the visible position */
        opacity: 1; /* Fade in during the transition */
      }

      &.b-toaster-leave-active {
        position: absolute;
        transform: translate(0, -100%); /* Move off-screen when leaving */
        opacity: 0; /* Fade out during the transition */
      }

      &.b-toaster-leave-to {
        opacity: 0; /* Ensure 0 opacity at the end of the transition */
      }
    }
  }
}
</style>