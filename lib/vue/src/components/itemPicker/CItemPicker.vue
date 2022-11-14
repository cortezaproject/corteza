<template>
  <div
    class="overflow-hidden d-flex flex-column w-100 vh-100"
  >
    <b-card
      no-body
      class="h-100"
    >
      <b-card-header
        class="bg-white p-0"
      >

      <c-input-search
        v-if="!hideFilter"
        v-model.trim="query"
        :disabled="_disabledFilter"
        :placeholder="labels.searchPlaceholder"
      />

      </b-card-header>
      <b-card-body
        class="d-flex p-0"
      >
        <b-card
          no-body
          class="col-sm-6 col-12 h-100 p-0"
        >
          <b-card-header
            class="bg-white py-2 pl-0 pr-2"
          >
            <div
              class="d-flex align-items-center"
            >
              <label
                class="text-primary mb-0"
              >
                {{ labels.availableItems }}
              </label>
              <b-button
                v-show="filteredAvailable.length && !disabled"
                variant="link"
                class="ml-auto px-0 text-muted"
                @click="selectAll()"
              >
                {{ labels.selectAllItems }}
              </b-button>
            </div>
          </b-card-header>
          <b-card-body
            class="overflow-auto py-0 pl-0 pr-2"
          >
            <b-list-group
              vertical
              class="h-100"
            >
              <draggable
                v-model="filteredAvailable"
                :sort="!_disabledSorting"
                :move="_disableDragging"
                draggable=".item"
                group="items"
                class="overflow-auto h-100"
              >
                <b-list-group-item
                  v-for="item in filteredAvailable"
                  :key="item.value"
                  class="item mb-3 border rounded"
                  @dblclick="select(item)"
                >
                  <c-item-picker-item
                    :item="item"
                    :disabled="disabled"
                    :disabled-dragging="disabledDragging"
                    :disabled-sorting="disabledSorting"
                    :hide-icons="hideIcons"
                    @select="select(item)"
                  >
                    <template
                      v-for="(_, slot) of $scopedSlots"
                      v-slot:[slot]="scope"
                    >
                      <slot
                        :name="slot"
                        :textField="textField"
                        :disabled="disabled"
                        :disabled-dragging="disabledDragging"
                        :disabled-sorting="disabledSorting"
                        :hide-icons="hideIcons"
                        v-bind="scope"
                      />
                    </template>
                  </c-item-picker-item>
                </b-list-group-item>

                <template #footer>
                  <h6
                    v-if="!filteredAvailable.length && query"
                    class="text-center my-4"
                  >
                    {{ labels.noItemsFound }}
                  </h6>
                </template>
              </draggable>
            </b-list-group>
          </b-card-body>
        </b-card>
        <b-card
          no-body
          class="h-100 pl-sm-0 col-sm-6 col-12 p-0"
        >
          <b-card-header
            class="bg-white py-2 pl-2 pr-0"
          >
            <div
              class="d-flex align-items-center"
            >
              <label
                class="text-primary mb-0"
              >
                {{ labels.selectedItems }}
              </label>
              <b-button
                v-show="filteredSelected.length && !disabled"
                variant="link"
                class="ml-auto px-0 text-muted"
                @click="unselectAll()"
              >
                {{ labels.unselectAllItems }}
              </b-button>
            </div>
          </b-card-header>
          <b-card-body
            class="overflow-auto py-0 pl-2 pr-0"
          >
            <b-list-group
              vertical
              class="h-100"
            >
              <draggable
                v-model="filteredSelected"
                :sort="!_disabledSorting"
                :move="_disableDragging"
                draggable=".item"
                group="items"
                class="overflow-auto h-100"
              >
                <b-list-group-item
                  v-for="item in filteredSelected"
                  :key="item.value"
                  class="item mb-3 border rounded"
                  @dblclick="unselect(item)"
                >
                  <c-item-picker-item
                    :item="item"
                    :disabled="disabled"
                    :disabled-dragging="disabledDragging"
                    :disabled-sorting="disabledSorting"
                    :hide-icons="hideIcons"
                    selected
                    @unselect="unselect(item)"
                  >
                    <template
                      v-for="(_, slot) of $scopedSlots"
                      v-slot:[slot]="scope"
                    >
                      <slot
                        :name="slot"
                        :textField="textField"
                        :disabled="disabled"
                        :disabled-dragging="disabledDragging"
                        :disabled-sorting="disabledSorting"
                        :hide-icons="hideIcons"
                        v-bind="scope"
                      />
                    </template>
                  </c-item-picker-item>
                </b-list-group-item>

                <template #footer>
                  <h6
                    v-if="!filteredSelected.length && query"
                    class="text-center my-4"
                  >
                    {{ labels.noItemsFound }}
                  </h6>
                </template>
              </draggable>
            </b-list-group>
          </b-card-body>
        </b-card>
      </b-card-body>
    </b-card>
  </div>
</template>
<script>
import draggable from 'vuedraggable'
import CItemPickerItem from './CItemPickerItem.vue'
import CInputSearch from '../input/CInputSearch.vue'
import { throttle } from 'lodash'

export default {
  name: 'CItemPicker',

  components: {
    draggable,
    CItemPickerItem,
    CInputSearch,
  },

  props: {
    /**
     * List of all items, available and selected/picked
     *
     * Internally, we'll build 2 arrays for each group (available + selected)
     * and work with them.
     * On the outside, we'll always deal input array of (full) items and
     * array of selected/picked values of those items.
     *
     * This component mimics the behaviour of <b-form-select> component from
     * Vue Bootstrap and could serve as a drop-in replacement!
     */
    options: {
      type: Array,
      required: true,
    },

    /**
     * List of values that can be found in the items
     */
    value: {
      type: Array,
      required: true,
    },

    valueField: {
      type: String,
      default: 'value',
    },

    textField: {
      type: String,
      default: 'text',
    },

    labels: {
      type: Object,
      default: () => ({
        searchPlaceholder: 'Filter items',
        availableItems: 'Available',
        selectAllItems: 'Select all',
        selectedItems: 'Selected',
        unselectAllItems: 'Unselect all',
        noItemsFound: 'No items found',
      }),
    },

    disabled: {
      type: Boolean,
    },

    disabledFilter: {
      type: Boolean,
    },

    disabledSorting: {
      type: Boolean,
    },

    disabledDragging: {
      type: Boolean,
    },

    hideIcons: {
      type: Boolean,
    },

    hideFilter: {
      type: Boolean,
    },
  },

  data () {
    return {
      query: '',
      available: [],
      selected: [],
    }
  },

  computed: {
    _disabledFilter () {
      return this.disabled || this.disabledFilter
    },

    _disabledSorting () {
      return this.disabled || this.disabledSorting
    },

    /**
     * Provides list of all available items, filtered by query
     * and a setter for draggable component to update available items on drag
     */
    filteredAvailable: {
      get () {
        const q = this.query.toLowerCase()
        return this.available.filter(i => i.text.toLowerCase().indexOf(q) > -1)
      },

      set (items) {
        this.available = items
      },
    },

    /**
     * Provides list of all selected items, filtered by query
     * and a setter for draggable component to update selected items on drag
     */
    filteredSelected: {
      get () {
        const q = this.query.toLowerCase()
        return this.selected.filter(i => i.text.toLowerCase().indexOf(q) > -1)
      },

      set (items) {
        this.selected = items
      },
    },
  },

  watch: {
    /**
     * Update parent component (if needed)
     * @param v
     */
    selected: {
      immediate: false,
      handler (items) {
        const value = items.map(i => i[this.valueField])

        // satisfy value.sync
        this.$emit('update:value', value)

        // satisfy v-model
        this.$emit('input', value)
      }
    },

    options: {
      deep: true,
      immediate: true,
      handler () {
        this.sync()
      },
    },

    value: {
      immediate: false,
      handler (value = [], oldValue = []) {
        /**
        * Make sure we do not fall into an infinite loop
        * 
        * If we update the value thenn sync will trigger recomputation of selected
        * Which then emits the update event and the loop will begin
        */
        if (value.length === oldValue.length) {
          if (value.filter(v => !oldValue.includes(v)).length === 0) {
            return
          }
        }

        this.sync()
      },
    },
  },

  created () {
    this.options.forEach(o => {
      if (typeof o !== 'object') {
        throw new Error('expecting array of objects for options prop')
      }
    })
  },

  methods: {
    _disableDragging (e) {
      if (this.disabledDragging && e.to !== e.from) {
        return false
      }
    },

    selectAll () {
      if (this.disabledSorting) {
        // sorting disabled, just reuse list from the frozen list
        this.selected = [...this.frozen()]
      } else {
        // sorting enabled, push all items left in available
        // at the end of the list
        this.selected.push(...this.available)
      }

      this.available = []
    },

    select: throttle(function (item) {
      // remove available
      this.available = this.available.filter(i => i !== item)

      // put the selected item in the place where it was (if sorting is disabled)
      // or at the end
      if (this.disabledSorting) {
        this.selected = this.frozen().filter(i => !this.available.includes(i))
      } else {
        if (!this.selected.some(({ value = '' }) => value === item.value)) {
          this.selected.push(item)
        }
      }
    }, 300),

    unselectAll () {
      this.available = [...this.frozen()]
      this.selected = []
    },

    unselect (item) {
      // filtering out the unselected item
      this.selected = this.selected.filter(i => i !== item)

      // sync available from the list of frozen items without selected
      this.available = this.frozen().filter(i => !this.selected.includes(i))
    },

    /**
     * Returns true given item is selected or not
     *
     * Item is considered selected if value of item's value-field is inside the selected array.
     *
     * @param item object
     * @returns boolean
     */
    isPicked (item) {
      return item[this.valueField] && this.value.includes(item[this.valueField])
    },

    sync () {
      /**
       * filter all unpicked options, freeze each item in the array and
       * build list of available items
       */
      this.available = this.frozen().filter(opt => !this.isPicked(opt))

      /**
       * filter all unpicked options, freeze each item in the array and
       * build list of selected items, traverse value to keep order
       */
      this.selected = this.value.map(v => {
        return this.frozen().find(item => item[this.valueField] === v)
      }).filter(f => f)
    },

    frozen () {
      return this.options.map(Object.freeze)
    },
  },
}
</script>
