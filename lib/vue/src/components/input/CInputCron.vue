<template>

</template>
<script lang="js">
export default {
  props: {
    value: {
      type: String,
      required: true,
    },

    simple: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    cronValue: {
      get () {
        return this.value
      },

      set (value) {
        this.$emit('input', value)
      },
    },
  },
}
</script>
