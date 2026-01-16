<template>
  <div
    class="nav-sidebar"
    :class="{ 'mt-1': root }"
  >
    <div
      v-for="({page = {}, params = {}, children = []}) of items"
      :key="pageKey(page)"
      :class="{ 'mb-1': root }"
    >
      <div class="d-flex align-items-start pointer pb-1">
        <b-button
          variant="link"
          active-class="nav-active"
          exact-active-class="nav-exact-active"
          :title="page.title"
          :to="{ name: page.name || defaultRouteName, params }"
          class="nav-item d-flex align-items-center text-decoration-none rounded flex-grow-1 text-left pl-1 py-1 gap-1"
          @click="onItemClick()"
        >
          <template v-if="page.icon">
            <font-awesome-icon
              v-if="Array.isArray(page.icon)"
              :icon="page.icon"
              class="icon"
              style="height: 1rem; width: 1rem;"
            />
            <img
              v-else
              :src="page.icon"
              class="mr-1"
              style="height: 1rem; width: 1rem;"
            >
          </template>

          <label
            class="title pointer mb-0"
            :class="{ 'root': root }"
          >
            {{ page.title }}
          </label>
        </b-button>

        <b-button
          v-if="children.length"
          variant="outline-light"
          class="p-0 border-0 ml-auto"
          style="min-width: 2rem; min-height: 2rem;"
          @click="toggle(page)"
        >
          <font-awesome-icon
            v-if="!collapses[pageKey(page)]"
            class="text-dark"
            :icon="['fas', 'chevron-down']"
          />
          <font-awesome-icon
            v-else
            class="text-primary"
            :icon="['fas', 'chevron-up']"
          />
        </b-button>
      </div>

      <b-collapse
        v-if="children.length"
        :visible="collapses[pageKey(page)]"
      >
        <c-sidebar-nav-items
          :items="children"
          :start-expanded="startExpanded"
          :default-route-name="defaultRouteName"
          :root="false"
          class="py-1 ml-2"
          v-on="$listeners"
        />
      </b-collapse>
    </div>
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
    root: {
      type: Boolean,
      default: true,
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
          const px = this.pageKey(page)
          // Apply startExpanded only if page isn't currently expanded
          this.$set(this.collapses, px, this.startExpanded || page.expanded || this.showChildren({ params, children }))
        })
      },
    },
  },

  methods: {
    onItemClick () {
      if (window.innerWidth < 1024) {
        this.$root.$emit('close-sidebar')
      }
    },

    pageKey (p) {
      return p.pageID || p.name || p.title
    },

    toggle (p) {
      const px = this.pageKey(p)
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
.nav-sidebar {
  .nav-item {
    transition: background-color 0.2s ease-out;

    .icon {
      color: var(--black);
      transition: color 0.3s ease-in-out;
    }

    .title {
      color: var(--black);
      font-family: var(--font-regular) !important;
      transition: color 0.3s ease-in-out;
      text-align: left;
    }

    &:hover {
      background-color: var(--light);
    }
  }

  .nav-active {
    .icon {
      color: var(--primary);
    }

    .title {
      font-family: var(--font-medium) !important;
      color: var(--primary);
    }
  }

}
</style>

