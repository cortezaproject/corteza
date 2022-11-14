<template>
  <div>
    <div
      v-if="singleInput"
      class="mb-2"
    >
      <slot name="single" />
    </div>

    <errors
      v-if="errors"
      :errors="errors"
      class="mb-1"
    />

    <draggable
      :list.sync="val"
      handle=".handle"
    >
      <div
        v-for="(v, index) of val"
        :key="index"
        class="d-flex w-100 align-items-center mb-1"
      >
        <font-awesome-icon
          v-b-tooltip.hover
          :icon="['fas', 'bars']"
          :title="$t('tooltip.dragAndDrop')"
          class="handle text-light ml-1 mr-2"
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
          class="pointer text-danger ml-2 mr-1"
          @click="removeValue(index)"
        />
      </div>
    </draggable>

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
