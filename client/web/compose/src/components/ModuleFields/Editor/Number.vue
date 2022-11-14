<template>
  <b-form-group
    label-class="text-primary"
    :class="formGroupStyleClasses"
  >
    <template
      v-if="!valueOnly"
      #label
    >
      <div
        class="d-flex align-items-top"
      >
        <label
          class="mb-0"
        >
          {{ label }}
        </label>

        <hint
          :id="field.fieldID"
          :text="hint"
        />
      </div>
      <small
        class="form-text font-weight-light text-muted"
      >
        {{ description }}
      </small>
    </template>

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
