<template>
  <b-card
    class="shadow-sm"
    data-test-id="card-user-info"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
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
          size="lg"
          size-confirm="lg"
          variant="danger"
          class="ml-2 h-100"
          @confirmed="$emit('resetAttachment', 'avatar')"
        >
          {{ $t('general:label.delete') }}
        </c-input-confirm>
      </div>

      <div class="form-row mt-3">
        <b-form-group
          :label="$t('initial.color')"
          class="col"
        >
          <b-form-input
            v-model="user.meta.avatarColor"
            type="color"
            data-test-id="input-handle"
          />
        </b-form-group>

        <b-form-group
          :label="$t('initial.backgroundColor')"
          class="col"
        >
          <b-form-input
            v-model="user.meta.avatarBgColor"
            type="color"
          />
        </b-form-group>
      </div>
    </b-form>

    <template #footer>
      <c-submit-button
        class="float-right"
        :processing="processing"
        :success="success"
        @submit="$emit('submit', user)"
      />
    </template>
  </b-card>
</template>

<script>
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'
import CUploaderWithPreview from 'corteza-webapp-admin/src/components/CUploaderWithPreview'

export default {
  name: 'CUserEditorAvatar',

  i18nOptions: {
    namespaces: 'system.users',
    keyPrefix: 'editor.avatar',
  },

  components: {
    CSubmitButton,
    CUploaderWithPreview,
  },

  props: {
    user: {
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
  },

  computed: {
    isKindAvatar () {
      return this.user.meta.avatarKind === 'avatar'
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
