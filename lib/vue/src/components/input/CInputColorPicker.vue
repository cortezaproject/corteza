<template>
  <div>
    <b-button
      :style="`color: ${pickedColor || 'rgba(0,0,0,0)'}; fill: ${pickedColor || 'rgba(0,0,0,0)'};`"
      class="p-0 rounded-circle btn"
      @click="toggleMenu"
    >
      <svg
        viewBox="0 0 24 24"
        style="width: 20px; height: 20px;"
      >
        <pattern
          id="checkerboard"
          width="6"
          height="6"
          patternUnits="userSpaceOnUse"
          fill="FFF"
        >
          <rect
            fill="#7080707f"
            x="0"
            width="3"
            height="3"
            y="0"
          />
          <rect
            fill="#7080707f"
            x="3"
            width="3"
            height="3"
            y="3"
          />
        </pattern>
        <circle
          cx="12"
          cy="12"
          r="12"
          fill="url(#checkerboard)"
        />
        <circle
          cx="12"
          cy="12"
          r="12"
        />
      </svg>
    </b-button>
    <b-modal
      :visible="showModal"
      title="Choose color"
      size="sm"
      centered
      body-class="p-3"
      cancel-variant="link"
      :ok-title="translations.saveBtnLabel"
      @ok="$emit('input', currentColor)"
      @hide="closeMenu"
    >
      <div
        class="d-flex justify-content-center"
      >
        <chrome
          :value="value || 'rgba(0,0,0,0)'"
          @input="updateColor"
        />
      </div>
      <small
        class="pt-2 px-4 form-text font-weight-light text-muted"
      >
        {{ translations.description }}
      </small>
    </b-modal>
  </div>
</template>

<script>
import { debounce } from 'lodash'
import { Chrome } from 'vue-color'

export default {
  name: 'CInputColorPicker',

  components: {
    Chrome,
  },

  props: {
    value: {
      type: String,
      default: '',
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
    updateColor: debounce(function (color = '') {
      this.currentColor = color.hex8
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
<style scoped>
.btn {
  background-color: white !important;
  border-color: white !important;
  box-shadow: none !important;
}

.vc-input__input {
  font-family: 'Poppins-Medium'
}
</style>
