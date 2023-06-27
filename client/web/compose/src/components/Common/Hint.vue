<template>
  <div
    v-if="text"
    :id="id"
    class="d-flex align-items-center ml-2"
  >
    <font-awesome-icon
      :icon="['far', 'question-circle']"
    />

    <b-tooltip
      :target="id"
      triggers="hover"
      placement="right"
      variant="dark"
      custom-class="m-1"
    >
      {{ text }}
    </b-tooltip>
  </div>
</template>

<script>
export default {
  props: {
    text: {
      type: String,
      default: '',
    },
    id: {
      type: String,
      required: true,
    },
  },
}
</script>
