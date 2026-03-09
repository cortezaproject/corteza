<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ pageTitle }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <Form
    v-else-if="route_"
    v-slot="$form"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4">
      <Card
        :pt="{
          body: { class: 'p-0 flex flex-col h-full min-h-0' },
          content: { class: 'p-0 flex flex-col h-full min-h-0' },
        }"
        class="overflow-hidden flex-1 min-h-0 flex flex-col"
      >
        <template #content>
          <Tabs v-model:value="activeTab" class="flex flex-col h-full min-h-0">
            <TabList class="rounded-t-lg shrink-0">
              <Tab value="route">{{ $t('system.apigw.editor.tabs.route', 'Route') }}</Tab>
              <Tab v-if="isEdit" value="filters">
                {{ $t('system.apigw.editor.tabs.filters', 'Filters') }}
              </Tab>
            </TabList>

            <TabPanels class="flex-1 overflow-y-auto min-h-0">
              <TabPanel value="route">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <FormField name="endpoint" class="flex flex-col gap-2 md:col-span-2">
                    <label for="endpoint" class="font-medium text-primary">
                      {{ $t('system.apigw.editor.info.endpoint', 'Endpoint') }} *
                    </label>
                    <InputText
                      id="endpoint"
                      name="endpoint"
                      v-model="route_.endpoint"
                      placeholder="/api/v1/..."
                    />
                    <Message
                      v-if="$form.endpoint?.invalid"
                      severity="error"
                      size="small"
                      variant="simple"
                    >
                      {{ $form.endpoint.error?.message }}
                    </Message>
                  </FormField>

                  <div class="flex flex-col gap-2">
                    <label for="method" class="font-medium text-primary">
                      {{ $t('system.apigw.editor.info.method', 'Method') }}
                    </label>
                    <Select
                      id="method"
                      v-model="route_.method"
                      :options="methodOptions"
                      option-label="label"
                      option-value="value"
                    />
                  </div>

                  <div class="flex flex-col gap-2">
                    <label for="group" class="font-medium text-primary">
                      {{ $t('system.apigw.editor.info.group', 'Group') }}
                    </label>
                    <InputText id="group" v-model="route_.group" />
                  </div>

                  <div class="flex flex-col gap-2 md:col-span-2">
                    <label for="description" class="font-medium text-primary">
                      {{ $t('system.apigw.editor.info.description', 'Description') }}
                    </label>
                    <Textarea id="description" v-model="route_.meta.description" rows="2" />
                  </div>

                  <div class="flex items-center gap-3">
                    <ToggleSwitch id="enabled" v-model="route_.enabled" />
                    <label for="enabled" class="font-medium text-primary cursor-pointer">
                      {{ $t('system.apigw.editor.info.enabled', 'Enabled') }}
                    </label>
                  </div>

                  <div class="flex items-center gap-3">
                    <ToggleSwitch id="async" v-model="route_.meta.async" />
                    <label for="async" class="font-medium text-primary cursor-pointer">
                      {{ $t('system.apigw.editor.info.async', 'Async') }}
                    </label>
                  </div>
                </div>
              </TabPanel>

              <TabPanel v-if="isEdit" value="filters">
                <CResourceList
                  primary-key="filterID"
                  :fields="filterFields"
                  :items="filterItems"
                  :filter="filterFilter"
                  :sorting="filterSorting"
                  :pagination="filterPagination"
                  :loading="filtersLoading"
                  :translations="{
                    searchPlaceholder: $t(
                      'system.apigw.editor.filters.searchPlaceholder',
                      'Search filters...',
                    ),
                    showingPagination: 'general.resourceList.pagination.showing',
                    singlePluralPagination: 'general.resourceList.pagination.single',
                    prevPagination: $t('general.resourceList.pagination.prev'),
                    nextPagination: $t('general.resourceList.pagination.next'),
                    recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
                    resourceSingle: $t('system.apigw.editor.filters.resourceSingle', 'filter'),
                    resourcePlural: $t('system.apigw.editor.filters.resourcePlural', 'filters'),
                  }"
                  clickable
                  @update:filter="Object.assign(filterFilter, $event)"
                  @sort="handleFilterSort"
                  @row-click="({ data }) => openFilterDialog(data)"
                  @page-change="handleFilterPageChange"
                >
                  <template #header>
                    <Button
                      :label="$t('system.apigw.editor.filters.add', 'Add Filter')"
                      icon="pi pi-plus"
                      size="small"
                      @click="openFilterDialog(null)"
                    />
                  </template>

                  <template #body-enabled="{ data }">
                    <Tag
                      :value="
                        data.enabled
                          ? $t('general.label.enabled', 'Enabled')
                          : $t('general.label.disabled', 'Disabled')
                      "
                      :severity="data.enabled ? 'success' : 'secondary'"
                    />
                  </template>

                  <template #body-actions="{ data }">
                    <Button
                      icon="pi pi-ellipsis-v"
                      text
                      severity="secondary"
                      size="small"
                      class="row-action-btn w-full"
                      @click.stop="toggleFilterActionsMenu($event, data)"
                    />
                  </template>
                </CResourceList>

                <TieredMenu ref="filterActionsMenu" :model="filterActionsMenuItems" popup>
                  <template #item="{ item, props }">
                    <a v-ripple v-bind="props.action" :class="item.class">
                      <span :class="item.icon" />
                      <span class="ml-2">{{ item.label }}</span>
                    </a>
                  </template>
                </TieredMenu>
              </TabPanel>
            </TabPanels>
          </Tabs>
        </template>
      </Card>
    </div>

    <!-- Filter Dialog -->
    <Dialog
      v-model:visible="filterDialogVisible"
      modal
      :header="
        activeFilter?.filterID
          ? $t('system.apigw.editor.filters.edit', 'Edit Filter')
          : $t('system.apigw.editor.filters.create', 'New Filter')
      "
      :style="{ width: '50vw' }"
      :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
      :pt="{ content: { class: '!overflow-hidden' } }"
    >
      <Form
        v-if="filterDialogVisible && activeFilter"
        v-slot="$filterForm"
        :resolver="filterResolver"
        :initialValues="filterInitialValues"
        @submit="handleFilterSubmit"
        class="flex flex-col gap-4 p-1"
      >
        <FormField name="kind" class="flex flex-col gap-2">
          <label for="filterKind" class="font-medium text-primary">
            {{ $t('system.apigw.editor.filters.kind', 'Kind') }} *
          </label>
          <InputText id="filterKind" name="kind" v-model="activeFilter.kind" />
          <Message v-if="$filterForm.kind?.invalid" severity="error" size="small" variant="simple">
            {{ $filterForm.kind.error?.message }}
          </Message>
        </FormField>

        <div class="grid grid-cols-2 gap-4">
          <div class="flex flex-col gap-2">
            <label for="filterWeight" class="font-medium text-primary">
              {{ $t('system.apigw.editor.filters.weight', 'Weight') }}
            </label>
            <InputNumber id="filterWeight" v-model="activeFilter.weight" :min="0" />
          </div>

          <div class="flex items-center gap-3 mt-6">
            <ToggleSwitch id="filterEnabled" v-model="activeFilter.enabled" />
            <label for="filterEnabled" class="font-medium text-primary cursor-pointer">
              {{ $t('general.label.enabled', 'Enabled') }}
            </label>
          </div>
        </div>

        <div class="flex flex-col gap-2">
          <label for="filterParams" class="font-medium text-primary">
            {{ $t('system.apigw.editor.filters.params', 'Params (JSON)') }}
          </label>
          <Textarea
            id="filterParams"
            name="params"
            v-model="rawFilterParams"
            rows="6"
            autoResize
            class="font-mono text-sm"
            @change="parseFilterParams"
          />
          <Message
            v-if="$filterForm.params?.invalid"
            severity="error"
            size="small"
            variant="simple"
          >
            {{ $filterForm.params.error?.message }}
          </Message>
        </div>

        <div class="border-t border-surface pt-3 flex justify-end gap-2">
          <Button
            :label="$t('general.label.cancel')"
            severity="secondary"
            size="small"
            outlined
            @click="filterDialogVisible = false"
          />
          <Button
            type="submit"
            :label="$t('general.label.save')"
            size="small"
            :loading="savingFilter"
          />
        </div>
      </Form>
    </Dialog>

    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="p-3 flex items-center justify-between">
        <Button
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="$router.push({ name: 'system.apiGateway' })"
        />
        <div class="flex gap-2">
          <CInputDelete
            v-if="isEdit && route_.canDeleteApigwRoute"
            :label="$t('system.apigw.editor.delete', 'Delete')"
            :message="$t('general.confirm.delete')"
            :header="route_.endpoint || route_.routeID"
            :disabled="deleting"
            @confirm="handleDelete"
          />
          <Button
            type="submit"
            :label="$t('general.label.save')"
            icon="pi pi-save"
            :loading="saving"
          />
        </div>
      </div>
    </div>
  </Form>
</template>

<script setup>
import { computed, inject, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { components, useConfirmDelete, useResourceList } from '@cortezaproject/corteza-vue-next'

const { CInputDelete, CResourceList } = components

const vueRoute = useRoute()
const router = useRouter()
const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const route_ = ref(null)
const activeTab = ref('route')

// Filter dialog state
const filterDialogVisible = ref(false)
const savingFilter = ref(false)
const activeFilter = ref(null)
const rawFilterParams = ref('{}')
const filterActionsMenu = ref()
const filterActionsMenuItems = ref([])

const methodOptions = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
  { label: 'PATCH', value: 'PATCH' },
]

const filterFields = [
  { key: 'kind', sortable: true, header: t('system.apigw.editor.filters.kind', 'Kind') },
  { key: 'weight', sortable: true, header: t('system.apigw.editor.filters.weight', 'Weight') },
  { key: 'enabled', sortable: false, header: t('system.apigw.editor.filters.enabled', 'Status') },
  {
    key: 'actions',
    class: 'text-right w-12',
    header: '',
    frozen: true,
    alignFrozen: 'right',
    pt: {
      headerCell: { class: 'border-l-0' },
      bodyCell: { class: 'p-0 border-l-0' },
    },
  },
]

const isEdit = computed(() => !!vueRoute.params.routeID)

const pageTitle = computed(() =>
  isEdit.value
    ? t('system.apigw.editor.title.edit', 'Edit Route')
    : t('system.apigw.editor.title.create', 'New Route'),
)

const initialValues = computed(() => ({
  endpoint: route_.value?.endpoint || '',
}))

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.endpoint || values.endpoint.trim().length === 0) {
    errors.endpoint = [{ message: t('general.label.required') }]
  }

  return { errors }
})

const filterInitialValues = computed(() => ({
  kind: activeFilter.value?.kind || '',
  params: rawFilterParams.value,
}))

const filterResolver = ref(({ values }) => {
  const errors = {}

  if (!values.kind || values.kind.trim().length === 0) {
    errors.kind = [{ message: t('general.label.required') }]
  }

  try {
    if (values.params) JSON.parse(values.params)
  } catch {
    errors.params = [
      { message: t('system.connections.editor.configurations.invalidJSON', 'Invalid JSON') },
    ]
  }

  return { errors }
})

// Filters resource list
const {
  items: filterItems,
  loading: filtersLoading,
  filter: filterFilter,
  sorting: filterSorting,
  pagination: filterPagination,
  handleSort: handleFilterSort,
  handlePageChange: handleFilterPageChange,
  filterList: reloadFilters,
} = useResourceList(
  params => {
    if (!isEdit.value || !vueRoute.params.routeID) return Promise.resolve({ set: [] })
    return $SystemAPI.apigwFilterListCancellable({ ...params, routeID: vueRoute.params.routeID })
  },
  {
    filter: { query: '' },
    sorting: { sortBy: 'weight', sortDesc: false },
    pagination: { limit: 50 },
  },
)

function newRoute() {
  return {
    endpoint: '',
    method: 'GET',
    enabled: true,
    group: '',
    meta: { description: '', async: false },
  }
}

function parseFilterParams() {
  try {
    const parsed = JSON.parse(rawFilterParams.value || '{}')
    if (parsed && activeFilter.value) {
      activeFilter.value.params = parsed
    }
  } catch {
    // caught by resolver
  }
}

function openFilterDialog(filter) {
  if (filter) {
    activeFilter.value = { ...filter }
    rawFilterParams.value = JSON.stringify(filter.params || {}, null, 2)
  } else {
    activeFilter.value = { kind: '', weight: 1, enabled: true, params: {} }
    rawFilterParams.value = '{\n  \n}'
  }
  filterDialogVisible.value = true
}

function toggleFilterActionsMenu(event, filter) {
  filterActionsMenuItems.value = [
    {
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => onConfirmFilterDelete(filter),
    },
  ]
  filterActionsMenu.value.toggle(event)
}

function onConfirmFilterDelete(filter) {
  confirmDelete({
    message: t('general.confirm.delete'),
    header: filter.kind || filter.filterID,
    onConfirm: () => handleFilterDelete(filter),
  })
}

async function handleFilterDelete(filter) {
  try {
    await $SystemAPI.apigwFilterDelete({
      routeID: vueRoute.params.routeID,
      filterID: filter.filterID,
    })
    $toast.toastSuccess(t('notification.gateway.filter.delete.success', 'Filter deleted'))
    reloadFilters()
  } catch (e) {
    $toast.toastErrorHandler(
      t('notification.gateway.filter.delete.error', 'Failed to delete filter'),
    )(e)
  }
}

async function handleFilterSubmit({ valid }) {
  if (!valid) return

  savingFilter.value = true
  try {
    parseFilterParams()

    const payload = {
      routeID: vueRoute.params.routeID,
      kind: activeFilter.value.kind,
      weight: activeFilter.value.weight,
      enabled: activeFilter.value.enabled,
      params: activeFilter.value.params || {},
    }

    if (activeFilter.value.filterID) {
      payload.filterID = activeFilter.value.filterID
      await $SystemAPI.apigwFilterUpdate(payload)
      $toast.toastSuccess(t('notification.gateway.filter.update.success', 'Filter updated'))
    } else {
      await $SystemAPI.apigwFilterCreate(payload)
      $toast.toastSuccess(t('notification.gateway.filter.create.success', 'Filter created'))
    }

    filterDialogVisible.value = false
    reloadFilters()
  } catch (e) {
    $toast.toastErrorHandler(t('notification.gateway.filter.save.error', 'Failed to save filter'))(
      e,
    )
  } finally {
    savingFilter.value = false
  }
}

async function loadRoute() {
  const routeID = vueRoute.params.routeID
  if (!routeID) {
    route_.value = newRoute()
    return
  }

  loading.value = true
  try {
    const raw = await $SystemAPI.apigwRouteRead({ routeID })
    route_.value = {
      ...raw,
      meta: { description: raw.meta?.description || '', async: raw.meta?.async || false },
    }
  } catch (e) {
    $toast.toastErrorHandler(t('notification.gateway.fetch.error', 'Failed to load route'))(e)
    router.push({ name: 'system.apiGateway' })
  } finally {
    loading.value = false
  }
}

async function handleSubmit({ valid }) {
  if (!valid) return

  saving.value = true
  try {
    const payload = {
      endpoint: route_.value.endpoint,
      method: route_.value.method,
      enabled: route_.value.enabled,
      group: route_.value.group,
      meta: route_.value.meta,
    }

    if (isEdit.value) {
      payload.routeID = route_.value.routeID
      const raw = await $SystemAPI.apigwRouteUpdate(payload)
      route_.value = {
        ...raw,
        meta: { description: raw.meta?.description || '', async: raw.meta?.async || false },
      }
      $toast.toastSuccess(t('notification.gateway.update.success', 'Route updated'))
    } else {
      const created = await $SystemAPI.apigwRouteCreate(payload)
      $toast.toastSuccess(t('notification.gateway.create.success', 'Route created'))
      router.push({ name: 'system.apiGateway.edit', params: { routeID: created.routeID } })
    }
  } catch (e) {
    $toast.toastErrorHandler(
      t(`notification.gateway.${isEdit.value ? 'update' : 'create'}.error`, 'Failed to save route'),
    )(e)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await $SystemAPI.apigwRouteDelete({ routeID: route_.value.routeID })
    $toast.toastSuccess(t('notification.gateway.delete.success', 'Route deleted'))
    router.push({ name: 'system.apiGateway' })
  } catch (e) {
    $toast.toastErrorHandler(t('notification.gateway.delete.error', 'Failed to delete route'))(e)
  } finally {
    deleting.value = false
  }
}

onMounted(() => loadRoute())
watch(
  () => vueRoute.params.routeID,
  () => loadRoute(),
)
</script>
