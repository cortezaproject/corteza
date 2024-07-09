<template>
  <div class="h-100">
    <b-card
      :id="blockID"
      no-body
      class="d-flex flex-column h-100 shadow-sm overflow-hidden position-static"
      :class="[blockClass, cardClass, customCSSClass]"
    >
      <b-card-header
        v-if="showHeader"
        :class="`border-bottom text-nowrap pr-2 ${headerClass}`"
        :header-text-variant="block.style.variants.headerText"
      >
        <div
          v-if="!headerSet"
          class="d-flex flex-column gap-1"
        >
          <div
            v-if="blockTitle || showOptions"
            class="d-flex"
          >
            <h4
              v-if="blockTitle"
              :title="blockTitle"
              class="text-truncate mb-0"
            >
              {{ blockTitle }}

              <slot name="title-badge" />
            </h4>

            <b-button-group
              v-if="showOptions"
              size="sm"
              class="ml-auto"
            >
              <b-button
                v-if="block.options.showRefresh"
                v-b-tooltip.noninteractive.hover="{ title: $t('general.label.refresh'), container: '#body' }"
                variant="outline-light"
                class="d-flex align-items-center text-secondary d-print-none border-0"
                @click="$emit('refreshBlock')"
              >
                <font-awesome-icon :icon="['fa', 'sync']" />
              </b-button>

              <b-button
                v-if="block.options.magnifyOption || isBlockMagnified"
                v-b-tooltip.noninteractive.hover="{ title: isBlockMagnified ? '' : $t('general.label.magnify'), container: '#body' }"
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
            :title="blockDescription"
            class="text-dark text-wrap"
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
        :body-class="bodyClass"
        class="p-0 flex-fill"
        :class="{ 'overflow-auto': scrollableBody }"
      >
        <slot
          name="default"
        />
      </b-card-body>

      <b-card-footer
        v-if="footerSet"
        class="p-0 bg-light"
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
