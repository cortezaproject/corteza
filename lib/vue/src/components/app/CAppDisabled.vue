<template>
  <div class="flex items-center justify-center h-screen">
    <div class="flex flex-col items-center gap-4 text-center max-w-md p-8">
      <i class="pi pi-ban text-6xl text-muted-color" />
      <h1 class="text-2xl font-semibold text-color">
        {{ t('general.appDisabled.title') }}
      </h1>
      <p class="text-muted-color">
        {{ t('general.appDisabled.message') }}
      </p>
      <a :href="homeURL" class="mt-2" @click="onAnchorClick($event, homeURL)">
        <Button
          :label="t('general.appDisabled.backToHome')"
          icon="pi pi-home"
          severity="secondary"
          outlined
        />
      </a>
    </div>
  </div>
</template>

<script setup>
import { useInternalLink } from '../../composables/useInternalLink'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const { onAnchorClick } = useInternalLink()
const homeURL = window.location.origin
</script>
