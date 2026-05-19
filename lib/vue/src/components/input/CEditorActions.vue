<template>
  <div class="shrink-0 border-t border-surface bg-surface">
    <div
      class="p-3 flex items-center"
      :class="hasLeft ? 'justify-between' : 'justify-end'"
    >
      <div v-if="hasLeft" class="flex items-center gap-2">
        <Button
          v-if="backTo"
          :label="backLabel || $t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="onBack"
        />
      </div>

      <div v-if="$slots.center" class="flex items-center gap-2">
        <slot name="center" />
      </div>

      <div class="flex gap-2">
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, useSlots } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  backTo: { type: [Object, String, Boolean, null], default: null },
  backLabel: { type: String, default: '' },
})

const emit = defineEmits(['back'])
const slots = useSlots()
const router = useRouter()

const hasLeft = computed(() => !!props.backTo)

function onBack() {
  emit('back')
  if (props.backTo && typeof props.backTo !== 'boolean') {
    router.push(props.backTo)
  }
}
</script>
