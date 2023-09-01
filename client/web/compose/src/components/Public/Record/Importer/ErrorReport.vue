<template>
  <b-card
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-form-group
      :label="$t('recordList.import.report.title')"
      label-class="text-primary"
    >
      <div
        class="small pl-2"
      >
        <b>{{ $t('recordList.import.report.startedAt') }}</b>: {{ datify(progress.startedAt) }}
      </div>

      <div
        class="small pl-2"
      >
        <b>{{ $t('recordList.import.report.finishedAt') }}</b>: {{ datify(progress.finishedAt) }}
      </div>

      <div
        class="small pl-2"
      >
        <b>{{ $t('recordList.import.report.totalRecords') }}</b>: {{ progress.entryCount }}
      </div>

      <div
        class="small pl-2"
      >
        <b>{{ $t('recordList.import.report.importedRecords') }}</b>: <span class="text-success">{{ progress.completed }}</span>
      </div>

      <div
        class="small pl-2"
      >
        <b>{{ $t('recordList.import.report.failedRecords') }}</b>: <span class="text-danger">{{ progress.failed }}</span>
      </div>
    </b-form-group>

    <b-form-group
      :label="$t('recordList.import.report.detectedErrors')"
      label-class="text-primary"
    >
      <b-table
        id="error-list"
        hover
        outlined
        fixed
        class="mb-0"
        head-variant="light"
        :items="errorList"
        :fields="fields"
      />
    </b-form-group>

    <b-form-group
      :label="$t('recordList.import.report.failedEntries')"
      label-class="text-primary"
    >
      <div
        v-for="(ee, ix) in progress.failLog.records"
        :key="ix"
        class="small pl-2"
      >
        <template v-if="ee.length == 1 || ee[0] === ee[1]">
          <b>{{ $t('recordList.import.report.failedEntriesLine') }}</b>: {{ ee[0] }}
        </template>
        <template v-else>
          <span>
            <b>{{ $t('recordList.import.report.failedEntriesLines') }}</b>:
            {{ ee[0] }}
            <font-awesome-icon
              :icon="['fas', 'arrow-right']"
              size="sm"
              class="mx-1"
            />
            {{ ee[1] }}
          </span>
        </template>
      </div>
    </b-form-group>

    <div slot="footer">
      <b-button
        variant="dark"
        class="float-right"
        @click="$emit('close')"
      >
        {{ $t('general.label.ok') }}
      </b-button>
    </div>
  </b-card>
</template>

<script>
import { fmt } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  props: {
    session: {
      type: Object,
      required: true,
      default: () => ({}),
    },

    noPool: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      expandedLogs: {},
      completedLogs: {},
    }
  },

  computed: {
    progress () {
      return this.session.progress
    },

    errorList () {
      return Object.entries(this.progress.failLog.errors)
        .map(([k, v]) => ({ k, v }))
    },

    fields () {
      return [
        {
          key: 'k',
          label: this.$t('recordList.import.report.error'),
          tdClass: 'border-top text-truncate pointer',
        },
        {
          key: 'v',
          label: this.$t('recordList.import.report.count'),
          tdClass: 'border-top text-truncate pointer',
        },
      ]
    },
  },

  methods: {
    datify (dt) {
      return fmt.fullDateTime(dt)
    },

    stringify (r) {
      return r.values.map(v => `${v.name}: ${v.value}`).join(', ')
    },
  },
}
</script>

<style lang="scss">
.progress-label {
  font-size: 15px;
}

.fit {
  white-space: nowrap;
  width: 15%;
}

.pointer {
  cursor: pointer;
}
</style>
