<template>
  <b-card
    header-bg-variant="white"
    footer-bg-variant="white"
    footer-class="d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>

    <b-form
      @submit.prevent="$emit('submit', node)"
    >
      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('name')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="node.name"
              :state="nameState"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('url')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="node.baseURL"
              placeholder="https://example.com/federation"
              type="url"
              :state="urlState"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('email')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="node.contact"
              placeholder="email@example.com"
              type="email"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            v-if="node.status"
            :label="$t('status')"
            label-class="text-primary"
          >
            {{ node.status }}
          </b-form-group>
        </b-col>
      </b-row>

      <c-system-fields
        :resource="node"
      />

      <!--
        include hidden input to enable
        trigger submit event w/ ENTER
      -->
      <input
        type="submit"
        class="d-none"
        :disabled="saveDisabled"
      >
    </b-form>

    <template #footer>
      <confirmation-toggle
        v-if="node && node.nodeID"
        @confirmed="$emit('delete')"
      >
        {{ getDeleteStatus }}
      </confirmation-toggle>

      <c-button-submit
        :disabled="saveDisabled"
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="$emit('submit', node)"
      />
    </template>
  </b-card>
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'
import ConfirmationToggle from 'corteza-webapp-admin/src/components/ConfirmationToggle'

export default {
  name: 'CFederationEditorInfo',

  i18nOptions: {
    namespaces: 'federation.nodes',
    keyPrefix: 'editor.info',
  },

  components: {
    ConfirmationToggle,
  },

  props: {
    node: {
      type: Object,
      required: true,
    },

    processing: {
      type: Boolean,
      value: false,
    },

    success: {
      type: Boolean,
      value: false,
    },

    canCreate: {
      type: Boolean,
      value: false,
    },
  },

  computed: {
    fresh () {
      return !this.node.nodeID || this.node.nodeID === NoID
    },

    editable () {
      return this.fresh ? this.canCreate : this.node.canManageNode
    },

    saveDisabled () {
      return !this.editable || [this.nameState, this.urlState].includes(false)
    },

    nameState () {
      const { name } = this.node
      return name ? null : false
    },

    urlState () {
      const { baseURL = '' } = this.node
      return /https?:\/\/*.*\/federation/gm.test(baseURL) ? null : false
    },

    getDeleteStatus () {
      return this.node.deletedAt ? this.$t('undelete') : this.$t('delete')
    },
  },
}
</script>
