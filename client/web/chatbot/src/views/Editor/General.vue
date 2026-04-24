<template>
  <Panel :header="$t('chatbot.editor.panels.general')" toggleable>
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <div class="flex flex-col gap-1">
        <label for="name" class="font-medium text-primary">
          {{ $t('chatbot.editor.name.label') }}
        </label>
        <small class="text-muted-color">{{ $t('chatbot.editor.name.help') }}</small>
        <InputText id="name" v-model="chatbot.name" />
      </div>

      <div class="flex flex-col gap-1">
        <label for="handle" class="font-medium text-primary">
          {{ $t('chatbot.editor.handle.label') }}
        </label>
        <small class="text-muted-color">{{ $t('chatbot.editor.handle.help') }}</small>
        <InputText id="handle" v-model="chatbot.handle" />
      </div>

      <CInputToggleCard
        v-model="chatbot.enabled"
        :label="$t('chatbot.editor.enabled.label')"
        :description="$t('chatbot.editor.enabled.help')"
        class="self-start"
      />

      <div class="hidden lg:block" />

      <Divider class="lg:col-span-2 !my-0" />

      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('chatbot.editor.embed.label') }}
        </label>
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
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('chatbot.editor.widgetKey.label') }}
        </label>
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
      </div>

      <div class="flex flex-col gap-2">
        <label class="font-medium text-primary">
          {{ $t('chatbot.editor.panels.origins') }}
        </label>
        <small class="text-muted-color">{{ $t('chatbot.editor.origins.help') }}</small>
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
            :label="$t('chatbot.editor.origins.add')"
            icon="pi pi-plus"
            severity="secondary"
            outlined
            size="small"
            @click="addOrigin"
          />
        </div>
      </div>

      <div class="flex flex-col gap-1 self-start">
        <label class="font-medium text-primary">
          {{ $t('chatbot.editor.sessionTTL.label') }}
        </label>
        <small class="text-muted-color">{{ $t('chatbot.editor.sessionTTL.help') }}</small>
        <InputText v-model="chatbot.sessionTTL" placeholder="2h" />
      </div>
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

function removeOrigin(idx) {
  props.chatbot.allowedOrigins.splice(idx, 1)
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
