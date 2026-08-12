<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.llmProviders.list.title') }}</span>
  </Teleport>

  <CViewContainer>
    <CResourceList
      ref="resourceListRef"
      primary-key="llmProviderID"
      :fields="fields"
      :items="items"
      :loading="loading"
      hide-search
      :action-items="getActionsMenuItems"
      :translations="{
        resourceSingle: $t('system.llmProviders.list.new'),
        resourcePlural: $t('system.llmProviders.list.title'),
      }"
      clickable
      class="h-full"
      @row-click="
        ({ data }) =>
          $router.push({
            name: 'system.llmProviders.edit',
            params: { llmProviderID: data.llmProviderID },
          })
      "
    >
      <template #header>
        <div class="flex gap-2">
          <Button
            :label="$t('system.llmProviders.list.new')"
            icon="pi pi-plus"
            size="small"
            @click="$router.push({ name: 'system.llmProviders.create' })"
          />
        </div>
      </template>

      <template #body-short="{ data }">
        {{ data.meta?.short || '—' }}
      </template>

      <template #body-provider="{ data }">
        <span v-if="data.provider">
          {{
            data.provider === 'openai'
              ? 'OpenAI'
              : data.provider === 'anthropic'
                ? 'Anthropic'
                : data.provider === 'mistral'
                  ? 'Mistral'
                  : 'Other'
          }}
        </span>
        <span v-else>—</span>
      </template>

      <template #body-status="{ data }">
        <Tag
          v-if="data.status"
          :value="$t(`system.llmProviders.editor.info.statusOptions.${data.status}`)"
          :severity="
            data.status === 'active'
              ? 'success'
              : data.status === 'unauthorized'
                ? 'danger'
                : 'warn'
          "
          :icon="
            data.status === 'active'
              ? 'pi pi-check-circle'
              : data.status === 'unauthorized'
                ? 'pi pi-times-circle'
                : 'pi pi-pause-circle'
          "
        />
        <span v-else>—</span>
      </template>

      <template #body-createdAt="{ data }">
        {{ locFullDateTime(data.createdAt) }}
      </template>
    </CResourceList>
  </CViewContainer>
</template>

<script setup>
import { inject, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { components, filters, useConfirmDelete } from '@planetcrust/human-vue'

const { CResourceList, CViewContainer } = components
const { locFullDateTime } = filters

const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const resourceListRef = ref()
const loading = ref(false)
const items = ref([])

const fields = [
  {
    key: 'short',
    sortable: false,
    header: t('system.llmProviders.list.columns.meta.short'),
  },
  {
    key: 'handle',
    sortable: false,
    header: t('system.llmProviders.list.columns.handle'),
  },
  {
    key: 'provider',
    sortable: false,
    header: t('system.llmProviders.list.columns.provider'),
  },
  {
    key: 'status',
    sortable: false,
    header: t('system.llmProviders.list.columns.status'),
  },
  {
    key: 'createdAt',
    sortable: false,
    header: t('system.llmProviders.list.columns.createdAt'),
    class: 'text-right',
    pt: { columnHeaderContent: 'justify-end' },
  },
]

async function fetchList() {
  loading.value = true
  try {
    const result = await $SystemAPI.llmProviderList({})
    items.value = result?.set || result || []
  } catch (e) {
    $toast.toastErrorHandler(t('notification.llmProvider.fetch.error'))(e)
    items.value = []
  } finally {
    loading.value = false
  }
}

function getActionsMenuItems(item) {
  const menuItems = []

  if (item.canDeleteLlmProvider) {
    menuItems.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => onConfirmDelete(item),
    })
  }

  return menuItems
}

function onConfirmDelete(item) {
  confirmDelete({
    message: t('general.confirm.delete'),
    header: item.meta?.short || item.handle || item.llmProviderID,
    onConfirm: () => handleDelete(item),
  })
}

async function handleDelete(item) {
  resourceListRef.value.hideActionsMenu()
  try {
    await $SystemAPI.llmProviderDelete({ llmProviderID: item.llmProviderID })
    $toast.toastSuccess(t('notification.llmProvider.delete.success'))
    fetchList()
  } catch (e) {
    $toast.toastErrorHandler(t('notification.llmProvider.delete.error'))(e)
  }
}

onMounted(() => fetchList())
</script>
