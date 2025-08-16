<template>
  <b-dropdown
    menu-class="text-center bg-white"
    variant="link"
    boundary="window"
    no-caret
  >
    <template #button-content>
      <span class="text-dark font-weight-bold">
        <span :class="{ 'text-primary': !!activeType && activeType !== 'left' }">
          <font-awesome-icon
            v-if="activeIcon"
            :icon="activeIcon"
          />
          <span v-else>
            {{ format.label }}
          </span>
        </span>
      </span>
    </template>

    <b-dropdown-item-button
      v-for="v of format.variants"
      :key="v.variant"
      @click="$emit('click', { type: 'alignment', attrs: v.attrs })"
    >
      <font-awesome-icon
        v-if="format.icon"
        :icon="v.icon"
      />
      <span v-else>
        {{ v.label }}
      </span>
    </b-dropdown-item-button>
  </b-dropdown>
</template>

<script>
import base from '../TNode/base.vue'

/**
 * Component is used to display node alignment formatting
 */
export default {
  name: 'TNattrAlignment',

  extends: base,

  props: {
    isActive: {
      type: Object,
      required: false,
      default: () => ({}),
    },
  },

  computed: {
    activeType () {
      const alignments = ['left', 'center', 'right', 'justify']

      return alignments.find(alignment =>
        this.editor.isActive({ textAlign: alignment }),
      )
    },

    activeIcon () {
      const alignmentMap = {
        left: 'align-left',
        center: 'align-center',
        right: 'align-right',
        justify: 'align-justify',
      }

      return alignmentMap[this.activeType] || 'align-left'
    },
  },
}
</script>

<style lang="scss">
</style>
