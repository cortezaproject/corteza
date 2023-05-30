<template>
  <wrap
    :scrollable-body="false"
    v-bind="$props"
    v-on="$listeners"
  >
    <div class="h-100 w-100 card overflow-hidden bg-transparent">
      <b-nav
        v-bind="{
          tabs: options.display.appearance === 'tabs',
          pills: options.display.appearance === 'pills',
          small: options.display.appearance === 'small',
          justified: options.display.justify === 'justify'
        }"
        :align="options.display.alignment"
        class="border-0 h-100"
      >
        <b-nav-item
          v-for="(navItem, index) in options.navigationItems"
          :key="`navItem-${index}`"
          :disabled="!navItem.options.enabled"
          :style="{ order: index, color: navItem.options.textColor, background: navItem.options.backgroundColor, justifyContent: options.display.alignment }"
          :link-attrs="{
            style: `color: ${navItem.options.textColor}`,
          }"
          link-classes="h-100 w-100 d-flex align-items-center justify-content-center"
          :href="generateHrefAttributeLink(navItem)"
          :to="generateToAttributeLink(navItem)"
          :target="selectTargetOption(navItem.options.item.target)"
          class="d-flex align-items-center"
        >
          <template v-if="navItem.type === 'dropdown' || isComposeDropdownPage(navItem)">
            <b-button
              :id="`dropdown-popover-${index}-${block.blockID}`"
              class="text-decoration-none p-0 w-100 h-100"
              variant="link"
              :style="{ color: navItem.options.textColor, background: navItem.options.backgroundColor }"
            >
              {{ displayDropdownText(navItem) }}
              <span class="ml-1">
                <font-awesome-icon
                  :icon="['fas', 'chevron-down']"
                  size="sm"
                />
              </span>
            </b-button>

            <b-popover
              ref="dropdown-popover"
              :target="`dropdown-popover-${index}-${block.blockID}`"
              :placement="navItem.options.item.align"
              delay="0"
              boundary="window"
              triggers="click blur"
            >
              <template
                v-if="navItem.type === 'dropdown'"
              >
                <div
                  v-for="(dropdown, dIndex) in navItem.options.item.dropdown.items"
                  :key="`dropdown-${dIndex}`"
                >
                  <a
                    class="dropdown-item"
                    :href="dropdown.url | checkValidURL"
                    :disabled="navItem.options.disabled"
                    :target="selectTargetOption(dropdown.target)"
                    :style="{ order: dIndex * 2 }"
                  >
                    {{ dropdown.label }}
                  </a>

                  <hr
                    v-if="dropdown.delimiter"
                    class="my-1"
                    :style="{ order: dIndex + 1 }"
                  >
                </div>
              </template>

              <template v-else>
                <b-link
                  :to="{ name: 'page', params: { pageID: navItem.options.item.pageID } }"
                  :target="selectTargetOption(navItem.options.item.target)"
                  :disabled="navItem.options.disabled"
                  class="dropdown-item"
                  style="white-space: normal;"
                >
                  {{ navItem.options.item.label }}
                </b-link>

                <hr
                  class="my-1"
                >

                <div
                  v-for="(dropdown, dIndex) in getSubPages(navItem.options.item.pageID)"
                  :key="`dropdown-${dIndex}`"
                >
                  <b-link
                    :to="{ name: 'page', params: { pageID: dropdown.pageID } }"
                    :target="selectTargetOption(navItem.options.item.target)"
                    :style="{ order: dIndex * 2 }"
                    :disabled="navItem.options.disabled"
                    class="dropdown-item"
                    style="white-space: normal;"
                  >
                    {{ dropdown.title }}
                  </b-link>
                </div>
              </template>
            </b-popover>
          </template>

          <template v-else>
            {{ navItem.options.item.label }}
          </template>
        </b-nav-item>
      </b-nav>
    </div>
  </wrap>
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'
import { mapGetters } from 'vuex'
import base from '../base'

export default {
  extends: base,

  computed: {
    ...mapGetters({
      pages: 'page/set',
    }),
  },

  methods: {
    isComposeDropdownPage (navItem) {
      return (navItem.type === 'compose' && navItem.options.item.displaySubPages)
    },

    getSubPages (selfID) {
      return this.pages.filter(value => value.selfID === selfID && value.moduleID === NoID) || []
    },

    selectTargetOption (target) {
      switch (target) {
        case 'sameTab':
          return '_self'
        case 'newTab':
          return '_blank'
      }
    },

    displayDropdownText (navItem) {
      if (navItem.type === 'dropdown') {
        return navItem.options.item.dropdown.label
      }

      return navItem.options.item.label
    },

    generateToAttributeLink (navItem) {
      if (['dropdown', 'text-section'].includes(navItem.type) || this.isComposeDropdownPage(navItem)) {
        return
      }

      let to = ''

      if (navItem.type === 'compose') {
        const pageID = navItem.options.item.pageID

        to = { name: 'page', params: { pageID } }
      }

      return to
    },

    generateHrefAttributeLink (navItem) {
      if (['dropdown', 'text-section'].includes(navItem.type) || this.isComposeDropdownPage(navItem)) {
        return
      }

      return navItem.type === 'url' ? this.$options.filters.checkValidURL(navItem.options.item.url) : ''
    },
  },
}
</script>
