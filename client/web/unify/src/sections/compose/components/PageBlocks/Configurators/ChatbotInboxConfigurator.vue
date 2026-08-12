<template>
  <div class="flex flex-col gap-3">
    <CFormGroup
      :label="$t('block.chatbotInbox.config.chatbots')"
      :description="$t('block.chatbotInbox.config.chatbotsHint')"
    >
      <MultiSelect
        v-model="chatbotIDs"
        :options="chatbots"
        option-label="name"
        option-value="chatbotID"
        :placeholder="$t('block.chatbotInbox.config.chatbotsPlaceholder')"
        :loading="loading"
        filter
        display="chip"
        class="w-full"
      />
    </CFormGroup>

    <CFormGroup :label="$t('block.chatbotInbox.config.status')">
      <div class="flex flex-wrap gap-3">
        <div v-for="s in statusChoices" :key="s.value" class="flex items-center gap-2">
          <Checkbox v-model="statusFilter" :input-id="`cb-inbox-st-${s.value}`" :value="s.value" />
          <label :for="`cb-inbox-st-${s.value}`" class="text-sm">
            {{ s.label }}
          </label>
        </div>
      </div>
    </CFormGroup>

    <CFormGroup
      :label="$t('block.chatbotInbox.config.refresh')"
      :description="$t('block.chatbotInbox.config.refreshHint')"
    >
      <InputNumber v-model="refreshRate" :min="1" :step="1" show-buttons class="w-full" />
    </CFormGroup>

    <div class="flex items-center gap-2">
      <Checkbox v-model="autoOpenFirst" :binary="true" input-id="cb-inbox-auto" />
      <label for="cb-inbox-auto" class="text-sm">
        {{ $t('block.chatbotInbox.config.autoOpen') }}
      </label>
    </div>

    <div class="flex items-center gap-2">
      <Checkbox v-model="showFilter" :binary="true" input-id="cb-inbox-show-filter" />
      <label for="cb-inbox-show-filter" class="text-sm">
        {{ $t('block.chatbotInbox.config.showFilter') }}
      </label>
    </div>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const block = inject('blockDraft')
const $SystemAPI = inject('$SystemAPI', null)
const { t } = useI18n()

const loading = ref(false)
const chatbots = ref([])

const statusChoices = computed(() => [
  { value: 'handoff_requested', label: t('block.chatbotInbox.statusHandoffRequested') },
  { value: 'handoff_active', label: t('block.chatbotInbox.statusHandoffActive') },
  { value: 'active', label: t('block.chatbotInbox.statusActive') },
  { value: 'closed', label: t('block.chatbotInbox.statusClosed') },
])

// PageBlockChatbotInbox in lib/js guarantees options is shaped on construction,
// so configurator can write straight through.
const opts = computed(() => block.value.options)

const chatbotIDs = computed({
  get: () => opts.value.chatbotIDs,
  set: v => {
    opts.value.chatbotIDs = v
  },
})

const statusFilter = computed({
  get: () => opts.value.statusFilter,
  set: v => {
    opts.value.statusFilter = v
  },
})

const refreshRate = computed({
  get: () => opts.value.refreshRate,
  set: v => {
    opts.value.refreshRate = Number(v) || opts.value.refreshRate
  },
})

const autoOpenFirst = computed({
  get: () => opts.value.autoOpenFirst,
  set: v => {
    opts.value.autoOpenFirst = !!v
  },
})

const showFilter = computed({
  get: () => opts.value.showFilter,
  set: v => {
    opts.value.showFilter = !!v
  },
})

onMounted(async () => {
  if (!$SystemAPI) return
  loading.value = true
  try {
    const res = await $SystemAPI.chatbotList({ limit: 0, sort: 'name ASC' })
    chatbots.value = (res?.set || []).map(cb => ({
      ...cb,
      name: cb.name || cb.handle || cb.chatbotID,
    }))
  } catch (err) {
    console.warn('[chatbot-inbox] chatbotList failed', err)
  } finally {
    loading.value = false
  }
})
</script>
