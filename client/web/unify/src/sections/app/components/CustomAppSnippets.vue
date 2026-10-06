<template>
  <Button
    :label="$t('app.snippets.insert')"
    icon="pi pi-angle-down"
    icon-pos="right"
    severity="secondary"
    variant="text"
    size="small"
    :disabled="disabled"
    data-test-id="custom-app-snippets"
    @click="menu.show($event, $event.currentTarget)"
  />
  <Menu ref="menu" :model="items" popup />
</template>

<script setup>
// The Insert menu over a custom app's or a Custom block's page: examples of
// the human bridge, written in with a module the page declared.
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { snippets } from '../hints'

const props = defineProps({
  // Module handles the page declared; the examples use the first.
  modules: { type: Array, default: () => [] },
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['insert'])

const { t } = useI18n()
const menu = ref(null)

const items = computed(() =>
  snippets(props.modules[0] || undefined).map(s => ({
    label: t(`app.snippets.${s.key}`),
    command: () => emit('insert', s.text),
  })),
)
</script>
