<template>
  <div>
    <h5
      v-if="!root"
      class="mb-0"
    >
      {{  group  }}
    </h5>

    <div
      v-if="subgroups.length > 0 || components.length > 0"
    >
      <div
        v-for="(cmp, i) in components"
        :key="i"
        @click="$emit('select', cmp)"
        class="component ml-2"
      >
        {{ cmp.name || cmp.component.name || 'Untitled' }}
        <b-badge
          v-if="cmp.wip"
          variant="warning"
          class="float-right"
        >
          wip
        </b-badge>
      </div>
      <component-list
        v-for="(group) in subgroups"
        :key="group"
        :catalogue="catalogue"
        :path="[...path, group]"
        @select="$emit('select', $event)"
        class="my-3"
      />
    </div>
  </div>
</template>
<script>
import {ExtractComponents, ExtractSubgroups} from './helpers.ts'

export default {
  name: 'ComponentList',
  props: {
    catalogue: {
      required: true,
      type: Object,
    },

    path: {
      type: Array,
      default: () => []
    },
  },

  computed: {
    // name of the current group
    root () {
      return this.path.length === 0
    },

    // name of the current group
    group () {
      return this.root ? undefined : this.path[this.path.length - 1]
    },

    // returns all groups at this level
    subgroups () {
      return ExtractSubgroups(this.catalogue, ...this.path)
    },

    components () {
      // return ExtractComponents(this.catalogue, ...this.path)
      return ExtractComponents(this.catalogue, ...this.path)
    }
  }
}
</script>
<style lang="scss" scoped>
.component {
  cursor: pointer;
}
</style>
