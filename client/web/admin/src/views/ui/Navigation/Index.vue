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

          <div v-if="$Settings.get('discovery.enabled', false)" class="flex items-center gap-2">
            <Checkbox v-model="hideSearch" :binary="true" inputId="hideSearch" />
            <label for="hideSearch">
              {{ $t('ui.settings.editor.topbar.search.hide') }}
            </label>
          </div>

          <div class="flex items-center gap-2">
            <Checkbox v-model="topbar.hideProfile" :binary="true" inputId="hideProfile" />
            <label for="hideProfile">
              {{ $t('ui.settings.editor.topbar.profile.hide') }}
            </label>
          </div>
        </div>

        <Fieldset :legend="$t('ui.settings.editor.topbar.profile.title')" class="mb-4">
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
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

            <CFormGroup :label="$t('ui.settings.editor.topbar.links.title')" class="lg:col-span-2">
              <template #actions>
                <Button
                  :label="$t('general.label.add')"
                  icon="pi pi-plus"
                  severity="secondary"
                  size="small"
                  @click="topbar.profileLinks.push({ handle: '', url: '', newTab: true })"
                />
              </template>
              <CFormList
                v-model="topbar.profileLinks"
                :columns="[
                  { label: $t('ui.settings.editor.topbar.links.handle'), width: '1fr' },
                  { label: $t('ui.settings.editor.topbar.links.url'), width: '1.5fr' },
                  {
                    label: $t('ui.settings.editor.topbar.links.new-tab'),
                    width: '5rem',
                    headerClass: 'text-center',
                  },
                ]"
              >
                <template #row="{ item }">
                  <InputText v-model="item.handle" size="small" class="w-full" />
                  <InputText v-model="item.url" size="small" class="w-full" />
                  <div class="flex items-center justify-center">
                    <Checkbox v-model="item.newTab" :binary="true" />
                  </div>
                </template>
              </CFormList>
            </CFormGroup>
          </div>
        </Fieldset>

        <!-- Page buttons sub-section -->
        <CFormGroup
          :label="$t('ui.settings.editor.topbar.page-buttons.title')"
          :description="$t('ui.settings.editor.topbar.page-buttons.description')"
        >
          <template #actions>
            <Button
              :label="$t('general.label.add')"
              icon="pi pi-plus"
              severity="secondary"
              size="small"
              @click="addPageButton"
            />
          </template>
          <CFormList
            v-model="topbar.pageButtons"
            :columns="[
              { label: $t('ui.settings.editor.topbar.page-buttons.label'), width: '1fr' },
              { label: $t('ui.settings.editor.topbar.page-buttons.url'), width: '1fr' },
              {
                label: $t('ui.settings.editor.topbar.page-buttons.url-match'),
                tooltip: $t('ui.settings.editor.topbar.page-buttons.url-match-description'),
                width: '1fr',
              },
              {
                label: $t('ui.settings.editor.topbar.page-buttons.new-tab'),
                width: '5rem',
                headerClass: 'text-center',
              },
            ]"
            class="mt-3"
          >
            <template #row="{ item }">
              <InputText v-model="item.label" size="small" class="w-full" />
              <InputText v-model="item.url" size="small" class="w-full" />
              <InputText
                v-model="item.urlMatch"
                size="small"
                class="w-full"
                placeholder="/builder"
              />
              <div class="flex items-center justify-center">
                <Checkbox v-model="item.newTab" :binary="true" />
              </div>
            </template>
            <template #extra="{ item }">
              <InputText
                v-model="item.description"
                size="small"
                class="w-full"
                :placeholder="
                  $t('ui.settings.editor.topbar.page-buttons.button-description-placeholder')
                "
              />
            </template>
          </CFormList>
        </CFormGroup>
      </Panel>
    </div>

    <CEditorActions>
      <Button
        :label="$t('general.label.save')"
        icon="pi pi-save"
        :loading="saving"
        @click="handleSave"
      />
    </CEditorActions>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')
const $Settings = inject('$Settings')

const loading = ref(false)
const saving = ref(false)

const topbar = reactive({
  hideAppSelector: false,
  hideNotifications: false,
  showSearch: true,
  hideProfileLink: false,
  hideChangePasswordLink: false,
  hideThemeSelector: false,
  profileLinks: [],
  pageButtons: [],
})

const hideDrafts = computed({
  get: () => topbar.showDrafts !== true,
  set: val => {
    topbar.showDrafts = !val
  },
})

const hideSearch = computed({
  get: () => topbar.showSearch !== true,
  set: val => {
    topbar.showSearch = !val
  },
})

function addPageButton() {
  topbar.pageButtons.push({ label: '', url: '', urlMatch: '', newTab: true, description: '' })
}

async function loadSettings() {
  loading.value = true
  try {
    const result = await $SystemAPI.settingsList({ prefix: 'ui.topbar' })
    for (const s of result || []) {
      if (s.name === 'ui.topbar' && s.value) {
        Object.assign(topbar, {
          ...s.value,
          profileLinks: s.value.profileLinks || [],
          pageButtons: s.value.pageButtons || [],
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
