<template>
  <b-card
    v-if="value && value.length > 0"
    data-test-id="card-external-auth-providers"
    header-class="border-bottom"
    body-class="p-0"
    class="shadow-sm"
  >
    <template #header>
      <h4 class="m-0">
        {{ $t('title') }}
      </h4>
    </template>

    <b-table
      :items="value"
      :fields="fields"
      head-variant="light"
      responsive
      hover
      class="mb-0"
      style="min-height: 200px;"
    >
      <template #cell(editor)="{ item }">
        <c-input-confirm
          data-test-id="button-remove-provider"
          show-icon
          @confirmed="$emit('delete', item.credentialsID)"
        />
      </template>
    </b-table>
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
