<template>
  <c311-app-shell
    mode="public"
    :brand="t('test.brand', 'City 311')"
    :title="t('test.title', 'C311 interaction test')"
  >
    <section aria-labelledby="c311-interaction-description">
      <p id="c311-interaction-description">{{ t('test.description', 'Shared interaction controls') }}</p>
      <button ref="opener" type="button" class="btn btn-primary" data-c311-modal-opener @click="modalOpen = true">
        {{ t('test.openDialog', 'Open dialog') }}
      </button>
    </section>
    <c311-focus-modal
      v-model="modalOpen"
      :title="t('test.dialogTitle', 'C311 dialog')"
      :close-label="t('action.close', 'Close')"
    >
      <p>{{ t('test.dialogBody', 'Dialog content') }}</p>
      <button type="button" class="btn btn-link">{{ t('test.secondaryAction', 'Secondary action') }}</button>
    </c311-focus-modal>
  </c311-app-shell>
</template>

<script>
import C311AppShell from './C311AppShell.vue'
import C311FocusModal from './C311FocusModal.vue'

export default {
  name: 'C311InteractionHarness',
  components: {
    C311AppShell,
    C311FocusModal,
  },
  data: () => ({ modalOpen: false }),
  methods: {
    t (key, fallback) {
      const translated = this.$t?.(`c311:${key}`)
      return translated && translated !== `c311:${key}` && translated !== key ? translated : fallback
    },
  },
}
</script>
