<template>
  <Panel :header="$t('chatbot.editor.panels.styling')" toggleable>
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
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
          :error="logoError || uploadDisabledLabel"
          :preview-url="logoPreviewUrl"
          :clearable="!!styling.logoURL"
          :disabled="!canUpload"
          :drop-label="
            canUpload ? $t('chatbot.editor.styling.logo.placeholder') : uploadDisabledLabel
          "
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
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
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
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('chatbot.editor.styling.fontSize.heading') }}
        </label>
        <small class="text-muted-color">
          {{ $t('chatbot.editor.styling.fontSizeHelp.heading') }}
        </small>
        <InputNumber v-model="fontHeadingPx" :min="8" :max="48" suffix=" px" fluid />
      </div>
      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('chatbot.editor.styling.fontSize.base') }}
        </label>
        <small class="text-muted-color">
          {{ $t('chatbot.editor.styling.fontSizeHelp.base') }}
        </small>
        <InputNumber v-model="fontBasePx" :min="8" :max="48" suffix=" px" fluid />
      </div>
    </div>

    <Divider align="left">
      <span class="text-xs uppercase tracking-wide text-muted-color">
        {{ $t('chatbot.editor.styling.buttonTitle') }}
      </span>
    </Divider>
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('chatbot.editor.styling.button.label.label') }}
        </label>
        <small class="text-muted-color">
          {{ $t('chatbot.editor.styling.button.label.help') }}
        </small>
        <InputText v-model="styling.launcher.buttonLabel" />
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
            {
              label: $t('chatbot.editor.styling.button.position.bottomRight'),
              value: 'bottom-right',
            },
            {
              label: $t('chatbot.editor.styling.button.position.bottomCenter'),
              value: 'bottom-center',
            },
            {
              label: $t('chatbot.editor.styling.button.position.bottomLeft'),
              value: 'bottom-left',
            },
            {
              label: $t('chatbot.editor.styling.button.position.leftMiddle'),
              value: 'left-middle',
            },
            {
              label: $t('chatbot.editor.styling.button.position.rightMiddle'),
              value: 'right-middle',
            },
            { label: $t('chatbot.editor.styling.button.position.topRight'), value: 'top-right' },
            { label: $t('chatbot.editor.styling.button.position.topCenter'), value: 'top-center' },
            { label: $t('chatbot.editor.styling.button.position.topLeft'), value: 'top-left' },
          ]"
          optionLabel="label"
          optionValue="value"
        />
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 mt-4">
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
          {{ $t('chatbot.editor.styling.button.size.label') }}
        </label>
        <small class="text-muted-color">
          {{ $t('chatbot.editor.styling.button.size.help') }}
        </small>
        <InputNumber v-model="launcherSizePx" :min="32" :max="96" suffix=" px" fluid />
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 mt-4">
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
          :error="iconError || uploadDisabledLabel"
          :preview-url="iconPreviewUrl"
          :clearable="!!styling.launcher.iconURL"
          :disabled="!canUpload"
          :drop-label="
            canUpload ? $t('chatbot.editor.styling.button.icon.placeholder') : uploadDisabledLabel
          "
          :label="$t('chatbot.editor.styling.button.icon.label')"
          compact
          preview-max-width="100%"
          preview-max-height="140px"
          @select="onIconSelect"
          @clear="onIconClear"
        />
      </div>

      <div class="flex flex-col gap-2">
        <CInputToggleCard
          v-model="styling.launcher.startOpen"
          :label="$t('chatbot.editor.styling.button.startOpen.label')"
          :description="$t('chatbot.editor.styling.button.startOpen.help')"
        />
      </div>
    </div>
  </Panel>
</template>

<script setup>
import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import { components, useFileUpload } from '@planetcrust/human-vue'

const { CInputColorPicker, CFileDropZone, CInputToggleCard } = components

const props = defineProps({
  styling: { type: Object, required: true },
  chatbotId: { type: String, default: '' },
})

const $SystemAPI = inject('$SystemAPI')
const { t } = useI18n()

const canUpload = computed(() => !!props.chatbotId)

function apiOrigin() {
  try {
    return new URL($SystemAPI.baseURL).origin
  } catch {
    return window.location.origin
  }
}
const uploadDisabledLabel = computed(() =>
  !canUpload.value ? t('chatbot.editor.styling.uploadNeedsSave') : '',
)

function pxModel(target, key) {
  return computed({
    get: () => {
      const v = target[key]
      if (typeof v === 'number') return v
      const n = parseFloat(v)
      return Number.isFinite(n) ? n : null
    },
    set: v => {
      target[key] = v == null ? '' : `${v}px`
    },
  })
}

const fontHeadingPx = pxModel(props.styling.fontSizes, 'heading')
const fontBasePx = pxModel(props.styling.fontSizes, 'base')
const launcherSizePx = pxModel(props.styling.launcher, 'size')

const contentColorSlots = [
  'header',
  'headerText',
  'background',
  'text',
  'userBubble',
  'agentBubble',
]

const {
  uploading: logoUploading,
  uploadError: logoUploadError,
  uploadFiles: uploadLogoFiles,
  reset: resetLogoUpload,
} = useFileUpload()

const logoError = computed(() => logoUploadError.value)

function absolutize(u) {
  if (!u) return ''
  if (/^(https?:|blob:|data:)/i.test(u)) return u
  if (u.startsWith('/')) return apiOrigin() + u
  return u
}

const logoPreviewUrl = computed(() => absolutize(props.styling.logoURL))

let logoBlobUrl = ''
async function onLogoSelect(files) {
  const file = files[0]
  if (!file || !canUpload.value) return
  if (logoBlobUrl) URL.revokeObjectURL(logoBlobUrl)
  logoBlobUrl = URL.createObjectURL(file)
  props.styling.logoURL = logoBlobUrl
  try {
    const results = await uploadLogoFiles([file], {
      api: $SystemAPI,
      endpoint: $SystemAPI.chatbotUploadAssetEndpoint({ chatbotID: props.chatbotId }),
    })
    const rsp = results[0]
    if (rsp) {
      props.styling.logoAttachmentID = rsp.attachmentID
      props.styling.logoURL = apiOrigin() + rsp.url
    }
  } catch {
    /* error set by composable */
  }
}

function onLogoClear() {
  if (logoBlobUrl) {
    URL.revokeObjectURL(logoBlobUrl)
    logoBlobUrl = ''
  }
  props.styling.logoAttachmentID = ''
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
const iconPreviewUrl = computed(() => absolutize(props.styling.launcher.iconURL))

let iconBlobUrl = ''
async function onIconSelect(files) {
  const file = files[0]
  if (!file || !canUpload.value) return
  if (iconBlobUrl) URL.revokeObjectURL(iconBlobUrl)
  iconBlobUrl = URL.createObjectURL(file)
  props.styling.launcher.iconURL = iconBlobUrl
  try {
    const results = await uploadIconFiles([file], {
      api: $SystemAPI,
      endpoint: $SystemAPI.chatbotUploadAssetEndpoint({ chatbotID: props.chatbotId }),
    })
    const rsp = results[0]
    if (rsp) {
      props.styling.launcher.iconAttachmentID = rsp.attachmentID
      props.styling.launcher.iconURL = apiOrigin() + rsp.url
    }
  } catch {
    /* error set by composable */
  }
}

function onIconClear() {
  if (iconBlobUrl) {
    URL.revokeObjectURL(iconBlobUrl)
    iconBlobUrl = ''
  }
  props.styling.launcher.iconAttachmentID = ''
  props.styling.launcher.iconURL = ''
  resetIconUpload()
}
</script>
