<template>
  <b-tab :title="$t('label')">
    <b-form-group
      :label="$t('kind')"
    >
      <b-form-select
        v-model="options.kind"
        :options="nylasComponentKinds"
      />
    </b-form-group>

    <b-form-group
      :label="$t('componentID')"
    >
      <b-form-input
        v-model="options.componentID"
      />
    </b-form-group>
  </b-tab>
</template>
<script>
import base from '../base'

export default {
  i18nOptions: {
    namespaces: 'block',
    keyPrefix: 'nylas.configurator',
  },

  extends: base,

  data () {
    return {
      nylasComponentKinds: [
        { value: 'Agenda', text: this.$t('kinds.agenda') },
        { value: 'Composer', text: this.$t('kinds.composer') },
        { value: 'ContactList', text: this.$t('kinds.contactList') },
        { value: 'Conversation', text: this.$t('kinds.conversation') },
        { value: 'Email', text: this.$t('kinds.email') },
        { value: 'Mailbox', text: this.$t('kinds.mailbox') },
      ],
    }
  },
}
</script>
