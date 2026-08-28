<template>
  <span v-if="allowed || explainWhenDenied" data-c311-capability-action>
    <b-button
      v-if="allowed"
      v-bind="$attrs"
      :disabled="disabled || busy"
      @click="$emit('click', $event)"
    >
      <span v-if="busy" class="spinner-border spinner-border-sm mr-1" aria-hidden="true" />
      <slot />
    </b-button>
    <span v-else class="text-muted small" role="note">
      {{ deniedLabel }}
    </span>
  </span>
</template>

<script>
export default {
  name: 'C311CapabilityAction',
  inheritAttrs: false,
  props: {
    capability: {
      type: String,
      required: true,
    },
    scope: {
      type: String,
      default: '',
    },
    allowAnonymous: {
      type: Boolean,
      default: false,
    },
    explainWhenDenied: {
      type: Boolean,
      default: false,
    },
    deniedLabel: {
      type: String,
      default: 'This action is unavailable for your role.',
    },
    disabled: Boolean,
    busy: Boolean,
  },
  computed: {
    allowed () {
      if (this.allowAnonymous && !this.$C311?.session?.authenticated) return true
      return !!this.$C311?.can(this.capability) && (!this.scope || this.$C311.hasScope(this.scope))
    },
  },
}
</script>
