<template>
  <b-table
    id="expression-table"
    fixed
    borderless
    hover
    head-row-variant="secondary"
    details-td-class="bg-white"
    class="mb-4"
    :items="items"
    :fields="fields"
    thead-tr-class="border-thick"
    :tbody-tr-class="rowClass"
    @row-clicked="item => $set(item, '_showDetails', !item._showDetails)"
  >
    <template #cell(target)="{ value }">
      <var>{{ value }}</var>
    </template>

    <template #cell(type)="{ value }">
      <var>{{ value }}</var>
    </template>

    <template #[`cell(${valueField})`]="{ value, item, index }">
      <div
        class="d-flex justify-content-between align-items-center"
      >
        <div
          class="text-truncate"
          :class="{ 'w-75': item._showDetails}"
        >
          <samp>{{ value }}</samp>
        </div>

        <b-button
          v-if="item._showDetails"
          variant="outline-danger"
          class="position-absolute trash border-0"
          @click="$emit('remove', index)"
        >
          <font-awesome-icon
            :icon="['far', 'trash-alt']"
          />
        </b-button>
      </div>
    </template>

    <template #row-details="{ item, index }">
      <div class="arrow-up" />

      <b-card
        class="bg-light"
        body-class="px-4 pb-3"
      >
        <b-form-group
          label-class="text-primary"
        >
          <b-form-input
            v-model="item.target"
            placeholder="Target"
            @input="$root.$emit('change-detected')"
          />
        </b-form-group>

        <b-form-group
          label-class="text-primary"
          :description="getTypeDescription(item.type)"
        >
          <vue-select
            v-model="item.type"
            :options="types"
            :clearable="false"
            :filter="varFilter"
            @input="$root.$emit('change-detected')"
          />
        </b-form-group>

        <b-form-group
          class="mb-0"
        >
          <expression-editor
            :value.sync="item[valueField]"
            lang="javascript"
            show-line-numbers
            @open="$emit('open-editor', index)"
            @input="$root.$emit('change-detected')"
          />
        </b-form-group>
      </b-card>
    </template>
  </b-table>
</template>

<script>
import ExpressionEditor from './ExpressionEditor.vue'
import { objectSearchMaker } from '../lib/filter'
import { VueSelect } from 'vue-select'

export default {
  components: {
    ExpressionEditor,
    VueSelect,
  },

  props: {
    valueField: {
      type: String,
      required: true,
    },

    items: {
      type: Array,
      required: true,
    },

    fields: {
      type: Array,
      required: true,
    },

    types: {
      type: Array,
      required: true,
    },
  },

  methods: {
    varFilter: objectSearchMaker('text'),

    rowClass (item, type) {
      return item._showDetails && type === 'row' ? 'border-thick' : 'border-thick-transparent'
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
