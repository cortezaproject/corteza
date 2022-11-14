<template>
  <div>
    <sortable-tree
      v-if="list.length"
      :draggable="namespace.canCreatePage"
      :data="{children:list}"
      tag="ul"
      mixin-parent-key="parent"
      class="list-group pb-3"
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
            class="text-right pr-2"
          >
            <router-link
              v-if="item.canUpdatePage"
              data-test-id="button-page-builder"
              :to="{name: 'admin.pages.builder', params: { pageID: item.pageID }}"
              class="btn btn-light mr-2"
            >
              {{ $t('block.general.label.pageBuilder') }}
            </router-link>
            <span class="view d-inline-block">
              <router-link
                v-if="!hideViewPageButton(item)"
                data-test-id="button-page-view"
                :to="pageViewer(item)"
                class="btn"
              >
                {{ $t('view') }}
              </router-link>
            </span>
            <span
              class="d-none d-md-inline-block edit text-left"
            >
              <router-link
                v-if="item.moduleID !== '0'"
                v-b-tooltip.hover.top
                data-test-id="button-module-edit"
                :title="moduleName(item)"
                class="btn text-primary"
                :to="{ name: 'admin.modules.edit', params: { moduleID: item.moduleID }}"
              >
                {{ $t('moduleEdit') }}
              </router-link>
              <router-link
                v-if="item.canUpdatePage && item.moduleID === '0'"
                :to="{name: 'admin.pages.edit', params: { pageID: item.pageID }}"
                data-test-id="button-page-edit"
                class="btn text-primary"
              >
                {{ $t('edit.edit') }}
              </router-link>

            </span>
            <c-permissions-button
              v-if="namespace.canGrant"
              :title="item.title"
              :target="item.title"
              :resource="`corteza::compose:page/${namespace.namespaceID}/${item.pageID}`"
              link
              class="btn px-2"
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

    pageViewer ({ pageID = NoID }) {
      return { name: 'page', params: { pageID } }
    },

    hideViewPageButton ({ blocks = {}, moduleID = NoID }) {
      return blocks && blocks.length >= 1 && moduleID !== NoID
    },

    handleChangePosition ({ beforeParent, data, afterParent }) {
      const { namespaceID } = this.namespace
      const beforeID = beforeParent.parent ? beforeParent.pageID : NoID
      const afterID = afterParent.parent ? afterParent.pageID : NoID

      const reorder = () => {
        const pageIDs = afterParent.children.map(p => p.pageID)
        if (pageIDs.length) {
          this.$ComposeAPI.pageReorder({ namespaceID, selfID: afterID, pageIDs: pageIDs }).then(() => {
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
$edit-width: 130px;
$view-width: 70px;

.edit {
  width: $edit-width;
}

.view {
  width: $view-width;
}

.grab {
  cursor: grab;
  z-index: 1;
}
</style>
