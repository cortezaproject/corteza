export function finder (translations, currentLanguage, resource) {
  return (key) => {
    return translations.find(t => t.key === key && t.lang === currentLanguage && t.resource === resource)
  }
}

export function getter (translations, currentLanguage, resource) {
  return (key) => {
    return (finder(translations, currentLanguage, resource)(key) || { message: undefined }).message
  }
}
