<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('automation.scripts.list.title') }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <div v-else class="flex flex-col h-full">
    <CViewContainer scroll>
      <Panel
        :header="$t('automation.scripts.list.title')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <!-- Search + Filters -->
        <div class="flex flex-col gap-4 mb-4">
          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{ $t('automation.scripts.list.filter.searchQuery') }}
            </label>
            <InputText v-model="filter.query" class="w-full md:w-1/2" />
          </div>

          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <CInputSwitch
              v-model="filter.incScriptsWithErrors"
              :label="
                $t('automation.scripts.list.filter.incScriptsWithErrors', {
                  count: totalScriptsWithErrors,
                })
              "
            />
            <CInputSwitch
              v-model="filter.incScriptsWithTriggers"
              :label="
                $t('automation.scripts.list.filter.incScriptsWithTriggers', {
                  count: totalScriptsWithTriggers,
                })
              "
            />
            <CInputSwitch
              v-model="filter.incScriptsWithIterator"
              :label="
                $t('automation.scripts.list.filter.incScriptsWithIterator', {
                  count: totalScriptsWithIterator,
                })
              "
            />
            <CInputSwitch
              v-model="filter.incScriptsWithSecurity"
              :label="
                $t('automation.scripts.list.filter.incScriptsWithSecurity', {
                  count: totalScriptsWithSecurity,
                })
              "
            />
          </div>
        </div>

        <Divider />

        <CResourceTable
          :items="filtered"
          :fields="scriptFields"
          primary-key="name"
          :empty-message="$t('general.resourceList.noItems')"
        >
          <template #body-name="{ data }">
            <div class="flex flex-col gap-1">
              <div class="flex items-center gap-2 flex-wrap">
                <span v-if="data.label" class="font-medium">{{ data.label }}</span>
                <span v-else class="text-muted-color italic">
                  {{ $t('automation.scripts.list.labelMissing') }}
                </span>

                <Tag
                  v-if="data.security"
                  :value="$t('automation.scripts.list.flags.security')"
                  severity="info"
                  class="text-xs"
                />
                <Tag
                  v-if="data.triggers"
                  :value="$t('automation.scripts.list.flags.triggers')"
                  severity="info"
                  class="text-xs"
                />
                <Tag
                  v-if="data.iterator"
                  :value="$t('automation.scripts.list.flags.iterator')"
                  severity="info"
                  class="text-xs"
                />
              </div>

              <span v-if="data.description" class="text-xs text-muted-color">
                {{ data.description }}
              </span>

              <code class="text-xs text-muted-color">{{ data.name }}</code>

              <div v-if="data.errors && data.errors.length" class="flex flex-col gap-1 mt-1">
                <Message
                  v-for="(error, i) in data.errors"
                  :key="i"
                  severity="warn"
                  :closable="false"
                  class="text-sm"
                >
                  {{ error }}
                </Message>
              </div>
            </div>
          </template>

          <template #body-updatedAt="{ data }">
            <span v-if="data.updatedAt" class="text-sm text-muted-color">
              {{ formatDate(data.updatedAt) }}
            </span>
          </template>
        </CResourceTable>
      </Panel>
    </CViewContainer>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@planetcrust/human-vue'

const { CResourceTable, CViewContainer } = components
const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const items = ref([])

const filter = reactive({
  query: '',
  incScriptsWithErrors: false,
  incScriptsWithTriggers: false,
  incScriptsWithIterator: false,
  incScriptsWithSecurity: false,
})

const scriptFields = [
  { key: 'name', header: t('automation.scripts.list.columns.name') },
  {
    key: 'updatedAt',
    header: t('automation.scripts.list.columns.updatedAt'),
    headerStyle: 'width: 12rem',
    headerClass: 'text-right',
    bodyClass: 'text-right',
  },
]

const filtered = computed(() => {
  const lcQuery = filter.query.toLocaleLowerCase()
  return items.value
    .filter(
      ({ name, label }) =>
        lcQuery.length === 0 || (name + ' ' + (label || '')).toLocaleLowerCase().includes(lcQuery),
    )
    .filter(({ errors }) => !filter.incScriptsWithErrors || (errors && errors.length > 0))
    .filter(({ triggers }) => !filter.incScriptsWithTriggers || !!triggers)
    .filter(({ iterator }) => !filter.incScriptsWithIterator || !!iterator)
    .filter(({ security }) => !filter.incScriptsWithSecurity || !!security)
})

const totalScriptsWithErrors = computed(
  () => items.value.filter(({ errors }) => errors && errors.length > 0).length,
)
const totalScriptsWithTriggers = computed(
  () => items.value.filter(({ triggers }) => !!triggers).length,
)
const totalScriptsWithIterator = computed(
  () => items.value.filter(({ iterator }) => !!iterator).length,
)
const totalScriptsWithSecurity = computed(
  () => items.value.filter(({ security }) => !!security).length,
)

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })
}

async function loadScripts() {
  loading.value = true
  try {
    const result = await $SystemAPI.automationList({})
    const set = Array.isArray(result) ? result : result?.set || []
    items.value = Array.isArray(set) ? set : []
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.fetch.error'))(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => loadScripts())
</script>
