<template>
  <b-card
    class="shadow-sm"
    :title="$t('title')"
  >
    <b-container
      v-if="canManage"
    >
      <b-row
        v-if="issues.length"
      >
        <b-col
          cols="12"
        >
          <label>{{ $t('connectivity-issues') }}</label>
          <b-alert
            v-for="issue in issues"
            :key="issue"
            show
            variant="danger"
          >
            {{ issue }}
          </b-alert>
        </b-col>
      </b-row>
      <b-row>
        <b-col
          cols="12"
          lg="12"
        >
          <b-form-group
            :label="$t('form.model-ident.label')"
            :description="$t('form.model-ident.description', { interpolation: { prefix: '{{{', suffix: '}}}' } })"
          >
            <b-form-input
              v-model="dal.modelIdent"
              :disabled="disabled"
              :placeholder="$t('form.model-ident.placeholder')"
            />
          </b-form-group>
        </b-col>
      </b-row>
      <b-row>
        <b-col
          cols="12"
          lg="12"
        >
          <b-form-group
            :label="$t('form.type.label')"
            :description="$t('form.type.description')"
          >
            <b-form-input
              v-model="dal.type"
              :disabled="disabled"
              :placeholder="$t('form.type.placeholder')"
            />
          </b-form-group>
        </b-col>
      </b-row>
      <b-row>
        <b-col
          cols="12"
          lg="12"
        >
          <b-form-group
            :label="$t('form.params.label')"
            :description="$t('form.params.description')"
          >
            <b-form-textarea
              v-model="paramsJson"
              :disabled="disabled"
              :class="paramsJsonEditorClass"
              :placeholder="$t('form.params.placeholder')"
              rows="5"
              @blur="processParamsJSON"
            />
          </b-form-group>
        </b-col>
      </b-row>
    </b-container>
    <b-container
      v-else
    >
      <b-alert
        variant="warning"
        show
      >
        {{ $t('no-access-warning') }}
      </b-alert>
    </b-container>
  </b-card>
</template>

<script>

export default {
  i18nOptions: {
    namespaces: 'system.connections',
    keyPrefix: 'editor.dal',
  },

  props: {
    disabled: { type: Boolean, default: false },
    canManage: { type: Boolean, default: false },

    dal: {
      type: Object,
      required: true,
    },

    issues: {
      type: Array,
      default: () => ([]),
    },
  },

  data () {
    return {
      /**
       * JSON string version of the connection params object
       * used to display nicely formatted JSON in the textarea
       */
      paramsJson: '',

      // holds JSON validation errors
      paramsJsonEditorClass: '',
    }
  },

  watch: {
    'dal': {
      handler: function (dal) {
        this.paramsJson = JSON.stringify(dal.params || { dsn: '' }, null, 2)
      },
      deep: true,
      immediate: true,
    },
  },

  methods: {
    /**
     * Validates JSON string and sets errors if any
     * In case of a valid JSON it will update the connection object
     */
    processParamsJSON () {
      this.paramsJsonEditorClass = ''

      try {
        // parse json string, ensure it's an object
        const json = JSON.parse(this.paramsJson)
        if (typeof json !== 'object') {
          throw new Error('JSON is not an object')
        }

        // iterate through all properties and assign them to connection object
        if (!this.dal.params) {
          this.$set(this.dal, 'params', {})
        }
        for (const key in json) {
          if (json.hasOwnProperty(key)) {
            this.$set(this.dal.params, key, json[key])
          }
        }
      } catch (e) {
        this.paramsJsonEditorClass = 'border-danger'
      }
    },
  },
}
</script>
