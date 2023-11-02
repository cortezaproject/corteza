<template>
  <div>
    <div class="d-flex align-items-center">
      <b-button
        :style="`color: ${pickedColor}; fill: ${pickedColor};`"
        class="p-0 rounded-circle bg-white border-white shadow-none"
        @click="toggleMenu"
      >
        <svg
          viewBox="0 0 32 32"
          :style="{ width: width, height: height }"
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
      <span v-if="showText" class="ml-2">
        {{ value }}
      </span>
    </div>
    <b-modal
      :visible="showModal"
      :title="translations.modalTitle"
      centered
      size="md"
      body-class="p-0"
      no-fade
      @hide="closeMenu"
    >
      <chrome
        :value="value"
        class="w-100 shadow-none"
        @input="updateColor"
      />

      <template #modal-footer>
        <slot name="footer" />

        <b-button
          variant="link"
          class="ml-auto"
          @click="closeMenu"
        >
          {{ translations.cancelBtnLabel }}
        </b-button>

        <b-button
          variant="primary"
          @click="saveColor"
        >
          {{ translations.saveBtnLabel }}
        </b-button>
      </template>
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

    width: {
      type: String,
      default: "32px",
    },

    height: {
      type: String,
      default: "32px",
    },

    showText: {
      type: Boolean,
      default: true,
    }
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

    saveColor () {
      this.$emit('input', this.currentColor || this.value)
      this.closeMenu()
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
