<template>
  <b-tabs
    nav-wrapper-class="bg-white border-bottom"
    active-tab-class="tab-content h-auto overflow-auto"
    card
    lazy
  >
    <b-tab
      active
      :title="$t('label.general')"
    >
      <basic
        :namespace="namespace"
        :module="module"
        :field.sync="f"
      />
    </b-tab>

    <b-tab
      v-if="fieldComponent"
      :title="$t(`fieldKinds.${field.kind}.label`)"
    >
      <component
        :is="fieldComponent"
        :namespace="namespace"
        :module="module"
        :field.sync="f"
      />
    </b-tab>

    <b-tab
      v-if="field.cap.multi"
      :disabled="!field.isMulti"
      :title="$t('label.multi')"
    >
      <multi
        :namespace="namespace"
        :field.sync="f"
      />
    </b-tab>

    <b-tab
      :title="$t('label.validation')"
    >
      <validation
        :namespace="namespace"
        :module="module"
        :field.sync="f"
      />
    </b-tab>

    <b-tab
      :title="$t('label.privacy')"
    >
      <data-privacy-settings
        v-if="connection"
        :resource="field"
        :connection="connection"
        :sensitivity-levels="sensitivityLevels"
        :max-level="maxLevelID"
        :translations="{
          sensitivity: {
            label: $t('field:privacy.sensitivity-level.label'),
            placeholder: $t('field:privacy.sensitivity-level.placeholder'),
          },
          usage: {
            label: $t('field:privacy.usage-disclosure.label'),
          },
        }"
      />
    </b-tab>
  </b-tabs>
</template>
<script>
import { NoID } from '@cortezaproject/corteza-js'
import base from './base'
import * as Configurators from './loader'
import multi from './multi'
import basic from './basic'
import validation from './validation'
import DataPrivacySettings from 'corteza-webapp-compose/src/components/Admin/Module/DataPrivacySettings'

export default {
  i18nOptions: {
    namespaces: 'general',
  },

  components: {
    ...Configurators,
    multi,
    basic,
    validation,
    DataPrivacySettings,
  },

  extends: base,

  props: {
    connection: {
      type: Object,
      required: true,
    },

    sensitivityLevels: {
      type: Array,
      default: () => [],
    },
  },

  computed: {
    fieldComponent () {
      // If field doesn't have a configurator, we show no field tab
      return Configurators[this.field.kind]
    },

    maxLevelID () {
      const { sensitivityLevelID = NoID } = this.module.config.privacy || {}
      return sensitivityLevelID
    },
  },
}
</script>

<style lang="scss" scoped>
.tab-content {
  max-height: 70vh;
}
</style>
