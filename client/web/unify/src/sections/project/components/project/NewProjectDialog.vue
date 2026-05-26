<template>
  <Dialog
    :visible="visible"
    @update:visible="$emit('update:visible', $event)"
    modal
    :header="step === 1 ? 'Start a new project' : 'Project details'"
    :style="{ width: '46rem' }"
    :pt="{ content: { class: '!pt-2' } }"
  >
    <!-- Step 1 — pick a build mode -->
    <template v-if="step === 1">
      <p class="text-sm text-muted-color mb-4">How would you like to build?</p>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <button
          type="button"
          class="text-left rounded-xl border-2 border-surface hover:border-primary transition-colors p-4 flex flex-col gap-2"
          @click="selectMode('gated')"
        >
          <div class="flex items-center gap-2">
            <i class="pi pi-shield text-xl text-amber-600 dark:text-amber-400" />
            <span class="font-semibold">Gated</span>
            <Tag value="Governance" severity="warn" class="!text-[10px] !py-0 ml-auto" />
          </div>
          <p class="text-sm text-muted-color">
            Steps flow through governance gates; approvers sign off and publishing requires full
            approval. Members &amp; roles are set up inside the project.
          </p>
        </button>

        <button
          type="button"
          class="text-left rounded-xl border-2 border-surface hover:border-primary transition-colors p-4 flex flex-col gap-2"
          @click="selectMode('free')"
        >
          <div class="flex items-center gap-2">
            <i class="pi pi-unlock text-xl text-emerald-600 dark:text-emerald-400" />
            <span class="font-semibold">Free</span>
            <Tag value="No review" severity="secondary" class="!text-[10px] !py-0 ml-auto" />
          </div>
          <p class="text-sm text-muted-color">
            No team, no governance, no gates. Build and publish freely. Best for solo work and quick
            prototypes.
          </p>
        </button>
      </div>
    </template>

    <!-- Step 2 — details + (gated) team -->
    <template v-else>
      <div class="flex flex-col gap-4 pt-1">
        <div class="flex items-center gap-2 text-sm text-muted-color">
          <Tag
            :value="mode.charAt(0).toUpperCase() + mode.slice(1)"
            :severity="mode === 'gated' ? 'warn' : 'secondary'"
            :icon="mode === 'gated' ? 'pi pi-shield' : 'pi pi-unlock'"
          />
        </div>

        <CFormGroup label="Name" required>
          <InputText v-model="name" placeholder="e.g. CRM & Sales" fluid autofocus />
        </CFormGroup>

        <CFormGroup label="Description">
          <Textarea
            v-model="description"
            rows="2"
            auto-resize
            fluid
            placeholder="Optional short description"
          />
        </CFormGroup>
      </div>
    </template>

    <template #footer>
      <div
        class="flex items-center gap-2 w-full"
        :class="step === 2 ? 'justify-between' : 'justify-end'"
      >
        <Button
          v-if="step === 2"
          label="Back"
          icon="pi pi-arrow-left"
          severity="secondary"
          text
          size="small"
          @click="step = 1"
        />
        <div class="flex gap-2">
          <Button label="Cancel" severity="secondary" outlined size="small" @click="close" />
          <Button
            v-if="step === 2"
            label="Create project"
            size="small"
            :disabled="!canCreate"
            @click="onCreate"
          />
        </div>
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed, ref, watch } from 'vue'

const props = defineProps({
  visible: { type: Boolean, default: false },
})
const emit = defineEmits(['update:visible', 'created'])

const store = useProjectsStore()

const step = ref(1)
const mode = ref('')
const name = ref('')
const description = ref('')

const canCreate = computed(() => !!name.value.trim())

function reset() {
  step.value = 1
  mode.value = ''
  name.value = ''
  description.value = ''
}

watch(
  () => props.visible,
  v => {
    if (v) reset()
  },
)

function selectMode(value) {
  mode.value = value
  step.value = 2
}

function close() {
  emit('update:visible', false)
}

function onCreate() {
  if (!canCreate.value) return
  const project = store.create({
    name: name.value,
    description: description.value,
    mode: mode.value,
  })
  emit('created', project)
  close()
}
</script>
