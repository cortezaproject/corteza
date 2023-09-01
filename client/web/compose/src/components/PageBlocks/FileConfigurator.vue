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

    <uploader
      :endpoint="endpoint"
      :max-filesize="$s('compose.Page.Attachments.MaxSize', 100)"
      :accepted-files="$s('compose.Page.Attachments.Mimetypes', ['*/*'])"
      @uploaded="appendAttachment"
    />
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
          sm="12"
          md="6"
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
          sm="12"
          md="6"
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
          sm="12"
          md="6"
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
          sm="12"
          md="6"
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
          sm="12"
          md="6"
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
          sm="12"
          md="6"
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
          sm="12"
          md="6"
        >
          <b-form-group
            :label="$t('kind.file.view.background')"
            label-class="text-primary"
          >
            <c-input-color-picker
              v-model="options.backgroundColor"
              :translations="{
                modalTitle: $t('kind.file.view.colorPicker'),
                cancelBtnLabel: $t('general:label.cancel'),
                saveBtnLabel: $t('general:label.saveAndClose')
              }"
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
  },

  methods: {
    appendAttachment ({ attachmentID } = {}) {
      this.options.attachments.push(attachmentID)
    },
  },
}
</script>
