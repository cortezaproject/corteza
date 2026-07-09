<template>
  <div class="flex flex-col items-center justify-center gap-4 py-10 text-center">
    <!-- Already live: the project has been published; offer the dashboard. -->
    <template v-if="isPublished">
      <i class="pi pi-check-circle text-4xl text-green-500" />
      <h1 class="text-xl font-medium">{{ $t('project.publishStep.publishedTitle') }}</h1>
      <p class="max-w-prose text-sm text-muted-color">
        {{ $t('project.publishStep.publishedDescription') }}
      </p>
      <Button
        :label="$t('project.publishStep.openDashboard')"
        icon="pi pi-arrow-right"
        icon-pos="right"
        size="small"
        @click="emit('open-dashboard')"
      />
    </template>

    <!-- Draft: review and publish. -->
    <template v-else>
      <i class="pi pi-cloud-upload text-4xl text-surface-400" />
      <h1 class="text-xl font-medium">{{ $t('project.publishStep.title') }}</h1>
      <p class="max-w-prose text-sm text-muted-color">
        {{ $t('project.publishStep.description') }}
      </p>
      <Button
        :label="$t('project.publishStep.action')"
        icon="pi pi-cloud-upload"
        size="small"
        :loading="publishing"
        @click="emit('publish')"
      />
    </template>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  project: { type: Object, required: true },
  publishing: { type: Boolean, default: false },
})

const emit = defineEmits(['publish', 'open-dashboard'])

// 'active' is the backend's live state; 'published' is treated the same for
// forward-compatibility. Either means the project is already live.
const isPublished = computed(() =>
  ['active', 'published'].includes(props.project?.status),
)
</script>
