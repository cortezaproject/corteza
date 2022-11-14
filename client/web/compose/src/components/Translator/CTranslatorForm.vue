<template>
  <b-form @submit.prevent="onSubmit">
    <div
      class="text-right mb-2"
    >
      <b-dropdown
        data-test-id="dropdown-add-language"
        :text="$t('add-language')"
        variant="light"
      >
        <b-dropdown-item-button
          v-for="lang in intLanguages"
          :key="lang.tag"
          :data-test-id="`dropdown-language-item-${lang.localizedName}`"
          :disabled="lang.default || lang.visible"
          @click="lang.visible = true"
        >
          {{ lang.localizedName }}
        </b-dropdown-item-button>
      </b-dropdown>
    </div>
    <b-table-simple>
      <b-thead>
        <b-tr>
          <b-th class="key p-1" />
          <b-th
            v-for="lang in visibleLanguages"
            :key="lang.tag"
            :data-test-id="`language-${lang.localizedName}`"
            class="text-truncate"
            :style="{ 'width': `${100 / visibleLanguages.length}%` }"
          >
            {{ lang.localizedName }}

            <b-button
              v-if="!lang.default"
              variant="link"
              class="float-right p-0 m-0"
              @click="lang.visible=false"
            >
              <font-awesome-icon
                :icon="['fas', 'times']"
              />
            </b-button>
          </b-th>
        </b-tr>
      </b-thead>
      <b-tbody
        v-for="(r, i) in resources()"
        :key="i"
      >
        <b-tr
          v-if="!r.isPrimary"
        >
          <b-th
            :data-test-id="`translation-title-${r.title}`"
            class="font-weight-bold border-top-0"
            :colspan="visibleLanguages.length + 1"
          >
            {{ r.title }}
          </b-th>
        </b-tr>
        <b-tr
          v-for="(key, k) in keys(r.resource)"
          :key="k"
          :class="{'bg-light': key === highlightKey }"
        >
          <b-td
            cols="2"
            class="text-break small"
          >
            <samp
              :data-test-id="`translation-field-${keyPrettifier(key)}`"
            >
              {{ keyPrettifier(key) }}
            </samp>
          </b-td>
          <b-td
            v-for="(lang, langIndex) in visibleLanguages"
            :key="lang.tag"
            :class="{'m-0 p-0': true, 'bg-warning': isDirty(r.resource, key, lang.tag) }"
          >
            <b-button
              v-if="isDirty(r.resource, key, lang.tag)"
              variant="link"
              class="float-right p-1 mt-2 mr-2"
              @click="reset(r.resource, key, lang.tag)"
            >
              <font-awesome-icon
                :icon="['fas', 'times']"
              />
            </b-button>
            <editable
              :data-test-id="`translation-value-${keyPrettifier(key)}-language-${lang.name}`"
              :value="msg(r.resource, key, lang.tag)"
              :placeholder="$t('missing-translation')"
              :tabindex="langIndex + 1"
              @input="onUpdate(r.resource, key, lang.tag, $event)"
            />
          </b-td>
        </b-tr>
      </b-tbody>
    </b-table-simple>
  </b-form>
</template>
<script lang="js">
import Editable from './Editable'

const lsKey = 'resource-translator.languages'

export default {
  i18nOptions: {
    namespaces: 'resource-translator',
    keyPrefix: 'translator',
  },

  components: {
    Editable,
  },

  props: {
    languages: {
      type: Array,
      required: true,
    },

    primaryResource: {
      type: String,
      required: true,
    },

    /**
     * Array of objects containing all relevant translations
     *
     * It is backend's responsibility to serve all relevant translations
     * for a certain resource
     */
    translations: {
      type: Array,
      required: true,
    },

    /**
     * Translation of resource types
     *
     * What we usually get for each translation is resource type like
     * compose:module/42 and this key/value allows us to map it to 'Module'
     * This is especially useful in case of fields when it's important
     * translator can differentiate between them
     */
    titles: {
      type: Object,
      default: () => ({}),
    },

    /**
     * When set, the row in translator is highlighted
     */
    highlightKey: {
      type: String,
      default: '',
    },

    disabled: {
      type: Boolean,
      default: () => false,
    },

    keyPrettifier: {
      type: Function,
      default: (key) => {
        return key
          .replace(/([A-Z])/, ' $1')
          .toLowerCase()
          .replace(/(\d+)/, '#$1')
          .split('.').map(s => s.substring(0, 1).toUpperCase() + s.substring(1))
          .join(' ')
      },
    },
  },

  data () {
    const preselected = (window.localStorage.getItem(lsKey) || '').split(',')

    return {
      intLanguages: this.languages.map((lang, i) => ({
        ...lang,
        // 1st one is default
        default: i === 0,
        // default is always visible
        // the rest, pick from the list from the local-store
        visible: i === 0 || preselected.includes(lang.tag),
      })),
      intTranslations: this.translations.map(t => ({ ...t, org: t.message, dirty: false })),
    }
  },

  computed: {
    visibleLanguages () {
      return this.intLanguages.filter(({ visible }) => visible)
    },
  },

  watch: {
    intLanguages: {
      deep: true,
      handler () {
        const selected = this.visibleLanguages.map(({ tag }) => tag).join(',')
        window.localStorage.setItem(lsKey, selected)
      },
    },
  },

  methods: {
    /**
     * Returns unique set of resources from the given values
     *
     * @returns object[]
     */
    resources () {
      return this.intTranslations
        .map(r => r.resource)
        .filter((r, i, rr) => rr.indexOf(r) === i)
        .map(resource => ({
          resource,
          title: this.titles[resource],
          isPrimary: resource === this.primaryResource,
        }))
    },

    /**
     * Returns unique set of keys from the given values for a specific resource
     *
     * @returns string[]
     */
    keys (resource) {
      return this.intTranslations
        .filter(r => r.resource === resource)
        .map(r => r.key)
        .filter((r, i, rr) => rr.indexOf(r) === i)
    },

    find (resource, key, lang) {
      return this.intTranslations.find(r => r.resource === resource && r.key === key && r.lang === lang) ||
          { dirty: false, message: '' }
    },

    isDirty (resource, key, lang) {
      return this.find(resource, key, lang).dirty
    },

    msg (resource, key, lang) {
      return this.find(resource, key, lang).message
    },

    reset (resource, key, lang) {
      const t = this.find(resource, key, lang)

      if (t.dirty) {
        t.message = t.org
        t.dirty = false
      }
    },

    onUpdate (resource, key, lang, message) {
      message = this.stripHtml(message)

      const v = this.intTranslations.find(r => r.resource === resource && r.key === key && r.lang === lang)
      if (v === undefined) {
        const fresh = { resource, key, lang, message, org: message, dirty: true }
        this.intTranslations.push(fresh)
      } else {
        // if new message is different as original, mark translation as dirty
        v.dirty = v.org !== message
        v.message = message
      }

      const dirty = this.intTranslations
        .filter(({ dirty }) => dirty)
        .map(({ resource, key, lang, message }) => ({ resource, key, lang, message }))

      this.$emit('change', dirty)
    },

    stripHtml (v) {
      const el = document.createElement('div')
      el.innerHTML = v
      return el.textContent || el.innerText || ''
    },
  },
}

</script>
<style lang="scss" scoped>
.key {
  min-width: 200px;
}
</style>
