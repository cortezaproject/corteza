<template>
  <b-card
    data-test-id="card-edit-authentication"
    header-bg-variant="white"
    footer-bg-variant="white"
    class="shadow-sm"
  >
    <template #header>
      <h3 class="m-0">
        {{ $t("title") }}
      </h3>
    </template>

    <b-row cols="12">
      <b-col cols="6">
        <div class="shadow-sm">
          <div class="d-flex justify-content-between">
            <h5>
              {{ $t("image.uploader.title") }}
            </h5>

            <b-button
              v-if="uploadedFile('auth.ui.background-image-src')"
              variant="link"
              class="d-flex align-items-top text-dark p-1"
              @click="$emit('resetAttachment', 'auth.ui.background-image-src')"
            >
              <font-awesome-icon :icon="['far', 'trash-alt']" />
            </b-button>
          </div>

          <c-uploader-with-preview
            :value="uploadedFile('auth.ui.background-image-src')"
            :endpoint="'/settings/auth.ui.background-image-src'"
            :disabled="!canManage"
            :labels="$t('image.uploader', { returnObjects: true })"
            @upload="$emit('onUpload')"
            @clear="$emit('resetAttachment', 'auth.ui.background-image-src')"
          />
        </div>
      </b-col>

      <b-col cols="6">
        <div class="shadow-sm">
          <h5>
            {{ $t("image.editor.title") }}
          </h5>

          <ace-editor
            data-test-id="auth-bg-image-styling-editor"
            :font-size="14"
            :show-print-margin="true"
            :show-gutter="true"
            :highlight-active-line="true"
            width="100%"
            height="200px"
            mode="css"
            theme="chrome"
            name="editor/css"
            :on-change="v => (settings['auth.ui.styles'] = v)"
            :value="settings['auth.ui.styles']"
            :editor-props="{
              $blockScrolling: false
            }"
          />

          <c-submit-button
            :disabled="!canManage"
            :processing="processing"
            :success="success"
            class="float-right mt-2"
            @submit="$emit('submit', settings['auth.ui.styles'])"
          />
        </div>
      </b-col>
    </b-row>
  </b-card>
</template>

<script>
import { Ace as AceEditor } from 'vue2-brace-editor'
import CUploaderWithPreview from 'corteza-webapp-admin/src/components/CUploaderWithPreview'
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'

import 'brace/mode/css'
import 'brace/theme/chrome'

export default {
  name: 'CSystemEditorAuthBgImage',

  i18nOptions: {
    namespaces: 'system.settings',
    keyPrefix: 'editor.bgScreen',
  },

  components: {
    CUploaderWithPreview,
    AceEditor,
    CSubmitButton,
  },

  props: {
    settings: {
      type: Object,
      required: true,
    },

    canManage: {
      type: Boolean,
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
  },

  methods: {
    uploadedFile (name) {
      const localAttachment = /^attachment:(\d+)/

      switch (true) {
        case this.settings[name] && localAttachment.test(this.settings[name]):
          const [, attachmentID] = localAttachment.exec(this.settings[name])

          return (
            this.$SystemAPI.baseURL +
            this.$SystemAPI.attachmentOriginalEndpoint({
              attachmentID,
              kind: 'settings',
              name,
            })
          )
      }

      return undefined
    },
  },
}
</script>
