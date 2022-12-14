<template>
  <b-row>
    <b-col
      cols="12"
      sm="6"
    >
      <b-form-group :label="$t('navigation.text')">
        <b-form-input
          v-model="column.options.itemOption.text"
          type="text"
          size="sm"
        />
      </b-form-group>
    </b-col>
    <b-col
      cols="12"
      sm="6"
    >
      <b-form-group :label="$t('navigation.composePage')">
        <b-form-select
          v-model="column.options.itemOption.referenceId"
          :options="composePageList"
          size="sm"
        />
      </b-form-group>
    </b-col>
  </b-row>
</template>

<script>
import { compose } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  props: {
    column: {
      type: Object,
      required: true,
    },

    namespace: {
      type: compose.Namespace,
      required: true,
    },
  },

  data () {
    return {
      tree: [],
    }
  },

  computed: {
    composePageList () {
      return this.tree.map(t => ({
        value: t.namespaceID,
        text: t.title,
      }))
    },
  },

  created () {
    this.loadTree()
  },

  methods: {
    loadTree () {
      const { namespaceID } = this.namespace
      this.$ComposeAPI
        .pageTree({ namespaceID })
        .then(tree => {
          this.tree = tree.map(p => new compose.Page(p))
        })
        .catch(this.toastErrorHandler(this.$t('notification:page.loadFailed')))
    },
  },
}
</script>
