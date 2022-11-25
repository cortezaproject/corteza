<template>
  <div class="h-100">
    <b-card
      no-body
      class="h-100 border-0 shadow"
      :class="blockClass"
    >
      <b-card-header
        v-if="headerSet || block.title || block.description || hasMagnifyIcon || block.options.refreshRate >= 5"
        class="border-0 text-nowrap px-3"
        header-bg-variant="white"
        :header-text-variant="block.style.variants.headerText"
      >
        <div
          v-if="!headerSet"
        >
          <div class="d-flex">
            <h5
              v-if="block.title"
              class="text-truncate mb-0"
            >
              {{ block.title }}

              <slot name="title-badge" />
            </h5>

            <div
              v-if="block.options.refreshRate >= 5 || hasMagnifyIcon"
              class="ml-auto"
            >
              <font-awesome-icon
                v-if="block.options.refreshRate >= 5"
                :icon="['fa', 'sync']"
                class="h6 text-secondary"
                role="button"
                @click="$emit('refreshBlock')"
              />

              <font-awesome-icon
                v-if="hasMagnifyIcon"
                :icon="['fas', isBlockOpened ? 'times' : 'search-plus']"
                :title="$t(isBlockOpened ? '' : 'general.label.magnify')"
                class="h6 text-secondary ml-2"
                role="button"
                @click="$root.$emit('magnify-page-block', isBlockOpened ? undefined : block.blockID)"
              />
            </div>
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
        class="p-0 overflow-hidden"
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
        this.block.kind,
      ]
    },

    hasMagnifyIcon () {
      return this.block.options.magnifyOption
    },

    isBlockOpened () {
      return this.block.blockID === this.$route.query.blockID
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
