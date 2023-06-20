<template>
  <b-modal
    v-model="showModal"
    title="Schema alterations"
    size="xl"
    body-class="p-0 border-top-0"
    header-class="p-3 pb-0 border-bottom-0"
    @change="$emit('change', $event)"
  >
    <b-table-simple
      borderless
      sticky-header
      responsive
      head-variant="light"
      style="max-height: 75vh;"
    >
      <b-thead>
        <b-tr>
          <b-th
            class="text-primary"
          >
            Alteration
          </b-th>
          <b-th
            class="text-primary"
          >
            Operation
          </b-th>

          <b-th
            class="text-primary"
          >
            Change
          </b-th>

          <b-th
            class="text-primary text-center"
          >
            Status
          </b-th>
          <b-th style="min-width: 200px;" />
        </b-tr>
      </b-thead>

      <b-tbody>
        <b-tr
          v-for="a in sortedAlterations"
          :key="a.alterationID"
          class="border-top"
          :class="{ 'bg-gray': a.alterationID === dependOnHover }"
          @mouseover="dependOnHover = a.dependsOn"
          @mouseleave="dependOnHover = undefined"
        >
          <b-td>
            {{ a.alterationID }}
          </b-td>

          <b-td>
            {{ $t(`schema-alteration.${a.kind}.label`) }}
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
              Resolved
            </b-badge>

            <b-badge
              v-else-if="a.dependsOn"
              variant="secondary"
            >
              Waiting for {{ a.dependsOn }}
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
                Resolve
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
                Dismiss
              </c-input-confirm>
            </template>
          </b-td>
        </b-tr>
      </b-tbody>
    </b-table-simple>

    <template #modal-footer>
      <b-button
        variant="link"
        size="sm"
        :disabled="processing"
        @click="$emit('cancel')"
      >
        Cancel
      </b-button>

      <c-input-confirm
        variant="primary"
        :disabled="processing"
        :borderless="false"
        @confirmed="onResolve()"
      >
        Resolve automatically
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
      type: String,
      required: false,
      default: undefined,
    },
  },

  data () {
    return {
      showModal: false,

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
  },

  watch: {
    modal: {
      immediate: true,
      handler (show = false) {
        this.showModal = show
      },
    },

    batch: {
      immediate: true,
      handler (batch) {
        this.load(batch)
      },
    },
  },

  methods: {
    async onDismiss (alteration = undefined) {
      if (!alteration) {
        alteration = this.alterations
      } else {
        alteration = [alteration]
      }

      const ids = []
      for (const a of alteration) {
        ids.push(a.alterationID)
        this.alterationProcessing[a.alterationID] = true
      }

      this.processing = true

      this.$SystemAPI.dalSchemaAlterationDismiss({ alterationID: ids }).then(() => {
        this.toastSuccess(this.$t('notification:module.alteration.dismiss.success'))
      }).catch(this.toastErrorHandler(this.$t('notification:alteration.dismiss.error')))
        .finally(() => {
          this.processing = false
          for (const i of ids) {
            this.$delete(this.alterationProcessing, i)
          }
          this.load(this.batch)
        })
    },

    async onResolve (alteration = undefined) {
      if (!alteration) {
        alteration = this.alterations
      } else {
        alteration = [alteration]
      }

      const ids = []
      for (const a of alteration) {
        ids.push(a.alterationID)
        this.alterationProcessing[a.alterationID] = true
      }

      this.processing = true

      this.$SystemAPI.dalSchemaAlterationApply({ alterationID: ids }).then(() => {
        this.toastSuccess(this.$t('notification:module.alteration.resolve.success'))
      }).catch(this.toastErrorHandler(this.$t('notification:alteration.resolve.error')))
        .finally(() => {
          this.processing = false
          for (const i of ids) {
            this.$delete(this.alterationProcessing, i)
          }
          this.load(this.batch)
        })
    },

    async load (batch) {
      if (!batch) {
        return
      }

      await this.$SystemAPI.dalSchemaAlterationList({ batchID: batch }).then(({ set }) => {
        this.alterations = set
      }).catch(this.toastErrorHandler(this.$t('notification:module.alteration.load.error')))
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

    stringifyAttributeAddParams (params) {
      return `Add column ${params.attr.ident} encoded as ${params.attr.store.type} of type ${params.attr.type.type}`
    },

    stringifyAttributeDeleteParams (params) {
      return `Delete column ${params.attr.ident} encoded as ${params.attr.store.type}`
    },

    stringifyAttributeReTypeParams (params) {
      return `Changing type of column ${params.attr.ident} to ${params.to.type}`
    },

    stringifyAttributeReEncodeParams (params) {
      return `Changing encoding of column ${params.attr.ident} to ${params.to.type}`
    },

    stringifyModelAddParams (params) {
      return `Add schema for model ${params.model.ident}`
    },

    stringifyModelDeleteParams (params) {
      return `Delete schema for model ${params.model.ident}`
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
