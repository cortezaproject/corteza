<template>
  <!-- A Build/Govern step reaches 'changes-requested' when a granter flags it
       via the wizard's per-step "Request changes" action (see Wizard.vue) —
       this banner only ever has something to show for that status (approved
       and draft carry no note). Shown to every member viewing the step (see
       Wizard.vue's showStatus), so whoever picks up the work sees exactly
       what needs fixing. -->
  <Message v-if="status === 'changes-requested' && reviewNote" severity="warn" :closable="false">
    <span class="font-medium">{{ $t('project.governance.banner.changesRequestedLabel') }}</span>
    {{ reviewNote }}
  </Message>
</template>

<script setup>
defineProps({
  status: { type: String, default: 'draft' },
  reviewNote: { type: String, default: '' },
})
</script>
