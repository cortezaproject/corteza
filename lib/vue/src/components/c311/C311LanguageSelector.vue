<template>
  <label class="d-inline-flex align-items-center mb-0">
    <span class="sr-only">{{ label }}</span>
    <select
      v-model="selected"
      class="custom-select custom-select-sm"
      :aria-label="label"
      @change="changeLanguage"
    >
      <option v-for="option in options" :key="option.value" :value="option.value">
        {{ option.label }}
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
    return {
      selected,
      options: [
        { value: 'en', label: 'English' },
        { value: 'es', label: 'Español' },
        { value: 'vi', label: 'Tiếng Việt' },
      ],
    }
  },
  watch: {
    value (value) { this.selected = value },
  },
  mounted () {
    if (this.selected !== this.value) this.$i18n?.i18next?.changeLanguage(this.selected)
  },
  methods: {
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
