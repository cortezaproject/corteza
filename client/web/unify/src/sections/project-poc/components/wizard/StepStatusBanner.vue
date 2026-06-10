<template>
  <Message v-if="status === 'changes-requested' && reviewNote" severity="warn" :closable="false">
    <span class="font-medium">Changes requested:</span> {{ reviewNote }}
  </Message>
  <Message v-else-if="status === 'submitted'" severity="info" :closable="false">
    Submitted for review — editing is locked until it's approved or sent back.
  </Message>
  <Message v-else-if="status === 'approved'" severity="success" :closable="false">
    Approved. Reopen from the toolbar to make changes.
  </Message>
  <Message v-else-if="status === 'draft' && reviewNote" severity="warn" :closable="false">
    <span class="font-medium">Reopened:</span> {{ reviewNote }}
  </Message>
</template>

<script setup>
defineProps({
  status: { type: String, default: 'draft' },
  reviewNote: { type: String, default: '' },
})
</script>
