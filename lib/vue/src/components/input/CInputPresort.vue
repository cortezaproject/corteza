<template>
  <c-form-table-wrapper
    :labels="{ addButton: labels.addButton }"
    :hide-add-button="textInput"
    @add-item="items.push({ field: undefined, descending: false })"
  >
    <b-form-group
      v-if="!textInput"
      :label="labels.title"
       label-class="text-primary"
    >
      <b-table-simple
        borderless
        small
        responsive="lg"
        class="mb-0"
      >
        <draggable
          :list.sync="items"
          group="sort"
          handle=".grab"
          tag="tbody"
        >
          <tr
            v-for="(column, index) in items"
            :key="index"
          >
            <td
              class="grab text-center align-middle"
              style="width: 40px;"
            >
              <font-awesome-icon
                :icon="['fas', 'bars']"
                class="text-secondary"
              />
            </td>
            <td
              class="align-middle"
              style="min-width: 250px;"
            >
              <c-input-select
                v-model="column.field"
                :options="availableFields"
                :reduce="o => o.name"
                :placeholder="labels.none"
                class="rounded"
              />
            </td>
            <td
              class="text-center align-middle"
              style="min-width: 200px;"
            >
              <b-form-radio-group
                v-model="column.descending"
                :options="sortDirections"
                buttons
                button-variant="outline-primary"
                class="bg-white"
              />
            </td>
            <td
              class="align-middle text-right"
              style="min-width: 80px; width: 80px;"
            >
              <c-input-confirm
                show-icon
                @confirmed="items.splice(index, 1)"
              />
            </td>
          </tr>
        </draggable>
      </b-table-simple>
    </b-form-group>

    <div
      v-else
    >
      <b-form-textarea
        v-model="presortValue"
        :placeholder="labels.placeholder"
      />
      <b-form-text>
        {{ labels.footnote }}
      </b-form-text>
    </div>

    <div
      v-if="allowTextInput"
      class="d-flex align-items-center mt-1"
    >
      <b-button
        variant="link"
        size="sm"
        class="text-decoration-none ml-auto"
        @click="textInput = !textInput"
      >
        {{ labels.toggleInput }}
      </b-button>
    </div>
  </c-form-table-wrapper>
</template>

<script>
import Draggable from 'vuedraggable'
import CInputSelect from './CInputSelect.vue'
import CFormTableWrapper from '../wrapper/CFormTableWrapper.vue'

export default {
  components: {
    Draggable,
    CInputSelect,
    CFormTableWrapper,
  },

  props: {
    value: {
      type: String,
      required: true,
    },

    fields: {
      type: Array,
      required: true,
    },

    labels: {
      type: Object,
      required: true,
    },

    allowTextInput: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      items: [],

      textInput: false,
    }
  },

  computed: {
    presortValue: {
      get () {
        return this.value
      },

      set (value) {
        this.$emit('input', value)
      },
    },

    sortDirections () {
      return [
        { value: false, text: this.labels.ascending },
        { value: true, text: this.labels.descending },
      ]
    },

    availableFields () {
      return this.fields.map(f => ({ ...f,label: `${f.label} (${f.name})` }))
    }
  },

  watch: {
    value: {
      immediate: true,
      handler (value) {
        if (value) {
          const sort = value.includes(',') ? value.split(',') : [value]

          this.items = sort.map(field => {
            let descending = false

            if (field.includes(' ')) {
              field = field.split(' ')[0]
              descending = true
            }

            return {
              field,
              descending: !!descending,
            }
          })
        } else {
          this.items = [{
            field: undefined,
            descending: false,
          }]
        }
      },
    },

    items: {
      deep: true,
      handler (items = [], oldItems = undefined) {
        if (oldItems) {
          this.$emit('input', items.filter(({ field }) => field).map(({ field, descending }) => {
            return descending ? `${field} DESC` : field
          }).join(','))
        }
      },
    },
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    setDefaultValues () {
      this.items = []
      this.textInput = false
    },
  },
}
</script>
