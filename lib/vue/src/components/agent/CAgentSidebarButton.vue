<template>
  <div class="relative" v-if="agentStore.availableAgents.length > 0">
    <Button
      v-tooltip.bottom="tooltipText"
      icon="pi pi-sparkles"
      severity="secondary"
      variant="text"
      rounded
      @click="agentStore.toggleVisibility()"
    />
  </div>
</template>

<script setup lang="ts">
import { inject, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAgentSidebarStore } from '../../stores/useAgentSidebarStore'

const { t } = useI18n()
const tooltipText = computed(() => t('agent.sidebar.title'))

const agentStore = useAgentSidebarStore()
const $SystemAPI = inject<any>('$SystemAPI')
const $Auth = inject<any>('$Auth')

onMounted(async () => {
  if (agentStore.availableAgents.length > 0) return // Already loaded

  try {
    const res = await $SystemAPI.agentList({ limit: 0 })
    const allAgents = res.set || []
    const userRoles = $Auth.user?.roles || []
    const configuredAgents = allAgents.filter((agent: any) => {
      // Must have user invocation enabled
      if (!agent.invocation?.user?.enabled) return false

      const sidebarRoles = agent.meta?.sidebarRoles || []
      if (!Array.isArray(sidebarRoles) || sidebarRoles.length === 0) return false
      
      // Superadmin bypass (role ID 2) or explicit role check
      return userRoles.includes('2') || sidebarRoles.some((r: string) => userRoles.includes(r))
    })
    
    agentStore.setAvailableAgents(configuredAgents)
  } catch (err) {
    console.warn('Failed to load agents for sidebar', err)
  }
})
</script>
