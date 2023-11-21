<template>
  <c-translator-button
    v-if="canManageResourceTranslations && resourceTranslationsEnabled"
    v-bind="$props"
    :tooltip="$t('tooltip')"
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

    pageLayouts: {
      type: Array,
      default: () => [],
    },

    pageLayout: {
      type: compose.PageLayout,
      required: false,
      default: undefined,
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

      if (this.pageLayout) {
        const { namespaceID, pageID, pageLayoutID, handle, meta } = this.pageLayout
        titles[`compose:page-layout/${namespaceID}/${pageID}/${pageLayoutID}`] = this.$t('layout.title', { handle: handle || meta.title || pageLayoutID })
      } else {
        this.pageLayouts.forEach(({ namespaceID, pageID, pageLayoutID, handle, meta }) => {
          titles[`compose:page-layout/${namespaceID}/${pageID}/${pageLayoutID}`] = this.$t('layout.title', { handle: handle || meta.title || pageLayoutID })
        })
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

            if (this.pageLayout) {
              set = set.filter(({ resource }) => resource.endsWith(`${pageID}`) || resource.endsWith(`/${this.pageLayout.pageLayoutID}`))
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

            const layoutResource = ({ pageLayoutID, pageID, namespaceID }) => `compose:page-layout/${namespaceID}/${pageID}/${pageLayoutID}`

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
              if (block.blockID === NoID) return block

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

            const updatePageLayoutTranslations = pageLayout => {
              if (pageLayout.pageLayoutID === NoID) return pageLayout

              const find = (key) => {
                return translations.find(t => t.key === key && t.lang === this.currentLanguage && t.resource === layoutResource(pageLayout))
              }

              let tr = find('title')
              if (tr !== undefined) {
                pageLayout.meta.title = tr.message
              }

              tr = find('description')
              if (tr !== undefined) {
                pageLayout.meta.description = tr.message
              }

              // Refresh page buttons for record pages
              if (pageLayout.moduleID && pageLayout.moduleID !== NoID) {
                tr = find('config.buttons.new.label')
                if (tr) {
                  pageLayout.config.buttons.new.label = tr.message
                }

                tr = find('config.buttons.edit.label')
                if (tr) {
                  pageLayout.config.buttons.edit.label = tr.message
                }

                tr = find('config.buttons.submit.label')
                if (tr) {
                  pageLayout.config.buttons.submit.label = tr.message
                }

                tr = find('config.buttons.delete.label')
                if (tr) {
                  pageLayout.config.buttons.delete.label = tr.message
                }

                tr = find('config.buttons.clone.label')
                if (tr) {
                  pageLayout.config.buttons.clone.label = tr.message
                }

                tr = find('config.buttons.back.label')
                if (tr) {
                  pageLayout.config.buttons.back.label = tr.message
                }
              }

              return pageLayout
            }

            if (this.block) {
              this.block = updateBlockTranslations(this.block)
            } else {
              this.page.blocks = this.page.blocks.map(block => updateBlockTranslations(block))
            }

            if (this.pageLayout) {
              this.pageLayout = updatePageLayoutTranslations(this.pageLayout)
            } else {
              this.pageLayouts = this.pageLayouts.map(pageLayout => updatePageLayoutTranslations(pageLayout))
            }

            return this.page
          })
          .then(page => {
            this.$emit('update:page', page)

            if (this.block) {
              this.$emit('update:block', this.block)
            }

            if (this.pageLayout) {
              this.$emit('update:pageLayout', this.pageLayout)
            } else {
              this.$emit('update:pageLayouts', this.pageLayouts)
            }
          })
      }
    },
  },
}
</script>
