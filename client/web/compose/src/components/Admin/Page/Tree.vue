<template>
  <div>
    <sortable-tree
      v-if="list.length"
      :draggable="namespace.canCreatePage"
      :data="{children:list}"
      tag="ul"
      mixin-parent-key="parent"
      class="list-group"
      @changePosition="handleChangePosition"
    >
      <template
        slot-scope="{item}"
      >
        <b-row
          v-if="item.pageID"
          no-gutters
          class="wrap d-flex pr-2"
        >
          <b-col
            cols="12"
            xl="6"
            lg="5"
            class="flex-fill pl-2 overflow-hidden"
            :class="{'grab': namespace.canCreatePage}"
          >
            {{ item.title }}
            <span
              v-if="!item.visible && item.moduleID == '0'"
              class="text-danger"
            >
              <font-awesome-icon
                :icon="['fas', 'eye-slash']"
                :title="$t('notVisible')"
              />
            </span>
            <b-badge
              v-if="!isValid(item)"
              variant="danger"
            >
              {{ $t('invalid') }}
            </b-badge>
          </b-col>
          <b-col
            cols="12"
            xl="6"
            lg="7"
            class="text-right"
          >
            <b-button-group
              v-if="item.canUpdatePage"
              size="sm"
              class="mr-1"
            >
              <b-button
                data-test-id="button-page-builder"
                variant="primary"
                size="sm"
                :to="{name: 'admin.pages.builder', params: { pageID: item.pageID }}"
              >
                {{ $t('block.general.label.pageBuilder') }}
                <font-awesome-icon
                  :icon="['fas', 'tools']"
                  class="ml-2"
                />
              </b-button>

              <b-button
                data-test-id="button-page-view"
                variant="primary"
                :title="$t('tooltip.view')"
                :to="pageViewer(item)"
                class="d-flex align-items-center"
                style="margin-left:2px;"
              >
                <font-awesome-icon
                  :icon="['far', 'eye']"
                />
              </b-button>

              <b-button
                data-test-id="button-page-edit"
                variant="primary"
                :title="$t('tooltip.edit.page')"
                :to="{name: 'admin.pages.edit', params: { pageID: item.pageID }}"
                class="d-flex align-items-center"
                style="margin-left:2px;"
              >
                <font-awesome-icon
                  :icon="['fas', 'pen']"
                />
              </b-button>
            </b-button-group>

            <c-permissions-button
              v-if="item.canGrant"
              :title="item.title || item.handle || item.pageID"
              :target="item.title || item.handle || item.pageID"
              :resource="`corteza::compose:page/${namespace.namespaceID}/${item.pageID}`"
              :tooltip="$t('permissions:resources.compose.page.tooltip')"
              button-variant="outline-light"
              class="text-dark d-print-none border-0"
            />
          </b-col>
        </b-row>
      </template>
    </sortable-tree>

    <div
      v-else
      class="text-center mt-5 mb-4 pb-1"
    >
      {{ $t('noPages') }}
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import SortableTree from 'vue-sortable-tree'
import { compose, NoID } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'page',
  },

  name: 'PageTree',

  components: {
    SortableTree,
  },

  props: {
    namespace: {
      type: compose.Namespace,
      required: true,
    },

    value: {
      type: Array,
      required: true,
    },

    parentID: {
      type: String,
      default: NoID,
    },

    level: {
      type: Number,
      default: 0,
    },
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
    }),

    list: {
      get () {
        return this.value
      },

      set (pages) {
        this.$emit('input', pages.filter(p => p))
      },
    },
  },

  methods: {
    ...mapActions({
      updatePage: 'page/update',
    }),

    moduleName ({ moduleID }) {
      if (moduleID === NoID) {
        return ''
      }

      return (this.getModuleByID(moduleID) || {}).name
    },

    pageViewer ({ pageID = NoID, moduleID = NoID }) {
      const name = moduleID !== NoID ? 'page.record.create' : 'page'
      return { name, params: { pageID } }
    },

    handleChangePosition ({ beforeParent, data, afterParent }) {
      const { namespaceID } = this.namespace
      const beforeID = beforeParent.parent ? beforeParent.pageID : NoID
      const afterID = afterParent.parent ? afterParent.pageID : NoID

      const reorder = () => {
        const pageIDs = afterParent.children.map(p => p.pageID)
        if (pageIDs.length) {
          this.$ComposeAPI.pageReorder({ namespaceID, selfID: afterID, pageIDs }).then(() => {
            return this.$store.dispatch('page/load', { namespaceID, clear: true, force: true })
          }).then(() => {
            this.toastSuccess(this.$t('reordered'))
            this.$emit('reorder')
          })
            .catch(this.toastErrorHandler(this.$t('pageMoveFailed')))
        }
      }

      if (beforeID !== afterID) {
        // Page moved to a different parent
        data.weight = 1
        data.selfID = afterID
        data.namespaceID = namespaceID

        this.updatePage(data).then(() => {
          reorder()
        }).catch(this.toastErrorHandler(this.$t('pageMoveFailed')))
      } else {
        reorder()
      }
    },

    /**
     * Validates page, returns true if there are no problems with it
     *
     * @param {compose.Page} page
     * @returns {boolean}
     */
    isValid (page) {
      if (typeof page.validate === 'function') {
        return page.validate().length === 0
      }

      return true
    },
  },
}
</script>
<style lang="scss" scoped>
.grab {
  cursor: grab;
  z-index: 1;
}
</style>
