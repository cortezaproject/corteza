<template>
  <div class="list-background w-100 p-3">
    <h5>{{ listTitle }}</h5>
    <slot></slot>

    <b-button
      v-if="useInternalButton"
      variant="primary"
      size="sm"
      @click="handleClick"
    >
      {{ actionTitle }}
    </b-button>
  </div>
</template>

<script>
export default {
  name: 'CInputList',

  props: {
    listTitle: {
      type: String,
      default: '',
    },

    actionTitle: {
      type: String,
      default: '',
    },

    useInternalButton: {
      type: Boolean,
      default: false,
    },
  },

  methods: {
    handleClick() {
      this.$emit('add-item');
    },
  },
};
</script>

<style scoped>
.list-background {
  background-color: var(--body-bg);
}
</style>
