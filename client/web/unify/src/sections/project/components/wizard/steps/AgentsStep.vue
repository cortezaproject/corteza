<template>
  <div class="h-full overflow-auto p-4">
    <CFormGroup :label="$t('project.agents.title')">
      <template #actions>
        <Button
          v-if="!disabled"
          icon="pi pi-plus"
          :label="$t('project.agents.add')"
          severity="secondary"
          size="small"
          @click="openCreate"
        />
      </template>

      <div class="mt-1">
        <CFormItemList
          :items="agents"
          item-key="id"
          reveal-on-hover
          :empty-message="$t('project.agents.empty')"
          :hide-remove="disabled"
          :remove-label="$t('project.agents.remove')"
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
    </CFormGroup>

    <ConfigureAgentDialog
      v-model="dialogOpen"
      :project-id="project.id"
      :agent="activeAgent"
      @saved="refresh(project.id)"
    />
  </div>
</template>

<script setup>
import ConfigureAgentDialog from '@/sections/project/components/agents/ConfigureAgentDialog.vue'
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

const cfg = kindConfig('agent')
const agents = computed(() => store.agentsFor(props.project.id))

const dialogOpen = ref(false)
const activeAgent = ref(null)

function openCreate() {
  activeAgent.value = null
  dialogOpen.value = true
}

// Clicking an agent opens its dialog (name/description + a button to open the
// full agent builder), rather than navigating straight to the builder.
function onSelect(item) {
  activeAgent.value = item
  dialogOpen.value = true
}

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
