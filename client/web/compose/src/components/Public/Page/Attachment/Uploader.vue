<template>
  <vue-dropzone
    id="dropzone"
    ref="dropzone"
    :use-custom-slot="true"
    :include-styling="false"
    :options="dzOptions"
    @vdropzone-file-added="onFileAdded"
    @vdropzone-file-added-manually="onFileAdded"
    @vdropzone-success="onSuccess"
    @vdropzone-error="onError"
    @vdropzone-upload-progress="onUploadProgress"
  >
    <div class="drop-container w-100 h-100 position-relative bg-light rounded">
      <template v-if="active">
        <div
          class="bg-primary h-100 progress-bar position-absolute"
          :style="progresBarStyle"
        />

        <span class="d-flex align-items-center h-100 w-100 uploading justify-content-center position-relative py-2">
          {{ $t('label.uploading') }} {{ active.file.name }} ({{ size(active.file) }})
        </span>
      </template>
      <div
        v-else-if="processing"
        class="d-flex justify-content-center py-1"
      >
        <b-spinner
          variant="primary"
        />
      </div>
      <div
        v-else
        class="d-flex align-items-center h-100 w-100 p-2 droparea justify-content-center"
        :class="{ 'bg-danger': error }"
      >
        {{ error || label || $t('label.dropFiles') }}
      </div>
    </div>
  </vue-dropzone>
</template>
<script>
import numeral from 'numeral'
import vueDropzone from 'vue2-dropzone'
import 'vue2-dropzone/dist/vue2Dropzone.min.css'
import { mixins } from '@cortezaproject/corteza-vue'

export default {
  i18nOptions: {
    namespaces: 'general',
  },

  components: {
    vueDropzone,
  },

  mixins: [
    mixins.files,
  ],

  props: {
    endpoint: {
      type: String,
      required: true,
    },
    acceptedFiles: {
      type: Array,
      default: () => [],
    },
    maxFilesize: {
      type: Number,
      default: 100,
    },
    label: {
      type: String,
      default: null,
    },
    formData: {
      type: Object,
      required: false,
      default: () => ({}),
    },
    paramName: {
      type: String,
      default: 'upload',
    },
  },

  data () {
    return {
      processing: false,
      active: null,
      error: null,
    }
  },

  computed: {
    dropzone () {
      return (this.$refs.dropzone && this.$refs.dropzone.dropzone) ? this.$refs.dropzone.dropzone : false
    },

    dzOptions () {
      const vm = this

      return {
        paramName: this.paramName,
        maxFilesize: this.maxFilesize, // mb
        url: () => this.baseUrl + this.endpoint,
        thumbnailMethod: 'contain',
        thumbnailWidth: 320,
        thumbnailHeight: 180,
        maxFiles: 1000,
        withCredentials: true,
        autoProcessQueue: true,
        disablePreview: true,
        uploadMultiple: false,
        parallelUploads: 1,
        acceptedFiles: null,
        init: function () {
          this.on('sending', function (file, xhr, formData) {
            for (const k in vm.formData || {}) {
              formData.append(k, vm.formData[k])
            }
          })
        },
        headers: {
          // https://github.com/enyo/dropzone/issues/1154
          'Cache-Control': '',
          'X-Requested-With': '',
          Authorization: 'Bearer ' + this.$auth.accessToken,
        },
        addedfile (file) {},
      }
    },

    baseUrl () {
      return window.CortezaAPI + '/compose'
    },

    progresBarStyle () {
      return {
        width: this.active.progress + '%',
      }
    },
  },

  methods: {
    size (a) {
      return numeral(a.size).format('0b')
    },

    onSuccess (file, { response, error }) {
      if (error) {
        return this.onError(error, error.message)
      }

      this.processing = false
      this.active = null
      this.error = null
      this.$emit('uploaded', response, file)
    },

    onFileAdded (file) {
      this.error = null
      this.processing = true

      // Check if file type is allowed
      let types = this.acceptedFiles
      if (!types || !types.length) {
        types = ['*/*']
      }
      if (!this.validateFileType(file.name, types)) {
        this.$refs.dropzone.removeFile(file)
        this.onError(null, this.$t('label.fileTypeNotAllowed'))
      }
    },

    onError (e, message) {
      this.error = this.$t('label.uploadError', { message })
      this.processing = false
      this.active = null
    },

    onUploadProgress (file, progress, bytesSent) {
      this.active = { file, progress, bytesSent }
    },
  },
}
</script>

<style lang="scss" scoped>
.drop-container {
  &:hover {
    background-color: var(--extra-light) !important;
  }
}

.droparea {
  cursor: pointer;
}

.progress-bar {
  width: 0;
  opacity: 0.3;
}

.uploading {
  background-size: 100% 100%;
  background-position: right bottom;
  cursor: wait;
}
</style>
