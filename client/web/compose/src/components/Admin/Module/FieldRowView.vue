<template>
  <tr>
    <td />
    <td class="pl-3">
      {{ field.name }}
      <b-badge
        v-if="field.omit"
        variant="info"
        class="ml-2 align-middle"
      >
        {{ $t('unavailable') }}
      </b-badge>
    </td>
    <td class="pl-3">
      {{ field.label }}
    </td>
    <td class="pl-3">
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
