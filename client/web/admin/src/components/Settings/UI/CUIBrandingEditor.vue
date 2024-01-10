<template>
  <b-card
    header-bg-variant="white"
    footer-bg-variant="white"
    footer-class="d-flex flex-wrap flex-fill-child gap-1"
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
                <template #footer>
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

        <!-- <b-row>
          <b-col>
            <b-form-group
              :label="$t('custom-css')"
              label-class="text-primary"
            >
              <c-ace-editor
                v-model="theme.values"
                lang="css"
                height="300px"
                font-size="14px"
                show-line-numbers
                :border="false"
                :show-popout="true"
              />
            </b-form-group>
          </b-col>
        </b-row> -->
      </b-tab>
    </b-tabs>

    <!-- <b-modal
      id="custom-css-editor"
      v-model="theme.showEditorModal"
      :title="$t('modal.editor')"
      cancel-variant="link"
      size="lg"
      :ok-title="$t('general:label.saveAndClose')"
      :cancel-title="$t('general:label.cancel')"
      body-class="p-0"
      @ok="saveCustomCSSInput(theme.id)"
      @hidden="resetCustomCSSInput(theme.id)"
    >
      <c-ace-editor
        v-model="theme.modalValue"
        lang="scss"
        height="500px"
        font-size="14px"
        show-line-numbers
        :border="false"
        :show-popout="false"
      />
    </b-modal> -->

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
      themeInputs: [
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
      },
      darkModeVariables: {
        'black': '#FBF7F4',
        'white': '#0B344E',
        'primary': '#FF9661',
        'secondary': '#758D9B',
        'success': '#43AA8B',
        'warning': '#E2A046',
        'danger': '#E54122',
        'light': '#768D9A',
        'extra-light': '#23495F',
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
              values: this.lightModeVariables,
            },
            {
              id: 'dark',
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

    resetColor (key, index, id) {
      this.themes.find(theme => theme.id === id).values[key] = id === 'light' ? this.lightModeVariables[key] : this.darkModeVariables[key]
      this.$refs.picker[index].closeMenu()
    },

  },
}
</script>
