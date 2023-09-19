import { NoID } from '@cortezaproject/corteza-js'
import { fetchID } from 'corteza-webapp-compose/src/lib/block'

export default {
  data () {
    return {
      processingClone: false,
    }
  },

  methods: {
    countDuplicateTitles (pages, pattern) {
      let substrTitle = this.page.title

      // if the old page begins with 'copy of' and it does not end with a number in brackets, remove copy of prefix
      if (this.page.title.includes(this.$t('copyOf')) && !pattern.test(this.page.title)) {
        substrTitle = this.page.title.substring(this.$t('copyOf').length, this.page.title.length)
      // if the old page begins with 'copy of' and it ends with a number in brackets, remove them
      } else if (this.page.title.includes(this.$t('copyOf')) && pattern.test(this.page.title)) {
        substrTitle = this.page.title.substring(this.$t('copyOf').length, this.page.title.lastIndexOf('(') - 1)
      }

      return pages.filter(p => p.title.includes(substrTitle)).length
    },

    handleClone () {
      this.processingClone = true
      const { namespaceID = NoID } = this.namespace
      const pattern = /\([0-9]+\)/

      let page = this.page.clone()

      this.loadPages({ namespaceID, force: true })
        .then(pages => {
          let newPageTitle = this.$t('copyOf', { title: this.page.title })
          const countDuplicates = this.countDuplicateTitles(pages, pattern)

          // if page with the same name already exists two times, add a number (2) sufix
          if (countDuplicates === 2) {
            newPageTitle = this.page.title.includes(this.$t('copyOf'))
              ? `${this.page.title} (${countDuplicates})`
              : `${this.$t('copyOf')} ${this.page.title} (${countDuplicates})`
          } else if (countDuplicates > 2 && this.page.title.includes(this.$t('copyOf'))) {
            /**
            * if page with the same name exists more than two times and it has "copy of" prefix, check old page's sufix
            * If page ends with a number inside parenthesis, replace it with the new number
            * If page does not have a number, add it
            */
            newPageTitle = pattern.test(page.title)
              ? `${this.page.title.substr(0, this.page.title.lastIndexOf('('))}(${countDuplicates})`
              : `${this.page.title} (${countDuplicates})`
          } else if (countDuplicates > 2 && !this.page.title.includes(this.$t('copyOf'))) {
            // if page with the same name exists more than two times and it does not have "copy of" prefix, append the prefix and the new number
            newPageTitle = `${this.$t('copyOf')}${this.page.title} (${countDuplicates})`
          }

          // Change tabbed blockID to use tempID's since they are persisted on save
          const blocks = page.blocks.map(block => {
            if (block.kind !== 'Tabs') return block
            const { tabs = [] } = block.options

            block.options.tabs = tabs.map(b => {
              const { tempID } = (page.blocks.find(({ blockID }) => blockID === b.blockID) || {}).meta || {}
              b.blockID = tempID
              return b
            })

            return block
          })

          page = {
            ...page,
            blocks,
            pageID: NoID,
            title: newPageTitle,
            handle: '',
          }

          return this.createPage({ namespaceID, ...page }).then((page) => {
            return this.cloneLayouts(page.pageID).then(() => {
              return this.updateTabbedBlockIDs(page)
            })
          })
        }).then(({ pageID }) => {
          this.$router.push({ name: this.$route.name, params: { pageID } })
        })
        .catch(this.toastErrorHandler(this.$t('notification:page.cloneFailed')))
        .finally(() => { this.processingClone = false })
    },

    cloneLayouts (pageID) {
      const layouts = [...this.layouts]
      return Promise.all(layouts.map(layout => {
        layout.pageID = pageID
        layout.pageLayoutID = NoID
        return this.createPageLayout(layout)
      }))
    },

    async updateTabbedBlockIDs (page) {
      // get the Tabs Block that still has tabs with tempIDs
      let updatePage = false

      page.blocks.filter(({ kind }) => kind === 'Tabs')
        .filter(({ options = {} }) => options.tabs.some(({ blockID }) => (blockID || '').startsWith('tempID-')))
        .forEach(b => {
          if (b.kind !== 'Tabs') return

          b.options.tabs.forEach((t, j) => {
            if (!t.blockID.startsWith('tempID-')) return false

            // find a block with the same tempID that should be updated by now and get its blockID
            const updatedBlock = page.blocks.find(block => block.meta.tempID === t.blockID)

            if (!updatedBlock) return false

            const tab = {
              // fetchID gets the blockID using the found block
              blockID: fetchID(updatedBlock),
              title: t.title,
            }

            b.options.tabs.splice(j, 1, tab)
            updatePage = true
          })
        })

      if (!updatePage) {
        return page
      }

      return this.updatePage(page)
    },
  },
}
