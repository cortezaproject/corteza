<template>
  <div>
    <b-form-group
      horizontal
      :label-cols="3"
      breakpoint="md"
      :label="$t('calendar.reminderFeed.colorLabel')"
    >
      <b-input-group>
        <b-form-input
          v-model="feed.options.color"
          style="max-width: 50px;"
          type="color"
        />
      </b-input-group>
    </b-form-group>
  </div>
</template>

<script>
import base from './base'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  extends: base,
}
</script>
