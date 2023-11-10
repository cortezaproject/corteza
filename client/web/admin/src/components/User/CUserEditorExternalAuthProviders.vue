<template>
  <b-card
    data-test-id="card-external-auth-providers"
    class="shadow-sm"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-table
      :items="value"
      :fields="fields"
    >
      <template #cell(editor)="{ item }">
        <c-input-confirm
          data-test-id="button-remove-provider"
          show-icon
          @confirmed="$emit('delete', item.credentialsID)"
        />
      </template>
    </b-table>
    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>
  </b-card>
</template>

<script>
export default {
  name: 'CUserEditorExternalAuthProviders',

  i18nOptions: {
    namespaces: 'system.users',
    keyPrefix: 'editor.external-auth-providers',
  },

  props: {
    value: {
      type: Array,
      required: true,
    },
  },

  computed: {
    fields () {
      return [
        { key: 'label', label: this.$t('label'), thStyle: { width: '350px' } },
        { key: 'type', label: this.$t('type'), thStyle: { width: '250px' } },
        { key: 'editor', label: '', tdClass: 'text-right' },
      ]
    },
  },
}
</script>
