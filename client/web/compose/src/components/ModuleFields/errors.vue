<template>
  <div>
    <b-form-invalid-feedback
      v-for="(error, i) in set"
      :key="i"
      force-show
      class="m-0"
    >
      {{ $t(error.message) }}
    </b-form-invalid-feedback>
  </div>
</template>
<script>
import { validator } from '@cortezaproject/corteza-js'

export default {
  props: {
    errors: {
      type: validator.Validated,
      required: true,
      default: undefined,
    },

    index: {
      type: Number,
      required: false,
      default: -1,
    },
  },

  computed: {
    set () {
      return this.index >= 0 ? this.errors.filterByMeta('index', this.index).get() : this.errors.get()
    },
  },
}
</script>
