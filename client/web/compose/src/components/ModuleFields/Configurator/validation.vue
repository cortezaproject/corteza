<template>
  <div
    v-if="loaded"
  >
    <b-form-group
      :label="$t('sanitizers.label')"
      label-size="lg"
    >
      <template #label>
        <div
          class="d-flex"
        >
          {{ $t('sanitizers.label') }}

          <b-button
            variant="link"
            class="p-0 ml-1 mr-auto"
            @click="field.expressions.sanitizers.push('')"
          >
            {{ $t('sanitizers.add') }}
          </b-button>

          <b-button
            variant="link"
            :href="`${documentationURL}#value-sanitizers`"
            target="_blank"
            class="p-0 ml-1"
          >
            {{ $t('sanitizers.examples') }}
          </b-button>
        </div>
      </template>

      <field-expressions
        v-model="field.expressions.sanitizers"
        :placeholder="$t('sanitizers.expression.placeholder')"
        @remove="field.expressions.sanitizers.splice($event,1)"
      />
      <b-form-text>
        {{ $t('sanitizers.description') }}
      </b-form-text>
    </b-form-group>

    <hr>

    <b-form-group
      class="mt-3"
      label-size="lg"
    >
      <template #label>
        <div
          class="d-flex"
        >
          {{ $t('validators.label') }}

          <b-button
            variant="link"
            class="p-0 ml-1 mr-auto"
            @click="add()"
          >
            {{ $t('sanitizers.add') }}
          </b-button>

          <b-button
            variant="link"
            :href="`${documentationURL}#value-validators`"
            target="_blank"
            class="p-0 ml-1"
          >
            {{ $t('sanitizers.examples') }}
          </b-button>
        </div>
      </template>

      <field-expressions
        v-slot:default="{ value }"
        v-model="field.expressions.validators"
        @remove="field.expressions.validators.splice($event,1)"
      >
        <b-form-input
          v-model="value.test"
          :placeholder="$t('validators.expression.placeholder')"
        />
        <b-input-group-prepend>
          <b-button variant="warning">
            !
          </b-button>
        </b-input-group-prepend>
        <b-form-input
          v-model="value.error"
          :placeholder="$t('validators.error.placeholder')"
        />
        <b-input-group-append>
          <field-translator
            v-if="field"
            :field="field"
            :module="module"
            :highlight-key="`expression.validator.${value.validatorID}.error`"
            :disabled="isNew(value)"
            button-variant="light"
          />
        </b-input-group-append>
      </field-expressions>

      <b-checkbox
        v-model="field.expressions.disableDefaultValidators"
        :disabled="!field.expressions.validators || field.expressions.validators.length === 0"
        :value="true"
        :unchecked-value="false"
        class="mt-2"
      >
        {{ $t('validators.disableBuiltIn') }}
      </b-checkbox>

      <b-form-text>
        {{ $t('validators.description') }}
      </b-form-text>
    </b-form-group>
  </div>
</template>

<script>
import FieldExpressions from 'corteza-webapp-compose/src/components/Common/Module/FieldExpressions'
import FieldTranslator from 'corteza-webapp-compose/src/components/Admin/Module/FieldTranslator'
import { compose, NoID } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  components: {
    FieldExpressions,
    FieldTranslator,
  },

  props: {
    field: {
      type: compose.ModuleField,
      required: true,
    },

    module: {
      type: compose.Module,
      required: true,
    },
  },

  data () {
    return {
      loaded: false,
    }
  },

  computed: {
    documentationURL () {
      // eslint-disable-next-line no-undef
      const [year, month] = VERSION.split('.')
      return `https://docs.cortezaproject.org/corteza-docs/${year}.${month}/integrator-guide/compose-configuration/index.html`
    },
  },

  mounted () {
    if (!this.field.expressions.sanitizers) {
      this.$set(this.field.expressions, 'sanitizers', [])
    }

    if (!this.field.expressions.validators) {
      this.$set(this.field.expressions, 'validators', [])
    }

    if (!this.field.expressions.disableDefaultValidators) {
      this.$set(this.field.expressions, 'disableDefaultValidators', false)
    }

    if (!this.field.expressions.formatters) {
      this.$set(this.field.expressions, 'formatters', [])
    }

    if (!this.field.expressions.disableDefaultFormatters) {
      this.$set(this.field.expressions, 'disableDefaultFormatters', false)
    }

    this.loaded = true
  },

  methods: {
    isNew (value) {
      return !(value.validatorID && value.validatorID !== NoID)
    },

    add () {
      this.field.expressions.validators
        .push({ test: '', error: '' })
    },
  },
}
</script>
