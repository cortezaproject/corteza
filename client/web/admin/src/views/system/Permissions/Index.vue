<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('permissions.ui.title.system') }}</span>
  </Teleport>

  <CPermissionGrid :api="$SystemAPI" component="system" />
</template>

<script setup>
import { inject } from 'vue'
import CPermissionGrid from '../../../components/Permissions/CPermissionGrid.vue'

const $SystemAPI = inject('$SystemAPI')
</script>
