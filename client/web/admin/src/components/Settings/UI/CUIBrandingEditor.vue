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

    <div
      v-if="!sassInstalled"
      class="bg-warning rounded p-2 mb-3"
    >
      {{ $t('sassNotInstalled') }}
      <a :href="installSassDocs">{{ $t('installSassDocs') }}</a>
    </div>

    <b-row>
      <b-col
        v-for="key in Object.keys(brandingVariables)"
        :key="key"
        md="6"
        cols="12"
      >
        <b-form-group
          :label="$t(`brandVariables.${key}`)"
          label-class="text-primary"
        >
          <c-input-color-picker
            v-model="brandingVariables[key]"
            :data-test-id="`input-${key}-color`"
            width="32px"
            height="32px"
            show-color-code-text
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
        'white': '#FFFFFF',
        'black': '#162425',
        'primary': '#0B344E',
        'secondary': '#758D9B',
        'success': '#43AA8B',
        'warning': '#E2A046',
        'danger': '#E54122',
        'light': '#E4E9EF',
        'extra-light': '#F3F5F7',
        'dark': '#162425',
        'tertiary': '#5E727E',
        'gray-200': '#F9FAFB',
        'body-bg': '#F9FAFB',
      },
      colorTranslations: {
        modalTitle: this.$t('colorPicker'),
        saveBtnLabel: this.$t('general:label.saveAndClose'),
      },
    }
  },

  computed: {
    sassInstalled () {
      return this.settings['ui.studio.sass-installed']
    },
    installSassDocs () {
      // eslint-disable-next-line no-undef
      const [year, month] = VERSION.split('.')
      return `https://docs.cortezaproject.org/corteza-docs/${year}.${month}/integrator-guide/corteza-studio/index.html`
    },
  },

  watch: {
    settings: {
      immediate: true,
      handler (settings) {
        if (settings['ui.studio.branding-sass']) {
          this.brandingVariables = JSON.parse(settings['ui.studio.branding-sass'])
        }
      },
    },
  },

  methods: {
    onSubmit () {
      this.$emit('submit', { 'ui.studio.branding-sass': JSON.stringify(this.brandingVariables) })
    },
  },
}
</script>
