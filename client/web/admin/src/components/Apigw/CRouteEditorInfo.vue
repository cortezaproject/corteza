<template>
  <b-card
    data-test-id="card-route-edit"
    header-class="border-bottom"
    footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <template #header>
      <h4 class="m-0">
        {{ $t('title') }}
      </h4>
    </template>

    <b-form
      @submit.prevent="$emit('submit', route)"
    >
      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :description="routeEndpointDescription"
            label-class="text-primary"
          >
            <template
              #label
            >
              <label
                label-for="endpoint"
                class="mb-0"
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
              data-test-id="input-endpoint"
              :state="isValidEndpoint"
              required
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('method')"
            label-class="text-primary"
          >
            <b-form-select
              v-model="route.method"
              data-test-id="select-method"
              :options="methods"
              required
            />
          </b-form-group>
        </b-col>

        <b-col
          v-if="route.meta"
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('description')"
            label-class="text-primary"
          >
            <b-form-textarea
              v-model="route.meta.description"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('enabled')"
            :class="{ 'mb-0': !route.routeID }"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="route.enabled"
              switch
              :labels="checkboxLabel"
              data-test-id="checkbox-enabled"
            />
          </b-form-group>
        </b-col>
      </b-row>

      <c-system-fields
        :id="route.routeID"
        :resource="route"
      />
    </b-form>

    <template #footer>
      <c-input-confirm
        v-if="route && route.routeID && route.canDeleteApigwRoute"
        :data-test-id="deletedButtonStatusCypressId"
        variant="danger"
        size="md"
        @confirmed="$emit('delete')"
      >
        {{ getDeleteStatus }}
      </c-input-confirm>

      <c-button-submit
        :disabled="saveDisabled"
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="$emit('submit', route)"
      />
    </template>
  </b-card>
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'

export default {
  name: 'CRouteEditorInfo',

  i18nOptions: {
    namespaces: [ 'system.apigw' ],
    keyPrefix: 'editor.info',
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
      checkboxLabel: {
        on: this.$t('general:label.general.yes'),
        off: this.$t('general:label.general.no'),
      },
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

    deletedButtonStatusCypressId () {
      return `button-${this.getDeleteStatus.toLowerCase()}`
    },
  },
}
</script>
