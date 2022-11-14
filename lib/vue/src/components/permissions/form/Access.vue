<template>
  <b-form-radio-group
    v-model="selected"
    data-test-id="toggle-role-permissions"
    buttons
    :disabled="!enabled"
    :button-variant="variant"
    :options="options"
    class="access rounded"
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

<style lang="scss">
.access {
  .btn {
    background-color: #E4E9EF;
    border: none;
  }

  .btn:nth-child(2), .btn:nth-child(3) {
    margin-left: 0.2rem !important;
  }
}
</style>
