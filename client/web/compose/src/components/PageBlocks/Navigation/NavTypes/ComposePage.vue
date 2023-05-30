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
    <td style="min-width: 200px;">
      <b-form-group :label="$t('navigation.composePage')">
        <vue-select
          key="pageID"
          v-model="options.item.pageID"
          :placeholder="$t('navigation.none')"
          :options="tree"
          append-to-body
          :get-option-key="getOptionKey"
          label="title"
          :calculate-position="calculateDropdownPosition"
          :reduce="f => f.pageID"
          option-value="pageID"
          option-text="title"
          class="nav-page-selector bg-white"
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
import { NoID, compose } from '@cortezaproject/corteza-js'

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
        .pageList({ namespaceID, sort: 'title' })
        .then(({ set: tree }) => {
          this.tree = tree.map(p => new compose.Page(p)).filter(p => p.moduleID === NoID)
        })
        .catch(this.toastErrorHandler(this.$t('notification:page.loadFailed')))
    },

    updateLabelValue (pageID) {
      if (!this.options.item.label) {
        const composePage = this.tree.find(t => t.pageID === pageID)
        this.options.item.label = composePage ? composePage.title : ''
      }
    },

    getOptionKey ({ pageID }) {
      return pageID
    },
  },
}
</script>

<style lang="scss" scoped>

.nav-page-selector {
  position: relative;

  &:not(.vs--open) .vs__selected + .vs__search {
    // force this to not use any space
    // we still need it to be rendered for the focus
    width: 0;
    padding: 0;
    margin: 0;
    border: none;
    height: 0;
  }

  .vs__selected-options {
    // do not allow growing
    width: 0;
  }

  .vs__selected {
    display: block;
    white-space: nowrap;
    text-overflow: ellipsis;
    max-width: 100%;
    overflow: hidden;
  }
}

th,
td {
  padding-left: 15px;
  padding-right: 15px;
}
</style>

<style lang="css">

.vs__dropdown-menu {
  min-width: auto;
}

.vs__dropdown-menu .vs__dropdown-option {
  text-overflow: ellipsis;
  overflow: hidden !important;
}
</style>
