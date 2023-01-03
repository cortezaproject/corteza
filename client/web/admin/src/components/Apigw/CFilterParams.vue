<template>
  <div>
    <b-form-group
      v-if="filter.params.length"
      class="w-100 mb-0"
    >
      <template
        v-for="(param, index) in filter.params"
      >
        <b-form-group
          :key="index"
        >
          <template slot="label">
            {{ $t(`filters.labels.${param.label}`) }}

            <template v-if="param.label === 'expr'">
              <a
                v-if="param.label === 'expr'"
                :href="documentationURL"
                target="_blank"
              >
                <font-awesome-icon
                  :icon="['far', 'question-circle']"
                />
              </a>
            </template>
          </template>

          <!-- TODO create multi field component-->
          <b-form-checkbox
            v-if="param.type === 'bool'"
            v-model="param.value"
          />

          <vue-select
            v-else-if="param.label === 'workflow'"
            v-model="param.value"
            :options="workflows"
            :reduce="wf => wf.workflowID"
            :placeholder="$t('filters.placeholders.workflow')"
            class="bg-white"
          />

          <b-form-select
            v-else-if="param.label === 'status'"
            v-model="param.value"
            :options="httpStatusOptions"
          >
            <template #first>
              <b-form-select-option
                :value="undefined"
              >
                {{ $t('filters.httpStatus.none') }}
              </b-form-select-option>
            </template>
          </b-form-select>

          <template v-else-if="filter.ref === 'response'">
            <template v-if="param.type === 'input'">
              <b-form-select
                v-model="param.value.type"
                :options="inputTypeOptions"
                class="mb-2"
              />

              <b-input-group>
                <b-input-group-prepend>
                  <b-button variant="dark">
                    ƒ
                  </b-button>
                </b-input-group-prepend>
                <b-form-input
                  v-model="param.value.expr"
                  :placeholder="$t('filters.help.expression.example')"
                />
              </b-input-group>
            </template>

            <template v-else>
              <b-input-group
                v-for="(header, hIndex) in param.value"
                :key="`header-${hIndex}`"
                class="mb-2"
              >
                <b-form-input
                  v-model="header.name"
                  :placeholder="$t('filters.labels.name')"
                />
                <b-form-input
                  v-model="header.expr"
                  :placeholder="$t('filters.labels.value')"
                />

                <b-input-group-append>
                  <b-button
                    variant="danger"
                    @click="param.value.splice(hIndex, 1)"
                  >
                    <font-awesome-icon
                      :icon="['far', 'trash-alt']"
                    />
                  </b-button>
                </b-input-group-append>
              </b-input-group>

              <b-button
                variant="link"
                class="text-decoration-none px-0"
                @click="param.value.push({ name: '', expr: '' })"
              >
                + {{ $t('filters.addHeader') }}
              </b-button>
            </template>
          </template>

          <template v-else>
            <b-form-textarea
              v-if="param.label === 'jsfunc'"
              v-model="param.value"
              max-rows="6"
            />
            <b-input-group v-else>
              <b-input-group-prepend
                v-if="param.label === 'expr'"
              >
                <b-button variant="dark">
                  ƒ
                </b-button>
              </b-input-group-prepend>
              <b-form-input
                v-if="param.label === 'expr'"
                v-model="param.value"
                :placeholder="$t('filters.help.expression.example')"
              />
              <b-form-input
                v-else
                v-model="param.value"
              />
            </b-input-group>
          </template>
        </b-form-group>
      </template>
    </b-form-group>
  </div>
</template>

<script>
import { VueSelect } from 'vue-select'

export default {
  components: {
    VueSelect,
  },

  props: {
    filter: {
      type: Object,
      default: () => ({}),
    },
  },

  data () {
    return {
      workflows: [],

      httpStatusOptions: [
        { value: 300, text: this.$t('filters.httpStatus.300') },
        { value: 301, text: this.$t('filters.httpStatus.301') },
        { value: 302, text: this.$t('filters.httpStatus.302') },
        { value: 303, text: this.$t('filters.httpStatus.303') },
        { value: 304, text: this.$t('filters.httpStatus.304') },
        { value: 307, text: this.$t('filters.httpStatus.307') },
        { value: 308, text: this.$t('filters.httpStatus.308') },
      ],

      inputTypeOptions: [
        'String',
        'Any',
        'Array',
        'KV',
        'DateTime',
        'Float',
        'Integer',
        'Reader',
        'Vars',
      ],
    }
  },

  computed: {
    documentationURL () {
      // eslint-disable-next-line no-undef
      const [year, month] = VERSION.split('.')
      return `https://docs.cortezaproject.org/corteza-docs/${year}.${month}/integrator-guide/expr/index.html`
    },
  },

  created () {
    if (this.filter.params.some(({ label = '' }) => label === 'workflow')) {
      this.$AutomationAPI.workflowList()
        .then(({ set: workflows = [] }) => {
          this.workflows = workflows.map(({ workflowID, handle, meta }) => {
            return { label: meta.name || handle, workflowID }
          })
        })
    }
  },
}
</script>
