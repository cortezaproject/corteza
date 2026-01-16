<template>
  <div class="h-100 d-flex flex-column">
    <div class="overflow-auto flex-grow-1 h-100">
      <div
        v-if="loading"
        class="d-flex justify-content-center p-5"
      >
        <b-spinner variant="primary" />
      </div>

      <div
        v-if="drafts.length > 0 && !loading"
        class="d-flex align-items-center justify-content-between p-2 pl-3 border-bottom bg-light"
      >
        <div class="text-secondary small font-weight-bold">
          {{ $t('count', { count: drafts.length }) }}
        </div>

        <div class="d-flex align-items-center gap-1">
          <c-input-select
            v-model="sortOrder"
            :options="sortOptions"
            :reduce="option => option.value"
            :clearable="false"
            :searchable="false"
            size="sm"
            style="width: 10rem;"
            class="border-0"
          >
            <template #option="option">
              <div class="d-flex align-items-center gap-1">
                <font-awesome-icon
                  v-if="option && option.icon"
                  :icon="['fas', option.icon]"
                  class="text-secondary"
                />
                {{ option.text }}
              </div>
            </template>
            <template #selected-option="option">
              <div
                v-if="option"
                class="d-flex align-items-center gap-1 text-secondary"
              >
                <font-awesome-icon
                  v-if="option.icon"
                  :icon="['fas', option.icon]"
                />
                {{ option.text }}
              </div>
            </template>
          </c-input-select>

          <c-input-confirm
            v-b-tooltip.hover.bottomleft="{ title: $t('deleteAllDrafts'), delay: { show: 500, hide: 0 } }"
            show-icon
            @confirmed="onClearAll"
          />
        </div>
      </div>

      <b-list-group v-if="sortedDrafts.length > 0 && !loading">
        <draft-item
          v-for="draft in sortedDrafts"
          :key="draft.revision.changeID"
          :draft="draft"
          :namespace="namespaces[draft.revision.resource.split('/')[1]]"
          :module="modules[draft.revision.resource.split('/')[2]]"
          :active="$route.query.draftID === String(draft.revision.changeID)"
          @click="onDraftClick(draft)"
          @view="onDraftClick(draft, true)"
          @delete="onDeleteDraft"
        />
      </b-list-group>

      <div
        v-else-if="!loading"
        class="text-center p-5"
      >
        <font-awesome-icon
          :icon="['far', 'edit']"
          class="text-secondary mb-3"
          size="3x"
        />
        <p class="text-secondary">
          {{ $t('empty') }}
        </p>
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions, mapMutations } from 'vuex'
import { compose } from '@cortezaproject/corteza-js'
import { components } from '@cortezaproject/corteza-vue'
import DraftItem from 'corteza-webapp-compose/src/components/Drafts/DraftItem.vue'

const { CInputConfirm, CInputSelect } = components

export default {
  i18nOptions: {
    namespaces: ['drafts', 'notifications', 'general'],
  },

  components: {
    DraftItem,
    CInputConfirm,
    CInputSelect,
  },

  data () {
    return {
      modules: {},
      namespaces: {},
      sortOrder: 'desc',
    }
  },

  computed: {
    ...mapGetters({
      drafts: 'drafts/getAllDrafts',
      loading: 'drafts/isLoading',
    }),

    sortedDrafts () {
      return [...this.drafts].sort((a, b) => {
        const aTime = new Date(a.revision.timestamp)
        const bTime = new Date(b.revision.timestamp)
        return this.sortOrder === 'desc' ? bTime - aTime : aTime - bTime
      })
    },

    sortOptions () {
      return [
        { value: 'desc', text: this.$t('newestFirst'), icon: 'sort-amount-down', label: this.$t('newestFirst') },
        { value: 'asc', text: this.$t('oldestFirst'), icon: 'sort-amount-up', label: this.$t('oldestFirst') },
      ]
    },
  },

  watch: {
    drafts: {
      immediate: true,
      handler (drafts) {
        this.fetchMetadata(drafts)
      },
    },
  },

  methods: {
    ...mapActions({
      removeDraft: 'drafts/removeDraft',
      clearDrafts: 'drafts/clearDrafts',
    }),

    ...mapMutations({
      setVisible: 'drafts/setVisible',
    }),

    async fetchMetadata (drafts) {
      try {
        for (const draft of drafts) {
          const parts = draft.revision.resource.split('/')
          const namespaceID = parts[1]
          const moduleID = parts[2]

          if (!this.namespaces[namespaceID]) {
            const ns = await this.$ComposeAPI.namespaceRead({ namespaceID })
            if (ns) {
              this.$set(this.namespaces, namespaceID, new compose.Namespace(ns))
            }
          }

          if (!this.modules[moduleID]) {
            const mod = await this.$ComposeAPI.moduleRead({ namespaceID, moduleID })
            if (mod) {
              this.$set(this.modules, moduleID, new compose.Module(mod))
            }
          }
        }
      } catch (e) {
        console.error('Failed to fetch metadata:', e)
      }
    },

    async onDraftClick (draft, view = false) {
      const { revision } = draft

      // Resource format: compose:record/{namespaceID}/{moduleID}/{recordID}
      const parts = revision.resource.split('/')
      if (parts.length < 4) return

      const namespaceID = parts[1]
      const moduleID = parts[2]
      const recordID = parts[3] === '0' ? undefined : parts[3]

      let pageID

      try {
        const pages = await this.$ComposeAPI.pageList({ namespaceID, moduleID })
        const recordPage = (pages.set || []).find(p => p.moduleID === moduleID)
        if (recordPage) {
          pageID = recordPage.pageID
        }
      } catch (e) {
        console.error('Failed to fetch page metadata:', e)
      }

      if (!pageID) {
        this.toastDanger(this.$t('notifications:recordRedirectError'))
        return
      }

      const isCompose = this.$router.app.$options.name === 'compose'
      const ns = this.namespaces[namespaceID]
      const slug = ns ? ns.slug || ns.namespaceID : namespaceID

      if (!isCompose) {
        const u = new URL(window.location)
        let url = `${u.origin}/compose/ns/${slug}/pages/${pageID}/record/`
        if (recordID) {
          url += recordID
          if (!view) {
            url += '/edit'
          }
        }
        if (!view) {
          url += `?draftID=${revision.changeID}`
        }

        window.location.assign(url)
        return
      }

      const route = {
        name: view ? 'page.record' : (recordID ? 'page.record.edit' : 'page.record.create'),
        params: {
          slug,
          pageID,
          recordID,
        },
        query: view
          ? {}
          : {
              draftID: revision.changeID,
            },
      }

      this.$router.push(route).catch(err => {
        console.error('Draft navigation failed:', err)
      })
    },

    onDeleteDraft ({ revision }) {
      this.removeDraft({ changeID: revision.changeID })
        .then(() => {
          this.toastSuccess(this.$t('deleted'))
        })
        .catch(() => {
          this.toastDanger(this.$t('deleteError'))
        })
    },

    onClearAll () {
      this.clearDrafts()
        .then(() => {
          this.toastSuccess(this.$t('allDeleted'))
        })
        .catch(() => {
          this.toastDanger(this.$t('deleteError'))
        })
    },
  },
}
</script>
