<template>
  <wrap
    :scrollable-body="false"
    v-bind="$props"
    v-on="$listeners"
  >
    <div
      ref="container"
      class="h-100 w-100 card overflow-hidden bg-transparent"
    >
      <div
        v-if="collapsed"
        class="d-flex align-items-center justify-content-start h-100 px-2"
      >
        <b-button
          :id="`navigation-menu-${block.blockID}`"
          data-test-id="button-navigation-menu"
          variant="outline-extra-light"
          class="d-flex align-items-center text-primary text-left border-0"
        >
          <font-awesome-icon :icon="['fas', 'bars']" />

          <span
            v-if="currentItem"
            data-test-id="navigation-menu-current"
            class="ml-2 text-body"
          >
            {{ displayDropdownText(currentItem) }}
          </span>
        </b-button>
      </div>

      <!-- Stays rendered while collapsed so it can be measured -->
      <b-nav
        ref="nav"
        v-bind="{
          tabs: options.display.appearance === 'tabs',
          pills: options.display.appearance === 'pills',
          small: options.display.appearance === 'small',
          justified: options.display.justify === 'justify'
        }"
        :align="options.display.alignment"
        class="border-0 h-100 overflow-auto"
        :class="{ 'flex-nowrap': isSmallScreen, 'position-absolute w-100 invisible': collapsed }"
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
          class="align-items-center"
          :class="{ 'd-flex': !hiddenItems.includes(index), 'd-none': hiddenItems.includes(index), 'text-nowrap': isSmallScreen }"
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

        <li
          v-if="hiddenItems.length"
          class="nav-item d-flex align-items-center"
          :style="{ order: options.navigationItems.length, flex: '0 0 auto' }"
        >
          <b-button
            :id="`navigation-menu-${block.blockID}`"
            data-test-id="button-navigation-more"
            variant="link"
            class="text-nowrap"
          >
            {{ $t('navigation.more', { count: hiddenItems.length }) }}
            <font-awesome-icon
              :icon="['fas', 'chevron-down']"
              size="sm"
              class="ml-1"
            />
          </b-button>
        </li>
      </b-nav>

      <b-popover
        v-if="collapsed || hiddenItems.length"
        :key="collapsed ? 'menu' : 'more'"
        :target="`navigation-menu-${block.blockID}`"
        placement="bottom"
        delay="0"
        boundary="window"
        triggers="click blur"
        custom-class="navigation-menu"
      >
        <template
          v-for="(navItem, index) in menuItems"
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

  data () {
    return {
      // On small screens, items that do not fit in one row go under More,
      // or all of them into a menu when not even one fits
      collapsed: false,
      hiddenItems: [],
    }
  },

  computed: {
    ...mapGetters({
      pages: 'page/set',
    }),

    // The item linking to the page being viewed, its exact layout first
    currentIndex () {
      const { pageID } = this.$route.params
      const { layoutID } = this.$route.query
      const items = this.options.navigationItems
      const onPage = ({ type, options }) => type === 'compose' && options.item.pageID === pageID

      const exact = items.findIndex(navItem => onPage(navItem) && navItem.options.item.pageLayoutID === layoutID)
      return exact >= 0 ? exact : items.findIndex(onPage)
    },

    currentItem () {
      return this.options.navigationItems[this.currentIndex]
    },

    menuItems () {
      const items = this.options.navigationItems
      return this.collapsed ? items : items.filter((_, index) => this.hiddenItems.includes(index))
    },
  },

  watch: {
    isSmallScreen () {
      this.$nextTick(this.checkNavFits)
    },

    currentIndex () {
      this.$nextTick(this.checkNavFits)
    },

    'options.navigationItems': {
      deep: true,
      handler () {
        this.$nextTick(this.checkNavFits)
      },
    },
  },

  mounted () {
    this.navObserver = new ResizeObserver(this.checkNavFits)
    this.navObserver.observe(this.$refs.container)
    this.$nextTick(this.checkNavFits)
  },

  beforeDestroy () {
    this.navObserver.disconnect()
  },

  methods: {
    async checkNavFits () {
      if (this.measuring) {
        this.remeasure = true
        return
      }

      this.measuring = true
      try {
        await this.fitNav()
      } finally {
        this.measuring = false
      }

      if (this.remeasure) {
        this.remeasure = false
        this.checkNavFits()
      }
    },

    // Keeps the items that fit in one row, the current page among them
    async fitNav () {
      // b-nav is functional, so the ref is its element
      const { nav } = this.$refs
      if (!nav || !this.isSmallScreen) {
        this.collapsed = false
        this.hiddenItems = []
        return
      }

      if (this.hiddenItems.length) {
        this.hiddenItems = []
        await this.$nextTick()
      }

      const overflows = () => nav.scrollWidth > nav.clientWidth + 1
      if (!overflows()) {
        this.collapsed = false
        return
      }

      const widths = [...nav.children].map(item => item.offsetWidth)
      const current = this.currentIndex
      const room = nav.clientWidth - (this.moreWidth || 0)
      const shown = current < 0 ? [] : [current]
      let used = current < 0 ? 0 : widths[current]

      for (let index = 0; index < widths.length; index++) {
        if (index === current) continue
        if (used + widths[index] > room) break
        shown.push(index)
        used += widths[index]
      }

      if (!shown.length || used > room) {
        this.collapsed = true
        return
      }

      this.collapsed = false
      this.hiddenItems = widths.map((_, index) => index).filter(index => !shown.includes(index))
      await this.$nextTick()

      // More can be wider than the room kept for it
      while (overflows() && shown.length > 1) {
        const last = Math.max(...shown.filter(index => index !== current))
        shown.splice(shown.indexOf(last), 1)
        this.hiddenItems = [...this.hiddenItems, last].sort((a, b) => a - b)
        await this.$nextTick()
      }

      const more = nav.querySelector('[data-test-id="button-navigation-more"]')
      this.moreWidth = more ? more.parentElement.offsetWidth : 0

      if (overflows()) {
        this.hiddenItems = []
        this.collapsed = true
      }
    },

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

  .dropdown-item {
    white-space: normal;
  }
}
</style>
