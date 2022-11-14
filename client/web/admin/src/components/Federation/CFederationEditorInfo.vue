<template>
  <b-card
    class="shadow-sm"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-form
      @submit.prevent="$emit('submit', node)"
    >
      <b-form-group
        :label="$t('name')"
        label-cols="2"
      >
        <b-form-input
          v-model="node.name"
          :state="nameState"
        />
      </b-form-group>

      <b-form-group
        :label="$t('url')"
        label-cols="2"
      >
        <b-form-input
          v-model="node.baseURL"
          placeholder="https://example.com/federation"
          type="url"
          :state="urlState"
        />
      </b-form-group>

      <b-form-group
        :label="$t('email')"
        label-cols="2"
      >
        <b-form-input
          v-model="node.contact"
          placeholder="email@example.com"
          type="email"
        />
      </b-form-group>

      <!-- <b-form-group
        :label="$t('tags.label')"
        label-cols="2"
      >
        <b-form-tags
          v-model="node.tags"
          tag-variant="warning"
          tag-class="rounded"
          input-class="h4"
          :placeholder="$t('tags.placeholder')"
          size="lg"
          class="py-1 px-2"
        />
      </b-form-group> -->

      <b-form-group
        v-if="node.status"
        :label="$t('status')"
        label-cols="2"
      >
        <b-form-input
          v-model="node.status"
          plaintext
          disabled
        />
      </b-form-group>

      <!-- <b-form-group
        label-cols="2"
        :class="{ 'mb-0': !node.nodeID }"
      >
        <b-form-checkbox
          v-model="node.enabled"
        >
          {{ $t('enabled') }}
        </b-form-checkbox>
      </b-form-group> -->

      <b-form-group
        v-if="node.updatedAt"
        :label="$t('updatedAt')"
        label-cols="2"
      >
        <b-form-input
          :value="node.updatedAt | locFullDateTime"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="node.deletedAt"
        :label="$t('deletedAt')"
        label-cols="2"
      >
        <b-form-input
          :value="node.deletedAt | locFullDateTime"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="node.createdAt"
        :label="$t('createdAt')"
        label-cols="2"
        class="mb-0"
      >
        <b-form-input
          :value="node.createdAt | locFullDateTime"
          plaintext
          disabled
        />
      </b-form-group>

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

    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>

    <template #footer>
      <c-submit-button
        class="float-right"
        :disabled="saveDisabled"
        :processing="processing"
        :success="success"
        @submit="$emit('submit', node)"
      />

      <confirmation-toggle
        v-if="node && node.nodeID"
        @confirmed="$emit('delete')"
      >
        {{ getDeleteStatus }}
      </confirmation-toggle>
    </template>
  </b-card>
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'
import ConfirmationToggle from 'corteza-webapp-admin/src/components/ConfirmationToggle'
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'

export default {
  name: 'CFederationEditorInfo',

  i18nOptions: {
    namespaces: 'federation.nodes',
    keyPrefix: 'editor.info',
  },

  components: {
    ConfirmationToggle,
    CSubmitButton,
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
      return !this.editable && [this.nameState, this.urlState].includes(false)
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
