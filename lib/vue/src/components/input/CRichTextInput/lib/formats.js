/* eslint-disable @typescript-eslint/explicit-function-return-type */

import { TaskItem, TaskList } from '@tiptap/extension-list'
import Mention from '@tiptap/extension-mention'
import { TableKit } from '@tiptap/extension-table'
import TextAlign from '@tiptap/extension-text-align'
import { TextStyleKit } from '@tiptap/extension-text-style'
import StarterKit from '@tiptap/starter-kit'
import suggestion from '../mentions/suggestion.js'

// Defines a set of formats that our document supports
export const getFormats = () => [
  StarterKit,
  TextStyleKit,
  TaskList,
  TaskItem.configure({
    nested: true,
  }),
  TextAlign.configure({
    types: ['heading', 'paragraph'],
  }),
  TableKit.configure({ table: { resizable: true } }),
  Mention.configure({
    HTMLAttributes: {
      class: 'mention',
    },
    suggestion,
  }),
]

// Defines the structure of our editor toolbar
export const getToolbar = () => [
  { type: 'bold', mark: true, icon: 'bold' },
  { type: 'italic', mark: true, icon: 'italic' },
  { type: 'underline', mark: true, icon: 'underline' },
  { type: 'strike', mark: true, icon: 'strikethrough' },
  { type: 'color', mark: true, component: 'Color' },
  { type: 'background', mark: true, component: 'Color', props: { background: true } },

  { type: 'blockquote', node: true, icon: 'quote-right' },
  { type: 'codeBlock', node: true, icon: 'code' },
  { type: 'heading', node: true, label: 'H1', attrs: { level: 1 } },
  { type: 'heading', node: true, label: 'H2', attrs: { level: 2 } },
  { type: 'heading', node: true, label: 'H3', attrs: { level: 3 } },
  { type: 'paragraph', node: true, icon: 'paragraph' },
  { type: 'orderedList', node: true, icon: 'list-ol' },
  { type: 'bulletList', node: true, icon: 'list-ul' },
  { type: 'taskList', node: true, icon: 'tasks' },

  {
    type: 'alignment',
    icon: 'align-left',
    nodeAttr: true,
    component: 'Alignment',
    variants: [
      { variant: 'left', icon: 'align-left', attrs: 'left' },
      { variant: 'center', icon: 'align-center', attrs: 'center' },
      { variant: 'right', icon: 'align-right', attrs: 'right' },
      { variant: 'justify', icon: 'align-justify', attrs: 'justify' },
    ],
  },

  {
    type: 'table',
    icon: 'table',
    nodeAttr: true,
    component: 'Table',
    variants: [
      { label: 'Insert Table', type: 'insertTable', attrs: { rows: 3, cols: 3, withHeaderRow: true } },
      { label: 'Insert Column Before', type: 'addColumnBefore' },
      { label: 'Insert Column After', type: 'addColumnAfter' },
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
    ],
  },

  { type: 'link', mark: true, component: 'Link', icon: 'link', attrs: { href: null } },

  // @note There is no free FA icon for this
  { type: 'horizontalRule', node: true, label: '__' },
]

export const nodeTypes = getToolbar().filter(({ node }) => node).map(({ type }) => type)
