<template>
  <b-card
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-form-group>
      <div class="mb-4">
        <label>
          {{ $t('recordList.import.matchFields') }}
        </label>
        <br>
        <small
          v-if="hasRequiredFileFields"
        >
          {{ $t('recordList.import.hasRequiredFileFields') }}: {{ showRequiredFields }}
        </small>
      </div>

      <b-table
        small
        :fields="tableFields"
        :items="rows"
        class="field-table"
      >
        <template v-slot:head(selected)>
          <b-form-checkbox
            class="pr-0"
            :checked="selectAll"
            @change="onSelectAll"
          />
        </template>

        <template v-slot:cell(selected)="data">
          <b-form-checkbox
            v-model="data.item.selected"
            class="pr-0"
          />
        </template>
        <template v-slot:cell(moduleField)="data">
          <b-form-select
            v-model="data.item.moduleField"
            :options="moduleFields"
            @change="moduleChanged(data)"
          >
            <template slot="first">
              <option
                :value="undefined"
                disabled
              >
                {{ $t('recordList.import.pickModuleField') }}
              </option>
            </template>
          </b-form-select>
        </template>
      </b-table>
    </b-form-group>

    <div slot="footer">
      <b-button
        variant="outline-dark"
        class="float-left"
        @click="$emit('back')"
      >
        {{ $t('general.label.back') }}
      </b-button>

      <b-button
        variant="dark"
        :disabled="!canContinue"
        class="float-right"
        @click="nextStep"
      >
        {{ $t('general.label.import') }}
      </b-button>
    </div>
  </b-card>
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
          thClass: 'pb-1',
        },
        {
          key: 'fileColumn',
          label: this.$t('recordList.import.fileColumns'),
          tdClass: 'align-middle',
          thClass: 'pb-1',
        },
        {
          key: 'moduleField',
          label: this.$t('recordList.import.moduleFields'),
          thClass: 'pb-1',
        },
      ]
    },

    moduleFields () {
      return this.module.fields
        .filter(({ kind }) => !['File'].includes(kind))
        .map(field => field.isRequired === true ? { ...field, label: field.label + '*' } : field)
        .map(({ name: value, label }) => ({ value, text: label || value }))
        .sort((a, b) => a.text.localeCompare(b.text))
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
    const moduleFields = {}
    this.module.fields.forEach(({ name }) => {
      moduleFields[name] = name
    })

    this.rows = Object.entries(fields)
      .map(([fileColumn, moduleField]) => ({
        selected: false,
        fileColumn,
        moduleField: moduleField || moduleFields[fileColumn],
      }))
  },

  methods: {
    moduleChanged (data) {
      const result = this.rows.find(row => row.fileColumn === data.item.fileColumn)
      result.selected = true
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
