<template>
  <b-row>
    <!-- Toolbox -->
    <b-col
      cols="12"
      lg="3"
      class="mb-3"
    >
      <editor-toolbox
        :template="template"
        :partials="partials"
      />
    </b-col>

    <!-- Preview configuration -->
    <b-col
      cols="12"
      lg="9"
      class="mb-3"
    >
      <b-card
        v-if="!template.partial && !fresh"
        body-class="p-0"
        header-class="border-bottom"
        footer-class="border-top d-flex justify-content-end flex-wrap flex-fill-child gap-1"
      >
        <!-- Partial templates can't be previewed -->
        <c-ace-editor
          v-model="previewData"
          data-test-id="template-preview-output"
          name="preview-data"
          lang="json"
          min-height="300px"
          show-line-numbers
          highlight-active-line
          resizable
          :border="false"
        />

        <template #header>
          <h4 class="m-0">
            {{ $t('preview.title') }}
          </h4>
        </template>

        <template #footer>
          <b-btn
            v-if="canPreviewHTML"
            data-test-id="button-preview-html-template"
            variant="light"
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
        </template>
      </b-card>
    </b-col>

    <!-- Content editor -->
    <b-col cols="12">
      <b-card
        body-class="p-0"
        header-class="d-flex align-items-center border-bottom"
        footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
        class="shadow-sm"
      >
        <b-row
          v-if="!template.partial"
          class="p-3 border-bottom no-gutters"
        >
          <b-col
            cols="12"
            lg="6"
            class="pr-lg-2"
          >
            <b-form-group
              :label="$t('headerTemplate')"
              :description="$t('headerTemplateDescription')"
              label-class="text-primary"
              class="mb-0"
            >
              <b-form-select
                v-model="headerTemplateID"
                data-test-id="select-header-template"
                :options="partialOptions"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            lg="6"
            class="pl-lg-2"
          >
            <b-form-group
              :label="$t('footerTemplate')"
              :description="$t('footerTemplateDescription')"
              label-class="text-primary"
              class="mb-0"
            >
              <b-form-select
                v-model="footerTemplateID"
                data-test-id="select-footer-template"
                :options="partialOptions"
              />
            </b-form-group>
          </b-col>
        </b-row>

        <component
          :is="editor"
          :template="template"
        />

        <template #header>
          <h4 class="m-0">
            {{ $t('title') }}
          </h4>
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
          <c-button-submit
            :disabled="!canCreate"
            :processing="processing"
            :success="success"
            :text="$t('admin:general.label.submit')"
            class="ml-auto"
            @submit="$emit('submit', template)"
          />
        </template>
      </b-card>
    </b-col>
  </b-row>
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'
import listHelpers from 'corteza-webapp-admin/src/mixins/listHelpers'
import EditorToolbox from './EditorToolbox'
import EditorTextHtml from './EditorTextHtml'
import EditorTextPlain from './EditorTextPlain'
import EditorUnsupported from './EditorUnsupported'
import { components } from '@cortezaproject/corteza-vue'

const { CAceEditor } = components

export default {

  components: {
    CAceEditor,
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
    fresh () {
      return !this.template.templateID || this.template.templateID === NoID
    },

    partialOptions () {
      const options = this.partials
        .filter(({ templateID }) => templateID !== this.template.templateID)
        .map(({ templateID, meta = {}, handle: h }) => ({ value: templateID, text: meta.short || h || templateID }))

      return [{ value: NoID, text: this.$t('general:label.none') }, ...options]
    },

    headerTemplateID: {
      get () {
        return this.template.meta.headerTemplateID || NoID
      },
      set (headerTemplateID) {
        this.$set(this.template.meta, 'headerTemplateID', headerTemplateID)
      },
    },

    footerTemplateID: {
      get () {
        return this.template.meta.footerTemplateID || NoID
      },
      set (footerTemplateID) {
        this.$set(this.template.meta, 'footerTemplateID', footerTemplateID)
      },
    },

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
