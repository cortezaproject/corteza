<template>
  <div>
    <fieldset class="form-group">
      <label
        class="text-primary"
      >
        {{ label }}
      </label>
      <div>{{ value }}</div>
    </fieldset>
  </div>
</template>
<script>
import multi from './multi'
import errors from '../errors'
import { compose, validator } from '@cortezaproject/corteza-js'
import Hint from 'corteza-webapp-compose/src/components/Common/Hint.vue'

export default {
  components: {
    // multi is used in the components that extends base
    // eslint-disable-next-line vue/no-unused-components
    multi,

    // errors is used in the components that extends base
    // eslint-disable-next-line vue/no-unused-components
    errors,

    // Hint is used in the components that extends base
    // eslint-disable-next-line vue/no-unused-components
    Hint,
  },

  props: {
    namespace: {
      type: compose.Namespace,
      required: true,
    },

    field: {
      type: compose.ModuleField,
      required: true,
    },

    record: {
      type: compose.Record,
      required: true,
    },

    errors: {
      type: validator.Validated,
      required: true,
    },

    valueOnly: {
      type: Boolean,
      default: false,
    },

    appendToBody: {
      type: Boolean,
      default: true,
    },
  },

  computed: {
    formGroupStyleClasses () {
      return {
        required: this.field.isRequired,

        // CSS class "small" is set on form group
        // because of font-size: inherit prop on .col-form-label on
        // wrapping element
        small: false,
      }
    },

    state () {
      if (!this.errors.valid()) {
        return null
      }

      return this.errors.valid() === true ? null : false
    },

    value: {
      get () {
        if (this.field.isSystem) {
          return this.record[this.field.name]
        }

        return this.record.values[this.field.name]
      },

      set (value) {
        if (this.field.isSystem) {
          this.record[this.field.name] = value
        } else {
          this.record.values[this.field.name] = value
        }

        this.$emit('change', value)
      },
    },

    showPopover: {
      get () {
        return this.preventPopoverClose
      },

      set (v) {
        this.$emit('update:preventPopoverClose', v)
      },
    },

    label () {
      if (this.valueOnly) {
        return ''
      }

      return this.field.label || this.field.name
    },

    description () {
      if (this.valueOnly) {
        return ''
      }

      const { view, edit } = this.field.options.description

      return edit || view
    },

    hint () {
      if (this.valueOnly) {
        return ''
      }
      const { view, edit } = this.field.options.hint

      return edit || view
    },
  },
}
</script>
