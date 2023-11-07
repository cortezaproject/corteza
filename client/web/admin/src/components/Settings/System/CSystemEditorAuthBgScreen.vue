<template>
  <b-card
    header-bg-variant="white"
    footer-bg-variant="white"
    footer-class="d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <template #header>
      <h3 class="m-0">
        {{ $t("title") }}
      </h3>
    </template>

    <b-row cols="12">
      <b-col
        cols="12"
        lg="6"
      >
        <b-form-group
          label-class="d-flex align-items-center text-primary"
        >
          <template #label>
            {{ $t('image.uploader.label') }}

            <c-input-confirm
              v-if="uploadedFile('auth.ui.background-image-src')"
              show-icon
              class="ml-auto"
              @confirmed="$emit('resetAttachment', 'auth.ui.background-image-src')"
            />
          </template>

          <c-uploader-with-preview
            :value="uploadedFile('auth.ui.background-image-src')"
            :endpoint="'/settings/auth.ui.background-image-src'"
            :disabled="!canManage"
            :labels="$t('image.uploader', { returnObjects: true })"
            @upload="$emit('onUpload')"
            @clear="$emit('resetAttachment', 'auth.ui.background-image-src')"
          />
        </b-form-group>
      </b-col>

      <b-col
        cols="12"
        lg="6"
      >
        <b-form-group
          :label="$t('image.editor.label')"
          label-class="d-flex align-items-center text-primary"
        >
          <ace-editor
            data-test-id="auth-bg-image-styling-editor"
            :font-size="14"
            :show-print-margin="true"
            :show-gutter="true"
            :highlight-active-line="true"
            mode="css"
            theme="chrome"
            height="h-100"
            name="editor/css"
            :on-change="v => (settings['auth.ui.styles'] = v)"
            :value="settings['auth.ui.styles']"
            :editor-props="{
              $blockScrolling: false
            }"
            class="flex-fill w-100"
            style="min-height: 300px;"
          />
        </b-form-group>
      </b-col>
    </b-row>

    <template #footer>
      <c-button-submit
        v-if="canManage"
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="$emit('submit', settings['auth.ui.styles'])"
      />
    </template>
  </b-card>
</template>

<script>
import 'brace/mode/css'
import 'brace/theme/chrome'

import { Ace as AceEditor } from 'vue2-brace-editor'
import CUploaderWithPreview from 'corteza-webapp-admin/src/components/CUploaderWithPreview'

export default {
  name: 'CSystemEditorAuthBgImage',

  i18nOptions: {
    namespaces: 'system.settings',
    keyPrefix: 'editor.bgScreen',
  },

  components: {
    CUploaderWithPreview,
    AceEditor,
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
