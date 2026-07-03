<template>
  <div class="h-full overflow-auto p-4">
    <div class="flex flex-col gap-2">
      <div>
        <Button
          v-if="!disabled"
          icon="pi pi-plus"
          :label="$t('project.chatbots.add')"
          severity="secondary"
          size="small"
          @click="createResource?.('chatbot')"
        />
      </div>

      <div class="mt-1">
        <CFormItemList
          :items="chatbots"
          item-key="id"
          reveal-on-hover
          :empty-message="$t('project.chatbots.empty')"
          :hide-remove="disabled"
          :remove-label="$t('project.chatbots.remove')"
          @select="item => inspectResource?.('chatbot', item.id)"
          @remove="onRemove"
        >
          <template #default="{ item }">
            <CFormItemContent :title="item.name" />
          </template>

          <template #actions="{ item }">
            <Tag
              :value="
                item.enabled
                  ? $t('project.chatbots.enabled')
                  : $t('project.chatbots.disabled')
              "
              :severity="item.enabled ? 'success' : 'secondary'"
              class="!text-xs shrink-0 me-2"
            />
          </template>

          <template #hover-actions="{ item }">
            <CRouterLinkButton
              :to="{ name: 'chatbot.edit', params: { chatbotID: item.id } }"
              icon="pi pi-external-link"
              severity="secondary"
              text
              size="small"
              :aria-label="$t('project.chatbots.openBuilder')"
              :title="$t('project.chatbots.openBuilder')"
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

// Detail/create dialogs are mounted once in Wizard.vue and opened through these
// injected helpers; the step never mounts them itself.
const inspectResource = inject('inspectResource', null)
const createResource = inject('createResource', null)

const chatbots = computed(() => store.chatbotsFor(props.project.id))

function onRemove(c) {
  confirmDelete({
    header: t('project.chatbots.removeConfirm.header'),
    message: t('project.chatbots.removeConfirm.message', { name: c.name }),
    onConfirm: async () => {
      try {
        await store.removeChatbot(props.project.id, c.id)
      } catch (err) {
        $toast.toastErrorHandler(t('project.chatbots.toastRemoveFailed'))(err)
      }
    },
  })
}

async function refresh(id) {
  try {
    await store.loadChatbots(id)
  } catch (err) {
    $toast.toastErrorHandler(t('project.chatbots.toastLoadFailed'))(err)
  }
}

onMounted(() => refresh(props.project.id))
watch(
  () => props.project.id,
  id => id && refresh(id),
)
</script>
