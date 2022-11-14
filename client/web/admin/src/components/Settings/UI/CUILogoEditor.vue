<template>
  <b-row cols="12">
    <b-col
      cols="6"
    >
      <b-card
        class="shadow-sm"
      >
        <div
          class="d-flex justify-content-between"
        >
          <h4
            class="card-title"
          >
            {{ $t('mainLogo.title') }}
          </h4>
          <b-button
            v-if="uploadedFile('ui.main-logo')"
            variant="link"
            class="d-flex align-items-top text-dark p-1"
            @click="resetAttachment('ui.main-logo')"
          >
            <font-awesome-icon
              :icon="['far', 'trash-alt']"
            />
          </b-button>
        </div>

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
      cols="6"
    >
      <b-card
        class="shadow-sm"
      >
        <div
          class="d-flex justify-content-between"
        >
          <h4
            class="card-title"
          >
            {{ $t('iconLogo.title') }}
          </h4>
          <b-button
            v-if="uploadedFile('ui.icon-logo')"
            variant="link"
            class="d-flex align-items-top text-dark p-1"
            @click="resetAttachment('ui.icon-logo')"
          >
            <font-awesome-icon
              :icon="['far', 'trash-alt']"
            />
          </b-button>
        </div>

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
