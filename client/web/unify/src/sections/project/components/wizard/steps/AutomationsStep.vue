<template>
  <div class="h-full overflow-auto p-4">
    <div class="flex flex-col gap-2">
      <div v-if="!disabled">
        <Button
          icon="pi pi-plus"
          :label="$t('project.automations.add')"
          size="small"
          @click="createResource?.('automation')"
        />
      </div>

      <div>
        <CFormItemList
          :items="automations"
          item-key="id"
          :reveal-on-hover="!disabled"
          :empty-message="$t('project.automations.empty')"
          :hide-remove="disabled"
          :remove-label="$t('project.automations.remove')"
          @select="item => inspectResource?.('automation', item.id)"
          @remove="onRemove"
        >
          <template #default="{ item }">
            <CFormItemContent :title="item.name" :subtitle="item.description || ''" />
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

// Detail/create dialogs live once in Wizard.vue; steps just ask it to open them.
const inspectResource = inject('inspectResource', null)
const createResource = inject('createResource', null)

const automations = computed(() => store.automationsFor(props.project.id))

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
