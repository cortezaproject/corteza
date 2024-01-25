<template>
  <b-card
    footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
    body-class="p-0"
    class="shadow-sm"
  >
    <template #header>
      <h4 class="m-0">
        {{ $t('title') }}
      </h4>
    </template>

    <div
      v-if="!sassInstalled"
      class="bg-warning rounded p-2 mb-3"
    >
      {{ $t('sassNotInstalled') }}
      <a :href="installSassDocs">{{ $t('installSassDocs') }}</a>
    </div>

    <b-tabs
      data-test-id="theme-tabs"
      card
    >
      <b-tab
        v-for="theme in themes"
        :key="theme.id"
        :title="theme.title"
      >
        <b-row v-if="theme.id !== 'general'">
          <b-col
            v-for="key in themeVariables"
            :key="key"
            md="6"
            cols="12"
          >
            <b-form-group
              :label="$t(`theme.variables.${key}`)"
              label-class="text-primary"
            >
              <c-input-color-picker
                ref="picker"
                v-model="theme.variables[key]"
                :default-value="theme.defaultVariables[key]"
                :data-test-id="`input-${key}-color`"
                :translations="{
                  modalTitle: $t('colorPicker'),
                  defaultBtnLabel: $t('label.default'),
                  light: $t('tabs.light'),
                  dark: $t('tabs.dark'),
                  cancelBtnLabel: $t('general:label.cancel'),
                  saveBtnLabel: $t('general:label.saveAndClose')
                }"
                :theme-settings="settings['ui.studio.themes']"
              />
            </b-form-group>
          </b-col>
        </b-row>
        <b-row>
          <b-col>
            <b-form-group
              :label="$t('custom-css')"
              label-class="text-primary"
              class="mb-0"
            >
              <c-ace-editor
                v-model="theme.customCSS"
                lang="css"
                height="400px"
                font-size="14px"
                show-line-numbers
                :show-popout="true"
                @open="openCustomCSSModal(theme.id)"
              />
            </b-form-group>
          </b-col>
        </b-row>
      </b-tab>
    </b-tabs>

    <b-modal
      id="custom-css-editor"
      v-model="customCSSModal.show"
      :title="$t('custom-css')"
      cancel-variant="link"
      size="lg"
      :ok-title="$t('general:label.saveAndClose')"
      :cancel-title="$t('general:label.cancel')"
      body-class="p-0"
      @ok="saveCustomCSSModal()"
      @hidden="resetCustomCSSModal()"
    >
      <c-ace-editor
        v-model="customCSSModal.value"
        lang="scss"
        height="500px"
        font-size="14px"
        show-line-numbers
        :border="false"
        :show-popout="false"
      />
    </b-modal>

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
const { CInputColorPicker, CAceEditor } = components

export default {
  name: 'CUIBrandingEditor',

  i18nOptions: {
    namespaces: 'ui.settings',
    keyPrefix: 'editor.corteza-studio',
  },

  components: {
    CInputColorPicker,
    CAceEditor,
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
      themeTabs: [
        'general',
        'light',
        'dark',
      ],
      themeVariables: [
        'black',
        'white',
        'primary',
        'secondary',
        'success',
        'warning',
        'danger',
        'light',
        'extra-light',
        'body-bg',
        'sidebar-bg',
        'topbar-bg',
      ],
      lightModeVariables: {
        'black': '#162425',
        'white': '#FFFFFF',
        'primary': '#0B344E',
        'secondary': '#758D9B',
        'success': '#43AA8B',
        'warning': '#E2A046',
        'danger': '#E54122',
        'light': '#F3F5F7',
        'extra-light': '#E4E9EF',
        'body-bg': '#F3F5F7',
        'sidebar-bg': '#FFFFFF',
        'topbar-bg': '#F3F5F7',
      },
      darkModeVariables: {
        'black': '#FBF7F4',
        'white': '#0B344E',
        'primary': '#FF9661',
        'secondary': '#758D9B',
        'success': '#43AA8B',
        'warning': '#E2A046',
        'danger': '#E54122',
        'light': '#23495F',
        'extra-light': '#768D9A',
        'body-bg': '#092B40',
        'sidebar-bg': '#0B344E',
        'topbar-bg': '#092B40',
      },

      themes: [],
      themeColors: [],

      customCSSModal: {
        show: false,
        id: '',
        value: '',
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
        const themes = settings['ui.studio.themes'] || []
        const customCSS = settings['ui.studio.custom-css'] || []

        this.themes = this.themeTabs.map((id) => {
          const { title, values = '' } = themes.find(t => t.id === id) || {}
          const defaultCustomCSS = customCSS.find(t => t.id === id) || {}

          let variables = JSON.parse(values || '{}')
          let defaultVariables

          if (['light', 'dark'].includes(id)) {
            if (!values) {
              variables = id === 'light' ? this.lightModeVariables : this.darkModeVariables
            }

            defaultVariables = id === 'light' ? this.lightModeVariables : this.darkModeVariables
          }

          return {
            id: id,
            title: title || this.$t(`tabs.${id}`),
            variables,
            defaultVariables,
            customCSS: defaultCustomCSS.values || '',
          }
        })
      },
    },
  },

  methods: {
    onSubmit () {
      this.$emit('submit', {
        'ui.studio.themes': this.themes.map(theme => {
          return {
            id: theme.id,
            title: theme.title,
            values: JSON.stringify(theme.variables),
          }
        }),
        'ui.studio.custom-css': this.themes.map(theme => {
          return {
            id: theme.id,
            title: theme.title,
            values: theme.customCSS,
          }
        }),
      })
    },

    openCustomCSSModal (id) {
      const { customCSS } = this.themes.find(t => t.id === id) || {}

      this.customCSSModal.id = id
      this.customCSSModal.value = customCSS
      this.customCSSModal.show = true
    },

    saveCustomCSSModal () {
      this.themes.find(t => t.id === this.customCSSModal.id).customCSS = this.customCSSModal.value
    },

    resetCustomCSSModal () {
      this.customCSSModal.id = ''
      this.customCSSModal.value = ''
      this.customCSSModal.show = false
    },

    resetColor (key, theme) {
      this.$set(theme.variables, key, theme.id === 'light' ? this.lightModeVariables[key] : this.darkModeVariables[key])
    },

  },
}
</script>
