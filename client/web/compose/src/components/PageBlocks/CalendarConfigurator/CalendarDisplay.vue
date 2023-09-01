<template>
  <fieldset class="form-group">
    <b-form-group
      :label="$t('calendar.calendarHeader')"
      horizontal
      label-class="text-primary"
    >
      <b-form-checkbox
        v-model="options.header.hide"
      >
        {{ $t('calendar.hideHeader') }}
      </b-form-checkbox>
      <b-form-checkbox
        v-model="options.header.hidePrevNext"
        :disabled="options.header.hide"
      >
        {{ $t('calendar.hideNavigation') }}
      </b-form-checkbox>
      <b-form-checkbox
        v-model="options.header.hideToday"
        :disabled="options.header.hide"
      >
        {{ $t('calendar.hideToday') }}
      </b-form-checkbox>
      <b-form-checkbox
        v-model="options.header.hideTitle"
        :disabled="options.header.hide"
      >
        {{ $t('calendar.hideTitle') }}
      </b-form-checkbox>
    </b-form-group>
    <b-form-group
      :label="$t('calendar.view.enabled')"
      horizontal
      label-class="text-primary"
    >
      <b-form-checkbox-group
        v-model="options.header.views"
        :disabled="options.header.hide"
        buttons
        button-variant="outline-secondary"
        size="sm"
        name="buttons2"
        :options="views"
      />
    </b-form-group>

    <b-form-group
      :description="$t('calendar.view.footnote')"
      :label="$t('calendar.view.default')"
      horizontal
      label-class="text-primary"
    >
      <b-form-radio-group
        v-model="options.defaultView"
        buttons
        button-variant="outline-secondary"
        size="sm"
        name="buttons2"
        :options="views"
      />
    </b-form-group>

    <b-row>
      <b-col
        cols="12"
        md="6"
      >
        <b-form-group
          :label="$t('calendar.view.onEventClick')"
          label-class="text-primary"
        >
          <b-form-select
            v-model="options.eventDisplayOption"
            :options="eventDisplayOptions"
          />
        </b-form-group>
      </b-col>
    </b-row>
  </fieldset>
</template>
<script>
import base from '../base'
import { compose } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  extends: base,

  computed: {
    views () {
      return compose.PageBlockCalendar.availableViews()
        .map(view => ({ value: view, text: this.$t(`calendar.view.${view}`) }))
    },

    eventDisplayOptions () {
      return [
        { value: 'sameTab', text: this.$t('calendar.view.openInSameTab') },
        { value: 'newTab', text: this.$t('calendar.view.openInNewTab') },
        { value: 'modal', text: this.$t('calendar.view.openInModal') },
      ]
    },
  },
}
</script>
