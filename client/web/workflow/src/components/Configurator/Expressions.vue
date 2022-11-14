<template>
  <div>
    <b-card
      class="flex-grow-1 rounded-0"
      body-class="p-0"
    >
      <b-card-header
        header-tag="header"
        class="d-flex align-items-center bg-white py-4"
      >
        <h5
          class="d-flex align-items-center mb-0"
        >
          {{ $t('steps:expressions.label') }}
          <a
            :href="documentationURL"
            target="_blank"
            class="d-flex align-items-center h6 mb-0 ml-1"
          >
            <font-awesome-icon
              :icon="['far', 'question-circle']"
            />
          </a>
        </h5>

        <portal to="sidebar-footer">
          <b-button
            variant="primary"
            class="align-top border-0 ml-auto"
            @click="addArgument()"
          >
            {{ $t('steps:expressions.configurator.add-expression') }}
          </b-button>
        </portal>
      </b-card-header>

      <b-card-body
        v-if="hasArguments"
        class="p-0"
      >
        <b-table
          id="arguments"
          fixed
          borderless
          hover
          head-row-variant="secondary"
          details-td-class="bg-white"
          :items="item.config.arguments"
          :fields="argumentFields"
          :tbody-tr-class="rowClass"
          @row-clicked="item=>$set(item, '_showDetails', !item._showDetails)"
        >
          <template #cell(target)="{ item: a }">
            <var>{{ a.target }}</var>
            <samp> ({{ a.type }})</samp>
          </template>

          <template #cell(type)="{ item: a }">
            <var>{{ a.type }}</var>
          </template>

          <template #cell(value)="{ item: a, index }">
            <div
              class="text-truncate"
              :class="{ 'w-75': a._showDetails}"
            >
              <samp>{{ a.expr }}</samp>
            </div>

            <b-button
              v-if="a._showDetails"
              variant="outline-danger"
              class="position-absolute trash border-0"
              @click="removeArgument(index)"
            >
              <font-awesome-icon
                :icon="['far', 'trash-alt']"
              />
            </b-button>
          </template>

          <template #row-details="{ item: a, index }">
            <div class="arrow-up" />

            <b-card
              class="bg-light"
              body-class="px-4 pb-3"
            >
              <b-form-group
                label-class="text-primary"
              >
                <b-form-input
                  v-model="a.target"
                  :placeholder="$t('configurator:target')"
                  @input="$root.$emit('change-detected')"
                />
              </b-form-group>

              <b-form-group
                label-class="text-primary"
                :description="getTypeDescription(a.type)"
              >
                <vue-select
                  v-model="a.type"
                  :options="fieldTypes"
                  :clearable="false"
                  :filter="varFilter"
                  @input="$root.$emit('change-detected')"
                />
              </b-form-group>

              <b-form-group
                class="mb-0"
              >
                <expression-editor
                  :value.sync="a.expr"
                  lang="javascript"
                  show-line-numbers
                  @open="openInEditor(index)"
                  @input="$root.$emit('change-detected')"
                />
              </b-form-group>
            </b-card>
          </template>
        </b-table>
      </b-card-body>
    </b-card>

    <b-modal
      id="expression-editor"
      :visible="!!expressionEditor.currentExpression"
      :title="$t('editor:editor')"
      size="lg"
      :ok-title="$t('general:save')"
      :cancel-title="$t('general:cancel')"
      body-class="p-0"
      @ok="saveExpression"
      @hidden="resetExpression"
    >
      <expression-editor
        :value.sync="currentExpressionValue"
        height="500"
        lang="javascript"
        font-size="18px"
        show-line-numbers
        :border="false"
        :show-popout="false"
      />
    </b-modal>
  </div>
</template>

<script>
import base from './base'
import { VueSelect } from 'vue-select'
import ExpressionEditor from '../ExpressionEditor.vue'
import { stringSearchMaker } from '../../lib/filter'

export default {
  components: {
    ExpressionEditor,
    VueSelect,
  },

  extends: base,

  data () {
    return {
      fieldTypes: [],

      expressionEditor: {
        currentIndex: undefined,
        currentExpression: undefined,
      },
    }
  },

  computed: {
    currentExpressionValue: {
      get () {
        return this.expressionEditor.currentExpression ? this.expressionEditor.currentExpression.expr : ''
      },

      set (value) {
        if (this.expressionEditor.currentExpression) {
          this.expressionEditor.currentExpression.expr = value
        }
      },
    },

    argumentFields () {
      return [
        {
          key: 'target',
          thClass: 'pl-3 py-2',
          tdClass: 'text-truncate pointer',
        },
        {
          key: 'value',
          label: this.$t('steps:expressions.configurator.expression'),
          thClass: 'py-2 pr-3',
          tdClass: 'position-relative pointer',
        },
      ]
    },

    hasArguments () {
      const { config } = this.item || {}
      return (config && (config.arguments || []).length) || []
    },

    documentationURL () {
      // eslint-disable-next-line no-undef
      const [year, month] = VERSION.split('.')
      return `https://docs.cortezaproject.org/corteza-docs/${year}.${month}/integrator-guide/expr/index.html`
    },
  },

  watch: {
    'item.config.stepID': {
      immediate: true,
      handler () {
        this.$set(this.item.config, 'arguments', this.item.config.arguments || [])
      },
    },
  },

  created () {
    this.getTypes()
  },

  methods: {
    varFilter: stringSearchMaker(),

    addArgument () {
      this.item.config.arguments.push({
        target: '',
        expr: '',
        type: 'Any',
        _showDetails: true,
      })
      this.$root.$emit('change-detected')
    },

    removeArgument (index) {
      this.item.config.arguments.splice(index, 1)
      this.$root.$emit('change-detected')
    },

    openInEditor (index = -1) {
      this.expressionEditor = {
        currentIndex: index >= -1 ? index : undefined,
        currentExpression: index >= 0 ? { ...this.item.config.arguments[index] } : undefined,
      }
    },

    saveExpression () {
      if (this.expressionEditor.currentIndex >= 0) {
        const args = [...this.item.config.arguments]
        args[this.expressionEditor.currentIndex] = this.expressionEditor.currentExpression
        this.$set(this.item.config, 'arguments', args)
        this.$root.$emit('change-detected')
      }

      this.resetExpression()
    },

    resetExpression () {
      this.expressionEditor = {
        currentIndex: undefined,
        currentExpression: undefined,
      }
    },

    async getTypes () {
      return this.$AutomationAPI.typeList()
        .then(({ set }) => {
          this.fieldTypes = set
        })
        .catch(this.defaultErrorHandler(this.$t('notification:fetch-types-failed')))
    },

    rowClass (item, type) {
      if (type === 'row') {
        return item._showDetails ? 'border-thick' : 'border-thick-transparent'
      } else if (type === 'row-details') {
        return ''
      }
    },

    getTypeDescription (type) {
      // This will be moved to backend field type information
      const typeDescriptions = {
        ID: 'Make sure to provide the ID in double quotes if you\'re using a literal value. Example "123"',
      }

      return typeDescriptions[type]
    },
  },
}
</script>

<style lang="scss">
tr.b-table-details > td {
  padding-top: 0;
}

.border-thick {
  border-left: 4px solid #A7D0E3;
}

.border-thick-transparent {
  border-left: none;
}

.border-thick-transparent td:first-child {
  padding-left: calc(0.75rem + 2px);
}
</style>

<style lang="scss" scoped>
.trash {
  right: 0;
  left: 1;
  top: 0;
  bottom: 0;
}

.arrow-up {
  width: 0;
  height: 0;
  margin: 0 auto;
  border-left: 10px solid transparent;
  border-right: 10px solid transparent;
  border-bottom: 10px solid $light;
}
</style>
