<template>
  <b-card no-body>
    <slot name="title" />

    <fieldset>
      <b-form-group
        :label="$t('metric.editStyle.color')"
        label-cols="4"
      >
        <b-form-input
          v-model="options.color"
          type="color"
          debounce="300"
          class="mb-1"
        />
      </b-form-group>

      <b-form-group
        :label="$t('metric.editStyle.backgroundColor')"
        label-cols="4"
      >
        <b-form-input
          v-model="options.backgroundColor"
          type="color"
          debounce="300"
          class="mb-1"
        />
      </b-form-group>

      <b-form-group
        :label="$t('metric.editStyle.fontSize')"
        label-cols="4"
      >
        <b-form-input
          v-model="options.fontSize"
          type="number"
          placeholder="16"
          class="mb-1"
        />
      </b-form-group>
    </fieldset>
  </b-card>
</template>

<script>
export default {
  i18nOptions: {
    namespaces: 'block',
  },

  props: {
    options: {
      type: Object,
      required: true,
      default: () => ({}),
    },
  },
}
</script>
