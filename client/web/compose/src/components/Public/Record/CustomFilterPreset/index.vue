<template>
  <div>
    <b-modal
      :visible="showModal"
      :title="$t('recordList.filterPresets.saveFilterAsPreset')"
      body-class="p-0"
      footer-class="d-flex w-100 align-items-center justify-content-between"
      centered
      no-fade
      @hide="onModalHide"
    >
      <b-card
        class="pt-0"
      >
        <b-form-group
          :label="$t('recordList.filterPresets.filterName')"
          label-class="primary"
        >
          <b-form-input
            v-model="filterName"
          />
        </b-form-group>
      </b-card>

      <template #modal-footer>
        <b-button
          variant="light"
          @click="onModalHide"
        >
          {{ $t('general.label.cancel') }}
        </b-button>

        <div>
          <b-button
            variant="primary"
            :disabled="!filterName"
            @click="onSave"
          >
            {{ $t('general.label.save') }}
          </b-button>
        </div>
      </template>
    </b-modal>
  </div>
</template>

<script>
export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'CustomFilterPreset',

  props: {
    visible: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      showModal: false,
      filterName: '',
    }
  },

  watch: {
    visible: {
      immediate: true,
      handler (val) {
        this.showModal = val
      },
    },
  },

  methods: {
    onModalHide () {
      this.showModal = false
      this.filterName = ''
      this.$emit('close')
    },

    onSave () {
      this.$emit('save', {
        name: this.filterName,
      })

      this.onModalHide()
    },
  },

}
</script>

<style lang="scss">
.position-initial {
  position: initial;
}
</style>
