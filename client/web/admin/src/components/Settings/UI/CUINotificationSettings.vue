<template>
  <b-card
    header-class="border-bottom"
    footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <template #header>
      <h4 class="m-0">
        {{ $t('title') }}
      </h4>
    </template>

    <b-form
      @submit.prevent="$emit('submit', settings)"
    >
      <b-form-group
        :label="$t('toastPosition.label')"
        :description="$t('toastPosition.description')"
        label-class="text-primary"
      >
        <c-input-select
          v-model="notificationSettings.toastPosition"
          :options="positionOptions"
          :reduce="o => o.value"
          label="text"
          :clearable="false"
        />
      </b-form-group>
    </b-form>

    <template #footer>
      <c-button-submit
        v-if="canManage"
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="onSubmit"
      />
    </template>
  </b-card>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'

const { CInputSelect } = components

export default {
  name: 'CUINotificationSettings',

  components: {
    CInputSelect,
  },

  i18nOptions: {
    namespaces: 'ui.settings',
    keyPrefix: 'editor.notifications',
  },

  props: {
    settings: {
      type: Object,
      required: true,
    },

    processing: {
      type: Boolean,
      value: false,
    },

    success: {
      type: Boolean,
      value: false,
    },

    canManage: {
      type: Boolean,
      required: true,
    },
  },

  data () {
    return {
      notificationSettings: {},

      positionOptions: [
        { value: 'b-toaster-top-right', text: this.$t('toastPosition.options.top-right') },
        { value: 'b-toaster-top-left', text: this.$t('toastPosition.options.top-left') },
        { value: 'b-toaster-top-center', text: this.$t('toastPosition.options.top-center') },
        { value: 'b-toaster-top-full', text: this.$t('toastPosition.options.top-full') },
        { value: 'b-toaster-bottom-right', text: this.$t('toastPosition.options.bottom-right') },
        { value: 'b-toaster-bottom-left', text: this.$t('toastPosition.options.bottom-left') },
        { value: 'b-toaster-bottom-center', text: this.$t('toastPosition.options.bottom-center') },
        { value: 'b-toaster-bottom-full', text: this.$t('toastPosition.options.bottom-full') },
      ],
    }
  },

  watch: {
    settings: {
      immediate: true,
      handler (settings) {
        this.notificationSettings = settings['ui.notifications'] || {}

        if (!this.notificationSettings.toastPosition) {
          this.$set(this.notificationSettings, 'toastPosition', 'b-toaster-bottom-right')
        }
      },
    },
  },

  methods: {
    onSubmit () {
      this.$emit('submit', { 'ui.notifications': this.notificationSettings })
    },
  },
}
</script>
