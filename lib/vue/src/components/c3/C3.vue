<template>
  <div class="layout">
    <aside class="sidebar p-2">
      <h5
        class="border-bottom">
        C3: Component Catalogue
      </h5>
      <component-list
        :catalogue="catalogue"
        @select="setCurrent($event)"
      />
    </aside>

    <main
      class="p-5"
    >
      <component
        :is="current.component"
        v-if="current"
        v-bind="current.props"
      />
      <p
        v-else
        class="text-center"
      >
        Select a component from the C3 Catalogue and start hacking :)
      </p>
    </main>

    <div
      class="controls px-5 py-2 mt-2"
      v-if="current"
    >
      <div
        class="control-group mr-2"
        v-for="(cg, g) in controlGroups"
        :key="`control-group-${g}`"
      >
        <h3>
          Controls
        </h3>
        <component
          :is="c.component"
          v-for="(c, i) in cg"
          :key="i"
          :value="c.value(current.props)"
          v-bind="c.props"
          @update="c.update(current.props, $event)"
        />
      </div>

      <div
        v-if="current.scenarios"
        class="control-group float-right"
      >
        <h3>
          Pre-set controls
        </h3>
        <ul
          class="pl-0"
        >
          <li
            v-for="(s, i) in current.scenarios"
            :key="i"
            class="list-unstyled scenario"
            @click="setScenario(s)"
          >
            {{ s.label }}
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>
<script>
import ComponentList from './ComponentList.vue'

export default {
  name: 'C3',

  components: {
    ComponentList,
  },

  props: {
    catalogue: {
      required: true,
      type: Object,
    },
  },

  computed: {
    controlGroups () {
      if (this.current.controls.length === 0) {
        return []
      }

      if (Array.isArray(this.current.controls[0])) {
        // already grouped
        return this.current.controls
      }

      // make one virtual group holding all controls
      return [this.current.controls]
    },
  },

  data () {
    return {
      current: undefined,
    }
  },

  methods: {
    setCurrent (component) {
      this.current = { props: {}, ...component }
      this.setScenario(this.current)
    },

    setScenario ({ props = {}, controls = [] }) {
      // create missing props from controls
      const apply = (c, props) => c.update(props, c.value(props) || null)
      controls.forEach(c => {
        if (Array.isArray(c)) {
          c.forEach(c => apply(c, props))
        } else {
          apply(c, props)
        }
      })

      this.current.props = props
    },
  },
}
</script>
<style lang="scss">
.layout {
  height: 100vh;
  width: 100vw;

  display: grid;
  grid-template-rows: auto 400px;
  grid-template-columns: 300px auto;
  align-content: stretch;
  grid-template-areas:
    "side main"
    "side controls"
  ;

  aside {
    grid-area: side;
    overflow: auto;
  }

  main {
    grid-area: main;
    background-image: linear-gradient(
      135deg,
      #F9FAFB 21.43%,
      #FFFFFF 21.43%,
      #FFFFFF 50%,
      #F9FAFB 50%,
      #F9FAFB 71.43%,
      #FFFFFF 71.43%,
      #FFFFFF 100%
    );
    background-size: 35.00px 35.00px;
    overflow: auto;
  }

  .controls {
    grid-area: controls;
    overflow: auto;

    .control-group {
      float: left;
    }
  }

  .scenario {
    cursor: pointer;
  }
}
</style>
