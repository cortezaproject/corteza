<template>
  <div class="flex flex-col gap-4">
    <Panel :header="$t('agent.editor.chatbot.panels.general')" toggleable>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <CInputToggleCard
          v-model="chatbot.enabled"
          :label="$t('agent.editor.chatbot.enabled.label')"
          :description="$t('agent.editor.chatbot.enabled.help')"
          class="md:col-span-2"
        />

        <div class="flex flex-col gap-1 md:col-span-2">
          <label class="font-medium text-primary">
            {{ $t('agent.editor.chatbot.widgetKey.label') }}
          </label>
          <small class="text-muted-color">
            {{ $t('agent.editor.chatbot.widgetKey.help') }}
          </small>
          <div class="flex gap-2 items-start">
            <InputText
              :model-value="chatbot.widgetKey || $t('agent.editor.chatbot.widgetKey.notGenerated')"
              readonly
              class="flex-1 font-mono"
            />
            <Button
              icon="pi pi-copy"
              severity="secondary"
              :disabled="!chatbot.widgetKey"
              size="small"
              @click="copyKey"
              v-tooltip.bottom="$t('agent.editor.chatbot.widgetKey.copy')"
            />
            <Button
              icon="pi pi-refresh"
              severity="warn"
              :disabled="!canRegenerate"
              size="small"
              @click="confirmRegenerate"
              v-tooltip.bottom="$t('agent.editor.chatbot.widgetKey.regenerate')"
            />
          </div>
        </div>

        <div class="flex flex-col gap-1 md:col-span-2">
          <label class="font-medium text-primary">
            {{ $t('agent.editor.chatbot.embed.label') }}
          </label>
          <Textarea :model-value="embedSnippet" readonly rows="2" class="font-mono text-xs" />
          <Button
            :label="$t('agent.editor.chatbot.embed.copy')"
            icon="pi pi-copy"
            severity="secondary"
            outlined
            size="small"
            class="self-start"
            :disabled="!chatbot.widgetKey"
            @click="copySnippet"
          />
        </div>

        <div class="flex flex-col gap-1">
          <label class="font-medium text-primary">
            {{ $t('agent.editor.chatbot.sessionTTL.label') }}
          </label>
          <small class="text-muted-color">
            {{ $t('agent.editor.chatbot.sessionTTL.help') }}
          </small>
          <InputText v-model="chatbot.sessionTTL" placeholder="2h" />
        </div>
      </div>
    </Panel>

    <Panel :header="$t('agent.editor.chatbot.panels.origins')" toggleable>
      <div class="flex flex-col gap-2">
        <small class="text-muted-color">
          {{ $t('agent.editor.chatbot.origins.help') }}
        </small>
        <div v-for="(_, idx) in chatbot.allowedOrigins" :key="idx" class="flex items-center gap-2">
          <InputText
            v-model="chatbot.allowedOrigins[idx]"
            placeholder="https://example.com"
            class="flex-1"
          />
          <Button
            icon="pi pi-trash"
            severity="danger"
            text
            size="small"
            @click="removeOrigin(idx)"
          />
        </div>
        <div>
          <Button
            :label="$t('agent.editor.chatbot.origins.add')"
            icon="pi pi-plus"
            severity="secondary"
            outlined
            size="small"
            @click="addOrigin"
          />
        </div>
      </div>
    </Panel>

    <Panel :header="$t('agent.editor.chatbot.panels.handoff')" toggleable>
      <Message severity="info" size="small" class="mb-3">
        {{ $t('agent.editor.chatbot.handoff.preview') }}
      </Message>
      <div class="flex flex-col gap-4">
        <CInputToggleCard
          v-model="chatbot.handoff.enabled"
          :label="$t('agent.editor.chatbot.handoff.enabled.label')"
          :description="$t('agent.editor.chatbot.handoff.enabled.help')"
        />
        <div
          class="flex flex-col gap-1"
          :class="{ 'opacity-50 pointer-events-none': !chatbot.handoff.enabled }"
        >
          <label class="font-medium text-primary">
            {{ $t('agent.editor.chatbot.handoff.targetRoles.label') }}
          </label>
          <CInputRole v-model="chatbot.handoff.targetRoles" :multiple="true" />
        </div>
      </div>
    </Panel>

    <ChatbotStyling :styling="chatbot.styling" />

    <ChatbotScenarios :scenarios="chatbot.scenarios" />
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import { useConfirm } from 'primevue/useconfirm'

import { components } from '@planetcrust/human-vue'

import ChatbotStyling from './ChatbotStyling.vue'
import ChatbotScenarios from './ChatbotScenarios.vue'

const { CInputRole, CInputToggleCard } = components

const props = defineProps({
  chatbot: { type: Object, required: true },
  agentID: { type: String, default: '' },
})

const emit = defineEmits(['regenerate-key'])

const { t } = useI18n()
const $toast = inject('$toast')
const confirm = useConfirm()

const canRegenerate = computed(() => !!props.agentID && !!props.chatbot.widgetKey)

const embedSnippet = computed(() => {
  if (!props.chatbot.widgetKey) return ''
  const origin = typeof window !== 'undefined' ? window.location.origin : ''
  return `<script src="${origin}/widget.js" data-widget-key="${props.chatbot.widgetKey}" async><\/script>`
})

function addOrigin() {
  props.chatbot.allowedOrigins.push('')
}

function removeOrigin(idx) {
  props.chatbot.allowedOrigins.splice(idx, 1)
}

async function copyKey() {
  await navigator.clipboard.writeText(props.chatbot.widgetKey)
  $toast?.toastSuccess(t('agent.editor.chatbot.widgetKey.copied'))
}

async function copySnippet() {
  await navigator.clipboard.writeText(embedSnippet.value)
  $toast?.toastSuccess(t('agent.editor.chatbot.embed.copied'))
}

function confirmRegenerate() {
  confirm.require({
    message: t('agent.editor.chatbot.widgetKey.regenerateConfirm'),
    header: t('agent.editor.chatbot.widgetKey.regenerate'),
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-warn',
    accept: () => emit('regenerate-key'),
  })
}
</script>
