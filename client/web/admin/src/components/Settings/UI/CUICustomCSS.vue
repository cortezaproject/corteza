<template>
  <b-card
    body-class="p-0"
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

    <b-tabs
      data-test-id="theme-tabs"
      nav-wrapper-class="bg-white rounded-0"
      card
    >
      <b-tab
        v-for="theme in themes"
        :key="theme.id"
        :title="$t(`tabs.${theme.id}`)"
      >
        <c-ace-editor
          v-model="theme.values"
          lang="css"
          height="300px"
          font-size="14px"
          show-line-numbers
          :border="false"
          :show-popout="true"
          @open="openEditorModal(theme.id)"
        />

        <b-modal
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
        </b-modal>
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
const { CAceEditor } = components

export default {
  name: 'CUIEditorCustomCSS',

  i18nOptions: {
    namespaces: 'ui.settings',
    keyPrefix: 'editor.custom-css',
  },

  components: {
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
      themes: [
        {
          'id': 'general',
          'title': 'General',
          'values': '',
          'modalValue': undefined,
          showEditorModal: false,
        },
        {
          'id': 'light',
          'title': 'Light mode',
          'values': '',
          'modalValue': undefined,
          showEditorModal: false,
        },
        {
          'id': 'dark',
          'title': 'Dark mode',
          'values': '',
          'modalValue': undefined,
          showEditorModal: false,
        },
      ],
    }
  },

  watch: {
    settings: {
      immediate: true,
      handler (settings) {
        if (settings['ui.studio.custom-css']) {
          this.themes = settings['ui.studio.custom-css'].map((theme) => {
            return {
              id: theme.id,
              title: theme.title,
              values: theme.values,
              modalValue: undefined,
              showEditorModal: false,
            }
          })
        }
      },
    },
  },

  methods: {
    onSubmit () {
      this.$emit('submit', {
        'ui.studio.custom-css': this.themes.map((theme) => {
          return {
            id: theme.id,
            title: theme.title,
            values: theme.values,
          }
        }),
      })
    },

    openEditorModal (themeId) {
      this.themes.forEach((theme) => {
        if (theme.id === themeId) {
          theme.modalValue = theme.values
          theme.showEditorModal = true
        }
      })
    },

    saveCustomCSSInput (themeId) {
      this.themes.forEach((theme) => {
        if (theme.id === themeId) {
          theme.values = theme.modalValue
        }
      })
    },

    resetCustomCSSInput (themeId) {
      this.themes.forEach((theme) => {
        if (theme.id === themeId) {
          theme.modalValue = undefined
        }
      })
    },
  },

}
</script>
