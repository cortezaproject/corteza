import { getter } from './util'

export default function (field, translations, currentLanguage, resource) {
  const get = getter(translations, currentLanguage, resource)

  field.label = get('label')
  field.options.description.view = get('meta.description.view')
  field.options.description.edit = get('meta.description.edit')
  field.options.hint.view = get('meta.hint.view')
  field.options.hint.edit = get('meta.hint.edit')

  if (field.expressions && Array.isArray(field.expressions.validators)) {
    for (const vld of field.expressions.validators) {
      vld.error = get(`expression.validator.${vld.validatorID}.error`)
    }
  }
}
