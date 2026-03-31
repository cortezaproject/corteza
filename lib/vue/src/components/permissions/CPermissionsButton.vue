<template>
  <Button
    :label="label"
    :icon="icon"
    :severity="severity"
    :size="size"
    :outlined="outlined"
    :text="textStyle"
    :class="buttonClass"
    @click="handleClick"
  >
    <template v-if="$slots.default" #default>
      <slot />
    </template>
  </Button>
</template>

<script setup>
import { usePermissions } from '../../composables/usePermissions'

const props = defineProps({
  /** RBAC resource string, e.g. 'corteza::system:role/12345' */
  resource: {
    type: String,
    required: true,
  },
  /** Display name for the resource */
  title: {
    type: String,
    default: '',
  },
  /** Target name for i18n interpolation */
  target: {
    type: String,
    default: '',
  },
  /** If true, use 'all-specific' i18n pattern */
  allSpecific: {
    type: Boolean,
    default: false,
  },
  /** Button label text */
  label: {
    type: String,
    default: undefined,
  },
  /** PrimeVue icon class */
  icon: {
    type: String,
    default: 'pi pi-lock',
  },
  /** PrimeVue button severity */
  severity: {
    type: String,
    default: 'secondary',
  },
  /** PrimeVue button size */
  size: {
    type: String,
    default: 'small',
  },
  /** Whether the button is outlined */
  outlined: {
    type: Boolean,
    default: true,
  },
  /** Whether to use text-style button (for menu items) */
  textStyle: {
    type: Boolean,
    default: false,
  },
  /** Extra CSS classes */
  buttonClass: {
    type: [String, Object, Array],
    default: undefined,
  },
})

const { open } = usePermissions()

function handleClick() {
  open({
    resource: props.resource,
    title: props.title,
    target: props.target,
    allSpecific: props.allSpecific,
  })
}
</script>
