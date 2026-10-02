<template>
  <div class="d-flex gap-2">
    <b-button
      v-for="{ b, i } in visibleButtons"
      :key="i"
      :variant="variant(b)"
      :disabled="!isValid(b) || processingIDs.includes(i)"
      :class="buttonClass"
      @click.prevent="handle(b, i)"
    >
      {{ buttonLabel(b.label) || '-' }}
    </b-button>
  </div>
</template>
<script>
import { compose, automation, NoID } from '@cortezaproject/corteza-js'
import { debounce } from 'lodash'
import { evaluatePrefilter } from 'corteza-webapp-compose/src/lib/record-filter'
import base from '../base'

export default {
  extends: base,

  props: {
    buttons: {
      type: Array,
      required: true,
    },

    automationScripts: {
      type: Array,
      required: false,
      default: () => [],
    },

    buttonClass: {
      type: String,
      default: '',
    },

    extraEventArgs: {
      type: Object,
      default: () => ({}),
    },
  },

  data () {
    return {
      processingIDs: [],

      // Results of button visibility expressions, keyed by button index
      visibility: {},
    }
  },

  computed: {
    hasVisibilityExpressions () {
      return this.buttons.some(b => this.buttonExpression(b))
    },

    visibleButtons () {
      return this.buttons
        .map((b, i) => ({ b, i }))
        .filter(({ b, i }) => !this.buttonExpression(b) || this.visibility[i])
    },
  },

  watch: {
    buttons: {
      immediate: true,
      handler () {
        this.evaluateVisibility()
      },
    },

    record () {
      this.evaluateVisibilityDebounced()
    },
  },

  created () {
    this.$root.$on('record-field-change', this.evaluateVisibilityDebounced)
  },

  beforeDestroy () {
    this.$root.$off('record-field-change', this.evaluateVisibilityDebounced)
  },

  methods: {
    buttonExpression (b) {
      return ((b || {}).visibility || {}).expression || ''
    },

    evaluateVisibilityDebounced: debounce(function () {
      this.evaluateVisibility()
    }, 300),

    async evaluateVisibility () {
      if (!this.hasVisibilityExpressions || this.$route.name === 'admin.pages.builder') {
        this.visibility = {}
        return
      }

      const expressions = {}
      this.buttons.forEach((b, i) => {
        const expression = this.buttonExpression(b)
        if (expression) {
          expressions[i] = expression
        }
      })

      const variables = {
        user: this.$auth.user,
        record: this.record ? this.record.serialize() : {},
        screen: {
          width: window.innerWidth,
          height: window.innerHeight,
          userAgent: navigator.userAgent,
          breakpoint: this.getBreakpoint(),
        },
      }

      return this.$SystemAPI.expressionEvaluate({ variables, expressions })
        .then(res => {
          this.visibility = res || {}
        })
        .catch(e => {
          this.toastErrorHandler(this.$t('notification:evaluate.failed'))(e)
          this.visibility = {}
        })
    },

    /**
     *
     */
    variant (b) {
      if (!b.script) {
        return b.variant
      }

      if (!this.isValid(b)) {
        // Does this script actually exist?
        return 'outline-danger'
      }

      return b.variant || 'primary'
    },

    isValid (b) {
      // Check if event resource types args exist
      const { resourceType } = b

      let paramsExist = true

      if (resourceType === 'compose:record') {
        paramsExist = this.record && this.module
      } else if (resourceType === 'compose:module') {
        paramsExist = !!this.module
      } else if (resourceType === 'compose:namespace') {
        paramsExist = !!this.namespace
      } else if (resourceType === 'compose:page') {
        paramsExist = !!this.page
      }

      if (!paramsExist) {
        return false
      } else if (b.workflowID) {
        return true
      } else if (b.script) {
        if (this.$UIHooks.FindByScript(b.script)) {
          return true
        }

        if (!this.automationScripts) {
          return false
        }

        return this.automationScripts.find(({ name }) => name === b.script)
      }

      return false
    },

    handle (b, i) {
      try {
        this.processingIDs.push(i)

        // Base of the raise event:
        // we'll attach all extra arguments passed to component to
        // be part of the generated event
        let ev = { args: this.extraEventArgs || {} }

        // @todo page event missing on backend
        switch (b.resourceType) {
          case 'compose:record':
            // Only generate event if record exists
            ev.args.namespace = this.namespace
            if (this.record && this.module) {
              ev.args.module = this.module
              ev = compose.RecordEvent(this.record, ev)
            }
            break
          case 'compose:module':
            ev.args.namespace = this.namespace
            ev = compose.ModuleEvent(this.module, ev)
            break
          case 'compose:namespace':
            ev = compose.NamespaceEvent(this.namespace, ev)
            break
          case 'compose:page':
            ev.args.namespace = this.namespace
            ev = compose.PageEvent(this.page.pageID, ev)
            break
          case 'compose':
            ev = compose.ComposeEvent(ev)
        }

        if (b.workflowID) {
          const { workflowID, stepID } = b
          const input = automation.Encode(ev.args)

          this.$AutomationAPI
            .workflowExec({
              workflowID,
              stepID,
              input,
            })
            .catch(this.toastErrorHandler(this.$t('notification:automation.scriptFailed')))
            .finally(() => {
              this.processingIDs = this.processingIDs.filter(id => id !== i)
            })

          return
        }

        if (!b.script) {
          return
        }

        // @todo this is not a complete implementation
        //       we need to do a proper filtering via constraint matching
        //       for now, all (configured) buttons are displayed

        // Passing events to eventbus
        //
        // The main reason to do this is because eventbus (or better, handlers registed there)
        // know how to handle each script - is it client or server script, what context to use
        // etc...
        this.$EventBus
          .Dispatch(ev, b.script)
          .catch(this.toastErrorHandler(this.$t('notification:automation.scriptFailed')))
          .finally(() => {
            this.processingIDs = this.processingIDs.filter(id => id !== i)
          })
      } catch (e) {
        this.toastErrorHandler(this.$t('notification:automation.scriptFailed'))(e)
        this.processingIDs = this.processingIDs.filter(id => id !== i)
      }
    },

    buttonLabel (label = '') {
      try {
        return evaluatePrefilter(label, {
          record: this.record,
          user: this.$auth.user || {},
          recordID: (this.record || {}).recordID || NoID,
          ownerID: (this.record || {}).ownedBy || NoID,
          userID: (this.$auth.user || {}).userID || NoID,
        })
      } catch (e) {
        return e
      }
    },
  },
}
</script>
