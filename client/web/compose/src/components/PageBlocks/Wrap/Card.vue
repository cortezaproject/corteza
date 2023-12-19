<template>
  <div class="h-100">
    <b-card
      no-body
      class="d-flex flex-column h-100 shadow overflow-hidden position-static"
      :class="[blockClass, cardClass]"
    >
      <b-card-header
        v-if="showHeader"
        class="border-0 text-nowrap pl-3 pr-2"
        header-bg-variant="white"
        :header-text-variant="block.style.variants.headerText"
      >
        <div v-if="!headerSet">
          <div class="d-flex">
            <h5
              v-if="blockTitle"
              class="text-truncate mb-0"
            >
              {{ blockTitle }}

              <slot name="title-badge" />
            </h5>

            <b-button-group
              v-if="showOptions"
              size="sm"
              class="ml-auto"
            >
              <b-button
                v-if="block.options.showRefresh"
                v-b-tooltip.hover="{ title: $t('general.label.refresh'), container: '#body' }"
                variant="outline-light"
                class="d-flex align-items-center text-secondary d-print-none border-0"
                @click="$emit('refreshBlock')"
              >
                <font-awesome-icon :icon="['fa', 'sync']" />
              </b-button>

              <b-button
                v-if="block.options.magnifyOption || isBlockMagnified"
                v-b-tooltip.hover="{ title: isBlockMagnified ? '' : $t('general.label.magnify'), container: '#body' }"
                variant="outline-light"
                class="d-flex align-items-center text-secondary d-print-none border-0"
                @click="$root.$emit('magnify-page-block', isBlockMagnified ? undefined : magnifyParams)"
              >
                <font-awesome-icon :icon="['fas', isBlockMagnified ? 'times' : 'search-plus']" />
              </b-button>
            </b-button-group>
          </div>

          <b-card-text
            v-if="blockDescription"
            class="text-dark text-wrap mt-1"
          >
            {{ blockDescription }}
          </b-card-text>
        </div>

        <slot
          v-else
          name="header"
        />
      </b-card-header>

      <div
        v-if="toolbarSet"
      >
        <slot
          name="toolbar"
        />
      </div>

      <b-card-body
        class="p-0 flex-fill"
        :class="{ 'overflow-auto': scrollableBody }"
      >
        <slot
          name="default"
        />
      </b-card-body>

      <b-card-footer
        v-if="footerSet"
        class="p-0 bg-white border-top"
      >
        <slot
          name="footer"
        />
      </b-card-footer>
    </b-card>
  </div>
</template>
<script>
import base from './base.vue'
export default {
  name: 'CardWrap',
  extends: base,
}
</script>
