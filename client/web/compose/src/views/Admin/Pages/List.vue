<template>
  <b-container
    fluid="xl"
    class="d-flex flex-column py-3"
  >
    <portal to="topbar-title">
      {{ $t('navigation.page') }}
    </portal>

    <b-card
      no-body
      class="shadow-sm h-100"
    >
      <b-card-header
        header-bg-variant="white"
        class="d-flex flex-column border-bottom gap-1"
      >
        <div class="d-flex align-items-stretch align-items-sm-center justify-content-between flex-column flex-sm-row gap-1">
          <div class="flex-grow-1">
            <b-input-group
              v-if="namespace.canCreatePage"
              class="h-100"
              style="min-width: 300px;"
            >
              <b-input
                id="name"
                v-model="page.title"
                data-test-id="input-name"
                required
                type="text"
                class="h-100"
                :placeholder="$t('newPlaceholder')"
              />
              <b-input-group-append>
                <b-button
                  data-test-id="button-create-page"
                  type="submit"
                  variant="primary"
                  @click="handleAddPageFormSubmit"
                >
                  {{ $t('createLabel') }}
                </b-button>
              </b-input-group-append>
            </b-input-group>
          </div>

          <div class="d-flex justify-content-sm-end flex-fill flex-grow-1">
            <b-dropdown
              v-if="namespace.canGrant"
              data-test-id="dropdown-permissions"
              size="lg"
              variant="light"
              class="d-flex align-items-center flex-sm-grow-0 flex-sm-shrink-0 flex-fill permissions-dropdown"
            >
              <template #button-content>
                <font-awesome-icon :icon="['fas', 'lock']" />
                <span>
                  {{ $t('label.permissions') }}
                </span>
              </template>

              <b-dropdown-item>
                <c-permissions-button
                  :resource="`corteza::compose:page/${namespace.namespaceID}/*`"
                  :button-label="$t('general:label.page')"
                  :show-button-icon="false"
                  button-variant="white text-left w-100"
                />
              </b-dropdown-item>

              <b-dropdown-item>
                <c-permissions-button
                  :resource="`corteza::compose:page-layout/${namespace.namespaceID}/*/*`"
                  :button-label="$t('general:label.pageLayout')"
                  :show-button-icon="false"
                  button-variant="white text-left w-100"
                />
              </b-dropdown-item>
            </b-dropdown>
          </div>
        </div>

        <span class="text-muted font-italic">
          {{ $t('instructions') }}
        </span>
      </b-card-header>

      <div
        v-if="processing"
        class="text-center text-muted m-5"
      >
        <div>
          <b-spinner
            class="align-middle m-2"
          />
        </div>
        {{ $t('loading') }}
      </div>

      <page-tree
        v-else
        v-model="tree"
        :namespace="namespace"
        class="card overflow-auto h-100"
        @reorder="handleReorder"
      />
    </b-card>
  </b-container>
</template>

<script>
import axios from 'axios'
import { mapActions } from 'vuex'
import PageTree from 'corteza-webapp-compose/src/components/Admin/Page/Tree'
import { compose } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'page',
  },

  name: 'PageList',

  components: {
    PageTree,
  },

  props: {
    namespace: {
      type: compose.Namespace,
      required: false,
      default: undefined,
    },
  },

  data () {
    return {
      tree: [],
      page: new compose.Page({ visible: true }),
      processing: false,
      abortableRequests: [],
    }
  },

  created () {
    this.loadTree()
  },

  beforeDestroy () {
    this.abortRequests()
    this.setDefaultValues()
  },

  methods: {
    ...mapActions({
      createPage: 'page/create',
      createPageLayout: 'pageLayout/create',
    }),

    loadTree () {
      this.processing = true
      const { namespaceID } = this.namespace

      const { response, cancel } = this.$ComposeAPI
        .pageTreeCancellable({ namespaceID })

      this.abortableRequests.push(cancel)

      response()
        .then((tree) => {
          this.tree = tree.map(p => new compose.Page(p))
        }).catch((e) => {
          if (!axios.isCancel(e)) {
            this.toastErrorHandler(this.$t('notification:page.loadFailed'))(e)
          }
        })
        .finally(() => {
          this.processing = false
        })
    },

    handleAddPageFormSubmit () {
      const { namespaceID } = this.namespace
      this.page.weight = this.tree.length
      this.createPage({ ...this.page, namespaceID }).then(({ pageID, title }) => {
        const pageLayout = new compose.PageLayout({ namespaceID, pageID, handle: 'primary', meta: { title } })
        return this.createPageLayout(pageLayout).then(() => {
          this.$router.push({ name: 'admin.pages.edit', params: { pageID } })
        })
      }).catch((e) => {
        if (!axios.isCancel(e)) {
          this.toastErrorHandler(this.$t('notification:page.saveFailed'))(e)
        }
      })
    },

    handleReorder () {
      this.loadTree()
    },

    setDefaultValues () {
      this.tree = []
      this.page = {}
      this.processing = false
      this.abortableRequests = []
    },

    abortRequests () {
      this.abortableRequests.forEach((cancel) => {
        cancel()
      })
    },
  },
}
</script>
