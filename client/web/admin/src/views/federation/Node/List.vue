<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('federation.nodes.list.title') }}</span>
  </Teleport>

  <div class="flex flex-col h-full">
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0">
      <CResourceList
        ref="resourceListRef"
        primary-key="nodeID"
        :fields="fields"
        :items="items"
        :filter="filter"
        :sorting="sorting"
        :pagination="pagination"
        :action-items="getActionsMenuItems"
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
            <CPermissionsButton
              v-if="canGrant"
              v-tooltip.bottom="$t('general.label.permissions')"
              resource="corteza::federation:node/*"
            />
          </div>
        </template>

        <template #body-status="{ data }">
          <Tag
            :value="data.status || 'unknown'"
            :severity="data.status === 'paired' ? 'success' : data.status === 'pair_requested' ? 'warn' : 'secondary'"
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
import {
  components,
  useResourceList,
  useRBACStore,
  usePermissions,
} from '@planetcrust/human-vue'
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { CResourceList } = components

const { t } = useI18n()
const $toast = inject('$toast')
const $FederationAPI = inject('$FederationAPI')

const rbac = useRBACStore()
const canGrant = computed(() => rbac.can('federation/', 'grant'))
const { open: openPermissions } = usePermissions()

const resourceListRef = ref()
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

function getActionsMenuItems(node) {
  const actions = []

  if (node.canGrant || canGrant.value) {
    actions.push({
      label: t('general.label.permissions'),
      icon: 'pi pi-lock',
      command: () => {
        resourceListRef.value?.hideActionsMenu?.()
        openPermissions({
          resource: `corteza::federation:node/${node.nodeID}`,
          title: node.name || node.nodeID,
        })
      },
    })
  }

  if (node.status === 'pair_requested') {
    actions.push({
      label: t('federation.nodes.pair.confirm'),
      icon: 'pi pi-check',
      class: 'text-orange-500',
      command: () => handleConfirmPending(node),
    })
  }

  return actions
}

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
