<template>
    <div
      class="d-flex"
    >
      <span
        class="mb-0 inline-block text-primary mr-2"
        :class="noClass"
      >
        {{ labels.off }}
      </span>

      <b-form-checkbox
        v-bind="$attrs"
        :checked="value"
        :value="!invert"
        :unchecked-value="invert"
        v-on="$listeners"
      />

      <span
        class="mb-0 inline-block strong text-primary ml-2"
        :class="yesClass"
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
      required: true,
    },

    labels: {
      type: Object,
      required: true,
      default: () => ({}),
    },

    invert: {
      type: Boolean,
      default: false,
    }
  },

  computed: {
    noClass () {
      const value = this.invert ? !this.value : this.value
      return {
        'text-muted': value,
        'font-weight-bold': !value,
      }
    },

    yesClass () {
      const value = this.invert ? !this.value : this.value
      return {
        'text-muted': !value,
        'font-weight-bold': value,
      }
    }
  },
}
</script>
