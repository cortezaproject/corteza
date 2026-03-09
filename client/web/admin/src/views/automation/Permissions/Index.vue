<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('permissions.ui.title.automation') }}</span>
  </Teleport>

  <CPermissionGrid :api="$AutomationAPI" component="automation" />
</template>

<script setup>
import { inject } from 'vue'
import CPermissionGrid from '../../../components/Permissions/CPermissionGrid.vue'

const $AutomationAPI = inject('$AutomationAPI')
</script>
