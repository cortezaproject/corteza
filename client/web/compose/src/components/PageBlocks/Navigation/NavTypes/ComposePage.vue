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
          key="pageID"
          v-model="options.item.pageID"
          :placeholder="$t('navigation.none')"
          :options="pageList"
          :get-option-key="getOptionKey"
          label="title"
          :reduce="f => f.pageID"
          option-value="pageID"
          option-text="title"
          @input="updateLabelValue"
        />
      </b-form-group>
    </td>

    <td>
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
      cols="12"
      sm="6"
      class="align-middle text-center"
    >
      <b-form-group
        :label="$t('navigation.displaySubPages')"
        label-class="text-nowrap text-primary"
        class="d-flex align-items-center justify-content-center mb-0"
      >
        <c-input-checkbox
          v-model="options.item.displaySubPages"
          switch
          size="sm"
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
      pageList: [],
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
          this.pageList = pages.map(p => new compose.Page(p)).filter(p => p.moduleID === NoID)
        })
        .catch(this.toastErrorHandler(this.$t('notification:page.loadFailed')))
    },

    updateLabelValue (pageID) {
      if (!this.options.item.label) {
        const composePage = this.pageList.find(t => t.pageID === pageID)
        this.options.item.label = composePage ? composePage.title : ''
      }
    },

    getOptionKey ({ pageID }) {
      return pageID
    },

    setDefaultValues () {
      this.pageList = []
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
