<template>
  <div class="h-100">
    <div
      class="d-flex flex-column card bg-transparent h-100 overflow-hidden position-static"
      :class="[blockClass, cardClass]"
    >
      <div
        v-if="showHeader"
        :class="`card-header bg-transparent border-0 text-nowrap pl-3 pr-2 text-${block.style.variants.headerText}`"
      >
        <div
          v-if="!headerSet"
        >
          <div
            class="d-flex"
          >
            <h5
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
      </div>

      <div
        v-if="toolbarSet"
      >
        <slot
          name="toolbar"
        />
      </div>

      <div
        class="card-body p-0 flex-fill"
        :class="{ 'overflow-auto': scrollableBody }"
        style="flex-shrink: 10;"
      >
        <slot
          name="default"
        />
      </div>

      <b-card-footer
        v-if="footerSet"
        class="card-footer bg-transparent p-0 border-top"
      >
        <slot
          name="footer"
        />
      </b-card-footer>
    </div>
  </div>
</template>
<script>
import base from './base.vue'
export default {
  name: 'PlainWrap',
  extends: base,
}
</script>
