<template>
  <div>
    <b-card
      data-test-id="card-general-info"
      class="shadow-sm"
      header-bg-variant="white"
      footer-bg-variant="white"
    >
      <template #header>
        <h3 class="m-0">
          {{ $t('general.label') }}
        </h3>
      </template>

      <b-form-group
        :label="$t('general.id')"
        label-cols="2"
      >
        <b-form-input
          v-model="hit.ID"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        :label="$t('general.route')"
        label-cols="2"
      >
        <b-form-input
          v-model="request.route"
          data-test-id="input-route"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        :label="$t('general.URL')"
        label-cols="2"
      >
        <b-form-input
          v-model="request.url"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        :label="$t('general.method')"
        label-cols="2"
      >
        <b-form-input
          v-model="request.method"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        :label="$t('general.statusCode')"
        label-cols="2"
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

      <b-form-group
        :label="$t('general.remoteAddress')"
        label-cols="2"
      >
        <b-form-input
          v-model="request.remoteAddress"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        :label="$t('general.duration')"
        label-cols="2"
      >
        <b-form-input
          v-model="request.duration"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        :label="$t('general.start')"
        label-cols="2"
      >
        <b-form-input
          v-model="request.start"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        :label="$t('general.end')"
        label-cols="2"
      >
        <b-form-input
          v-model="request.end"
          plaintext
          disabled
        />
      </b-form-group>

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
    </b-card>

    <b-card
      header-bg-variant="white"
      footer-bg-variant="white"
      class="shadow-sm mt-3"
    >
      <template #header>
        <h3 class="m-0">
          {{ $t('headers.label') }}
        </h3>
      </template>

      <b-form-group
        v-for="header in request.headers"
        :key="header.label"
        :label="header.label"
        label-cols="2"
      >
        <b-form-input
          v-model="header.value"
          plaintext
          disabled
        />
      </b-form-group>
    </b-card>

    <b-card
      header-bg-variant="white"
      footer-bg-variant="white"
      class="shadow-sm mt-3"
    >
      <template #header>
        <h3 class="m-0">
          {{ $t('body.label') }}
        </h3>
      </template>

      <ace-editor
        :font-size="14"
        width="100%"
        mode="json"
        theme="chrome"
        show-print-margin
        read-only
        show-gutter
        highlight-active-line
        :value="request.body"
        :editor-props="{
          $blockScrolling: false,
        }"
        :set-options="{
          useWorker: false,
        }"
      />
    </b-card>
  </div>
</template>

<script>
import { fmt, NoID } from '@cortezaproject/corteza-js'
import { Ace as AceEditor } from 'vue2-brace-editor'

import 'brace/mode/json'
import 'brace/theme/chrome'

export default {
  name: 'CProfilerHitInfo',

  i18nOptions: {
    namespaces: [ 'system.apigw' ],
    keyPrefix: 'profiler.hit',
  },

  components: {
    AceEditor,
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
