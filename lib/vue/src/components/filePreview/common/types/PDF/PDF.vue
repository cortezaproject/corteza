<template>
  <div
    :style="previewStyle"
    :class="[...previewClass, 'pdf-preview', inline ? 'inline' : '', $listeners.click ? 'clickable' : '']"
    @click="onPreviewClick"
  >
    <!-- Container for pdf's pages -->
    <div
      v-show="show"
      ref="pages"
      class="pages"
    />

    <div
      v-if="loadError"
      class="doc-msg doc-err"
    >
      <div>
        <h4 class="err-message">
          {{ loadError.message }}
        </h4>
        <p v-if="labels.clickToRetry">
          {{ labels.clickToRetry }}
        </p>
      </div>
    </div>
    <div
      v-else-if="!show && labels.loading"
      class="doc-msg doc-err"
    >
      <p>{{ labels.loading }}</p>
    </div>
    <div
      v-else-if="!pageCount && labels.noPages"
      class="doc-msg doc-err"
    >
      <p>{{ labels.noPages }}</p>
    </div>
    <template v-else>
      <div v-if="!inline && labels.downloadForAll && show && hasMore">
        <p>{{ labels.downloadForAll }}</p>
      </div>
      <div v-else-if="inline && labels.firstPagePreview && show">
        <p>{{ labels.firstPagePreview }}</p>
      </div>
    </template>
  </div>
</template>

<script lang="js">
import pdfjs from 'pdfjs-dist'
import pdfjsWorker from 'pdfjs-dist/build/pdf.worker.entry'
import base from '../base.vue'
import { makePlaceholder, makeFailedPage, Page, Document } from './helpers'

pdfjs.GlobalWorkerOptions.workerSrc = pdfjsWorker

function sleep (t) {
  return new Promise(resolve => setTimeout(resolve, t))
}

export default {
  extends: base,

  props: {
    retryBackoff: {
      type: Number,
      required: false,
      default: 300,
    },
    maxRetries: {
      type: Number,
      required: false,
      default: 10,
    },

    maxPages: {
      required: false,
      type: Number,
      default: 25,
    },

    initialScale: {
      required: false,
      type: Number,
      default: 1,
    },
  },

  data () {
    return {
      document: null,
      pages: [],
      show: false,
      loadError: undefined,
    }
  },

  computed: {
    /**
     * Helper to provide the number of pages for the given PDF
     * @returns {Number}
     */
    pageCount () {
      if (!this.document || !this.document.pdf) {
        return 0
      }
      return this.document.pdf.numPages
    },

    /**
     * Helper to determine if the given PDF has more pages then we are willing to show
     * @returns {Boolean}
     */
    hasMore () {
      return this.maxPages < this.pageCount
    },
  },

  created () {
    if (!this.src) {
      this.stdErr(new Error('src.missing'))
      return
    }

    this.$nextTick(() => this.init())
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    /**
     * Helper to handle on preview click. It either requests a retry or
     * emits an open event
     */
    onPreviewClick () {
      if (this.loadError) {
        this.init()
      } else {
        this.$emit('openPreview', { document: this.document })
      }
    },

    /**
     * Initializes the state, loads the document & renderes it's pages
     */
    async init () {
      this.document = null
      this.pages = []
      this.show = false
      this.loadError = undefined

      return this.loadDocument(this.src)
        .then(this.renderDocument)
        .catch(this.stdErr)
    },

    /**
     * Helper method to load the given document. Needed for test stubbing
     * @param {String} src Document's src
     * @returns {Promise<PDFDocumentProxy>}
     */
    async pdfjsLoad (src) {
      return pdfjs.getDocument(src).promise
    },

    /**
     * Loads the given PDF. It can load it from API or re-use existing document
     * @param {Document|String} src PDF's source or a Document object
     * @returns {Document}
     */
    async loadDocument (src) {
      if (src instanceof Document) {
        this.document = new Document({ ...src, scale: this.initialScale })
      } else if (typeof src === 'string') {
        let retries = 0
        let err
        const pdfl = async () => {
          return sleep(retries * this.retryBackoff)
            .then(() => this.pdfjsLoad(src))
            .then(pdf => {
              this.document = new Document({ pdf, src, scale: this.initialScale })
            })
        }

        // Retry with backoff it it fails to load
        while (!this.document && retries < this.maxRetries) {
          await pdfl().catch(e => {
            retries++
            err = e
          })
        }

        if (!this.document) {
          throw err
        }
      } else {
        throw new Error('src.notValid')
      }
      return this.document
    },

    /**
     * Renders the given PDF
     * @param {Document} doc The Document to render
     * @return {Promise<undefined>}
     */
    async renderDocument (doc) {
      const rf = this.$refs.pages

      const pgCount = Math.min(this.pageCount, this.maxPages)
      this.pages = [...new Array(pgCount)].map((_, i) => new Page({ index: i }))

      if (pgCount <= 0) {
        this.show = true
        return
      }

      // Loadup pages
      for (let i = 0; i < pgCount; i++) {
        const node = makePlaceholder(this.labels)
        rf.appendChild(node)
        this.pages.splice(i, 1, new Page({ ...this.pages[i], node, loading: true }))

        this.renderPage(this.pages[i], doc, rf)
          .then(page => {
            this.pages.splice(page.index, 1, page)

            if (page.index === 0) {
              this.show = true
            }
          })
          .catch(this.stdErr)
      }
    },

    /**
     * Renders the given page
     * @param {Page} page The page in question
     * @param {Document} doc Page source
     * @param {Node} rf PDF container
     * @returns {Promise<undefined>}
     */
    async renderPage (page, doc, rf) {
      // pdfjs starts with 1!
      return doc.pdf.getPage(page.index + 1).then(p => {
        const np = new Page(page)
        np.loading = false
        np.loaded = true
        np.page = p

        // Render page
        const canvas = document.createElement('canvas')
        const scale = doc.scale
        const viewport = np.page.getViewport(scale)
        const canvasContext = canvas.getContext('2d')
        const renderContext = { canvasContext, viewport }
        canvas.height = viewport.height
        canvas.width = viewport.width

        return np.page.render(renderContext).then(() => {
          np.node = canvas
          np.rendered = true
          if (this.inline) {
            canvas.classList.add('inline')
          }
          return np
        })

          .catch(() => {
            const node = makeFailedPage(this.labels)
            np.node = node
            np.rendered = false
            np.failed = true
            return np
          })

          .then(np => {
            rf.replaceChild(np.node, page.node)
            return np
          })
      })
    },

    /**
     * Standard error handler
     * @param {Error} err The error
     */
    stdErr (err) {
      this.loadError = err
      this.$emit('error', err)
    },

    setDefaultValues () {
      this.document = null
      this.pages = []
      this.show = false
      this.loadError = undefined
    },
  },
}
</script>

<style lang="scss" scoped>
$white: #FFFFFF !default;
$danger: #E54122 !default;

.doc-msg {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  width: 100%;
  background-color: $white;
}
.doc-err {
  cursor: pointer;

  .err-message {
    color: $danger;
  }
}

</style>

<style lang="scss">
// Style not scoped, since pages are manually rendered

.pdf-preview {
  text-align: center;
  &:not(.inline) {
    padding-top: 20px;
    padding-bottom: 20px;
    canvas {
      box-shadow: 0 0 3px #1E1E1E41;
    }
  }
  &.inline {
    height: 200px;
    overflow-y: scroll;
    display: inline-block;
    cursor: zoom-in;
    width: 100%;
    max-width: 500px;
  }

  canvas {
    margin-bottom: 10px;
    &.inline {
      width: 100%;
    }

    &:not(.inline) {
      margin: 0 auto 10px auto;
      display: block;
    }

    &:last-of-type {
      margin-bottom: unset;
    }
  }

  .loader {
    margin-bottom: 10px;
  }
}

</style>
