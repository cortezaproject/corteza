<template>
  <b-form-group
    label-class="text-primary"
    :state="state"
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
