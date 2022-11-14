<template>
  <b-row>
    <!-- Toolbox -->
    <b-col
      cols="12"
      lg="3"
      class="mb-3 mb-lg-0"
    >
      <editor-toolbox
        :template="template"
        :partials="partials"
      />
    </b-col>

    <!-- Content editor -->
    <b-col
      cols="12"
      lg="9"
    >
      <b-card
        class="shadow-sm"
        header-bg-variant="white"
        header-class="d-flex align-items-center"
        footer-bg-variant="white"
      >
        <component
          :is="editor"
          :template="template"
        />

        <template #header>
          <h3 class="m-0">
            {{ $t('title') }}
          </h3>
          <b-badge
            v-if="template.partial"
            data-test-id="badge-partial-template"
            variant="primary"
            class="ml-2"
          >
            {{ $t('partial') }}
          </b-badge>
        </template>

        <template #footer>
          <c-submit-button
            :disabled="!canCreate"
            class="float-right"
            :processing="processing"
            :success="success"
            @submit="$emit('submit', template)"
          />
        </template>
      </b-card>

      <!-- Preview configuration -->
      <b-card
        v-if="!template.partial"
        class="shadow-sm mt-3"
        header-bg-variant="white"
        footer-bg-variant="white"
      >
        <!-- Partial templates can't be previewed -->
        <ace-editor
          data-test-id="template-preview-output"
          :font-size="14"
          :show-print-margin="true"
          :show-gutter="true"
          :highlight-active-line="true"
          class="mt-1"
          width="100%"
          height="500px"
          mode="json"
          theme="chrome"
          name="preview-data"
          :on-change="(v) => previewData = v"
          :value="previewData"
          :editor-props="{
            $blockScrolling: false,
          }"
        />

        <template #header>
          <h3 class="m-0">
            {{ $t('preview.title') }}
          </h3>
        </template>

        <template #footer>
          <div
            class="float-right"
          >
            <b-btn
              v-if="canPreviewHTML"
              data-test-id="button-preview-html-template"
              variant="light"
              class="mr-2"
              @click="openPreview('html')"
            >
              {{ $t('preview.html') }}
            </b-btn>

            <b-btn
              v-if="canPreviewPDF"
              data-test-id="button-preview-pdf-template"
              variant="light"
              @click="openPreview('pdf')"
            >
              {{ $t('preview.pdf') }}
            </b-btn>
          </div>
        </template>
      </b-card>
    </b-col>
  </b-row>
</template>

<script>
import listHelpers from 'corteza-webapp-admin/src/mixins/listHelpers'
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'
import EditorToolbox from './EditorToolbox'
import EditorTextHtml from './EditorTextHtml'
import EditorTextPlain from './EditorTextPlain'
import EditorUnsupported from './EditorUnsupported'
import { Ace as AceEditor } from 'vue2-brace-editor'

import 'brace/mode/json'
import 'brace/theme/chrome'

export default {

  components: {
    CSubmitButton,
    AceEditor,
    EditorToolbox,
  },
  mixins: [
    listHelpers,
  ],

  i18nOptions: {
    namespaces: 'system.templates',
    keyPrefix: 'editor.content',
  },

  props: {
    template: {
      type: Object,
      required: true,
      default: () => ({}),
    },

    partials: {
      type: Array,
      required: false,
      default: () => [],
    },

    processing: {
      type: Boolean,
      value: false,
    },

    success: {
      type: Boolean,
      value: false,
    },

    canCreate: {
      type: Boolean,
      required: true,
    },
  },

  data () {
    return {
      // @todo i18n this?
      previewData: '{\n  "variables": {\n    "param1": "value1",\n    "param2": {\n      "nestedParam1": "value2"\n    }\n  },\n  "options": {\n    "documentSize": "A4",\n    "contentScale": "1",\n    "orientation": "portrait",\n    "margin": "0.3"\n  }\n}\n',

      previewBlob: '',

      availableDrivers: [],
    }
  },

  computed: {
    editor () {
      switch (this.template.type) {
        case 'text/html':
          return EditorTextHtml
        case 'text/plain':
          return EditorTextPlain
        default:
          return EditorUnsupported
      }
    },

    canPreviewHTML () {
      return this.availableDrivers.find(({ outputTypes }) => outputTypes.includes('text/html'))
    },

    canPreviewPDF () {
      return this.availableDrivers.find(({ outputTypes }) => outputTypes.includes('application/pdf'))
    },
  },

  async created () {
    this.availableDrivers = await this.$SystemAPI.templateRenderDrivers()
      .then(rsp => rsp.set)
      .catch(this.toastErrorHandler(this.$t('notification:template.driver.error')))
  },

  methods: {
    openPreview (ext) {
      this.incLoader()

      const cfg = {
        method: 'post',
        responseType: 'blob',
        url: this.$SystemAPI.templateRenderEndpoint({
          templateID: this.template.templateID,
          filename: 'preview',
          ext,
        }),
        data: JSON.parse(this.previewData),
      }

      this.$SystemAPI.api().request(cfg)
        .then(r => {
          this.previewBlob = window.URL.createObjectURL(r.data)
          window.open(this.previewBlob, '_newtab')
        })
        .catch(this.toastErrorHandler(this.$t('notification:template.preview.error')))
        .finally(() => {
          this.decLoader()
        })
    },
  },
}
</script>
