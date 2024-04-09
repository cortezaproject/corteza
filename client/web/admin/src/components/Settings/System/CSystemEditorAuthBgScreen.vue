<template>
  <b-card
    header-class="border-bottom"
    footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <template #header>
      <h4 class="m-0">
        {{ $t("title") }}
      </h4>
    </template>

    <b-row>
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
          <c-ace-editor
            v-model="settings['auth.ui.styles']"
            data-test-id="auth-bg-image-styling-editor"
            name="editor/css"
            lang="css"
            height="300px"
            font-size="14px"
            show-line-numbers
            class="flex-fill w-100"
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
import { components } from '@cortezaproject/corteza-vue'
import CUploaderWithPreview from 'corteza-webapp-admin/src/components/CUploaderWithPreview'

const { CAceEditor } = components

export default {
  name: 'CSystemEditorAuthBgImage',

  i18nOptions: {
    namespaces: 'system.settings',
    keyPrefix: 'editor.bgScreen',
  },

  components: {
    CUploaderWithPreview,
    CAceEditor,
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
