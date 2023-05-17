<template>
  <div>
    <div class="color-picker">
      <b-button
        ref="guide"
        :style="`color: ${currentColor || 'rgba(0,0,0,0)'}; fill: ${currentColor || 'rgba(0,0,0,0)'};`"
        class="p-2 rounded-circle btn-white"
        @click="toggleMenu"
      >
        <svg
          viewBox="0 0 24 24"
          class="icon" 
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
      <div
        v-if="isMenuActive"
        class="color-menu"
      >
        <b-button
          class="button-close"
          @click="closeMenu"
        >
          <font-awesome-icon
            :icon="['fas', 'times']"
            class="icon-small"
          />
        </b-button>
        <chrome
          ref="menu"
          :value="value || 'rgba(0,0,0,0)'"
          class=""
          @input="updateColor"
        />
      </div>
    </div>
    <verte
      :value="value || 'rgba(0,0,0,0)'"
      @input="updateColor"
    />
  </div>
</template>

<script>
import { debounce } from 'lodash'
import { Chrome } from 'vue-color'
import Verte from 'verte'

export default {
  name: 'CInputColorPicker',

  components: {
    Chrome,
    Verte,
  },

  props: {
    value: {
      type: String,
      default: '',
    },
  },

  data () {
    return {
      isMenuActive: false,
      // closeCallback: '',c
    }
  },

  computed: {
    currentColor: {
      get () {
        return this.value
      },
      set (val) {
        this.currentColor = val
        this.updateColor(val)
      },
    },
  },

  methods: {
    updateColor: debounce(function (color = '') {
      this.$emit('input', color.hex8)
    }, 300),

    toggleMenu () {
      if (this.isMenuActive) {
        this.closeMenu()
      } else {
        this.openMenu()
      }
    },

    openMenu () {
      this.isMenuActive = true
      // this.closeCallback = (evnt) => {
      //   if (
      //     console.log(this.$refs)
      //     // console.log('!this.isElementClosest(evnt.target, this.$refs.menu)', !this.isElementClosest(evnt.target, this.$refs.menu))
      //     // console.log('!this.isElementClosest(evnt.target, this.$refs.guide)', !this.isElementClosest(evnt.target, this.$refs.guide) && !this.isElementClosest(evnt.target, this.$refs.menu))
      //     // !this.isElementClosest(evnt.target, this.$refs.menu) &&
      //     // !this.isElementClosest(evnt.target, this.$refs.guide)
      //   ) {
      //     this.closeMenu()
      //   }
      // }
      // document.addEventListener('mousedown', this.closeCallback)
    },

    closeMenu () {
      this.isMenuActive = false
      document.removeEventListener('mousedown', this.closeCallback);
      this.$emit('close', this.currentColor);
    },

    isElementClosest (element, wrapper) {
      while (element !== document && element !== null) {
        if (element === wrapper) return true
        element = element.parentNode
      }

      return false
    },
  },
}
</script>
<style>

.verte {
  position: relative;
  display: flex;
  justify-content: center;
}
.verte * {
    box-sizing: border-box;
}
.verte--loading {
  opacity: 0;
}
.verte__guide {
  width: 24px;
  height: 24px;
  padding: 0;
  border: 0;
  background: transparent;
}
.verte__guide:focus {
    outline: 0;
}
.verte__guide svg {
    width: 100%;
    height: 100%;
    fill: inherit;
}
.verte__menu {
  flex-direction: column;
  justify-content: center;
  align-items: stretch;
  width: 250px;
  border-radius: 6px;
  background-color: #fff;
  will-change: transform;
  box-shadow: 0 8px 15px rgba(0, 0, 0, 0.1);
}
.verte__menu:focus {
    outline: none;
}
.verte__menu-origin {
  display: none;
  position: absolute;
  z-index: 10;
}
.verte__menu-origin--active {
    display: flex;
}
.verte__menu-origin--static {
    position: static;
    z-index: initial;
}
.verte__menu-origin--top {
    bottom: 50px;
}
.verte__menu-origin--bottom {
    top: 50px;
}
.verte__menu-origin--right {
    right: 0;
}
.verte__menu-origin--left {
    left: 0;
}
.verte__menu-origin--center {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    justify-content: center;
    align-items: center;
    background-color: rgba(0, 0, 0, 0.1);
}
.verte__menu-origin:focus {
    outline: none;
}
.verte__controller {
  padding: 0 20px 20px;
}
.verte__recent {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  align-items: center;
  width: 100%;
}
.verte__recent-color {
    margin: 4px;
    width: 27px;
    height: 27px;
    border-radius: 50%;
    background-color: #fff;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    background-image: linear-gradient(45deg, rgba(112, 128, 144, 0.5) 25%, transparent 25%), linear-gradient(45deg, transparent 75%, rgba(112, 128, 144, 0.5) 75%), linear-gradient(-45deg, rgba(112, 128, 144, 0.5) 25%, transparent 25%), linear-gradient(-45deg, transparent 75%, rgba(112, 128, 144, 0.5) 75%);
    background-size: 6px 6px;
    background-position: 0 0, 3px -3px, 0 3px, -3px 0px;
    overflow: hidden;
}
.verte__recent-color:after {
      content: '';
      display: block;
      width: 100%;
      height: 100%;
      background-color: currentColor;
}
.verte__value {
  padding: 0.6em;
  width: 100%;
  border: 1px solid #708090;
  border-radius: 6px 0 0 6px;
  text-align: center;
  font-size: 12px;
  -webkit-appearance: none;
  -moz-appearance: textfield;
}
.verte__value:focus {
    outline: none;
    border-color: #1a3aff;
}
.verte__icon {
  width: 20px;
  height: 20px;
}
.verte__icon--small {
    width: 12px;
    height: 12px;
}
.verte__input {
  padding: 5px;
  margin: 0 3px;
  min-width: 0;
  text-align: center;
  border-width: 0 0 1px 0;
  appearance: none;
  -moz-appearance: textfield;
}
.verte__input::-webkit-inner-spin-button, .verte__input::-webkit-outer-spin-button {
    -webkit-appearance: none;
    margin: 0;
}
.verte__inputs {
  display: flex;
  font-size: 16px;
  margin-bottom: 5px;
}
.verte__draggable {
  border-radius: 6px 6px 0 0;
  height: 8px;
  width: 100%;
  cursor: grab;
  background: linear-gradient(90deg, #fff 2px, transparent 1%) center, linear-gradient(#fff 2px, transparent 1%) center, rgba(112, 128, 144, 0.2);
  background-size: 4px 4px;
}
.verte__model,
.verte__submit {
  position: relative;
  display: inline-flex;
  justify-content: center;
  align-items: center;
  padding: 1px;
  border: 0;
  text-align: center;
  cursor: pointer;
  background-color: transparent;
  font-weight: 700;
  color: #708090;
  fill: #708090;
  outline: none;
}
.verte__model:hover,
  .verte__submit:hover {
    fill: #1a3aff;
    color: #1a3aff;
}
.verte__close {
  position: absolute;
  top: 1px;
  right: 1px;
  z-index: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 4px;
  cursor: pointer;
  border-radius: 50%;
  border: 0;
  transform: translate(50%, -50%);
  background-color: rgba(0, 0, 0, 0.4);
  fill: #fff;
  outline: none;
  box-shadow: 1px 1px 1px rgba(0, 0, 0, 0.2);
}
.verte__close:hover {
    background-color: rgba(0, 0, 0, 0.6);
}

/*# sourceMappingURL=Verte.vue.map */
.verte-picker {
  width: 100%;
  margin: 0 auto 10px;
  display: flex;
  flex-direction: column;
}
.verte-picker--wheel {
    margin-top: 20px;
}
.verte-picker__origin {
    user-select: none;
    position: relative;
    margin: 0 auto;
    overflow: hidden;
}
.verte-picker__slider {
    margin: 20px 20px 0;
}
.verte-picker__canvas {
    display: block;
}
.verte-picker__cursor {
    position: absolute;
    top: 0;
    left: 0;
    margin: -6px;
    width: 12px;
    height: 12px;
    border: 1px solid #fff;
    border-radius: 50%;
    will-change: transform;
    pointer-events: none;
    background-color: transparent;
    box-shadow: #fff 0px 0px 0px 1.5px, rgba(0, 0, 0, 0.3) 0px 0px 1px 1px inset, rgba(0, 0, 0, 0.4) 0px 0px 1px 2px;
}
.verte-picker__input {
    display: flex;
    margin-bottom: 10px;
}

/*# sourceMappingURL=Picker.vue.map */
.slider {
  position: relative;
  display: flex;
  align-items: center;
  box-sizing: border-box;
  margin-bottom: 10px;
  font-size: 20px;
}
.slider:hover .slider-label, .slider--dragging .slider-label {
    visibility: visible;
    opacity: 1;
}
.slider__input {
  margin-bottom: 0;
  padding: 0.3em;
  margin-left: 0.2em;
  max-width: 70px;
  width: 20%;
  border: 0;
  text-align: center;
  font-size: 12px;
  -webkit-appearance: none;
  -moz-appearance: textfield;
}
.slider__input::-webkit-inner-spin-button, .slider__input::-webkit-outer-spin-button {
    -webkit-appearance: none;
    margin: 0;
}
.slider__input:focus {
    outline: none;
    border-color: #1a3aff;
}
.slider__track {
  position: relative;
  flex: 1;
  margin: 3px;
  width: auto;
  height: 8px;
  background: #fff;
  will-change: transfom;
  background-image: linear-gradient(45deg, rgba(112, 128, 144, 0.5) 25%, transparent 25%), linear-gradient(45deg, transparent 75%, rgba(112, 128, 144, 0.5) 75%), linear-gradient(-45deg, rgba(112, 128, 144, 0.5) 25%, transparent 25%), linear-gradient(-45deg, transparent 75%, rgba(112, 128, 144, 0.5) 75%);
  background-size: 6px 6px;
  background-position: 0 0, 3px -3px, 0 3px, -3px 0px;
  border-radius: 10px;
}
.slider__handle {
  position: relative;
  position: absolute;
  top: 0;
  left: 0;
  will-change: transform;
  color: #000;
  margin: -2px 0 0 -8px;
  width: 12px;
  height: 12px;
  border: 2px solid #fff;
  background-color: currentColor;
  border-radius: 50%;
  box-shadow: 0 1px 4px -2px black;
}
.slider__label {
  position: absolute;
  top: -3em;
  left: 0.4em;
  z-index: 999;
  visibility: hidden;
  padding: 6px;
  min-width: 3em;
  border-radius: 6px;
  background-color: #000;
  color: #fff;
  text-align: center;
  font-size: 12px;
  line-height: 1em;
  opacity: 0;
  transform: translate(-50%, 0);
  white-space: nowrap;
}
.slider__label:before {
    position: absolute;
    bottom: -0.6em;
    left: 50%;
    display: block;
    width: 0;
    height: 0;
    border-width: 0.6em 0.6em 0 0.6em;
    border-style: solid;
    border-color: #000 transparent transparent transparent;
    content: '';
    transform: translate3d(-50%, 0, 0);
}
.slider__fill {
  width: 100%;
  height: 100%;
  transform-origin: left top;
  border-radius: 10px;
}

/*# sourceMappingURL=Slider.vue.map */
</style>
<style lang="scss">
.verte__menu-origin {
  display: none;
  position: absolute;
  z-index: 10;
}
.verte__menu-origin--active {
    display: flex;
}
.verte__menu-origin--center {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    justify-content: center;
    align-items: center;
    background-color: rgba(0, 0, 0, 0.1);
}
.verte__menu-origin:focus {
    outline: none;
}
.btn-white:hover,
.btn-white:focus,
.btn-white:active {
  background-color: white !important;
  border-color: white !important;
  box-shadow: none !important;
}

.color-picker * {
  box-sizing: border-box;
}

.color-picker {
  position: relative
}

.icon {
  width: 20px;
  height: 20px;
}

.icon-small {
  width: 12px !important;
  height: 12px;
}

.color-menu {
  position: absolute;
  z-index: 10;
}

.color-menu-active {
  display: block;
}

.button-close {
  position: absolute;
  top: 1px;
  right: 1px;
  z-index: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 4px;
  cursor: pointer;
  border-radius: 50%;
  border: 0;
  transform: translate(50%, -50%);
  background-color: rgba(0, 0, 0, 0.4);
  fill: #fff;
  outline: none;
  box-shadow: 1px 1px 1px rgba(0, 0, 0, 0.2);
}

.button-close:hover {
  background-color: rgba(0, 0, 0, 0.6);
}

.vc-chrome-fields .vc-input__input {
  font-family: "Poppins-Medium";
}
</style>
