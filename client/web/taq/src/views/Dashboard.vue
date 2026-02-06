<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('dashboard.title') }}</span>
  </Teleport>

  <div class="p-6">
    <p class="text-muted-color mb-8">{{ $t('dashboard.welcome') }}</p>

    <!-- Loading state -->
    <div v-if="loading" class="flex justify-center py-12">
      <ProgressSpinner />
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <!-- Create new TAQ -->
      <Card class="cursor-pointer hover:shadow-lg transition-shadow" @click="openCreateDialog">
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-plus" />
            {{ $t('dashboard.createTaq') }}
          </div>
        </template>
        <template #content>
          <p class="text-sm text-muted-color">{{ $t('dashboard.createNew') }}</p>
        </template>
      </Card>

      <!-- Saved TAQs -->
      <Card
        v-for="automation in automations"
        :key="automation.automationID"
        class="cursor-pointer hover:shadow-lg transition-shadow relative group"
        @click="openAutomation(automation.automationID)"
      >
        <template #title>
          <div class="flex items-center justify-between">
            <span>{{ automation.meta?.short || $t('dashboard.untitled') }}</span>
            <Tag v-if="automation.enabled" :value="$t('dashboard.active')" severity="success" class="text-xs" />
          </div>
        </template>
        <template #content>
          <p class="text-sm text-muted-color">
            {{ $t('dashboard.triggers', { count: automation.triggers.length }) }}, {{ $t('dashboard.steps', { count: automation.steps.length }) }}
          </p>
          <p class="text-xs text-muted-color opacity-70">
            {{ $t('dashboard.updated', { date: formatDate(automation.updatedAt) }) }}
          </p>
          <!-- Delete button -->
          <Button
            icon="pi pi-trash"
            severity="danger"
            text
            rounded
            size="small"
            class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity"
            @click.stop="confirmDelete(automation)"
          />
        </template>
      </Card>
    </div>

    <!-- Empty state -->
    <div v-if="!loading && automations.length === 0" class="text-center py-12 text-muted-color">
      <i class="pi pi-inbox text-4xl mb-4" />
      <p>{{ $t('dashboard.noSaved') }}</p>
    </div>

    <!-- Create TAQ Dialog -->
    <Dialog
      v-model:visible="showCreateDialog"
      modal
      :header="$t('dashboard.dialog.create.header')"
      :style="{ width: '500px' }"
      @hide="resetCreateForm"
    >
      <div class="flex flex-col gap-6">
        <!-- Name field -->
        <div class="flex flex-col gap-2">
          <label for="automation-name" class="font-medium text-primary"
            >{{ $t('dashboard.dialog.create.name') }} <span class="text-red-500">*</span></label
          >
          <InputText
            id="automation-name"
            v-model="newAutomationName"
            :placeholder="$t('dashboard.dialog.create.namePlaceholder')"
            class="w-full"
            :invalid="!!nameError"
            autofocus
            @keyup.enter="createAutomation"
          />
          <small v-if="nameError" class="text-red-500">{{ nameError }}</small>
        </div>

        <!-- Description field -->
        <div class="flex flex-col gap-2">
          <label for="automation-description" class="font-medium text-primary">{{ $t('dashboard.dialog.create.description') }}</label>
          <Textarea
            id="automation-description"
            v-model="newAutomationDescription"
            :placeholder="$t('dashboard.dialog.create.descriptionPlaceholder')"
            class="w-full"
            rows="3"
          />
        </div>
      </div>
      <template #footer>
        <Button :label="$t('dashboard.button.cancel')" text @click="showCreateDialog = false" />
        <Button
          :label="$t('dashboard.button.create')"
          icon="pi pi-check"
          :disabled="!newAutomationName.trim() || creating"
          :loading="creating"
          @click="createAutomation"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { useAutomationStore } from '@/stores/automation'
import { useConfirm } from 'primevue/useconfirm'
import { computed, inject, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const { t } = useI18n()
const router = useRouter()
const confirm = useConfirm()
const automationStore = useAutomationStore()

// Inject API client and auth
const $AutomationAPI = inject('$AutomationAPI')
const $Auth = inject('$Auth')

// Store state
const loading = computed(() => automationStore.loading)
const automations = computed(() => automationStore.list)

// Create dialog state
const showCreateDialog = ref(false)
const newAutomationName = ref('')
const newAutomationDescription = ref('')
const creating = ref(false)
const nameError = ref('')

function resetCreateForm() {
  newAutomationName.value = ''
  newAutomationDescription.value = ''
  nameError.value = ''
}

function openCreateDialog() {
  resetCreateForm()
  showCreateDialog.value = true
}

function validateForm() {
  nameError.value = ''

  if (!newAutomationName.value.trim()) {
    nameError.value = t('dashboard.dialog.create.nameRequired')
    return false
  }

  return true
}

async function createAutomation() {
  if (!validateForm()) return

  creating.value = true

  try {
    const automation = await automationStore.create($AutomationAPI, {
      meta: {
        short: newAutomationName.value.trim(),
        description: newAutomationDescription.value.trim() || undefined,
      },
      enabled: false,
      triggers: [],
      steps: [],
      paths: [],
      ownedBy: $Auth?.user?.userID,
    })

    openAutomation(automation.automationID)
  } catch (e) {
    console.error('Failed to create automation:', e)
  } finally {
    creating.value = false
  }
}

async function loadAutomations() {
  try {
    await automationStore.fetchList($AutomationAPI)
  } catch (e) {
    console.error('Failed to load automations:', e)
  }
}

function openAutomation(id) {
  router.push(`/builder/${id}`)
}

function confirmDelete(automation) {
  const name = automation.meta?.short || t('dashboard.untitled')
  confirm.require({
    message: t('dashboard.dialog.delete.confirm', { name }),
    header: t('dashboard.dialog.delete.header'),
    icon: 'pi pi-exclamation-triangle',
    rejectProps: {
      label: t('dashboard.button.cancel'),
      severity: 'secondary',
      outlined: true,
    },
    acceptProps: {
      label: t('dashboard.button.delete'),
      severity: 'danger',
    },
    accept: async () => {
      try {
        await automationStore.remove($AutomationAPI, automation.automationID)
      } catch (e) {
        console.error('Failed to delete automation:', e)
      }
    },
  })
}

function formatDate(dateStr) {
  if (!dateStr) return t('dashboard.unknown')
  const date = new Date(dateStr)
  return date.toLocaleDateString()
}

onMounted(() => {
  loadAutomations()
})
</script>
