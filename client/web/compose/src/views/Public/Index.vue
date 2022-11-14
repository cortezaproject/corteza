<template>
  <div class="d-flex flex-column h-100 w-100">
    <portal to="topbar-title">
      {{ pageTitle }}
    </portal>

    <portal to="topbar-tools">
      <b-button-group
        v-if="page && page.canUpdatePage"
        size="sm"
        class="mr-1"
      >
        <b-button
          variant="primary"
          class="d-flex align-items-center"
          :to="pageBuilder"
        >
          {{ $t('label.pageBuilder') }}
          <font-awesome-icon
            :icon="['fas', 'cogs']"
            class="ml-2"
          />
        </b-button>
        <page-translator
          v-if="trPage"
          :page.sync="trPage"
          style="margin-left:2px;"
        />
        <b-button
          variant="primary"
          style="margin-left:2px;"
          class="d-flex align-items-center"
          :to="pageEditor"
        >
          <font-awesome-icon
            :icon="['far', 'edit']"
          />
        </b-button>
      </b-button-group>
    </portal>

    <div
      v-if="showSteps"
      class="d-flex flex-column m-5 vh-75"
    >
      <h1 class="display-3">
        {{ $t('label.welcome') }}
      </h1>
      <p class="lead">
        {{ $t('message.noPages') }}
        <span v-if="namespace.canManageNamespace">
          {{ $t('message.startBuilding') }}
        </span>
        <span v-else>
          {{ $t('message.notifyAdministrator') }}
        </span>
      </p>
      <b-container
        v-if="namespace.canManageNamespace"
        fluid="xl"
        class="align-items-center border-top steps"
      >
        <b-row
          align-v="center"
          class="text-center justify-content-between"
        >
          <b-col>
            <circle-step
              step-number="1"
              :done="hasModules"
            >
              <b-button
                v-if="!hasModules"
                data-test-id="button-module-create"
                :disabled="!namespace.canCreateModule"
                variant="outline-primary"
                size="lg"
                @click="createNewModule"
              >
                {{ $t('step.module.create') }}
              </b-button>
              <router-link
                v-else
                :to="{ name: 'admin.modules' }"
              >
                <b-button
                  data-test-id="button-module-view"
                  :disabled="!namespace.canManageNamespace"
                  variant="primary"
                  size="lg"
                >
                  {{ $t('step.module.view') }}
                </b-button>
              </router-link>
            </circle-step>
          </b-col>
          <b-col>
            <hr>
          </b-col>
          <b-col>
            <circle-step
              :done="hasCharts"
              :disabled="!hasModules"
              optional
            >
              <b-button
                v-if="!hasCharts"
                :disabled="!hasModules || !namespace.canCreateChart"
                variant="outline-primary"
                size="lg"
                @click="createNewChart"
              >
                {{ $t('step.chart.create') }}
              </b-button>
              <router-link
                v-else
                :to="{ name: 'admin.charts' }"
              >
                <b-button
                  :disabled="!namespace.canManageNamespace"
                  variant="primary"
                  size="lg"
                >
                  {{ $t('step.chart.view') }}
                </b-button>
              </router-link>
            </circle-step>
          </b-col>
          <b-col>
            <hr>
          </b-col>
          <b-col>
            <circle-step
              step-number="2"
              :done="hasPages"
              :disabled="!hasModules"
            >
              <b-button
                v-if="!hasPages"
                data-test-id="button-page-build"
                :disabled="!hasModules || !namespace.canCreatePage"
                variant="outline-primary"
                size="lg"
                @click="createNewPage"
              >
                {{ $t('step.page.create') }}
              </b-button>
              <router-link
                v-else
                :to="{ name: 'admin.pages' }"
              >
                <b-button
                  :disabled="!namespace.canManageNamespace"
                  data-test-id="button-page-view"
                  variant="primary"
                  size="lg"
                >
                  {{ $t('step.page.view') }}
                </b-button>
              </router-link>
            </circle-step>
          </b-col>
        </b-row>
      </b-container>
    </div>
    <router-view
      v-else
      class="flex-grow-1 overflow-auto"
      :namespace="namespace"
      :page="page"
    />
    <portal-target
      name="toolbar"
    />
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import CircleStep from 'corteza-webapp-compose/src/components/Common/CircleStep'
import PageTranslator from 'corteza-webapp-compose/src/components/Admin/Page/PageTranslator'
import { compose, NoID } from '@cortezaproject/corteza-js'

const pushContentAbove = 610
const demoPageHandle = 'demo_page'

export default {
  i18nOptions: {
    namespaces: 'onboarding',
  },

  name: 'PublicRoot',

  components: {
    CircleStep,
    PageTranslator,
  },

  props: {
    pageID: {
      type: String,
      required: false,
      default: '',
    },

    namespace: { // via router-view
      type: compose.Namespace,
      required: true,
    },
  },

  data () {
    return {
      navVisible: false,
      documentWidth: 0,
      loaded: false,
    }
  },

  computed: {
    ...mapGetters({
      modules: 'module/set',
      pages: 'page/set',
      charts: 'chart/set',
    }),

    routerViewClass () {
      return {
        'compose-content': true,
        padded: this.navVisible && this.canPushContent,
      }
    },

    canPushContent () {
      return this.documentWidth > pushContentAbove
    },

    trPage: {
      get () {
        return this.page.clone()
      },
      set (v) {
        this.updatePageSet(v)
      },
    },

    page () {
      return this.$store.getters['page/getByID'](this.pageID) || new compose.Page()
    },

    showSteps () {
      return !this.pageID && this.loaded
    },

    hasModules () {
      return !!this.modules.length
    },

    hasCharts () {
      return !!this.charts.length
    },

    hasPages () {
      return this.pages.filter(p => p.visible || p.handle === demoPageHandle).length > 0
    },

    pageTitle () {
      if (this.page.pageID !== NoID) {
        const { title = '', handle = '' } = this.page
        return title || handle || this.$t('navigation.noPageTitle')
      }

      return ''
    },

    pageEditorDisabled () {
      return this.page.moduleID !== NoID
    },

    pageEditor () {
      return { name: 'admin.pages.edit', params: { pageID: this.pageID } }
    },

    pageBuilder () {
      return { name: 'admin.pages.builder', params: { pageID: this.pageID } }
    },
  },

  watch: {
    pageID: {
      immediate: true,
      handler (pageID) {
        // If we redirect to page index, try to find & redirect to a first
        // available public page.
        if (!this.pageID) {
          const { pageID } = this.$store.getters['page/homePage'] || {}
          if (pageID) {
            // Use replace so we don't push to history stack
            this.$router.replace({ name: 'page', params: { pageID } })
          } else {
            this.loaded = true
          }
        }
      },
    },
  },

  created () {
    this.documentWidth = document.body.offsetWidth
    window.onresize = () => {
      this.documentWidth = document.body.offsetWidth
    }
  },

  methods: {
    ...mapActions({
      createModule: 'module/create',
      createPage: 'page/create',
      createChart: 'chart/create',
      updatePageSet: 'page/updateSet',
    }),

    createNewModule () {
      const newModule = new compose.Module({
        name: 'Demo Module',
        handle: 'demo_module',
        fields: [
          new compose.ModuleFieldString({ fieldID: '0', name: 'Sample' }),
        ],
      }, this.namespace)

      this.createModule(newModule).then((module) => {
        this.$router.push({ name: 'admin.modules.edit', params: { moduleID: module.moduleID } })
      }).catch(this.toastErrorHandler(this.$t('notification:module.createFailed')))
    },

    createNewChart () {
      const { namespaceID } = this.namespace
      const { moduleID = '' } = this.modules.find(m => m.moduleID) || {}
      const newChart = new compose.Chart({
        namespaceID,
        name: 'Demo Chart',
        handle: 'demo_chart',
        config: {
          reports: [{ moduleID }],
        },
      })

      this.createChart(newChart).then((chart) => {
        this.$router.push({ name: 'admin.charts.edit', params: { chartID: chart.chartID } })
      }).catch(this.toastErrorHandler(this.$t('notification:chart.createFailed')))
    },

    createNewPage () {
      const { namespaceID } = this.namespace
      const newPage = new compose.Page({
        namespaceID,
        title: 'Demo Page',
        handle: demoPageHandle,
        blocks: [],
      })

      this.createPage(newPage).then((page) => {
        this.$router.push({ name: 'admin.pages.builder', params: { pageID: page.pageID } })
      }).catch(this.toastErrorHandler(this.$t('notification:page.saveFailed')))
    },
  },
}
</script>
<style lang="scss" scoped>
.steps {
  padding: 0;
  padding-top: 20vh;
}
</style>
