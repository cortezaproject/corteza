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
      <c-uploader
        ref="uploader"
        :endpoint="fileUploadEndpoint"
        :accepted-files="mimetypes"
        :max-filesize="maxSize"
        :form-data="uploaderFormData"
        :labels="{
          uploading: $t('general:label.uploading'),
          placeholder: $t('general:label.dropFiles'),
          fileTypeNotAllowed: $t('general:label.fileTypeNotAllowed'),
        }"
        class="flex-grow-1"
        @upload="appendAttachment"
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
import { components } from '@cortezaproject/corteza-vue'
import ListLoader from 'corteza-webapp-compose/src/components/Public/Page/Attachment/ListLoader'
import { NoID } from '@cortezaproject/corteza-js'
const { CUploader } = components

export default {
  i18nOptions: {
    namespaces: 'general',
  },

  components: {
    CUploader,
    ListLoader,
  },

  extends: base,

  computed: {
    fileUploadEndpoint () {
      const { moduleID, recordID } = this.record
      const { namespaceID } = this.namespace

      return this.$ComposeAPI.baseURL + this.$ComposeAPI.recordUploadEndpoint({
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
    appendAttachment (response, _file) {
      if (response && response.duplicate) {
        this.handleDuplicate(response)
        return
      }

      const { attachmentID } = response || {}
      if (this.field.isMulti) {
        this.value = [attachmentID, ...this.value]
      } else {
        this.value = attachmentID
      }
    },

    handleDuplicate (response) {
      const {
        existingRecordID,
        conflictAction,
        recordPageID,
        namespaceSlug,
        namespaceID,
      } = response

      // getByUrlPart in the namespace store accepts either slug or namespaceID,
      // so fall back to namespaceID when the namespace has no slug set.
      const nsUrlPart = namespaceSlug || namespaceID

      // canNavigate requires a resolved record and page. Namespace is always
      // available (we always have at least namespaceID from the response).
      const canNavigate = existingRecordID && existingRecordID !== '0' &&
        recordPageID && recordPageID !== '0'

      const routeParams = canNavigate
        ? {
            name: 'page.record',
            params: {
              recordID: existingRecordID,
              pageID: recordPageID,
              slug: nsUrlPart,
            },
          }
        : null

      const showModal = () => {
        this.$root.$emit('show-record-modal', {
          recordID: existingRecordID,
          recordPageID,
          edit: false,
        })
      }

      const action = canNavigate ? conflictAction : 'alert'

      switch (action) {
        case 'modal':
          showModal()
          break

        case 'newTab': {
          const resolved = this.$router.resolve(routeParams)
          window.open(resolved.href, '_blank')
          break
        }

        case 'sameTab':
          this.$router.push(routeParams)
          break

        case 'alert':
        default: {
          const h = this.$createElement
          const vNodeMsg = h('div', [
            h('span', this.$t('notification:field.duplicateAlert', { recordID: existingRecordID })),
            existingRecordID && existingRecordID !== '0' && recordPageID && recordPageID !== '0'
              ? h(
                'b-button',
                {
                  props: { variant: 'link', size: 'sm' },
                  on: {
                    click: () => {
                      this.$root.$emit('bv::hide::toast', 'duplicate-file-alert')
                      showModal()
                    },
                  },
                  class: 'p-0 ml-2',
                },
                this.$t('notification:field.openRecord'),
              )
              : null,
          ])
          this.$bvToast.toast([vNodeMsg], {
            id: 'duplicate-file-alert',
            title: this.$t('notification:general.warning'),
            variant: 'warning',
            solid: true,
            toaster: 'b-toaster-top-right',
          })
          break
        }
      }
    },

    uploadWebcamImage (file) {
      const uploader = this.$refs.uploader
      uploader.$refs.dropzone.addFile(file)
    },
  },
}
</script>
