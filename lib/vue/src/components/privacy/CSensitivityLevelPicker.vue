<template>
  <vue-select
    key="type"
    :value="_value"
    :disabled="_disabled"
    :options="sensitivityLevels"
    :get-option-label="getLabel"
    :placeholder="placeholder"
    :reduce="l => l.sensitivityLevelID"
    append-to-body
    class="bg-white"
    @input="onInput"
  />
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'
import { VueSelect } from 'vue-select'

export default {
  components: {
    VueSelect,
  },

  props: {
    value: {
      type: String,
      default: '',
    },

    options: {
      type: Array,
      required: true
    },

    placeholder: {
      type: String,
      default: '',
    },

    // ID of sensitivityLevel with the maximum allowed level
    maxLevel: {
      type: String,
      default: undefined
    },

    disabled: {
      type: Boolean,
      default: false
    },
  },

  computed: {
    _value () {
      return this.value === NoID ? undefined : this.value
    },

    _disabled () {
      return this.disabled || this.maxLevel === NoID
    },

    sensitivityLevels () {
      if (this.maxLevel === NoID) {
        return []
      }

      if (this.maxLevel) {
        const maxLevelConnection = this.options.find(({ sensitivityLevelID }) => sensitivityLevelID === this.maxLevel)
        if (maxLevelConnection) {
          return this.options.filter(({ level }) => level <= maxLevelConnection.level)
        }
      }

      return this.options
    },
  },

  watch: {
    // If sensitivityLevels change because of maxLevel change, then assert if value is still allowed, otherwise reset it
    sensitivityLevels: {
      immediate: true,
      handler () {
        const isValueCompatible = this.sensitivityLevels.some(({ sensitivityLevelID }) => sensitivityLevelID === this.value)
        if (!isValueCompatible) {
          this.$emit('input', NoID)
        }
      }
    },
  },

  methods: {
    getLabel ({ handle, meta = {} }) {
      return meta.name || handle
    },

    onInput (sensitivityLevelID) {
      this.$emit('input', sensitivityLevelID ? sensitivityLevelID : NoID)
    }
  }
}
</script>

<style>

</style>
