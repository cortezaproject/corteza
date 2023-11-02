<template>
  <b-form
    @submit.prevent="$emit('upload')"
  >
    <div
      v-if="value"
      class="preview d-flex w-100 mb-2"
    >
      <b-img
        :src="value"
        class="w-100 h-auto"
      />
    </div>

    <c-uploader
      v-if="!disabled"
      :endpoint="endpoint"
      :labels="labels"
      :accepted-files="['image/*']"
      @upload="$emit('upload', $event)"
    />
  </b-form>
</template>

<script>
import CUploader from 'corteza-webapp-admin/src/components/CUploader'

export default {
  name: 'CUploaderWithPreview',

  components: {
    CUploader,
  },

  props: {
    value: {
      type: String,
      default: () => undefined,
    },

    disabled: {
      type: Boolean,
      default: () => false,
    },

    labels: {
      type: Object,
      default: () => ({}),
    },

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
  },

  data () {
    return {
      panels: [],
    }
  },
}
</script>
<style lang="scss">
.preview {
  height: 300px;
  background: transparent;
}
</style>
