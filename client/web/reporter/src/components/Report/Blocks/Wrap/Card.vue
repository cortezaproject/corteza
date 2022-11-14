<template>
  <div class="h-100 p-2 overflow-hidden">
    <b-card
      no-body
      class="h-100 border-0 shadow-sm overflow-hidden"
    >
      <b-card-header
        class="sticky-top border-0 border-bottom p-0"
        header-bg-variant="white"
      >
        <slot
          name="header"
        />
      </b-card-header>

      <b-card-body
        v-if="bodySet"
        class="p-0 h-100"
      >
        <slot
          name="default"
        />
      </b-card-body>
    </b-card>
  </div>
</template>

<script>
export default {
  computed: {
    bodySet () {
      return !!this.$scopedSlots.default
    },
  },
}
</script>
