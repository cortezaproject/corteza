<template>
  <div class="p-5">
    <h1 class="text-2xl font-semibold mb-4">
      {{ isNew ? $t('admin.charts.create') : $t('admin.charts.edit') }}
    </h1>
    <p class="text-muted-color">{{ $t('admin.charts.editor.placeholder') }}</p>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const props = defineProps({
  namespace: {
    type: Object,
    required: true,
  },
})

const route = useRoute()

const isNew = computed(() => !route.params.chartID)
</script>
