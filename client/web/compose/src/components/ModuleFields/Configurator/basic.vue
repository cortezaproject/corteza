<template>
  <div>
    <b-form-checkbox
      v-model="field.isRequired"
      :disabled="!field.cap.required || showValueExpr"
    >
      {{ $t('general:label.required') }}
    </b-form-checkbox>

    <b-form-checkbox
      v-model="field.isMulti"
      :disabled="!field.cap.multi"
    >
      {{ $t('label.multi') }}
    </b-form-checkbox>

    <b-form-checkbox
      v-model="showValueExpr"
      :disabled="field.isRequired || defaultValueEnabled"
    >
      {{ $t('valueExpr.label') }}
    </b-form-checkbox>

    <b-form-checkbox
      v-if="showDefaultValue"
      :checked="defaultValueEnabled"
      :disabled="!!showValueExpr"
      @change="toggleDefaultValue()"
    >
      {{ $t('defaultValue') }}
    </b-form-checkbox>

    <hr>

    <b-form-group
      v-if="showValueExpr"
      :label="$t('valueExpr.label')"
      class="mt-2"
    >
      <b-input-group>
        <b-input-group-append>
          <b-button variant="dark">
            ƒ
          </b-button>
        </b-input-group-append>
        <b-form-input
          v-model="field.expressions.value"
          :placeholder="$t('valueExpr.placeholder')"
        />
        <b-input-group-append>
          <b-button
            variant="outline-secondary"
            :href="documentationURL"
            class="d-flex justify-content-center align-items-center"
            target="_blank"
          >
            ?
          </b-button>
        </b-input-group-append>
      </b-input-group>
      <b-form-text>
        {{ $t('valueExpr.description') }}
      </b-form-text>
    </b-form-group>

    <b-form-group
      v-else-if="showDefaultField"
      :label="$t('defaultFieldValue')"
      class="mt-3 mb-0"
    >
      <field-editor
        value-only
        v-bind="mock"
      />
    </b-form-group>

    <hr v-if="showValueExpr || showDefaultField">

    <b-form-group
      :label="$t(`options.description.label.${noDescriptionEdit ? 'default' : 'view'}`)"
      class="mt-2"
    >
      <b-input-group>
        <b-form-input
          v-model="field.options.description.view"
          :placeholder="$t(`options.description.placeholder.${noDescriptionEdit ? 'default' : 'view'}`)"
        />
        <b-input-group-append>
          <field-translator
            v-if="field"
            :field="field"
            :module="module"
            :disabled="isNew"
            :highlight-key="`meta.description.view`"
            button-variant="light"
          />
        </b-input-group-append>
      </b-input-group>
    </b-form-group>

    <b-form-group
      v-if="!noDescriptionEdit"
      :label="$t('options.description.label.edit')"
      class="mt-2"
    >
      <b-input-group>
        <b-form-input
          v-model="field.options.description.edit"
          :placeholder="$t('options.description.placeholder.edit')"
        />
        <b-input-group-append>
          <field-translator
            v-if="field"
            :field="field"
            :module="module"
            :disabled="isNew"
            :highlight-key="`meta.description.edit`"
            button-variant="light"
          />
        </b-input-group-append>
      </b-input-group>
    </b-form-group>

    <b-form-checkbox
      :checked="noDescriptionEdit"
      tabindex="-1"
      @change="$set(field.options.description, 'edit', $event ? undefined : '')"
    >
      {{ $t('options.description.same') }}
    </b-form-checkbox>

    <hr>

    <b-form-group
      :label="$t(`options.hint.label.${noHintEdit ? 'default' : 'view'}`)"
      class="mt-2"
    >
      <b-input-group>
        <b-form-input
          v-model="field.options.hint.view"
          :placeholder="$t(`options.hint.placeholder.${noHintEdit ? 'default' : 'view'}`)"
        />
        <b-input-group-append>
          <field-translator
            v-if="field"
            :field="field"
            :module="module"
            :disabled="isNew"
            :highlight-key="`meta.hint.view`"
            button-variant="light"
          />
        </b-input-group-append>
      </b-input-group>
    </b-form-group>

    <b-form-group
      v-if="!noHintEdit"
      :label="$t('options.hint.label.edit')"
      class="mt-2"
    >
      <b-input-group>
        <b-form-input
          v-model="field.options.hint.edit"
          :placeholder="$t('options.hint.placeholder.edit')"
        />
        <b-input-group-append>
          <field-translator
            v-if="field"
            :field="field"
            :module="module"
            :disabled="isNew"
            :highlight-key="`meta.hint.edit`"
            button-variant="light"
          />
        </b-input-group-append>
      </b-input-group>
    </b-form-group>

    <b-form-checkbox
      :checked="noHintEdit"
      tabindex="-1"
      @change="$set(field.options.hint, 'edit', $event ? undefined : '')"
    >
      {{ $t('options.hint.same') }}
    </b-form-checkbox>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { compose, validator, NoID } from '@cortezaproject/corteza-js'
import FieldEditor from '../Editor'
import FieldTranslator from 'corteza-webapp-compose/src/components/Admin/Module/FieldTranslator'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  components: {
    FieldEditor,
    FieldTranslator,
  },

  props: {
    namespace: {
      type: compose.Namespace,
      required: true,
    },

    module: {
      type: compose.Module,
      required: true,
    },

    field: {
      type: compose.ModuleField,
      required: true,
    },
  },

  data () {
    return {
      showValueExpr: false,

      mock: {
        show: true,
        namespace: undefined,
        module: undefined,
        field: undefined,
        record: undefined,
        errors: new validator.Validated(),
      },
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
    }),

    noDescriptionEdit () {
      return this.field.options.description.edit === undefined
    },

    noHintEdit () {
      return this.field.options.hint.edit === undefined
    },

    showDefaultValue () {
      return !['File'].includes(this.field.kind)
    },

    defaultValueEnabled () {
      return !!this.field.defaultValue.length
    },

    showDefaultField () {
      return this.defaultValueEnabled && this.mock.show && this.mock.field
    },

    documentationURL () {
      // eslint-disable-next-line no-undef
      const [year, month] = VERSION.split('.')
      return `https://docs.cortezaproject.org/corteza-docs/${year}.${month}/integrator-guide/expr/index.html`
    },

    isNew () {
      return this.field.fieldID === NoID
    },
  },

  watch: {
    defaultValueEnabled: {
      handler (val) {
        if (val) {
          this.showValueExpr = false
        }
      },
    },

    'mock.record.values': {
      handler ({ defValField: dv }) {
        if (!Array.isArray(dv)) {
          dv = [dv]
        } else if (!dv.length) {
          dv = [undefined]
        }

        // Transform to backend value struct
        this.field.defaultValue = dv.map(v => {
          if (v !== undefined && v.toString) {
            v = v.toString()
          }

          const defaultValue = {
            name: this.field.name,
          }

          if (v) {
            defaultValue.value = v
          }

          return defaultValue
        })
      },
      deep: true,
    },

    'field.options': {
      handler (options) {
        if (this.mock.field) {
          this.mock.field.options = options
        }
      },
      deep: true,
    },

    'field.isMulti': {
      handler () {
        // Only init mocks if default value exists, to convert to multi mock
        if (this.field.defaultValue.length) {
          this.initMocks(this.field.defaultValue)
        }
      },
    },
  },

  /**
   * Prepare mock values for default-value field editor
   */
  created () {
    let { defaultValue, expressions } = this.field

    if (!defaultValue) {
      defaultValue = []
    }

    if (defaultValue.length) {
      this.initMocks(defaultValue)
    }

    // when loading, assume empty strings are same as undefined
    if (!this.field.options.hint.edit) {
      this.$set(this.field.options.hint, 'edit', undefined)
    }

    if (!this.field.options.description.edit) {
      this.$set(this.field.options.description, 'edit', undefined)
    }

    this.showValueExpr = expressions.value && expressions.value.length > 0
    if (!this.field.expressions.value) {
      this.$set(this.field.expressions, 'value', '')
    }
  },

  beforeDestroy () {
    // Sanitize expression/required flag
    if (this.showValueExpr) {
      this.field.required = false
      this.field.defaultValue = []
    } else {
      this.field.expressions.value = undefined
    }
  },

  methods: {
    initMocks (defaultValue = []) {
      if (this.field.isMulti) {
        defaultValue = defaultValue.map(v => (v || {}).value).filter(v => v)
      } else {
        defaultValue = (defaultValue[0] || {}).value
      }

      this.mock.namespace = this.namespace
      this.mock.field = compose.ModuleFieldMaker(this.field)
      this.mock.field.apply({ label: this.mock.field.label || 'Default value' })
      this.mock.field.apply({ name: 'defValField' })
      this.mock.module = new compose.Module({ fields: [this.mock.field] }, this.namespace)
      this.mock.record = new compose.Record(this.mock.module, { defValField: defaultValue })
    },

    toggleDefaultValue () {
      if (this.defaultValueEnabled) {
        this.field.defaultValue = []
      } else {
        this.initMocks()
      }
    },
  },
}
</script>
