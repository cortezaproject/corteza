<template>
  <b-card
    class="shadow-sm"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-form
      @submit.prevent="$emit('submit', route)"
    >
      <b-form-group
        v-if="route.routeID"
        :label="$t('id')"
        label-cols="2"
      >
        <b-form-input
          v-model="route.routeID"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        label-cols="2"
        :description="routeEndpointDescription"
      >
        <template
          #label
        >
          <label
            label-for="endpoint"
          >
            {{ $t('endpoint') }}
          </label>

          <font-awesome-icon
            id="endpoint_info"
            class="ml-1"
            :icon="['far', 'question-circle']"
          />
          <b-tooltip
            target="endpoint_info"
            triggers="hover"
          >
            {{ $t('tooltip') }}
          </b-tooltip>
        </template>

        <b-form-input
          id="endpoint"
          v-model="route.endpoint"
          :state="isValidEndpoint"
          required
        />
      </b-form-group>

      <b-form-group
        :label="$t('method')"
        label-cols="2"
      >
        <b-form-select
          v-model="route.method"
          :options="methods"
          required
        />
      </b-form-group>

      <b-form-group
        label-cols="2"
        :class="{ 'mb-0': !route.routeID }"
      >
        <b-form-checkbox
          v-model="route.enabled"
        >
          {{ $t('enabled') }}
        </b-form-checkbox>
      </b-form-group>

      <b-form-group
        v-if="route.updatedAt"
        :label="$t('updatedAt')"
        label-cols="2"
      >
        <b-form-input
          :value="route.updatedAt | locFullDateTime"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="route.deletedAt"
        :label="$t('deletedAt')"
        label-cols="2"
      >
        <b-form-input
          :value="route.deletedAt | locFullDateTime"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="route.createdAt"
        :label="$t('createdAt')"
        label-cols="2"
        class="mb-0"
      >
        <b-form-input
          :value="route.createdAt | locFullDateTime"
          plaintext
          disabled
        />
      </b-form-group>
    </b-form>

    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>

    <template #footer>
      <c-submit-button
        class="float-right"
        :processing="processing"
        :success="success"
        :disabled="saveDisabled"
        @submit="$emit('submit', route)"
      />

      <confirmation-toggle
        v-if="route && route.routeID && route.canDeleteApigwRoute"
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
  name: 'CRouteEditorInfo',

  i18nOptions: {
    namespaces: [ 'system.apigw' ],
    keyPrefix: 'editor.info',
  },

  components: {
    ConfirmationToggle,
    CSubmitButton,
  },

  props: {
    route: {
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
      required: true,
    },
  },

  data () {
    return {
      methods: ['GET', 'POST', 'PUT', 'DELETE'],
    }
  },

  computed: {
    fresh () {
      return !this.route.routeID || this.route.routeID === NoID
    },

    editable () {
      return this.fresh ? this.canCreate : true // this.route.canUpdateRoute
    },

    saveDisabled () {
      return !this.editable || [this.isValidEndpoint].includes(false)
    },

    getDeleteStatus () {
      return this.route.deletedAt ? this.$t('undelete') : this.$t('delete')
    },

    isValidEndpoint () {
      const { endpoint } = this.route

      return (!!endpoint && /^(\/[\w-]+)+$/.test(endpoint)) ? null : false
    },

    startsWithSlash () {
      return this.route.endpoint ? /^\//.test(this.route.endpoint) : null
    },

    routeEndpointDescription () {
      if (this.isValidEndpoint === false) {
        if (!this.startsWithSlash) {
          return this.$t('validation.slash')
        } else if (this.route.endpoint.length < 2) {
          return this.$t('validation.minLength')
        } else if (!/^([\w-]+)+$/.test(this.route.endpoint)) {
          return this.$t('validation.invalidCharacters')
        }
      }
      return ''
    },
  },
}
</script>
