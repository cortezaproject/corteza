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
      style="width: 25%;"
    >
      <b-form-input
        v-model="value.name"
        required
        :readonly="disabled"
        :state="nameState"
        type="text"
        class="form-control"
      />
    </td>
    <td>
      <b-input-group>
        <b-form-input
          v-model="value.label"
          type="text"
          class="form-control"
        />
        <b-input-group-append>
          <field-translator
            :field.sync="value"
            :module="module"
            :disabled="isNew"
            highlight-key="label"
            button-variant="light"
          />
        </b-input-group-append>
      </b-input-group>
    </td>
    <td>
      <b-input-group class="field-type">
        <b-select
          v-model="value.kind"
          :disabled="disabled"
        >
          <option
            v-for="({ kind, label }) in fieldKinds"
            :key="kind"
            :value="kind"
          >
            {{ label }}
          </option>
        </b-select>
        <b-input-group-append>
          <b-button
            variant="light"
            :title="$t('tooltip.field')"
            :disabled="!value.cap.configurable"
            class="px-2"
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
      class="text-right align-middle pr-2"
      style="min-width: 100px;"
    >
      <c-input-confirm
        :no-prompt="!value.name"
        class="mr-2"
        @confirmed="$emit('delete')"
      />
      <c-permissions-button
        v-if="canGrant && exists"
        class="text-dark px-0"
        button-variant="link"
        :title="value.label || value.name || value.fieldID"
        :target="value.label || value.name || value.fieldID"
        :tooltip="$t('permissions:resources.compose.module-field.tooltip')"
        :resource="`corteza::compose:module-field/${module.namespaceID}/${module.moduleID}/${value.fieldID}`"
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
