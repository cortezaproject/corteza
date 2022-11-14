/* eslint-disable @typescript-eslint/explicit-function-return-type */

// maps legacy Quill alignment classes to new classes
export const qAlignments = {
  'ql-align-right': 'right',
  'ql-align-left': 'left',
  'ql-align-center': 'center',
  'ql-align-justify': 'justify',
}

// helper to construct node's attributes
export function toAttrs (node) {
  if (node.attrs.alignment) {
    return { style: `text-align: ${node.attrs.alignment}` }
  }
  return {}
}

// helper to determine node's alignment
export function alignmentParser (node) {
  // Covers current structure
  let alignment = node.style.textAlign
  if (alignment) {
    return { alignment }
  }

  // Covers legacy structure
  node.classList.forEach((c) => {
    alignment = alignment || qAlignments[c]
  })

  return { alignment: alignment || undefined }
}

// helper to construct DOM parsers with alignment attrs
export const makeDOMParser = pd => ({
  ...pd,
  getAttrs: (node) => {
    let rtr = {}
    if (pd.getAttrs) {
      rtr = pd.getAttrs(node)
    }
    return { ...rtr, ...alignmentParser(node) }
  },
})
