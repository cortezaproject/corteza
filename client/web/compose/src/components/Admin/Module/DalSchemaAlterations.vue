<template>
  <b-modal
    v-model="showModal"
    :title="$t('title')"
    size="xl"
    body-class="p-0 border-top-0 position-relative"
    header-class="p-3 pb-0 border-bottom-0"
    no-fade
    v-on="$listeners"
  >
    <b-table-simple
      borderless
      sticky-header
      responsive
      hover
      head-variant="light"
      class="mb-0"
      style="min-height: 300px; max-height: 75vh;"
    >
      <b-thead>
        <b-tr>
          <b-th>
            {{ $t('columns.alteration') }}
          </b-th>

          <b-th style="max-width: 300px;">
            {{ $t('columns.change') }}
          </b-th>

          <b-th
            class="text-center"
          >
            {{ $t('columns.status') }}
          </b-th>

          <b-th style="min-width: 200px;" />
        </b-tr>
      </b-thead>

      <b-tbody v-if="sortedAlterations.length && !loading">
        <b-tr
          v-for="a in sortedAlterations"
          :key="a.alterationID"
          class="border-top"
          :class="{ 'bg-extra-light': a.alterationID === dependOnHover }"
          @mouseover="dependOnHover = a.dependsOn"
          @mouseleave="dependOnHover = undefined"
        >
          <b-td>
            {{ a.alterationID }}
          </b-td>

          <b-td>
            {{ stringifyParams(a.params) }}
          </b-td>

          <b-td class="text-center align-top">
            <b-badge
              v-if="a.error"
              variant="danger"
            >
              {{ a.error || '' }}
            </b-badge>

            <b-badge
              v-else-if="a.completedAt"
              variant="success"
            >
              {{ $t('resolved') }}
            </b-badge>

            <b-badge
              v-else-if="a.dependsOn"
              variant="extra-light"
            >
              {{ $t('waitingFor', { id: a.dependsOn }) }}
            </b-badge>
          </b-td>

          <b-td class="text-right">
            <b-spinner
              v-if="a.processing"
              variant="primary"
              small
            />

            <template v-else>
              <c-input-confirm
                v-if="!a.completedAt"
                :disabled="!canResolve(a) || a.processing || processing"
                variant="primary"
                size="sm"
                class="mx-1"
                @click.stop
                @confirmed="onResolve(a)"
              >
                {{ $t('resolve') }}
              </c-input-confirm>

              <c-input-confirm
                v-if="!a.completedAt"
                :disabled="!canDismiss(a) || a.processing || processing"
                variant="light"
                size="sm"
                class="mx-1"
                @click.stop
                @confirmed="onDismiss(a)"
              >
                {{ $t('dismiss') }}
              </c-input-confirm>
            </template>
          </b-td>
        </b-tr>
      </b-tbody>

      <div
        v-else
        class="position-absolute text-center mt-5 d-print-none"
        style="left: 0; right: 0;"
      >
        <b-spinner
          v-if="loading"
        />

        <p
          v-else-if="!sortedAlterations.length"
          class="mb-0 mx-2"
        >
          {{ $t('noAlterations') }}
        </p>
      </div>
    </b-table-simple>

    <template #modal-footer>
      <b-button
        :disabled="processing"
        class="text-primary border-0"
        variant="light"
        @click="showModal = false"
      >
        {{ canResolveAlterations ? $t('general:label.cancel') : $t('general:label.close') }}
      </b-button>

      <c-input-confirm
        v-if="canResolveAlterations"
        variant="primary"
        :disabled="processing"
        size="md"
        @confirmed="onResolve()"
      >
        {{ $t('resolveAuto') }}
        <b-spinner
          v-if="processing"
          variant="white"
          small
        />
      </c-input-confirm>
    </template>
  </b-modal>
</template>
<script>
import { compose } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'module',
    keyPrefix: 'edit.schemaAlterations',
  },

  props: {
    modal: {
      type: Boolean,
      required: false,
    },

    module: {
      type: compose.Module,
      required: true,
    },

    batch: {
      type: Array,
      required: false,
      default: undefined,
    },
  },

  data () {
    return {
      showModal: false,

      loading: false,
      processing: false,

      dependOnHover: undefined,

      alterations: [],

      alterationProcessing: {},
    }
  },

  computed: {
    sortedAlterations () {
      return this.alterations.toSorted((a, b) => (a.batchID || '').localeCompare(b.batchID || '') || (a.dependsOn || '').localeCompare(b.dependsOn || ''))
    },

    canResolveAlterations () {
      return this.sortedAlterations.some(this.canResolve)
    },
  },

  watch: {
    batch: {
      handler (batch) {
        if (batch && batch.length) {
          this.load(...batch)
        }
      },
    },
  },

  methods: {
    async onDismiss (alteration) {
      this.processing = true

      alteration = alteration ? [alteration] : this.alterations

      const alterationID = []
      for (const a of alteration) {
        alterationID.push(a.alterationID)
        a.processing = true
      }

      this.$SystemAPI.dalSchemaAlterationDismiss({ alterationID }).then(() => {
        this.toastSuccess(this.$t('notification:module.schemaAlterations.dismiss.success'))
      }).catch(this.toastErrorHandler(this.$t('notification:module.schemaAlterations.dismiss.error')))
        .finally(() => {
          for (const a of alteration) {
            a.processing = false
          }
          this.load(...this.batch)
          this.processing = false
        })
    },

    async onResolve (alteration) {
      this.processing = true

      alteration = alteration ? [alteration] : this.alterations

      const alterationID = []
      for (const a of alteration) {
        alterationID.push(a.alterationID)
        a.processing = true
      }

      this.$SystemAPI.dalSchemaAlterationApply({ alterationID }).then(() => {
        this.toastSuccess(this.$t('notification:module.schemaAlterations.resolve.success'))
      }).catch(this.toastErrorHandler(this.$t('notification:module.schemaAlterations.resolve.error')))
        .finally(() => {
          for (const a of alteration) {
            a.processing = false
          }
          this.load(...this.batch)
          this.processing = false
        })
    },

    async load (...batchID) {
      if (!batchID || (batchID && !batchID.length)) {
        this.alterations = []
        return
      }

      this.loading = true

      return this.$SystemAPI.dalSchemaAlterationList({ batchID }).then(({ set }) => {
        this.alterations = set

        if (this.alterations.length) {
          this.showModal = true
        }
      }).catch(this.toastErrorHandler(this.$t('notification:module.schemaAlterations.load.error')))
        .finally(() => {
          this.loading = false
        })
    },

    stringifyParams (params) {
      switch (true) {
        case !!params.attributeAdd:
          return this.stringifyAttributeAddParams(params.attributeAdd)

        case !!params.attributeDelete:
          return this.stringifyAttributeDeleteParams(params.attributeDelete)

        case !!params.attributeReType:
          return this.stringifyAttributeReTypeParams(params.attributeReType)

        case !!params.attributeReEncode:
          return this.stringifyAttributeReEncodeParams(params.attributeReEncode)

        case !!params.modelAdd:
          return this.stringifyModelAddParams(params.modelAdd)

        case !!params.modelDelete:
          return this.stringifyModelDeleteParams(params.modelDelete)
      }

      throw new Error('Unknown alteration type')
    },

    stringifyAttributeAddParams ({ attr = {} }) {
      return this.$t('params.attribute.add', { ident: attr.ident, storeType: attr.store.type, attrType: attr.type.type })
    },

    stringifyAttributeDeleteParams ({ attr = {} }) {
      return this.$t('params.attribute.delete', { ident: attr.ident, storeType: attr.store.type })
    },

    stringifyAttributeReTypeParams ({ attr = {}, to = {} }) {
      return this.$t('params.attribute.reType', { ident: attr.ident, toType: to.type })
    },

    stringifyAttributeReEncodeParams ({ attr = {}, to = {} }) {
      return this.$t('params.attribute.reEncode', { ident: attr.ident, toType: to.type })
    },

    stringifyModelAddParams ({ model = {} }) {
      return this.$t('params.model.add', { ident: model.ident })
    },

    stringifyModelDeleteParams ({ model = {} }) {
      return this.$t('params.model.delete', { ident: model.ident })
    },

    canDismiss (alteration) {
      if (alteration.completedAt) {
        return false
      }

      if (alteration.dependsOn) {
        return this.alterations.some(a => a.alterationID === alteration.dependsOn && !a.completedAt)
      }

      return true
    },

    canResolve (alteration) {
      return this.canDismiss(alteration)
    },
  },
}
</script>
