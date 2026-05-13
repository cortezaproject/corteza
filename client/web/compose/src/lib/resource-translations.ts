// client/web/compose/src/lib/resource-translations.ts

export interface ResourceTranslation {
  resource: string
  key: string
  lang: string
  message: string
}

// Returns the translated message for a given key/lang/resource, or undefined.
function get(
  translations: ResourceTranslation[],
  resource: string,
  key: string,
  lang: string,
): string | undefined {
  return translations.find(t => t.resource === resource && t.key === key && t.lang === lang)
    ?.message
}

// Applies translations to a module field object in-place.
export function applyFieldTranslations(
  field: any,
  translations: ResourceTranslation[],
  lang: string,
  resource: string,
): void {
  const g = (key: string) => get(translations, resource, key, lang)

  const label = g('label')
  if (label !== undefined) field.label = label

  if (!field.options) field.options = {}
  const descView = g('meta.description.view')
  if (descView !== undefined) field.options.description = { ...field.options.description, view: descView }
  const descEdit = g('meta.description.edit')
  if (descEdit !== undefined) field.options.description = { ...field.options.description, edit: descEdit }
  const hintView = g('meta.hint.view')
  if (hintView !== undefined) field.options.hint = { ...field.options.hint, view: hintView }
  const hintEdit = g('meta.hint.edit')
  if (hintEdit !== undefined) field.options.hint = { ...field.options.hint, edit: hintEdit }

  if (field.expressions?.validators) {
    for (const vld of field.expressions.validators) {
      const err = g(`expression.validator.${vld.validatorID}.error`)
      if (err !== undefined) vld.error = err
    }
  }
}

// Applies select option translations to a field in-place.
export function applySelectTranslations(
  field: any,
  translations: ResourceTranslation[],
  lang: string,
  resource: string,
): void {
  if (!Array.isArray(field.options?.options)) return
  field.options.options = field.options.options.map((opt: any) => {
    const value = typeof opt === 'string' ? opt : opt.value
    const key = `meta.options.${value}.text`
    const text = get(translations, resource, key, lang)
    if (typeof opt === 'string') return text ? { value: opt, text } : opt
    return { ...opt, text: text ?? opt.text ?? opt.value }
  })
}

// Applies bool label translations to a field in-place.
export function applyBoolTranslations(
  field: any,
  translations: ResourceTranslation[],
  lang: string,
  resource: string,
): void {
  const trueLabel = get(translations, resource, 'meta.bool.true.label', lang)
  if (trueLabel !== undefined && field.options) field.options.trueLabel = trueLabel
  const falseLabel = get(translations, resource, 'meta.bool.false.label', lang)
  if (falseLabel !== undefined && field.options) field.options.falseLabel = falseLabel
}

// Applies translations to an entire module (name + all fields) in-place.
export function applyModuleTranslations(
  module: any,
  translations: ResourceTranslation[],
  lang: string,
): void {
  const moduleResource = `compose:module/${module.namespaceID}/${module.moduleID}`
  const name = get(translations, moduleResource, 'name', lang)
  if (name !== undefined) module.name = name

  for (const field of module.fields || []) {
    const fRes = `compose:module-field/${module.namespaceID}/${module.moduleID}/${field.fieldID}`
    applyFieldTranslations(field, translations, lang, fRes)
    if (field.kind === 'Select') applySelectTranslations(field, translations, lang, fRes)
    if (field.kind === 'Bool') applyBoolTranslations(field, translations, lang, fRes)
  }
}

// Applies translations to a page layout object in-place.
export function applyPageLayoutTranslations(
  layout: any,
  translations: ResourceTranslation[],
  lang: string,
): void {
  const res = `compose:page-layout/${layout.namespaceID}/${layout.pageID}/${layout.pageLayoutID}`
  const g = (key: string) => get(translations, res, key, lang)

  const title = g('title')
  if (title !== undefined && layout.meta) layout.meta.title = title

  const desc = g('description')
  if (desc !== undefined && layout.meta) layout.meta.description = desc

  for (const btn of ['new', 'edit', 'submit', 'delete', 'clone', 'back']) {
    const label = g(`config.buttons.${btn}.label`)
    if (label !== undefined && layout.config?.buttons?.[btn]) {
      layout.config.buttons[btn].label = label
    }
  }
}

// Applies page-level and block-level translations in-place.
// Pass a non-empty `layouts` array to also apply layout translations.
export function applyPageTranslations(
  page: any,
  layouts: any[],
  translations: ResourceTranslation[],
  lang: string,
): void {
  const pageRes = `compose:page/${page.namespaceID}/${page.pageID}`
  const g = (key: string) => get(translations, pageRes, key, lang)

  const title = g('title')
  if (title !== undefined) page.title = title

  const desc = g('description')
  if (desc !== undefined) page.description = desc

  // Record-page toolbar buttons
  const zeroID = '0'
  if (page.moduleID && page.moduleID !== zeroID) {
    for (const btn of ['new', 'edit', 'submit', 'delete', 'clone', 'back']) {
      const label = g(`recordToolbar.${btn}.label`)
      if (label !== undefined && page.config?.buttons?.[btn]) {
        page.config.buttons[btn].label = label
      }
    }
  }

  // Block translations
  for (const block of page.blocks || []) {
    if (!block.blockID || block.blockID === zeroID) continue

    const bt = g(`pageBlock.${block.blockID}.title`)
    if (bt !== undefined) block.title = bt

    const bd = g(`pageBlock.${block.blockID}.description`)
    if (bd !== undefined) block.description = bd

    // Automation block buttons
    if (block.kind === 'Automation') {
      ;(block.options?.buttons || []).forEach((btn: any, i: number) => {
        const label = g(`pageBlock.${block.blockID}.button.${btn.buttonID || i}.label`)
        if (label !== undefined) btn.label = label
      })
    }

    // RecordList selection buttons
    if (block.kind === 'RecordList') {
      ;(block.options?.selectionButtons || []).forEach((btn: any, i: number) => {
        const label = g(`pageBlock.${block.blockID}.button.${btn.buttonID || i}.label`)
        if (label !== undefined) btn.label = label
      })
    }

    // Content block body
    if (block.kind === 'Content') {
      const body = g(`pageBlock.${block.blockID}.content.body`)
      if (body !== undefined && block.options) block.options.body = body
    }
  }

  // Layout translations
  for (const layout of layouts) {
    applyPageLayoutTranslations(layout, translations, lang)
  }
}
