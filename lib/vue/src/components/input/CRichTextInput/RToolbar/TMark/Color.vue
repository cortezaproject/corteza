<template>
  <div>
    <b-button
      :id="`color-popover-${format.type}`"
      variant="link"
      class="text-dark font-weight-bold text-decoration-none"
      @click.stop.prevent="showPicker"
    >
      <span
        class="icon"
        :class="typeStyle"
        :style="{
          backgroundColor: background ? selectedColor : 'transparent',
          borderColor: background ? 'transparent' : selectedColor,
        }"
      >
        A
      </span>
    </b-button>

    <c-input-color-picker
      ref="picker"
      class="d-none"
      :value="selectedColor"
      :default-value="getDefaultColor()"
      :show-text="false"
      :width="'0px'"
      :height="'0px'"
      @input="applyFromPicker"
    />
  </div>
</template>

<script>
import CInputColorPicker from '../../../CInputColorPicker.vue'
import base from './base.vue'

export default {
  name: 'TMarkColor',

  components: {
    CInputColorPicker,
  },

  extends: base,

  props: {
    background: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      selectedColor: this.getDefaultColor(),
    }
  },

  computed: {
    typeStyle () {
      return `${this.format.type}-icon`
    },
  },

  methods: {
    getComputedColor (cssVar) {
      try {
        const computedStyle = getComputedStyle(document.documentElement)
        return computedStyle.getPropertyValue(cssVar).trim()
      } catch (error) {
        return null
      }
    },

    showPicker () {
      if (this.$refs.picker && this.$refs.picker.openMenu) {
        this.$refs.picker.openMenu()
      }
    },

    applyFromPicker (val) {
      if (!val) return
      this.selectedColor = val
      this.onClick(this.format.type, { color: val })
    },

    getDefaultColor () {
      return this.background ? this.getComputedColor('--white') : this.getComputedColor('--dark')
    },
  },
}
</script>

<style lang="scss" scoped>
// Since we don't have icons for color related formatters,
// this extra styling makes up for that
.icon {
  padding: 0 2px;

  &.color-icon {
    border-bottom: 2px solid var(--dark);
  }

  &.background-icon {
    padding: 2px 3px;
  }
}
</style>

