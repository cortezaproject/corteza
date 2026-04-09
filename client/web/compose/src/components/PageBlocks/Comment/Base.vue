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
      <p class="mb-0 my-3">
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
          class="flex-grow-1 py-2 px-1 overflow-auto"
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
            class="date-group d-flex flex-column gap-2 mt-2"
          >
            <div
              v-if="comments.length > 1"
              class="d-flex align-items-center justify-content-center gap-3 mx-2 text-muted gap-2"
            >
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
                :attachment-field="attachmentField"
                :reactions-field="reactionsField"
                :emoji-data="emojiData"
                :namespace="namespace"
                :show-header="ci === 0"
                :show-title="showTitle(comment)"
                :show-content="showContent(comment)"
                :highlighted="highlightedCommentId === comment.recordID"
                :current-user-i-d="currentUserID"
                :find-user-by-i-d="findUserByID"
                class="mb-1"
                @reply="replyToComment(comment)"
                @edit="onEditComment(comment, $event)"
                @react="onReact(comment, $event)"
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
          <p class="mb-0 my-3">
            {{ $t('comment.noComments') }}
          </p>
        </div>

        <section
          v-if="canAddRecord"
          class="d-flex flex-column bg-white border-top"
        >
          <div
            v-if="newRecord.replyTo"
            class="reply-to-container p-3"
          >
            <p class="text-muted">
              Replying to
            </p>

            <div class="position-relative">
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
            :placeholder="$t('comment.title.placeholder')"
          />

          <c-rich-text-input
            ref="richTextInput"
            v-model="newRecord.content"
            :placeholder="$t('comment.content.placeholder')"
            :labels="{
              urlPlaceholder: $t('content.urlPlaceholder'),
              ok: $t('content.ok'),
              emojiPicker: {
                search: $t('content.emojiPicker.search'),
                searchResults: $t('content.emojiPicker.searchResults'),
                frequentlyUsed: $t('content.emojiPicker.frequentlyUsed'),
                noResults: $t('content.emojiPicker.noResults'),
                quickReactions: $t('content.emojiPicker.quickReactions'),
              },
            }"
            min-body-height="4rem"
            max-body-height="10rem"
            body-class="overflow-auto"
            style="border: none !important;"
            @upload="handleFileUpload"
          />

          <c-uploader
            v-if="attachmentField"
            ref="uploader"
            :endpoint="fileUploadEndpoint"
            :form-data="uploaderFormData"
            :accepted-files="mimetypes"
            :max-filesize="maxSize"
            :max-files="attachmentField.isMulti ? undefined : 1"
            class="d-none"
            @upload="appendAttachment"
          />

          <list-loader
            v-if="attachmentField && newRecord.attachmentIDs.length"
            kind="record"
            :set.sync="newRecord.attachmentIDs"
            :namespace="namespace"
            :enable-order="attachmentField.isMulti"
            enable-delete
            mode="list"
            :hide-file-name="attachmentField.options.hideFileName"
            :preview-options="attachmentField.options"
            class="px-2"
          />

          <div class="d-flex align-items-center justify-content-end m-2 gap-1">
            <b-button
              v-if="attachmentField"
              v-b-tooltip.hover="{ title: $t('comment.tooltip.attach'), delay: { show: 300, hide: 0 } }"
              variant="outline-light"
              class="text-secondary border-0"
              @click="openFileUpload"
            >
              <font-awesome-icon :icon="['fas', 'paperclip']" />
            </b-button>

            <c-button-submit
              :text="$t('comment.submit')"
              :disabled="!isValid || isProcessing"
              :processing="submitting"
              @submit="submitComment()"
            />
          </div>
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
              :attachment-field="attachmentField"
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
import { evaluatePrefilter, getFieldFilter, isFieldInFilter } from 'corteza-webapp-compose/src/lib/record-filter'
import records from 'corteza-webapp-compose/src/mixins/records'
import users from 'corteza-webapp-compose/src/mixins/users'
import { mapGetters } from 'vuex'
import base from '../base'
import CommentItem from './Item.vue'
import CommentReply from './Reply.vue'
import ListLoader from 'corteza-webapp-compose/src/components/Public/Page/Attachment/ListLoader'
const { CRichTextInput, CUploader, emojiData } = components

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    CRichTextInput,
    CommentReply,
    CommentItem,
    ListLoader,
    CUploader,
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
        attachmentIDs: [],
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

      commentRefreshInterval: null,
      autoFetching: false,

      emojiData,
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
      pages: 'page/set',
      findUserByID: 'user/findByID',
      findRecordByID: 'record/findByID',
    }),

    lastCommentTimestamp () {
      if (this.comments.length === 0) {
        return null
      }

      const { messages = [] } = this.comments[this.comments.length - 1] || {}

      if (messages.length === 0) {
        return null
      }

      const { comments = [] } = messages[messages.length - 1] || {}

      if (comments.length === 0) {
        return null
      }

      const comment = comments[comments.length - 1]

      if (!comment) {
        return null
      }

      return comment.createdAt
    },

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

    attachmentField () {
      const { attachmentField } = this.options

      if (!attachmentField) {
        return undefined
      }

      const f = this.roModule.fields.find(f => f.name === attachmentField)
      if (!f) {
        return undefined
      }

      const af = compose.ModuleFieldMaker(f)
      af.options.mode = 'list'
      return af
    },

    reactionsField () {
      const { reactionsField } = this.options

      if (!reactionsField) {
        return undefined
      }

      return this.roModule.fields.find(f => f.name === reactionsField)
    },

    currentUserID () {
      return (this.$auth.user || {}).userID || ''
    },

    fileUploadEndpoint () {
      if (!this.attachmentField) {
        return undefined
      }

      const moduleID = this.moduleID
      const recordID = NoID
      const { namespaceID } = this.namespace

      return this.$ComposeAPI.baseURL + this.$ComposeAPI.recordUploadEndpoint({
        namespaceID,
        moduleID,
        recordID,
        fieldName: this.attachmentField.name,
      })
    },

    uploaderFormData () {
      if (!this.attachmentField) {
        return {}
      }

      return {
        fieldName: this.attachmentField.name,
      }
    },

    mimetypes () {
      if (!this.attachmentField) {
        return []
      }

      const a = (this.attachmentField.options.mimetypes || '').trim()
      if (!a) {
        return []
      }

      return a.split(',').map(p => p.trim())
    },

    maxSize () {
      if (!this.attachmentField) {
        return 100
      }

      return this.attachmentField.options.maxSize || 100
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
      return (!!this.newRecord.title || !!this.newRecord.content || this.newRecord.attachmentIDs.length) && !this.isNewRecord
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
    this.startAutoRefresh()
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

    startAutoRefresh () {
      this.commentRefreshInterval = setInterval(() => {
        // Skip auto-refresh if:
        // - Already fetching
        // - Currently submitting a comment
        // - Currently loading more messages via pagination
        if (this.autoFetching || this.submitting || this.loadingMore) {
          return
        }

        // For oldest-first mode, only auto-refresh if we're at the latest page (no nextPage)
        if (!this.showNewestFirst && this.filter.nextPage) {
          return
        }

        this.autoFetching = true

        Promise.all([
          this.loadNewComments(),
          this.loadUpdatedComments(),
        ]).finally(() => {
          this.autoFetching = false
        })
      }, 5000)
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

    loadNewComments () {
      const filter = [
        this.expandFilter(),
        this.lastCommentTimestamp ? `${getFieldFilter('createdAt', 'DateTime', this.lastCommentTimestamp, '>')}` : '',
      ].filter(Boolean).join(' AND ')

      // Remember if scroll was at bottom before updating comments
      const wasAtBottom = this.isScrollAtBottom()

      return this.fetchCommentRecords(this.roModule, filter, false).then(newComments => {
        this.comments = this.mergeMessageGroups(this.comments, newComments, true)

        // Auto-scroll to bottom if scroll was at bottom before the update
        if (wasAtBottom) {
          this.$nextTick(() => {
            this.scrollToLatest()
          })
        }
      })
    },

    loadUpdatedComments () {
      if (!this.reactionsField) return Promise.resolve()

      // Collect all record IDs currently displayed
      const recordIDs = []
      this.comments.forEach(dateGroup => {
        (dateGroup.messages || []).forEach(messageGroup => {
          (messageGroup.comments || []).forEach(comment => {
            recordIDs.push(comment.recordID)
          })
        })
      })

      if (!recordIDs.length) return Promise.resolve()

      // Re-fetch these specific records by ID
      const idFilter = recordIDs.map(id => `recordID = '${id}'`).join(' OR ')
      const filter = `(${idFilter})`

      return this.fetchCommentRecords(this.roModule, filter, false).then(updatedGroups => {
        this.updateExistingComments(updatedGroups)
      })
    },

    updateExistingComments (updatedGroups) {
      if (!updatedGroups || !updatedGroups.length) return

      const reactionsFieldName = this.reactionsField.name

      // Collect all updated comments into a map for fast lookup
      const updatedMap = {}
      updatedGroups.forEach(dateGroup => {
        (dateGroup.messages || []).forEach(messageGroup => {
          (messageGroup.comments || []).forEach(comment => {
            updatedMap[comment.recordID] = comment
          })
        })
      })

      // Walk through existing comments and replace only those whose reactions changed
      this.comments.forEach(dateGroup => {
        dateGroup.messages.forEach(messageGroup => {
          messageGroup.comments.forEach((comment, index) => {
            const updated = updatedMap[comment.recordID]
            if (!updated) return

            const oldReactions = comment.values[reactionsFieldName]
            const newReactions = updated.values[reactionsFieldName]
            if (oldReactions === newReactions) return

            // Preserve author and reply that were already resolved
            updated.author = updated.author || comment.author
            updated.reply = updated.reply || comment.reply
            messageGroup.comments.splice(index, 1, updated)
          })
        })
      })
    },

    mergeMessageGroups (existing, newGroups, showNewestFirst = this.showNewestFirst) {
      if (!existing.length || !newGroups.length) {
        return showNewestFirst ? [...existing, ...newGroups] : [...newGroups, ...existing]
      }

      const [existingGroup, newGroup] = showNewestFirst
        ? [existing[existing.length - 1], newGroups[0]]
        : [existing[0], newGroups[newGroups.length - 1]]

      if (existingGroup.date === newGroup.date) {
        if (showNewestFirst) {
          // Merge messages from newGroup into existingGroup
          newGroup.messages.forEach(newMessage => {
            const lastExistingMessage = existingGroup.messages[existingGroup.messages.length - 1]

            // If the last message in existing group has the same author, merge the comments
            if (lastExistingMessage && lastExistingMessage.authorId === newMessage.authorId) {
              lastExistingMessage.comments = [...lastExistingMessage.comments, ...newMessage.comments]
            } else {
              // Add as a new message group
              existingGroup.messages.push(newMessage)
            }
          })
        } else {
          existingGroup.messages = [...newGroup.messages, ...existingGroup.messages]
        }

        showNewestFirst ? newGroups.shift() : newGroups.pop()
      }

      return showNewestFirst
        ? [...existing, ...newGroups]
        : [...newGroups, ...existing]
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
        this.comments = this.mergeMessageGroups(this.comments, newGroups, !this.showNewestFirst)
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

    isScrollAtBottom () {
      const container = this.$refs.chatContainer
      if (!container) return false

      // Consider it "at bottom" if within 25px of the bottom
      const threshold = 25
      return container.scrollTop + container.clientHeight >= container.scrollHeight - threshold
    },

    handleFileUpload (files) {
      if (!this.attachmentField) {
        return
      }

      const uploader = this.$refs.uploader
      if (uploader && uploader.$refs.dropzone) {
        Array.from(files).forEach(file => {
          uploader.$refs.dropzone.addFile(file)
        })
      }
    },

    openFileUpload () {
      const uploader = this.$refs.uploader
      if (uploader && uploader.$refs.dropzone) {
        uploader.$refs.dropzone.dropzone.hiddenFileInput.click()
      }
    },

    appendAttachment ({ attachmentID } = {}) {
      if (attachmentID) {
        if (this.attachmentField && this.attachmentField.isMulti) {
          this.newRecord.attachmentIDs = [...this.newRecord.attachmentIDs, attachmentID]
        } else {
          this.newRecord.attachmentIDs = [attachmentID]
        }
      }
    },

    submitComment () {
      if (!this.isValid) {
        return
      }

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

      if (this.replyField && this.newRecord.replyTo) {
        record.values[this.replyField.name] = this.newRecord.replyTo.recordID
      }

      if (this.attachmentField && this.newRecord.attachmentIDs.length) {
        record.values[this.attachmentField.name] = this.attachmentField.isMulti ? this.newRecord.attachmentIDs : this.newRecord.attachmentIDs[0]
      }

      return this.$ComposeAPI.recordCreate(record).then(rec => {
        rec = new compose.Record(this.roModule, rec)

        this.newRecord.title = ''
        this.newRecord.content = ''
        this.newRecord.replyTo = null
        this.newRecord.attachmentIDs = []

        if (this.showNewestFirst) {
          return this.loadNewComments()
        } else {
          this.showNewestFirst = true
          this.filter.nextPage = ''

          return this.fetchCommentRecords(this.roModule, this.expandFilter()).then(groupedRecords => {
            this.comments = groupedRecords
          })
        }
      })
        .catch(this.toastErrorHandler(this.$t('notification:record.createFailed')))
        .finally(() => {
          this.submitting = false
          this.$nextTick(() => {
            this.scrollToLatest()
          })
        })
    },

    onEditComment (comment, { title, content }) {
      const record = new compose.Record(this.roModule, { ...comment })

      if (this.titleField) {
        record.values[this.titleField.name] = title
      }

      if (this.contentField) {
        record.values[this.contentField.name] = content
      }

      return this.$ComposeAPI.recordUpdate(record).then(rec => {
        const updatedRecord = new compose.Record(this.roModule, rec)
        updatedRecord.author = comment.author
        updatedRecord.reply = comment.reply

        this.comments.forEach(dateGroup => {
          dateGroup.messages.forEach(messageGroup => {
            const index = messageGroup.comments.findIndex(c => c.recordID === updatedRecord.recordID)
            if (index > -1) {
              messageGroup.comments.splice(index, 1, updatedRecord)
            }
          })
        })
      })
        .catch(this.toastErrorHandler(this.$t('notification:record.updateFailed')))
    },

    onReact (comment, emoji) {
      if (!this.reactionsField) return

      const record = new compose.Record(this.roModule, { ...comment })
      const fieldName = this.reactionsField.name
      let reactions = {}

      try {
        reactions = JSON.parse(record.values[fieldName] || '{}') || {}
      } catch {
        reactions = {}
      }

      const userID = this.currentUserID
      if (!userID) return

      // Toggle: add or remove current user
      if (!reactions[emoji]) {
        reactions[emoji] = []
      }

      const idx = reactions[emoji].indexOf(userID)
      if (idx > -1) {
        reactions[emoji].splice(idx, 1)
        if (reactions[emoji].length === 0) {
          delete reactions[emoji]
        }
      } else {
        reactions[emoji].push(userID)
      }

      record.values[fieldName] = JSON.stringify(reactions)

      return this.$ComposeAPI.recordUpdate(record).then(rec => {
        const updatedRecord = new compose.Record(this.roModule, rec)
        updatedRecord.author = comment.author
        updatedRecord.reply = comment.reply

        // Resolve any new user IDs from reactions
        this.resolveReactionUsers(updatedRecord)

        this.comments.forEach(dateGroup => {
          dateGroup.messages.forEach(messageGroup => {
            const index = messageGroup.comments.findIndex(c => c.recordID === updatedRecord.recordID)
            if (index > -1) {
              messageGroup.comments.splice(index, 1, updatedRecord)
            }
          })
        })
      })
        .catch(this.toastErrorHandler(this.$t('notification:record.updateFailed')))
    },

    resolveReactionUsers (record) {
      if (!this.reactionsField) return

      try {
        const reactions = JSON.parse(record.values[this.reactionsField.name] || '{}') || {}
        const userIDs = [...new Set(Object.values(reactions).flat())].filter(Boolean)
        if (userIDs.length) {
          this.$store.dispatch('user/resolveUsers', userIDs)
        }
      } catch {
        // ignore
      }
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

    async fetchCommentRecords (module, query, useNextPage = true) {
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

      if (useNextPage && this.filter.nextPage) {
        sort = ''
      }

      const { moduleID, namespaceID } = module

      const params = {
        namespaceID,
        moduleID,
        query,
        sort,
        limit: useNextPage ? this.filter.limit : 500,
        pageCursor: useNextPage ? this.filter.nextPage : '',
      }

      const { response, cancel } = this.$ComposeAPI.recordListCancellable(params)
      this.abortableRequests.push(cancel)

      return response().then(({ set = [], filter = {} }) => {
        if (useNextPage) {
          this.filter.nextPage = filter.nextPage || ''
        }

        const comments = set.map(r => new compose.Record(module, r))

        return Promise.all([
          this.fetchUsers([{ name: 'createdBy', kind: 'User', isSystem: true, isMulti: false }], comments),
          this.fetchReplyRecords(comments),
        ]).then(() => {
          // Resolve user IDs from reactions
          comments.forEach(c => this.resolveReactionUsers(c))
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

      this.$nextTick(() => {
        const rti = this.$refs.richTextInput
        if (rti) {
          if (typeof rti.focus === 'function') {
            rti.focus()
          } else if (rti.editor && typeof rti.editor.focus === 'function') {
            rti.editor.focus()
          }
        }
      })
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
        attachmentIDs: [],
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

      this.autoFetching = false
      if (this.commentRefreshInterval) {
        clearInterval(this.commentRefreshInterval)
        this.commentRefreshInterval = null
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
