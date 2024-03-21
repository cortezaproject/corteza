<template>
  <tr>
    <td />

    <td>
      {{ field.name }}
      <b-badge
        v-if="field.omit"
        variant="info"
        class="ml-2 align-middle"
      >
        {{ $t('unavailable') }}
      </b-badge>
    </td>

    <td>
      {{ field.label }}
    </td>

    <td>
      {{ field.kind }}
    </td>
  </tr>
</template>

<script>

export default {
  props: {
    field: {
      type: Object,
      required: true,
    },
  },
}
</script>
