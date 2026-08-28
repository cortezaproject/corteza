<template>
  <div>
    <b-form-checkbox
      v-model="f.options.dynamicMode"
      v-b-tooltip.noninteractive.hover="{ title: $t('kind.link.dynamicModeTooltip'), boundary: 'viewport' }"
    >
      {{ $t('kind.link.dynamicMode') }}
    </b-form-checkbox>

    <div
      v-if="!f.options.dynamicMode"
      class="ml-4 mb-3"
    >
      <div>
        <label>{{ $t('kind.link.protocol') }}</label>
        <b-form-select
          v-model="f.options.tempProtocol"
          :options="protocolOptions"
          @change="onProtocolChange"
        />
      </div>

      <div
        v-if="f.options.tempProtocol === '__CUSTOM__'"
        class="mt-2"
      >
        <div>
          <label>{{ $t('kind.link.customProtocol') }}</label>
          <b-form-input
            v-model="f.options.customProtocol"
            placeholder="e.g., skype:"
          />
        </div>
      </div>
    </div>

    <div v-if="showURLOptions">
      <div>
        <b-form-checkbox v-model="f.options.trimPath">
          {{ $t('kind.url.label') }}
        </b-form-checkbox>
      </div>
      <div>
        <b-form-checkbox v-model="f.options.trimFragment">
          {{ $t('kind.url.trimHash') }}
        </b-form-checkbox>
      </div>
      <div>
        <b-form-checkbox v-model="f.options.trimQuery">
          {{ $t('kind.url.trimQuestionMark') }}
        </b-form-checkbox>
      </div>
      <div>
        <b-form-checkbox v-model="f.options.onlySecure">
          {{ $t('kind.url.sshOnly') }}
        </b-form-checkbox>
      </div>
    </div>

    <div>
      <b-form-checkbox v-model="f.options.outputPlain">
        {{ $t('kind.link.preventToLink') }}
      </b-form-checkbox>
    </div>
  </div>
</template>

<script>
import base from './base'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  extends: base,

  data () {
    return {
      tempProtocol: '',
    }
  },

  computed: {
    isStaticURL () {
      return !this.f.options.dynamicMode && this.f.options.tempProtocol === 'https://'
    },

    showURLOptions () {
      return this.f.options.dynamicMode || this.isStaticURL
    },

    protocolOptions () {
      return [
        { text: this.$t('kind.link.protocolOptions.none'), value: '' },
        { text: this.$t('kind.link.protocolOptions.mailto'), value: 'mailto:' },
        { text: this.$t('kind.link.protocolOptions.tel'), value: 'tel:' },
        { text: this.$t('kind.link.protocolOptions.url'), value: 'https://' },
        { text: this.$t('kind.link.protocolOptions.skype'), value: 'skype:' },
        { text: this.$t('kind.link.protocolOptions.msteams'), value: 'msteams:' },
        { text: this.$t('kind.link.protocolOptions.slack'), value: 'slack://' },
        { text: this.$t('kind.link.protocolOptions.sms'), value: 'sms:' },
        { text: this.$t('kind.link.protocolOptions.facetime'), value: 'facetime:' },
        { text: this.$t('kind.link.protocolOptions.zoommtg'), value: 'zoommtg://' },
        { text: this.$t('kind.link.protocolOptions.whatsapp'), value: 'whatsapp://' },
        { text: this.$t('kind.link.protocolOptions.signal'), value: 'sgnl:' },
        { text: this.$t('kind.link.protocolOptions.fb'), value: 'fb://' },
        { text: this.$t('kind.link.protocolOptions.fbmessenger'), value: 'fb-messenger://' },
        { text: this.$t('kind.link.protocolOptions.instagram'), value: 'instagram://' },
        { text: this.$t('kind.link.protocolOptions.twitter'), value: 'twitter://' },
        { text: this.$t('kind.link.protocolOptions.youtube'), value: 'youtube://' },
        { text: this.$t('kind.link.protocolOptions.spotify'), value: 'spotify://' },
        { text: this.$t('kind.link.protocolOptions.tiktok'), value: 'tiktok://' },
        { text: this.$t('kind.link.protocolOptions.discord'), value: 'discord://' },
        { text: this.$t('kind.link.protocolOptions.custom'), value: '__CUSTOM__' },
      ]
    },
  },

  watch: {
    'f.options.customProtocol': {
      handler (val) {
        if (val && val !== '__CUSTOM__') {
          this.tempProtocol = val
        }
      },
    },

    'f.options.dynamicMode': {
      handler (val) {
        if (!val && this.f.options.customProtocol) {
          this.tempProtocol = this.f.options.customProtocol
        }
      },
    },
  },

  mounted () {
    if (!this.f.options.tempProtocol) {
      if (this.f.options.customProtocol) {
        this.f.options.tempProtocol = this.f.options.customProtocol
      } else {
        this.f.options.tempProtocol = ''
      }
    } else {
      this.tempProtocol = this.f.options.tempProtocol
    }
  },

  methods: {
    onProtocolChange () {
      if (this.tempProtocol !== '__CUSTOM__') {
        this.$set(this.field.options, 'customProtocol', this.tempProtocol)
      }
    },
  },
}
</script>
