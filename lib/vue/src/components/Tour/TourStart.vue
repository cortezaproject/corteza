<template>
  <b-modal
    id="tourStart"
    size="lg"
    :title="$t('start.welcome')"
    :ok-title="$t('start.show')"
    :cancel-title="$t('start.skip')"
    body-class="p-0"
    no-fade
    @ok="onSave"
    @cancel="onCancel"
  >
    <div
      class="card-body"
    >
      {{ $t('start.instructions') }}
    </div>
  </b-modal>
</template>

<script>
export default {
  i18nOptions: {
    namespaces: 'onboarding-tour',
  },

  created () {
    this.$nextTick(() => {
      if (window.localStorage.getItem('corteza.tour') === null) {
        this.$bvModal.show('tourStart')
      }
    })
  },

  methods: {
    onSave () {
      this.updateStorage(true)
      this.$emit('start')
      this.showModal = false
    },

    onCancel () {
      this.updateStorage(false)
      this.showModal = false
    },

    updateStorage (status) {
      window.localStorage.setItem('corteza.tour', JSON.stringify(status))
    },
  },
}
</script>
