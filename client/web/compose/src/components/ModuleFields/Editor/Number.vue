<template>
  <b-form-group
    label-class="d-flex align-items-center text-primary p-0"
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
      <b-input-group
        :prepend="field.options.prefix"
        :append="field.options.suffix"
      >
        <b-form-input
          v-model="value[ctx.index]"
          autocomplete="off"
          type="number"
          number
          class="mr-2"
        />
      </b-input-group>
    </multi>

    <b-input-group
      v-else
      :prepend="field.options.prefix"
      :append="field.options.suffix"
    >
      <b-form-input
        v-model="value"
        autocomplete="off"
        type="number"
        number
      />
    </b-input-group>
    <errors :errors="errors" />
  </b-form-group>
</template>
<script>
import base from './base'

export default {
  extends: base,
}
</script>
