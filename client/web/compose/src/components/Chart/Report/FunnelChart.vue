<template>
  <ReportEdit
    :report="report"
    :chart="chart"
    :modules="modules"
    :supported-metrics="supportedMetrics"
    @update:report="$emit('update:report', $event)"
  >
    <template #metric-options="{ metric }">
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 mt-4">
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.metric.labelLabel') }}
          </label>
          <InputText v-model="metric.label" class="w-full" />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.metric.options.label') }}
          </label>
          <div class="flex flex-col gap-2">
            <div class="flex items-center gap-2">
              <Checkbox v-model="metric.fixTooltips" :binary="true" input-id="funnelFixTooltips" />
              <label for="funnelFixTooltips">{{ $t('chart.edit.metric.fixTooltips') }}</label>
            </div>
            <div class="flex items-center gap-2">
              <Checkbox v-model="metric.cumulative" :binary="true" input-id="funnelCumulative" />
              <label for="funnelCumulative">{{ $t('chart.edit.metric.cumulative') }}</label>
            </div>
          </div>
        </div>
      </div>

      <Divider />

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.formatting.prefix.label') }}
          </label>
          <InputText
            v-model="metric.formatting.prefix"
            :placeholder="$t('chart.edit.formatting.prefix.placeholder')"
            class="w-full"
          />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.formatting.suffix.label') }}
          </label>
          <InputText
            v-model="metric.formatting.suffix"
            :placeholder="$t('chart.edit.formatting.suffix.placeholder')"
            class="w-full"
          />
        </div>
      </div>
    </template>
  </ReportEdit>
</template>

<script setup>
import ReportEdit from './ReportEdit.vue'

defineProps({
  report: {
    type: Object,
    required: true,
  },
  chart: {
    type: Object,
    default: () => ({}),
  },
  modules: {
    type: Array,
    required: true,
  },
  supportedMetrics: {
    type: Number,
    default: -1,
  },
})

defineEmits(['update:report'])
</script>
