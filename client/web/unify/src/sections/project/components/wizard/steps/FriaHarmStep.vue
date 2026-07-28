<template>
  <FriaScenarioEditorShell :project="project" :disabled="disabled" v-slot="{ scenario, update }">
    <div class="flex flex-col gap-6 max-w-3xl">
      <div class="flex items-start gap-3 rounded-lg border border-surface bg-emphasis p-4">
        <i class="pi pi-info-circle text-primary mt-0.5" />
        <div class="text-xs leading-relaxed">
          <p class="font-medium mb-1">{{ $t('fria.harm.exampleTitle') }}</p>
          <p class="text-muted-color italic">{{ $t('fria.harm.exampleBody') }}</p>
        </div>
      </div>

      <CFormGroup :label="$t('fria.harm.titleLabel')" required>
        <InputText
          :model-value="scenario.title"
          fluid
          maxlength="120"
          :disabled="disabled"
          :placeholder="$t('fria.harm.titlePlaceholder')"
          @update:model-value="v => update({ title: v })"
        />
      </CFormGroup>

      <CFormGroup :label="$t('fria.harm.descriptionLabel')" required>
        <Textarea
          :model-value="scenario.description"
          rows="6"
          auto-resize
          fluid
          maxlength="2000"
          :disabled="disabled"
          :placeholder="$t('fria.harm.descriptionPlaceholder')"
          @update:model-value="v => update({ description: v })"
        />
      </CFormGroup>

      <CFormGroup :label="$t('fria.harm.severityLabel')" required>
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
          <button
            v-for="level in FRIA_SEVERITY_LEVELS"
            :key="level"
            type="button"
            class="rounded-lg border-2 p-4 text-center transition-colors"
            :class="[
              scenario.severity === level
                ? FRIA_SEVERITY_CARD_CLASSES[level]
                : 'border-surface hover:border-primary/60',
              disabled ? 'cursor-not-allowed opacity-70' : 'cursor-pointer',
            ]"
            :disabled="disabled"
            @click="update({ severity: level })"
          >
            <div class="text-sm font-semibold mb-1">{{ $t(`fria.severity.${level}.label`) }}</div>
            <div class="text-xs text-muted-color leading-snug">
              {{ $t(`fria.severity.${level}.description`) }}
            </div>
          </button>
        </div>
      </CFormGroup>
    </div>
  </FriaScenarioEditorShell>
</template>

<script setup>
import FriaScenarioEditorShell from '@/sections/project/components/wizard/fria/FriaScenarioEditorShell.vue'
import {
  FRIA_SEVERITY_CARD_CLASSES,
  FRIA_SEVERITY_LEVELS,
} from '@/sections/project/config/friaScenario'

defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})
</script>
