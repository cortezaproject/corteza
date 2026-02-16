<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ pageTitle }}</span>
  </Teleport>

  <Teleport to="#topbar-tools" defer>
    <div v-if="isEdit" class="flex gap-1">
      <Button
        :label="$t('namespace.visit')"
        icon="pi pi-external-link"
        size="small"
        :disabled="!namespace?.enabled"
        @click="visitNamespace"
      />
    </div>
  </Teleport>

  <!-- Loading -->
  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <!-- Form -->
  <div v-else-if="namespace" class="flex flex-col h-full">
    <div class="flex-1 overflow-auto py-6">
      <div class="container mx-auto max-w-3xl px-4">
        <Card>
          <template #content>
            <form class="flex flex-col gap-5" @submit.prevent="handleSubmit">
              <!-- Name -->
              <div class="flex flex-col gap-2">
                <label for="name" class="font-medium text-primary">
                  {{ $t('namespace.name.label') }}
                </label>
                <InputText
                  id="name"
                  v-model="namespace.name"
                  :placeholder="$t('namespace.name.placeholder')"
                  :invalid="!nameValid"
                />
                <small v-if="!nameValid" class="text-red-500">
                  {{ $t('general.label.required') }}
                </small>
              </div>

              <!-- Slug -->
              <div class="flex flex-col gap-2">
                <label for="slug" class="font-medium text-primary">
                  {{ $t('namespace.slug.label') }}
                </label>
                <InputText
                  id="slug"
                  v-model="namespace.slug"
                  :placeholder="$t('namespace.slug.placeholder')"
                  :invalid="slugState === false"
                />
                <small class="text-muted-color">
                  {{ $t('namespace.slug.description') }}
                </small>
                <small v-if="slugState === false" class="text-red-500">
                  {{ $t('namespace.slug.invalid-handle-characters') }}
                </small>
              </div>

              <!-- Enabled -->
              <div class="flex items-center gap-2">
                <Checkbox id="enabled" v-model="namespace.enabled" binary />
                <label for="enabled">{{ $t('namespace.enabled.label') }}</label>
              </div>

              <Divider />

              <!-- Logo -->
              <div class="flex items-center gap-2">
                <Checkbox id="logoEnabled" v-model="namespace.meta.logoEnabled" binary />
                <label for="logoEnabled">{{ $t('namespace.logo.show') }}</label>
              </div>

              <!-- Subtitle -->
              <div class="flex flex-col gap-2">
                <label for="subtitle" class="font-medium text-primary">
                  {{ $t('namespace.subtitle.label') }}
                </label>
                <InputText
                  id="subtitle"
                  v-model="namespace.meta.subtitle"
                  :placeholder="$t('namespace.subtitle.placeholder')"
                />
              </div>

              <!-- Description -->
              <div class="flex flex-col gap-2">
                <label for="description" class="font-medium text-primary">
                  {{ $t('namespace.description.label') }}
                </label>
                <Textarea
                  id="description"
                  v-model="namespace.meta.description"
                  :placeholder="$t('namespace.description.placeholder')"
                  rows="3"
                  auto-resize
                />
              </div>

              <Divider />

              <!-- Sidebar -->
              <div class="flex items-center gap-2">
                <Checkbox id="hideSidebar" v-model="namespace.meta.hideSidebar" binary />
                <label for="hideSidebar">{{ $t('namespace.sidebar.hide') }}</label>
              </div>
            </form>
          </template>
        </Card>
      </div>
    </div>

    <!-- Toolbar -->
    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="flex items-center justify-between p-3">
        <Button
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="$router.back()"
        />
        <div class="flex gap-2">
          <CInputDelete
            v-if="isEdit && namespace.canDeleteNamespace"
            :label="$t('general.label.delete')"
            :message="$t('namespace.deleteConfirm')"
            :header="namespace.name"
            :disabled="deleting"
            @confirm="handleDelete"
          />
          <Button
            :label="$t('general.label.save')"
            icon="pi pi-save"
            :loading="saving"
            :disabled="!canSave"
            @click="handleSubmit"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useNamespaceStore } from '@/stores/namespace'
import { compose } from '@cortezaproject/corteza-js-next'
import { components } from '@cortezaproject/corteza-vue-next'
import { computed, inject, onMounted, ref, watch } from 'vue'

const { CInputDelete } = components
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const $toast = inject('$toast')
const namespaceStore = useNamespaceStore()

// State
const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const namespace = ref(null)

// Computed
const isEdit = computed(() => !!route.params.slug)

const pageTitle = computed(() => {
  return isEdit.value ? t('namespace.edit') : t('namespace.create')
})

const nameValid = computed(() => {
  return namespace.value?.name?.length > 0
})

const slugState = computed(() => {
  const slug = namespace.value?.slug
  if (!slug) return null
  // Valid handle: alphanumeric, underscore, must start with letter
  return /^[a-zA-Z][a-zA-Z0-9_]*$/.test(slug)
})

const canSave = computed(() => {
  if (isEdit.value && !namespace.value?.canUpdateNamespace) return false
  return nameValid.value && slugState.value !== false
})

// Methods
async function loadNamespace() {
  const slug = route.params.slug
  if (!slug) {
    // Create new
    namespace.value = new compose.Namespace({
      enabled: true,
      meta: {
        subtitle: '',
        description: '',
        hideSidebar: false,
        logoEnabled: false,
      },
    })
    return
  }

  loading.value = true
  try {
    // Find by slug or ID
    const found = namespaceStore.getByUrlPart(slug)
    if (found) {
      namespace.value = new compose.Namespace({ ...found })
    } else {
      // Load from API
      await namespaceStore.load({ force: true })
      const ns = namespaceStore.getByUrlPart(slug)
      if (ns) {
        namespace.value = new compose.Namespace({ ...ns })
      } else {
        $toast.toastDanger(t('notification.namespace.loadFailed'))
        router.push({ name: 'namespace.manage' })
      }
    }
  } catch (e) {
    console.error('Failed to load namespace:', e)
    $toast.toastDanger(t('notification.namespace.loadFailed'))
  } finally {
    loading.value = false
  }
}

async function handleSubmit() {
  if (!canSave.value) return

  saving.value = true
  try {
    const payload = {
      name: namespace.value.name,
      slug: namespace.value.slug,
      enabled: namespace.value.enabled,
      meta: namespace.value.meta,
    }

    if (isEdit.value) {
      payload.namespaceID = namespace.value.namespaceID
      await namespaceStore.update(payload)
      $toast.toastSuccess(t('notification.namespace.saved'))
    } else {
      const created = await namespaceStore.create(payload)
      $toast.toastSuccess(t('notification.namespace.saved'))
      // Navigate to edit view
      router.push({
        name: 'namespace.edit',
        params: { slug: created.slug || created.namespaceID },
      })
    }
  } catch (e) {
    console.error('Failed to save namespace:', e)
    $toast.toastDanger(t('notification.namespace.saveFailed'))
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await namespaceStore.delete({ namespaceID: namespace.value.namespaceID })
    $toast.toastSuccess(t('notification.namespace.deleted'))
    router.push({ name: 'namespace.manage' })
  } catch (e) {
    console.error('Failed to delete namespace:', e)
    $toast.toastDanger(t('notification.namespace.deleteFailed'))
  } finally {
    deleting.value = false
  }
}

function visitNamespace() {
  router.push({
    name: 'namespace.view',
    params: { slug: namespace.value.slug || namespace.value.namespaceID },
  })
}

// Lifecycle
onMounted(() => {
  loadNamespace()
})

watch(
  () => route.params.slug,
  () => {
    loadNamespace()
  },
)
</script>
