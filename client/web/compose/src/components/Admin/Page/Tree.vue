<template>
  <div>
    <sortable-tree
      v-if="list.length"
      :draggable="namespace.canCreatePage"
      :data="{ children: list }"
      tag="ul"
      mixin-parent-key="parent"
      class="list-group"
      @changePosition="handleChangePosition"
    >
      <template
        slot-scope="{item}"
      >
        <div
          v-if="item.pageID"
          no-gutters
          class="d-flex flex-wrap align-content-center justify-content-between"
        >
          <div
            class="px-2 flex-fill overflow-hidden text-truncate gap-1 rounded-lg"
            :class="{'grab': namespace.canCreatePage }"
          >
            {{ item.title }}
            <font-awesome-icon
              v-if="!item.visible && item.moduleID == '0'"
              v-b-tooltip.noninteractive.hover="{ title: $t('notVisible'), container: '#body' }"
              class="text-danger"
              :icon="['fas', 'eye-slash']"
            />
            <b-badge
              v-if="!isValid(item)"
              variant="danger"
              class="ml-1"
            >
              {{ $t('invalid') }}
            </b-badge>
          </div>

          <div class="actions px-2">
            <b-button-group
              v-if="item.canUpdatePage"
              size="sm"
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
                v-b-tooltip.noninteractive.hover="{ title: $t('tooltip.view'), container: '#body' }"
                data-test-id="button-page-view"
                variant="primary"
                :to="pageViewer(item)"
                class="d-flex align-items-center"
                style="margin-left:2px;"
              >
                <font-awesome-icon
                  :icon="['far', 'eye']"
                />
              </b-button>

              <b-button
                v-b-tooltip.noninteractive.hover="{ title: $t('tooltip.edit.page'), container: '#body' }"
                data-test-id="button-page-edit"
                variant="primary"
                :to="{name: 'admin.pages.edit', params: { pageID: item.pageID }}"
                class="d-flex align-items-center"
                style="margin-left:2px;"
              >
                <font-awesome-icon
                  :icon="['far', 'edit']"
                />
              </b-button>
            </b-button-group>

            <b-dropdown
              v-if="item.canGrant || namespace.canGrant"
              v-b-tooltip.noninteractive.hover="{ title: $t('permissions:resources.compose.page.tooltip'), container: '#body' }"
              data-test-id="dropdown-permissions"
              variant="extra-light"
              size="sm"
              class="permissions-dropdown ml-2"
            >
              <template #button-content>
                <font-awesome-icon :icon="['fas', 'lock']" />
              </template>

              <b-dropdown-item-button>
                <c-permissions-button
                  v-if="namespace.canGrant"
                  :title="item.title || item.handle || item.pageID"
                  :target="item.title || item.handle || item.pageID"
                  :resource="`corteza::compose:page/${namespace.namespaceID}/${item.pageID}`"
                  :button-label="$t('general:label.page')"
                  :show-button-icon="false"
                />
              </b-dropdown-item-button>

              <b-dropdown-item-button>
                <c-permissions-button
                  v-if="item.canGrant"
                  :title="item.title || item.handle || item.pageID"
                  :target="item.title || item.handle || item.pageID"
                  :resource="`corteza::compose:page-layout/${namespace.namespaceID}/${item.pageID}/*`"
                  :button-label="$t('general:label.pageLayout')"
                  :show-button-icon="false"
                  all-specific
                />
              </b-dropdown-item-button>
            </b-dropdown>

            <template v-if="item.canDeletePage">
              <b-dropdown
                v-if="hasChildren(item)"
                data-test-id="dropdown-delete"
                variant="outline-danger"
                :disabled="item.processingDelete"
                size="sm"
                toggle-class="border-0"
                class="ml-2"
              >
                <template #button-content>
                  <font-awesome-icon
                    v-if="!item.processingDelete"
                    :icon="['far', 'trash-alt']"
                  />

                  <b-spinner
                    v-else
                    small
                  />
                </template>

                <b-dropdown-item-button
                  data-test-id="dropdown-item-delete-update-parent-of-sub-pages"
                  @click="handleDeletePage(item, 'rebase')"
                >
                  {{ $t('delete.rebase') }}
                </b-dropdown-item-button>

                <b-dropdown-item-button
                  data-test-id="dropdown-item-delete-sub-pages"
                  @click="handleDeletePage(item, 'cascade')"
                >
                  {{ $t('delete.cascade') }}
                </b-dropdown-item-button>
              </b-dropdown>

              <c-input-confirm
                v-else
                show-icon
                size="md"
                button-class="px-2"
                :processing="item.processingDelete"
                class="ml-2"
                @confirmed="handleDeletePage(item, 'cascade')"
              />
            </template>
          </div>
        </div>
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
      deletePage: 'page/delete',
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
            this.$emit('reorder')

            setTimeout(() => {
              this.$store.dispatch('page/load', { namespaceID, clear: true, force: true })
            }, 1000)
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

    handleDeletePage (page, strategy = 'abort') {
      this.$set(page, 'processingDelete', true)

      this.deletePage({ ...page, strategy }).then(() => {
        setTimeout(() => {
          this.$emit('reorder')
          this.toastSuccess(this.$t('notification:page.deleted'))
        }, 300)
      })
        .catch(this.toastErrorHandler(this.$t('notification:page.deleteFailed')))
        .finally(() => {
          setTimeout(() => {
            this.$set(page, 'processingDelete', false)
          }, 300)
        })
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

    hasChildren (page) {
      return page.children && page.children.length > 0
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

<style lang="scss">
//!important usage to over-ride library styling
$input-height: 2.4rem;
$content-height: 2.4rem;
$blank-li-height: 1rem;
$padding: 0.5rem;
$border-color: var(--light);
$hover-color: var(--extra-light);
$dropping-color: var(--primary);

.page-name-input {
  height: $input-height;
}

.list-group {
  margin-right: $padding;

  .content {
    height: 0 !important;
  }

  .droper {
    background-color: $dropping-color !important;
    border-radius: 0.3rem !important;
    margin-left: 0 !important;

    .blank-li.first-child {
      background-color: $dropping-color !important;
    }
  }

  .draging {
    background-color: var(--white) !important;
  }

  ul {
    .content {
      height: 100% !important;
      min-height: $content-height !important;
      line-height: $content-height !important;
      background-color: var(--light);
      border-radius: 0.3rem !important;
      margin-left: 12px;

      .actions {
        display: none;
      }

      &:hover {
        background-color: $hover-color;

        .actions {
          display: block;
        }
      }
    }
  }

  li {
    white-space: nowrap;
    background: var(--white);

    &.blank-li {
      height: $blank-li-height !important;

      .sortable-tree {
        max-height: 100%;
      }

      &:nth-last-of-type(1)::before {
        border-left-color: var(--white) !important;
        height: 0;
      }
    }

    &::before {
      top: calc($content-height / -2) !important;
      border-left-color: var(--white) !important;
    }

    &::after {
      height: $content-height !important;
      top: calc($content-height / 2) !important;
      border-color: var(--white) !important;
    }

    &.parent-li:nth-last-child(2)::before {
      height: $content-height !important;
      top: calc($content-height / -2) !important;
    }
  }

  .parent-li {
    border-top: 1px solid $border-color;
    padding-bottom: $padding !important;
    padding-top: $padding !important;

    .exist-li,
    .blank-li {
      border-top: none;
      padding-top: 0 !important;
      padding-bottom: 0 !important;

      &::after {
        border-top: 2px solid $border-color !important;
        margin-left: 0;
      }

      &::before {
        border-left: 2px solid $border-color !important;
      }
    }

    &.blank-li {
      &::before {
        border-left: 2px solid $border-color !important;
      }
    }

    &.exist-li {
      &::before {
        border-color: var(--white) !important;
      }

      .parent-li {
        &.exist-li {
          &::before {
            border-color: $border-color !important;
          }
        }
      }
    }
  }
}

.pages-list-header {
  min-height: $content-height;
  background-color: var(--gray-200);
  margin-bottom: -1.8rem !important;
  border-bottom: 2px solid var(--light);
  z-index: 1;
}
</style>
