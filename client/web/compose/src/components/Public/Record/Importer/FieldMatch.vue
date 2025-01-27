<template>
  <div class="d-flex flex-column h-100 overflow-hidden">
    <b-form-group
      :label="$t('recordList.import.matchFields')"
      label-class="text-primary p-3"
      class="flex-fill overflow-auto mb-0"
    >
      <div
        v-if="hasRequiredFileFields"
        class="small px-3 mb-3"
      >
        {{ $t('recordList.import.hasRequiredFileFields') }}: {{ showRequiredFields }}
      </div>

      <b-table
        :fields="tableFields"
        :items="rows"
        head-variant="light"
        sticky-header
        style="max-height: 70vh;"
        class="field-table mb-0"
      >
        <template #head(selected)>
          <b-form-checkbox
            class="pr-0"
            :checked="selectAll"
            @change="onSelectAll"
          />
        </template>

        <template #cell(selected)="data">
          <b-form-checkbox
            v-model="data.item.selected"
            class="pr-0"
          />
        </template>
        <template #cell(moduleField)="data">
          <c-input-select
            v-model="data.item.moduleField"
            :options="moduleFields"
            :reduce="o => o.name"
            :placeholder="$t('recordList.import.pickModuleField')"
            :class="{ 'border border-danger': data.item.selected && !data.item.moduleField }"
            @input="moduleChanged(data)"
          />
          <span
            v-if="data.item.fileColumn === 'id'"
            class="small text-muted"
          >
            {{ $t('recordList.import.idFieldDescription') }}
          </span>
        </template>
      </b-table>
    </b-form-group>

    <div class="mt-auto p-3">
      <b-button
        variant="light"
        class="float-left"
        @click="$emit('back')"
      >
        {{ $t('general.label.back') }}
      </b-button>

      <b-button
        variant="primary"
        :disabled="!canContinue"
        class="float-right"
        @click="nextStep"
      >
        {{ $t('general.label.import') }}
      </b-button>
    </div>
  </div>
</template>

<script>
export default {
  i18nOptions: {
    namespaces: 'block',
  },

  props: {
    session: {
      type: Object,
      required: true,
      default: () => ({}),
    },

    module: {
      type: Object,
      required: true,
      default: () => ({}),
    },
  },

  data () {
    return {
      rows: [],
      unsupportedFields: [],
    }
  },

  computed: {
    // @note don't use selectAll directly as v-model, since it fires
    // when ever data changes, eaven by this getter.
    // Above onSelectAll doesn't fire if getter changes
    selectAll: {
      get () {
        // if all rows selected
        return this.rows.reduce((acc, { selected }) => acc && selected, true)
      },

      set (v) {
        // set all rows to what this one is set to
        this.rows.forEach((r) => { r.selected = v })
      },
    },

    canContinue () {
      // has anything selected && all selected rows have mapped module fields
      const selected = this.rows.filter(({ selected }) => selected)
      const named = selected.filter(({ moduleField }) => !!moduleField)
      return !!selected.length && selected.length === named.length && !this.hasRequiredFileFields
    },

    tableFields () {
      return [
        {
          key: 'selected',
          label: '',
          tdClass: 'picker align-middle',
          thStyle: 'width: 30px',
        },
        {
          key: 'fileColumn',
          label: this.$t('recordList.import.fileColumns'),
          tdClass: 'align-middle',
        },
        {
          key: 'moduleField',
          label: this.$t('recordList.import.moduleFields'),
          thStyle: 'width: 25rem',
        },
      ]
    },

    moduleFields () {
      return [
        ...this.module.fields,
        ...this.module.systemFields().map(({ name }) => {
          return {
            label: this.$t(`field:system.${name}`),
            name,
          }
        }),
      ].filter(({ kind }) => !['File'].includes(kind))
        .map(field => field.isRequired === true ? { ...field, label: field.label + '*' } : field)
    },

    requiredFields () {
      return this.module.fields.filter(field => field.isRequired === true)
    },

    filteredRows () {
      // if required field is selected
      const result = this.rows.filter(row => {
        return this.requiredFields.some(field => {
          return row.moduleField === field.name
        })
      })
      // filter duplicated selected required fields, if user clicks one multiple times
      return result.filter((value, index, self) => self.findIndex(v => v.moduleField === value.moduleField) === index)
    },

    hasRequiredFileFields () {
      return !(this.requiredFields.length === this.filteredRows.length)
    },

    showRequiredFields () {
      const array = []
      // do not show required fields that have been already selected
      let result = this.requiredFields.filter(field => {
        return this.filteredRows.some(row => {
          return field.name !== row.moduleField
        })
      })
      if (result.length === 0) result = this.requiredFields
      result.map(field => array.push(field.label))
      return array.join(', ').toString()
    },
  },

  created () {
    // Prep row object for us to alter
    const { fields = {} } = this.session

    const moduleFields = {
      id: 'recordID',
    }

    this.module.fields.forEach(({ name }) => {
      moduleFields[name] = name
    })

    this.module.systemFields().forEach(({ name }) => {
      moduleFields[name] = name
    })

    this.rows = Object.entries(fields)
      .map(([fileColumn, moduleField]) => {
        moduleField = moduleField || moduleFields[fileColumn]

        return {
          selected: false,
          fileColumn,
          moduleField,
        }
      })
  },

  methods: {
    moduleChanged (data) {
      if (data.item.moduleField) {
        const result = this.rows.find(row => row.moduleField === data.item.moduleField)
        result.selected = true
      }
    },

    nextStep () {
      if (!this.canContinue) {
        return
      }

      // convert to api's structure
      const rtr = {}
      this.rows.forEach(({ selected, fileColumn, moduleField }) => {
        if (selected) {
          rtr[fileColumn] = moduleField
        }
      })

      this.$emit('fieldsMatched', rtr)
    },

    onSelectAll (e) {
      this.selectAll = e
    },
  },

}
</script>
