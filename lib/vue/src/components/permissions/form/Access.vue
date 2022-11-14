<template>
  <b-form-radio-group
    v-model="selected"
    buttons
    :disabled="!enabled"
    :button-variant="variant"
    :options="options"
  />
</template>
<script lang="js">
export default {
  props: {
    enabled: {
      type: Boolean,
      default: true,
    },

    access: {
      type: String,
      required: false,
      default: undefined,
    },

    current: {
      type: String,
      required: false,
      default: undefined,
    },
  },

  computed: {
    options () {
      return ['allow', 'inherit', 'deny'].map(value => ({
        value,
        text: this.$t('permissions:ui.access.' + value)
      }))
    },

    isChanged () {
      return this.selected !== this.current
    },

    variant () {
      return this.isChanged ? 'outline-warning' : 'outline-primary'
    },

    selected: {
      get () {
        return this.access
      },

      set (sel) {
        if (this.access !== sel) {
          this.$emit('update', sel)
        }

        this.$emit('update:access', sel)
      },
    },
  },
}
</script>
