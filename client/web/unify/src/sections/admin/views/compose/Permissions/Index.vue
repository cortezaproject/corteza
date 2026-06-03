<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('permissions.ui.title.compose') }}</span>
  </Teleport>

  <CPermissionGrid :api="$ComposeAPI" component="compose" />
</template>

<script setup>
import { inject } from 'vue'
import CPermissionGrid from '../../../components/Permissions/CPermissionGrid.vue'

const $ComposeAPI = inject('$ComposeAPI')
</script>
