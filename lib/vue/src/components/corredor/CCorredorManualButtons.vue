<template>
  <b-button-group>
    <b-button
      v-for="(b, i) in buttons"
      :key="i"
      :variant="b.variant || defaultVariant"
      :class="buttonClass"
      @click="$emit('click', b)"
    >
      {{ b.label }}
    </b-button>
  </b-button-group>
</template>

<script lang="js">
export default {
  props: {
    resourceType: {},

    uiSlot: {
      type: String,
      required: true,
    },

    uiPage: {
      type: String,
      required: true,
    },

    defaultVariant: {
      type: String,
      default: () => 'link'
    },

    buttonClass: {
      type: String,
      default: () => 'mr-1'
    },
  },

  computed: {
    buttons () {
      return this.$UIHooks.Find(this.resourceType, this.uiPage, this.uiSlot)
    }
  }
}
</script>
