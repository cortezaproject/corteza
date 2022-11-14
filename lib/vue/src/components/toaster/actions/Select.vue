<template>
  <b-form-select
    v-bind="options"
    :value="null"
    class="w-auto"
    @change="$emit('action', $event)"
  >
    <template v-slot:first>
      <option
        :value="null"
        disabled
        v-html="label"
      />
    </template>

    <option
      v-for="(opt, i) in options.items || []"
      :key="i"
      v-bind="opt"
      :value="opt.value"
      v-html="opt.label"
    />
  </b-form-select>
</template>

<script lang="js">
import base from './base.vue'

export default {
  extends: base,
}

</script>
