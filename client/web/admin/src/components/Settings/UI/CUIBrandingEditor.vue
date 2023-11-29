<template>
  <b-card
    header-bg-variant="white"
    footer-bg-variant="white"
    footer-class="d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
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

    <b-tabs
      data-test-id="theme-tabs"
      nav-wrapper-class="bg-white white border-bottom rounded-0"
      card
    >
      <b-tab
        v-for="theme in themes"
        :key="theme.id"
        :title="$t(`tabs.${theme.id}`)"
      >
        <b-row>
          <b-col
            v-for="(key, index) in themeInputs"
            :key="key"
            md="6"
            cols="12"
          >
            <b-form-group
              :label="$t(`theme.values.${key}`)"
              label-class="text-primary"
            >
              <c-input-color-picker
                ref="picker"
                v-model="theme.values[key]"
                :data-test-id="`input-${key}-color`"
                :translations="{
                  modalTitle: $t('colorPicker'),
                  cancelBtnLabel: $t('general:label.cancel'),
                  saveBtnLabel: $t('general:label.saveAndClose')
                }"
              >
                <template v-slot:footer>
                  <b-button
                    variant="outline-primary"
                    class="mr-auto"
                    @click="resetColor(key, index, theme.id)"
                  >
                    {{ $t('label.default') }}
                  </b-button>
                </template>
              </c-input-color-picker>
            </b-form-group>
          </b-col>
        </b-row>
      </b-tab>
    </b-tabs>

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
const { CInputColorPicker } = components

export default {
  name: 'CUIBrandingEditor',

  i18nOptions: {
    namespaces: 'ui.settings',
    keyPrefix: 'editor.corteza-studio',
  },

  components: {
    CInputColorPicker,
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
      themeInputs: [
        'white',
        'black',
        'primary',
        'secondary',
        'success',
        'warning',
        'danger',
        'light',
        'extra-light',
        'dark',
        'tertiary',
        'gray-200',
        'body-bg',
      ],
      lightModeVariables: {
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
      darkModeVariables: {
        'white': '#162425',
        'black': '#FFFFFF',
        'primary': '#43AA8B',
        'secondary': '#E4E9EF',
        'success': '#E2A046',
        'warning': '#758D9B',
        'danger': '#E54122',
        'light': '#5E727E',
        'extra-light': '#F3F5F7',
        'dark': '#0B344E',
        'tertiary': '#F9FAFB',
        'gray-200': '#162425',
        'body-bg': '#162425',
      },
      themes: [],
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
        if (settings['ui.studio.themes']) {
          this.themes = settings['ui.studio.themes'].map(theme => {
            theme.values = JSON.parse(theme.values)
            return theme
          })
        } else {
          this.themes = [
            {
              id: 'light',
              title: this.$t('light'),
              values: this.lightModeVariables,
            },
            {
              id: 'dark',
              title: this.$t('dark'),
              values: this.darkModeVariables,
            },
          ]
        }
      },
    },
  },

  methods: {
    onSubmit () {
      this.$emit('submit', { 'ui.studio.themes': this.themes.map(theme => {
        theme.values = JSON.stringify(theme.values)
        return theme
      }),
      })
    },

    resetColor (key, index, mode) {
      this.themes.forEach(theme => {
        theme.values[key] = mode === 'light' ? this.lightModeVariables[key] : this.darkModeVariables[key]
      })
      this.$refs.picker[index].closeMenu()
    },

  },
}
</script>
