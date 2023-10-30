<template>
  <div
    class="d-flex align-items-center"
  >
    <span
      class="mb-0 inline-block text-primary"
      :class="offClass"
    >
      {{ labels.off }}
    </span>

    <b-form-checkbox
      v-bind="$attrs"
      :checked="value"
      :value="!invert"
      :unchecked-value="invert"
      v-on="$listeners"
    >
      <slot />
    </b-form-checkbox>

    <span
      class="mb-0 inline-block strong text-primary"
      :class="onClass"
    >
      {{ labels.on }}
    </span>
  </div>
</template>

<script>
export default {
  props: {
    value: {
      type: Boolean,
    },

    labels: {
      type: Object,
      default: () => ({}),
    },

    invert: {
      type: Boolean,
      default: false,
    }
  },

  computed: {
    offClass () {
      const value = this.invert ? !this.value : this.value
      return {
        'text-muted': value,
        'font-weight-bold': !value,
        'mr-2': !!this.labels.off
      }
    },

    onClass () {
      const value = this.invert ? !this.value : this.value
      return {
        'text-muted': !value,
        'font-weight-bold': value,
        'ml-1': !!this.labels.on
      }
    }
  },
}
</script>
