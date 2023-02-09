<template>
  <div class="h-100 p-2">
    <b-card
      no-body
      class="h-100 shadow"
      :class="blockClass"
    >
      <b-card-header
        v-if="headerSet || block.title || block.description"
        class="border-0 text-nowrap px-3"
        :class="{ 'p-0': !(block.title || block.description)}"
        header-bg-variant="white"
        :header-text-variant="block.style.variants.headerText"
      >
        <div
          v-if="!headerSet"
        >
          <h5
            v-if="block.title"
            class="text-truncate mb-0"
          >
            {{ block.title }}
          </h5>

          <b-card-text
            v-if="block.description"
            class="text-dark text-truncate mt-1"
          >
            {{ block.description }}
          </b-card-text>
        </div>

        <slot
          v-else
          name="header"
        />
      </b-card-header>

      <div
        v-if="toolbarSet"
        class="overflow-hidden"
      >
        <slot
          name="toolbar"
        />
      </div>

      <b-card-body
        class="p-0"
        :class="{ 'overflow-auto': scrollableBody }"
        style="flex-shrink: 10;"
      >
        <slot
          name="default"
        />
      </b-card-body>

      <b-card-footer
        v-if="footerSet"
        class="p-0 overflow-hidden bg-white border-top"
      >
        <slot
          name="footer"
        />
      </b-card-footer>
    </b-card>
  </div>
</template>
<script>
import { compose } from '@cortezaproject/corteza-js'

export default {
  props: {
    block: {
      type: compose.PageBlock,
      required: true,
    },

    scrollableBody: {
      type: Boolean,
      required: false,
      default: () => true,
    },
  },

  computed: {
    blockClass () {
      return [
        'block',
        { border: this.block.style.border.enabled },
        this.block.kind,
      ]
    },

    headerSet () {
      return !!this.$scopedSlots.header
    },

    toolbarSet () {
      return !!this.$scopedSlots.toolbar
    },

    footerSet () {
      return !!this.$scopedSlots.footer
    },
  },
}
</script>
