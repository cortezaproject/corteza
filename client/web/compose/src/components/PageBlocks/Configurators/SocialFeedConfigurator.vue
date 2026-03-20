<template>
  <div class="flex flex-col gap-3">
    <div class="flex flex-col gap-1">
      <label class="text-primary font-medium text-sm">{{ $t('block.socialFeed.twitterProfileLabel') }}</label>
      <InputText
        v-model="profileUrl"
        placeholder="https://twitter.com/..."
        class="w-full"
      />
    </div>

    <div v-if="isRecordPage" class="flex flex-col gap-1">
      <label class="text-primary font-medium text-sm">{{ $t('block.socialFeed.twitterProfileField') }}</label>
      <Select
        v-model="profileSourceField"
        :options="urlFields"
        option-label="label"
        option-value="name"
        class="w-full"
        show-clear
      />
    </div>

    <Divider />

    <div class="flex flex-col gap-1">
      <label class="text-primary font-medium text-sm">{{ $t('block.general.refreshRate') }}</label>
      <InputNumber
        v-model="refreshRate"
        :min="0"
        suffix=" s"
        class="w-full"
      />
    </div>

    <div class="flex items-center gap-2">
      <Checkbox v-model="showRefresh" binary input-id="showRefreshSocialFeed" />
      <label for="showRefreshSocialFeed" class="text-sm">{{ $t('block.general.showRefresh') }}</label>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useModuleStore } from '@/stores/module'

const moduleStore = useModuleStore()

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const emit = defineEmits(['update:block'])

const isRecordPage = computed(() => !!props.page?.moduleID && props.page.moduleID !== '0')

const urlFields = computed(() => {
  if (!props.page?.moduleID) return []
  const mod = moduleStore.getByID(props.page.moduleID)
  if (!mod) return []
  return mod.fields
    .filter(f => f.kind === 'Url' || f.kind === 'String')
    .map(f => ({ name: f.name, label: f.label || f.name }))
})

function updateOptions(key, value) {
  emit('update:block', {
    ...props.block,
    options: { ...props.block.options, [key]: value },
  })
}

const profileUrl = computed({
  get: () => props.block.options?.profileUrl || '',
  set: v => updateOptions('profileUrl', v),
})

const profileSourceField = computed({
  get: () => props.block.options?.profileSourceField || '',
  set: v => updateOptions('profileSourceField', v),
})

const refreshRate = computed({
  get: () => props.block.options?.refreshRate ?? 0,
  set: v => updateOptions('refreshRate', v),
})

const showRefresh = computed({
  get: () => !!props.block.options?.showRefresh,
  set: v => updateOptions('showRefresh', v),
})
</script>
