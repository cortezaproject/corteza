<template>
  <c-form-table-wrapper
    :labels="{
      addButton: $t('general:label.add')
    }"
    @add-item="addParam"
  >
    <b-table-simple
      v-if="groups.length"
      responsive
      borderless
      small
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
          v-for="(group, index) in groups"
          :key="index"
        >
          <b-td>
            <b-form-input
              v-model="group.name"
              :placeholder="$t('datasources:new.name')"
            />
          </b-td>

          <b-td>
            <b-form-input
              v-model="group.label"
              :placeholder="$t('datasources:new.label')"
            />
          </b-td>

          <b-td>
            <b-form-input
              v-model="group.def.raw"
              :placeholder="$t('datasources:expression')"
            />
          </b-td>

          <b-td class="align-middle">
            <c-input-confirm
              show-icon
              @confirmed="deleteGroup(index)"
            />
          </b-td>
        </b-tr>
      </b-tbody>
    </b-table-simple>
  </c-form-table-wrapper>
</template>

<script>
export default {
  props: {
    groupBy: {
      type: Array,
      required: true,
    },
  },

  data () {
    return {
    }
  },

  computed: {
    groups: {
      get () {
        return this.groupBy || []
      },

      set (groupBy) {
        this.$emit('update:groupBy', groupBy)
      },
    },
  },

  methods: {
    addParam () {
      this.groups.push({
        name: '',
        label: '',
        def: {
          raw: '',
        },
      })
    },

    deleteGroup (index) {
      this.groups.splice(index, 1)
    },
  },
}
</script>

<style lang="scss" scoped>
.table td.fit,
.table th.fit {
  white-space: nowrap;
  width: 1%;
}

.btn-add-group {
  &:hover, &:active {
    background-color: var(--primary) !important;
    color: var(--white) !important;
  }
}

.filter-border {
  background-image: linear-gradient(to left, lightgray, lightgray);
  background-repeat: no-repeat;
  background-size: 100% 1px;
  background-position: center;
}
</style>
