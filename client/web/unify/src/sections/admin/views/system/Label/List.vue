<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.labels.list.title') }}</span>
  </Teleport>

  <CViewContainer>
    <CResourceList
      ref="resourceListRef"
      primary-key="_rowKey"
      :fields="fields"
      :items="items"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :translations="{
        searchPlaceholder: $t('system.labels.list.filterForm.query.placeholder'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('system.labels.list.resourceSingle'),
        resourcePlural: $t('system.labels.list.title'),
      }"
      class="h-full"
      clickable
      @update:filter="Object.assign(filter, $event)"
      @sort="handleSort"
      @page-change="handlePageChange"
      @row-click="onRowClick"
    >
      <template #header>
        <Button
          :label="$t('system.labels.list.createLabel')"
          icon="pi pi-plus"
          size="small"
          @click="showCreateDialog = true"
        />
      </template>

      <template #body-resourceCount="{ data }">
        <Tag :value="String(data.resourceCount || 0)" severity="secondary" rounded />
      </template>
    </CResourceList>
  </CViewContainer>

  <Dialog
    v-model:visible="showCreateDialog"
    :header="$t('system.labels.list.createLabel')"
    modal
    :style="{ width: '28rem' }"
    @hide="newLabelName = ''"
  >
    <Form
      ref="createFormRef"
      :resolver="createResolver"
      :initialValues="{ name: newLabelName }"
      @submit="submitCreate"
    >
      <CFormGroup
        name="name"
        input-id="new-label-name"
        :label="$t('system.labels.editor.info.name')"
        class="pt-2"
      >
        <InputText
          id="new-label-name"
          name="name"
          v-model="newLabelName"
          :placeholder="$t('system.labels.create.namePlaceholder')"
          autofocus
          fluid
        />
      </CFormGroup>
    </Form>

    <template #footer>
      <div class="flex justify-end gap-2">
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          text
          size="small"
          @click="showCreateDialog = false"
        />
        <Button
          :label="$t('general.label.save')"
          size="small"
          @click="createFormRef?.submit?.()"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import {
  components,
  useResourceList,
} from '@planetcrust/human-vue'

const { CResourceList, CViewContainer } = components

const { t } = useI18n()
const router = useRouter()

const $SystemAPI = inject('$SystemAPI')

const resourceListRef = ref()
const showCreateDialog = ref(false)
const newLabelName = ref('')
const createFormRef = ref(null)

const fields = [
  {
    key: 'name',
    sortable: false,
    header: t('system.labels.list.columns.name'),
  },
  {
    key: 'resourceCount',
    sortable: false,
    header: t('system.labels.list.columns.resources'),
  },
]

function labelListCancellable(params) {
  const { response, cancel } = $SystemAPI.labelListCancellable({
    ...params,
    name: params.query || undefined,
  })

  return {
    response: async () => {
      const result = await response()
      const set = Array.isArray(result) ? result : result.set || []

      const enriched = set.map((item, idx) => ({
        ...item,
        _rowKey: `${item.kind || ''}_${item.name || ''}_${idx}`,
      }))

      if (Array.isArray(result)) {
        return enriched
      }

      return {
        ...result,
        set: enriched,
      }
    },
    cancel,
  }
}

const { items, loading, filter, sorting, pagination, handleSort, handlePageChange } =
  useResourceList(params => labelListCancellable(params), {
    filter: { query: '' },
    sorting: {},
    pagination: { limit: 100 },
  })

function onRowClick({ data }) {
  if (data?.name) {
    router.push({
      name: 'system.labels.edit',
      params: { labelID: encodeURIComponent(data.name) },
    })
  }
}

const HANDLE_PATTERN = /^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$|^[A-Za-z]$/

function createResolver({ values }) {
  const errors = {}
  const name = (values?.name || '').trim()
  if (!name) {
    errors.name = [{ message: t('general.label.required') }]
  } else if (!HANDLE_PATTERN.test(name)) {
    errors.name = [{ message: t('system.labels.create.invalid-handle-characters') }]
  }
  return { errors }
}

function submitCreate({ valid, values }) {
  if (!valid) return
  const name = (values?.name || newLabelName.value).trim()
  if (!name) return
  showCreateDialog.value = false
  router.push({
    name: 'system.labels.edit',
    params: { labelID: encodeURIComponent(name) },
  })
}
</script>
