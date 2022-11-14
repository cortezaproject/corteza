/**
 * Decodes special flag to help the backend understand language of the payload sent
 *
 * @param resourceTranslationLanguage
 * @returns {{headers: {'Content-Language'}}|{}}
 */
export function config ({ resourceTranslationLanguage }) {
  if (resourceTranslationLanguage === undefined) {
    return {}
  }

  return { headers: { 'Content-Language': resourceTranslationLanguage } }
}
