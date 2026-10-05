<template>
  <wrap
    :scrollable-body="false"
    v-bind="$props"
    v-on="$listeners"
  >
    <div class="h-100 w-100 card overflow-hidden bg-transparent">
      <div
        v-if="isSmallScreen"
        class="d-flex align-items-center justify-content-end h-100 px-2"
      >
        <b-button
          :id="`navigation-menu-${block.blockID}`"
          data-test-id="button-navigation-menu"
          variant="outline-extra-light"
          class="text-primary border-0"
        >
          <font-awesome-icon :icon="['fas', 'bars']" />
        </b-button>

        <b-popover
          :target="`navigation-menu-${block.blockID}`"
          placement="bottom"
          delay="0"
          boundary="window"
          triggers="click blur"
          custom-class="navigation-menu"
        >
          <template
            v-for="(navItem, index) in options.navigationItems"
          >
            <template v-if="navItem.type === 'dropdown'">
              <h6
                :key="`navItem-${index}`"
                class="dropdown-header"
              >
                {{ navItem.options.item.dropdown.label }}
              </h6>

              <a
                v-for="(dropdown, dIndex) in navItem.options.item.dropdown.items"
                :key="`navItem-${index}-${dIndex}`"
                :href="dropdown.url | checkValidURL"
                :target="selectTargetOption(dropdown.target)"
                class="dropdown-item pl-4"
              >
                {{ dropdown.label }}
              </a>
            </template>

            <template v-else-if="isComposeDropdownPage(navItem)">
              <b-link
                :key="`navItem-${index}`"
                :to="{ name: 'page', params: { pageID: navItem.options.item.pageID } }"
                :target="selectTargetOption(navItem.options.item.target)"
                :disabled="!navItem.options.enabled"
                class="dropdown-item"
              >
                {{ navItem.options.item.label }}
              </b-link>

              <b-link
                v-for="(subPage, dIndex) in getSubPages(navItem.options.item.pageID)"
                :key="`navItem-${index}-${dIndex}`"
                :to="{ name: 'page', params: { pageID: subPage.pageID } }"
                :target="selectTargetOption(navItem.options.item.target)"
                :disabled="!navItem.options.enabled"
                class="dropdown-item pl-4"
              >
                {{ subPage.title }}
              </b-link>
            </template>

            <span
              v-else-if="navItem.type === 'text-section'"
              :key="`navItem-${index}`"
              class="dropdown-item-text"
            >
              {{ navItem.options.item.label }}
            </span>

            <b-link
              v-else
              :key="`navItem-${index}`"
              :href="generateHrefAttributeLink(navItem)"
              :to="generateToAttributeLink(navItem)"
              :target="selectTargetOption(navItem.options.item.target)"
              :disabled="!navItem.options.enabled"
              :style="{ color: navItem.options.textColor }"
              class="dropdown-item"
            >
              {{ navItem.options.item.label }}
            </b-link>
          </template>
        </b-popover>
      </div>

      <b-nav
        v-else
        v-bind="{
          tabs: options.display.appearance === 'tabs',
          pills: options.display.appearance === 'pills',
          small: options.display.appearance === 'small',
          justified: options.display.justify === 'justify'
        }"
        :align="options.display.alignment"
        class="border-0 h-100 overflow-auto"
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
              class="p-0 w-100 h-100"
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
                  v-if="getSubPages(navItem.options.item.pageID).length > 0"
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
import smallScreen from '../../../mixins/smallScreen'

export default {
  extends: base,

  mixins: [
    smallScreen,
  ],

  computed: {
    ...mapGetters({
      pages: 'page/set',
    }),
  },

  methods: {
    isComposeDropdownPage (navItem) {
      return (navItem.type === 'compose' && navItem.options.item.displaySubPages)
    },

    getSubPages (pageID) {
      return this.pages.filter(value => value.selfID === pageID && value.moduleID === NoID) || []
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

      if (navItem.type === 'compose') {
        const pageID = navItem.options.item.pageID
        const pageLayoutID = navItem.options.item.pageLayoutID
        const moduleID = navItem.options.item.moduleID

        // Handle modal context - update modal layout ID
        if (this.inModal && pageID === this.$route.query.recordPageID) {
          return {
            ...this.$route,
            query: {
              ...this.$route.query,
              modalLayoutID: pageLayoutID,
            },
          }
        }

        // Determine if we're staying on the same page
        const isSamePage = pageID === this.$route.params.pageID

        // Handle record pages
        if (moduleID) {
          return isSamePage
            ? {
                ...this.$route,
                query: { layoutID: pageLayoutID },
              }
            : {
                name: 'page.record.create',
                params: { pageID },
                query: { layoutID: pageLayoutID },
              }
        }

        // Handle regular pages
        return isSamePage
          ? {
              ...this.$route,
              query: { layoutID: pageLayoutID },
            }
          : {
              name: 'page',
              params: { pageID },
              query: { layoutID: pageLayoutID },
            }
      }
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

<style lang="scss" scoped>
.nav-link:hover {
  text-decoration: underline !important;
}
</style>

<style lang="scss">
.navigation-menu .popover-body {
  max-height: 70vh;
  overflow-y: auto;
  padding: 0.5rem 0;
}
</style>
