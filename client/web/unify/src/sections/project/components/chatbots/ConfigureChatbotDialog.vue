<template>
  <Dialog
    v-model:visible="visible"
    modal
    :header="isEdit ? $t('project.configureChatbot.editTitle') : $t('project.configureChatbot.createTitle')"
    :style="{ width: '40rem' }"
    :pt="{ content: { class: '!pt-2' }, footer: { class: 'flex justify-between gap-2' } }"
  >
    <p class="text-sm text-muted-color mb-4">
      {{ $t('project.configureChatbot.blurb') }}
    </p>

    <div class="flex flex-col gap-4">
      <CFormGroup :label="$t('project.configureChatbot.name')" required>
        <InputText
          v-model="name"
          fluid
          :placeholder="$t('project.configureChatbot.namePlaceholder')"
          @keyup.enter="save()"
        />
      </CFormGroup>

      <CInputToggleCard
        v-if="isEdit"
        v-model="enabled"
        :label="$t('project.configureChatbot.enabled')"
        :description="$t('project.configureChatbot.enabledHint')"
      />
    </div>

    <template #footer>
      <CRouterLinkButton
        v-if="isEdit"
        :to="{ name: 'chatbot.edit', params: { chatbotID: chatbot.id } }"
        target="_blank"
        rel="noopener"
        :label="$t('project.configureChatbot.openBuilder')"
        icon="pi pi-external-link"
        severity="secondary"
        size="small"
      />
      <span v-else />

      <div class="flex gap-2">
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          text
          size="small"
          @click="visible = false"
        />
        <Button
          v-if="isEdit"
          :label="$t('general.label.save')"
          size="small"
          :loading="saving"
          :disabled="!name.trim()"
          @click="save()"
        />
        <template v-else>
          <Button
            :label="$t('project.configureChatbot.create')"
            outlined
            size="small"
            :loading="saving && !openingBuilder"
            :disabled="!name.trim() || saving"
            @click="save(false)"
          />
          <Button
            :label="$t('project.configureChatbot.createAndOpenBuilder')"
            icon="pi pi-external-link"
            icon-pos="right"
            size="small"
            :loading="saving && openingBuilder"
            :disabled="!name.trim() || saving"
            @click="save(true)"
          />
        </template>
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { useProjectsStore } from '@/sections/project/stores/projects'
import { components } from '@planetcrust/human-vue'
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const { CRouterLinkButton } = components

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  projectId: { type: [String, Number], required: true },
  // Existing chatbot to edit; null opens the dialog in create mode.
  chatbot: { type: Object, default: null },
})
const emit = defineEmits(['update:modelValue', 'saved'])

const store = useProjectsStore()
const { t } = useI18n()
const router = useRouter()
const $toast = inject('$toast')

const visible = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v),
})

const isEdit = computed(() => !!props.chatbot?.id)

const name = ref('')
const enabled = ref(false)
const saving = ref(false)
// Which create button is in flight, so only it shows a spinner.
const openingBuilder = ref(false)

// (Re)seed the form whenever the dialog opens.
watch(visible, open => {
  if (!open) return
  name.value = props.chatbot?.name || ''
  enabled.value = !!props.chatbot?.enabled
})

// Create just creates the chatbot; editing only updates name/enabled. In create
// mode `openBuilder` also opens the new chatbot in the full builder (new tab, so
// the wizard stays put). Editing reaches the builder via the footer link instead.
async function save(openBuilder = false) {
  if (!name.value.trim() || saving.value) return
  saving.value = true
  openingBuilder.value = openBuilder
  // Open the tab synchronously inside the click so it isn't blocked as a popup
  // after the await; we point it at the builder once the chatbot exists, or
  // close it if the create fails.
  const builderTab = openBuilder ? window.open('', '_blank') : null
  try {
    if (isEdit.value) {
      await store.updateChatbot(props.projectId, props.chatbot.id, {
        name: name.value.trim(),
        enabled: enabled.value,
      })
      emit('saved')
    } else {
      const id = await store.addChatbot(props.projectId, { name: name.value.trim() })
      if (builderTab) {
        const { href } = router.resolve({ name: 'chatbot.edit', params: { chatbotID: id } })
        builderTab.location.href = new URL(href, window.location.href).href
      }
      emit('saved', id)
    }
    visible.value = false
  } catch (err) {
    builderTab?.close()
    $toast.toastErrorHandler(t('project.configureChatbot.toastFailed'))(err)
  } finally {
    saving.value = false
    openingBuilder.value = false
  }
}
</script>
