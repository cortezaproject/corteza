<template>
  <b-tab :title="$t('label')">
    <b-row>
      <b-col
        cols="12"
        lg="6"
      >
        <b-form-group
          :label="$t('kind')"
          label-class="text-primary"
        >
          <b-form-select
            v-model="options.kind"
            :options="nylasComponentKinds"
          />
        </b-form-group>
      </b-col>

      <b-col
        cols="12"
        lg="6"
      >
        <b-form-group
          :label="$t('componentID')"
          label-class="text-primary"
        >
          <b-form-input
            v-model="options.componentID"
          />
        </b-form-group>
      </b-col>

      <b-col
        cols="12"
        lg="6"
      >
        <b-form-group
          :label="$t('accessTokenRequired')"
          label-class="text-primary"
        >
          <c-input-checkbox
            v-model="accessTokenRequired"
            switch
            :labels="checkboxLabels"
          />
        </b-form-group>
      </b-col>
    </b-row>

    <template v-if="showPreviewSection">
      <hr class="my-3">

      <div>
        <h5 class="mb-3">
          {{ $t('prefill.title') }}
        </h5>

        <b-row>
          <template v-if="options.kind === 'Composer'">
            <b-col
              cols="12"
              sm="6"
            >
              <b-form-group
                :label="$t('prefill.to')"
                label-class="text-primary"
              >
                <c-input-select
                  v-model="options.prefill.to"
                  :options="moduleTextFields"
                  :get-option-label="getFieldLabel"
                  :get-option-key="getOptionKey"
                  :placeholder="$t('prefill.selectField')"
                  :reduce="field => field.fieldID"
                  class="bg-white rounded"
                />
              </b-form-group>
            </b-col>
            <b-col
              cols="12"
              sm="6"
            >
              <b-form-group
                :label="$t('prefill.subject')"
                label-class="text-primary"
              >
                <c-input-select
                  v-model="options.prefill.subject"
                  :options="moduleTextFields"
                  :get-option-label="getFieldLabel"
                  :get-option-key="getOptionKey"
                  :placeholder="$t('prefill.selectField')"
                  :reduce="getOptionKey"
                  :calculate-position="calculateDropdownPosition"
                />
              </b-form-group>
            </b-col>
            <b-col
              cols="12"
              sm="6"
            >
              <b-form-group
                :label="$t('prefill.body')"
                label-class="text-primary"
              >
                <c-input-select
                  v-model="options.prefill.body"
                  :options="moduleTextFields"
                  :get-option-label="getFieldLabel"
                  :get-option-key="getOptionKey"
                  :placeholder="$t('prefill.selectField')"
                  :reduce="getOptionKey"
                  :calculate-position="calculateDropdownPosition"
                />
              </b-form-group>
            </b-col>
          </template>

          <template v-if="options.kind === 'Mailbox'">
            <b-col
              cols="12"
              sm="6"
            >
              <b-form-group
                :label="$t('prefill.queryString')"
                label-class="text-primary"
              >
                <c-input-select
                  v-model="options.prefill.queryString"
                  :options="moduleTextFields"
                  :get-option-label="getFieldLabel"
                  :get-option-key="getOptionKey"
                  :placeholder="$t('prefill.selectField')"
                  :reduce="getOptionKey"
                  :calculate-position="calculateDropdownPosition"
                />
              </b-form-group>
            </b-col>
          </template>
        </b-row>
      </div>
    </template>
  </b-tab>
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'
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

      checkboxLabels: {
        on: this.$t('general:label.yes'),
        off: this.$t('general:label.no'),
      },
    }
  },

  computed: {
    moduleTextFields () {
      return this.module.fields.filter(({ kind }) => ['String', 'Email'].includes(kind))
    },

    showPreviewSection () {
      return ['Composer', 'Mailbox'].includes(this.options.kind) && (this.page.moduleID !== NoID)
    },

    accessTokenRequired: {
      get () {
        return !!this.options.accessTokenRequired
      },

      set (v) {
        this.options.accessTokenRequired = v
      },
    },
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    getFieldLabel ({ name, label }) {
      return label || name
    },

    getOptionKey ({ fieldID, name }) {
      return fieldID !== NoID ? fieldID : name
    },

    setDefaultValues () {
      this.nylasComponentKinds = []
      this.checkboxLabels = {}
    },
  },
}
</script>
