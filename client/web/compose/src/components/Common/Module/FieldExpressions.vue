<template>
  <b-form-group
    class="p-0 m-0"
  >
    <b-form-row
      v-for="(expr, e) in value"
      :key="e"
      class="mb-2"
      no-gutters
    >
      <b-input-group>
        <b-input-group-prepend>
          <b-button
            v-b-tooltip.hover="{ title: $t('validators.expression.tooltip'), container: '#body' }"
            variant="dark"
          >
            ƒ
          </b-button>
        </b-input-group-prepend>
        <slot :value="value[e]">
          <b-form-input
            v-model="value[e]"
            :placeholder="placeholder"
          />
        </slot>
        <b-input-group-addon
          class="m-1"
        >
          <!-- no prompt/confirmation on empty input -->
          <c-input-confirm
            :no-prompt="value[e].length === 0"
            show-icon
            @confirmed="$emit('remove', e)"
          />
        </b-input-group-addon>
      </b-input-group>
    </b-form-row>
  </b-form-group>
</template>
<script>

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  props: {
    value: {
      type: Array,
      default: () => ([]),
    },

    placeholder: {
      type: String,
      default: () => {},
    },
  },
}
</script>
