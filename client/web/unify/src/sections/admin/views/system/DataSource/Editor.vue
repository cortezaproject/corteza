<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ pageTitle }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <Form
    v-else-if="dataSource"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <CViewContainer scroll>
      <div v-if="isEdit" class="flex justify-end gap-2 shrink-0">
        <CPermissionsButton
          v-tooltip.bottom="$t('general.label.permissions')"
          :resource="`corteza::system:dal-connection/${dataSource.connectionID}`"
          :title="dataSource.meta?.name || dataSource.handle"
          :target="dataSource.meta?.name || dataSource.handle"
        />
      </div>

      <Panel :header="$t('system.data-sources.editor.basic.title')" toggleable :collapsed="false">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <CFormGroup
            name="name"
            :label="$t('system.data-sources.editor.basic.form.name.label')"
            required
          >
            <InputText
              id="name"
              name="name"
              v-model="dataSource.meta.name"
              :placeholder="$t('system.data-sources.editor.basic.form.name.placeholder')"
            />
          </CFormGroup>

          <CFormGroup
            name="handle"
            :label="$t('system.data-sources.editor.basic.form.handle.label')"
          >
            <InputText
              id="handle"
              name="handle"
              v-model="dataSource.handle"
              :placeholder="$t('system.data-sources.editor.basic.form.handle.placeholder')"
            />
          </CFormGroup>

          <CFormGroup
            :label="$t('system.data-sources.editor.basic.form.location-name.label')"
            :description="$t('system.data-sources.editor.basic.form.location-name.description')"
            input-id="locationName"
          >
            <InputText
              id="locationName"
              v-model="dataSource.meta.location.properties.name"
              :placeholder="$t('system.data-sources.editor.basic.form.location-name.placeholder')"
            />
          </CFormGroup>

          <CFormGroup
            :label="$t('system.data-sources.editor.basic.form.location-geometry.label')"
            :description="$t('system.data-sources.editor.basic.form.location-geometry.description')"
            input-id="locationCoords"
          >
            <CInputLocation
              :model-value="locationPoint"
              :dialog-header="
                dataSource.meta?.name ||
                $t('system.data-sources.editor.basic.form.location-geometry.label')
              "
              @update:model-value="onLocationUpdate"
            />
          </CFormGroup>

          <CFormGroup
            :label="$t('system.data-sources.editor.basic.form.ownership.label')"
            :description="$t('system.data-sources.editor.basic.form.ownership.description')"
            input-id="ownership"
          >
            <InputText
              id="ownership"
              v-model="dataSource.meta.ownership"
              :placeholder="$t('system.data-sources.editor.basic.form.ownership.placeholder')"
            />
          </CFormGroup>
        </div>
      </Panel>

      <Panel
        v-if="isEdit && dataSource.meta.properties"
        :header="$t('system.data-sources.editor.properties.title')"
        toggleable
        :collapsed="false"
      >
        <small class="block text-muted-color mb-4">
          {{ $t('system.data-sources.editor.properties.intro') }}
        </small>
        <div class="flex flex-col">
          <template v-for="(prop, index) in propertyKeys" :key="prop">
            <Divider v-if="index > 0" class="my-4" />
            <div class="flex flex-col gap-3">
              <CInputToggleCard
                v-model="dataSource.meta.properties[prop].enabled"
                :label="
                  $t(`system.data-sources.editor.properties.form.${kebabCase(prop)}.checkbox.label`)
                "
                :description="
                  $t(
                    `system.data-sources.editor.properties.form.${kebabCase(prop)}.checkbox.description`,
                  )
                "
              />
              <CFormGroup
                :label="
                  $t(`system.data-sources.editor.properties.form.${kebabCase(prop)}.notes.label`)
                "
                :description="
                  $t(
                    `system.data-sources.editor.properties.form.${kebabCase(prop)}.notes.description`,
                  )
                "
                :input-id="`${prop}Notes`"
                class="ml-2"
              >
                <Textarea
                  :id="`${prop}Notes`"
                  v-model="dataSource.meta.properties[prop].notes"
                  rows="3"
                  autoResize
                />
              </CFormGroup>
            </div>
          </template>
        </div>
      </Panel>

      <Panel
        v-if="isEdit && canManageDal"
        :header="$t('system.data-sources.editor.dal.title')"
        toggleable
        :collapsed="false"
      >
        <div class="flex flex-col">
          <template v-if="dataSource.issues?.length">
            <Message
              v-for="issue in dataSource.issues"
              :key="issue.issue"
              severity="error"
              :closable="false"
              class="mb-3"
            >
              {{ issue.issue }}
            </Message>
          </template>

          <CFormGroup
            :label="$t('system.data-sources.editor.dal.form.model-ident.label')"
            :description="
              $t('system.data-sources.editor.dal.form.model-ident.description', {
                interpolation: { prefix: '{{{', suffix: '}}}' },
              })
            "
            input-id="modelIdent"
          >
            <InputText
              id="modelIdent"
              v-model="dataSource.config.dal.modelIdent"
              :placeholder="$t('system.data-sources.editor.dal.form.model-ident.placeholder')"
            />
          </CFormGroup>

          <Divider class="my-4" />

          <CFormGroup
            :label="$t('system.data-sources.editor.dal.form.type.label')"
            :description="$t('system.data-sources.editor.dal.form.type.description')"
            input-id="dalType"
          >
            <InputText
              id="dalType"
              v-model="dataSource.config.dal.type"
              :placeholder="$t('system.data-sources.editor.dal.form.type.placeholder')"
            />
          </CFormGroup>

          <Divider class="my-4" />

          <CFormGroup
            name="dalParams"
            :label="$t('system.data-sources.editor.dal.form.params.label')"
            :description="$t('system.data-sources.editor.dal.form.params.description')"
          >
            <Textarea
              id="dalParams"
              name="dalParams"
              v-model="rawDalParams"
              rows="5"
              autoResize
              class="font-mono text-sm"
              :placeholder="$t('system.data-sources.editor.dal.form.params.placeholder')"
              @blur="parseDalParams"
            />
          </CFormGroup>
        </div>
      </Panel>

      <Message v-if="isEdit && !canManageDal" severity="warn" :closable="false">
        {{ $t('system.data-sources.editor.dal.no-access-warning') }}
      </Message>
    </CViewContainer>

    <CEditorActions :back-to="{ name: 'system.dataSources' }">
      <CInputDelete
        v-if="isEdit && dataSource.canDeleteConnection"
        :label="$t('system.data-sources.editor.delete')"
        :message="$t('general.confirm.delete')"
        :header="dataSource.meta?.name || dataSource.handle"
        :disabled="deleting"
        @confirm="handleDelete"
      />
      <Button type="submit" :label="$t('general.label.save')" icon="pi pi-save" :loading="saving" />
    </CEditorActions>
  </Form>
</template>

<script setup>
import { computed, inject, nextTick, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { system } from '@planetcrust/human-js'
import { components, useUnsavedGuard } from '@planetcrust/human-vue'
import { cloneDeep, isEqual, kebabCase } from 'lodash-es'

const { CInputDelete, CInputLocation, CInputToggleCard, CViewContainer } = components

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const dataSource = ref(null)
const initialDataSource = ref(null)
const rawDalParams = ref('{}')
const initialRawDalParams = ref('{}')

const propertyKeys = [
  'dataAtRestEncryption',
  'dataAtRestProtection',
  'dataAtTransitEncryption',
  'dataRestoration',
]

const isEdit = computed(() => !!route.params.connectionID)

const canManageDal = computed(() => !!dataSource.value?.canManageDalConfig)

const pageTitle = computed(() =>
  isEdit.value
    ? t('system.data-sources.editor.title.edit')
    : t('system.data-sources.editor.title.create'),
)

const locationPoint = computed(() => {
  const coords = dataSource.value?.meta?.location?.geometry?.coordinates
  if (!Array.isArray(coords) || coords.length !== 2) return null
  const [lng, lat] = coords
  if (typeof lng !== 'number' || typeof lat !== 'number') return null
  return { type: 'Point', coordinates: [lng, lat] }
})

function onLocationUpdate(point) {
  if (!dataSource.value) return
  ensureLocationShape()
  if (!point?.coordinates) {
    dataSource.value.meta.location.geometry.coordinates = []
    return
  }
  dataSource.value.meta.location.geometry.coordinates = [...point.coordinates]
}

function ensureLocationShape() {
  const meta = dataSource.value.meta
  if (!meta.location || typeof meta.location !== 'object') {
    meta.location = {
      type: 'Feature',
      geometry: { type: 'Point', coordinates: [] },
      properties: { name: '' },
    }
    return
  }
  if (!meta.location.type) meta.location.type = 'Feature'
  if (!meta.location.geometry || typeof meta.location.geometry !== 'object') {
    meta.location.geometry = { type: 'Point', coordinates: [] }
  } else if (!meta.location.geometry.type) {
    meta.location.geometry.type = 'Point'
  }
  if (!meta.location.properties || typeof meta.location.properties !== 'object') {
    meta.location.properties = { name: '' }
  }
}

const initialValues = computed(() => ({
  name: dataSource.value?.meta?.name || '',
  handle: dataSource.value?.handle || '',
  dalParams: rawDalParams.value,
}))

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.name || values.name.trim().length === 0) {
    errors.name = [{ message: t('general.label.required') }]
  }

  if (values.handle && !/^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$|^[A-Za-z]$/.test(values.handle)) {
    errors.handle = [
      { message: t('system.data-sources.editor.basic.form.handle.invalid-characters') },
    ]
  }

  if (values.dalParams) {
    try {
      const parsed = JSON.parse(values.dalParams)
      if (typeof parsed !== 'object' || parsed === null || Array.isArray(parsed)) {
        errors.dalParams = [
          { message: t('system.data-sources.editor.dal.form.params.description') },
        ]
      }
    } catch {
      errors.dalParams = [{ message: t('system.data-sources.editor.dal.form.params.description') }]
    }
  }

  return { errors }
})

function parseDalParams() {
  try {
    const parsed = JSON.parse(rawDalParams.value || '{}')
    if (parsed && typeof parsed === 'object' && !Array.isArray(parsed) && dataSource.value) {
      if (!dataSource.value.config.dal) dataSource.value.config.dal = {}
      dataSource.value.config.dal.params = parsed
    }
  } catch {
    // resolver will surface the error
  }
}

function initRawFields() {
  if (!dataSource.value) return
  rawDalParams.value = JSON.stringify(dataSource.value.config?.dal?.params || { dsn: '' }, null, 2)
  initialRawDalParams.value = rawDalParams.value
  ensureLocationShape()
}

async function loadDataSource() {
  const connectionID = route.params.connectionID
  if (!connectionID) {
    dataSource.value = new system.DalConnection({})
    initRawFields()
    initialDataSource.value = cloneDeep(dataSource.value)
    return
  }

  loading.value = true
  try {
    const raw = await $SystemAPI.dalConnectionRead({ connectionID })
    dataSource.value = new system.DalConnection(raw)
    initRawFields()
    initialDataSource.value = cloneDeep(dataSource.value)
  } catch (e) {
    $toast.toastErrorHandler(t('notification.data-source.fetch.error'))(e)
    router.push({ name: 'system.dataSources' })
  } finally {
    loading.value = false
  }
}

async function handleSubmit({ valid }) {
  if (!valid) {
    $toast.toastWarning(t('general.notification.formErrors'))
    nextTick(() => {
      document
        .querySelector('.p-message-error')
        ?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    })
    return
  }

  saving.value = true
  try {
    parseDalParams()

    const payload = {
      handle: dataSource.value.handle,
      type: dataSource.value.type,
      meta: dataSource.value.meta,
      config: dataSource.value.config,
    }

    if (isEdit.value) {
      payload.connectionID = dataSource.value.connectionID
      const raw = await $SystemAPI.dalConnectionUpdate(payload)
      dataSource.value = new system.DalConnection(raw)
      initRawFields()
      initialDataSource.value = cloneDeep(dataSource.value)
      $toast.toastSuccess(t('notification.data-source.update.success'))
    } else {
      const created = await $SystemAPI.dalConnectionCreate(payload)
      $toast.toastSuccess(t('notification.data-source.create.success'))
      markSaved()
      router.push({
        name: 'system.dataSources.edit',
        params: { connectionID: created.connectionID },
      })
    }
  } catch (e) {
    $toast.toastErrorHandler(
      t(
        `notification.data-source.${isEdit.value ? 'update' : 'create'}.error`,
        'Failed to save data source',
      ),
    )(e)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await $SystemAPI.dalConnectionDelete({ connectionID: dataSource.value.connectionID })
    $toast.toastSuccess(t('notification.data-source.delete.success'))
    router.push({ name: 'system.dataSources' })
  } catch (e) {
    $toast.toastErrorHandler(t('notification.data-source.delete.error'))(e)
  } finally {
    deleting.value = false
  }
}

const { markSaved } = useUnsavedGuard({
  isDirty: () =>
    !saving.value &&
    !deleting.value &&
    !!dataSource.value &&
    !!initialDataSource.value &&
    (!isEqual(dataSource.value, initialDataSource.value) ||
      rawDalParams.value !== initialRawDalParams.value),
  messageKey: 'general.editor.unsavedChanges',
})

onMounted(() => loadDataSource())

watch(
  () => route.params.connectionID,
  () => loadDataSource(),
)
</script>
