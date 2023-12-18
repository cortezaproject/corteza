
<template>
  <b-card
    data-test-id="card-primary-database"
    class="shadow-sm"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <template>
      <h3 class="d-flex justify-content-between m-0 mb-3">
        {{ $t('title') }}

        <b-button
          v-if="connection"
          data-test-id="button-edit"
          size="sm"
          variant="outline-light"
          class="d-flex align-items-center text-primary border-0"
          :to="{ name: 'system.connection.edit', params: { connectionID: (connection || {}).connectionID } }"
        >
          <font-awesome-icon
            :icon="['far', 'edit']"
          />
        </b-button>
      </h3>
    </template>

    <div
      v-if="connection"
    >
      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('name')"
            label-class="text-primary"
            class="mb-3"
          >
            {{ connection.meta.name || '-' }}
          </b-form-group>
        </b-col>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('handle')"
            label-class="text-primary"
            class="mb-3"
          >
            {{ connection.handle || '-' }}
          </b-form-group>
        </b-col>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group>
            <label
              class="d-flex align-items-center text-primary"
            >
              {{ $t('location') }}
              <c-location
                v-model="locationCoordinates"
                class="ml-1"
              />
            </label>

            {{ locationName || '-' }}
          </b-form-group>
        </b-col>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('ownership')"
            label-class="text-primary"
            class="mb-0"
          >
            {{ connection.meta.ownership || '-' }}
          </b-form-group>
        </b-col>
      </b-row>

      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('sensitivity-level')"
            label-class="text-primary"
          >
            {{ sensitivityLevelName || '-' }}
          </b-form-group>
        </b-col>
      </b-row>
    </div>
  </b-card>
</template>

<script>
import CLocation from 'corteza-webapp-admin/src/components/CLocation'
import { NoID } from '@cortezaproject/corteza-js'

export default {
  components: {
    CLocation,
  },

  i18nOptions: {
    namespaces: 'system.connections',
    keyPrefix: 'primary',
  },

  data () {
    return {
      loading: false,

      connection: undefined,
      sensitivityLevel: undefined,
    }
  },

  computed: {
    locationCoordinates () {
      return this.connection.meta.location.geometry.coordinates || []
    },

    locationName () {
      return this.connection.meta.location.properties.name || 'Unnamed location'
    },

    sensitivityLevelName () {
      const { sensitivityLevelID, handle, meta = {} } = this.sensitivityLevel || {}
      return meta.name || handle || sensitivityLevelID || 'N/A'
    },
  },

  created () {
    this.fetchPrimaryConnection()
  },

  methods: {
    fetchPrimaryConnection () {
      this.loading = true

      return this.$SystemAPI.dalConnectionList({ type: 'corteza::system:primary-dal-connection' }).then(({ set = [] }) => {
        this.connection = set.find(({ type }) => type === 'corteza::system:primary-dal-connection')
        const { sensitivityLevelID } = this.connection.config.privacy || {}

        if (sensitivityLevelID && sensitivityLevelID !== NoID) {
          return this.$SystemAPI.dalSensitivityLevelRead({ sensitivityLevelID })
            .then(sensitivityLevel => {
              this.sensitivityLevel = sensitivityLevel
            })
        }
      }).catch(this.toastErrorHandler(this.$t('notification:fetch.error')))
        .finally(async () => {
          this.loading = false
        })
    },
  },
}
</script>
