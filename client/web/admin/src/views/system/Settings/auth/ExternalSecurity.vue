<template>
  <div class="flex flex-col gap-4">
    <Divider />
    <h5 class="text-sm font-semibold">
      {{ $t('system.settings.editor.external.security.title') }}
    </h5>

    <div class="flex flex-col gap-1">
      <label class="font-medium text-sm">
        {{ $t('system.settings.editor.external.security.permitted-roles.label') }}
      </label>
      <span class="text-xs text-muted-color">
        {{ $t('system.settings.editor.external.security.permitted-roles.description') }}
      </span>
      <CInputRole
        :multiple="true"
        :model-value="modelValue.permittedRoles || []"
        filter-context-roles
        @update:model-value="updateList('permittedRoles', $event)"
      />
    </div>

    <div class="flex flex-col gap-1">
      <label class="font-medium text-sm">
        {{ $t('system.settings.editor.external.security.prohibited-roles.label') }}
      </label>
      <span class="text-xs text-muted-color">
        {{ $t('system.settings.editor.external.security.prohibited-roles.description') }}
      </span>
      <CInputRole
        :multiple="true"
        :model-value="modelValue.prohibitedRoles || []"
        filter-context-roles
        @update:model-value="updateList('prohibitedRoles', $event)"
      />
    </div>

    <div class="flex flex-col gap-1">
      <label class="font-medium text-sm">
        {{ $t('system.settings.editor.external.security.forced-roles.label') }}
      </label>
      <span class="text-xs text-muted-color">
        {{ $t('system.settings.editor.external.security.forced-roles.description') }}
      </span>
      <CInputRole
        :multiple="true"
        :model-value="modelValue.forcedRoles || []"
        filter-context-roles
        @update:model-value="updateList('forcedRoles', $event)"
      />
    </div>
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
