<template>
  <div>
    <b-modal
      v-model="showModal"
      hide-footer
      size="xl"
      :title="translatedTitle"
      lazy
      scrollable
      class="overflow-hidden"
      @hide="onHide"
    >
      <c-permissions-form
        v-if="resource"
        :resource="resource"
        :target="target"
        :all-specific="allSpecific"
        class="h-100"
      />
    </b-modal>
  </div>
</template>
<script lang="js">
import {modalOpenEventName, split} from './def.ts'
import CPermissionsForm from './CPermissionsForm.vue'

export default {
  i18nOptions: {
    namespaces: 'permissions',
  },

  components: {
    CPermissionsForm,
  },

  data () {
    return {
      resource: undefined,
      title: undefined,
      target: undefined,
      allSpecific: false,
    }
  },

  computed: {
    showModal: {
      get () {
        return !!this.resource
      },

      set (open) {
        if (!open) {
          this.clear()
        }
      },
    },

    translatedTitle () {
      if (this.resource) {
        const { i18nPrefix } = split(this.resource)

        let target
        if (this.allSpecific) {
          target = this.$t(`permissions:${i18nPrefix}.all-specific`, { target: this.title, interpolation: { escapeValue: false } })
        } else if (this.title) {
          target = this.$t(`permissions:${i18nPrefix}.specific`, { target: this.title, interpolation: { escapeValue: false } })
        } else {
          target = this.$t(`permissions:${i18nPrefix}.all`,)
        }

        return this.$t('permissions:ui.set-for', { target: target, interpolation: { escapeValue: false } })
      }

      return undefined
    },
  },

  mounted () {
    this.$root.$on(modalOpenEventName, ({ resource, title, target, allSpecific }) => {
      this.resource = resource
      this.title = title
      this.target = target
      this.allSpecific = allSpecific
    })
  },

  destroyed () {
    this.$root.$off(modalOpenEventName)
  },

  methods: {
    onHide () {
      this.clear()
    },

    clear () {
      this.resource = undefined
      this.title = undefined
      this.target = undefined
    },
  },
}
</script>
