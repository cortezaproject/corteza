<template>
  <div
    class="d-flex flex-column"
    @click="handleRecordNavigation"
  >
    <h5 class="font-weight-bold text-break">
      {{ title }}
    </h5>

    <div class="text-secondary mb-1 text-break">
      {{ description }}
    </div>
  </div>
</template>

<script>
export default {
  props: {
    notification: {
      type: Object,
      required: true,
    },
  },

  computed: {
    title () {
      if (!this.notification || !this.notification.config) {
        return ''
      }

      return this.notification.config.title
    },

    description () {
      if (!this.notification || !this.notification.config) {
        return ''
      }

      return this.notification.config.description || ''
    },

    isOnPagesRouteOrChild () {
      // Check if route exists and is 'pages' or starts with 'page.'
      return this.$route && (['pages', 'page', 'page.record', 'page.record.edit', 'page.record.create'].includes(this.$route.name))
    },
  },

  methods: {
    async handleRecordNavigation () {
      const { namespaceID, recordID, moduleID, openMode, edit } = this.notification.config

      try {
        const namespace = await this.$ComposeAPI.namespaceRead({ namespaceID })

        if (!namespace) {
          this.toastDanger(this.$t('namespaceNotFound'))
          return
        }

        const slug = namespace.slug || namespace.namespaceID

        const recordPages = await this.$ComposeAPI.pageList({ moduleID, namespaceID }).then(({ set = [] }) => set)

        if (!recordPages || recordPages.length === 0) {
          this.toastDanger(this.$t('pageNotFound'))
          return
        }

        const record = await this.$ComposeAPI.recordRead({ recordID, moduleID, namespaceID })

        if (!record) {
          this.toastDanger(this.$t('recordNotFound'))
          return
        }

        const { pageID } = recordPages[0]

        if (this.$router.app.$options.name !== 'compose') {
          const u = new URL(window.location)
          const url = `${u.origin}/compose/ns/${slug}/pages/${pageID}/record/${recordID}/${edit ? 'edit' : ''}`

          if (openMode === 'newTab') {
            window.open(url, '_blank')
          } else {
            window.location = url
          }

          return
        }

        let routeName = 'page.record'

        if (!recordID || recordID === '0') {
          routeName += '.create'
        } else if (edit) {
          routeName += '.edit'
        }

        const routeParams = {
          name: routeName,
          params: {
            recordID,
            pageID,
            slug,
            edit,
          },
        }

        if (openMode === 'newTab') {
          window.open(this.$router.resolve(routeParams).href, '_blank')
        } else if (this.isOnPagesRouteOrChild && openMode === 'modal' && slug === this.$route.params.slug) {
          this.$root.$emit('show-record-modal', {
            recordID: !recordID || recordID === '0' ? '0' : recordID,
            recordPageID: pageID,
            edit,
          })

          return
        } else {
          this.$router.push(routeParams)
        }
      } catch (error) {
        this.toastErrorHandler(this.$t('recordRedirectError'))(error)
      }
    },
  },
}
</script>
