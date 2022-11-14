import { getter } from './util'
import moduleFieldResTr from './module-field'
import moduleFieldSelectResTr from './module-field-select'

export default function (mod, translations, currentLanguage, resource) {
  const get = getter(translations, currentLanguage, resource)

  mod.name = get('name')

  for (const f of mod.fields) {
    const fRes = `compose:module-field/${mod.namespaceID}/${mod.moduleID}/${f.fieldID}`
    moduleFieldResTr(f, translations, currentLanguage, fRes)

    if (f.kind === 'Select') {
      moduleFieldSelectResTr(f, translations, currentLanguage, fRes)
    }
  }
}
