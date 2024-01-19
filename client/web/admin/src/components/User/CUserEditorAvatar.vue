<template>
  <b-card
    data-test-id="card-user-profile-avatar"
    header-bg-variant="white"
    footer-bg-variant="white"
    header-class="border-bottom"
    footer-class="d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <template #header>
      <h4 class="m-0">
        {{ $t('title') }}
      </h4>
    </template>

    <b-form
      enctype="multipart/form-data"
      @submit.prevent="$emit('submit', user)"
    >
      <img
        :src="uploadedAvatar('avatar')"
        style="height: 4rem; width: 4rem;"
        class="rounded-circle mb-4"
      >

      <div
        class="d-flex align-items-center"
      >
        <c-uploader-with-preview
          :endpoint="`/users/${user.userID}/avatar`"
          :labels="$t('uploader', { returnObjects: true })"
          @upload="$emit('onUpload')"
          @clear="$emit('resetAttachment', 'avatar')"
        />

        <c-input-confirm
          v-if="uploadedAvatar('avatar') && isKindAvatar"
          :processing="processingAvatar"
          :text="$t('general:label.delete')"
          size="lg"
          size-confirm="lg"
          variant="danger"
          class="ml-2 h-100"
          @confirmed="$emit('resetAttachment', 'avatar')"
        />
      </div>

      <b-row class="mt-3">
        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('initial.color')"
            label-class="text-primary"
          >
            <c-input-color-picker
              v-model="user.meta.avatarColor"
              data-test-id="input-text-color"
              :translations="{
                modalTitle: $t('colorPicker'),
                light: $t('ui.settings:editor.corteza-studio.tabs.light'),
                dark: $t('ui.settings:editor.corteza-studio.tabs.dark'),
                cancelBtnLabel: $t('general:label.cancel'),
                saveBtnLabel: $t('general:label.saveAndClose')
              }"
              :color-tooltips="colorSchemeTooltips"
              :swatchers="themeColors"
              :swatcher-labels="swatcherLabels"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('initial.backgroundColor')"
            label-class="text-primary"
          >
            <c-input-color-picker
              v-model="user.meta.avatarBgColor"
              data-test-id="input-background-color"
              :translations="{
                modalTitle: $t('colorPicker'),
                light: $t('ui.settings:editor.corteza-studio.tabs.light'),
                dark: $t('ui.settings:editor.corteza-studio.tabs.dark'),
                cancelBtnLabel: $t('general:label.cancel'),
                saveBtnLabel: $t('general:label.saveAndClose')
              }"
              :color-tooltips="colorSchemeTooltips"
              :swatchers="themeColors"
              :swatcher-labels="swatcherLabels"
            />
          </b-form-group>
        </b-col>
      </b-row>
    </b-form>

    <template #footer>
      <c-button-submit
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="$emit('submit', user)"
      />
    </template>
  </b-card>
</template>

<script>
import CUploaderWithPreview from 'corteza-webapp-admin/src/components/CUploaderWithPreview'
import { components } from '@cortezaproject/corteza-vue'
const { CInputColorPicker } = components

export default {
  name: 'CUserEditorAvatar',

  i18nOptions: {
    namespaces: 'system.users',
    keyPrefix: 'editor.avatar',
  },

  components: {
    CUploaderWithPreview,
    CInputColorPicker,
  },

  props: {
    user: {
      type: Object,
      required: true,
    },

    processing: {
      type: Boolean,
    },

    processingAvatar: {
      type: Boolean,
    },

    success: {
      type: Boolean,
    },
  },

  data () {
    return {
      swatcherLabels: [
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
    }
  },

  computed: {
    isKindAvatar () {
      return this.user.meta.avatarKind === 'avatar'
    },

    colorSchemeTooltips () {
      return this.swatcherLabels.reduce((acc, label) => {
        acc[label] = this.$t(`ui.settings:editor.corteza-studio.theme.variables.${label}`)
        return acc
      }, {})
    },

    themeColors () {
      const theme = this.$Settings.get('ui.studio.themes', [])
      if (!theme.length) {
        return theme
      }

      return theme.map(theme => {
        console.log(theme.values)
        theme.values = JSON.parse(theme.values)
        return theme
      })
    },
  },

  methods: {
    uploadedAvatar (name) {
      const attachmentID = this.user.meta.avatarID

      if (attachmentID !== '0') {
        return (
          this.$SystemAPI.baseURL +
            this.$SystemAPI.attachmentOriginalEndpoint({
              attachmentID,
              kind: 'avatar',
              name,
            })
        )
      }
    },
  },
}
</script>
