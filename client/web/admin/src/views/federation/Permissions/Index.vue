<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('permissions.ui.title.federation') }}</span>
  </Teleport>

  <CPermissionGrid :api="$FederationAPI" component="federation" />
</template>

<script setup>
import { inject } from 'vue'
import CPermissionGrid from '../../../components/Permissions/CPermissionGrid.vue'

const $FederationAPI = inject('$FederationAPI')
</script>
