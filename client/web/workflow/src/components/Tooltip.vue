<template>
  <b-tooltip
    :id="kind"
    :target="kind"
    triggers="hover"
    placement="right"
    variant="light"
    custom-class="ml-2"
    boundary="window"
    noninteractive
    @hide="show = false"
    @show="show = true"
  >
    <transition name="fade">
      <b-card
        v-if="show"
        :img-src="img"
        :img-alt="kind"
        img-height="100px"
        img-top
        :title="title"
        class="text-left"
      >
        <b-card-text>
          {{ text }}
        </b-card-text>
      </b-card>
    </transition>
  </b-tooltip>
</template>

<script>
export default {
  props: {
    title: {
      type: String,
      required: true,
    },
    kind: {
      type: String,
      required: true,
    },
    img: {
      type: String,
      required: true,
    },
    text: {
      type: String,
      required: true,
    },
  },

  data () {
    return {
      show: false,
    }
  },
}
</script>

<style lang="scss">
.tooltip.b-tooltip-light > .tooltip-inner {
  padding: 0;
  max-width: 20rem;
  min-width: 20rem;
  background-color: transparent !important;
  opacity: 1;
}

.tooltip.b-tooltip-light.bs-tooltip-right .arrow::before {
  border-right-color: #2D2D2D !important;
}

.tooltip.b-tooltip {
  opacity: 1;
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */ {
  opacity: 0;
}
</style>
