<template>
  <div class="h-full flex flex-col">
    <header class="px-6 py-5 border-b border-surface shrink-0">
      <h1 class="text-xl font-semibold flex items-center gap-2">
        <i v-if="step === 2 && activeCat" :class="['pi', activeCat.icon, activeCat.color]" />
        {{ title }}
      </h1>
      <p class="text-sm text-muted-color mt-1">{{ hint }}</p>
    </header>

    <div class="flex-1 min-h-0 overflow-y-auto p-6">
      <!-- Step 1 — pick a category (icon + label tiles, per the demo) -->
      <div v-if="step === 1" class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-3">
        <button
          v-for="cat in EVENT_CATEGORIES"
          :key="cat.key"
          type="button"
          class="rounded-xl border border-surface bg-surface-50 dark:bg-surface-950 p-4 flex flex-col items-center justify-center gap-2 text-center hover:border-primary hover:bg-primary/5 transition-colors"
          @click="selectCategory(cat.key)"
        >
          <i :class="['pi text-[22px]', cat.icon, cat.color]" />
          <span class="text-xs font-semibold text-muted-color">
            {{ $t(`project.dashboard.event.categories.${cat.key}.label`) }}
          </span>
        </button>
      </div>

      <!-- Step 2 — the selected category's form (two-column, like the demo) -->
      <div v-else class="max-w-5xl">
        <GovernanceForm :schema="schema" v-model="model" :columns="2" />
      </div>
    </div>

    <footer v-if="step === 2" class="px-6 py-4 border-t border-surface shrink-0 flex items-center gap-2">
      <Button
        :label="$t('project.dashboard.event.back')"
        icon="pi pi-arrow-left"
        severity="secondary"
        text
        @click="step = 1"
      />
      <Button
        :label="$t('project.dashboard.event.create')"
        icon="pi pi-check"
        class="ml-auto"
        @click="submit"
      />
    </footer>
  </div>
</template>

<script setup>
import GovernanceForm from '@/sections/project/components/wizard/GovernanceForm.vue'
import { EVENT_CATEGORIES, EVENT_FORMS } from '@/sections/project/config/eventForm'
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const $toast = inject('$toast')

const step = ref(1)
const category = ref('')
const model = ref({})

const schema = computed(() => EVENT_FORMS[category.value] || [])
const activeCat = computed(() => EVENT_CATEGORIES.find(c => c.key === category.value) || null)

// Title + description reflect the picked category once we're on step 2.
const title = computed(() =>
  step.value === 1
    ? t('project.dashboard.event.title')
    : t('project.dashboard.event.step2', {
        category: t(`project.dashboard.event.categories.${category.value}.label`),
      }),
)
const hint = computed(() =>
  step.value === 1
    ? t('project.dashboard.event.step1')
    : t(`project.dashboard.event.categories.${category.value}.desc`),
)

function selectCategory(key) {
  category.value = key
  model.value = {}
  step.value = 2
}

// No backend yet: acknowledge and return to the Events view; data is discarded.
function submit() {
  $toast.toastSuccess(t('project.dashboard.event.toast.created'))
  router.push({ name: 'project.overview.events', params: { projectId: route.params.projectId } })
}
</script>
