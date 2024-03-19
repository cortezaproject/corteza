<template>
  <div class="list-background w-100 p-3 rounded border border-light">
    <slot></slot>

    <b-button
      v-if="!hideAddButton"
      variant="primary"
      size="sm"
      :data-test-id="testID"
      :class="addButtonClass"
      :disabled="disableAddButton"
      @click="$emit('add-item')"
    >
       <font-awesome-icon
         :icon="['fas', 'plus']"
         class="mr-1"
       />

      {{ labels.addButton || '' }}
    </b-button>
  </div>
</template>

<script>
export default {
  name: 'CFormTableWrapper',

  props: {
    labels: {
      type: Object,
      default: () => {},
    },

    hideAddButton: {
      type: Boolean,
      default: false,
    },

    disableAddButton: {
      type: Boolean,
      default: false,
    },

    addButtonClass: {
      type: String,
      default: '',
    },

    buttonTestId : {
      type: String,
      default: '',
    }
  },
}
</script>

<style scoped>
.list-background {
  background-color: var(--body-bg);
}
</style>
