<template>
  <div>
    <b-button
      v-b-tooltip.noninteractive.hover="{ title: labels.tooltip, container: '#body' }"
      variant="light"
      :class="buttonClass"
      @click.prevent="openWebcamModal"
    >
      <slot />
    </b-button>

    <b-modal
      ref="webcamModal"
      :title="labels.modalTitle"
      size="lg"
      centered
      body-class="p-0"
      @show="initializeWebcam"
      @hidden="closeCamera"
    >
      <div
        v-if="showErrorMessage"
        class="p-3 text-danger"
      >
        {{ labels.cameraErrorMessage }}
      </div>

      <div
        v-else
        class="embed-responsive embed-responsive-4by3 d-flex justify-content-center align-items-center"
      >
        <b-spinner
          v-if="processingWebcam"
          variant="primary"
        />

        <video
          v-show="!processingWebcam && !hasCapturedImage"
          ref="video"
          autoplay
          playsinline
        />

        <img
          v-if="hasCapturedImage"
          :src="capturedImage"
          alt="Captured image"
          class="embed-responsive-item"
        >
      </div>

      <template #modal-footer>
        <div class="d-flex align-items-center gap-2">
          <b-button
            variant="light"
            @click="handleCloseClick"
          >
            {{ labels.cancelButtonLabel }}
          </b-button>

          <b-button
            :disabled="processingWebcam"
            variant="primary"
            @click="handleCaptureClick"
          >
            {{ hasCapturedImage ? labels.confirmButtonLabel : labels.captureButtonLabel }}
          </b-button>
        </div>
      </template>
    </b-modal>
  </div>
</template>

<script>
export default {
  name: 'CWebcamModal',

  props: {
    buttonClass: {
      type: String,
      default: 'd-flex align-items-center h-100',
    },
    labels: {
      type: Object,
      default: () => ({}),
    },
  },

  data () {
    return {
      video: null,
      stream: null,
      capturedImage: null,
      processingWebcam: true,
      showErrorMessage: false,
    }
  },

  computed: {
    hasCapturedImage () {
      return !!this.capturedImage
    },
  },

  methods: {
    openWebcamModal () {
      this.$refs.webcamModal.show()
    },

    async initializeWebcam () {
      await this.$nextTick()

      this.startWebcam()
    },

    startWebcam () {
      this.showErrorMessage = false

      // Get access to the camera
      navigator.mediaDevices.getUserMedia({
        video: {
          facingMode: 'user',
        },
      })
        .then(stream => {
          this.stream = stream
          this.$refs.video.srcObject = stream
        })
        .catch(err => {
          console.error('Error accessing the camera:', err)

          this.showErrorMessage = true
        }).finally(() => {
          this.processingWebcam = false
        })
    },

    capturePhoto () {
      const video = this.$refs.video
      const canvas = document.createElement('canvas')
      canvas.width = video.videoWidth
      canvas.height = video.videoHeight
      canvas.getContext('2d').drawImage(video, 0, 0, canvas.width, canvas.height)

      this.capturedImage = canvas.toDataURL('image/jpeg')
    },

    uploadCapturedImage () {
      if (!this.capturedImage) {
        return
      }

      fetch(this.capturedImage).then(res => res.blob()).then(blob => {
        const imageSuffix = new Date().toISOString().replace(/[:.]/g, '-')
        const file = new File([blob], `webcam-image-${imageSuffix}.jpg`, { type: 'image/jpeg' })

        this.$emit('upload', file)
      })

      this.$refs.webcamModal.hide()
    },

    stopWebcam () {
      if (!this.stream) {
        return
      }

      this.stream.getTracks().forEach(track => track.stop())
    },

    handleCaptureClick () {
      if (this.hasCapturedImage) {
        this.uploadCapturedImage()
      } else {
        this.capturePhoto()
      }
    },

    handleCloseClick () {
      if (this.hasCapturedImage) {
        this.discardCapturedImage()
      } else {
        this.$refs.webcamModal.hide()
      }
    },

    closeCamera () {
      this.stopWebcam()
      this.capturedImage = null
      this.processingWebcam = true
    },

    discardCapturedImage () {
      this.capturedImage = null
      this.startWebcam()
    },
  },
}
</script>
