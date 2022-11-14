export default {
  computed: {
    /**
     * Keep this inside methods otherwise it (no-value) gets caught by Vue
     * since settings are actually fetched a tad later...
     *
     * @returns {string[]}
     */
    resourceTranslationLanguages () {
      const ll = this.$Settings.get('resourceTranslations.languages')
      if (!ll || !Array.isArray(ll) || ll.length === 0) {
        return ['en']
      }

      return ll
    },

    /**
     * Resource translations are enabled when there is more than 1 language
     * @returns {boolean}
     */
    resourceTranslationsEnabled () {
      return this.resourceTranslationLanguages.length > 1
    },

    /**
     * 1st language is always default
     * @returns {*}
     */
    defaultTranslationLanguage () {
      return this.resourceTranslationLanguages[0]
    },

    currentLanguage () {
      return this.$i18n.i18next.language
    },

  },

  methods: {
    /**
     * Returns directionality of the page according to the language
     *  - Arabic (ar)
     *  - Hebrew (he)
     *  - Pashto (pa)
     *  - Persian (fa)
     *  - Urdu (ur)
     *  - Sindhi (sd)
     * @returns {string} rtl | ltr
     */
    textDirectionality (language = this.currentLanguage) {
      switch (language) {
        case 'ar':
        case 'he':
        case 'pa':
        case 'fa':
        case 'ur':
        case 'sd':
          return 'rtl'
        default:
          return 'ltr'
      }
    },
  },
}
