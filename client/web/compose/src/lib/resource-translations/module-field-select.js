import { getter } from './util'

// @note copied from FieldSelectTranslator.vue
const keyPrefix = 'meta.options.'
const keySuffix = '.text'

export default function (field, translations, currentLanguage, resource) {
  const get = getter(translations, currentLanguage, resource)

  field.options.options = field.options.options.map(opt => {
    if (typeof opt === 'string') {
      return { value: opt, text: get(`${keyPrefix}${opt}${keySuffix}`) || opt }
    }

    if (typeof opt === 'object' && opt.value) {
      opt.text = get(`${keyPrefix}${opt.value}${keySuffix}`) || opt.value
    }

    return opt
  })
}
