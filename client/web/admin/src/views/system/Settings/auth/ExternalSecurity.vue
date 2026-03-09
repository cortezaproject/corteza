<template>
  <div class="flex flex-col gap-4">
    <Divider />
    <h5 class="text-sm font-semibold">
      {{ $t('system.settings.editor.external.security.title', 'Security') }}
    </h5>

    <div class="flex flex-col gap-1">
      <label class="font-medium text-sm">
        {{
          $t('system.settings.editor.external.security.permitted-roles.label', 'Permitted roles')
        }}
      </label>
      <span class="text-xs text-surface-500">
        {{
          $t(
            'system.settings.editor.external.security.permitted-roles.description',
            'Only roles in this list will be added into security context when authenticates with this provider',
          )
        }}
      </span>
      <CInputRole
        v-for="(roleID, index) in modelValue.permittedRoles"
        :key="'permitted-' + index"
        :model-value="roleID"
        class="mb-1"
        @update:model-value="updateRole('permittedRoles', index, $event)"
      />
      <Button
        :label="$t('general.label.plus-add', '+ Add')"
        severity="secondary"
        text
        size="small"
        class="self-start"
        @click="addRole('permittedRoles')"
      />
    </div>

    <div class="flex flex-col gap-1">
      <label class="font-medium text-sm">
        {{
          $t('system.settings.editor.external.security.prohibited-roles.label', 'Prohibited roles')
        }}
      </label>
      <span class="text-xs text-surface-500">
        {{
          $t(
            'system.settings.editor.external.security.prohibited-roles.description',
            'Roles from this list will be removed from security context when authenticates with this provider',
          )
        }}
      </span>
      <CInputRole
        v-for="(roleID, index) in modelValue.prohibitedRoles"
        :key="'prohibited-' + index"
        :model-value="roleID"
        class="mb-1"
        @update:model-value="updateRole('prohibitedRoles', index, $event)"
      />
      <Button
        :label="$t('general.label.plus-add', '+ Add')"
        severity="secondary"
        text
        size="small"
        class="self-start"
        @click="addRole('prohibitedRoles')"
      />
    </div>

    <div class="flex flex-col gap-1">
      <label class="font-medium text-sm">
        {{ $t('system.settings.editor.external.security.forced-roles.label', 'Forced roles') }}
      </label>
      <span class="text-xs text-surface-500">
        {{
          $t(
            'system.settings.editor.external.security.forced-roles.description',
            'Roles from this list will be always added to security context when authenticates with this provider',
          )
        }}
      </span>
      <CInputRole
        v-for="(roleID, index) in modelValue.forcedRoles"
        :key="'forced-' + index"
        :model-value="roleID"
        class="mb-1"
        @update:model-value="updateRole('forcedRoles', index, $event)"
      />
      <Button
        :label="$t('general.label.plus-add', '+ Add')"
        severity="secondary"
        text
        size="small"
        class="self-start"
        @click="addRole('forcedRoles')"
      />
    </div>
  </div>
</template>

<script setup>
import { CInputRole } from '@cortezaproject/corteza-vue-next'

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

function addRole(listKey) {
  const updated = { ...props.modelValue }
  updated[listKey] = [...(updated[listKey] || []), '']
  emit('update:modelValue', updated)
}

function updateRole(listKey, index, value) {
  const updated = { ...props.modelValue }
  const list = [...(updated[listKey] || [])]
  if (value) {
    list[index] = value
  } else {
    list.splice(index, 1)
  }
  updated[listKey] = list
  emit('update:modelValue', updated)
}
</script>
