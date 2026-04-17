<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.apigw.list.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0 flex flex-col gap-4">
    <!-- Profiler & Proxy Settings -->
    <Panel :header="$t('system.apigw.settings.title')" toggleable :collapsed="true">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="flex flex-col gap-2">
          <label class="font-medium text-primary text-sm">
            {{ $t('system.apigw.settings.profiler.label') }}
          </label>
          <SelectButton
            v-model="profilerSetting"
            :options="profilerOptions"
            option-label="label"
            option-value="value"
            @change="onSettingsChanged"
          />
        </div>

        <div class="flex flex-col gap-2">
          <label class="font-medium text-primary text-sm">
            {{ $t('system.apigw.settings.proxy.label') }}
          </label>
          <div class="flex items-center gap-3">
            <ToggleSwitch
              id="proxyFollow"
              v-model="proxyFollowRedirects"
              @change="onSettingsChanged"
            />
            <label for="proxyFollow" class="cursor-pointer">
              {{ $t('system.apigw.settings.proxy.follow') }}
            </label>
          </div>
        </div>
      </div>

      <div class="flex justify-end mt-4">
        <Button
          :label="$t('general.label.save')"
          icon="pi pi-save"
          size="small"
          :loading="savingSettings"
          @click="saveSettings"
        />
      </div>
    </Panel>
    <CResourceList
      ref="resourceListRef"
      primary-key="routeID"
      :fields="fields"
      :items="items"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :action-items="getActionsMenuItems"
      :translations="{
        searchPlaceholder: $t('system.apigw.list.filterForm.query.placeholder'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('system.apigw.list.new'),
        resourcePlural: $t('system.apigw.list.title'),
      }"
      clickable
      class="h-full"
      @update:filter="Object.assign(filter, $event)"
      @sort="handleSort"
      @row-click="
        ({ data }) =>
          $router.push({ name: 'system.apiGateway.edit', params: { routeID: data.routeID } })
      "
      @page-change="handlePageChange"
    >
      <template #header>
        <div class="flex gap-2">
          <Button
            :label="$t('system.apigw.list.new')"
            icon="pi pi-plus"
            size="small"
            @click="$router.push({ name: 'system.apiGateway.create' })"
          />
          <CPermissionsButton
            v-if="canGrant"
            resource="corteza::system:apigw-route/*"
            v-tooltip.bottom="$t('general.label.permissions')"
          />
        </div>
      </template>

      <template #body-method="{ data }">
        <Tag :value="data.method" severity="info" />
      </template>

      <template #body-enabled="{ data }">
        <Tag
          :value="data.enabled ? $t('general.label.enabled') : $t('general.label.disabled')"
          :severity="data.enabled ? 'success' : 'secondary'"
        />
      </template>

      <template #body-createdAt="{ data }">
        {{ locFullDateTime(data.deletedAt || data.updatedAt || data.createdAt) }}
      </template>

      <template #filter>
        <Button
          icon="pi pi-filter"
          severity="secondary"
          size="small"
          text
          @click="toggleFilterMenu"
        />
      </template>
    </CResourceList>

    <Popover ref="filterMenu">
      <div class="flex flex-col gap-4 p-2 w-64">
        <div class="flex flex-col gap-2">
          <span class="font-medium text-sm text-primary">
            {{ $t('system.apigw.list.filterForm.deleted.label') }}
          </span>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del0" value="0" />
            <label for="del0" class="text-sm cursor-pointer">
              {{ $t('system.apigw.list.filterForm.excluded.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del1" value="1" />
            <label for="del1" class="text-sm cursor-pointer">
              {{ $t('system.apigw.list.filterForm.inclusive.label') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <RadioButton v-model="filter.deleted" inputId="del2" value="2" />
            <label for="del2" class="text-sm cursor-pointer">
              {{ $t('system.apigw.list.filterForm.exclusive.label') }}
            </label>
          </div>
        </div>
      </div>
    </Popover>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  components,
  filters,
  useConfirmDelete,
  useResourceList,
  useRBACStore,
  usePermissions,
} from '@planetcrust/human-vue'

const { CResourceList } = components
const { locFullDateTime } = filters

const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const rbac = useRBACStore()
const canGrant = computed(() => rbac.can('system/', 'grant'))
const { open: openPermissions } = usePermissions()

const resourceListRef = ref()
const filterMenu = ref()

function toggleFilterMenu(event) {
  filterMenu.value.toggle(event)
}

const fields = [
  {
    key: 'endpoint',
    sortable: true,
    header: t('system.apigw.list.columns.endpoint'),
  },
  {
    key: 'method',
    sortable: true,
    header: t('system.apigw.list.columns.method'),
  },
  {
    key: 'enabled',
    sortable: false,
    header: t('system.apigw.list.columns.enabled'),
  },
  {
    key: 'group',
    sortable: true,
    header: t('system.apigw.list.columns.group'),
  },
  {
    key: 'createdAt',
    sortable: true,
    header: t('system.apigw.list.columns.createdAt'),
    class: 'text-right',
    pt: { columnHeaderContent: 'justify-end' },
  },
]

const { items, loading, filter, sorting, pagination, handleSort, handlePageChange, filterList } =
  useResourceList(params => $SystemAPI.apigwRouteListCancellable({ ...params }), {
    filter: { query: '', deleted: '0' },
    sorting: { sortBy: 'createdAt', sortDesc: true },
    pagination: { limit: 50 },
  })

function getActionsMenuItems(item) {
  const items = []

  if (item.canGrant || canGrant.value) {
    items.push({
      label: t('general.label.permissions'),
      icon: 'pi pi-lock',
      command: () => {
        resourceListRef.value?.hideActionsMenu?.()
        openPermissions({
          resource: `corteza::system:apigw-route/${item.routeID}`,
          title: item.endpoint || item.routeID,
        })
      },
    })
  }

  if (item.canDeleteApigwRoute) {
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => onConfirmDelete(item),
    })
  }

  return items
}

function onConfirmDelete(item) {
  confirmDelete({
    message: t('general.confirm.delete'),
    header: item.endpoint || item.routeID,
    onConfirm: () => handleDelete(item),
  })
}

async function handleDelete(item) {
  resourceListRef.value.hideActionsMenu()
  try {
    await $SystemAPI.apigwRouteDelete({ routeID: item.routeID })
    $toast.toastSuccess(t('notification.gateway.delete.success'))
    filterList()
  } catch (e) {
    $toast.toastErrorHandler(t('notification.gateway.delete.error'))(e)
  }
}

// --- Settings ---
const savingSettings = ref(false)
const settingsItems = ref([])
const profilerSetting = ref('disabled')
const proxyFollowRedirects = ref(false)

const profilerOptions = [
  { value: 'disabled', label: 'Disabled' },
  { value: 'filter', label: 'Enabled as filter' },
  { value: 'global', label: 'Enabled for all routes' },
]

function onSettingsChanged() {
  // Just marks dirty, save button handles persistence
}

async function fetchSettings() {
  try {
    const settings = await $SystemAPI.settingsList()
    settingsItems.value = settings || []

    const map = settings.reduce((m, { name, value }) => {
      m[name] = value
      return m
    }, {})

    // Derive profiler setting
    if (map['apigw.profiler.enabled']) {
      profilerSetting.value = map['apigw.profiler.global'] ? 'global' : 'filter'
    } else {
      profilerSetting.value = 'disabled'
    }

    proxyFollowRedirects.value = !!map['apigw.proxy.follow-redirects']
  } catch (e) {
    console.error('Failed to fetch settings:', e)
  }
}

async function saveSettings() {
  savingSettings.value = true
  try {
    const values = [
      {
        name: 'apigw.profiler.enabled',
        value: ['filter', 'global'].includes(profilerSetting.value),
      },
      { name: 'apigw.profiler.global', value: profilerSetting.value === 'global' },
      { name: 'apigw.proxy.follow-redirects', value: proxyFollowRedirects.value },
    ]

    await $SystemAPI.settingsUpdate({ values })
    $toast.toastSuccess(t('notification.settings.system.apigw.success'))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.system.apigw.error'))(e)
  } finally {
    savingSettings.value = false
  }
}

onMounted(() => fetchSettings())
</script>
