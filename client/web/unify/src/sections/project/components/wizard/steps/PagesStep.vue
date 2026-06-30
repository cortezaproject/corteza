<template>
  <div class="h-full overflow-auto p-4">
    <CFormGroup :label="$t('project.pages.title')">
      <template #actions>
        <Button
          v-if="!disabled"
          icon="pi pi-plus"
          :label="$t('project.pages.add')"
          severity="secondary"
          size="small"
          @click="openCreate"
        />
      </template>

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
          @select="onSelect"
        >
          <template #default="{ item }">
            <div class="flex items-center gap-3 min-w-0">
              <span
                class="inline-flex items-center justify-center w-8 h-8 rounded-md ring-1 shrink-0"
                :class="[cfg.bg, cfg.ring]"
              >
                <i :class="[cfg.icon, cfg.text]" />
              </span>
              <div class="min-w-0">
                <div class="font-medium truncate">{{ item.name }}</div>
                <div
                  v-if="item.isRecordPage"
                  class="text-xs text-muted-color truncate flex items-center gap-1"
                >
                  <i class="pi pi-database text-[10px]" />
                  {{ $t('project.pages.recordSubtitle') }}
                </div>
              </div>
            </div>
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
    </CFormGroup>

    <ConfigurePageDialog
      v-model="dialogOpen"
      :project-id="project.id"
      :namespace-id="project.namespaceID"
      :page="activePage"
      @saved="refresh(project.id)"
    />
  </div>
</template>

<script setup>
import ConfigurePageDialog from '@/sections/project/components/pages/ConfigurePageDialog.vue'
import { kindConfig } from '@/sections/project/config/kinds'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { components, useConfirmDelete } from '@planetcrust/human-vue'
import { computed, inject, onMounted, ref, watch } from 'vue'
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

const cfg = kindConfig('page')
// Record (module detail) pages first, then standalone — both in one list.
const pages = computed(() => {
  const all = store.pagesFor(props.project.id)
  return [...all.filter(p => p.isRecordPage), ...all.filter(p => !p.isRecordPage)]
})

const dialogOpen = ref(false)
const activePage = ref(null)

function openCreate() {
  activePage.value = null
  dialogOpen.value = true
}

function onSelect(item) {
  activePage.value = item
  dialogOpen.value = true
}

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
