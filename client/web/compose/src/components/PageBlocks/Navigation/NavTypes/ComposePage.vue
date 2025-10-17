<template>
  <tr>
    <td />

    <td>
      <b-form-group
        :label="$t('navigation.fieldLabel')"
        label-class="text-primary"
      >
        <b-form-input
          v-model="options.item.label"
          type="text"
        />
      </b-form-group>
    </td>

    <td style="min-width: 200px;">
      <b-form-group
        :label="$t('navigation.composePage')"
        label-class="text-primary"
      >
        <c-input-select
          v-model="options.item.pageID"
          :placeholder="$t('navigation.none')"
          :options="pageList"
          :get-option-key="getOptionKey"
          label="title"
          :reduce="f => f.pageID"
          option-value="pageID"
          @input="updateLabelValue"
        />
      </b-form-group>
    </td>

    <td style="min-width: 200px;">
      <b-form-group
        :label="$t('navigation.pageLayout')"
        label-class="text-primary"
      >
        <c-input-select
          v-model="options.item.pageLayoutID"
          :placeholder="$t('navigation.defaultLayout')"
          :options="pageLayoutList"
          :get-option-key="getLayoutOptionKey"
          :get-option-label="getLayoutOptionLabel"
          :reduce="f => f.pageLayoutID"
          option-value="pageLayoutID"
          :loading="loadingPageLayouts"
          :disabled="!options.item.pageID"
        />
      </b-form-group>
    </td>

    <td style="min-width: 200px;">
      <b-form-group
        :label="$t('navigation.openIn')"
        label-class="text-primary"
      >
        <b-form-select
          v-model="options.item.target"
          :options="targetOptions"
        />
      </b-form-group>
    </td>

    <td
      v-if="selectedPageChildren(options.item.pageID).length > 0"
      class="align-middle text-center"
    >
      <b-form-group
        :label="$t('navigation.displaySubPages')"
        label-class="text-nowrap text-primary"
      >
        <c-input-checkbox
          v-model="options.item.displaySubPages"
          switch
        />
      </b-form-group>
    </td>

    <td />
  </tr>
</template>

<script>
import base from './base'
import { NoID, compose } from '@cortezaproject/corteza-js'

export default {
  extends: base,

  props: {
    namespace: {
      type: compose.Namespace,
      required: true,
    },
  },

  data () {
    return {
      loadingPageLayouts: false,
      pageList: [],
      pageLayoutList: [],
      checkboxLabel: {
        on: this.$t('general:label.yes'),
        off: this.$t('general:label.no'),
      },
      targetOptions: [
        { value: 'sameTab', text: this.$t('navigation.sameTab') },
        { value: 'newTab', text: this.$t('navigation.newTab') },
      ],
    }
  },

  watch: {
    'options.item.pageID': {
      handler (pageID) {
        if (pageID && pageID !== NoID) {
          this.loadPageLayouts(pageID)
        } else {
          this.pageLayoutList = []
        }
      },
      immediate: true,
    },
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  created () {
    this.loadPages()
  },

  methods: {
    selectedPageChildren (pageID) {
      return this.pageList.filter(value => value.selfID === pageID && value.moduleID === NoID) || []
    },

    loadPages () {
      const { namespaceID } = this.namespace
      this.$ComposeAPI
        .pageList({ namespaceID, sort: 'title' })
        .then(({ set: pages }) => {
          this.pageList = pages.map(p => new compose.Page(p))
        })
        .catch(this.toastErrorHandler(this.$t('notification:page.listFailed')))
    },

    loadPageLayouts (pageID) {
      const { namespaceID } = this.namespace
      this.loadingPageLayouts = true
      this.$ComposeAPI
        .pageLayoutList({ namespaceID, pageID, sort: 'weight ASC' })
        .then(({ set: layouts }) => {
          this.pageLayoutList = layouts.map(pl => new compose.PageLayout(pl))
        })
        .catch(this.toastErrorHandler(this.$t('notification:page-layout.listFailed')))
        .finally(() => {
          this.loadingPageLayouts = false
        })
    },

    updateLabelValue (pageID) {
      const composePage = this.pageList.find(t => t.pageID === pageID)

      if (!this.options.item.label) {
        this.options.item.label = composePage ? composePage.title : ''
      }

      this.options.item.pageLayoutID = ''
      this.options.item.moduleID = composePage && composePage.moduleID !== NoID ? composePage.moduleID : ''
    },

    getOptionKey ({ pageID }) {
      return pageID
    },

    getLayoutOptionKey ({ pageLayoutID }) {
      return pageLayoutID
    },

    getLayoutOptionLabel ({ handle, meta = {}, pageLayoutID }) {
      return meta.title || handle || pageLayoutID
    },

    setDefaultValues () {
      this.pageList = []
      this.pageLayoutList = []
      this.checkboxLabel = {}
      this.targetOptions = []
    },
  },
}
</script>

<style lang="scss" scoped>
th,
td {
  padding-left: 15px;
  padding-right: 15px;
}
</style>
