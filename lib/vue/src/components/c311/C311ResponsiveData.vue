<template>
  <div class="c311-responsive-data" data-c311-responsive-data>
    <div class="c311-responsive-data__table" tabindex="0" role="region" :aria-label="label">
      <table class="table table-hover mb-0">
        <thead>
          <tr>
            <th v-for="column in visibleColumns" :key="column.key" scope="col">
              {{ column.label }}
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in items" :key="item[rowKey]">
            <td v-for="column in visibleColumns" :key="column.key" :data-label="column.label">
              {{ valueFor(item, column) }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="c311-responsive-data__cards">
      <article v-for="item in items" :key="item[rowKey]" class="border rounded p-3 mb-2">
        <dl class="mb-0">
          <div v-for="column in visibleColumns" :key="column.key" class="d-flex justify-content-between py-1">
            <dt class="font-weight-bold mr-3">{{ column.label }}</dt>
            <dd class="mb-0 text-right">{{ valueFor(item, column) }}</dd>
          </div>
        </dl>
      </article>
    </div>
  </div>
</template>

<script>
export default {
  name: 'C311ResponsiveData',
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    columns: {
      type: Array,
      required: true,
    },
    rowKey: {
      type: String,
      default: 'id',
    },
    label: {
      type: String,
      default: 'Data table',
    },
  },
  computed: {
    visibleColumns () {
      return this.columns.filter(column => column.visible !== false)
    },
  },
  methods: {
    valueFor (item, column) {
      return typeof column.format === 'function' ? column.format(item[column.key], item) : item[column.key]
    },
  },
}
</script>

<style scoped>
.c311-responsive-data__cards {
  display: none;
}

.c311-responsive-data__table {
  overflow-x: auto;
  max-width: 100%;
}

@media (max-width: 767px) {
  .c311-responsive-data__table {
    display: none;
  }

  .c311-responsive-data__cards {
    display: block;
  }
}

@media (min-width: 768px) and (max-width: 1023px) {
  .c311-responsive-data__table th:nth-child(n+4),
  .c311-responsive-data__table td:nth-child(n+4) {
    display: none;
  }
}
</style>
