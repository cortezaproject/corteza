<template>
  <vue-select
    v-model="_value"
    v-bind="$attrs"
    data-test-id="select"
    :clearable="clearable"
    :options="options"
    :searchable="searchable"
    :calculate-position="calculateDropdownPosition"
    :append-to-body="appendToBody"
    class="bg-white rounded"
    :class="sizeClass"
    @search="($event) => $emit('search', $event)"
  >
    <template
      v-for="(_, name) in $scopedSlots"
      v-slot:[name]="data"
    >
      <slot
        :name="name"
        v-bind="data"
      />
    </template>
  </vue-select>
</template>

<script>
import { VueSelect } from 'vue-select'
import { createPopper } from '@popperjs/core'
import 'vue-select/dist/vue-select.css';

export default {
  name: 'CInputSelect',

  components: {
    VueSelect,
  },

  props: {
    value: {
      type: [String, Array],
      default: () => '',
    },

    options: {
      type: Array,
      default: () => {
        return []
      },
    },

    clearable: {
      type: Boolean,
      default: true,
    },

    searchable: {
      type: Boolean,
      default: true,
    },

    appendToBody: {
      type: Boolean,
      default: true,
    },

    defaultValue: {
      type: [String, Array],
      default: () => '',
    },

    size: {
      type: String,
      default: 'md',
    },
  },

  computed: {
    _value: {
      get () {
        const fallbackValue = this.$attrs.multiple ? [] : ''
        return !!this.defaultValue && (this.value === this.defaultValue) ? fallbackValue : this.value
      },

      set (v) {
        this.$emit('input', !v ? this.defaultValue : v)
      }
    },

    sizeClass () {
      return this.size === 'sm' ? 'c-input-sm' : this.size === 'lg' ? 'c-input-lg' : ''
    },
  },

  methods: {
    calculateDropdownPosition (dropdownList, component, { width }) {
      /**
       * We need to explicitly define the dropdown width since
       * it is usually inherited from the parent with CSS.
       */
      dropdownList.style.width = width

      /**
       * Here we position the dropdownList relative to the $refs.toggle Element.
       *
       * The 'offset' modifier aligns the dropdown so that the $refs.toggle and
       * the dropdownList overlap by 1 pixel.
       *
       * The 'toggleClass' modifier adds a 'drop-up' class to the Vue Select
       * wrapper so that we can set some styles for when the dropdown is placed
       * above.
       */
      const popper = createPopper(component.$refs.toggle, dropdownList, {
        placement: 'bottom',
        modifiers: [
          {
            name: 'offset',
            options: {
              offset: [0, -1],
            },
          },
          {
            name: 'toggleClass',
            enabled: true,
            phase: 'write',
            fn ({ state }) {
              component.$el.classList.toggle('drop-up', state.placement === 'top')
            },
          }],
      })

      /**
       * To prevent memory leaks Popper needs to be destroyed.
       * If you return function, it will be called just before dropdown is removed from DOM.
       */
      return () => popper.destroy()
    },
  },
}
</script>

<style lang="scss">
.v-select {
  min-width: auto;
  position: relative;
  -ms-flex: 1 1 auto;
  flex: 1 1 auto;
  margin-bottom: 0;
  font-size: .9rem !important;
  border-radius: 0.25rem;
  color: #495057 !important;

  .vs__selected-options {
    // do not allow growing
    width: 0;
    padding: 0;
  }

  .vs__selected {
    display: block;
    white-space: nowrap;
    text-overflow: ellipsis;
    max-width: 100%;
    overflow: hidden;
    border: 0;
  }

  .vs__search {
    font-size: .9rem;
    border: 0;
    padding: 0 2px;
    padding-top: 0.375rem;
    margin: 0;
  }

  &:not(.vs--open) .vs__selected + .vs__search {
    // force this to not use any space
    // we still need it to be rendered for the focus
    width: 0;
    padding: 0;
    margin: 0;
    border: none;
    height: 0;
  }

  .vs__dropdown-toggle {
    min-height: calc(1.5em + 0.75rem + 4px);
    padding: 0.375rem calc(0.75rem - 2px);
    padding-top: 0 !important;
    border-width: 2px;
    border-color: var(--light);

    .vs__selected {
      margin-top: 0.375rem;
      padding: 0 2px;
    }

    .vs__actions {
      padding-top: 0.375rem;
      padding-right: 0;
    }
  }

  .vs__clear,
  .vs__open-indicator {
    fill: var(--gray-900);
    display: inline-flex;
  }

  .vs__clear {
    padding: 0;
    border: 0;
    background-color: transparent;
    cursor: pointer;
    margin-right: 8px
  }

  &.vs--single {
    .vs__selected {
      margin-left: 0;
      margin-right: 0;
    }
  }
}

.vs__spinner, .vs__spinner::after {
  width: 4em;
  height: 4em;
}

.vs__dropdown-menu {
  z-index: 1090;
}

.c-input-sm {
  font-size: 0.7875rem !important;

  .vs__search {
    font-size: 0.7875rem;
    padding-top: 0.25rem;
  }

  .vs__dropdown-toggle {
    min-height: calc(1.5em + 0.5rem + 4px);
    padding: 0.25rem calc(0.5rem - 2px);

    .vs__selected {
      margin-top: 0.25rem;
    }

    .vs__actions {
      padding-top: 0.25rem;
    }
  }
}
.c-input-lg {
  font-size: 1.125rem !important;

  .vs__search {
    font-size: 1.125rem;
    padding-top: .5rem;
  }

  .vs__dropdown-toggle {
    min-height: calc(1.5em + 1rem + 4px);
    padding: .5rem calc(1rem - 2px);

    .vs__selected {
      margin-top: .5rem;
    }

    .vs__actions {
      padding-top: .5rem;
    }
  }
}

</style>
