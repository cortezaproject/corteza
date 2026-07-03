<template>
  <div class="h-full overflow-auto p-4">
    <div class="flex flex-col gap-2">
      <div>
        <Button
          v-if="!disabled"
          icon="pi pi-plus"
          :label="$t('project.pages.add')"
          severity="secondary"
          size="small"
          @click="createResource?.('page')"
        />
      </div>

      <!-- Record (module detail) pages and standalone pages share one list. The
           built-in remove is disabled (hide-remove) so we can render a delete
           button only for standalone pages — detail pages go with their module. -->
      <div class="mt-1">
        <CFormItemList
          :items="pages"
          item-key="id"
          reveal-on-hover
          hide-remove
          :empty-message="$t('project.pages.empty')"
          @select="item => inspectResource?.('page', item.id)"
        >
          <template #default="{ item }">
            <CFormItemContent :title="item.name">
              <template v-if="item.isRecordPage" #subtitle>
                <div class="text-xs text-muted-color truncate flex items-center gap-1">
                  <i class="pi pi-database text-[10px]" />
                  {{ $t('project.pages.recordSubtitle') }}
                </div>
              </template>
            </CFormItemContent>
          </template>

          <template #actions="{ item }">
            <Tag
              :value="item.visible ? $t('project.pages.visible') : $t('project.pages.hidden')"
              :severity="item.visible ? 'success' : 'secondary'"
              class="!text-xs shrink-0 me-2"
            />
          </template>

          <template #hover-actions="{ item }">
            <CRouterLinkButton
              :to="{ name: 'admin.pages.builder', params: { slug: project.namespaceID, pageID: item.id } }"
              icon="pi pi-external-link"
              severity="secondary"
              text
              size="small"
              :aria-label="$t('project.pages.openBuilder')"
              :title="$t('project.pages.openBuilder')"
              @click.stop
            />
            <!-- Detail pages are tied to their module; only standalone pages can be
                 removed from here. -->
            <Button
              v-if="!disabled && !item.isRecordPage"
              icon="pi pi-trash"
              severity="danger"
              text
              size="small"
              :aria-label="$t('project.pages.remove')"
              :title="$t('project.pages.remove')"
              @click.stop="onRemove(item)"
            />
          </template>
        </CFormItemList>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useProjectsStore } from '@/sections/project/stores/projects'
import { components, useConfirmDelete } from '@planetcrust/human-vue'
import { computed, inject, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { CRouterLinkButton } = components

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()
const { t } = useI18n()
const $toast = inject('$toast')
const { confirmDelete } = useConfirmDelete()

// Detail/create dialogs are mounted once in the wizard; open them via injection.
const inspectResource = inject('inspectResource', null)
const createResource = inject('createResource', null)

// Record (module detail) pages first, then standalone — both in one list.
const pages = computed(() => {
  const all = store.pagesFor(props.project.id)
  return [...all.filter(p => p.isRecordPage), ...all.filter(p => !p.isRecordPage)]
})

function onRemove(p) {
  confirmDelete({
    header: t('project.pages.removeConfirm.header'),
    message: t('project.pages.removeConfirm.message', { name: p.name }),
    onConfirm: async () => {
      try {
        await store.removePage(props.project.id, p.id)
      } catch (err) {
        $toast.toastErrorHandler(t('project.pages.toastRemoveFailed'))(err)
      }
    },
  })
}

async function refresh(id) {
  try {
    await store.loadPages(id)
  } catch (err) {
    $toast.toastErrorHandler(t('project.pages.toastLoadFailed'))(err)
  }
}

onMounted(() => refresh(props.project.id))
watch(
  () => props.project.id,
  id => id && refresh(id),
)
</script>
