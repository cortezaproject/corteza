<template>
  <div class="h-100 p-2">
    <div
      class="d-flex flex-column position-relative h-100 border-0"
      :class="blockClass"
    >
      <div
        v-if="headerSet || block.title || block.description"
        :class="`card-header bg-transparent border-0 text-nowrap px-3 text-${block.style.variants.headerText}`"
      >
        <div
          v-if="!headerSet"
        >
          <div
            v-if="block.title"
            class="d-flex justify-content-between align-items-center"
          >
            <h5
              class="text-truncate mb-0"
            >
              {{ block.title }}

              <slot name="title-badge" />
            </h5>

            <font-awesome-icon
              v-if="block.options.refreshRate >= 5"
              :icon="['fa', 'sync']"
              class="h6 text-secondary mb-0"
              role="button"
              @click="$emit('refreshBlock')"
            />
          </div>

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
      </div>

      <div
        v-if="toolbarSet"
        class="overflow-hidden"
      >
        <slot
          name="toolbar"
        />
      </div>

      <div
        class="card-body p-0"
        :class="{ 'overflow-auto': scrollableBody }"
        style="flex-shrink: 10;"
      >
        <slot
          name="default"
        />
      </div>

      <b-card-footer
        v-if="footerSet"
        class="card-footer bg-transparent p-0 overflow-hidden"
      >
        <slot
          name="footer"
        />
      </b-card-footer>
    </div>
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
