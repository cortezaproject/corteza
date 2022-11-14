<template>
  <b-container
    v-if="template"
    class="py-3"
  >
    <c-content-header
      :title="title"
    >
      <span
        class="text-nowrap"
      >
        <b-button
          v-if="templateID && canCreate"
          data-test-id="button-new-template"
          variant="primary"
          class="mr-2"
          :to="{ name: 'system.template.new' }"
        >
          {{ $t('new') }}
        </b-button>
        <c-permissions-button
          v-if="templateID && canGrant"
          :title="template.handle"
          :target="template.handle"
          :resource="'corteza::system:template/'+templateID"
          button-variant="light"
        >
          <font-awesome-icon :icon="['fas', 'lock']" />
          {{ $t('permissions') }}
        </c-permissions-button>
      </span>
      <c-corredor-manual-buttons
        ui-page="template/editor"
        ui-slot="toolbar"
        resource-type="system:template"
        default-variant="link"
        class="mr-1"
        @click="dispatchCortezaSystemTemplateEvent($event, { template })"
      />
    </c-content-header>

    <c-template-editor-info
      :template="template"
      :processing="info.processing"
      :success="info.success"
      :can-create="canCreate"
      @delete="onDelete"
      @submit="onInfoSubmit"
    />

    <c-template-editor-content
      v-if="template && template.templateID != '0'"
      class="mt-3"
      :template="template"
      :partials="partials"
      :processing="info.processing"
      :success="info.success"
      :can-create="canCreate"
      @submit="onInfoSubmit"
    />
  </b-container>
</template>

<script>
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CTemplateEditorInfo from 'corteza-webapp-admin/src/components/Template/CTemplateEditorInfo'
import CTemplateEditorContent from 'corteza-webapp-admin/src/components/Template/CTemplateEditorContent/Index'
import { system } from '@cortezaproject/corteza-js'
import { mapGetters } from 'vuex'

export default {
  components: {
    CTemplateEditorInfo,
    CTemplateEditorContent,
  },

  i18nOptions: {
    namespaces: 'system.templates',
    keyPrefix: 'editor',
  },

  mixins: [
    editorHelpers,
  ],

  props: {
    templateID: {
      type: String,
      required: false,
      default: undefined,
    },
  },

  data () {
    return {
      template: undefined,

      info: {
        processing: false,
        success: false,
      },

      partials: [],
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canCreate () {
      return this.can('system/', 'template.create')
    },

    canGrant () {
      return this.can('system/', 'grant')
    },

    title () {
      return this.templateID ? this.$t('title.edit') : this.$t('title.create')
    },
  },

  watch: {
    templateID: {
      immediate: true,
      handler () {
        this.fetchPartials()
        if (this.templateID) {
          this.fetchTemplate()
        } else {
          this.template = new system.Template()
        }
      },
    },
  },

  methods: {
    fetchTemplate () {
      this.incLoader()

      this.$SystemAPI.templateRead({ templateID: this.templateID })
        .then(t => {
          this.template = new system.Template(t)
        })
        .catch(this.toastErrorHandler(this.$t('notification:template.fetch.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    fetchPartials () {
      this.incLoader()

      this.$SystemAPI.templateList({ partial: true })
        .then(({ set: tt }) => {
          this.partials = tt.map(t => new system.Template(t))
        })
        .catch(this.toastErrorHandler(this.$t('notification:template.fetchPartials.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    onDelete () {
      this.incLoader()
      if (this.template.deletedAt) {
        this.$SystemAPI.templateUndelete({ templateID: this.templateID })
          .then(() => {
            this.fetchTemplate()

            this.toastSuccess(this.$t('notification:template.undelete.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:template.undelete.error')))
          .finally(() => {
            this.decLoader()
          })
      } else {
        this.$SystemAPI.templateDelete({ templateID: this.templateID })
          .then(() => {
            this.fetchTemplate()

            this.toastSuccess(this.$t('notification:template.delete.success'))
            this.$router.push({ name: 'system.template' })
          })
          .catch(this.toastErrorHandler(this.$t('notification:template.delete.error')))
          .finally(() => {
            this.decLoader()
          })
      }
    },

    onInfoSubmit (template) {
      this.incLoader()

      if (this.templateID) {
        this.$SystemAPI.templateUpdate(template)
          .then(template => {
            this.template = template

            this.toastSuccess(this.$t('notification:template.update.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:template.update.error')))
          .finally(() => {
            this.decLoader()
          })
      } else {
        this.$SystemAPI.templateCreate(template)
          .then(({ templateID }) => {
            this.animateSuccess('info')
            this.toastSuccess(this.$t('notification:template.create.success'))

            this.$router.push({ name: 'system.template.edit', params: { templateID } })
          })
          .catch(this.toastErrorHandler(this.$t('notification:template.create.error')))
          .finally(() => {
            this.decLoader()
          })
      }
    },
  },
}
</script>
