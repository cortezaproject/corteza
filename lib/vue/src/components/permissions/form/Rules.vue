<template>
  <div
    data-test-id="role-permissions-list"
  >
    <rule
      v-for="(p, i) in rules"
      :key="p.resource + p.operation"
      v-bind="p"
      :class="{ 'mt-4': i > 0 }"
      @update="onUpdate"
    />
  </div>
</template>
<script lang="js">
import Rule from './Rule.vue'

export default {
  components: {
    Rule,
  },

  props: {
    rules: {
      type: Array,
      required: true,
    },
  },

  methods: {
    onUpdate ({ resource, operation, access }) {
      let rr = this.rules
      let ri = rr.findIndex(r => r.resource === resource && r.operation === operation)
      if (ri > -1) {
        rr[ri].access = access
        this.$emit('update:rules', rr)
      }
    },
  },
}
</script>
