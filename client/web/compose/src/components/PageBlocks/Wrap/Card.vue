<template>
  <div class="h-100">
    <b-card
      no-body
      class="h-100 shadow overflow-hidden"
      :class="blockClass"
    >
      <b-card-header
        v-if="showHeader"
        class="border-0 text-nowrap pl-3 pr-2 mr-1"
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

            <b-button-group
              v-if="showOptions"
              size="sm"
              class="ml-auto"
            >
              <b-button
                v-if="block.options.showRefresh"
                :title="$t('general.label.refresh')"
                variant="outline-light"
                class="d-flex align-items-center text-secondary d-print-none border-0"
                @click="$emit('refreshBlock')"
              >
                <font-awesome-icon :icon="['fa', 'sync']" />
              </b-button>

              <b-button
                v-if="block.options.magnifyOption"
                :title="isBlockOpened ? '' : $t('general.label.magnify')"
                variant="outline-light"
                class="d-flex align-items-center text-secondary d-print-none border-0"
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
import base from './base.vue'
export default {
  name: 'CardWrap',
  extends: base,
}
</script>
