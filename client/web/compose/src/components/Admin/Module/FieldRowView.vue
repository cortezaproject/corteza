<template>
  <tr>
    <td />
    <td class="pl-3">
      {{ field.name }}
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
