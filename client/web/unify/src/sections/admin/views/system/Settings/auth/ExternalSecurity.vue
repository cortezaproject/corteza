<template>
  <div class="flex flex-col gap-4">
    <Divider />
    <h5 class="text-sm font-semibold">
      {{ $t('system.settings.editor.external.security.title') }}
    </h5>

    <CFormGroup
      :label="$t('system.settings.editor.external.security.permitted-roles.label')"
      :description="$t('system.settings.editor.external.security.permitted-roles.description')"
    >
      <CInputRole
        :multiple="true"
        :model-value="modelValue.permittedRoles || []"
        filter-context-roles
        @update:model-value="updateList('permittedRoles', $event)"
      />
    </CFormGroup>

    <CFormGroup
      :label="$t('system.settings.editor.external.security.prohibited-roles.label')"
      :description="$t('system.settings.editor.external.security.prohibited-roles.description')"
    >
      <CInputRole
        :multiple="true"
        :model-value="modelValue.prohibitedRoles || []"
        filter-context-roles
        @update:model-value="updateList('prohibitedRoles', $event)"
      />
    </CFormGroup>

    <CFormGroup
      :label="$t('system.settings.editor.external.security.forced-roles.label')"
      :description="$t('system.settings.editor.external.security.forced-roles.description')"
    >
      <CInputRole
        :multiple="true"
        :model-value="modelValue.forcedRoles || []"
        filter-context-roles
        @update:model-value="updateList('forcedRoles', $event)"
      />
    </CFormGroup>
  </div>
</template>

<script setup>
const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({
      permittedRoles: [],
      prohibitedRoles: [],
      forcedRoles: [],
    }),
  },
})

const emit = defineEmits(['update:modelValue'])

function updateList(listKey, value) {
  emit('update:modelValue', { ...props.modelValue, [listKey]: value || [] })
}
</script>
