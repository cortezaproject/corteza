<template>
  <PageBlock :block="block" @refreshBlock="refresh">
    <!-- Not configured -->
    <div v-if="!isConfigured" class="flex items-center justify-center h-full">
      <p class="text-muted-color m-3">
        {{ $t('block.general.noConfiguration') }}
      </p>
    </div>

    <!-- Configured -->
    <template v-else-if="roModule">
      <div class="flex flex-col h-full">
        <!-- Loading -->
        <div v-if="processing" class="flex items-center justify-center h-full">
          <ProgressSpinner style="width: 24px; height: 24px" />
        </div>

        <!-- Comments list -->
        <section
          v-else-if="comments.length"
          ref="chatContainer"
          class="flex-1 py-2 px-1 overflow-auto"
        >
          <!-- Load older (newest-first mode) -->
          <div v-if="showNewestFirst && hasNextPage" class="text-center mb-1">
            <Button
              :label="$t('block.comment.load.older')"
              :loading="loadingMore"
              size="small"
              severity="secondary"
              text
              @click="loadMoreMessages"
            />
          </div>

          <!-- Date groups -->
          <div v-for="dateGroup in comments" :key="dateGroup.date" class="flex flex-col gap-2 mt-2">
            <!-- Date separator -->
            <div
              v-if="comments.length > 1"
              class="flex items-center justify-center gap-3 mx-2 text-muted-color"
            >
              <Divider class="flex-1" />
              <span class="text-sm whitespace-nowrap">{{ dateGroup.date }}</span>
              <Divider class="flex-1" />
            </div>

            <!-- Message groups (by author) -->
            <div v-for="(messageGroup, mgi) in dateGroup.messages" :key="mgi">
              <CommentItem
                v-for="(comment, ci) in messageGroup.comments"
                :id="`comment-${comment.recordID}`"
                :key="comment.recordID"
                :comment="comment"
                :title-field="titleField"
                :content-field="contentField"
                :attachment-field="attachmentField"
                :reactions-field="reactionsField"
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

          <!-- Load newer (oldest-first mode) -->
          <div v-if="!showNewestFirst && hasNextPage" class="text-center mt-1">
            <Button
              :label="$t('block.comment.load.newer')"
              :loading="loadingMore"
              size="small"
              severity="secondary"
              text
              @click="loadMoreMessages"
            />
          </div>
        </section>

        <!-- No comments -->
        <div v-else class="flex items-center justify-center h-full">
          <p class="text-muted-color m-3">
            {{ $t('block.comment.noComments') }}
          </p>
        </div>

        <!-- Input section -->
        <section v-if="canAddRecord" class="flex flex-col bg-surface border-t">
          <!-- Reply preview -->
          <div v-if="newComment.replyTo" class="reply-to-container p-3">
            <p class="text-muted-color">
              {{ $t('block.comment.replyingTo') }}
            </p>

            <div class="relative">
              <div class="reply-to-close">
                <Button
                  icon="pi pi-times"
                  text
                  size="small"
                  severity="secondary"
                  @click="newComment.replyTo = null"
                />
              </div>

              <CommentReply
                :reply="newComment.replyTo"
                :title-field="titleField"
                :content-field="contentField"
                @click="handleReplyClick(newComment.replyTo.recordID)"
              />
            </div>
          </div>

          <!-- Title input -->
          <InputText
            v-if="titleField"
            v-model="newComment.title"
            :placeholder="$t('block.comment.title.placeholder')"
            class="w-full mb-2"
          />

          <!-- Content input -->
          <CRichTextInput
            ref="contentInput"
            v-model="newComment.content"
            :placeholder="$t('block.comment.content.placeholder')"
            min-body-height="4rem"
            max-body-height="10rem"
            body-class="overflow-auto"
            submit-on-enter
            @upload="handleFileUpload"
            @submit="onComposerSubmit"
          />

          <!-- Attachment previews -->
          <div v-if="attachmentField && newComment.attachmentIDs.length" class="flex flex-wrap gap-2 px-2 py-1">
            <div
              v-for="(attID, idx) in newComment.attachmentIDs"
              :key="attID"
              class="flex items-center gap-1 bg-highlight rounded-lg px-2 py-1 text-sm"
            >
              <i class="pi pi-file text-muted-color" />
              <span class="text-muted-color">{{ $t('block.comment.attachment.file') }} {{ idx + 1 }}</span>
              <Button
                icon="pi pi-times"
                text
                size="small"
                severity="secondary"
                class="p-0 w-5 h-5"
                @click="removeAttachment(idx)"
              />
            </div>
          </div>

          <!-- Submit row -->
          <div class="flex items-center justify-end m-2 gap-1">
            <Button
              v-if="attachmentField"
              v-tooltip.top="{ value: $t('block.comment.tooltip.attach'), showDelay: 300 }"
              icon="pi pi-paperclip"
              text
              severity="secondary"
              @click="openFileUpload"
            />

            <input
              v-if="attachmentField"
              ref="fileInput"
              type="file"
              :multiple="attachmentField.isMulti"
              class="hidden"
              @change="onFileSelected"
            />

            <Button
              :label="$t('block.comment.submit')"
              :disabled="!isValid || submitting"
              :loading="submitting"
              size="small"
              @click="submitComment"
            />
          </div>
        </section>

        <!-- Reply modal for off-screen comments -->
        <Dialog
          v-model:visible="replyModal.show"
          modal
          :header="$t('block.comment.replyModalTitle')"
          :style="{ width: '50rem' }"
          :breakpoints="{ '960px': '75vw', '641px': '90vw' }"
        >
          <div v-if="!replyModal.comment" class="flex items-center justify-center p-6">
            <ProgressSpinner style="width: 24px; height: 24px" />
          </div>

          <div v-else>
            <CommentItem
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
        </Dialog>
      </div>
    </template>
  </PageBlock>
</template>

<script setup>
import { ref, computed, watch, inject, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { compose } from '@planetcrust/human-js'
import { useModuleStore } from '@planetcrust/human-vue'
import { useUserStore } from '@planetcrust/human-vue'
import { useRecordStore } from '@planetcrust/human-vue'
import PageBlock from './PageBlock.vue'
import CommentItem from './Comment/CommentItem.vue'
import CommentReply from './Comment/CommentReply.vue'
import { components } from '@planetcrust/human-vue'
import { useI18n } from 'vue-i18n'
import { evaluatePrefilter, getFieldFilter } from '../../../lib/record-filter'

const { CRichTextInput } = components

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const $ComposeAPI = inject('$ComposeAPI')
const $Auth = inject('$Auth')
const $eventBus = inject('$eventBus', null)
const { t } = useI18n()
const moduleStore = useModuleStore()
const userStore = useUserStore()
const recordStore = useRecordStore()

const chatContainer = ref(null)
const contentInput = ref(null)
const fileInput = ref(null)
const processing = ref(false)
const submitting = ref(false)
const loadingMore = ref(false)
const comments = ref([])
const highlightedCommentId = ref(null)
const abortableRequests = ref([])

const filter = ref({
  limit: 50,
  nextPage: '',
})

const newComment = ref({
  title: '',
  content: '',
  replyTo: null,
  attachmentIDs: [],
})

const replyModal = ref({
  show: false,
  comment: null,
})

let refreshInterval = null
let autoFetching = false

const options = computed(() => props.block.options || {})
const showNewestFirst = computed(() => options.value.sortDirection === 'asc')
const hasNextPage = computed(() => !!filter.value.nextPage)

const currentUserID = computed(() => ($Auth?.user || {}).userID || '')

const roModule = computed(() => {
  if (!options.value.moduleID) return null
  return moduleStore.getByID(options.value.moduleID)
})

const titleField = computed(() => {
  if (!options.value.titleField || !roModule.value) return undefined
  return roModule.value.fields.find(f => f.name === options.value.titleField)
})

const contentField = computed(() => {
  if (!options.value.contentField || !roModule.value) return undefined
  return roModule.value.fields.find(f => f.name === options.value.contentField)
})

const referenceField = computed(() => {
  if (!options.value.referenceField || !roModule.value) return undefined
  return roModule.value.fields.find(f => f.name === options.value.referenceField)
})

const replyField = computed(() => {
  if (!options.value.replyField || !roModule.value) return undefined
  return roModule.value.fields.find(f => f.name === options.value.replyField)
})

const attachmentField = computed(() => {
  if (!options.value.attachmentField || !roModule.value) return undefined
  return roModule.value.fields.find(f => f.name === options.value.attachmentField)
})

const reactionsField = computed(() => {
  if (!options.value.reactionsField || !roModule.value) return undefined
  return roModule.value.fields.find(f => f.name === options.value.reactionsField)
})

const canAddRecord = computed(() => roModule.value?.canCreateRecord)

const isValid = computed(() =>
  (!!newComment.value.title || !!newComment.value.content || newComment.value.attachmentIDs.length > 0),
)

const isConfigured = computed(() => !!contentField.value)

const reference = computed(() => {
  if (props.record?.recordID && props.record.recordID !== '0') {
    return props.record.recordID
  }
  return '0'
})

// ---- Date formatting ----

function getFormattedDate(timestamp) {
  const date = new Date(timestamp)
  const today = new Date()
  const yesterday = new Date()
  yesterday.setDate(yesterday.getDate() - 1)

  const cmp = d => new Date(d.getFullYear(), d.getMonth(), d.getDate()).getTime()

  if (cmp(date) === cmp(today)) {
    return t('block.comment.today')
  } else if (cmp(date) === cmp(yesterday)) {
    return t('block.comment.yesterday')
  } else {
    return date.toLocaleDateString(undefined, { dateStyle: 'long' })
  }
}

// ---- Author helpers ----

function getAuthor(userID) {
  const user = userStore.findByID(userID)
  const name = user?.name || user?.handle || user?.email || ''

  let initials = '?'
  if (name) {
    const words = name.trim().split(/\s+/)
    initials =
      words.length === 1
        ? words[0].substring(0, 2).toUpperCase()
        : words
            .slice(0, 2)
            .map(w => w.charAt(0).toUpperCase())
            .join('')
  }

  return {
    name,
    initials,
    isCurrentUser: Boolean($Auth?.user && $Auth.user.userID === userID),
  }
}

function findUserByID(userID) {
  return userStore.findByID(userID)
}

// ---- Fetch comments ----

async function fetchCommentRecords(query, useNextPage = true) {
  if (!roModule.value || !$ComposeAPI) return []

  const mod = roModule.value

  // Add reference filter
  if (referenceField.value && reference.value !== '0') {
    const refFilter = `${referenceField.value.name} = '${reference.value}'`
    query = query ? `${query} AND ${refFilter}` : refFilter
  }

  let sort = showNewestFirst.value ? 'createdAt DESC' : 'createdAt ASC'

  if (useNextPage && filter.value.nextPage) {
    sort = ''
  }

  const params = {
    namespaceID: props.namespace.namespaceID,
    moduleID: mod.moduleID,
    query,
    sort,
    limit: useNextPage ? filter.value.limit : 500,
    pageCursor: useNextPage ? filter.value.nextPage : '',
  }

  try {
    const { response, cancel } = $ComposeAPI.recordListCancellable(params)
    abortableRequests.value.push(cancel)

    const { set = [], filter: paging = {} } = await response()

    if (useNextPage) {
      filter.value.nextPage = paging.nextPage || ''
    }

    const records = set.map(r => new compose.Record(mod, r))

    // Resolve users for author info
    const userIDs = [...new Set(records.map(r => r.createdBy).filter(Boolean))]
    if (userIDs.length) {
      await userStore.resolveUsers(userIDs).catch(() => {})
    }

    // Resolve reply records
    await fetchReplyRecords(records)

    // Resolve reaction user IDs
    records.forEach(c => resolveReactionUsers(c))

    // Group by date then by author
    if (showNewestFirst.value) {
      records.reverse()
    }

    const groups = {}
    records.forEach(comment => {
      const date = getFormattedDate(comment.createdAt)
      const authorId = comment.createdBy
      comment.author = getAuthor(authorId)
      comment.reply = getReplyComment(comment)

      if (!groups[date]) {
        groups[date] = { date, messages: [] }
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
  } catch (e) {
    // Don't log cancelled requests
    if (e?.message !== 'canceled' && e?.code !== 'ERR_CANCELED') {
      console.error('Failed to fetch comments:', e)
    }
    return []
  }
}

// ---- Reply resolution ----

async function fetchReplyRecords(records) {
  if (!replyField.value || records.length === 0) return

  // Collect all reply record IDs
  const replyIDs = records
    .map(r => r.values[replyField.value.name])
    .filter(Boolean)
    .filter(id => id !== '0')

  if (replyIDs.length === 0) return

  const uniqueIDs = [...new Set(replyIDs)]

  // Fetch each reply record into the record store
  const mod = roModule.value
  await Promise.all(
    uniqueIDs.map(recordID =>
      recordStore.findByID({
        namespaceID: props.namespace.namespaceID,
        moduleID: mod.moduleID,
        recordID,
      }).catch(() => null),
    ),
  )
}

function getReplyComment(comment) {
  if (!replyField.value) return null

  const replyID = comment.values[replyField.value.name]
  if (!replyID || replyID === '0') return null

  let replyRecord = recordStore.getByID(replyID)
  if (!replyRecord) return null

  replyRecord = new compose.Record(roModule.value, replyRecord)
  replyRecord.author = getAuthor(replyRecord.createdBy)

  return replyRecord
}

function showTitle(comment) {
  return Boolean(titleField.value && titleField.value.canReadRecordValue && comment.values[titleField.value.name])
}

function showContent(comment) {
  return Boolean(contentField.value && contentField.value.canReadRecordValue && comment.values[contentField.value.name])
}

// ---- Merge message groups for auto-refresh ----

function mergeMessageGroups(existing, newGroups) {
  if (!existing.length || !newGroups.length) {
    return showNewestFirst.value ? [...existing, ...newGroups] : [...newGroups, ...existing]
  }

  const existingGroup = showNewestFirst.value ? existing[existing.length - 1] : existing[0]
  const newGroup = showNewestFirst.value ? newGroups[0] : newGroups[newGroups.length - 1]

  if (existingGroup.date === newGroup.date) {
    if (showNewestFirst.value) {
      newGroup.messages.forEach(newMessage => {
        const lastExisting = existingGroup.messages[existingGroup.messages.length - 1]
        if (lastExisting && lastExisting.authorId === newMessage.authorId) {
          lastExisting.comments = [...lastExisting.comments, ...newMessage.comments]
        } else {
          existingGroup.messages.push(newMessage)
        }
      })
      newGroups.shift()
    } else {
      existingGroup.messages = [...newGroup.messages, ...existingGroup.messages]
      newGroups.pop()
    }
  }

  return showNewestFirst.value ? [...existing, ...newGroups] : [...newGroups, ...existing]
}

// ---- Last comment timestamp for auto-refresh ----

const lastCommentTimestamp = computed(() => {
  if (!comments.value.length) return null
  const lastGroup = comments.value[comments.value.length - 1]
  if (!lastGroup?.messages?.length) return null
  const lastMsgGroup = lastGroup.messages[lastGroup.messages.length - 1]
  if (!lastMsgGroup?.comments?.length) return null
  return lastMsgGroup.comments[lastMsgGroup.comments.length - 1]?.createdAt || null
})

// ---- Load new comments (auto-refresh) ----

async function loadNewComments() {
  const filter = [
    expandFilter(),
    lastCommentTimestamp.value ? `${getFieldFilter('createdAt', 'DateTime', lastCommentTimestamp.value, '>')}` : '',
  ].filter(Boolean).join(' AND ')

  const wasAtBottom = isScrollAtBottom()
  const newGroups = await fetchCommentRecords(filter, false)
  comments.value = mergeMessageGroups(comments.value, newGroups)

  if (wasAtBottom) {
    nextTick(() => scrollToLatest())
  }
}

// ---- Refresh ----

async function refresh() {
  if (!options.value.moduleID || !roModule.value || !contentField.value) return

  const isInitialLoad = !comments.value.length
  if (isInitialLoad) {
    processing.value = true
  }
  filter.value.nextPage = ''

  try {
    const groupedRecords = await fetchCommentRecords(expandFilter())

    if (showNewestFirst.value) {
      comments.value = groupedRecords.sort((a, b) => new Date(a.date) - new Date(b.date))
    } else {
      comments.value = groupedRecords.sort((a, b) => new Date(b.date) - new Date(a.date))
    }
  } catch (e) {
    console.error('Comment refresh error:', e)
  } finally {
    if (isInitialLoad) {
      setTimeout(() => {
        processing.value = false
        nextTick(() => scrollToPosition())
      }, 300)
    }
  }
}

// ---- Load more (pagination) ----

async function loadMoreMessages() {
  loadingMore.value = true
  const container = chatContainer.value
  const prevScrollTop = container ? container.scrollTop : 0
  const prevScrollHeight = container ? container.scrollHeight : 0

  try {
    const newGroups = await fetchCommentRecords(expandFilter())
    comments.value = mergeMessageGroups(comments.value, newGroups)
  } finally {
    nextTick(() => {
      if (container && showNewestFirst.value) {
        const heightDiff = container.scrollHeight - prevScrollHeight
        container.scrollTop = prevScrollTop + heightDiff
      }
    })
    loadingMore.value = false
  }
}

// ---- Submit ----

function onComposerSubmit() {
  if (submitting.value) return
  submitComment()
}

async function submitComment() {
  if (!isValid.value || !roModule.value || submitting.value) return

  submitting.value = true

  try {
    const record = new compose.Record(roModule.value)

    if (titleField.value) {
      record.values[titleField.value.name] = newComment.value.title
    }
    if (contentField.value) {
      record.values[contentField.value.name] = newComment.value.content
    }
    if (referenceField.value) {
      record.values[referenceField.value.name] = reference.value
    }
    if (replyField.value && newComment.value.replyTo) {
      record.values[replyField.value.name] = newComment.value.replyTo.recordID
    }
    if (attachmentField.value && newComment.value.attachmentIDs.length) {
      record.values[attachmentField.value.name] = attachmentField.value.isMulti
        ? newComment.value.attachmentIDs
        : newComment.value.attachmentIDs[0]
    }

    await $ComposeAPI.recordCreate(record)

    newComment.value.title = ''
    newComment.value.content = ''
    newComment.value.replyTo = null
    newComment.value.attachmentIDs = []
    contentInput.value?.clear()

    if (showNewestFirst.value) {
      await loadNewComments()
    } else {
      await refresh()
    }
  } catch (e) {
    console.error('Failed to submit comment:', e)
  } finally {
    submitting.value = false
    nextTick(() => scrollToLatest())
  }
}

// ---- Edit ----

async function onEditComment(comment, { title, content }) {
  const mod = roModule.value
  if (!mod) return

  try {
    const record = new compose.Record(mod, { ...comment })

    if (titleField.value) {
      record.values[titleField.value.name] = title
    }
    if (contentField.value) {
      record.values[contentField.value.name] = content
    }

    const updatedRaw = await $ComposeAPI.recordUpdate(record)
    const updatedRecord = new compose.Record(mod, updatedRaw)
    updatedRecord.author = comment.author
    updatedRecord.reply = comment.reply

    // Update in-place
    comments.value.forEach(dateGroup => {
      dateGroup.messages.forEach(messageGroup => {
        const idx = messageGroup.comments.findIndex(c => c.recordID === updatedRecord.recordID)
        if (idx > -1) {
          messageGroup.comments.splice(idx, 1, updatedRecord)
        }
      })
    })
  } catch (e) {
    console.error('Failed to update comment:', e)
  }
}

// ---- Reactions ----

async function onReact(comment, emoji) {
  if (!reactionsField.value) return

  const mod = roModule.value
  if (!mod) return

  try {
    const record = new compose.Record(mod, { ...comment })
    const fieldName = reactionsField.value.name
    let reactions = {}

    try {
      reactions = JSON.parse(record.values[fieldName] || '{}') || {}
    } catch {
      reactions = {}
    }

    const userID = currentUserID.value
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

    const updatedRaw = await $ComposeAPI.recordUpdate(record)
    const updatedRecord = new compose.Record(mod, updatedRaw)
    updatedRecord.author = comment.author
    updatedRecord.reply = comment.reply

    resolveReactionUsers(updatedRecord)

    comments.value.forEach(dateGroup => {
      dateGroup.messages.forEach(messageGroup => {
        const idx = messageGroup.comments.findIndex(c => c.recordID === updatedRecord.recordID)
        if (idx > -1) {
          messageGroup.comments.splice(idx, 1, updatedRecord)
        }
      })
    })
  } catch (e) {
    console.error('Failed to update reaction:', e)
  }
}

function resolveReactionUsers(record) {
  if (!reactionsField.value) return

  try {
    const reactions = JSON.parse(record.values[reactionsField.value.name] || '{}') || {}
    const userIDs = [...new Set(Object.values(reactions).flat())].filter(Boolean)
    if (userIDs.length) {
      userStore.resolveUsers(userIDs).catch(() => {})
    }
  } catch {
    // ignore
  }
}

// ---- File attachments ----

function openFileUpload() {
  fileInput.value?.click()
}

async function onFileSelected(event) {
  const files = event.target.files
  if (!files || files.length === 0) return

  for (const file of files) {
    await uploadFile(file)
  }

  // Reset file input
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

function handleFileUpload(files) {
  if (!attachmentField.value || !files) return
  Array.from(files).forEach(file => uploadFile(file))
}

async function uploadFile(file) {
  if (!attachmentField.value || !$ComposeAPI) return

  try {
    const url = $ComposeAPI.recordUploadEndpoint({
      namespaceID: props.namespace.namespaceID,
      moduleID: roModule.value.moduleID,
    })

    const formData = new FormData()
    formData.append('fieldName', attachmentField.value.name)
    formData.append('upload', file, file.name)

    const { data } = await $ComposeAPI
      .api()
      .post(url, formData, { headers: { 'Content-Type': undefined } })

    if (data?.error) throw new Error(data.error)
    const attachment = data?.response ?? data
    if (!attachment?.attachmentID) throw new Error('Upload failed: no attachmentID')

    if (attachmentField.value.isMulti) {
      newComment.value.attachmentIDs = [...newComment.value.attachmentIDs, attachment.attachmentID]
    } else {
      newComment.value.attachmentIDs = [attachment.attachmentID]
    }
  } catch (e) {
    console.error('Failed to upload file:', e)
  }
}

function removeAttachment(index) {
  newComment.value.attachmentIDs.splice(index, 1)
}

// ---- Reply ----

function replyToComment(comment) {
  newComment.value.replyTo = comment
  replyModal.value.show = false
  nextTick(() => {
    contentInput.value?.$el?.focus?.()
  })
}

function handleReplyClick(recordID) {
  const el = document.getElementById(`comment-${recordID}`)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'center' })
    highlightedCommentId.value = recordID
  } else {
    openReplyInModal(recordID)
  }
}

async function openReplyInModal(recordID) {
  if (!roModule.value) return

  replyModal.value.show = true
  replyModal.value.comment = null

  try {
    let comment = recordStore.getByID(recordID)

    if (!comment) {
      comment = await recordStore.findByID({
        namespaceID: props.namespace.namespaceID,
        moduleID: roModule.value.moduleID,
        recordID,
      })
    }

    if (!comment) {
      replyModal.value.show = false
      return
    }

    comment = new compose.Record(roModule.value, comment)
    await fetchReplyRecords([comment])
    comment.reply = getReplyComment(comment)
    comment.author = getAuthor(comment.createdBy)

    replyModal.value.comment = comment
  } catch (e) {
    console.error('Failed to load reply record:', e)
    replyModal.value.show = false
  }
}

function resetHighlightedComment(recordID) {
  if (highlightedCommentId.value === recordID) {
    highlightedCommentId.value = null
  }
}

// ---- Scroll helpers ----

function scrollToPosition() {
  const container = chatContainer.value
  if (!container) return
  if (showNewestFirst.value) {
    container.scrollTop = container.scrollHeight
  } else {
    container.scrollTop = 0
  }
}

function scrollToLatest() {
  const container = chatContainer.value
  if (container) {
    container.scrollTop = container.scrollHeight
  }
}

function isScrollAtBottom() {
  const container = chatContainer.value
  if (!container) return false
  return container.scrollTop + container.clientHeight >= container.scrollHeight - 25
}

// ---- Prefilter ----

function expandFilter() {
   
  if (!props.record) {
    // If there is no current record and we are using recordID/ownerID variable in (pre)filter
    // we should disable the block
    if ((options.value.filter || '').includes('${record')) {
      throw Error(t('block.comment.invalidRecordVar'))
    }

    if ((options.value.filter || '').includes('${ownerID}')) {
      throw Error(t('block.comment.invalidOwnerVar'))
    }
  }

  if (options.value.filter) {
    try {
      return evaluatePrefilter(options.value.filter, {
        record: props.record,
        user: $Auth?.user || {},
        recordID: (props.record || {}).recordID || '0',
        ownerID: (props.record || {}).ownedBy || '0',
        userID: ($Auth?.user || {}).userID || '0',
      })
    } catch (e) {
      return e
    }
  }

  return ''
}

// ---- Watchers ----

watch(
  () => props.record?.recordID,
  () => refresh(),
  { immediate: true },
)

watch(
  () => options.value,
  () => refresh(),
  { deep: true },
)

// ---- Lifecycle ----

onMounted(() => {
  refreshInterval = setInterval(() => {
    if (autoFetching || submitting.value || loadingMore.value) return
    if (!showNewestFirst.value && filter.value.nextPage) return

    autoFetching = true
    loadNewComments()
      .catch(() => {})
      .finally(() => {
        autoFetching = false
      })
  }, 5000)
})

onBeforeUnmount(() => {
  // Cancel pending requests
  abortableRequests.value.forEach(cancel => {
    if (typeof cancel === 'function') cancel()
  })
  abortableRequests.value = []

  if (refreshInterval) {
    clearInterval(refreshInterval)
    refreshInterval = null
  }

  offRefetch?.()
})

const offRefetch = $eventBus?.on('refetch-records', () => refresh())
</script>

<style scoped>
.reply-to-container .reply-to-close {
  position: absolute;
  top: 0;
  right: 0;
  opacity: 0;
  transition: opacity 0.2s ease;
  z-index: 1;
}

.reply-to-container:hover .reply-to-close {
  opacity: 1;
}
</style>
