<template>
  <div
    class="d-flex align-items-center"
  >
    <font-awesome-icon
      v-if="!disabled && !disabledDragging && !disabledSorting && !hideIcons"
      :icon="['fas', 'grip-vertical']"
      class="handle align-baseline mr-3 text-primary"
    />
    <b
      class="cursor-default text-truncate"
    >
      <slot
        v-bind="item"
      >
        {{ item[textField] }}
      </slot>
    </b>
    <template
      v-if="_hideIcons"
    >
      <font-awesome-icon
        v-if="selected"
        :icon="['far', 'eye']"
        class="align-baseline ml-auto text-muted pointer"
        @click="$emit('unselect')"
      />
      <font-awesome-icon
        v-else
        :icon="['fas', 'eye-slash']"
        class="align-baseline ml-auto text-muted pointer"
        @click="$emit('select')"
      />
    </template>
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
<style scoped>
.handle {
  cursor: grab;
}

.handle:active {
  cursor: grabbing;
}

.cursor-default {
  cursor: default;
}
</style>
