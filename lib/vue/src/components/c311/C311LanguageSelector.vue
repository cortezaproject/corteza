<template>
  <label class="d-inline-flex align-items-center mb-0">
    <span class="sr-only">{{ translatedLabel }}</span>
    <select
      v-model="selected"
      class="custom-select custom-select-sm"
      data-c311-language
      :aria-label="translatedLabel"
      @change="changeLanguage"
    >
      <option v-for="option in options" :key="option.value" :value="option.value">
        {{ translate(option.key, option.label) }}
      </option>
    </select>
  </label>
</template>

<script>
import { readC311Locale } from '../../libs/c311-i18n'

export default {
  name: 'C311LanguageSelector',
  props: {
    value: { type: String, default: 'en' },
    actorID: { type: String, default: '' },
    label: { type: String, default: 'Language' },
  },
  data () {
    const selected = readC311Locale(this.actorID) || this.value
    return { selected }
  },
  computed: {
    translatedLabel () {
      return this.translate('language.label', this.label)
    },
    options () {
      return [
        { value: 'en', key: 'language.english', label: 'English' },
        { value: 'es', key: 'language.spanish', label: 'Español' },
        { value: 'vi', key: 'language.vietnamese', label: 'Tiếng Việt' },
      ]
    },
  },
  watch: {
    value (value) { this.selected = value },
  },
  mounted () {
    if (this.selected !== this.value) this.$i18n?.i18next?.changeLanguage(this.selected)
  },
  methods: {
    translate (key, fallback) {
      const translated = this.$t?.(`c311:${key}`)
      return translated && translated !== `c311:${key}` && translated !== key ? translated : fallback
    },
    changeLanguage () {
      const locale = this.selected
      if (typeof localStorage !== 'undefined') localStorage.setItem(`c311.locale.${this.actorID || 'anonymous'}`, locale)
      this.$i18n?.i18next?.changeLanguage(locale)
      this.$emit('input', locale)
      this.$emit('change', locale)
    },
  },
}
</script>
