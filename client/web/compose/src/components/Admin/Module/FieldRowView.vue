<template>
  <tr>
    <td class="pr-2" />

    <td style="width: 250px;">
      {{ field.name }}
      <b-badge
        v-if="field.omit"
        variant="info"
        class="ml-2 align-middle"
      >
        {{ $t('unavailable') }}
      </b-badge>
    </td>

    <td style="width: 250px;">
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
