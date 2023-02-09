<template>
  <div class="h-100">
    <div
      class="card bg-transparent h-100"
      :class="blockClass"
    >
      <div
        v-if="showHeader"
        :class="`card-header bg-transparent border-0 text-nowrap px-3 text-${block.style.variants.headerText}`"
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

            <div
              v-if="showOptions"
              class="ml-auto"
            >
              <font-awesome-icon
                v-if="block.options.showRefresh"
                :icon="['fa', 'sync']"
                class="h6 text-secondary"
                role="button"
                @click="$emit('refreshBlock')"
              />

              <font-awesome-icon
                v-if="block.options.magnifyOption"
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
