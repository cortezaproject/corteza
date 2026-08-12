<template>
  <div class="flex flex-col gap-4">
    <div class="flex items-center gap-2">
      <Checkbox v-model="field.options.trimPath" inputId="trimPath" :binary="true" />
      <label for="trimPath" class="cursor-pointer">{{ $t('field.kind.url.trimToDomain') }}</label>
    </div>
    <div class="flex items-center gap-2">
      <Checkbox v-model="field.options.trimFragment" inputId="trimFragment" :binary="true" />
      <label for="trimFragment" class="cursor-pointer">{{ $t('field.kind.url.trimHash') }}</label>
    </div>
    <div class="flex items-center gap-2">
      <Checkbox v-model="field.options.trimQuery" inputId="trimQuery" :binary="true" />
      <label for="trimQuery" class="cursor-pointer">
        {{ $t('field.kind.url.trimQuestionMark') }}
      </label>
    </div>
    <div class="flex items-center gap-2">
      <Checkbox v-model="field.options.onlySecure" inputId="onlySecure" :binary="true" />
      <label for="onlySecure" class="cursor-pointer">{{ $t('field.kind.url.sshOnly') }}</label>
    </div>
    <div class="flex items-center gap-2">
      <Checkbox v-model="field.options.outputPlain" inputId="outputPlain" :binary="true" />
      <label for="outputPlain" class="cursor-pointer">
        {{ $t('field.kind.url.preventToLink') }}
      </label>
    </div>

    <CConfiguratorMultiDelimiter />
  </div>
</template>

<script setup>
import { inject } from 'vue'
import CConfiguratorMultiDelimiter from '../CConfiguratorMultiDelimiter.vue'

const field = inject('fieldDraft')
</script>
