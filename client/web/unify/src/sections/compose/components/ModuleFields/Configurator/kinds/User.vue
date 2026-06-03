<template>
  <div class="flex flex-col gap-6">
    <!-- Preset with current user -->
    <div class="flex items-center gap-2">
      <Checkbox v-model="field.options.presetWithAuthenticated" inputId="presetWithAuth" :binary="true" />
      <label for="presetWithAuth" class="cursor-pointer">{{ $t('field.kind.user.presetWithCurrentUser') }}</label>
    </div>

    <!-- Role filter -->
    <div class="flex flex-col gap-2">
      <label class="font-medium text-muted-color text-sm">{{ $t('field.kind.user.roles.label') }}</label>
      <MultiSelect
        v-model="selectedRoles"
        :options="roleOptions"
        option-label="name"
        option-value="roleID"
        :placeholder="$t('field.kind.user.roles.placeholder')"
        :loading="loadingRoles"
        filter
        class="w-full"
        @update:model-value="field.options.roles = $event"
      />
    </div>

    <!-- Multi-value select type -->
    <template v-if="field.isMulti">
      <div>
        <label class="font-medium text-muted-color text-sm block mb-3">
          {{ $t('field.kind.select.optionType.label') }}
        </label>
        <div class="flex flex-col gap-2">
          <div
            v-for="opt in selectTypeOptions"
            :key="opt.value"
            class="flex items-center gap-2"
          >
            <RadioButton
              :input-id="`userSelectType-${opt.value}`"
              v-model="field.options.selectType"
              :value="opt.value"
              @update:model-value="onSelectTypeChange"
            />
            <label :for="`userSelectType-${opt.value}`" class="cursor-pointer">{{ opt.label }}</label>
          </div>
        </div>
      </div>

      <div v-if="showAllowDuplicates" class="flex items-center gap-2">
        <Checkbox
          :model-value="!field.options.isUniqueMultiValue"
          inputId="userAllowDuplicates"
          :binary="true"
          @update:model-value="field.options.isUniqueMultiValue = !$event"
        />
        <label for="userAllowDuplicates" class="cursor-pointer">{{ $t('field.kind.select.allow-duplicates') }}</label>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const field = inject('fieldDraft')

const $SystemAPI = inject('$SystemAPI')

const roleOptions = ref([])
const loadingRoles = ref(false)

// Roles where duplicates can be allowed
const duplicatesAllowedTypes = ['default', 'each']
const showAllowDuplicates = computed(() =>
  duplicatesAllowedTypes.includes(field.value.options.selectType),
)

const selectTypeOptions = computed(() => [
  { value: 'default', label: t('field.kind.select.optionType.default') },
  { value: 'multiple', label: t('field.kind.select.optionType.multiple') },
  { value: 'each', label: t('field.kind.select.optionType.each') },
])

// selectedRoles is initialized from field.options.roles (array of roleIDs)
const selectedRoles = computed({
  get: () => field.value.options.roles || [],
  set: val => { field.value.options.roles = val },
})

function onSelectTypeChange(val) {
  const allowDuplicates = duplicatesAllowedTypes.includes(val)
  if (!allowDuplicates) {
    field.value.options.isUniqueMultiValue = true
  }
}

onMounted(async () => {
  if (!$SystemAPI) return
  loadingRoles.value = true
  try {
    const { set } = await $SystemAPI.roleList({ limit: 200, sort: 'name ASC' })
    roleOptions.value = set || []
  } catch {
    // ignore
  } finally {
    loadingRoles.value = false
  }
})
</script>
