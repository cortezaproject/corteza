<template>
  <div
    class="d-flex flex-column gap-1"
    @click="handleRecordNavigation"
  >
    <h5 class="font-weight-bold text-break">
      {{ title }}
    </h5>

    <div
      class="text-secondary mb-1 text-break"
    >
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

      const namespace = await this.$ComposeAPI.namespaceRead({ namespaceID })

      if (!namespace) {
        return
      }

      const slug = namespace.slug || namespace.namespaceID

      const recordPages = await this.$ComposeAPI.pageList({ moduleID, namespaceID }).then(({ set = [] }) => set)

      if (!recordPages || recordPages.length === 0) {
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

      // If name and params match, make sure to refresh page instead of pushing
      const sameRoute = routeName === this.$route.name &&
                      slug === this.$route.params.slug &&
                      pageID === this.$route.params.pageID &&
                      recordID === this.$route.params.recordID

      if (openMode === 'newTab') {
        window.open(this.$router.resolve(routeParams).href, '_blank')
      } else if (this.isOnPagesRouteOrChild && !sameRoute && openMode === 'modal') {
        this.$root.$emit('show-record-modal', {
          recordID: !recordID || recordID === '0' ? '0' : recordID,
            recordPageID: pageID,
            edit,
        })

        return
      } else {
        this.$router.push(routeParams)
      }
    },
  },
}
</script>
