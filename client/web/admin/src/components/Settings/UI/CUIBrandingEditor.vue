<template>
  <b-card
    class="shadow-sm"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>

    <b-row class="mx-2">
      <b-col
        v-for="key in Object.keys(brandingVariables)"
        :key="key"
        cols="6"
        class="p-2"
      >
        <b-form-group
          :label="$t(`brandVariables.${key}`)"
          class="row px-3"
        >
          <c-input-color-picker
            v-model="brandingVariables[key]"
            :data-test-id="`input-${key}-color`"
            width="24px"
            height="24px"
            :show-color-code-text="true"
            :translations="colorTranslations"
          />
        </b-form-group>
      </b-col>
    </b-row>

    <template #footer>
      <c-submit-button
        :disabled="!canManage"
        :processing="processing"
        :success="success"
        class="float-right mt-2"
        @submit="onSubmit"
      />
    </template>
  </b-card>
</template>

<script>
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'
import { components } from '@cortezaproject/corteza-vue'
const { CInputColorPicker } = components

export default {
  name: 'CUIBrandingEditor',

  i18nOptions: {
    namespaces: 'ui.settings',
    keyPrefix: 'editor.corteza-studio',
  },

  components: {
    CInputColorPicker,
    CSubmitButton,
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
      brandingVariables: {
        white: '#FFFFFF',
        black: '#162425',
        primary: '#0B344E',
        secondary: '#758D9B',
        success: '#43AA8B',
        warning: '#E2A046',
        danger: '#E54122',
        light: '#E4E9EF',
        extraLight: '#F3F5F7',
        dark: '#162425',
        tertiary: '#5E727E',
        gray200: '#F9FAFB',
        bodyBg: '#F9FAFB',
      },
      colorTranslations: {
        modalTitle: this.$t('colorPicker'),
        saveBtnLabel: this.$t('general:label.saveAndClose'),
      },
    }
  },

  watch: {
    settings: {
      immediate: true,
      handler (settings) {
        if (settings['ui.branding-sass']) {
          this.brandingVariables = JSON.parse(settings['ui.branding-sass'])
        }
      },
    },
  },

  methods: {
    onSubmit () {
      this.$emit('submit', { 'ui.branding-sass': JSON.stringify(this.brandingVariables) })
    },
  },
}
</script>
