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

    <b-form-row
      @submit.prevent="$emit('submit', settings)"
    >
      <b-col
        cols="12"
        md="6"
        class="mb-3 mb-md-0"
      >
        <b-form-group
          :label="$t('profiler.label')"
          label-class="text-primary"
          class="mb-0"
        >
          <b-form-radio-group
            v-model="profilerSetting"
            :options="profilerOptions"
            button-variant="outline-primary"
            size="sm"
            buttons
          />
        </b-form-group>
      </b-col>

      <b-col
        cols="12"
        md="6"
      >
        <b-form-group
          :label="$t('proxy.label')"
          label-class="text-primary"
          class="mb-0"
        >
          <b-form-checkbox
            v-model="settings['apigw.proxy.follow-redirects']"
          >
            {{ $t('proxy.follow') }}
          </b-form-checkbox>
        </b-form-group>
      </b-col>
    </b-form-row>

    <template #footer>
      <c-button-submit
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="float-right"
        @submit="$emit('submit', settings)"
      />
    </template>
  </b-card>
</template>

<script>

export default {
  name: 'CSettingsEditor',

  i18nOptions: {
    namespaces: 'system.apigw',
    keyPrefix: 'settings',
  },

  props: {
    settings: {
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
  },

  data () {
    return {
      profilerOptions: [
        { value: 'disabled', text: 'Disabled' },
        { value: 'filter', text: 'Enabled as filter' },
        { value: 'global', text: 'Enabled for all routes' },
      ],
    }
  },

  computed: {
    profilerSetting: {
      get () {
        if (this.settings['apigw.profiler.enabled']) {
          return this.settings['apigw.profiler.global'] ? 'global' : 'filter'
        }

        return 'disabled'
      },

      set (setting) {
        this.settings['apigw.profiler.enabled'] = ['filter', 'global'].includes(setting)
        this.settings['apigw.profiler.global'] = setting === 'global'
      },
    },
  },
}
</script>
