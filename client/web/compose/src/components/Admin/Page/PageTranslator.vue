<template>
  <c-translator-button
    v-if="canManageResourceTranslations && resourceTranslationsEnabled"
    v-bind="$props"
    :resource="resource"
    :titles="titles"
    :fetcher="fetcher"
    :updater="updater"
  />
</template>

<script>
import { compose, NoID } from '@cortezaproject/corteza-js'
import { mapGetters } from 'vuex'
import CTranslatorButton from 'corteza-webapp-compose/src/components/Translator/CTranslatorButton'

export default {
  components: {
    CTranslatorButton,
  },

  i18nOptions: {
    namespaces: 'resource-translator',
    keyPrefix: 'resources.page',
  },

  props: {
    page: {
      type: compose.Page,
      required: true,
    },

    block: {
      type: compose.PageBlock,
      required: false,
      default: undefined,
    },

    buttonVariant: {
      type: String,
      default: () => 'primary',
    },

    highlightKey: {
      type: String,
      default: '',
    },

    disabled: {
      type: Boolean,
      default: () => false,
    },
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canManageResourceTranslations () {
      return this.can('compose/', 'resource-translations.manage')
    },

    resource () {
      const { pageID, namespaceID } = this.page
      return `compose:page/${namespaceID}/${pageID}`
    },

    titles () {
      const titles = {}
      if (this.block) {
        const { title, blockID } = this.block

        titles[this.resource] = this.$t('block.title', { title, blockID })
      } else {
        const { pageID, handle } = this.page

        titles[this.resource] = this.$t('title', { handle: handle || pageID })
      }

      return titles
    },

    fetcher () {
      const { pageID, namespaceID } = this.page

      return () => {
        return this.$ComposeAPI
          .pageListTranslations({ namespaceID, pageID })
          .then(set => {
            if (this.block) {
              /**
               * When block is set, intercept the resolved request and filter out the
               * translations that are relevant for that block
               */
              set = set.filter(({ key }) => key.startsWith(`pageBlock.${this.block.blockID}.`))
            }

            return set
          })
      }
    },

    updater () {
      const { pageID, namespaceID } = this.page

      return translations => {
        return this.$ComposeAPI
          .pageUpdateTranslations({ namespaceID, pageID, translations })
          // re-fetch translations, sanitized and stripped
          .then(() => this.fetcher())
          .then((translations) => {
            // When translations are successfully saved,
            // scan changes and apply them back to the passed object
            // not the most elegant solution but is saves us from
            // handling the resource on multiple places
            //
            // @todo move this to Namespace* classes
            // the logic there needs to be implemented; the idea is to encode
            // values from the set of translations back to the resource object

            const find = (key) => {
              return translations.find(t => t.key === key && t.lang === this.currentLanguage && t.resource === this.resource)
            }

            let tr = find('title')
            if (tr !== undefined) {
              this.page.title = tr.message
            }

            tr = find('description')
            if (tr !== undefined) {
              this.page.description = tr.message
            }

            // Refresh page buttons for record pages
            if (this.page.moduleID && this.page.moduleID !== NoID) {
              tr = find('recordToolbar.new.label')
              if (tr) {
                this.$set(this.page.config.buttons.new, 'label', tr.message)
              }

              tr = find('recordToolbar.edit.label')
              if (tr) {
                this.$set(this.page.config.buttons.edit, 'label', tr.message)
              }

              tr = find('recordToolbar.submit.label')
              if (tr) {
                this.$set(this.page.config.buttons.submit, 'label', tr.message)
              }

              tr = find('recordToolbar.delete.label')
              if (tr) {
                this.$set(this.page.config.buttons.delete, 'label', tr.message)
              }

              tr = find('recordToolbar.clone.label')
              if (tr) {
                this.$set(this.page.config.buttons.clone, 'label', tr.message)
              }

              tr = find('recordToolbar.back.label')
              if (tr) {
                this.$set(this.page.config.buttons.back, 'label', tr.message)
              }
            }

            const updateBlockTranslations = block => {
              block.title = (find(`pageBlock.${block.blockID}.title`) || {}).message
              block.description = (find(`pageBlock.${block.blockID}.description`) || {}).message

              switch (true) {
                case block instanceof compose.PageBlockAutomation:
                  block.options.buttons.forEach((btn, index) => {
                    tr = find(`pageBlock.${block.blockID}.button.${btn.buttonID || index}.label`)
                    if (tr) {
                      btn.label = tr.message
                    }
                  })
                  break

                case block instanceof compose.PageBlockRecordList:
                  block.options.selectionButtons.forEach((btn, index) => {
                    tr = find(`pageBlock.${block.blockID}.button.${btn.buttonID || index}.label`)
                    if (tr) {
                      btn.label = tr.message
                    }
                  })
                  break

                case block instanceof compose.PageBlockContent:
                  tr = find(`pageBlock.${block.blockID}.content.body`)
                  if (tr) {
                    block.options.body = tr.message
                  }
                  break
              }

              return block
            }

            if (this.block) {
              this.block = updateBlockTranslations(this.block)
            }

            this.page.blocks = this.page.blocks.map(block => updateBlockTranslations(block))

            return this.page
          })
          .then(page => {
            this.$emit('update:page', page)

            if (this.block) {
              this.$emit('update:block', this.block)
            }
          })
      }
    },
  },
}
</script>
