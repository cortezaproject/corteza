<template>
  <div>
    <b-card
      data-test-id="card-general-info"
      header-class="border-bottom"
      class="shadow-sm"
    >
      <template #header>
        <h4 class="m-0">
          {{ $t('general.label') }}
        </h4>
      </template>

      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('general.id')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="hit.ID"
              plaintext
              disabled
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('general.route')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="request.route"
              data-test-id="input-route"
              plaintext
              disabled
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('general.URL')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="request.url"
              plaintext
              disabled
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('general.method')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="request.method"
              plaintext
              disabled
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('general.statusCode')"
            label-class="text-primary"
          >
            <div
              class="d-flex align-items-center h-100"
            >
              <h5 class="mb-0">
                <b-badge :variant="getStatusCodeVariant(request.statusCode)">
                  {{ request.statusCode }}
                </b-badge>
              </h5>
            </div>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('general.remoteAddress')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="request.remoteAddress"
              plaintext
              disabled
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('general.duration')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="request.duration"
              plaintext
              disabled
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('general.start')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="request.start"
              plaintext
              disabled
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('general.end')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="request.end"
              plaintext
              disabled
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            v-if="showOpenRoute"
            label-cols="2"
          >
            <b-button
              variant="light"
              :to="{ name: 'system.apigw.edit', params: { routeID: request.routeID } }"
            >
              {{ $t('general.openRoute') }}
            </b-button>
          </b-form-group>
        </b-col>
      </b-row>
    </b-card>

    <b-card
      header-class="border-bottom"
      class="shadow-sm mt-3"
    >
      <template #header>
        <h4 class="m-0">
          {{ $t('headers.label') }}
        </h4>
      </template>

      <b-row>
        <b-col
          v-for="header in request.headers"
          :key="header.label"
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="header.label"
            label-class="text-primary"
          >
            <b-form-input
              v-model="header.value"
              plaintext
              disabled
            />
          </b-form-group>
        </b-col>
      </b-row>
    </b-card>

    <b-card
      header-class="border-bottom"
      body-class="p-0"
      class="shadow-sm mt-3 overflow-hidden"
    >
      <template #header>
        <h4 class="m-0">
          {{ $t('body.label') }}
        </h4>
      </template>

      <c-ace-editor
        :value="request.body"
        :border="false"
        lang="json"
        height="400px"
        show-line-numbers
        read-only
      />
    </b-card>
  </div>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'
import { fmt, NoID } from '@cortezaproject/corteza-js'

const { CAceEditor } = components

export default {
  name: 'CProfilerHitInfo',

  i18nOptions: {
    namespaces: [ 'system.apigw' ],
    keyPrefix: 'profiler.hit',
  },

  components: {
    CAceEditor,
  },

  props: {
    hit: {
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

  computed: {
    request () {
      // eslint-disable-next-line camelcase
      const { request = {}, body = '', route = NoID, time_duration = 0, time_start, time_finish, http_status_code } = this.hit || {}
      const { URL = {}, RequestURI, Method, RemoteAddr, Header = {} } = request
      const { Path } = URL
      const headers = Object.entries(Header).map(([key, value = []]) => {
        return { label: key, value: value.join('') }
      })

      // Try to parse body as json
      let jsonBody = atob(body)
      try {
        jsonBody = JSON.stringify(JSON.parse(jsonBody), null, 2)
      } catch (e) {}

      return {
        routeID: route,
        route: Path,
        url: RequestURI,
        method: Method,
        statusCode: http_status_code,
        remoteAddress: RemoteAddr,
        duration: `${time_duration.toFixed(2)} ms`,
        start: fmt.fullDateTime(time_start),
        end: fmt.fullDateTime(time_finish),
        headers,
        body: jsonBody,
      }
    },

    showOpenRoute () {
      return this.request.routeID !== NoID
    },
  },

  methods: {
    getStatusCodeVariant (statusCode = '') {
      const codeVariants = {
        '2': 'success',
        '3': 'info',
        '4': 'danger',
        '5': 'warning',
      }

      return codeVariants[statusCode[0]]
    },
  },
}
</script>
