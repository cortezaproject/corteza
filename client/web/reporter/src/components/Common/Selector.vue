<template>
  <b-container fluid>
    <b-row>
      <b-col
        cols="3"
      >
        <b-list-group>
          <b-list-group-item
            v-for="element in items"
            :key="element.kind"
            button
            @mouseover="currentValue = element.value"
            @click="$emit('select', element.kind)"
          >
            {{ element.label }}
          </b-list-group-item>
        </b-list-group>
      </b-col>
      <b-col
        v-if="currentValue"
        cols="9"
        :class="{ 'my-auto': displayMode === 'image' }"
      >
        <b-img
          v-if="displayMode === 'image'"
          fluid
          thumbnail
          :src="currentValue"
        />

        <p
          v-else-if="displayMode === 'text'"
        >
          {{ currentValue }}
        </p>
      </b-col>
    </b-row>
  </b-container>
</template>
<script>
export default {
  props: {
    // { label, kind, value }
    items: {
      type: Array,
      required: true,
    },

    // image or text
    displayMode: {
      type: String,
      default: 'image',
    },
  },

  data () {
    return {
      currentValue: undefined,
    }
  },

  watch: {
    items: {
      immediate: true,
      handler (items = []) {
        const { value } = items[0]

        if (value) {
          this.currentValue = value
        }
      },
    },
  },
}
</script>
