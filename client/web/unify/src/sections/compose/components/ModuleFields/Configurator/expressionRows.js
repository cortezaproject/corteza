// A sanitizer or validator row the author left blank. The server parses every
// stored expression, and an empty one is a parse error ("unexpected EOF"), so
// each record save on the module fails; a validator without a message rejects
// the record with nothing to show for it.

export function isBlankExpressionRow(value) {
  return !String(value ?? '').trim()
}

export function incompleteExpressionRows(expressions) {
  const { sanitizers = [], validators = [] } = expressions || {}

  return (
    sanitizers.some(isBlankExpressionRow) ||
    validators.some(v => isBlankExpressionRow(v?.test) || isBlankExpressionRow(v?.error))
  )
}
