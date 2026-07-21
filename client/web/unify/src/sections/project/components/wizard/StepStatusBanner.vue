<template>
  <!-- The only governance status a non-publish step can reach is
       'changes-requested', raised anytime by a granter via the wizard's
       per-step "Request changes" action (see Wizard.vue) — every other step
       stays in 'draft' forever, carrying no note. Shown to every member
       viewing the step, gated-mode only (see Wizard.vue's showStatus), so
       whoever picks up the work sees exactly what needs fixing. -->
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
