/* eslint-disable @typescript-eslint/explicit-function-return-type */

import {
  Paragraph,
  Heading,
  TextColor,
  TextBackground,
  Link,
  HtmlPaster,
} from '../extensions'

import {
  Blockquote,
  CodeBlock,
  HorizontalRule,
  OrderedList,
  BulletList,
  TodoList,
  Bold,
  Italic,
  Strike,
  Underline,
  ListItem,
  TodoItem,
  History,
  Table,
  TableRow,
  TableHeader,
  TableCell,
} from 'tiptap-extensions'

// Defines a set of formats that our document supports
export const getFormats = () => [
  new HtmlPaster(),
  new Bold(),
  new Italic(),
  new Underline(),
  new Strike(),
  new Blockquote(),
  new CodeBlock(),
  new OrderedList(),
  new BulletList(),
  new TodoList(),
  new Heading({ levels: [1, 2, 3] }),
  new Paragraph({ alignments: ['left', 'right', 'center', 'justify'] }),
  new Link(),
  new HorizontalRule(),
  new ListItem(),
  new TodoItem(),
  new History(),
  new TextBackground(),
  new TextColor(),
  new Table({
    resizable: true,
  }),
  new TableHeader(),
  new TableCell(),
  new TableRow(),
]

// Defines the structure of our editor toolbar
export const getToolbar = () => [
  { type: 'bold', mark: true, icon: 'bold' },
  { type: 'italic', mark: true, icon: 'italic' },
  { type: 'underline', mark: true, icon: 'underline' },
  { type: 'strike', mark: true, icon: 'strikethrough' },
  { type: 'color', mark: true, component: 'Color' },
  { type: 'background', mark: true, component: 'Color' },

  { type: 'blockquote', node: true, icon: 'quote-right' },
  { type: 'code_block', node: true, icon: 'code' },
  { type: 'heading', node: true, label: 'H1', attrs: { level: 1 } },
  { type: 'heading', node: true, label: 'H2', attrs: { level: 2 } },
  { type: 'heading', node: true, label: 'H3', attrs: { level: 3 } },
  { type: 'paragraph', node: true, icon: 'paragraph' },
  { type: 'ordered_list', node: true, icon: 'list-ol' },
  { type: 'bullet_list', node: true, icon: 'list-ul' },
  { type: 'todo_list', node: true, icon: 'tasks' },

  {
    type: 'alignment',
    icon: 'align-left',
    nodeAttr: true,
    component: 'Alignment',
    variants: [
      { variant: 'left', icon: 'align-left', attrs: { alignment: 'left' } },
      { variant: 'center', icon: 'align-center', attrs: { alignment: 'center' } },
      { variant: 'right', icon: 'align-right', attrs: { alignment: 'right' } },
      { variant: 'justify', icon: 'align-justify', attrs: { alignment: 'justify' } },
    ],
  },

  {
    type: 'table',
    icon: 'table',
    nodeAttr: true,
    component: 'Table',
    variants: [
      { variant: 'insert', icon: '', label: 'Insert Table', type: 'createTable', attrs: {rowsCount: 3, colsCount: 3, withHeaderRow: true } },
      { variant: 'splitCell', icon: '', label: 'Split Cell', type: 'splitCell', },
      { variant: 'mergeCells', icon: '', label: 'Merge Cells', type: 'mergeCells', },
      { variant: 'deleteColumn', icon: '', label: 'Delete Column', type: 'deleteColumn', },
      { variant: 'insertRowBefore', icon: '', label: 'Add Row Before', type: 'addRowBefore', },
      { variant: 'insertColumnAfter', icon: '', label: 'Insert Column After', type: 'addColumnAfter', },
      { variant: 'insertHeaderRow', icon: '', label: 'Insert Header Row', type: 'toggleHeaderRow', },
      { variant: 'insertHeaderCell', icon: '', label: 'Insert Header Cell', type: 'toggleHeaderCell', },
      { variant: 'insertColumnBefore', icon: '', label: 'Insert Column Before', type: 'addColumnBefore', },
      { variant: 'insertHeaderColumn', icon: '', label: 'Insert Header Column', type: 'toggleHeaderColumn', },
      { variant: 'deleteTable', icon: '', label: 'Delete Table', type: 'deleteTable', },
    ]
  },

  { type: 'link', mark: true, component: 'Link', icon: 'link', attrs: { href: null } },

  // @note There is no free FA icon for this
  { type: 'horizontal_rule', node: true, label: '__' },
]

export const nodeTypes = getToolbar().filter(({ node }) => node).map(({ type }) => type)
