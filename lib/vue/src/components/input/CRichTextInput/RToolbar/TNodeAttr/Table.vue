<template>
  <b-dropdown
    menu-class="text-center"
    variant="link"
    boundary="window"
  >
    <template slot="button-content">
      <span class="text-dark font-weight-bold">
        <span :class="rootActiveClasses()">
          <font-awesome-icon
            v-if="format.icon"
            :icon="format.icon"
          />
          <span v-else>
            {{ format.label }}
          </span>
        </span>
      </span>
    </template>

    <b-dropdown-item
      v-for="v of format.variants"
      :key="v.variant"
      @click="emitClick(v)"
    >
      {{ v.label }}
    </b-dropdown-item>
  </b-dropdown>
</template>

<script>
import base from '../TNode/base.vue'
import { nodeTypes } from '../../lib/formats'

/**
 * Component is used to display node alignment formatting
 */
export default {
  name: 'TNattrTable',
  extends: base,

  props: {
    isActive: {
      type: Object,
      required: false,
      default: () => ({}),
    },
  },

  methods: {
    activeClasses (attrs) {
      const an = this.activeNode(nodeTypes, attrs)
      if (!an || !an.node) {
        return undefined
      }

      const ac = (type, attrs) => {
        const b = (this.isActive[type])
        return b && (b(attrs))
      }
      if (ac(an.node.type.name, { ...an.node.attrs, ...attrs })) {
        return ['text-success']
      }

      return undefined
    },

    /**
     * dispatches node attr update for all affected nodes
     * use a single transaction, so ctrl + z works as intended
     */
    dispatchTransaction (v) {
      const ann = this.activeNodes(nodeTypes)
      const tr = this.$attrs.editor.state.tr
      for (const an of ann) {
        tr.setNodeMarkup(an.position, an.node.type, { ...an.node.attrs, ...v.attrs })
      }
      this.$attrs.editor.dispatchTransaction(tr)
    },

    emitClick (v) {
      this.$emit('click', { type: v.type, attrs: { ...v.attrs } })
    },

    /**
     * Helper method to determine if the root formater should be shown as active
     * @returns {Array|undefined}
     */
    rootActiveClasses (v) {
      if (this.format.variants.find(({ type, attrs }) => this.activeClasses(attrs))) {
        return ['text-success']
      }
    },
  },
}
</script>
