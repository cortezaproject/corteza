<template>
  <Panel :header="$t('agent.editor.chatbot.panels.styling')" toggleable>
    <!-- Branding: logo, title, theme colors, typography -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('agent.editor.chatbot.styling.logo.label') }}
        </label>
        <small class="text-muted-color">
          {{ $t('agent.editor.chatbot.styling.logo.help') }}
        </small>
        <CFileDropZone
          accept="image/*"
          :uploading="logoUploading"
          :error="logoError"
          :preview-url="logoPreviewUrl"
          :clearable="!!styling.logoURL"
          :drop-label="$t('agent.editor.chatbot.styling.logo.placeholder')"
          :label="$t('agent.editor.chatbot.styling.logo.label')"
          compact
          preview-max-width="100%"
          preview-max-height="140px"
          @select="onLogoSelect"
          @clear="onLogoClear"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('agent.editor.chatbot.styling.title.label') }}
        </label>
        <small class="text-muted-color">
          {{ $t('agent.editor.chatbot.styling.title.help') }}
        </small>
        <InputText v-model="styling.launcher.label" />
      </div>
    </div>

    <Divider align="left">
      <span class="text-xs uppercase tracking-wide text-muted-color">
        {{ $t('agent.editor.chatbot.styling.colors') }}
      </span>
    </Divider>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div v-for="slot in contentColorSlots" :key="slot" class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t(`agent.editor.chatbot.styling.color.${slot}`) }}
        </label>
        <small class="text-muted-color">
          {{ $t(`agent.editor.chatbot.styling.colorHelp.${slot}`) }}
        </small>
        <CInputColorPicker v-model="styling.colors[slot]" />
      </div>
    </div>

    <Divider align="left">
      <span class="text-xs uppercase tracking-wide text-muted-color">
        {{ $t('agent.editor.chatbot.styling.typography') }}
      </span>
    </Divider>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('agent.editor.chatbot.styling.fontFamily.label') }}
        </label>
        <small class="text-muted-color">
          {{ $t('agent.editor.chatbot.styling.fontFamily.help') }}
        </small>
        <InputText v-model="styling.fontFamily" placeholder="Inter, sans-serif" />
      </div>
      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('agent.editor.chatbot.styling.fontSize.heading') }}
        </label>
        <small class="text-muted-color">
          {{ $t('agent.editor.chatbot.styling.fontSizeHelp.heading') }}
        </small>
        <PxInput v-model="styling.fontSizes.heading" :min="8" :max="48" />
      </div>
      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('agent.editor.chatbot.styling.fontSize.base') }}
        </label>
        <small class="text-muted-color">
          {{ $t('agent.editor.chatbot.styling.fontSizeHelp.base') }}
        </small>
        <PxInput v-model="styling.fontSizes.base" :min="8" :max="48" />
      </div>
    </div>

    <!-- Button (was Launcher) -->
    <Divider align="left">
      <span class="text-xs uppercase tracking-wide text-muted-color">
        {{ $t('agent.editor.chatbot.styling.buttonTitle') }}
      </span>
    </Divider>
    <!-- Row 1: Icon + Label -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="flex flex-col gap-2">
        <CInputToggleCard
          v-model="styling.launcher.iconVisible"
          :label="$t('agent.editor.chatbot.styling.button.icon.label')"
          :description="$t('agent.editor.chatbot.styling.button.icon.help')"
        />
        <CFileDropZone
          v-if="styling.launcher.iconVisible"
          accept="image/*"
          :uploading="iconUploading"
          :error="iconError"
          :clearable="!!styling.launcher.iconURL"
          :drop-label="$t('agent.editor.chatbot.styling.button.icon.placeholder')"
          :label="$t('agent.editor.chatbot.styling.button.icon.label')"
          compact
          @select="onIconSelect"
          @clear="onIconClear"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('agent.editor.chatbot.styling.button.label.label') }}
        </label>
        <small class="text-muted-color">
          {{ $t('agent.editor.chatbot.styling.button.label.help') }}
        </small>
        <InputText v-model="styling.launcher.buttonLabel" />
      </div>
    </div>

    <!-- Row 2: Colors -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('agent.editor.chatbot.styling.button.color.label') }}
        </label>
        <small class="text-muted-color">
          {{ $t('agent.editor.chatbot.styling.button.color.help') }}
        </small>
        <CInputColorPicker v-model="styling.colors.primary" />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('agent.editor.chatbot.styling.button.colorText.label') }}
        </label>
        <small class="text-muted-color">
          {{ $t('agent.editor.chatbot.styling.button.colorText.help') }}
        </small>
        <CInputColorPicker v-model="styling.colors.primaryText" />
      </div>
    </div>

    <!-- Row 3: Size + Shape -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('agent.editor.chatbot.styling.button.size.label') }}
        </label>
        <small class="text-muted-color">
          {{ $t('agent.editor.chatbot.styling.button.size.help') }}
        </small>
        <PxInput v-model="styling.launcher.size" :min="32" :max="96" />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('agent.editor.chatbot.styling.button.shape.label') }}
        </label>
        <small class="text-muted-color">
          {{ $t('agent.editor.chatbot.styling.button.shape.help') }}
        </small>
        <Select
          v-model="styling.launcher.shape"
          :options="[
            { label: $t('agent.editor.chatbot.styling.button.shape.circle'), value: 'circle' },
            { label: $t('agent.editor.chatbot.styling.button.shape.square'), value: 'square' },
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

const contentColorSlots = ['background', 'text', 'userBubble', 'agentBubble']


// Logo (header image) upload
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

// Button icon upload
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
