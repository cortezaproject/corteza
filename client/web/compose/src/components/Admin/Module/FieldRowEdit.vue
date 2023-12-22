<template>
  <tr>
    <td
      v-b-tooltip.hover
      class="handle align-middle"
    >
      <font-awesome-icon
        :icon="['fas', 'bars']"
        class="text-light grab"
      />
    </td>

    <td
      style="min-width: 250px;"
    >
      <b-form-input
        v-model="value.name"
        required
        :readonly="disabled"
        :state="nameState"
        type="text"
      />
    </td>

    <td
      style="min-width: 250px;"
    >
      <b-input-group>
        <b-form-input
          v-model="value.label"
          type="text"
        />

        <b-input-group-append>
          <field-translator
            :field.sync="value"
            :module="module"
            :disabled="isNew"
            highlight-key="label"
          />
        </b-input-group-append>
      </b-input-group>
    </td>

    <td
      style="min-width: 250px;"
    >
      <b-input-group class="field-type">
        <c-input-select
          v-model="value.kind"
          :options="fieldKinds"
          :reduce="kind => kind.kind"
          :disabled="disabled"
          :clearable="false"
          @input="$emit('updateKind')"
        />

        <b-input-group-append>
          <b-button
            v-b-tooltip.hover="{ title: $t('tooltip.field'), container: '#body' }"
            data-test-id="button-configure-field"
            variant="light"
            :disabled="!value.cap.configurable"
            @click.prevent="$emit('edit')"
          >
            <font-awesome-icon
              :icon="['fas', 'wrench']"
            />
          </b-button>
        </b-input-group-append>
      </b-input-group>
    </td>

    <td />
    <td />

    <td
      class="align-middle text-center"
    >
      <b-form-checkbox
        v-model="value.isRequired"
        :disabled="!value.cap.required"
        switch
      />
    </td>

    <td
      class="align-middle text-center"
    >
      <b-form-checkbox
        v-model="value.isMulti"
        :disabled="!value.cap.multi"
        switch
        class="ml-2"
      />
    </td>

    <td
      class="text-right align-middle pr-2"
      style="min-width: 100px;"
    >
      <c-permissions-button
        v-if="canGrant && exists"
        button-variant="outline-light"
        size="sm"
        :title="value.label || value.name || value.fieldID"
        :target="value.label || value.name || value.fieldID"
        :tooltip="$t('permissions:resources.compose.module-field.tooltip')"
        :resource="`corteza::compose:module-field/${module.namespaceID}/${module.moduleID}/${value.fieldID}`"
        class="text-dark border-0 mr-2"
      />

      <c-input-confirm
        show-icon
        @confirmed="$emit('delete')"
      />
    </td>
  </tr>
</template>

<script>
import FieldTranslator from 'corteza-webapp-compose/src/components/Admin/Module/FieldTranslator'
import { compose, NoID } from '@cortezaproject/corteza-js'

export default {
  components: {
    FieldTranslator,
  },

  i18nOptions: {
    namespaces: 'general',
  },

  props: {
    value: {
      type: Object,
      required: true,
    },

    module: {
      type: compose.Module,
      required: true,
    },

    canGrant: {
      type: Boolean,
      required: false,
    },

    hasRecords: {
      type: Boolean,
      required: true,
    },

    isDuplicate: {
      type: Boolean,
      required: false,
    },
  },

  data () {
    return {
      updateField: null,
    }
  },

  computed: {
    nameState () {
      if (this.disabled) {
        return null
      }

      if (this.isDuplicate) {
        return false
      }

      return this.value.isValid ? null : false
    },

    disabled () {
      return this.value.fieldID !== NoID && this.hasRecords
    },

    isNew () {
      return this.module.moduleID === NoID || this.value.fieldID === NoID
    },

    fieldKinds () {
      return [...compose.ModuleFieldRegistry.keys()]
        // for now this field is hidden, since it's implementation is mia.
        .map(kind => {
          return { kind, label: this.$t('fieldKinds.' + kind + '.label') }
        }).sort((a, b) => a.label.localeCompare(b.text))
    },

    exists () {
      return this.module.ID !== NoID && this.value.fieldID !== NoID
    },
  },
}
</script>
<style lang="scss" scoped>
td {
  input, .input-group {
    min-width: 150px;
  }

  .handle {
    width: 30px;
  }
}
</style>
