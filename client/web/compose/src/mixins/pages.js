import { NoID } from '@cortezaproject/corteza-js'

export default {
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

          page = {
            ...page,
            pageID: NoID,
            title: newPageTitle,
            handle: '',
          }

          return this.createPage({ namespaceID, ...page })
        })
        .then(({ pageID }) => {
          this.$router.push({ name: this.$route.name, params: { pageID } })
        })
        .catch(this.toastErrorHandler(this.$t('notification:page.cloneFailed')))
    },

  },

}
