<template>
  <div class="h-full overflow-auto p-4">
    <CFormGroup :label="$t('project.automations.title')">
      <template #actions>
        <Button
          v-if="!disabled"
          icon="pi pi-plus"
          :label="$t('project.automations.add')"
          severity="secondary"
          size="small"
          @click="openCreate"
        />
      </template>

      <div class="mt-1">
        <CFormItemList
          :items="automations"
          item-key="id"
          :reveal-on-hover="!disabled"
          :empty-message="$t('project.automations.empty')"
          :hide-remove="disabled"
          :remove-label="$t('project.automations.remove')"
          @select="onSelect"
          @remove="onRemove"
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
                <div v-if="item.description" class="text-xs text-muted-color truncate">
                  {{ item.description }}
                </div>
              </div>
            </div>
          </template>

          <template #actions="{ item }">
            <Tag
              :value="
                item.enabled
                  ? $t('project.automations.enabled')
                  : $t('project.automations.disabled')
              "
              :severity="item.enabled ? 'success' : 'secondary'"
              class="!text-xs shrink-0"
            />
          </template>

          <template #hover-actions="{ item }">
            <CRouterLinkButton
              :to="{ name: 'taq.builder-edit', params: { id: item.id } }"
              icon="pi pi-external-link"
              severity="secondary"
              text
              size="small"
              :aria-label="$t('project.automations.openBuilder')"
              :title="$t('project.automations.openBuilder')"
              @click.stop
            />
          </template>
        </CFormItemList>
      </div>
    </CFormGroup>

    <ConfigureAutomationDialog
      v-model="dialogOpen"
      :project-id="project.id"
      :automation="activeAutomation"
      @saved="refresh(project.id)"
    />
  </div>
</template>

<script setup>
import ConfigureAutomationDialog from '@/sections/project/components/automations/ConfigureAutomationDialog.vue'
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

const cfg = kindConfig('automation')
const automations = computed(() => store.automationsFor(props.project.id))

const dialogOpen = ref(false)
const activeAutomation = ref(null)

function openCreate() {
  activeAutomation.value = null
  dialogOpen.value = true
}

// Clicking an automation opens its dialog (name/description + a button to open
// the full TAQ builder), rather than navigating straight to the builder.
function onSelect(item) {
  activeAutomation.value = item
  dialogOpen.value = true
}

function onRemove(a) {
  confirmDelete({
    header: t('project.automations.removeConfirm.header'),
    message: t('project.automations.removeConfirm.message', { name: a.name }),
    onConfirm: async () => {
      try {
        await store.removeAutomation(props.project.id, a.id)
      } catch (err) {
        $toast.toastErrorHandler(t('project.automations.toastRemoveFailed'))(err)
      }
    },
  })
}

async function refresh(id) {
  try {
    await store.loadAutomations(id)
  } catch (err) {
    $toast.toastErrorHandler(t('project.automations.toastLoadFailed'))(err)
  }
}

onMounted(() => refresh(props.project.id))
watch(
  () => props.project.id,
  id => id && refresh(id),
)
</script>
