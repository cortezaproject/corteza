<template>
  <div
    class="d-flex align-items-center"
  >
    <font-awesome-icon
      v-if="!disabled && !disabledSorting && !hideIcons"
      :icon="['fas', 'grip-vertical']"
      :class="{
        'text-muted': disabledDragging,
      }"
      class="align-baseline mr-3 text-primary"
    />
    <b
      class="text-truncate"
    >
      <slot
        v-bind="item"
      >
        {{ item[textField] }}
      </slot>
    </b>
    <b-button
      v-if="_hideIcons"
      :data-test-id="`button-${selected ? 'unselect' : 'select'}`"
      variant="link"
      class="text-decoration-none d-flex align-items-center align-baseline ml-auto px-2"
    >
      <font-awesome-icon
        :icon="[selected ? 'far' : 'fas', selected ? 'eye' : 'eye-slash']"
        class="text-muted"
        @click="$emit(selected ? 'unselect' : 'select')"
      />
    </b-button>
  </div>
</template>

<script>
export default {
  name: 'CItemPickerItem',

  props: {
    item: {
      type: Object,
      required: true,
    },

    textField: {
      type: String,
      default: 'text',
    },

    selected: {
      type: Boolean,
    },

    disabled: {
      type: Boolean,
    },

    disabledDragging: {
      type: Boolean,
    },

    disabledSorting: {
      type: Boolean,
    },

    hideIcons: {
      type: Boolean,
    },
  },

  computed: {
    _hideIcons () {
      return !this.disabled && !this.hideIcons
    },
  },
}
</script>
