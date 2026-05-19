<template>
  <div class="flex flex-col gap-3">
    <div class="flex flex-col gap-1">
      <label class="text-primary font-medium text-sm">
        {{ $t('block.chatbotInbox.config.chatbots') }}
      </label>
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
      <small class="text-muted-color">
        {{ $t('block.chatbotInbox.config.chatbotsHint') }}
      </small>
    </div>

    <div class="flex flex-col gap-1">
      <label class="text-primary font-medium text-sm">
        {{ $t('block.chatbotInbox.config.status') }}
      </label>
      <div class="flex flex-wrap gap-3">
        <div
          v-for="s in statusChoices"
          :key="s.value"
          class="flex items-center gap-2"
        >
          <Checkbox
            v-model="statusFilter"
            :input-id="`cb-inbox-st-${s.value}`"
            :value="s.value"
          />
          <label :for="`cb-inbox-st-${s.value}`" class="text-sm">
            {{ s.label }}
          </label>
        </div>
      </div>
    </div>

    <div class="flex flex-col gap-1">
      <label class="text-primary font-medium text-sm">
        {{ $t('block.chatbotInbox.config.refresh') }}
      </label>
      <InputNumber
        v-model="refreshRate"
        :min="1"
        :step="1"
        show-buttons
        class="w-full"
      />
      <small class="text-muted-color">
        {{ $t('block.chatbotInbox.config.refreshHint') }}
      </small>
    </div>

    <div class="flex items-center gap-2">
      <Checkbox v-model="autoOpenFirst" :binary="true" input-id="cb-inbox-auto" />
      <label for="cb-inbox-auto" class="text-sm">
        {{ $t('block.chatbotInbox.config.autoOpen') }}
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
  { value: 'handoff_requested', label: t('block.chatbotInbox.status.handoff_requested') },
  { value: 'handoff_active', label: t('block.chatbotInbox.status.handoff_active') },
  { value: 'active', label: t('block.chatbotInbox.status.active') },
  { value: 'closed', label: t('block.chatbotInbox.status.closed') },
])

// PageBlockChatbotInbox in lib/js guarantees options is shaped on construction,
// so configurator can write straight through.
const opts = computed(() => block.value.options)

const chatbotIDs = computed({
  get: () => opts.value.chatbotIDs,
  set: v => { opts.value.chatbotIDs = v },
})

const statusFilter = computed({
  get: () => opts.value.statusFilter,
  set: v => { opts.value.statusFilter = v },
})

const refreshRate = computed({
  get: () => opts.value.refreshRate,
  set: v => { opts.value.refreshRate = Number(v) || opts.value.refreshRate },
})

const autoOpenFirst = computed({
  get: () => opts.value.autoOpenFirst,
  set: v => { opts.value.autoOpenFirst = !!v },
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
