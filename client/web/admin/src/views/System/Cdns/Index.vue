<template>
  <b-container
    class="pt-2 pb-3"
  >
    <c-content-header
      :title="$t('title')"
    />
    <b-card
      body-class="p-0"
      header-class="border-bottom"
      footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
      class="shadow-sm"
    >
      <div class="align-items-center gap-1 p-3">
        <b-button
          variant="primary"
          size="lg"
          @click="newCDN()"
        >
          {{ $t('cdns-provider.add') }}
        </b-button>
      </div>

      <b-table
        :items="providerItems"
        :fields="cdnProviderFields"
        head-variant="light"
        hover
        class="mb-0"
      >
        <template #cell(provider)="{ item }">
          {{ item.provider || item.tag }}
        </template>

        <template #cell(editor)="{ item }">
          <c-input-confirm
            v-if="item.delete"
            :icon="item.deleted ? ['fas', 'trash-restore'] : undefined"
            @confirmed="item.delete()"
          />

          <b-button
            variant="link"
            @click="openEditor(item.editor)"
          >
            <font-awesome-icon
              :icon="['fas', 'wrench']"
            />
          </b-button>
        </template>
      </b-table>

      <b-modal
        id="modal-cdn"
        v-model="modal.open"
        :title="modal.title"
        scrollable
        size="lg"
        title-class="text-capitalize"
        @ok="modal.updater(modal.data)"
      >
        <component
          :is="modal.component"
          v-model="modal.data"
        />
        <template #modal-footer="{ ok, cancel}">
          <b-button
            size="sm"
            variant="danger"
            @click="deleteCDN(modal.index)"
          >
            Delete
          </b-button>
          <b-button
            size="sm"
            variant="secondary"
            @click="cancel()"
          >
            Cancel
          </b-button>
          <b-button
            size="sm"
            variant="primary"
            @click="ok()"
          >
            OK
          </b-button>
        </template>
      </b-modal>

      <template #footer>
        <c-button-submit
          :disabled="!canManage"
          :processing="processing"
          :success="success"
          :text="$t('admin:general.label.submit')"
          class="ml-auto"
          @submit="onSubmit()"
        />
      </template>
    </b-card>
  </b-container>
</template>

<script>
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CdnProviderEditor from 'corteza-webapp-admin/src/components/Settings/System/CDNs/CSystemCdnProviderEditor'
import { mapGetters } from 'vuex'

export default {
  name: 'CSystemCdnEditor',

  i18nOptions: {
    namespaces: 'system.cdns',
    keyPrefix: 'editor',
  },

  components: {
    CdnProviderEditor,
  },

  mixins: [
    editorHelpers,
  ],

  data () {
    return {
      cdns: [],
      modal: {
        open: false,
        editor: null,
        title: null,
        data: [],
        index: null,
      },

      cdn: {
        processing: false,
        success: false,
      },
      originalCdns: [],
    }
  },

  computed: {
    ...mapGetters({
      canManage: 'rbac/can',
    }),

    cdnProviderFields () {
      return [
        { key: 'provider', label: this.$t('cdns-provider.table-headers.provider'), thStyle: { width: '200px' }, tdClass: 'text-capitalize' },
        { key: 'value', label: this.$t('cdns-provider.table-headers.value'), tdClass: 'td-content-overflow' },
        { key: 'editor', label: '', thStyle: { width: '200px' }, tdClass: 'text-right' },
      ]
    },

    providerItems () {
      return this.cdns.map((s, i) => ({
        provider: s.name,
        value: s.cdnScript,

        editor: {
          component: 'cdn-provider-editor',
          data: s,
          index: i,
          title: s.name,
          updater: (changed) => {
            this.cdns[i] = changed
          },
        },
      }))
    },
  },

  created () {
    this.fetchSettings()

    this.originalCdns = [...this.cdns]
  },
  methods: {
    openEditor ({ component, title, data, updater }) {
      this.modal.open = true
      this.modal.component = component
      this.modal.title = title
      this.modal.updater = updater

      // deref
      this.modal.data = data
    },

    newCDN () {
      this.openEditor({
        component: 'cdn-provider-editor',
        title: this.$t('cdns-provider.add'),
        data: {},
        updater: (changed) => {
          this.cdns.push(changed)
        },
      })
    },
    fetchSettings () {
      this.incLoader()
      this.$Settings.fetch()
      return this.$SystemAPI.settingsList({ prefix: 'cdns' })
        .then(settings => {
          if (settings[0]) {
            this.cdns = settings[0].value
          } else {
            this.cdns = []
          }
        })
        .catch(this.toastErrorHandler(this.$t('notification:settings.cdn.fetch.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    settingsUpdate (action) {
      this.cdn.processing = true
      this.$SystemAPI.settingsUpdate({ values: [{ name: 'cdns', value: this.cdns }] })
        .then(() => {
          if (action === 'delete') {
            this.animateSuccess('cdn')
            this.toastSuccess(this.$t('notification:settings.cdn.delete.success'))
          } else {
            this.animateSuccess('cdn')
            this.toastSuccess(this.$t('notification:settings.cdn.update.success'))
          }
        })
        .catch(this.toastErrorHandler(this.$t('notification:settings.cdn.update.error')))
        .finally(() => {
          this.cdn.processing = false
        })
    },
    onSubmit () {
      this.settingsUpdate('update')
    },

    deleteCDN (i) {
      this.cdns.splice(i, 1)
      this.settingsUpdate('delete')
      this.$bvModal.hide('modal-cdn')
    },
  },
}

</script>

<style>
  .td-content-overflow {
    max-width: 100px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
</style>
