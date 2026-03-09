<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('ui.settings.editor.navigation.title') }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <div v-else class="flex flex-col h-full">
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-5 overflow-y-auto">
      <Panel
        :header="$t('ui.settings.editor.topbar.title')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <!-- General -->
        <h5 class="text-sm font-semibold mb-3">
          {{ $t('ui.settings.editor.topbar.general') }}
        </h5>
        <div class="flex flex-col gap-3 mb-4">
          <div class="flex items-center gap-2">
            <Checkbox v-model="topbar.hideAppSelector" :binary="true" inputId="hideAppSelector" />
            <label for="hideAppSelector">
              {{ $t('ui.settings.editor.topbar.app-selector.hide') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <Checkbox v-model="hideDrafts" :binary="true" inputId="hideDrafts" />
            <label for="hideDrafts">
              {{ $t('ui.settings.editor.topbar.drafts.hide') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <Checkbox
              v-model="topbar.hideNotifications"
              :binary="true"
              inputId="hideNotifications"
            />
            <label for="hideNotifications">
              {{ $t('ui.settings.editor.topbar.notifications.hide') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <Checkbox v-model="topbar.hideHelp" :binary="true" inputId="hideHelp" />
            <label for="hideHelp">
              {{ $t('ui.settings.editor.topbar.help.hide') }}
            </label>
          </div>
          <div class="flex items-center gap-2">
            <Checkbox v-model="topbar.hideProfile" :binary="true" inputId="hideProfile" />
            <label for="hideProfile">
              {{ $t('ui.settings.editor.topbar.profile.hide') }}
            </label>
          </div>
        </div>

        <Divider />

        <!-- Help sub-section -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-4">
          <div>
            <h5 class="text-sm font-semibold mb-3">
              {{ $t('ui.settings.editor.topbar.help.title') }}
            </h5>
            <div class="flex flex-col gap-3">
              <div class="flex items-center gap-2">
                <Checkbox v-model="topbar.hideForumLink" :binary="true" inputId="hideForumLink" />
                <label for="hideForumLink">
                  {{ $t('ui.settings.editor.topbar.help.hide-forum-link') }}
                </label>
              </div>
              <div class="flex items-center gap-2">
                <Checkbox
                  v-model="topbar.hideDocumentationLink"
                  :binary="true"
                  inputId="hideDocLink"
                />
                <label for="hideDocLink">
                  {{ $t('ui.settings.editor.topbar.help.hide-documentation-link') }}
                </label>
              </div>
              <div class="flex items-center gap-2">
                <Checkbox
                  v-model="topbar.hideFeedbackLink"
                  :binary="true"
                  inputId="hideFeedbackLink"
                />
                <label for="hideFeedbackLink">
                  {{ $t('ui.settings.editor.topbar.help.hide-feedback-link') }}
                </label>
              </div>
            </div>
          </div>

          <div class="lg:col-span-2">
            <div class="flex flex-col gap-2">
              <label class="font-medium text-sm text-primary">
                {{ $t('ui.settings.editor.topbar.links.title') }}
              </label>
              <DataTable :value="topbar.helpLinks" class="border rounded" size="small">
                <Column :header="$t('ui.settings.editor.topbar.links.handle')" class="w-1/3">
                  <template #body="{ data }">
                    <InputText v-model="data.handle" size="small" class="w-full" />
                  </template>
                </Column>
                <Column :header="$t('ui.settings.editor.topbar.links.url')" class="w-1/2">
                  <template #body="{ data }">
                    <InputText v-model="data.url" size="small" class="w-full" />
                  </template>
                </Column>
                <Column
                  :header="$t('ui.settings.editor.topbar.links.new-tab')"
                  class="w-20 text-center"
                >
                  <template #body="{ data }">
                    <Checkbox v-model="data.newTab" :binary="true" />
                  </template>
                </Column>
                <Column class="w-16 text-right">
                  <template #body="{ index }">
                    <Button
                      icon="pi pi-trash"
                      severity="danger"
                      text
                      rounded
                      size="small"
                      @click="topbar.helpLinks.splice(index, 1)"
                    />
                  </template>
                </Column>
              </DataTable>
              <div>
                <Button
                  :label="$t('general.label.add')"
                  icon="pi pi-plus"
                  size="small"
                  severity="secondary"
                  @click="topbar.helpLinks.push({ handle: '', url: '', newTab: true })"
                />
              </div>
            </div>
          </div>
        </div>

        <Divider />

        <!-- Profile sub-section -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div>
            <h5 class="text-sm font-semibold mb-3">
              {{ $t('ui.settings.editor.topbar.profile.title') }}
            </h5>
            <div class="flex flex-col gap-3">
              <div class="flex items-center gap-2">
                <Checkbox
                  v-model="topbar.hideProfileLink"
                  :binary="true"
                  inputId="hideProfileLink"
                />
                <label for="hideProfileLink">
                  {{ $t('ui.settings.editor.topbar.profile.hide-profile-link') }}
                </label>
              </div>
              <div class="flex items-center gap-2">
                <Checkbox
                  v-model="topbar.hideChangePasswordLink"
                  :binary="true"
                  inputId="hideChangePwLink"
                />
                <label for="hideChangePwLink">
                  {{ $t('ui.settings.editor.topbar.profile.hide-change-password-link') }}
                </label>
              </div>
              <div class="flex items-center gap-2">
                <Checkbox
                  v-model="topbar.hideThemeSelector"
                  :binary="true"
                  inputId="hideThemeSelector"
                />
                <label for="hideThemeSelector">
                  {{ $t('ui.settings.editor.topbar.profile.hide-theme-selector') }}
                </label>
              </div>
            </div>
          </div>

          <div class="lg:col-span-2">
            <div class="flex flex-col gap-2">
              <label class="font-medium text-sm text-primary">
                {{ $t('ui.settings.editor.topbar.links.title') }}
              </label>
              <DataTable :value="topbar.profileLinks" class="border rounded" size="small">
                <Column :header="$t('ui.settings.editor.topbar.links.handle')" class="w-1/3">
                  <template #body="{ data }">
                    <InputText v-model="data.handle" size="small" class="w-full" />
                  </template>
                </Column>
                <Column :header="$t('ui.settings.editor.topbar.links.url')" class="w-1/2">
                  <template #body="{ data }">
                    <InputText v-model="data.url" size="small" class="w-full" />
                  </template>
                </Column>
                <Column
                  :header="$t('ui.settings.editor.topbar.links.new-tab')"
                  class="w-20 text-center"
                >
                  <template #body="{ data }">
                    <Checkbox v-model="data.newTab" :binary="true" />
                  </template>
                </Column>
                <Column class="w-16 text-right">
                  <template #body="{ index }">
                    <Button
                      icon="pi pi-trash"
                      severity="danger"
                      text
                      rounded
                      size="small"
                      @click="topbar.profileLinks.splice(index, 1)"
                    />
                  </template>
                </Column>
              </DataTable>
              <div>
                <Button
                  :label="$t('general.label.add')"
                  icon="pi pi-plus"
                  size="small"
                  severity="secondary"
                  @click="topbar.profileLinks.push({ handle: '', url: '', newTab: true })"
                />
              </div>
            </div>
          </div>
        </div>
      </Panel>
    </div>

    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="p-3 flex items-center justify-end">
        <Button
          :label="$t('general.label.save')"
          icon="pi pi-save"
          :loading="saving"
          @click="handleSave"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const saving = ref(false)

const topbar = reactive({
  hideAppSelector: false,
  hideNotifications: false,
  hideHelp: false,
  hideProfile: false,
  showDrafts: true,
  showSearch: true,
  hideForumLink: false,
  hideDocumentationLink: false,
  hideFeedbackLink: false,
  hideProfileLink: false,
  hideChangePasswordLink: false,
  hideThemeSelector: false,
  helpLinks: [],
  profileLinks: [],
})

const hideDrafts = computed({
  get: () => topbar.showDrafts !== true,
  set: val => {
    topbar.showDrafts = !val
  },
})

async function loadSettings() {
  loading.value = true
  try {
    const result = await $SystemAPI.settingsList({ prefix: 'ui.topbar' })
    for (const s of result || []) {
      if (s.name === 'ui.topbar' && s.value) {
        Object.assign(topbar, {
          ...s.value,
          helpLinks: s.value.helpLinks || [],
          profileLinks: s.value.profileLinks || [],
        })
      }
    }
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.navigation.fetch.error'))(e)
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  try {
    await $SystemAPI.settingsUpdate({
      values: [{ name: 'ui.topbar', value: { ...topbar } }],
    })
    $toast.toastSuccess(t('notification.settings.navigation.update.success'))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.navigation.update.error'))(e)
  } finally {
    saving.value = false
  }
}

onMounted(() => loadSettings())
</script>
