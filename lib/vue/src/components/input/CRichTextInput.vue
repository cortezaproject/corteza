<template>
  <div class="c-rich-text-input rounded">
    <!-- Toolbar -->
    <div v-if="editor && !hideToolbar" class="rt-toolbar flex flex-wrap items-center p-1 gap-0.5">
      <!-- Bold -->
      <Button
        v-tooltip.top="'Bold'"
        text
        size="small"
        :severity="editor.isActive('bold') ? 'primary' : 'secondary'"
        class="rt-toolbar-btn"
        @click="editor.chain().focus().toggleBold().run()"
      >
        <template #icon><FontAwesomeIcon :icon="faBold" /></template>
      </Button>
      <!-- Italic -->
      <Button
        v-tooltip.top="'Italic'"
        text
        size="small"
        :severity="editor.isActive('italic') ? 'primary' : 'secondary'"
        class="rt-toolbar-btn"
        @click="editor.chain().focus().toggleItalic().run()"
      >
        <template #icon><FontAwesomeIcon :icon="faItalic" /></template>
      </Button>
      <!-- Underline -->
      <Button
        v-tooltip.top="'Underline'"
        text
        size="small"
        :severity="editor.isActive('underline') ? 'primary' : 'secondary'"
        class="rt-toolbar-btn"
        @click="editor.chain().focus().toggleUnderline().run()"
      >
        <template #icon><FontAwesomeIcon :icon="faUnderline" /></template>
      </Button>
      <!-- Strikethrough -->
      <Button
        v-tooltip.top="'Strikethrough'"
        text
        size="small"
        :severity="editor.isActive('strike') ? 'primary' : 'secondary'"
        class="rt-toolbar-btn"
        @click="editor.chain().focus().toggleStrike().run()"
      >
        <template #icon><FontAwesomeIcon :icon="faStrikethrough" /></template>
      </Button>

      <!-- Text color -->
      <CInputColorPicker
        v-model="textColor"
        :show-text="false"
        width="22px"
        height="22px"
        @update:model-value="applyTextColor"
      >
        <template #trigger="{ toggle }">
          <Button
            v-tooltip.top="'Text color'"
            text
            size="small"
            severity="secondary"
            class="rt-toolbar-btn"
            @click="toggle"
          >
            <span :style="{ borderBottom: `3px solid ${displayTextColor}` }" class="rt-color-label">A</span>
          </Button>
        </template>
      </CInputColorPicker>

      <!-- Background color -->
      <CInputColorPicker
        v-model="bgColor"
        :show-text="false"
        width="22px"
        height="22px"
        @update:model-value="applyBgColor"
      >
        <template #trigger="{ toggle }">
          <Button
            v-tooltip.top="'Highlight color'"
            text
            size="small"
            severity="secondary"
            class="rt-toolbar-btn"
            @click="toggle"
          >
            <span :style="{ backgroundColor: displayBgColor }" class="rt-color-label rt-bg-label">A</span>
          </Button>
        </template>
      </CInputColorPicker>

      <div class="rt-toolbar-separator" />

      <!-- Blockquote -->
      <Button
        v-tooltip.top="'Blockquote'"
        text
        size="small"
        :severity="editor.isActive('blockquote') ? 'primary' : 'secondary'"
        class="rt-toolbar-btn"
        @click="editor.chain().focus().toggleBlockquote().run()"
      >
        <template #icon><FontAwesomeIcon :icon="faQuoteRight" /></template>
      </Button>
      <!-- Code block -->
      <Button
        v-tooltip.top="'Code block'"
        text
        size="small"
        :severity="editor.isActive('codeBlock') ? 'primary' : 'secondary'"
        class="rt-toolbar-btn"
        @click="editor.chain().focus().toggleCodeBlock().run()"
      >
        <template #icon><FontAwesomeIcon :icon="faCode" /></template>
      </Button>

      <div class="rt-toolbar-separator" />

      <!-- Headings -->
      <Button
        text
        size="small"
        :severity="editor.isActive('heading', { level: 1 }) ? 'primary' : 'secondary'"
        class="rt-toolbar-btn font-bold"
        label="H1"
        @click="editor.chain().focus().toggleHeading({ level: 1 }).run()"
      />
      <Button
        text
        size="small"
        :severity="editor.isActive('heading', { level: 2 }) ? 'primary' : 'secondary'"
        class="rt-toolbar-btn font-bold"
        label="H2"
        @click="editor.chain().focus().toggleHeading({ level: 2 }).run()"
      />
      <Button
        text
        size="small"
        :severity="editor.isActive('heading', { level: 3 }) ? 'primary' : 'secondary'"
        class="rt-toolbar-btn font-bold"
        label="H3"
        @click="editor.chain().focus().toggleHeading({ level: 3 }).run()"
      />
      <!-- Paragraph -->
      <Button
        v-tooltip.top="'Paragraph'"
        text
        size="small"
        :severity="editor.isActive('paragraph') ? 'primary' : 'secondary'"
        class="rt-toolbar-btn"
        @click="editor.chain().focus().setParagraph().run()"
      >
        <template #icon><FontAwesomeIcon :icon="faParagraph" /></template>
      </Button>

      <div class="rt-toolbar-separator" />

      <!-- Ordered list -->
      <Button
        v-tooltip.top="'Ordered list'"
        text
        size="small"
        :severity="editor.isActive('orderedList') ? 'primary' : 'secondary'"
        class="rt-toolbar-btn"
        @click="editor.chain().focus().toggleOrderedList().run()"
      >
        <template #icon><FontAwesomeIcon :icon="faListOl" /></template>
      </Button>
      <!-- Bullet list -->
      <Button
        v-tooltip.top="'Bullet list'"
        text
        size="small"
        :severity="editor.isActive('bulletList') ? 'primary' : 'secondary'"
        class="rt-toolbar-btn"
        @click="editor.chain().focus().toggleBulletList().run()"
      >
        <template #icon><FontAwesomeIcon :icon="faListUl" /></template>
      </Button>
      <!-- Task list -->
      <Button
        v-tooltip.top="'Task list'"
        text
        size="small"
        :severity="editor.isActive('taskList') ? 'primary' : 'secondary'"
        class="rt-toolbar-btn"
        @click="editor.chain().focus().toggleTaskList().run()"
      >
        <template #icon><FontAwesomeIcon :icon="faListCheck" /></template>
      </Button>

      <div class="rt-toolbar-separator" />

      <!-- Alignment dropdown -->
      <div class="relative">
        <Button
          v-tooltip.top="'Alignment'"
          text
          size="small"
          severity="secondary"
          class="rt-toolbar-btn"
          @click="toggleAlignMenu"
        >
          <template #icon><FontAwesomeIcon :icon="activeAlignIcon" /></template>
        </Button>
        <div v-if="showAlignMenu" class="rt-dropdown">
          <button
            v-for="a in alignments"
            :key="a.value"
            class="rt-dropdown-item"
            :class="{ 'rt-active': editor.isActive({ textAlign: a.value }) }"
            @click="applyAlign(a.value)"
          >
            <FontAwesomeIcon :icon="a.icon" />
          </button>
        </div>
      </div>

      <div class="rt-toolbar-separator" />

      <!-- Table dropdown -->
      <div class="relative">
        <Button
          v-tooltip.top="'Table'"
          text
          size="small"
          severity="secondary"
          class="rt-toolbar-btn"
          @click="toggleTableMenu"
        >
          <template #icon><FontAwesomeIcon :icon="faTable" /></template>
        </Button>
        <div v-if="showTableMenu" class="rt-dropdown rt-dropdown-wide">
          <button
            v-for="op in tableOps"
            :key="op.type"
            class="rt-dropdown-item rt-dropdown-item-text"
            @click="execTableOp(op)"
          >
            {{ op.label }}
          </button>
        </div>
      </div>

      <div class="rt-toolbar-separator" />

      <!-- Link -->
      <div class="relative">
        <Button
          v-tooltip.top="'Link'"
          text
          size="small"
          :severity="editor.isActive('link') ? 'primary' : 'secondary'"
          class="rt-toolbar-btn"
          @click="toggleLinkInput"
        >
          <template #icon><FontAwesomeIcon :icon="faLink" /></template>
        </Button>
        <div v-if="showLinkInput" class="rt-dropdown rt-link-dropdown">
          <input
            ref="linkUrlInput"
            v-model="linkUrl"
            type="url"
            placeholder="https://..."
            class="rt-link-input"
            @keydown.enter.prevent="applyLink"
            @keydown.escape.prevent="showLinkInput = false"
          >
          <Button
            size="small"
            severity="success"
            label="OK"
            class="rt-link-btn"
            @click="applyLink"
          />
        </div>
      </div>

      <div class="rt-toolbar-separator" />

      <!-- Emoji picker -->
      <Button
        v-tooltip.top="'Emoji'"
        text
        size="small"
        severity="secondary"
        icon="pi pi-face-smile"
        class="rt-toolbar-btn"
        @click="toggleEmojiPicker"
      />
      <Popover ref="emojiPopoverRef" @show="onEmojiPopoverShow">
        <CEmojiPicker
          ref="emojiPicker"
          :emojis="allEmojis"
          :labels="emojiLabels"
          :show-quick-reactions="false"
          @select="onEmojiSelect"
        />
      </Popover>

      <div class="rt-toolbar-separator" />

      <!-- Horizontal rule -->
      <Button
        v-tooltip.top="'Horizontal rule'"
        text
        size="small"
        severity="secondary"
        label="__"
        class="rt-toolbar-btn font-bold"
        @click="editor.chain().focus().setHorizontalRule().run()"
      />

      <div class="rt-toolbar-separator" />

      <!-- Clear formatting -->
      <Button
        v-tooltip.top="'Clear formatting'"
        text
        size="small"
        severity="secondary"
        class="rt-toolbar-btn"
        @click="editor.chain().focus().unsetAllMarks().run()"
      >
        <template #icon><FontAwesomeIcon :icon="faTextSlash" /></template>
      </Button>
    </div>

    <!-- Editor content -->
    <editor-content
      v-if="editor"
      :editor="editor"
      :class="bodyClass"
      class="rt-editor-content rt-content p-2"
      :style="{ minHeight: minBodyHeight, maxHeight: maxBodyHeight }"
      @drop="onDrop"
      @paste="onPaste"
      @dragover.prevent
    />
  </div>
</template>

<script setup>
import { ref, watch, nextTick, computed } from 'vue'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Placeholder from '@tiptap/extension-placeholder'
import Underline from '@tiptap/extension-underline'
import TextStyle from '@tiptap/extension-text-style'
import { Color } from '@tiptap/extension-color'
import Highlight from '@tiptap/extension-highlight'
import TextAlign from '@tiptap/extension-text-align'
import Link from '@tiptap/extension-link'
import Table from '@tiptap/extension-table'
import TableRow from '@tiptap/extension-table-row'
import TableHeader from '@tiptap/extension-table-header'
import TableCell from '@tiptap/extension-table-cell'
import TaskList from '@tiptap/extension-task-list'
import TaskItem from '@tiptap/extension-task-item'
import Emoji, { emojis as emojiData } from '@tiptap/extension-emoji'
import emojiSuggestion from './CRichTextInput/emoji/suggestion.js'
import CEmojiPicker from '../CEmojiPicker.vue'
import CInputColorPicker from './CInputColorPicker.vue'

// Font Awesome — on-demand imports (tree-shaken)
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import {
  faBold,
  faItalic,
  faUnderline,
  faStrikethrough,
  faQuoteRight,
  faCode,
  faParagraph,
  faListOl,
  faListUl,
  faListCheck,
  faAlignLeft,
  faAlignCenter,
  faAlignRight,
  faAlignJustify,
  faTable,
  faLink,
  faTextSlash,
} from '@fortawesome/free-solid-svg-icons'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  placeholder: {
    type: String,
    default: '',
  },
  hideToolbar: {
    type: Boolean,
    default: false,
  },
  minBodyHeight: {
    type: String,
    default: '4rem',
  },
  maxBodyHeight: {
    type: String,
    default: '',
  },
  bodyClass: {
    type: String,
    default: '',
  },
  emojiLabels: {
    type: Object,
    default: () => ({}),
  },
})

const emit = defineEmits(['update:modelValue', 'upload'])

// Refs
const linkUrlInput = ref(null)
const emojiPicker = ref(null)
const emojiPopoverRef = ref(null)

// State
const textColor = ref('#000000FF')
const bgColor = ref('transparent')
const showAlignMenu = ref(false)
const showTableMenu = ref(false)
const showLinkInput = ref(false)
const linkUrl = ref('')

// Emoji data for the picker
const allEmojis = computed(() => emojiData || [])

// Display colors for toolbar A labels
const displayTextColor = computed(() => {
  if (!textColor.value || textColor.value === 'transparent') return '#000000'
  return textColor.value.substring(0, 7)
})

const displayBgColor = computed(() => {
  if (!bgColor.value || bgColor.value === 'transparent') return 'transparent'
  return bgColor.value.substring(0, 7)
})

// Alignment config
const alignments = [
  { value: 'left', icon: faAlignLeft },
  { value: 'center', icon: faAlignCenter },
  { value: 'right', icon: faAlignRight },
  { value: 'justify', icon: faAlignJustify },
]

const activeAlignIcon = computed(() => {
  if (!editor.value) return faAlignLeft
  for (const a of alignments) {
    if (editor.value.isActive({ textAlign: a.value })) return a.icon
  }
  return faAlignLeft
})

// Table operations
const tableOps = [
  { label: 'Insert Table', type: 'insertTable', attrs: { rows: 3, cols: 3, withHeaderRow: true } },
  { label: 'Add Column Before', type: 'addColumnBefore' },
  { label: 'Add Column After', type: 'addColumnAfter' },
  { label: 'Delete Column', type: 'deleteColumn' },
  { label: 'Add Row Before', type: 'addRowBefore' },
  { label: 'Add Row After', type: 'addRowAfter' },
  { label: 'Delete Row', type: 'deleteRow' },
  { label: 'Merge Cells', type: 'mergeCells' },
  { label: 'Split Cell', type: 'splitCell' },
  { label: 'Toggle Header Row', type: 'toggleHeaderRow' },
  { label: 'Toggle Header Cell', type: 'toggleHeaderCell' },
  { label: 'Toggle Header Column', type: 'toggleHeaderColumn' },
  { label: 'Delete Table', type: 'deleteTable' },
]

// Editor
const editor = useEditor({
  extensions: [
    StarterKit,
    Underline,
    TextStyle,
    Color,
    Highlight.configure({ multicolor: true }),
    TextAlign.configure({ types: ['heading', 'paragraph'] }),
    Link.configure({ openOnClick: false }),
    Table.configure({ resizable: true }),
    TableRow,
    TableHeader,
    TableCell,
    TaskList,
    TaskItem.configure({ nested: true }),
    Placeholder.configure({ placeholder: props.placeholder }),
    Emoji.configure({
      enableEmoticons: true,
      suggestion: emojiSuggestion,
    }),
  ],
  content: props.modelValue || '',
  parseOptions: {
    preserveWhitespace: 'full',
  },
  onUpdate: ({ editor: e }) => {
    const html = e.getHTML().replace(/<p><\/p>/g, '<p><br></p>')
    const value = html === '<p><br></p>' ? '' : html
    emittedContent = true
    emit('update:modelValue', value)
  },
})

let emittedContent = false

watch(() => props.modelValue, (val) => {
  if (!emittedContent && editor.value) {
    editor.value.commands.setContent(val || '', false)
  }
  emittedContent = false
})

// Color methods
function applyTextColor(color) {
  textColor.value = color
  const hex6 = color === 'transparent' ? '#000000' : color.substring(0, 7)
  editor.value?.chain().focus().setColor(hex6).run()
}

function applyBgColor(color) {
  bgColor.value = color
  if (color === 'transparent') {
    editor.value?.chain().focus().unsetHighlight().run()
  } else {
    const hex6 = color.substring(0, 7)
    editor.value?.chain().focus().toggleHighlight({ color: hex6 }).run()
  }
}

// Alignment methods
function toggleAlignMenu() {
  showAlignMenu.value = !showAlignMenu.value
  showTableMenu.value = false
  showLinkInput.value = false
}

function applyAlign(alignment) {
  editor.value?.chain().focus().setTextAlign(alignment).run()
  showAlignMenu.value = false
}

// Table methods
function toggleTableMenu() {
  showTableMenu.value = !showTableMenu.value
  showAlignMenu.value = false
  showLinkInput.value = false
}

function execTableOp(op) {
  const chain = editor.value?.chain().focus()
  if (!chain) return

  switch (op.type) {
    case 'insertTable': chain.insertTable(op.attrs).run(); break
    case 'addColumnBefore': chain.addColumnBefore().run(); break
    case 'addColumnAfter': chain.addColumnAfter().run(); break
    case 'deleteColumn': chain.deleteColumn().run(); break
    case 'addRowBefore': chain.addRowBefore().run(); break
    case 'addRowAfter': chain.addRowAfter().run(); break
    case 'deleteRow': chain.deleteRow().run(); break
    case 'mergeCells': chain.mergeCells().run(); break
    case 'splitCell': chain.splitCell().run(); break
    case 'toggleHeaderRow': chain.toggleHeaderRow().run(); break
    case 'toggleHeaderCell': chain.toggleHeaderCell().run(); break
    case 'toggleHeaderColumn': chain.toggleHeaderColumn().run(); break
    case 'deleteTable': chain.deleteTable().run(); break
  }

  showTableMenu.value = false
}

// Link methods
function toggleLinkInput() {
  showLinkInput.value = !showLinkInput.value
  showAlignMenu.value = false
  showTableMenu.value = false

  if (showLinkInput.value) {
    // Pre-fill with existing link href
    const attrs = editor.value?.getAttributes('link')
    linkUrl.value = attrs?.href || ''
    nextTick(() => linkUrlInput.value?.focus())
  }
}

function applyLink() {
  if (linkUrl.value) {
    editor.value?.chain().focus().setLink({ href: linkUrl.value }).run()
  } else {
    editor.value?.chain().focus().unsetLink().run()
  }
  showLinkInput.value = false
  linkUrl.value = ''
}

// Emoji picker methods
function toggleEmojiPicker(event) {
  emojiPopoverRef.value?.toggle(event)
  showAlignMenu.value = false
  showTableMenu.value = false
  showLinkInput.value = false
}

function onEmojiPopoverShow() {
  nextTick(() => {
    if (emojiPicker.value) {
      emojiPicker.value.reset()
    }
  })
}

function onEmojiSelect(emoji) {
  if (emoji && emoji.name) {
    editor.value?.chain().focus().insertContent({
      type: 'emoji',
      attrs: { name: emoji.name },
    }).run()
  }
  emojiPopoverRef.value?.hide()
}

// Utility
function focus() {
  if (editor.value) {
    editor.value.commands.focus()
  }
}

function onDrop(event) {
  if (event.dataTransfer?.files?.length > 0) {
    event.preventDefault()
    emit('upload', event.dataTransfer.files)
  }
}

function onPaste(event) {
  if (event.clipboardData?.files?.length > 0) {
    event.preventDefault()
    emit('upload', event.clipboardData.files)
  }
}

defineExpose({ focus, editor, allEmojis })
</script>

<style>
.c-rich-text-input {
  display: flex;
  flex-direction: column;
}

.c-rich-text-input .rt-toolbar {
  border-bottom: 1px solid var(--p-surface-200, #e2e8f0);
}

.c-rich-text-input .rt-editor-content {
  height: 100%;
}

.c-rich-text-input .rt-editor-content .tiptap {
  height: 100%;
  outline: none;
}

.c-rich-text-input .rt-editor-content .tiptap p.is-editor-empty:first-child::before {
  color: var(--p-text-muted-color, #adb5bd);
  content: attr(data-placeholder);
  float: left;
  height: 0;
  pointer-events: none;
}

.c-rich-text-input .rt-editor-content .tiptap p {
  margin: 0;
}

/* Toolbar button styling — matches Corteza's 2.25rem */
.rt-toolbar-btn {
  width: 2.25rem !important;
  height: 2.25rem !important;
  padding: 0.25rem !important;
  border-radius: 0.25rem !important;
}

.rt-toolbar-btn:hover {
  background-color: var(--p-surface-100, #f1f5f9) !important;
}

.rt-toolbar-separator {
  width: 1px;
  height: 1.5rem;
  background: var(--p-surface-200, #e2e8f0);
  margin: 0 0.25rem;
}

/* Color buttons */
.rt-color-label {
  font-weight: 700;
  font-size: 0.875rem;
  display: inline-block;
  padding: 0 0.15rem;
  line-height: 1.2;
}

.rt-bg-label {
  border-radius: 0.15rem;
  padding: 0 0.2rem;
}

/* Dropdown menus */
.rt-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  z-index: 100;
  background: var(--p-surface-0, #fff);
  border: 1px solid var(--p-surface-200, #e2e8f0);
  border-radius: 0.375rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  padding: 0.25rem;
  display: flex;
  gap: 0.15rem;
}

.rt-dropdown-wide {
  flex-direction: column;
  min-width: 12rem;
  padding: 0.25rem 0;
}

.rt-dropdown-item {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  border: none;
  background: none;
  border-radius: 0.25rem;
  cursor: pointer;
  color: var(--p-text-color, #1e293b);
}

.rt-dropdown-item:hover {
  background: var(--p-surface-100, #f1f5f9);
}

.rt-dropdown-item.rt-active {
  color: var(--p-primary-color, #3b82f6);
}

.rt-dropdown-item-text {
  width: auto;
  justify-content: flex-start;
  padding: 0.4rem 0.75rem;
  font-size: 0.85rem;
}

/* Link input dropdown */
.rt-link-dropdown {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.35rem;
  min-width: 16rem;
}

.rt-link-input {
  flex: 1;
  border: 1px solid var(--p-surface-300, #cbd5e1);
  border-radius: 0.25rem;
  padding: 0.3rem 0.5rem;
  font-size: 0.85rem;
  outline: none;
}

.rt-link-input:focus {
  border-color: var(--p-primary-color, #3b82f6);
}

.rt-link-btn {
  flex-shrink: 0;
}

/*
 * Rich text content styles are provided globally by rt-content.css
 * (imported in lib/vue/src/index.ts). The editor content area has
 * class="rt-content" so those styles apply automatically.
 *
 * Only editor-specific overrides below:
 */

/* In the editor, checkboxes should be interactive */
.c-rich-text-input .rt-editor-content .tiptap input[type="checkbox"] {
  pointer-events: auto !important;
  cursor: pointer !important;
}
</style>
