<template>
  <div>
    <b-table-simple
      v-if="columns.length"
      responsive
      borderless
      small
      class="mb-0"
    >
      <b-thead>
        <b-tr>
          <b-th
            class="w-25"
          >
            {{ $t('datasources:name') }}
          </b-th>
          <b-th
            class="w-25"
          >
            {{ $t('datasources:label') }}
          </b-th>
          <b-th
            class="w-50"
          >
            {{ $t('datasources:expression') }}
          </b-th>
        </b-tr>
      </b-thead>

      <b-tbody>
        <b-tr
          v-for="(column, index) in columns"
          :key="index"
        >
          <b-td>
            <b-form-input
              v-model="column.name"
              :placeholder="$t('datasources:new.name')"
            />
          </b-td>
          <b-td>
            <b-form-input
              v-model="column.label"
              :placeholder="$t('datasources:new.label')"
            />
          </b-td>
          <b-td>
            <b-form-input
              v-model="column.def.raw"
              :placeholder="$t('datasources:expression')"
            />
          </b-td>
          <b-td
            class="d-flex align-items-center justify-content-center pl-2 pr-0"
          >
            <c-input-confirm
              variant="link"
              size="lg"
              button-class="text-dark px-0"
              @confirmed="deleteColumn(index)"
            />
          </b-td>
        </b-tr>
      </b-tbody>
    </b-table-simple>

    <b-button
      variant="link text-decoration-none"
      class="px-0"
      @click="addColumn()"
    >
      <font-awesome-icon
        :icon="['fas', 'plus']"
        size="sm"
        class="mr-1"
      />
      {{ $t('datasources:add') }}
    </b-button>
  </div>
</template>

<script>
export default {
  props: {
    aggregate: {
      type: Array,
      required: true,
    },
  },

  data () {
    return {
    }
  },

  computed: {
    columns: {
      get () {
        return this.aggregate || []
      },

      set (aggregate) {
        this.$emit('update:aggregate', aggregate)
      },
    },
  },

  methods: {
    addColumn () {
      this.columns.push({
        name: '',
        label: '',
        def: {
          raw: '',
        },
      })
    },

    deleteColumn (index) {
      this.columns.splice(index, 1)
    },
  },
}
</script>
