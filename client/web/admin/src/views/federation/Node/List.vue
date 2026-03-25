<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('federation.nodes.list.title') }}</span>
  </Teleport>

  <div class="flex flex-col h-full">
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0">
      <CResourceList
        primary-key="nodeID"
        :fields="fields"
        :items="items"
        :filter="filter"
        :sorting="sorting"
        :pagination="pagination"
        :loading="loading"
        :translations="{
          searchPlaceholder: $t('federation.nodes.list.filter.query.placeholder'),
          showingPagination: 'general.resourceList.pagination.showing',
          singlePluralPagination: 'general.resourceList.pagination.single',
          prevPagination: $t('general.resourceList.pagination.prev'),
          nextPagination: $t('general.resourceList.pagination.next'),
          recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
          resourceSingle: $t('federation.nodes.list.new'),
          resourcePlural: $t('federation.nodes.list.title'),
        }"
        clickable
        @update:filter="Object.assign(filter, $event)"
        @sort="handleSort"
        @row-click="({ data }) => $router.push({ name: 'federation.nodes.edit', params: { nodeID: data.nodeID } })"
        @page-change="handlePageChange"
      >
        <template #header>
          <div class="flex gap-2">
            <Button
              :label="$t('federation.nodes.list.new')"
              icon="pi pi-plus"
              size="small"
              @click="$router.push({ name: 'federation.nodes.create' })"
            />
            <Button
              :label="$t('federation.nodes.list.pair')"
              icon="pi pi-link"
              size="small"
              severity="secondary"
              @click="pairDialogVisible = true"
            />
          </div>
        </template>

        <template #body-status="{ data }">
          <Tag
            :value="data.status || 'unknown'"
            :severity="data.status === 'paired' ? 'success' : data.status === 'pair_requested' ? 'warn' : 'secondary'"
          />
        </template>

        <template #body-actions="{ data }">
          <Button
            v-if="data.status === 'pair_requested'"
            :label="$t('federation.nodes.pair.confirm')"
            icon="pi pi-check"
            size="small"
            severity="warn"
            text
            @click.stop="handleConfirmPending(data)"
          />
        </template>
      </CResourceList>
    </div>
  </div>

  <!-- Pair Dialog -->
  <Dialog
    v-model:visible="pairDialogVisible"
    modal
    :header="$t('federation.nodes.pair.title')"
    :style="{ width: '500px' }"
  >
    <div class="flex flex-col gap-4 p-2">
      <p class="text-sm text-muted-color">{{ $t('federation.nodes.pair.description') }}</p>
      <InputText
        v-model="pairURL"
        placeholder="https://..."
        class="w-full"
      />
      <div class="flex justify-end gap-2">
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          @click="pairDialogVisible = false"
        />
        <Button
          :label="$t('federation.nodes.pair.button')"
          icon="pi pi-link"
          :disabled="!pairURL"
          :loading="pairing"
          @click="handlePair"
        />
      </div>
    </div>
  </Dialog>
</template>

<script setup>
import { components, useResourceList } from '@cortezaproject/corteza-vue-next'
import { inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { CResourceList } = components

const { t } = useI18n()
const $toast = inject('$toast')
const $FederationAPI = inject('$FederationAPI')

const pairDialogVisible = ref(false)
const pairURL = ref('')
const pairing = ref(false)

const fields = [
  { key: 'name', sortable: true, header: t('federation.nodes.list.columns.name') },
  { key: 'status', sortable: true, header: t('federation.nodes.list.columns.status') },
  { key: 'baseURL', header: t('federation.nodes.list.columns.baseURL') },
  {
    key: 'createdAt',
    sortable: true,
    header: t('federation.nodes.list.columns.createdAt'),
    formatter: v => (v ? new Date(v).toLocaleDateString() : ''),
  },
  { key: 'actions', header: '', class: 'text-right w-24' },
]

const {
  items,
  loading,
  filter,
  sorting,
  pagination,
  handleSort,
  handlePageChange,
  filterList,
} = useResourceList(
  params => $FederationAPI.nodeSearchCancellable(params),
  {
    filter: { query: '' },
    sorting: { sortBy: 'createdAt', sortDesc: true },
    pagination: { limit: 20 },
  },
)

async function handlePair() {
  if (!pairURL.value) return

  pairing.value = true
  try {
    const node = await $FederationAPI.nodeCreate({ pairingURI: pairURL.value })
    await $FederationAPI.nodePair(node)
    $toast.toastSuccess(t('federation.nodes.pair.success'))
    pairDialogVisible.value = false
    pairURL.value = ''
    filterList()
  } catch (e) {
    $toast.toastErrorHandler(t('federation.nodes.pair.error'))(e)
  } finally {
    pairing.value = false
  }
}

async function handleConfirmPending(node) {
  try {
    await $FederationAPI.nodeHandshakeConfirm({ nodeID: node.nodeID })
    $toast.toastSuccess(t('federation.nodes.pair.confirmSuccess'))
    filterList()
  } catch (e) {
    $toast.toastErrorHandler(t('federation.nodes.pair.confirmError'))(e)
  }
}
</script>
