<template>
  <Dialog
    :visible="visible"
    @update:visible="$emit('update:visible', $event)"
    modal
    :header="headerText"
    :style="{ width: '32rem' }"
    :pt="{ content: { class: '!pt-2' } }"
  >
    <!-- Step 1 — details -->
    <template v-if="step === 1">
      <div class="flex flex-col gap-4 pt-1">
        <CFormGroup :label="$t('general.label.name')" required>
          <InputText
            v-model="name"
            :placeholder="$t('project.newDialog.namePlaceholder')"
            fluid
            autofocus
          />
        </CFormGroup>

        <CFormGroup :label="$t('general.label.description')">
          <Textarea
            v-model="description"
            rows="2"
            auto-resize
            fluid
            :placeholder="$t('project.newDialog.descriptionPlaceholder')"
          />
        </CFormGroup>
      </div>
    </template>

    <!-- Step 2 — Deployer categories (drive whether a FRIA is required) -->
    <template v-else>
      <div class="flex flex-col gap-4 pt-1">
        <p class="text-sm text-muted-color">
          {{ $t('project.deployer.prompt') }}
        </p>

        <div
          v-for="(q, i) in DEPLOYER_QUESTIONS"
          :key="q.key"
          class="flex items-start gap-4 rounded-lg border border-surface p-3"
        >
          <p class="text-sm flex-1 min-w-0">
            <span class="font-medium mr-1">{{ i + 1 }}.</span>
            {{ $t(q.labelKey) }}
          </p>
          <SelectButton
            v-model="deployer[q.key]"
            :options="YES_NO"
            option-label="label"
            option-value="value"
            :allow-empty="false"
            class="shrink-0"
          />
        </div>
      </div>
    </template>

    <template #footer>
      <div
        class="flex items-center gap-2 w-full"
        :class="step > 1 ? 'justify-between' : 'justify-end'"
      >
        <Button
          v-if="step > 1"
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          text
          size="small"
          @click="step = step - 1"
        />
        <div class="flex gap-2">
          <Button
            :label="$t('general.label.cancel')"
            severity="secondary"
            text
            size="small"
            @click="close"
          />
          <Button
            v-if="step === 1"
            :label="$t('general.label.next')"
            icon="pi pi-arrow-right"
            icon-pos="right"
            size="small"
            :disabled="!canCreate"
            @click="step = 2"
          />
          <Button
            v-else
            :label="$t('project.newDialog.createProject')"
            size="small"
            :disabled="!canCreate"
            :loading="creating"
            @click="onCreate"
          />
        </div>
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed, inject, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  visible: { type: Boolean, default: false },
})
const emit = defineEmits(['update:visible', 'created'])

const { t } = useI18n()
const store = useProjectsStore()
const $toast = inject('$toast')

// AI Act Deployer categories. Answering "yes" to any makes a Fundamental
// Rights Impact Assessment (FRIA) required for the project. `labelKey` is an
// i18n key resolved with $t at render. Every project collects these — there
// is no build-mode gate any more.
const DEPLOYER_QUESTIONS = [
  { key: 'publicAuthorityAnnex3', labelKey: 'project.deployer.questions.publicAuthorityAnnex3' },
  {
    key: 'privateEssentialServices',
    labelKey: 'project.deployer.questions.privateEssentialServices',
  },
  { key: 'insuranceBanking', labelKey: 'project.deployer.questions.insuranceBanking' },
]
const YES_NO = [
  { label: t('general.label.yes'), value: true },
  { label: t('general.label.no'), value: false },
]

const step = ref(1)
const name = ref('')
const description = ref('')
const deployer = reactive(defaultDeployer())

function defaultDeployer() {
  return Object.fromEntries(DEPLOYER_QUESTIONS.map(q => [q.key, false]))
}

const canCreate = computed(() => !!name.value.trim())

const headerText = computed(
  () =>
    ({
      1: t('project.newDialog.headerDetails'),
      2: t('project.newDialog.headerDeployer'),
    })[step.value],
)

function reset() {
  step.value = 1
  name.value = ''
  description.value = ''
  Object.assign(deployer, defaultDeployer())
}

watch(
  () => props.visible,
  v => {
    if (v) reset()
  },
)

function close() {
  emit('update:visible', false)
}

const creating = ref(false)

async function onCreate() {
  if (!canCreate.value || creating.value) return
  creating.value = true
  try {
    const project = await store.create({
      name: name.value,
      description: description.value,
      deployer: { ...deployer },
    })
    emit('created', project)
    close()
  } catch (err) {
    // Dialog stays open so nothing typed is lost.
    $toast.toastErrorHandler(t('project.newDialog.toastCreateFailed'))(err)
  } finally {
    creating.value = false
  }
}
</script>
