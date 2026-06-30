<template>
  <div class="h-full overflow-auto p-4">
    <CFormGroup :label="$t('project.chatbots.title')">
      <template #actions>
        <Button
          v-if="!disabled"
          icon="pi pi-plus"
          :label="$t('project.chatbots.add')"
          severity="secondary"
          size="small"
          @click="openCreate"
        />
      </template>

      <div class="mt-1">
        <CFormItemList
          :items="chatbots"
          item-key="id"
          reveal-on-hover
          :empty-message="$t('project.chatbots.empty')"
          :hide-remove="disabled"
          :remove-label="$t('project.chatbots.remove')"
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
              </div>
            </div>
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
    </CFormGroup>

    <ConfigureChatbotDialog
      v-model="dialogOpen"
      :project-id="project.id"
      :chatbot="activeChatbot"
      @saved="refresh(project.id)"
    />
  </div>
</template>

<script setup>
import ConfigureChatbotDialog from '@/sections/project/components/chatbots/ConfigureChatbotDialog.vue'
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

const cfg = kindConfig('chatbot')
const chatbots = computed(() => store.chatbotsFor(props.project.id))

const dialogOpen = ref(false)
const activeChatbot = ref(null)

function openCreate() {
  activeChatbot.value = null
  dialogOpen.value = true
}

// Clicking a chatbot opens its dialog (name + a button to open the full chatbot
// builder), rather than navigating straight to the builder.
function onSelect(item) {
  activeChatbot.value = item
  dialogOpen.value = true
}

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
