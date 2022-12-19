<template>
  <div>
    <draggable
      v-if="!textInput"
      :list.sync="items"
      group="sort"
      handle=".grab"
    >
      <b-form-row
        v-for="(column, index) in items"
        :key="index"
        class="mb-1"
      >
        <b-col
          cols="1"
          class="d-flex align-items-center justify-content-center"
        >
          <font-awesome-icon
            :icon="['fas', 'bars']"
            class="grab text-grey"
          />
        </b-col>

        <b-col
          cols="5"
        >
          <b-form-select
            v-model="column.field"
            :options="availableFields"
            text-field="label"
            value-field="name"
            class="rounded"
          >
            <template #first>
              <b-form-select-option
                :value="undefined"
                disabled
              >
                {{ labels.none }}
              </b-form-select-option>
            </template>
          </b-form-select>
        </b-col>

        <b-col
          cols="6"
          class="d-flex align-items-center justify-content-around"
        >
          <b-form-radio-group
            v-model="column.descending"
            :options="sortDirections"
            buttons
            size="sm"
            button-variant="outline-primary"
          />
          <c-input-confirm
            variant="link"
            size="lg"
            button-class="text-dark px-0"
            @confirmed="items.splice(index, 1)"
          />
        </b-col>
      </b-form-row>
    </draggable>

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
      class="d-flex align-items-center"
    >
      <b-button
        v-if="!textInput"
        variant="link"
        class="d-flex align-items-center px-0 text-decoration-none"
        @click="items.push({ field: undefined, descending: false })"
      >
        <font-awesome-icon
          :icon="['fas', 'plus']"
          size="sm"
          class="mr-1"
        />
        {{ labels.add }}
      </b-button>

      <b-button
        v-if="allowTextInput"
        variant="link"
        size="sm"
        class="text-decoration-none ml-auto"
        @click="textInput = !textInput"
      >
        {{ labels.toggleInput }}
      </b-button>
    </div>
  </div>
</template>

<script>
import Draggable from 'vuedraggable'

export default {
  components: {
    Draggable,
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
}
</script>
