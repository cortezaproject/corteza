<template>
  <Panel :header="$t('chatbot.editor.panels.styling')" toggleable>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('chatbot.editor.styling.title.label') }}
        </label>
        <small class="text-muted-color">
          {{ $t('chatbot.editor.styling.title.help') }}
        </small>
        <InputText v-model="styling.launcher.label" />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('chatbot.editor.styling.logo.label') }}
        </label>
        <small class="text-muted-color">
          {{ $t('chatbot.editor.styling.logo.help') }}
        </small>
        <CFileDropZone
          accept="image/*"
          :uploading="logoUploading"
          :error="logoError"
          :preview-url="logoPreviewUrl"
          :clearable="!!styling.logoURL"
          :drop-label="$t('chatbot.editor.styling.logo.placeholder')"
          :label="$t('chatbot.editor.styling.logo.label')"
          compact
          preview-max-width="100%"
          preview-max-height="140px"
          @select="onLogoSelect"
          @clear="onLogoClear"
        />
      </div>
    </div>

    <Divider align="left">
      <span class="text-xs uppercase tracking-wide text-muted-color">
        {{ $t('chatbot.editor.styling.colors') }}
      </span>
    </Divider>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div v-for="slot in contentColorSlots" :key="slot" class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t(`chatbot.editor.styling.color.${slot}`) }}
        </label>
        <CInputColorPicker v-model="styling.colors[slot]" />
      </div>
    </div>

    <Divider align="left">
      <span class="text-xs uppercase tracking-wide text-muted-color">
        {{ $t('chatbot.editor.styling.typography') }}
      </span>
    </Divider>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('chatbot.editor.styling.fontFamily.label') }}
        </label>
        <small class="text-muted-color">
          {{ $t('chatbot.editor.styling.fontFamily.help') }}
        </small>
        <InputText v-model="styling.fontFamily" placeholder="Inter, sans-serif" />
      </div>
      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('chatbot.editor.styling.fontSize.heading') }}
        </label>
        <small class="text-muted-color">
          {{ $t('chatbot.editor.styling.fontSizeHelp.heading') }}
        </small>
        <PxInput v-model="styling.fontSizes.heading" :min="8" :max="48" />
      </div>
      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('chatbot.editor.styling.fontSize.base') }}
        </label>
        <small class="text-muted-color">
          {{ $t('chatbot.editor.styling.fontSizeHelp.base') }}
        </small>
        <PxInput v-model="styling.fontSizes.base" :min="8" :max="48" />
      </div>
    </div>

    <Divider align="left">
      <span class="text-xs uppercase tracking-wide text-muted-color">
        {{ $t('chatbot.editor.styling.buttonTitle') }}
      </span>
    </Divider>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('chatbot.editor.styling.button.label.label') }}
        </label>
        <small class="text-muted-color">
          {{ $t('chatbot.editor.styling.button.label.help') }}
        </small>
        <InputText v-model="styling.launcher.buttonLabel" />
      </div>

      <div class="flex flex-col gap-2">
        <CInputToggleCard
          v-model="styling.launcher.iconVisible"
          :label="$t('chatbot.editor.styling.button.icon.label')"
          :description="$t('chatbot.editor.styling.button.icon.help')"
        />
        <CFileDropZone
          v-if="styling.launcher.iconVisible"
          accept="image/*"
          :uploading="iconUploading"
          :error="iconError"
          :clearable="!!styling.launcher.iconURL"
          :drop-label="$t('chatbot.editor.styling.button.icon.placeholder')"
          :label="$t('chatbot.editor.styling.button.icon.label')"
          compact
          @select="onIconSelect"
          @clear="onIconClear"
        />
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('chatbot.editor.styling.button.color.label') }}
        </label>
        <CInputColorPicker v-model="styling.colors.primary" />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('chatbot.editor.styling.button.colorText.label') }}
        </label>
        <CInputColorPicker v-model="styling.colors.primaryText" />
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('chatbot.editor.styling.button.size.label') }}
        </label>
        <small class="text-muted-color">
          {{ $t('chatbot.editor.styling.button.size.help') }}
        </small>
        <PxInput v-model="styling.launcher.size" :min="32" :max="96" />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('chatbot.editor.styling.button.shape.label') }}
        </label>
        <small class="text-muted-color">
          {{ $t('chatbot.editor.styling.button.shape.help') }}
        </small>
        <Select
          v-model="styling.launcher.shape"
          :options="[
            { label: $t('chatbot.editor.styling.button.shape.circle'), value: 'circle' },
            { label: $t('chatbot.editor.styling.button.shape.square'), value: 'square' },
          ]"
          optionLabel="label"
          optionValue="value"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('chatbot.editor.styling.button.position.label') }}
        </label>
        <small class="text-muted-color">
          {{ $t('chatbot.editor.styling.button.position.help') }}
        </small>
        <Select
          v-model="styling.launcher.position"
          :options="[
            { label: $t('chatbot.editor.styling.button.position.bottomRight'), value: 'bottom-right' },
            { label: $t('chatbot.editor.styling.button.position.bottomCenter'), value: 'bottom-center' },
            { label: $t('chatbot.editor.styling.button.position.bottomLeft'), value: 'bottom-left' },
            { label: $t('chatbot.editor.styling.button.position.leftMiddle'), value: 'left-middle' },
            { label: $t('chatbot.editor.styling.button.position.rightMiddle'), value: 'right-middle' },
            { label: $t('chatbot.editor.styling.button.position.topRight'), value: 'top-right' },
            { label: $t('chatbot.editor.styling.button.position.topCenter'), value: 'top-center' },
            { label: $t('chatbot.editor.styling.button.position.topLeft'), value: 'top-left' },
          ]"
          optionLabel="label"
          optionValue="value"
        />
      </div>
    </div>
  </Panel>
</template>

<script setup>
import { computed, inject } from 'vue'
import { components, useFileUpload } from '@planetcrust/human-vue'

import PxInput from './PxInput.vue'

const { CInputColorPicker, CFileDropZone, CInputToggleCard } = components

const props = defineProps({
  styling: { type: Object, required: true },
})

const $SystemAPI = inject('$SystemAPI')

const contentColorSlots = ['header', 'headerText', 'background', 'text', 'userBubble', 'agentBubble']

const {
  uploading: logoUploading,
  uploadError: logoUploadError,
  uploadFiles: uploadLogoFiles,
  reset: resetLogoUpload,
} = useFileUpload()

const logoError = computed(() => logoUploadError.value)
const logoPreviewUrl = computed(() => props.styling.logoURL)

async function onLogoSelect(files) {
  const file = files[0]
  if (!file) return
  try {
    const results = await uploadLogoFiles([file], {
      api: $SystemAPI,
      endpoint: $SystemAPI.applicationUploadEndpoint(),
    })
    const rsp = results[0]
    if (rsp) props.styling.logoURL = $SystemAPI.baseURL + rsp.url
  } catch {
    /* error set by composable */
  }
}

function onLogoClear() {
  props.styling.logoURL = ''
  resetLogoUpload()
}

const {
  uploading: iconUploading,
  uploadError: iconUploadError,
  uploadFiles: uploadIconFiles,
  reset: resetIconUpload,
} = useFileUpload()

const iconError = computed(() => iconUploadError.value)

async function onIconSelect(files) {
  const file = files[0]
  if (!file) return
  try {
    const results = await uploadIconFiles([file], {
      api: $SystemAPI,
      endpoint: $SystemAPI.applicationUploadEndpoint(),
    })
    const rsp = results[0]
    if (rsp) props.styling.launcher.iconURL = $SystemAPI.baseURL + rsp.url
  } catch {
    /* error set by composable */
  }
}

function onIconClear() {
  props.styling.launcher.iconURL = ''
  resetIconUpload()
}
</script>
