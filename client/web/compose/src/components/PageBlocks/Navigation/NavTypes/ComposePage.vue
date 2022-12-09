<template>
  <tr>
    <td />
    <td>
      <b-form-group :label="$t('navigation.fieldLabel')">
        <b-form-input
          v-model="options.item.label"
          type="text"
        />
      </b-form-group>
    </td>
    <td>
      <b-form-group :label="$t('navigation.composePage')">
        <vue-select
          key="pageID"
          v-model="options.item.pageID"
          :placeholder="$t('navigation.none')"
          :options="tree"
          append-to-body
          label="title"
          :reduce="f => f.pageID"
          option-value="pageID"
          option-text="title"
          class="bg-white"
          @input="updateLabelValue"
        />
      </b-form-group>
    </td>
    <td>
      <b-form-group :label="$t('navigation.openIn')">
        <b-form-select
          v-model="options.item.target"
          :options="targetOptions"
        />
      </b-form-group>
    </td>
    <td
      v-if="selectedPageChildren(options.item.pageID)"
      cols="12"
      sm="6"
      class="align-middle text-center"
    >
      <b-form-group class="m-0">
        <b-form-checkbox
          v-model="options.item.displaySubPages"
          switch
          size="sm"
        >
          {{ $t('navigation.displaySubPages') }}
        </b-form-checkbox>
      </b-form-group>
    </td>
    <td />
  </tr>
</template>

<script>
import base from './base'

import { VueSelect } from 'vue-select'
import { compose } from '@cortezaproject/corteza-js'

export default {
  components: {
    VueSelect,
  },

  extends: base,

  props: {
    namespace: {
      type: compose.Namespace,
      required: true,
    },
  },

  data () {
    return {
      tree: [],
      targetOptions: [
        { value: 'sameTab', text: this.$t('navigation.sameTab') },
        { value: 'newTab', text: this.$t('navigation.newTab') },
      ],
    }
  },

  created () {
    this.loadTree()
  },

  methods: {
    selectedPageChildren (pageID) {
      const tree = this.tree ? this.tree.find(t => t.pageID === pageID) : {}
      return tree && tree.children ? tree.children.length > 0 : false
    },

    loadTree () {
      const { namespaceID } = this.namespace
      this.$ComposeAPI
        .pageTree({ namespaceID })
        .then(tree => {
          this.tree = tree.map(p => new compose.Page(p))
        })
        .catch(this.toastErrorHandler(this.$t('notification:page.loadFailed')))
    },

    updateLabelValue (pageID) {
      if (!this.options.item.label) {
        const composePage = this.tree.find(t => t.pageID === pageID)
        this.options.item.label = composePage ? composePage.title : ''
      }
    },
  },
}
</script>
