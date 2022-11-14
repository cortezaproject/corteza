<template>
  <div
    class="d-flex justify-content-end mb-2"
  >
    <portal to="topbar-title">
      {{ title }}
    </portal>

    <b-button-toolbar>
      <slot />
    </b-button-toolbar>
  </div>
</template>

<script>
export default {
  name: 'CContentHeader',

  props: {
    title: {
      type: String,
      required: false,
      default: '',
    },
  },
}
</script>
