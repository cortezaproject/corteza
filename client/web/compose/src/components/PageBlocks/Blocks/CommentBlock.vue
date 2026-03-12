<template>
  <PageBlock :block="block">
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
                :namespace="namespace"
                :show-header="ci === 0"
                :show-title="!!titleField"
                :show-content="!!contentField"
                :highlighted="highlightedCommentId === comment.recordID"
                class="mb-1"
                @reply="replyToComment(comment)"
                @edit="onEditComment(comment, $event)"
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
        <section v-if="canAddRecord" class="flex flex-col bg-surface-0 border-t">
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
          />

          <!-- Submit -->
          <div class="flex items-center justify-end m-2 gap-1">
            <Button
              :label="$t('block.comment.submit')"
              :disabled="!isValid || submitting"
              :loading="submitting"
              size="small"
              @click="submitComment"
            />
          </div>
        </section>
      </div>
    </template>
  </PageBlock>
</template>

<script setup>
import { ref, computed, watch, inject, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { compose } from '@cortezaproject/corteza-js-next'
import { useModuleStore } from '@/stores/module'
import { useUserStore } from '@/stores/user'
import PageBlock from './PageBlock.vue'
import CommentItem from './Comment/CommentItem.vue'
import CommentReply from './Comment/CommentReply.vue'
import { components } from '@cortezaproject/corteza-vue-next'

const { CRichTextInput } = components

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const $ComposeAPI = inject('$ComposeAPI')
const $Auth = inject('$Auth')
const moduleStore = useModuleStore()
const userStore = useUserStore()

const chatContainer = ref(null)
const contentInput = ref(null)
const processing = ref(false)
const submitting = ref(false)
const loadingMore = ref(false)
const comments = ref([])
const highlightedCommentId = ref(null)

const filter = ref({
  limit: 50,
  nextPage: '',
})

const newComment = ref({
  title: '',
  content: '',
  replyTo: null,
})

let refreshInterval = null

const options = computed(() => props.block.options || {})
const showNewestFirst = computed(() => options.value.sortDirection === 'asc')
const hasNextPage = computed(() => !!filter.value.nextPage)

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

const canAddRecord = computed(() => roModule.value?.canCreateRecord)

const isValid = computed(() => !!newComment.value.title || !!newComment.value.content)

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

// ---- Fetch comments ----

async function fetchCommentRecords(query, useNextPage = true) {
  if (!roModule.value || !$ComposeAPI) return []

  const mod = roModule.value

  // Add reference filter
  if (referenceField.value && reference.value !== '0') {
    const refFilter = `${referenceField.value.name} = '${reference.value}'`
    query = query ? `${query} AND ${refFilter}` : refFilter
  }

  const sort = showNewestFirst.value ? 'createdAt DESC' : 'createdAt ASC'

  const params = {
    namespaceID: props.namespace.namespaceID,
    moduleID: mod.moduleID,
    query,
    sort: useNextPage && filter.value.nextPage ? '' : sort,
    limit: useNextPage ? filter.value.limit : 500,
    pageCursor: useNextPage ? filter.value.nextPage : '',
  }

  try {
    const { set = [], filter: paging = {} } = await $ComposeAPI.recordList(params)

    if (useNextPage) {
      filter.value.nextPage = paging.nextPage || ''
    }

    const records = set.map(r => new compose.Record(mod, r))

    // Resolve users for author info
    const userIDs = [...new Set(records.map(r => r.createdBy).filter(Boolean))]
    if (userIDs.length) {
      await userStore.resolveUsers(userIDs).catch(() => {})
    }

    // Group by date then by author
    if (showNewestFirst.value) {
      records.reverse()
    }

    const groups = {}
    records.forEach(comment => {
      const date = getFormattedDate(comment.createdAt)
      const authorId = comment.createdBy
      comment.author = getAuthor(authorId)

      // Resolve reply if replyField is set
      // (simplified: no deep reply resolution)

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
    console.error('Failed to fetch comments:', e)
    return []
  }
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
  return lastMsgGroup.comments[lastMsgGroup.comments.length - 1]?.createdAt
})

// ---- Load new comments (auto-refresh) ----

async function loadNewComments() {
  let query = expandFilter()

  if (lastCommentTimestamp.value) {
    const tsFilter = `createdAt > '${lastCommentTimestamp.value}'`
    query = query ? `${query} AND ${tsFilter}` : tsFilter
  }

  const wasAtBottom = isScrollAtBottom()
  const newGroups = await fetchCommentRecords(query, false)
  comments.value = mergeMessageGroups(comments.value, newGroups)

  if (wasAtBottom) {
    nextTick(() => scrollToLatest())
  }
}

// ---- Refresh ----

async function refresh() {
  if (!options.value.moduleID || !roModule.value || !contentField.value) return

  processing.value = true
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
    setTimeout(() => {
      processing.value = false
      nextTick(() => scrollToPosition())
    }, 300)
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

async function submitComment() {
  if (!isValid.value || !roModule.value) return

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

    await $ComposeAPI.recordCreate(record)

    newComment.value.title = ''
    newComment.value.content = ''
    newComment.value.replyTo = null

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

// ---- Reply ----

function replyToComment(comment) {
  newComment.value.replyTo = comment
  nextTick(() => {
    contentInput.value?.$el?.focus?.()
  })
}

function handleReplyClick(recordID) {
  const el = document.getElementById(`comment-${recordID}`)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'center' })
    highlightedCommentId.value = recordID
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
  return options.value.filter || ''
}

// ---- i18n ----
import { useI18n } from 'vue-i18n'
const { t } = useI18n()

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
    if (submitting.value || loadingMore.value) return
    if (!showNewestFirst.value && filter.value.nextPage) return
    loadNewComments().catch(() => {})
  }, 5000)
})

onBeforeUnmount(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
    refreshInterval = null
  }
})
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
