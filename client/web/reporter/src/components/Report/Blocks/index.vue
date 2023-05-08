<template>
  <wrap
    v-if="block"
    v-bind="$props"
    v-on="$listeners"
  >
    <template
      #header
    >
      <div
        v-if="block.title || block.description"
        class="px-3"
        style="padding-top: 0.75rem; padding-bottom: 0.75rem;"
      >
        <h5
          v-if="block.title"
          class="text-primary text-truncate mb-0"
        >
          {{ block.title }}
        </h5>

        <b-card-text
          v-if="block.description"
          class="text-dark text-truncate"
          :class="{ 'mt-1': block.title }"
        >
          {{ block.description }}
        </b-card-text>
      </div>
    </template>

    <template #default>
      <split
        v-if="showDisplayElements"
        ref="split"
        :direction="block.layout"
        :gutter-size="12"
        class="h-100"
        @onDragEnd="setDisplayElementSizes"
      >
        <split-area
          v-for="(element, displayElementIndex) in block.elements"
          :key="displayElementIndex"
          :size="element.meta.size"
          :min-size="0"
          :class="{
            'overflow-hidden h-100': element.kind !== 'Text',
            'w-100': block.elements.length === 1,
          }"
          class="position-relative"
        >
          <div
            v-if="processing"
            class="d-flex align-items-center justify-content-center h-100"
          >
            <b-spinner />
          </div>

          <display-element
            v-else
            :display-element="element"
            :labels="{
              previous: $t('display-element:table.view.previous'),
              next: $t('display-element:table.view.next'),
            }"
            @update="updateDataframes({ displayElementIndex, definition: $event })"
          />
        </split-area>
      </split>
    </template>
  </wrap>
</template>

<script>
import Wrap from './Wrap'
import { Split, SplitArea } from 'vue-split-panel'
import DisplayElement from './DisplayElements/Viewers'
import { reporter } from '@cortezaproject/corteza-js'

export default {
  name: 'Block',

  components: {
    Split,
    SplitArea,
    Wrap,
    DisplayElement,
  },

  props: {
    index: {
      type: Number,
      default: () => -1,
    },

    block: {
      type: Object,
      required: true,
    },

    scenario: {
      type: Object,
      default: () => ({}),
    },

    reportID: {
      type: String,
      required: false,
      default: '0',
    },
  },

  data () {
    return {
      processing: false,

      dataframes: {},

      showDisplayElements: false,
    }
  },

  watch: {
    'block.blockID': {
      immediate: true,
      handler () {
        this.renderBlock()
      },
    },

    'block.elements.length': {
      handler (length, oldLength) {
        const addedOrRemoved = length !== oldLength && oldLength

        if (addedOrRemoved) {
          const defaultSize = Math.floor(100 / length)

          // Reset sizes to default if element was added or removed
          this.block.elements = this.block.elements.map(e => {
            e.meta.size = defaultSize
            return e
          })
        }

        this.renderBlock()
      },
    },
  },

  methods: {
    renderBlock () {
      const { elements = [] } = this.block || {}
      if (elements.length) {
        this.runReport()

        // Hack around split not rerendering
        this.showDisplayElements = false
        this.$nextTick().then(() => {
          this.showDisplayElements = true
        })
      }
    },

    setDisplayElementSizes (sizes = []) {
      sizes.forEach((size, index) => {
        this.block.elements[index].meta.size = size
      })

      this.$emit('item-updated', this.index)
    },

    getScenarioDefinition (element) {
      const scenarioDefinition = {}

      // Generate filter for each load datasource
      if (this.scenario.filters) {
        element.options.datasources.forEach(({ name }) => {
          scenarioDefinition[name] = {
            ref: name,
            filter: this.scenario.filters[name] || {},
          }
        })
      }

      return scenarioDefinition
    },

    runReport () {
      this.processing = true
      this.dataframes = {}
      const frames = []

      this.block.elements.forEach((element, key) => {
        element = reporter.DisplayElementMaker(element)

        if (element && element.kind !== 'Text') {
          if (element.elementID === '0') {
            element.elementID = `${key}`
          }

          const { dataframes = [] } = element.reportDefinitions(this.getScenarioDefinition(element))

          frames.push(...dataframes.filter(({ source }) => source))
        }
      })

      if (frames.length) {
        this.$SystemAPI.reportRun({ frames, reportID: this.reportID })
          .then(({ frames = [] }) => {
            this.block.elements = this.block.elements.map((element, key) => {
              if (element.elementID === '0') {
                element.elementID = `${key}`
              }
              const dataframes = frames.filter(({ name }) => name === element.elementID)
              return { ...element, dataframes }
            })
          }).catch((e) => {
            this.toastErrorHandler(this.$t('notification:report.runFailed'))(e)
          }).finally(() => {
            this.processing = false
          })
      } else {
        this.processing = false
      }
    },

    updateDataframes ({ displayElementIndex, definition }) {
      const element = reporter.DisplayElementMaker(this.block.elements[displayElementIndex])
      const frames = []

      if (element && element.kind !== 'Text') {
        const scenarioDefinition = this.getScenarioDefinition(element)
        Object.entries(definition).forEach(([key, value]) => {
          definition[key] = { ...value, ...scenarioDefinition[key] }
        })

        const { dataframes = [] } = element.reportDefinitions(definition)

        frames.push(...dataframes.filter(({ source }) => source))

        if (frames.length) {
          this.$SystemAPI.reportRun({ frames, reportID: this.reportID })
            .then(({ frames = [] }) => {
              this.block.elements.find(({ elementID }) => elementID === element.elementID).dataframes = frames
            }).catch((e) => {
              this.toastErrorHandler(this.$t('notification:report.runFailed'))(e)
            })
        }
      }
    },
  },
}
</script>

<style lang="scss">
.split .gutter {
  background-color: transparent;
}
</style>
