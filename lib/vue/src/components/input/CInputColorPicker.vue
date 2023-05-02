<template>
  <div>
    <b-button
      :style="`color: ${pickedColor}; fill: ${pickedColor};`"
      class="p-0 rounded-circle bg-white border-white shadow-none"
      @click="toggleMenu"
    >
      <svg
        viewBox="0 0 32 32"
        style="width: 32px; height: 32px;"
        class="border border-light rounded-circle"
      >
        <pattern
          id="checkerboard"
          width="12"
          height="12"
          patternUnits="userSpaceOnUse"
          fill="FFF"
        >
          <rect
            fill="#7080707f"
            x="0"
            width="6"
            height="6"
            y="0"
          />
          <rect
            fill="#7080707f"
            x="6"
            width="6"
            height="6"
            y="6"
          />
        </pattern>

        <circle
          cx="16"
          cy="16"
          r="16"
          fill="url(#checkerboard)"
        />

        <circle
          cx="16"
          cy="16"
          r="16"
        />
      </svg>
    </b-button>

    <b-modal
      :visible="showModal"
      :title="translations.modalTitle"
      :ok-title="translations.saveBtnLabel"
      centered
      size="sm"
      body-class="p-0"
      cancel-variant="link"
      @ok="$emit('input', currentColor || value)"
      @hide="closeMenu"
    >
      <chrome
        :value="value"
        class="w-100 shadow-none"
        @input="updateColor"
      />
    </b-modal>
  </div>
</template>

<script>
import { Chrome } from 'vue-color'
import { debounce } from 'lodash'

export default {
  name: 'CInputColorPicker',

  components: {
    Chrome,
  },

  props: {
    value: {
      type: String,
      default: 'rgba(0,0,0,0)',
    },

    translations: {
      type: Object,
    },
  },

  data () {
    return {
      showModal: false,
      currentColor: '',
    }
  },

  computed: {
    pickedColor: {
      get () {
        return this.value
      },

      set (pickedColor) {
        this.pickedColor = pickedColor
      },
    },
  },

  methods: {
    updateColor: debounce(function ({ hex8 = '' }) {
      this.currentColor = hex8
    }, 300),

    toggleMenu () {
      if (this.showModal) {
        this.closeMenu()
      } else {
        this.openMenu()
      }
    },

    openMenu () {
      this.showModal = true
    },

    closeMenu () {
      this.showModal = false
    },
  },
}
</script>

<style>
.vc-chrome {
  font-family: 'Poppins-Medium' !important;
}
</style>
