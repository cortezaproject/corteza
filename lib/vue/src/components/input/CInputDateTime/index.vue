<template>
  <div class="c-input-date-time d-flex flex-wrap w-100">
    <b-form-datepicker
      v-if="!noDate"
      v-model="date"
      data-test-id="picker-date"
      :placeholder="labels.none"
      :date-format-options="{ year: 'numeric', month: 'numeric', day: 'numeric' }"
      :min="minDate"
      :max="maxDate"
      :label-reset-button="labels.clear"
      :label-today-button="labels.today"
      label-help=""
      today-variant="info"
      selected-variant="secondary"
      boundary="window"
      hide-header
      reset-button
      today-button
      class="h-100 overflow-hidden"
    />

    <b-form-timepicker
      v-if="!noTime"
      v-model="time"
      data-test-id="picker-time"
      :placeholder="labels.none"
      :label-reset-button="labels.clear"
      :label-now-button="labels.now"
      boundary="window"
      hide-header
      no-close-button
      reset-button
      now-button
      class="h-100 overflow-hidden"
    />

    <slot />
  </div>
</template>
<script lang="js">
import { getDate, setDate, getTime, setTime } from './lib/index.ts'

export default {
  props: {
    value: {
      type: [String, Date],
      required: false,
    },

    noTime: {
      type: Boolean,
      default: false,
    },

    noDate: {
      type: Boolean,
      default: false,
    },

    onlyFuture: {
      type: Boolean,
      default: false,
    },

    onlyPast: {
      type: Boolean,
      default: false,
    },

    size: {
      type: String,
      default: 'md',
    },

    labels: {
      type: Object,
      required: true,
    },
  },

  computed: {
    date: {
      get () {
        return getDate(this.value)
      },

      set (date) {
        this.$emit('input', setDate(date, this.value, this.noDate, this.noTime))
      },
    },

    time: {
      get () {
        return getTime(this.value)
      },

      set (time) {
        this.$emit('input', setTime(time, this.value, this.noDate, this.noTime))
      },
    },

    minDate () {
      return this.onlyFuture ? new Date() : undefined
    },

    maxDate () {
      return this.onlyPast ? new Date() : undefined
    },
  },
}
</script>

<style lang="scss">
.c-input-date-time {
  min-width: 120px;

  .btn {
    padding: 0.25rem 0.5rem;
  }

  label {
    font-family: "Poppins-Regular";
    color: #495057 !important;
  }

  .b-form-datepicker, .b-form-timepicker {
    flex: 1 0 130px;
  }
}

.b-calendar-inner {
  background-color: white;
}
</style>
