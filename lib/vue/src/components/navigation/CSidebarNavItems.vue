<template>
  <div class="nav-sidebar">
    <b-button
      v-for="({page = {}, params = {}, children = []}) of items"
      :key="pageIndex(page)"
      variant="link"
      class="w-100 text-dark text-decoration-none p-0 pt-2 nav-item"
      active-class="nav-active"
      exact-active-class="nav-active"
      :to="{ name: page.name || defaultRouteName, params }"
    >
      <span
        class="d-inline-block w-75 text-nowrap text-truncate"
        @click="closeSidebar()"
      >
        <font-awesome-icon
          v-if="page.icon"
          class="icon"
          :icon="page.icon"
        />
        <span
          class="title"
        >
          {{ page.title }}
        </span>
      </span>

      <template
        v-if="children.length"
      >
        <b-button
          variant="link"
          class="p-0 float-right"
          :disabled="showChildren({ params, children })"
          @click.self.stop.prevent="toggle(page)"
        >
          <font-awesome-icon
            v-if="!collapses[pageIndex(page)]"
            class="pointer-none"
            :icon="['fas', 'chevron-down']"
          />
          <font-awesome-icon
            v-else
            class="pointer-none"
            :icon="['fas', 'chevron-up']"
          />
        </b-button>

        <b-collapse
          :visible="collapses[pageIndex(page)] || showChildren({ params, children })"
          @click.stop.prevent
        >
          <c-sidebar-nav-items
            class="ml-2"
            :items="children"
            :start-expanded="startExpanded"
            :default-route-name="defaultRouteName"
            v-on="$listeners"
          />
        </b-collapse>
      </template>
    </b-button>
  </div>
</template>

<script>
export default {
  name: 'CSidebarNavItems',

  props: {
    /*
    * {
        page: { name, title }
        params: {...}
      }
    */
    items: {
      type: Array,
      required: true,
      default: () => [],
    },
    defaultRouteName: {
      type: String,
      required: true,
    },
    startExpanded: {
      type: Boolean,
      required: false,
    },
  },

  data () {
    return {
      collapses: {},
    }
  },

  watch: {
    items: {
      immediate: true,
      handler (items = []) {
        items.forEach(({ page }) => {
          const px = this.pageIndex(page)
          // Apply startExpanded only if page isn't currently expanded
          if (!this.collapses[px]) {
            this.$set(this.collapses, px, this.startExpanded)
          }
        })
      },
    },
  },

  methods: {
    closeSidebar () {
      if (window.innerWidth < 576) {
        this.$root.$emit('close-sidebar')
      }
    },

    pageIndex (p) {
      return p.pageID || p.name || p.title
    },

    toggle (p) {
      const px = this.pageIndex(p)
      this.$set(this.collapses, px, !this.collapses[px])
    },

    // Recursively check for child pages that are open, so that parents can open aswell
    showChildren ({ params = {}, children = [] }) {
      const partialParamsMatch = Object.entries(params).some(([key, value]) => {
        return this.$route.params[key] === value
      })

      if (partialParamsMatch) {
        return partialParamsMatch
      }

      return children.map(c => this.showChildren(c)).some(isOpen => isOpen)
    },
  },
}
</script>

<style scoped lang="scss">
// This has to be there, so chevrons are clickable inside the button
.pointer-none {
  pointer-events: none;
}

.svg-inline--fa {
  width: 30px;
}

.nav-item > span > {
  .title {
    font-family: 'Poppins-Regular'
  }
}

.nav-active > span > {
  .icon {
    color: #4D7281;
  }

  .title {
    font-family: 'Poppins-SemiBold'
  }
}

.nav-item {
  text-align: left;
}

[dir="rtl"] {
  .nav-item {
    text-align: right;
  }
}
</style>
