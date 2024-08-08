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
          @click="newCodeSnippet()"
        >
          {{ $t('code-snippets.add') }}
        </b-button>
      </div>

      <b-table
        :items="providerItems"
        :fields="codeSnippetProviderFields"
        head-variant="light"
        show-empty
        hover
        class="mb-0"
        style="min-height: 50px;"
      >
        <template #cell(provider)="{ item }">
          {{ item.provider || item.tag }}
        </template>

        <template #empty>
          <p
            data-test-id="no-matches"
            class="text-center text-dark"
            style="margin-top: 1vh;"
          >
            {{ $t('code-snippets.empty') }}
          </p>
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
        id="modal-codeSnippet"
        v-model="modal.open"
        :title="modal.title"
        scrollable
        size="lg"
        title-class="text-capitalize"
        @ok="modal.updater(modal.data)"
      >
        <b-form-group
          :label="$t('code-snippets.form.provider.label')"
          label-class="text-primary"
        >
          <b-input-group>
            <b-form-input v-model="modal.data.name" />
          </b-input-group>
        </b-form-group>

        <div>
          <div class="mb-2">
            <h5>
              {{ $t('code-snippets.add') }}
            </h5>
            <span class="text-muted">
              {{ $t('code-snippets.form.value.description') }}
            </span>
          </div>

          <c-ace-editor
            v-model="modal.data.script"
            lang="javascript"
            height="500px"
            font-size="14px"
            show-line-numbers
            :border="false"
            :show-popout="false"
          />
        </div>

        <template #modal-footer="{ ok, cancel}">
          <c-input-confirm
            size="md"
            variant="danger"
            @confirmed="deleteCodeSnippet(modal.index)"
          >
            {{ $t('general:label.delete') }}
          </c-input-confirm>

          <b-button
            variant="light"
            class="ml-auto"
            @click="cancel()"
          >
            {{ $t('general:label.cancel') }}
          </b-button>

          <b-button
            variant="primary"
            @click="ok()"
          >
            {{ $t('general:label.saveAndClose') }}
          </b-button>
        </template>
      </b-modal>

      <template #footer>
        <c-button-submit
          :disabled="!canManage"
          :processing="codeSnippet.processing"
          :success="codeSnippet.success"
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
import { components } from '@cortezaproject/corteza-vue'
import { mapGetters } from 'vuex'
const { CAceEditor } = components

export default {
  name: 'CSystemCodeSnippetEditor',

  i18nOptions: {
    namespaces: 'system.code-snippets',
    keyPrefix: 'editor',
  },

  components: {
    CAceEditor,
  },

  mixins: [
    editorHelpers,
  ],

  data () {
    return {
      codeSnippets: [],
      modal: {
        open: false,
        editor: null,
        title: null,
        data: [],
        index: null,
      },

      codeSnippet: {
        processing: false,
        success: false,
      },
      originalCodeSnippets: [],
    }
  },

  computed: {
    ...mapGetters({
      canManage: 'rbac/can',
    }),

    codeSnippetProviderFields () {
      return [
        { key: 'provider', label: this.$t('code-snippets.table-headers.provider'), thStyle: { width: '200px' }, tdClass: 'text-capitalize' },
        { key: 'value', label: this.$t('code-snippets.table-headers.value'), tdClass: 'td-content-overflow' },
        { key: 'editor', label: '', thStyle: { width: '200px' }, tdClass: 'text-right' },
      ]
    },

    providerItems () {
      return this.codeSnippets.map((s, i) => ({
        provider: s.name,
        value: s.script,

        editor: {
          data: s,
          index: i,
          title: s.name,
          updater: (changed) => {
            this.codeSnippets[i] = changed
          },
        },
      }))
    },
  },

  created () {
    this.fetchSettings()

    this.originalCodeSnippets = [...this.codeSnippets]
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

    newCodeSnippet () {
      this.openEditor({
        title: this.$t('code-snippets.add'),
        data: {
          name: '',
          script: '<' + 'script> ' + '</' + 'script>',
        },
        updater: (changed) => {
          this.codeSnippets.push(changed)
        },
      })
    },
    fetchSettings () {
      this.incLoader()
      this.$Settings.fetch()
      return this.$SystemAPI.settingsList({ prefix: 'code-snippets' })
        .then(settings => {
          if (settings && settings[0]) {
            this.codeSnippets = settings[0].value
          } else {
            this.codeSnippets = []
          }
        })
        .catch(this.toastErrorHandler(this.$t('notification:settings.codeSnippet.fetch.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    settingsUpdate (action) {
      this.codeSnippet.processing = true
      this.$SystemAPI.settingsUpdate({ values: [{ name: 'code-snippets', value: this.codeSnippets }] })
        .then(() => {
          this.animateSuccess('codeSnippet')
          if (action === 'delete') {
            this.toastSuccess(this.$t('notification:settings.codeSnippet.delete.success'))
          } else {
            this.toastSuccess(this.$t('notification:settings.codeSnippet.update.success'))
          }
        })
        .catch(this.toastErrorHandler(this.$t('notification:settings.codeSnippet.update.error')))
        .finally(() => {
          this.codeSnippet.processing = false
        })
    },
    onSubmit () {
      this.settingsUpdate('update')
    },

    deleteCodeSnippet (i) {
      this.codeSnippets.splice(i, 1)
      this.settingsUpdate('delete')
      this.$bvModal.hide('modal-codeSnippet')
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
