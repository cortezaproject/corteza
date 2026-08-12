<template>
  <div class="flex flex-col gap-6">
    <!-- View mode -->
    <div class="flex flex-col gap-2">
      <label class="font-medium text-muted-color text-sm">
        {{ $t('field.kind.file.view.modeLabel') }}
      </label>
      <div class="flex flex-col gap-2">
        <div class="flex items-center gap-2">
          <RadioButton inputId="modeList" v-model="field.options.mode" value="list" />
          <label for="modeList" class="cursor-pointer">{{ $t('field.kind.file.view.list') }}</label>
        </div>
        <div class="flex items-center gap-2">
          <RadioButton inputId="modeGallery" v-model="field.options.mode" value="gallery" />
          <label for="modeGallery" class="cursor-pointer">
            {{ $t('field.kind.file.view.gallery') }}
          </label>
        </div>
      </div>
      <small class="text-muted-color">{{ $t('field.kind.file.view.modeFootnote') }}</small>
    </div>

    <!-- Max size -->
    <div class="flex flex-col gap-2">
      <label class="font-medium text-muted-color text-sm">
        {{ $t('field.kind.file.view.maxSizeLabel') }}
      </label>
      <InputNumber v-model="field.options.maxSize" :min="0" show-buttons class="w-full md:w-1/2" />
    </div>

    <!-- MIME types -->
    <div class="flex flex-col gap-2">
      <label class="font-medium text-muted-color text-sm">
        {{ $t('field.kind.file.view.mimetypesLabel') }}
      </label>
      <InputText v-model="field.options.mimetypes" class="w-full" />
      <small class="text-muted-color">{{ $t('field.kind.file.view.mimetypesFootnote') }}</small>
    </div>

    <!-- General options -->
    <div class="flex flex-col gap-3">
      <div class="flex items-center gap-2">
        <Checkbox v-model="field.options.clickToView" inputId="clickToView" :binary="true" />
        <label for="clickToView" class="cursor-pointer">
          {{ $t('field.kind.file.view.clickToView') }}
        </label>
      </div>
      <div class="flex items-center gap-2">
        <Checkbox v-model="field.options.enableDownload" inputId="enableDownload" :binary="true" />
        <label for="enableDownload" class="cursor-pointer">
          {{ $t('field.kind.file.view.enableDownload') }}
        </label>
      </div>
      <div class="flex items-center gap-2">
        <Checkbox v-model="field.options.enableWebcam" inputId="enableWebcam" :binary="true" />
        <label for="enableWebcam" class="cursor-pointer">
          {{ $t('field.kind.file.view.webcam.enable.label') }}
        </label>
      </div>
      <div v-if="field.options.mode === 'gallery'" class="flex items-center gap-2">
        <Checkbox v-model="field.options.hideFileName" inputId="hideFileName" :binary="true" />
        <label for="hideFileName" class="cursor-pointer">
          {{ $t('field.kind.file.view.showName') }}
        </label>
      </div>
    </div>

    <!-- Gallery preview styling -->
    <div v-if="field.options.mode === 'gallery'" class="flex flex-col gap-4">
      <h4 class="font-semibold text-base m-0">{{ $t('field.kind.file.view.previewStyle') }}</h4>
      <small class="text-muted-color -mt-2">{{ $t('field.kind.file.view.description') }}</small>

      <div class="grid grid-cols-2 gap-4">
        <div class="flex flex-col gap-1">
          <label class="text-sm text-muted-color">{{ $t('field.kind.file.view.height') }}</label>
          <InputText v-model="field.options.height" placeholder="200px" class="w-full" />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm text-muted-color">{{ $t('field.kind.file.view.width') }}</label>
          <InputText v-model="field.options.width" placeholder="200px" class="w-full" />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm text-muted-color">{{ $t('field.kind.file.view.maxHeight') }}</label>
          <InputText v-model="field.options.maxHeight" placeholder="300px" class="w-full" />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm text-muted-color">{{ $t('field.kind.file.view.maxWidth') }}</label>
          <InputText v-model="field.options.maxWidth" placeholder="300px" class="w-full" />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm text-muted-color">
            {{ $t('field.kind.file.view.borderRadius') }}
          </label>
          <InputText v-model="field.options.borderRadius" placeholder="4px" class="w-full" />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm text-muted-color">{{ $t('field.kind.file.view.margin') }}</label>
          <InputText v-model="field.options.margin" placeholder="auto" class="w-full" />
        </div>
      </div>

      <div class="flex items-center gap-3">
        <label class="text-sm text-muted-color">{{ $t('field.kind.file.view.background') }}</label>
        <CInputColorPicker v-model="field.options.backgroundColor" show-text />
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'
import { components } from '@planetcrust/human-vue'

const { CInputColorPicker } = components

const field = inject('fieldDraft')
</script>
