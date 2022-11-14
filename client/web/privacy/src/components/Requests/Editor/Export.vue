<template>
  <div
    class="d-flex flex-column h-100"
  >
    <b-card
      class="shadow-sm mb-3"
    >
      <b-form-group
        :label="$t('data-type.label')"
        label-class="text-primary"
      >
        <b-form-checkbox
          v-model="payload.profile"
          class="ml-2 mb-1"
        >
          {{ $t('data-type.profile-information') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="payload.application"
          class="ml-2"
        >
          {{ $t('data-type.application-data') }}
        </b-form-checkbox>
      </b-form-group>

      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('date-range.label')"
            label-class="text-primary"
          >
            <b-form-select
              v-model="payload.range"
              :options="rangeOptions"
            />
          </b-form-group>
        </b-col>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('file-format.label')"
            label-class="text-primary"
          >
            <b-form-select
              v-model="payload.format"
              :options="formatOptions"
            />
          </b-form-group>
        </b-col>
      </b-row>
    </b-card>

    <portal to="editor-toolbar">
      <editor-toolbar
        :processing="processing"
        :back-link="{ name: 'data-overview.application' }"
        submit-show
        :submit-label="$t('submit')"
        :submit-disabled="!(payload.profile || payload.application)"
        @submit="$emit('submit', { kind: 'export', payload })"
      />
    </portal>
  </div>
</template>

<script>
import EditorToolbar from 'corteza-webapp-privacy/src/components/Common/EditorToolbar'

export default {
  name: 'ApplicationDataDelete',

  i18nOptions: {
    namespaces: 'request',
    keyPrefix: 'edit.export',
  },

  components: {
    EditorToolbar,
  },

  data () {
    return {
      processing: false,

      payload: {
        profile: false,
        application: false,
        range: 'all',
        format: 'json',
      },

      rangeOptions: [
        { text: this.$t('date-range.all'), value: 'all' },
      ],

      formatOptions: [
        { text: this.$t('file-format.json'), value: 'json' },
        { text: this.$t('file-format.csv'), value: 'csv' },
      ],
    }
  },
}
</script>
