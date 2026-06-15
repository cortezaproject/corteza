<template>
  <Message v-if="status === 'changes-requested' && reviewNote" severity="warn" :closable="false">
    <span class="font-medium">{{ $t('project.governance.banner.changesRequestedLabel') }}</span> {{ reviewNote }}
  </Message>
  <Message v-else-if="status === 'submitted'" severity="info" :closable="false">
    {{ $t('project.governance.banner.submitted') }}
  </Message>
  <Message v-else-if="status === 'approved'" severity="success" :closable="false">
    {{ $t('project.governance.banner.approved') }}
  </Message>
  <Message v-else-if="status === 'draft' && reviewNote" severity="warn" :closable="false">
    <span class="font-medium">{{ $t('project.governance.banner.reopenedLabel') }}</span> {{ reviewNote }}
  </Message>
</template>

<script setup>
defineProps({
  status: { type: String, default: 'draft' },
  reviewNote: { type: String, default: '' },
})
</script>
