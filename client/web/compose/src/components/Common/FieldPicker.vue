<template>
  <div
    class="d-flex flex-column border border-light rounded p-2 "
  >
    <c-item-picker
      :value.sync="selected"
      :options="options"
      :labels="{
        searchPlaceholder: $t('selector.search'),
        availableItems: $t('available-items'),
        selectAllItems: $t('selector.selectAll'),
        selectedItems: $t('selected-items'),
        unselectAllItems: $t('selector.unselectAll'),
        noItemsFound: $t('no-items-found'),
      }"
    >
      <template
        #default="{ field }"
      >
        <b class="cursor-default text-dark">
          <template
            v-if="field.label"
          >
            {{ field.label }}
            <template v-if="field.originalName">
              ({{ field.originalName }})
            </template>
          </template>
          <template
            v-else
          >
            {{ field.originalName || field.name }}
          </template>
          <template
            v-if="field.isRequired"
          >
            *
          </template>
        </b>
        <small
          v-if="field.isSystem"
          class="cursor-default ml-1 text-truncate"
        >
          {{ $t('selector.systemField') }}
        </small>
        <small
          v-if="field.isParentField"
          class="cursor-default ml-1 text-truncate text-primary"
        >
          {{ field.parentModuleName }}
        </small>
      </template>
    </c-item-picker>
  </div>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'
const { CItemPicker } = components

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  components: {
    CItemPicker,
  },

  props: {
    // source of fields
    module: {
      type: Object,
      required: true,
    },

    // array of objects, list of all selected fields
    fields: {
      type: Array,
      required: true,
    },

    disabledTypes: {
      type: Array,
      default: () => [],
    },

    disableSystemFields: {
      type: Boolean,
      default: false,
    },

    disableSorting: {
      type: Boolean,
    },

    // source of additional (system fields) we'll use
    systemFields: {
      type: Array,
      default: null,
    },

    fieldSubset: {
      type: Array,
      default: null,
    },

    // NEW: optional parent module object
    extraModule: {
      type: Object,
      default: null,
    },

    // NEW: array of field objects from the parent module
    extraModuleFields: {
      type: Array,
      default: () => [],
    },
  },

  computed: {
    /**
     * Converting old-way of passing data to picker components
     * where arrays of objects where exchanged and stored
     *
     * Item picker does not do that any more so we need to adapt
     * whatever we get from the outside to array of picked field names
     * and vice versa!
     */
    selected: {
      get () {
        // Only need names of the fields
        return this.fields.map(({ name }) => name)
      },

      set (selected) {
        // take list of field names passed to the setter
        // and filter out the options to recreate the list
        // module field objects
        const fields = selected.map(s => {
          return (this.options.find(({ value }) => value === s) || {}).field
        }).filter(f => f)

        this.$emit('update:fields', fields)
      },
    },

    options () {
      let mFields = [...(this.fieldSubset ? this.module.filterFields(this.fieldSubset) : this.module.fields)]

      if (this.disabledTypes.length > 0) {
        mFields = mFields.filter(({ kind }) => !this.disabledTypes.find(t => t === kind))
      }

      // Filter out fields the user cannot read (RBAC)
      mFields = mFields.filter(({ canReadRecordValue }) => canReadRecordValue !== false)

      let sysFields = []

      if (this.disableSystemFields && mFields) {
        mFields = mFields.filter(({ isSystem }) => !isSystem)
      } else if (!this.fieldSubset) {
        sysFields = this.module.systemFields().map(sf => {
          sf.label = this.$t(`field:system.${sf.name}`)
          return sf
        })

        if (this.systemFields) {
          sysFields = sysFields.filter(({ name }) => this.systemFields.find(sf => sf === name))
        }
      }

      if (!this.disableSorting && mFields) {
        mFields.sort((a, b) => a.label.localeCompare(b.label))
      }

      // Build base options from child module fields
      let baseOptions = []
      if (mFields && sysFields) {
        baseOptions = [
          ...[...mFields],
          ...sysFields,
        ].map(field => ({
          value: field.name,
          text: [
            field.name,
            field.label,
            field.kind,
            field.isSystem ? 'system' : '',
            field.isRequired ? 'required' : '',
          ].join(' '),
          field,
        }))
      } else {
        baseOptions = Object.keys(this.module).map(key => {
          return this.module[key]
        }).map(f => ({
          ...f,
          field: {
            name: f.text,
          },
        }))
      }

      // Build parent module field options (if extraModule and extraModuleFields are provided)
      const extraFieldOptions = (this.extraModule && this.extraModuleFields.length > 0)
        ? this.extraModuleFields.map(field => {
          // Use moduleID::fieldName as the unique key to avoid name collisions
          const uniqueName = `${this.extraModule.moduleID}::${field.name}`
          return {
            value: uniqueName,
            text: [field.name, field.label, field.kind, this.extraModule.name].join(' '),
            field: {
              ...field,
              // Override name with the namespaced version for storage/lookup
              name: uniqueName,
              label: field.label || field.name,
              // Flags and metadata for downstream use
              isParentField: true,
              parentModuleID: this.extraModule.moduleID,
              parentModuleName: this.extraModule.name,
              originalName: field.name,
            },
          }
        })
        : []

      // Concat extra options after the base options
      return [...baseOptions, ...extraFieldOptions]
    },
  },
}
</script>
