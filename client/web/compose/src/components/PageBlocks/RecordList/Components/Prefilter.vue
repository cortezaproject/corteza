<template>
  <c-form-table-wrapper hide-add-button>
    <b-form-group
      :label="$t('recordList.record.prefilterCommand')"
      label-class="text-primary"
      class="m-0"
    >
      <b-row v-if="textInput">
        <b-col>
          <b-form-group label-class="text-primary">
            <c-filter-field
              v-model="options.prefilter"
              auto-complete
              height="59px"
              lang="javascript"
              :suggestion-tree="filterSuggestionTree"
            />

            <i18next
              path="recordList.record.prefilterFootnote"
              tag="small"
              class="text-muted"
            >
              <code>${record.values.fieldName}</code>
              <code>${recordID}</code>
              <code>${ownerID}</code>
              <span><code>${userID}</code>, <code>${user.name}</code></span>
            </i18next>
          </b-form-group>
        </b-col>
      </b-row>

      <c-form-table-wrapper
        v-else
        hide-add-button
      >
        <filter-toolbox
          v-model="filterGroup"
          :module="module"
          :namespace="namespace"
          :mock.sync="mock"
          reset-filter-on-created
        />
      </c-form-table-wrapper>

      <div class="mt-1 d-flex align-items-center">
        <b-button
          variant="link"
          size="sm"
          class="ml-auto text-decoration-none"
          @click="toggleFilterView"
        >
          {{ $t('recordList.prefilter.toggleInputType') }}
        </b-button>
      </div>
    </b-form-group>
  </c-form-table-wrapper>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'
import { compose, validator } from '@cortezaproject/corteza-js'
import {
  getRecordListFilterSql,
  trimChar,
} from 'corteza-webapp-compose/src/lib/record-filter.js'
import FilterToolbox from 'corteza-webapp-compose/src/components/Common/FilterToolbox.vue'

const { CFilterField } = components

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'RecordListConfiguratorPrefilter',

  components: {
    FilterToolbox,
    CFilterField,
  },

  props: {
    options: {
      type: Object,
      required: true,
    },

    namespace: {
      type: compose.Namespace,
      required: true,
    },

    module: {
      type: compose.Module,
      required: true,
    },
  },

  data () {
    return {
      textInput: true,
      filterGroup: [],
    }
  },

  computed: {
    filterSuggestionTree () {
      const moduleFields = (this.module.fields || []).map(({ name }) => name)
      const userProperties = Object.keys(this.$auth.user)

      return {
        '': ['record', 'recordID', 'ownerID', 'userID', 'user'].concat(moduleFields),
        'record': ['values'],
        'record.values': moduleFields,
        'user': userProperties,
      }
    },
  },

  created () {
    // Change all module fields to single value to keep multi value fields and single value
    const module = JSON.parse(JSON.stringify(this.module || {}))

    module.fields = [
      ...[...module.fields].map((f) => {
        f.multi = f.isMulti
        f.isMulti = false

        // Disable edge case options
        if (f.kind === 'DateTime') {
          f.options.onlyFutureValues = false
          f.options.onlyPastValues = false
        }

        return f
      }),
      ...this.module.systemFields().map((sf) => {
        return { ...sf, label: this.$t(`field:system.${sf.name}`) }
      }),
    ]

    this.mock = {
      namespace: this.namespace,
      module,
      errors: new validator.Validated(),
    }
  },

  methods: {
    toggleFilterView () {
      if (!this.textInput) {
        this.options.prefilter = this.parseFilter()
      }

      this.textInput = !this.textInput
    },

    getOptionKey ({ name }) {
      return name
    },

    processFilter (filterGroup = this.value) {
      return filterGroup.map(({ groupCondition, filter = [], name }) => {
        filter = filter.map(({ record, ...f }) => {
          if (record) {
            f.value = record[f.name] || record.values[f.name]
          }

          if (this.isBetweenOperator(f.operator)) {
            f.value = {
              start: this.getField(f.name).isSystem
                ? record[`${f.name}-start`]
                : record.values[`${f.name}-start`],
              end: this.getField(f.name).isSystem
                ? record[`${f.name}-end`]
                : record.values[`${f.name}-end`],
            }
          }

          return f
        })

        return { groupCondition, filter, name }
      })
    },

    isBetweenOperator (op) {
      return ['BETWEEN', 'NOT BETWEEN'].includes(op)
    },

    parseFilter (filterGroup = this.filterGroup) {
      const filter = this.processFilter(filterGroup)

      const filterSqlArray = filter
        .map(({ groupCondition, filter = [] }) => {
          groupCondition = groupCondition ? ` ${groupCondition} ` : ''
          filter = getRecordListFilterSql(filter)

          return filter ? `${filter}${groupCondition}` : ''
        })
        .filter((filter) => filter)

      const filterSql = trimChar(
        trimChar(filterSqlArray.join(''), ' AND '),
        ' OR '
      )

      return filterSql
    },
  },
}
</script>
