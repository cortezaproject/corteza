<template>
  <div class="d-flex">
    <b-btn
      variant="light"
      size="lg"
      class="flex-fill"
      @click="showModal=true"
    >
      {{ $t('label.import') }}
    </b-btn>

    <b-modal
      v-model="showModal"
      size="lg"
      :title="$t('label.import')"
    >
      <b-input-group>
        <!-- To handle file upload -->
        <template v-if="!importObj">
          <b-form-file
            :placeholder="$t('label.importPlaceholder')"
            :browse-text="$t('label.browse')"
            @change="loadFile"
          />

          <h6
            v-if="processing"
            class="my-auto ml-3 "
          >
            {{ $t('label.processing') }}
          </h6>
        </template>

        <!-- To confirm selection & import -->
        <template v-else>
          <b-container class="p-0">
            <b-row
              no-gutters
              class="mb-3"
            >
              <b-button
                variant="light"
                @click="selectAll(true)"
              >
                {{ $t('field.selectAll') }}
              </b-button>
              <b-button
                class="ml-2"
                variant="light"
                @click="selectAll(false)"
              >
                {{ $t('field.unselectAll') }}
              </b-button>
            </b-row>
            <b-row no-gutters>
              <b-col
                v-for="(o, index) in importObj.list"
                :key="index"
                cols="12"
                md="6"
                lg="4"
              >
                <b-form-checkbox v-model="o.import">
                  {{ o.name || o.title }}
                </b-form-checkbox>
              </b-col>
            </b-row>
          </b-container>
        </template>
      </b-input-group>

      <div slot="modal-footer">
        <b-button
          :disabled="!importObj || !importObj.list.filter(i => i.import).length > 0"
          variant="primary"
          @click="jsonImport(importObj)"
        >
          {{ $t('label.import') }}
        </b-button>
      </div>
    </b-modal>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { compose } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'general',
  },

  props: {
    namespace: {
      type: Object,
      required: true,
    },
    type: {
      type: String,
      required: true,
    },
  },

  data () {
    return {
      showModal: false,
      importObj: null,
      processing: false,
      classes: {
        module: compose.Module,
        chart: compose.Chart,
      },
    }
  },

  computed: {
    ...mapGetters({
      modules: 'module/set',
    }),
  },

  methods: {
    async jsonImport ({ list, type }) {
      this.processing = true
      const { namespaceID } = this.namespace
      const ItemClass = this.classes[type]
      try {
        for (let item of list.filter(i => i.import)) {
          if (this.importObj) {
            item = new ItemClass(item).import(this.getModuleID)
            item.namespaceID = namespaceID
            await this.$store.dispatch(`${this.type}/create`, item)
          } else {
            break
          }
        }
        this.$emit('importSuccessful')
      } catch (e) {
        this.toastErrorHandler(this.$t('notification:general.import.failed'))(e)
      }
      this.cancelImport()
    },

    getModuleID (moduleName) {
      // find module, replaceID
      const matchedModules = this.modules.filter(m => m.name === moduleName)
      if (matchedModules.length > 0) {
        return matchedModules[0].moduleID
      } else {
        return null
      }
    },

    selectAll (selectAll) {
      this.importObj.list = this.importObj.list.map(i => {
        i.import = selectAll && true
        return i
      })
    },

    cancelImport () {
      this.importObj = null
      this.processing = false
      this.showModal = false
    },

    loadFile (e = {}) {
      const { files = [] } = (e.type === 'drop' ? e.dataTransfer : e.target) || {}

      if (files[0]) {
        this.processing = true
        const reader = new FileReader()
        reader.readAsText(files[0])
        reader.onload = (evt) => {
          try {
            this.importObj = JSON.parse(evt.target.result)
            if (!this.importObj.list) {
              throw new Error(this.$t('notification:general.import.readingError'))
            } else {
              this.importObj.list = this.importObj.list.map(i => {
                return { import: true, ...i }
              })
            }
          } catch (err) {
            this.toastErrorHandler(this.$t('notification:general.import.failed'))(err)
            this.importObj = null
          } finally {
            this.processing = false
          }
        }
        reader.onerror = (evt) => {
          this.toastErrorHandler(this.$t('notification:general.import.readingError'))
          this.processing = false
        }
      }
    },
  },
}
</script>
