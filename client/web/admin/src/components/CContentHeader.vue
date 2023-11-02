<template>
  <div
    class="d-flex justify-content-end flex-wrap gap-1 mb-2 flex-fill-child"
  >
    <portal to="topbar-title">
      {{ title }}
    </portal>
    <slot />
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
