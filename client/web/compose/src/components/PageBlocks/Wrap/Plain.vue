<template>
  <div class="h-100">
    <div
      class="card bg-transparent h-100 overflow-hidden"
      :class="blockClass"
    >
      <div
        v-if="showHeader"
        :class="`card-header bg-transparent border-0 text-nowrap pl-3 pr-2 mr-1 text-${block.style.variants.headerText}`"
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
              {{ block.title }}

              <slot name="title-badge" />
            </h5>

            <b-button-group
              v-if="showOptions"
              size="sm"
              class="ml-auto"
            >
              <b-button
                v-if="block.options.showRefresh"
                :title="$t('general.label.refresh')"
                variant="outline-light"
                class="d-flex align-items-center text-primary d-print-none border-0"
                @click="$emit('refreshBlock')"
              >
                <font-awesome-icon :icon="['fa', 'sync']" />
              </b-button>

              <b-button
                v-if="block.options.magnifyOption"
                :title="isBlockOpened ? '' : $t('general.label.magnify')"
                variant="outline-light"
                class="d-flex align-items-center text-primary d-print-none border-0"
                @click="$root.$emit('magnify-page-block', isBlockOpened ? undefined : { blockID: block.blockID })"
              >
                <font-awesome-icon :icon="['fas', isBlockOpened ? 'times' : 'search-plus']" />
              </b-button>
            </b-button-group>
          </div>

          <b-card-text
            v-if="block.description"
            class="text-dark text-wrap mt-1"
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
        class="card-footer bg-transparent p-0 overflow-hidden border-top"
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
