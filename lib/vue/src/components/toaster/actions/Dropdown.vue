<template>
  <b-dropdown v-bind="options">
    <template v-slot:button-content>
      <span v-html="label" />
    </template>

    <component
      :is="fromKind(opt)"
      v-for="(opt, i) in options.items || []"
      :key="i"
      v-bind="opt"
      @click="$emit('action', opt.value)"
    >
      <span v-html="opt.label" />
    </component>
  </b-dropdown>
</template>

<script lang="js">
import base from './base.vue'

export default {
  extends: base,

  methods: {
    fromKind ({ kind }) {
      return `b-dropdown-${kind}`
    },
  },
}

</script>
