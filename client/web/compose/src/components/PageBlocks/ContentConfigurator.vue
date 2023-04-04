<template>
  <b-tab :title="$t('content.label')">
    <b-button
      variant="primary"
      @click="editBlockTest"
    >
      Edit Block
    </b-button>

    <b-form-group>
      <c-rich-text-input
        v-model="options.body"
        :labels="{
          urlPlaceholder: $t('content.urlPlaceholder'),
          ok: $t('content.ok'),
          openLinkInNewTab: $t('content.openLinkInNewTab'),
        }"
      />
    </b-form-group>
  </b-tab>
</template>
<script>
import base from './base'
import { components } from '@cortezaproject/corteza-vue'
const { CRichTextInput } = components

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'Content',

  components: {
    CRichTextInput,
  },

  extends: base,

  methods: {
    editBlockTest () {
      this.$root.$emit('tab-editRequest', '4' /* Your block ID goes here for manual testing */)
    },
  },
}
</script>
