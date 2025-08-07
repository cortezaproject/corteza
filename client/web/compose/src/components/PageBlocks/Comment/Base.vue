<template>
  <wrap
    v-bind="$props"
    :scrollable-body="false"
    v-on="$listeners"
    @refreshBlock="refresh"
  >
    <div
      v-if="!isBlockConfigured"
      class="d-flex h-100 align-items-center justify-content-center"
    >
      <p class="mb-0">
        {{ $t('noConfiguration') }}
      </p>
    </div>

    <template v-else-if="roModule">
      <div class="d-flex flex-column h-100">
        <div
          v-if="isProcessing"
          class="d-flex align-items-center justify-content-center h-100"
        >
          <b-spinner />
        </div>

        <section
          v-else-if="comments.length"
          ref="chatContainer"
          class="flex-grow-1 py-3 px-1 overflow-auto"
        >
          <div
            v-if="showNewestFirst && hasNextPage"
            class="text-center"
          >
            <c-button-submit
              :text="$t('comment.load.older')"
              :processing="loadingMore"
              variant="extra-light"
              class="mb-1"
              @submit="loadMoreMessages"
            />
          </div>

          <div
            v-for="dateGroup in comments"
            :key="dateGroup.date"
            class="date-group d-flex flex-column gap-1 mt-2"
          >
            <div class="d-flex align-items-center justify-content-center gap-3 mx-2 text-muted gap-2">
              <hr class="flex-grow-1 m-0">
              <span>{{ dateGroup.date }}</span>
              <hr class="flex-grow-1 m-0">
            </div>

            <div
              v-for="(messageGroup, index) in dateGroup.messages"
              :key="index"
              class="message-group"
            >
              <comment-item
                v-for="(comment, ci) in messageGroup.comments"
                :id="`comment-${comment.recordID}`"
                :key="comment.recordID"
                :comment="comment"
                :title-field="titleField"
                :content-field="contentField"
                :namespace="namespace"
                :show-header="ci === 0"
                :show-title="showTitle(comment)"
                :show-content="showContent(comment)"
                :highlighted="highlightedCommentId === comment.recordID"
                @reply="replyToComment(comment)"
                @reply-click="handleReplyClick"
                @mouseleave="resetHighlightedComment(comment.recordID)"
              />
            </div>
          </div>

          <div
            v-if="!showNewestFirst && hasNextPage"
            class="text-center"
          >
            <c-button-submit
              :text="$t('comment.load.newer')"
              :processing="loadingMore"
              variant="extra-light"
              class="mt-1"
              @submit="loadMoreMessages"
            />
          </div>
        </section>

        <div
          v-else
          class="d-flex align-items-center justify-content-center h-100"
        >
          <p class="mb-0">
            {{ $t('comment.noComments') }}
          </p>
        </div>

        <section
          v-if="canAddRecord"
          class="d-flex flex-column px-3 py-3 bg-white border-top"
        >
          <div
            v-if="newRecord.replyTo"
            class="reply-to-container"
          >
            <p class="text-muted">
              Replying to
            </p>

            <div class="position-relative mb-3">
              <div class="reply-to-toolbox">
                <b-button
                  variant="extra-light"
                  size="sm"
                  class="py-1"
                  @click="newRecord.replyTo = null"
                >
                  <font-awesome-icon :icon="['fas', 'times']" />
                </b-button>
              </div>

              <comment-reply
                :reply="newRecord.replyTo"
                :title-field="titleField"
                :content-field="contentField"
                :namespace="namespace"
                @click.native="handleReplyClick(newRecord.replyTo.recordID)"
              />
            </div>
          </div>

          <b-form-input
            v-if="titleField"
            v-model="newRecord.title"
            class="mb-2"
            :placeholder="$t('comment.titleInput')"
          />

          <c-rich-text-input
            v-model="newRecord.content"
            :placeholder="$t('comment.contentInput')"
            :labels="{
              urlPlaceholder: $t('content.urlPlaceholder'),
              ok: $t('content.ok'),
            }"
            min-body-height="3rem"
            max-body-height="10rem"
            body-class="overflow-auto"
          />

          <c-button-submit
            :text="$t('comment.submit')"
            :disabled="!isValid || isProcessing"
            :processing="submitting"
            class="ml-auto mt-2"
            @submit="createNewRecord()"
          />
        </section>

        <b-modal
          id="comment-modal"
          v-model="replyModal.show"
          size="xl"
          centered
          scrollable
          hide-header
          hide-footer
          body-class="p-2"
        >
          <div
            v-if="!replyModal.comment"
            class="d-flex align-items-center justify-content-center p-3"
          >
            <b-spinner />
          </div>

          <div v-else>
            <div class="d-flex align-items-center justify-content-center gap-3 mx-2 text-muted gap-2">
              <hr class="flex-grow-1 m-0">
              <span>{{ getFormattedDate((replyModal.comment || {}).createdAt) }}</span>
              <hr class="flex-grow-1 m-0">
            </div>

            <comment-item
              :comment="replyModal.comment"
              :title-field="titleField"
              :content-field="contentField"
              :namespace="namespace"
              :show-time-always="true"
              :show-title="showTitle(replyModal.comment)"
              :show-content="showContent(replyModal.comment)"
              :highlighted="false"
              :disable-hover="true"
              @reply="replyToComment(replyModal.comment)"
              @reply-click="openReplyInModal"
            />
          </div>
        </b-modal>
      </div>
    </template>
  </wrap>
</template>

<script>
import { NoID, compose, fmt } from '@cortezaproject/corteza-js'
import { components } from '@cortezaproject/corteza-vue'
import axios from 'axios'
import { evaluatePrefilter, isFieldInFilter } from 'corteza-webapp-compose/src/lib/record-filter'
import records from 'corteza-webapp-compose/src/mixins/records'
import users from 'corteza-webapp-compose/src/mixins/users'
import { mapGetters } from 'vuex'
import base from '../base'
import CommentReply from './Reply.vue'
import CommentItem from './Item.vue'
const { CRichTextInput } = components

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    CRichTextInput,
    CommentReply,
    CommentItem,
  },

  extends: base,

  mixins: [
    users,
    records,
  ],

  data () {
    return {
      filter: {
        sort: '',
        filter: '',
        limit: 50,
        pageCursor: '',
        prevPage: '',
        nextPage: '',
      },

      comments: [],
      newRecord: {
        title: '',
        content: '',
        replyTo: null,
      },

      submitting: false,
      loadingMore: false,
      abortableRequests: [],

      showNewestFirst: true,
      highlightedCommentId: null,

      replyModal: {
        show: false,
        comment: null,
      },
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
      pages: 'page/set',
      findUserByID: 'user/findByID',
      findRecordByID: 'record/findByID',
    }),

    roModule () {
      return this.getModuleByID(this.moduleID)
    },

    moduleID () {
      return this.options.moduleID
    },

    titleField () {
      const { titleField } = this.options

      if (!titleField) {
        return undefined
      }

      return this.roModule.fields.find(f => f.name === titleField)
    },

    contentField () {
      const { contentField } = this.options

      if (!contentField) {
        return undefined
      }

      return this.roModule.fields.find(f => f.name === contentField)
    },

    referenceField () {
      const { referenceField } = this.options

      if (!referenceField) {
        return undefined
      }

      return this.roModule.fields.find(f => f.name === referenceField) || {}
    },

    replyField () {
      const { replyField } = this.options

      if (!replyField) {
        return undefined
      }

      return this.roModule.fields.find(f => f.name === replyField) || {}
    },

    canAddRecord () {
      return this.roModule && this.roModule.canCreateRecord
    },

    isValid () {
      return !!this.newRecord.content && !this.isNewRecord
    },

    isNewRecord () {
      if (this.record) {
        return this.record.recordID === NoID
      }
      return false
    },

    reference () {
      if (this.record) {
        return this.record.recordID !== NoID ? this.record.recordID : NoID
      }

      return NoID
    },

    isBlockConfigured () {
      return !!this.contentField
    },

    hasPrevPage () {
      return this.filter.prevPage
    },

    hasNextPage () {
      return this.filter.nextPage
    },
  },

  watch: {
    'record.recordID': {
      immediate: true,
      handler () {
        this.showNewestFirst = this.options.sortDirection === 'asc'
        this.refresh()
      },
    },

    options: {
      deep: true,
      handler () {
        this.showNewestFirst = this.options.sortDirection === 'asc'
        this.refresh()
      },
    },
  },

  created () {
    this.refreshBlock(this.refresh)
  },

  mounted () {
    this.createEvents()
  },

  beforeDestroy () {
    this.abortRequests()
    this.destroyEvents()
    this.setDefaultValues()
  },

  methods: {
    createEvents () {
      this.$root.$on('module-records-updated', this.refreshOnRelatedRecordsUpdate)
      this.$root.$on('record-field-change', this.refetchOnPrefilterValueChange)
      this.$root.$on('refetch-records', this.refresh)
    },

    refetchOnPrefilterValueChange ({ fieldName }) {
      const { filter } = this.options

      if (isFieldInFilter(fieldName, filter)) {
        this.refresh()
      }
    },

    getFormattedDateTime (date) {
      return fmt.fullDateTime(date)
    },

    getFormattedDate (timestamp) {
      const date = new Date(timestamp)
      const today = new Date()
      const yesterday = new Date()
      yesterday.setDate(yesterday.getDate() - 1)

      const compareDate = new Date(date.getFullYear(), date.getMonth(), date.getDate())
      const compareToday = new Date(today.getFullYear(), today.getMonth(), today.getDate())
      const compareYesterday = new Date(yesterday.getFullYear(), yesterday.getMonth(), yesterday.getDate())

      if (compareDate.getTime() === compareToday.getTime()) {
        return this.$t('comment.today')
      } else if (compareDate.getTime() === compareYesterday.getTime()) {
        return this.$t('comment.yesterday')
      } else {
        return fmt.date(timestamp, { dateStyle: 'long' })
      }
    },

    getFormattedTime (date) {
      return fmt.time(date)
    },

    mergeMessageGroups (existing, newGroups) {
      if (!existing.length || !newGroups.length) {
        return this.showNewestFirst ? [...newGroups, ...existing] : [...existing, ...newGroups]
      }

      const [existingGroup, newGroup] = this.showNewestFirst
        ? [existing[0], newGroups[newGroups.length - 1]]
        : [existing[existing.length - 1], newGroups[0]]

      if (existingGroup.date === newGroup.date) {
        existingGroup.messages = this.showNewestFirst
          ? [...newGroup.messages, ...existingGroup.messages]
          : [...existingGroup.messages, ...newGroup.messages]

        this.showNewestFirst ? newGroups.pop() : newGroups.shift()
      }

      return this.showNewestFirst
        ? [...newGroups, ...existing]
        : [...existing, ...newGroups]
    },

    getAuthor (userID) {
      const user = this.findUserByID(userID) || {}
      const name = user.name || user.handle || user.email || ''

      let initials = '?'
      if (name) {
        const words = name.trim().split(/\s+/)
        initials = words.length === 1
          ? words[0].substring(0, 2).toUpperCase()
          : words.slice(0, 2).map(w => w.charAt(0).toUpperCase()).join('')
      }

      return {
        name,
        initials,
        isCurrentUser: Boolean(this.$auth.user && this.$auth.user.userID === userID),
      }
    },

    loadMoreMessages () {
      this.loadingMore = true

      const container = this.$refs.chatContainer
      const currentScrollTop = container ? container.scrollTop : 0
      const currentScrollHeight = container ? container.scrollHeight : 0

      this.fetchCommentRecords(this.roModule, this.expandFilter()).then(newGroups => {
        this.comments = this.mergeMessageGroups(this.comments, newGroups)
      }).finally(() => {
        this.$nextTick(() => {
          if (container && this.showNewestFirst) {
            const newScrollHeight = container.scrollHeight
            const heightDifference = newScrollHeight - currentScrollHeight

            container.scrollTop = currentScrollTop + heightDifference
          }
        })

        this.loadingMore = false
      })
    },

    refreshOnRelatedRecordsUpdate ({ moduleID } = {}) {
      if (this.options.moduleID === moduleID) {
        this.refresh()
      }
    },

    refresh () {
      if (!this.options.moduleID) {
      // Make sure block is properly configured
        throw Error(this.$t('record.moduleOrPageNotSet'))
      }

      if (this.roModule && this.contentField) {
        this.processing = true

        this.filter.nextPage = ''

        return this.fetchCommentRecords(this.roModule, this.expandFilter())
          .then(groupedRecords => {
            if (this.showNewestFirst) {
              this.comments = groupedRecords.sort((a, b) => {
                return new Date(a.date) - new Date(b.date)
              })
            } else {
              this.comments = groupedRecords.sort((a, b) => {
                return new Date(b.date) - new Date(a.date)
              })
            }
          })
          .catch(e => {
            console.error(e)
          })
          .finally(() => {
            setTimeout(() => {
              this.processing = false
              this.$nextTick(() => {
                this.scrollToPosition()
              })
            }, 300)
          })
      }
    },

    scrollToPosition () {
      const container = this.$refs.chatContainer

      if (!container) {
        return
      }

      if (this.showNewestFirst) {
        container.scrollTop = container.scrollHeight
      } else {
        container.scrollTop = 0
      }
    },

    scrollToLatest () {
      const container = this.$refs.chatContainer
      if (!container) return
      container.scrollTop = container.scrollHeight
    },

    createNewRecord () {
      this.submitting = true

      const record = new compose.Record(this.roModule)

      if (this.titleField) {
        record.values[this.titleField.name] = this.newRecord.title
      }

      if (this.contentField) {
        record.values[this.contentField.name] = this.newRecord.content
      }

      if (this.referenceField) {
        record.values[this.referenceField.name] = this.reference
      }

      if (this.replyField) {
        record.values[this.replyField.name] = this.newRecord.replyTo.recordID
      }

      this.$ComposeAPI.recordCreate(record).then(rec => {
        rec = new compose.Record(this.roModule, rec)

        this.newRecord.title = ''
        this.newRecord.content = ''
        this.newRecord.replyTo = null

        this.showNewestFirst = true
        this.filter.nextPage = ''

        return this.fetchCommentRecords(this.roModule, this.expandFilter()).then(groupedRecords => {
          this.comments = groupedRecords
        })
      })
        .catch(this.toastErrorHandler(this.$t('notification:record.createFailed')))
        .finally(() => {
          this.submitting = false
          this.$nextTick(() => {
            this.scrollToLatest()
          })
        })
    },

    expandFilter () {
      /* eslint-disable no-template-curly-in-string */
      if (!this.record) {
        // If there is no current record and we are using recordID/ownerID variable in (pre)filter
        // we should disable the block
        if ((this.options.filter || '').includes('${record')) {
          throw Error(this.$t('record.invalidRecordVar'))
        }

        if ((this.options.filter || '').includes('${ownerID}')) {
          throw Error(this.$t('record.invalidOwnerVar'))
        }
      }

      if (this.options.filter) {
        try {
          return evaluatePrefilter(this.options.filter, {
            record: this.record,
            user: this.$auth.user || {},
            recordID: (this.record || {}).recordID || NoID,
            ownerID: (this.record || {}).ownedBy || NoID,
            userID: (this.$auth.user || {}).userID || NoID,
          })
        } catch (e) {
          return e
        }
      }

      return ''
    },

    async fetchCommentRecords (module, query) {
      if (module.moduleID !== this.options.moduleID) {
        throw Error('Module incompatible, module mismatch')
      }

      if (this.referenceField) {
        if (query.length) {
          query += ' AND '
        }
        query += `${this.referenceField.name} = '${this.reference}' `
      }

      let sort = this.showNewestFirst ? 'createdAt DESC' : 'createdAt ASC'

      if (this.filter.nextPage) {
        sort = ''
      }

      const { moduleID, namespaceID } = module

      const params = {
        namespaceID,
        moduleID,
        query,
        sort,
        limit: this.filter.limit,
        pageCursor: this.filter.nextPage,
      }

      const { response, cancel } = this.$ComposeAPI.recordListCancellable(params)
      this.abortableRequests.push(cancel)

      return response().then(({ set = [], filter = {} }) => {
        this.filter.nextPage = filter.nextPage || ''

        const comments = set.map(r => new compose.Record(module, r))

        return Promise.all([
          this.fetchUsers([{ name: 'createdBy', kind: 'User', isSystem: true, isMulti: false }], comments),
          this.fetchReplyRecords(comments),
        ]).then(() => {
          const groups = {}

          if (this.showNewestFirst) {
            comments.reverse()
          }

          comments.forEach(comment => {
            const date = this.getFormattedDate(comment.createdAt)
            const authorId = comment.createdBy
            comment.reply = this.getReplyComment(comment)
            comment.author = this.getAuthor(authorId)

            if (!groups[date]) {
              groups[date] = {
                date,
                messages: [],
              }
            }

            const lastMessage = groups[date].messages[groups[date].messages.length - 1]

            if (lastMessage && lastMessage.authorId === authorId) {
              lastMessage.comments.push(comment)
            } else {
              groups[date].messages.push({
                authorId,
                comments: [comment],
              })
            }
          })

          return Object.values(groups)
        })
      }).catch(e => {
        if (!axios.isCancel(e)) {
          console.error(e)
        }
        return []
      })
    },

    fetchReplyRecords (records) {
      if (!this.replyField || records.length === 0) {
        return Promise.resolve()
      }

      const fields = [this.replyField]

      return this.fetchRecords(this.namespace.namespaceID, fields, records)
    },

    replyToComment (comment) {
      this.newRecord.replyTo = comment

      this.replyModal.show = false
    },

    handleReplyClick (recordID) {
      const commentElement = document.getElementById(`comment-${recordID}`)

      if (commentElement) {
        commentElement.scrollIntoView({
          behavior: 'smooth',
          block: 'center',
        })

        this.highlightedCommentId = recordID
      } else {
        this.openReplyInModal(recordID)
      }
    },

    openReplyInModal (recordID) {
      const { namespaceID, moduleID } = this.roModule || {}
      if (!namespaceID || !moduleID) return

      this.replyModal.show = true
      this.replyModal.comment = null

      let comment = this.findRecordByID(recordID)

      if (!comment) {
        return
      }

      comment = new compose.Record(this.roModule, comment)

      this.fetchReplyRecords([comment])
        .then(() => {
          comment.reply = this.getReplyComment(comment)
          comment.author = this.getAuthor(comment.createdBy)

          this.replyModal.comment = comment
        })
        .catch(e => {
          this.replyModal.show = false
          this.toastErrorHandler(this.$t('notification:record.loadFailed'))(e)
        })
    },

    resetHighlightedComment (recordID) {
      if (this.highlightedCommentId === recordID) {
        this.highlightedCommentId = null
      }
    },

    showTitle (comment) {
      return Boolean(this.titleField && this.titleField.canReadRecordValue && comment.values[this.titleField.name])
    },

    showContent (comment) {
      return Boolean(this.contentField && this.contentField.canReadRecordValue && comment.values[this.contentField.name])
    },

    showReply (comment) {
      return Boolean(this.replyField && this.replyField.canReadRecordValue && comment.values[this.replyField.name])
    },

    getReplyComment (comment) {
      if (!this.showReply(comment)) {
        return null
      }

      let replyRecord = this.findRecordByID(comment.values[this.replyField.name])

      if (!replyRecord) {
        return null
      }

      replyRecord = new compose.Record(this.roModule, replyRecord)
      replyRecord.author = this.getAuthor(replyRecord.createdBy)

      return replyRecord
    },

    setDefaultValues () {
      this.filter = {
        sort: '',
        filter: '',
        limit: 50,
        pageCursor: '',
        prevPage: '',
        nextPage: '',
      }
      this.comments = []
      this.newRecord = {
        title: '',
        content: '',
        replyTo: null,
      }
      this.abortableRequests = []
      this.submitting = false
      this.loadingMore = false
      this.showNewestFirst = true
      this.highlightedCommentId = null
      this.replyModal = {
        show: false,
        comment: null,
      }
    },

    abortRequests () {
      this.abortableRequests.forEach((cancel) => {
        cancel()
      })
    },

    destroyEvents () {
      this.$root.$off('module-records-updated', this.refreshOnRelatedRecordsUpdate)
      this.$root.$off('record-field-change', this.refetchOnPrefilterValueChange)
      this.$root.$off('refetch-records', this.refresh)
    },
  },
}
</script>

<style lang="scss" scoped>
.reply-to-container {
  .reply-to-toolbox {
    position: absolute;
    top: 0;
    right: 0;
    opacity: 0;
    transition: opacity 0.2s ease;
    z-index: 1;
  }

  &:hover {
    .reply-to-toolbox {
      opacity: 1;
    }
  }
}
</style>
