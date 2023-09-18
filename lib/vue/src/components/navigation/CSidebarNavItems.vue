<template>
  <div class="nav-sidebar">
    <b-button
      v-for="({page = {}, params = {}, children = []}) of items"
      :key="pageIndex(page)"
      variant="link"
      class="w-100 text-dark text-decoration-none p-0 pt-2 nav-item"
      active-class="nav-active"
      exact-active-class="nav-active"
      :title="page.title"
      :to="{ name: page.name || defaultRouteName, params }"
    >
      <span
        class="d-inline-block w-75 text-nowrap text-truncate"
        @click="closeSidebar()"
      >
        <template
          v-if="page.icon"
        >
          <font-awesome-icon
            v-if="Array.isArray(page.icon)"
            class="icon"
            :icon="page.icon"
          />
          <template v-else>
            <img
              :src="page.icon"
              class="mr-1"
              style="height: 1.5em; width: 1.5em;"
            />
          </template>
        </template>
        <label
          class="title mb-0 pointer"
        >
          {{ page.title }}
        </label>
      </span>

      <template
        v-if="children.length"
      >
        <b-button
          variant="outline-light"
          size="sm"
          class="text-primary p-0 border-0 float-right mr-1"
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
          :visible="collapses[pageIndex(page)]"
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
        items.forEach(({ page, params, children }) => {
          const px = this.pageIndex(page)
          // Apply startExpanded only if page isn't currently expanded
          this.$set(this.collapses, px, this.startExpanded || page.expanded || this.showChildren({ params, children }))
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

    // Recursively check for child pages that are open, so that parents can open as well
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

.nav-item > span {
  .title {
    color: var(--tertiary);
  }
}

.nav-item:hover > span {
  .title {
    color: var(--primary);
    transition: color 0.25s;
  }
}

.nav-active > span > {
  .icon {
    color: var(--primary)
  }

  .title {
    font-family: 'Poppins-SemiBold';
    color: var(--primary);
    transition: color 0.5s
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
