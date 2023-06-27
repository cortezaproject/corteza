<template>
  <b-form-group
    label-class="d-flex align-items-center text-primary p-0"
    :state="state"
    :class="formGroupStyleClasses"
  >
    <template
      v-if="!valueOnly"
      #label
    >
      <span class="d-inline-block text-truncate mw-100 py-1">
        {{ label }}
      </span>

      <hint
        :id="field.fieldID"
        :text="hint"
      />
    </template>

    <div
      class="small text-muted"
      :class="{ 'mb-1': description }"
    >
      {{ description }}
    </div>

    <multi
      v-if="field.isMulti"
      v-slot="ctx"
      :value.sync="value"
      :errors="errors"
    >
      <b-form-input
        v-model="value[ctx.index]"
        type="email"
        class="mr-2"
      />
    </multi>

    <template
      v-else
    >
      <b-form-input
        v-model="value"
        type="email"
        :state="state"
      />
      <errors :errors="errors" />
    </template>
  </b-form-group>
</template>
<script>
import base from './base'

export default {
  extends: base,
}
</script>
