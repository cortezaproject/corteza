<template>
  <div
    class="d-flex"
  >
    <c-item-picker
      :value.sync="selected"
      :options="options"
      :labels="{
        searchPlaceholder: $t('list:searchPlaceholder'),
        availableItems: $t('available-columns'),
        selectAllItems: $t('select-all'),
        selectedItems: $t('selected-columns'),
        unselectAllItems: $t('unselect-all'),
        noItemsFound: $t('no-columns-found'),
      }"
      style="max-height: 40vh;"
    >
      <template
        v-slot:default="{ field }"
      >
        <b class="cursor-default text-dark">
          <template
            v-if="field.label"
          >
            {{ field.label }} ({{ field.name }})
          </template>
          <template
            v-else
          >
            {{ field.name }}
          </template>
        </b>
      </template>
    </c-item-picker>
  </div>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'
const { CItemPicker } = components

export default {
  i18nOptions: {
    namespaces: 'builder',
  },

  components: {
    CItemPicker,
  },

  props: {
    // source of fields
    allItems: {
      type: Array,
      required: true,
    },

    // array of objects, list of all selected fields
    selectedItems: {
      type: Array,
      required: true,
    },

    disableSorting: {
      type: Boolean,
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
        return this.selectedItems.map(({ name }) => name)
      },

      set (selected) {
        // take list of field names passed to the setter
        // and filter out the options to recreate the list
        // module field objects
        const fields = selected.map(s => {
          return (this.options.find(({ value }) => value === s) || {}).field
        }).filter(f => f)

        this.$emit('update:selected-items', fields)
      },
    },

    options () {
      const fields = this.allItems

      if (!this.disableSorting) {
        fields.sort((a, b) => a.label.localeCompare(b.label))
      }

      return [...fields,
      ].map(field => ({
        value: field.name,
        text: [
          field.name,
          field.label,
          field.kind,
          field.system ? 'system' : '',
        ].join(' '),
        field,
      }))
    },
  },
}
</script>
