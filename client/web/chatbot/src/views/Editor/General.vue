<template>
  <Panel :header="$t('chatbot.editor.panels.general')" toggleable>
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <CFormGroup
        :label="$t('chatbot.editor.name.label')"
        :description="$t('chatbot.editor.name.help')"
        input-id="name"
      >
        <InputText id="name" v-model="chatbot.name" />
      </CFormGroup>

      <CFormGroup
        :label="$t('chatbot.editor.handle.label')"
        :description="$t('chatbot.editor.handle.help')"
        input-id="handle"
      >
        <InputText id="handle" v-model="chatbot.handle" />
      </CFormGroup>

      <CInputToggleCard
        v-model="chatbot.enabled"
        :label="$t('chatbot.editor.enabled.label')"
        :description="$t('chatbot.editor.enabled.help')"
        class="self-start"
      />

      <div class="hidden lg:block" />

      <Divider class="lg:col-span-2 !my-0" />

      <CFormGroup :label="$t('chatbot.editor.embed.label')">
        <div class="flex gap-2 items-start">
          <Textarea
            :model-value="embedSnippet"
            readonly
            rows="2"
            class="flex-1 font-mono text-xs"
          />
          <Button
            v-tooltip.bottom="$t('chatbot.editor.embed.copy')"
            icon="pi pi-copy"
            severity="secondary"
            :disabled="!chatbot.widgetKey"
            @click="copySnippet"
          />
        </div>
      </CFormGroup>

      <CFormGroup :label="$t('chatbot.editor.widgetKey.label')">
        <div class="flex gap-2 items-start">
          <InputGroup class="flex-1">
            <InputText
              :model-value="chatbot.widgetKey || $t('chatbot.editor.widgetKey.notGenerated')"
              readonly
              class="font-mono"
            />
            <Button
              v-tooltip.bottom="$t('chatbot.editor.widgetKey.regenerate')"
              icon="pi pi-refresh"
              severity="warn"
              :disabled="!canRegenerate"
              @click="confirmRegenerate"
            />
          </InputGroup>
          <Button
            v-tooltip.bottom="$t('chatbot.editor.widgetKey.copy')"
            icon="pi pi-copy"
            severity="secondary"
            :disabled="!chatbot.widgetKey"
            @click="copyKey"
          />
        </div>
      </CFormGroup>

      <CFormGroup
        :label="$t('chatbot.editor.panels.origins')"
        :description="$t('chatbot.editor.origins.help')"
      >
        <template #actions>
          <Button
            :label="$t('chatbot.editor.origins.add')"
            icon="pi pi-plus"
            severity="secondary"
            size="small"
            @click="addOrigin"
          />
        </template>
        <CFormList
          v-model="chatbot.allowedOrigins"
          :columns="[{ label: '', width: '1fr' }]"
        >
          <template #row="{ index }">
            <InputText
              v-model="chatbot.allowedOrigins[index]"
              placeholder="https://example.com"
              size="small"
              class="w-full"
            />
          </template>
        </CFormList>
      </CFormGroup>

      <CFormGroup
        :label="$t('chatbot.editor.sessionTTL.label')"
        :description="$t('chatbot.editor.sessionTTL.help')"
        class="self-start"
      >
        <InputText v-model="chatbot.sessionTTL" placeholder="2h" />
      </CFormGroup>
    </div>
  </Panel>
</template>

<script setup>
import { components } from '@planetcrust/human-vue'
import { computed, inject } from 'vue'
import { useConfirm } from 'primevue/useconfirm'
import { useI18n } from 'vue-i18n'

const { CInputToggleCard } = components

const props = defineProps({
  chatbot: { type: Object, required: true },
  isCreate: { type: Boolean, default: false },
})

const emit = defineEmits(['regenerate-key'])

const { t } = useI18n()
const $toast = inject('$toast')
const confirm = useConfirm()

const canRegenerate = computed(() => !props.isCreate && !!props.chatbot.widgetKey)

const embedSnippet = computed(() => {
  if (!props.chatbot.widgetKey) return ''
  const origin = typeof window !== 'undefined' ? window.location.origin : ''
  return `<script src="${origin}/chatbot/widget.js" data-widget-key="${props.chatbot.widgetKey}" async><\/script>`
})

function addOrigin() {
  props.chatbot.allowedOrigins.push('')
}

async function copyKey() {
  await navigator.clipboard.writeText(props.chatbot.widgetKey)
  $toast?.toastSuccess(t('chatbot.editor.widgetKey.copied'))
}

async function copySnippet() {
  await navigator.clipboard.writeText(embedSnippet.value)
  $toast?.toastSuccess(t('chatbot.editor.embed.copied'))
}

function confirmRegenerate() {
  confirm.require({
    message: t('chatbot.editor.widgetKey.regenerateConfirm'),
    header: t('chatbot.editor.widgetKey.regenerate'),
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-warn',
    accept: () => emit('regenerate-key'),
  })
}
</script>
