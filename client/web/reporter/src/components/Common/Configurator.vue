<template>
  <b-container
    fluid
    class="p-0"
  >
    <b-row>
      <b-col
        cols="3"
      >
        <b-list-group>
          <draggable
            :list.sync="items"
            :disabled="!draggable"
            group="items"
            handle=".grab"
          >
            <b-list-group-item
              v-for="(item, index) in items"
              :key="index"
              button
              :active="currentIndex ? currentIndex === index : index === 0"
              class="item d-flex align-items-center justify-content-between"
              @click="$emit('select', index)"
            >
              <slot
                name="label"
                :item="item"
              />
            </b-list-group-item>
          </draggable>

          <b-list-group-item
            button
            class="text-primary rounded-top"
            :class="{ 'border-top-0': items.length }"
            @click="$emit('add')"
          >
            <font-awesome-icon
              :icon="['fas', 'plus']"
              size="sm"
              class="mr-1"
            />
            {{ $t('general:label.add') }}
          </b-list-group-item>
        </b-list-group>
      </b-col>
      <b-col
        v-if="currentIndex !== undefined"
        cols="9"
      >
        <slot
          name="configurator"
        />

        <c-input-confirm
          variant="danger"
          size="lg"
          size-confirm="lg"
          :borderless="false"
          :text="$t('general:label.delete')"
          class="d-flex"
          @confirmed="$emit('delete')"
        />
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import Draggable from 'vuedraggable'

export default {
  components: {
    Draggable,
  },

  props: {
    items: {
      type: Array,
      required: true,
    },

    currentIndex: {
      type: Number,
      default: undefined,
    },

    draggable: {
      type: Boolean,
      default: false,
    },
  },
}
</script>

<style lang="scss">
.item {
  .grab {
    opacity: 0;
  }

  &:hover {
    .grab {
      opacity: 1;
    }
  }
}
</style>
