<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('list.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <CResourceList
      ref="resourceListRef"
      primary-key="automationID"
      :fields="listFields"
      :items="items"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :action-items="getActionsMenuItems"
      :translations="{
        searchPlaceholder: $t('list.filterForm.query.placeholder'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('list.new'),
        resourcePlural: $t('list.title'),
      }"
      clickable
      class="h-full"
      @update:filter="Object.assign(filter, $event)"
      @sort="handleSort"
      @row-click="handleRowClick"
      @page-change="handlePageChange"
    >
      <template #header>
        <div class="flex gap-2">
          <Button
            :label="$t('list.new')"
            icon="pi pi-plus"
            size="small"
            @click="openCreateDialog"
          />

          <CPermissionsButton
            v-if="canGrant"
            v-tooltip.bottom="$t('general.label.permissions')"
            resource="corteza::automation:ng-automation/*"
          />
        </div>
      </template>

      <template #body-name="{ data }">
        <div class="flex flex-col">
          <span>{{ data.meta?.short || $t('list.untitled') }}</span>
          <span v-if="data.meta?.description" class="text-xs text-muted-color truncate max-w-full">
            {{ data.meta.description }}
          </span>
        </div>
      </template>

      <template #body-enabled="{ data }">
        <Tag
          :value="data.enabled ? $t('list.active') : $t('list.disabled')"
          :severity="data.enabled ? 'success' : 'secondary'"
        />
      </template>

      <template #body-updatedAt="{ data }">
        {{ locFullDateTime(data.deletedAt || data.updatedAt || data.createdAt) }}
      </template>

      <template #filter>
        <Button
          icon="pi pi-filter"
          severity="secondary"
          outlined
          size="small"
          @click="toggleFilterMenu"
        />
      </template>
    </CResourceList>

    <Popover ref="filterMenu">
      <div class="flex flex-col gap-4 p-2 w-64">
        <div class="flex flex-col gap-2">
          <span class="font-medium text-sm text-primary">
            {{ $t('list.filterForm.deleted.label') }}
          </span>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del0" value="0" />
            <label for="del0" class="text-sm cursor-pointer">
              {{ $t('list.filterForm.excluded.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del1" value="1" />
            <label for="del1" class="text-sm cursor-pointer">
              {{ $t('list.filterForm.inclusive.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del2" value="2" />
            <label for="del2" class="text-sm cursor-pointer">
              {{ $t('list.filterForm.exclusive.label') }}
            </label>
          </div>
        </div>

        <div class="flex flex-col gap-2">
          <span class="font-medium text-sm text-primary">
            {{ $t('list.filterForm.disabled.label') }}
          </span>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.disabled" inputId="dis0" value="0" />
            <label for="dis0" class="text-sm cursor-pointer">
              {{ $t('list.filterForm.excluded.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.disabled" inputId="dis1" value="1" />
            <label for="dis1" class="text-sm cursor-pointer">
              {{ $t('list.filterForm.inclusive.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.disabled" inputId="dis2" value="2" />
            <label for="dis2" class="text-sm cursor-pointer">
              {{ $t('list.filterForm.exclusive.label') }}
            </label>
          </div>
        </div>
      </div>
    </Popover>

    <!-- Create TAQ Dialog -->
    <Dialog
      v-model:visible="showCreateDialog"
      modal
      :header="$t('list.dialog.create.header')"
      :style="{ width: '500px' }"
      @hide="resetCreateForm"
    >
      <div class="flex flex-col gap-6">
        <!-- Name field -->
        <div class="flex flex-col gap-2">
          <label for="automation-name" class="font-medium text-primary"
            >{{ $t('list.dialog.create.name') }} <span class="text-red-500">*</span></label
          >
          <InputText
            id="automation-name"
            v-model="newAutomationName"
            :placeholder="$t('list.dialog.create.namePlaceholder')"
            class="w-full"
            :invalid="!!nameError"
            autofocus
            @keyup.enter="createAutomation"
          />
          <small v-if="nameError" class="text-red-500">{{ nameError }}</small>
        </div>

        <!-- Description field -->
        <div class="flex flex-col gap-2">
          <label for="automation-description" class="font-medium text-primary">{{ $t('list.dialog.create.description') }}</label>
          <Textarea
            id="automation-description"
            v-model="newAutomationDescription"
            :placeholder="$t('list.dialog.create.descriptionPlaceholder')"
            class="w-full"
            rows="3"
          />
        </div>
      </div>
      <template #footer>
        <Button :label="$t('list.button.cancel')" text @click="showCreateDialog = false" />
        <Button
          :label="$t('list.button.create')"
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
import {
  components,
  filters,
  useConfirmDelete,
  useRBACStore,
  useResourceList,
  usePermissions,
} from '@cortezaproject/corteza-vue-next'
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const { CResourceList } = components
const { locFullDateTime } = filters

const router = useRouter()
const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()

const $toast = inject('$toast')
const $AutomationAPI = inject('$AutomationAPI')
const $Auth = inject('$Auth')
const rbac = useRBACStore()
const canGrant = computed(() => rbac.can('automation/', 'grant'))
const { open: openPermissions } = usePermissions()

const resourceListRef = ref()

// Filter menu
const filterMenu = ref()
function toggleFilterMenu(event) {
  filterMenu.value.toggle(event)
}

// Column definitions
const listFields = [
  {
    key: 'name',
    sortable: false,
    header: t('list.columns.name'),
  },
  {
    key: 'enabled',
    sortable: false,
    header: t('list.columns.enabled'),
  },
  {
    key: 'updatedAt',
    sortable: true,
    header: t('list.columns.updatedAt'),
    class: 'text-right',
    pt: {
      columnHeaderContent: 'justify-end',
    },
  },
]

// Resource list composable
const {
  items,
  loading,
  filter,
  sorting,
  pagination,
  handleSort,
  handlePageChange,
  filterList,
} = useResourceList(params => $AutomationAPI.ngAutomationListCancellable(params), {
  filter: { query: '', deleted: '0', disabled: '1' },
  sorting: { sortBy: 'createdAt', sortDesc: true },
  pagination: { limit: 50 },
})

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
    nameError.value = t('list.dialog.create.nameRequired')
    return false
  }

  return true
}

async function createAutomation() {
  if (!validateForm()) return

  creating.value = true

  try {
    const response = await $AutomationAPI.ngAutomationCreate({
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

    showCreateDialog.value = false
    router.push(`/builder/${response.automationID}`)
  } catch (e) {
    console.error('Failed to create automation:', e)
  } finally {
    creating.value = false
  }
}

// Row click navigation
function handleRowClick({ data }) {
  router.push(`/builder/${data.automationID}`)
}

// Actions menu methods
function getActionsMenuItems(automation) {
  const items = []

  if (automation.canGrant || canGrant.value) {
    items.push({
      label: t('general.label.permissions'),
      icon: 'pi pi-lock',
      command: () => {
        resourceListRef.value?.hideActionsMenu?.()
        openPermissions({
          resource: `corteza::automation:ng-automation/${automation.automationID}`,
          title: automation.meta?.short || automation.automationID,
        })
      },
    })
  }

  items.push({
    label: t('general.label.delete'),
    icon: 'pi pi-trash',
    class: 'text-red-500',
    command: () => onConfirmDelete(automation),
  })

  return items
}

function onConfirmDelete(automation) {
  confirmDelete({
    message: t('general.confirm.delete'),
    header: automation.meta?.short || t('list.untitled'),
    onConfirm: () => handleDelete(automation),
  })
}

async function handleDelete(automation) {
  resourceListRef.value.hideActionsMenu()
  try {
    await $AutomationAPI.ngAutomationDelete({
      automationID: automation.automationID,
    })
    $toast.toastSuccess(t('notification.automation.delete.success'))
    filterList()
  } catch (e) {
    console.error('Failed to delete automation:', e)
    $toast.toastErrorHandler(t('notification.automation.delete.error'))(e)
  }
}
</script>
