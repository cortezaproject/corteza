<template>
  <b-tab :title="$t('kind.file.label')">
    <b-form-group
      horizontal
      :description="$t('kind.file.view.modeFootnote')"
      :label="$t('kind.file.view.modeLabel')"
      label-class="text-primary"
    >
      <b-form-radio-group
        v-model="options.mode"
        buttons
        button-variant="outline-secondary"
        size="sm"
        name="buttons2"
        :options="modes"
      />
    </b-form-group>

    <b-form-group>
      <b-form-checkbox
        v-if="enablePreviewStyling"
        v-model="options.hideFileName"
      >
        {{ $t('kind.file.view.showName') }}
      </b-form-checkbox>

      <b-form-checkbox
        v-model="options.clickToView"
      >
        {{ $t('kind.file.view.clickToView') }}
      </b-form-checkbox>

      <b-form-checkbox
        v-model="options.enableDownload"
      >
        {{ $t('kind.file.view.enableDownload') }}
      </b-form-checkbox>
    </b-form-group>

    <div class="d-flex gap-1">
      <uploader
        ref="uploader"
        :endpoint="endpoint"
        :max-filesize="$s('compose.Page.Attachments.MaxSize', 100)"
        :accepted-files="$s('compose.Page.Attachments.Mimetypes', ['*/*'])"
        class="flex-grow-1"
        @uploaded="appendAttachment"
      />

      <c-webcam
        :labels="{
          tooltip: $t('general:webcam.tooltip'),
          modalTitle: $t('general:webcam.title'),
          cancelButtonLabel: $t('general:webcam.buttons.cancel'),
          confirmButtonLabel: $t('general:webcam.buttons.confirm'),
          captureButtonLabel: $t('general:webcam.buttons.capture'),
          cameraErrorMessage: $t('general:webcam.errors.camera')
        }"
        @upload="uploadWebcamImage"
      >
        <font-awesome-icon
          class="text-primary"
          :icon="['fas', 'camera']"
        />
      </c-webcam>
    </div>

    <list-loader
      kind="page"
      enable-delete
      :namespace="namespace"
      :set.sync="options.attachments"
      mode="list"
      class="mt-2"
    />

    <template v-if="enablePreviewStyling">
      <hr>

      <h5 class="mb-2">
        {{ $t('kind.file.view.previewStyle') }}
      </h5>

      <small>{{ $t('kind.file.view.description' ) }}</small>

      <b-row
        align-v="center"
        class="mb-2 mt-2"
      >
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('kind.file.view.height')"
            label-class="text-primary"
          >
            <b-input-group>
              <b-form-input
                v-model="options.height"
              />
            </b-input-group>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('kind.file.view.width')"
            label-class="text-primary"
          >
            <b-input-group>
              <b-form-input
                v-model="options.width"
              />
            </b-input-group>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('kind.file.view.maxHeight')"
            label-class="text-primary"
          >
            <b-input-group>
              <b-form-input
                v-model="options.maxHeight"
              />
            </b-input-group>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('kind.file.view.maxWidth')"
            label-class="text-primary"
          >
            <b-input-group>
              <b-form-input
                v-model="options.maxWidth"
              />
            </b-input-group>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('kind.file.view.borderRadius')"
            label-class="text-primary"
          >
            <b-input-group>
              <b-form-input
                v-model="options.borderRadius"
              />
            </b-input-group>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('kind.file.view.margin')"
            label-class="text-primary"
          >
            <b-input-group>
              <b-form-input
                v-model="options.margin"
              />
            </b-input-group>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('kind.file.view.background')"
            label-class="text-primary"
          >
            <c-input-color-picker
              v-model="options.backgroundColor"
              :translations="{
                modalTitle: $t('kind.file.view.colorPicker'),
                light: $t('general:themes.labels.light'),
                dark: $t('general:themes.labels.dark'),
                cancelBtnLabel: $t('general:label.cancel'),
                saveBtnLabel: $t('general:label.saveAndClose')
              }"
              :theme-settings="themeSettings"
            />
          </b-form-group>
        </b-col>
      </b-row>
    </template>
  </b-tab>
</template>
<script>
import base from './base'
import Uploader from 'corteza-webapp-compose/src/components/Public/Page/Attachment/Uploader'
import ListLoader from 'corteza-webapp-compose/src/components/Public/Page/Attachment/ListLoader'
import { components } from '@cortezaproject/corteza-vue'
const { CInputColorPicker } = components

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  name: 'File',

  components: {
    Uploader,
    ListLoader,
    CInputColorPicker,
  },

  extends: base,

  computed: {
    endpoint () {
      const { pageID } = this.page

      return this.$ComposeAPI.pageUploadEndpoint({
        namespaceID: this.namespace.namespaceID,
        pageID,
      })
    },

    modes () {
      return [
        { value: 'list', text: this.$t('kind.file.view.list') },
        { value: 'gallery', text: this.$t('kind.file.view.gallery') },
      ]
    },

    enablePreviewStyling () {
      const { mode } = this.options
      return mode === 'gallery'
    },

    themeSettings () {
      return this.$Settings.get('ui.studio.themes', [])
    },
  },

  methods: {
    appendAttachment ({ attachmentID } = {}) {
      this.options.attachments.push(attachmentID)
    },

    uploadWebcamImage (file) {
      const uploader = this.$refs.uploader
      uploader.$refs.dropzone.addFile(file)
    },
  },
}
</script>
