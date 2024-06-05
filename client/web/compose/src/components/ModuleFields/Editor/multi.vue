<template>
  <div>
    <div
      v-if="singleInput"
      class="mb-2"
    >
      <slot name="single" />
    </div>

    <draggable
      v-if="showList"
      :list.sync="val"
      handle=".handle"
    >
      <div
        v-for="(v, index) of val"
        :key="index"
        class="d-flex w-100 align-items-center mb-1"
      >
        <font-awesome-icon
          v-b-tooltip.noninteractive.hover="{ title: $t('tooltip.dragAndDrop'), container: '#body' }"
          :icon="['fas', 'bars']"
          class="handle text-secondary mr-3"
        />

        <div
          class="flex-grow-1"
        >
          <slot
            :index="index"
          />
        </div>

        <font-awesome-icon
          v-if="removable"
          :icon="['fas', 'times']"
          class="pointer text-danger ml-3"
          @click="removeValue(index)"
        />
      </div>
    </draggable>

    <errors :errors="errors" />

    <b-button
      v-if="!singleInput"
      variant="primary"
      size="sm"
      :class="{ 'mt-2': val.length }"
      @click="val.push(defaultValue)"
    >
      + {{ $t('label.addValue') }}
    </b-button>
  </div>
</template>
<script>
import errors from '../errors'
import draggable from 'vuedraggable'
import { validator } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'general',
  },

  components: {
    draggable,
    errors,
  },

  props: {
    value: {
      type: Array,
      required: true,
      default: () => [],
    },

    removable: {
      type: Boolean,
      default: true,
    },

    singleInput: {
      type: Boolean,
      default: false,
    },

    showList: {
      type: Boolean,
      default: true,
    },

    errors: {
      type: validator.Validated,
      required: true,
    },

    defaultValue: {
      type: undefined,
      default: undefined,
    },
  },

  computed: {
    val: {
      get () {
        return this.value
      },

      set (val) {
        this.$emit('update:value', val)
      },
    },
  },

  methods: {
    removeValue (index) {
      if (index > -1) {
        this.val.splice(index, 1)
      }
    },
  },

}
</script>
<style lang="scss" scoped>
.handle {
  cursor: grab;
}

.pointer {
  cursor: pointer;
}
</style>
