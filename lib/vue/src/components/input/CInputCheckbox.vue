<template>
    <div
      class="d-flex"
    >
      <span
        class="mb-0 inline-block text-primary mr-2"
        :class="{ 
          'text-muted': value,
          'font-weight-bold': !value,
       }"
      >
        {{ labels.off }}
      </span>

      <b-form-checkbox
        v-bind="$attrs"
        :checked="value"
        v-on="$listeners"
      />

      <span
        class="mb-0 inline-block strong text-primary ml-2"
        :class="{ 
          'text-muted': !value,
          'font-weight-bold': value,
       }"
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
  },
}
</script>
