<template>
  <div class="h-full overflow-auto p-4">
    <div class="flex flex-col gap-2">
      <div v-if="!disabled">
        <Button
          icon="pi pi-plus"
          :label="$t('project.agents.add')"
          size="small"
          @click="createResource?.('agent')"
        />
      </div>

      <div>
        <CFormItemList
          :items="agents"
          item-key="id"
          reveal-on-hover
          :empty-message="$t('project.agents.empty')"
          :hide-remove="disabled"
          :remove-label="$t('project.agents.remove')"
          @select="item => inspectResource?.('agent', item.id)"
          @remove="onRemove"
        >
          <template #default="{ item }">
            <CFormItemContent :title="item.name" :subtitle="item.description || ''" />
          </template>

          <template #actions="{ item }">
            <Tag
              :value="
                item.status === 'active'
                  ? $t('project.agents.active')
                  : $t('project.agents.inactive')
              "
              :severity="item.status === 'active' ? 'success' : 'secondary'"
              class="!text-xs shrink-0 me-2"
            />
          </template>

          <template #hover-actions="{ item }">
            <CRouterLinkButton
              :to="{ name: 'agentic.edit', params: { agentID: item.id } }"
              icon="pi pi-external-link"
              severity="secondary"
              text
              size="small"
              :aria-label="$t('project.agents.openBuilder')"
              :title="$t('project.agents.openBuilder')"
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

// Detail/create dialogs are mounted once by the wizard and opened through these
// injected openers; the step never mounts them itself.
const inspectResource = inject('inspectResource', null)
const createResource = inject('createResource', null)

const agents = computed(() => store.agentsFor(props.project.id))

function onRemove(a) {
  confirmDelete({
    header: t('project.agents.removeConfirm.header'),
    message: t('project.agents.removeConfirm.message', { name: a.name }),
    onConfirm: async () => {
      try {
        await store.removeAgent(props.project.id, a.id)
      } catch (err) {
        $toast.toastErrorHandler(t('project.agents.toastRemoveFailed'))(err)
      }
    },
  })
}

async function refresh(id) {
  try {
    await store.loadAgents(id)
  } catch (err) {
    $toast.toastErrorHandler(t('project.agents.toastLoadFailed'))(err)
  }
}

onMounted(() => refresh(props.project.id))
watch(
  () => props.project.id,
  id => id && refresh(id),
)
</script>
