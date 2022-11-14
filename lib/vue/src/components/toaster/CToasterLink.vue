<template>
  <div>
    <router-link
      :to="link"
    >
      {{ link.label }}
    </router-link>
  </div>
</template>

<script lang="js">
export default {
  props: {
    link: {
      type: Object,
      required: true,
    },
  },
}
</script>
