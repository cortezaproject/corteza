<template>
  <b-row>
    <b-col
      cols="12"
      lg="6"
    >
      <b-card
        header-class="d-flex align-items-center justify-content-between pr-2"
        class="shadow-sm mb-3"
      >
        <template #header>
          <h3 class="mb-0">
            {{ $t('mainLogo.title') }}
          </h3>

          <c-input-confirm
            v-if="uploadedFile('ui.main-logo')"
            show-icon
            @confirmed="resetAttachment('ui.main-logo')"
          />
        </template>

        <c-uploader-with-preview
          :value="uploadedFile('ui.main-logo')"
          :endpoint="'/settings/ui.main-logo'"
          :disabled="!canManage"
          :labels="$t('mainLogo.uploader', { returnObjects: true })"
          @upload="onUpload($event)"
        />
      </b-card>
    </b-col>

    <b-col
      cols="12"
      lg="6"
    >
      <b-card
        header-class="d-flex align-items-center justify-content-between"
        class="shadow-sm mb-3"
      >
        <template #header>
          <h3 class="mb-0">
            {{ $t('iconLogo.title') }}
          </h3>

          <c-input-confirm
            v-if="uploadedFile('ui.icon-logo')"
            show-icon
            @confirmed="resetAttachment('ui.icon-logo')"
          />
        </template>

        <c-uploader-with-preview
          :value="uploadedFile('ui.icon-logo')"
          :endpoint="'/settings/ui.icon-logo'"
          :disabled="!canManage"
          :labels="$t('iconLogo.uploader', { returnObjects: true })"
          @upload="onUpload($event)"
          @clear="resetAttachment('ui.icon-logo')"
        />
      </b-card>
    </b-col>
  </b-row>
</template>

<script>
import CUploaderWithPreview from 'corteza-webapp-admin/src/components//CUploaderWithPreview'

export default {
  name: 'CUILogoEditor',

  components: {
    CUploaderWithPreview,
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
  },

  methods: {
    onUpload ({ name, value }) {
      this.$set(this.settings, name, value)
    },

    resetAttachment (name) {
      this.$SystemAPI.settingsUpdate({ values: [{ name, value: undefined }], upload: {} })
        .then(() => {
          this.$set(this.settings, name, undefined)
        })
    },

    uploadedFile (name) {
      const localAttachment = /^attachment:(\d+)/

      switch (true) {
        case this.settings[name] && localAttachment.test(this.settings[name]):
          const [, attachmentID] = localAttachment.exec(this.settings[name])

          return this.$SystemAPI.baseURL +
            this.$SystemAPI.attachmentOriginalEndpoint({
              attachmentID,
              kind: 'settings',
              name,
            })
      }

      return undefined
    },
  },
}
</script>
