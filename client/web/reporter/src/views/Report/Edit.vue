<template>
  <div class="py-3">
    <portal to="topbar-title">
      <div
        class="d-flex w-100"
      >
        {{ pageTitle }}
      </div>
    </portal>

    <portal to="topbar-tools">
      <b-button-group
        v-if="report && !isNew"
        size="sm"
        class="mr-1"
      >
        <b-button
          variant="primary"
          data-test-id="button-report-builder"
          class="d-flex align-items-center justify-content-center"
          :disabled="!report.canUpdateReport"
          :to="reportBuilder"
        >
          {{ $t('report.builder') }}
          <font-awesome-icon
            class="ml-2"
            :icon="['fas', 'tools']"
          />
        </b-button>
        <b-button
          v-b-tooltip.hover="{ title: $t('tooltip.view-report'), container: '#body' }"
          variant="primary"
          style="margin-left:2px;"
          :disabled="!canRead"
          :to="reportViewer"
        >
          <font-awesome-icon
            :icon="['far', 'eye']"
          />
        </b-button>
      </b-button-group>
    </portal>

    <b-container v-if="report">
      <b-row no-gutters>
        <b-col>
          <b-card
            no-body
            show
            class="shadow-sm"
          >
            <b-card-header
              v-if="!isNew"
              header-bg-variant="white border-bottom"
              class="py-3"
            >
              <b-row
                no-gutters
                class="align-items-center"
              >
                <div>
                  <b-button
                    v-if="canCreate"
                    data-test-id="button-create-report"
                    variant="primary"
                    size="lg"
                    class="mr-1"
                    :to="{ name: 'report.create' }"
                  >
                    {{ $t('new-report') }}
                  </b-button>

                  <c-permissions-button
                    v-if="canGrant"
                    :title="report.meta.name || report.handle || report.reportID"
                    :target="report.meta.name || report.handle || report.reportID"
                    :resource="`corteza::system:report/${report.reportID}`"
                    :button-label="$t('permissions')"
                    class="btn-lg ml-1"
                  />
                </div>
              </b-row>
            </b-card-header>

            <b-container
              fluid
              class="px-4 pt-3"
            >
              <b-row>
                <b-col
                  cols="12"
                  md="6"
                  xl="4"
                >
                  <b-form-group
                    :label="$t('name-with-star')"
                    label-class="text-primary"
                  >
                    <b-form-input
                      v-model="report.meta.name"
                      data-test-id="input-name"
                      :placeholder="$t('name')"
                      required
                      :state="nameState"
                      @input="handleDetectStateChange"
                    />
                  </b-form-group>
                </b-col>
                <b-col
                  cols="12"
                  md="6"
                  xl="4"
                >
                  <b-form-group
                    :label="$t('handle-with-star')"
                    label-class="text-primary"
                  >
                    <b-form-input
                      v-model="report.handle"
                      data-test-id="input-handle"
                      :placeholder="$t('placeholder-handle')"
                      required
                      :state="handleState"
                      @input="handleDetectStateChange"
                    />
                    <b-form-invalid-feedback
                      data-test-id="input-handle-invalid-state"
                      :state="handleState"
                    >
                      {{ $t('invalid-handle-characters') }}
                    </b-form-invalid-feedback>
                  </b-form-group>
                </b-col>
              </b-row>

              <b-form-group
                :label="$t('description')"
                label-class="text-primary"
              >
                <b-form-textarea
                  v-model="report.meta.description"
                  data-test-id="input-description"
                  :placeholder="$t('report.description')"
                  rows="5"
                  @input="handleDetectStateChange"
                />
              </b-form-group>

              <!-- <b-form-group
                label="Tags"
                class="text-primary"
              >
                <b-form-tags
                  v-model="report.meta.tags"
                />
              </b-form-group> -->
            </b-container>
          </b-card>
        </b-col>
      </b-row>
    </b-container>

    <portal to="report-toolbar">
      <editor-toolbar
        :back-link="{ name: 'report.list' }"
        :hide-delete="isNew"
        :delete-disabled="!canDelete"
        :save-disabled="!canSave"
        :processing="processing"
        :processing-save="processingSave"
        :processing-delete="processingDelete"
        @delete="handleDelete"
        @save="handleSave"
      />
    </portal>
  </div>
</template>

<script>
import { system } from '@cortezaproject/corteza-js'
import { handle } from '@cortezaproject/corteza-vue'
import report from 'corteza-webapp-reporter/src/mixins/report'
import EditorToolbar from 'corteza-webapp-reporter/src/components/EditorToolbar'
import { mapGetters } from 'vuex'
import { isEqual } from 'lodash'

export default {
  name: 'EditReport',

  i18nOptions: {
    namespaces: 'edit',
  },

  components: {
    EditorToolbar,
  },

  mixins: [
    report,
  ],

  data () {
    return {
      processing: false,

      report: undefined,
      initialReportState: undefined,

      detectStateChange: false,
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canGrant () {
      return this.can('system/', 'grant')
    },

    canCreate () {
      return this.can('system/', 'report.create')
    },

    canRead () {
      return this.report ? this.report.canReadReport : false
    },

    canDelete () {
      return this.report ? this.report.canDeleteReport : false
    },

    canUpdate () {
      return this.isNew ? this.canCreate : (this.report && this.report.canUpdateReport) || false
    },

    reportID () {
      return this.$route.params.reportID
    },

    isNew () {
      return !this.reportID
    },

    pageTitle () {
      return this.isNew ? this.$t('report.create') : this.$t('report.edit')
    },

    reportBuilder () {
      return this.report ? { name: 'report.builder', params: { reportID: this.report.reportID } } : undefined
    },

    reportViewer () {
      return this.report ? { name: 'report.view', params: { reportID: this.report.reportID } } : undefined
    },

    nameState () {
      const { name = '' } = this.report.meta
      return name.length ? null : false
    },

    handleState () {
      return handle.handleState(this.report.handle)
    },

    canSave () {
      return this.canUpdate && ![this.nameState, this.handleState].includes(false)
    },
  },

  watch: {
    reportID: {
      immediate: true,
      handler (reportID) {
        // Fetch report or make a new one
        if (reportID) {
          this.fetchReport(reportID)
        } else {
          this.report = new system.Report()
          this.initialReportState = new system.Report()
        }
      },
    },
  },

  beforeRouteUpdate (to, from, next) {
    this.checkUnsavedChart(next)
  },

  beforeRouteLeave (to, from, next) {
    this.checkUnsavedChart(next)
  },

  methods: {
    handleDetectStateChange () {
      this.detectStateChange = true
    },

    checkUnsavedChart (next) {
      if (this.report.deletedAt) {
        return next(true)
      }

      const reportState = {
        handle: this.report.handle,
        meta: {
          name: this.report.meta.name,
          description: this.report.meta.description,
        },
      }

      const initialReportState = {
        handle: this.initialReportState.handle,
        meta: {
          name: this.initialReportState.meta.name,
          description: this.initialReportState.meta.description,
        },
      }

      next(!isEqual(reportState, initialReportState) ? window.confirm(this.$t('unsavedChanges')) : true)
    },
  },
}
</script>
