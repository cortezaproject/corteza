<template>
  <b-form-group
    :label-cols-md="horizontal && '5'"
    :label-cols-xl="horizontal && '4'"
    :content-cols-md="horizontal && '7'"
    :content-cols-xl="horizontal && '8'"
    :class="formGroupStyleClasses"
  >
    <template
      #label
    >
      <div
        v-if="!valueOnly"
        class="d-flex align-items-center text-primary p-0"
      >
        <span
          :title="label"
          class="d-inline-block mw-100"
        >
          {{ label }}
        </span>

        <c-hint :tooltip="hint" />

        <slot name="tools" />
      </div>
      <div
        class="small text-muted"
        :class="{ 'mb-1': description }"
      >
        {{ description }}
      </div>
    </template>

    <div class="d-flex gap-1">
      <uploader
        ref="uploader"
        :endpoint="endpoint"
        :accepted-files="mimetypes"
        :max-filesize="maxSize"
        :form-data="uploaderFormData"
        class="flex-grow-1"
        @uploaded="appendAttachment"
      />

      <c-webcam
        v-if="field.options.enableWebcam"
        :labels="{
          tooltip: $t('webcam.tooltip'),
          modalTitle: $t('webcam.title'),
          cancelButtonLabel: $t('webcam.buttons.cancel'),
          confirmButtonLabel: $t('webcam.buttons.confirm'),
          captureButtonLabel: $t('webcam.buttons.capture'),
          cameraErrorMessage: $t('webcam.errors.camera')
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
      v-if="set.length > 0"
      kind="record"
      :set.sync="set"
      :namespace="namespace"
      :enable-order="field.isMulti"
      enable-delete
      mode="list"
      class="mt-3"
    />

    <errors :errors="errors" />
  </b-form-group>
</template>
<script>
import base from './base'
import Uploader from 'corteza-webapp-compose/src/components/Public/Page/Attachment/Uploader'
import ListLoader from 'corteza-webapp-compose/src/components/Public/Page/Attachment/ListLoader'
import { NoID } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'general',
  },

  components: {
    Uploader,
    ListLoader,
  },

  extends: base,

  computed: {
    endpoint () {
      const { moduleID, recordID } = this.record
      const { namespaceID } = this.namespace

      return this.$ComposeAPI.recordUploadEndpoint({
        namespaceID,
        moduleID,
        recordID,
        fieldName: this.field.name,
      })
    },

    uploaderFormData () {
      const fd = {
        fieldName: this.field.name,
      }

      if (this.record && this.record.recordID !== NoID) {
        fd.recordID = this.record.recordID
      }

      return fd
    },

    mimetypes () {
      const a = (this.field.options.mimetypes || '').trim()
      if (!a) {
        return this.$s('compose.Record.Attachments.Mimetypes', ['*/*'])
      }

      return a.split(',').map(p => p.trim())
    },

    maxSize () {
      return this.field.options.maxSize || this.$s('compose.Record.Attachments.MaxSize', 100)
    },

    set: {
      get () {
        return this.field.isMulti ? this.value : [this.value].filter(id => !!id)
      },

      set (v) {
        if (this.field.isMulti) {
          this.value = v
        } else {
          this.value = (Array.isArray(v) && v.length > 0) ? v[0] : undefined
        }
      },
    },
  },

  methods: {
    appendAttachment ({ attachmentID } = {}) {
      if (this.field.isMulti) {
        this.value.push(attachmentID)
      } else {
        this.value = attachmentID
      }
    },

    uploadWebcamImage (file) {
      const uploader = this.$refs.uploader
      uploader.$refs.dropzone.addFile(file)
    },
  },
}
</script>
