<template>
  <b-modal
    v-model="showModal"
    :title="federationModalTitle"
    :ok-title="$t('general.label.saveAndClose')"
    ok-only
    ok-variant="dark"
    size="lg"
    body-class="p-0 border-top-0"
    header-class="p-3 pb-0 border-bottom-0"
    @ok="handleFederationSettingsSave()"
    @change="$emit('change', $event)"
  >
    <b-tabs
      active-nav-item-class="bg-grey"
      nav-wrapper-class="bg-white border-bottom"
      active-tab-class="tab-content h-auto overflow-auto"
      card
    >
      <!-- <b-tab
        :title="$t('edit.federationSettings.general.title')"
        active
      >
        <b-form-group
          class="mb-0"
        >
          <b-form-checkbox
            v-model="upstream.enabled"
          >
            {{ $t('edit.federationSettings.general.send') }}
          </b-form-checkbox>

          <b-form-checkbox
            v-model="downstream.enabled"
          >
            {{ $t('edit.federationSettings.general.receive') }}
          </b-form-checkbox>
        </b-form-group>
      </b-tab> -->
      <b-tab
        :title="$t('edit.federationSettings.upstream.title')"
        active
      >
        <b-list-group
          vertical
          class="overflow-auto server-list"
        >
          <b-list-group-item
            v-for="f in servers"
            :key="f.nodeID"
            href="#"
            :class="{ 'border border-primary': f.nodeID === upstream.active }"
            class="border d-flex flex-column"
            @click="upstream.active = f.nodeID"
          >
            <p
              class="mb-0 text-truncate"
            >
              {{ f.name }}
            </p>
            <small
              class="text-truncate"
            >
              {{ f.baseURL }}
            </small>
          </b-list-group-item>
        </b-list-group>

        <div
          v-if="upstream.processing"
          class="d-flex flex-grow-1 justify-content-center align-items-center"
        >
          <b-spinner
            variant="primary"
          />
        </div>

        <div
          v-else-if="upstream[upstream.active]"
          class="list-group flex-grow-1 ml-4"
        >
          <div
            v-if="upstream[upstream.active].canManageModule"
          >
            {{ $t('edit.federationSettings.upstream.description') }}
            <b-form-group
              label-cols-sm="4"
              label-cols-lg="5"
              :label="$t('edit.federationSettings.upstream.copyFrom')"
            >
              <b-form-select
                :key="upstream.active"
                v-model="upstream[upstream.active].copy"
                :options="upstream[upstream.active].options"
                value-field="nodeID"
                text-field="name"
                @input="copyUpstreamFrom"
              />
            </b-form-group>

            <b-form-checkbox
              :checked="upstream[upstream.active].allFields"
              class="mb-2"
              @change="selectAllFields($event, 'upstream')"
            >
              <strong>{{ $t('edit.federationSettings.upstream.allFields') }}</strong>
            </b-form-checkbox>

            <div
              class="overflow-auto"
            >
              <b-form-checkbox
                v-for="f in upstream[upstream.active].fields"
                :key="`${upstream.active}${f.name}`"
                v-model="f.value"
                class="mb-2"
                @change="checkChange($event, 'upstream')"
              >
                {{ f.label }}
              </b-form-checkbox>
            </div>
          </div>

          <div
            v-else
            class="d-flex flex-grow-1 align-items-center justify-content-center"
          >
            {{ $t('edit.federationSettings.noPermission') }}
          </div>
        </div>

        <div
          v-else
          class="d-flex flex-grow-1 align-items-center justify-content-center"
        >
          {{ $t('edit.federationSettings.noNodes') }}
        </div>
      </b-tab>

      <!-- downstream tab -->
      <b-tab
        :title="$t('edit.federationSettings.downstream.title')"
      >
        <b-list-group
          vertical
          class="overflow-auto server-list"
        >
          <b-list-group-item
            v-for="f in servers"
            :key="f.nodeID"
            href="#"
            :class="{ 'border border-primary': f.nodeID === downstream.active }"
            class="border d-flex flex-column"
            @click="downstream.active = f.nodeID"
          >
            <p
              class="mb-0 text-truncate"
            >
              {{ f.name }}
            </p>
            <small
              class="text-truncate"
            >
              {{ f.baseURL }}
            </small>
          </b-list-group-item>
        </b-list-group>

        <div
          v-if="downstream.processing"
          class="d-flex flex-grow-1 justify-content-center align-items-center"
        >
          <b-spinner
            variant="primary"
          />
        </div>

        <div
          v-else-if="downstream[downstream.active]"
          class="list-group flex-grow-1 ml-4"
        >
          <!-- dropdown list of federated shared modules -->
          <b-form-group>
            <b-form-select
              :key="downstream.active"
              v-model="downstream[downstream.active].module"
              :options="downstream[downstream.active].options"
              value-field="moduleID"
              text-field="name"
              class="w-50"
            />
          </b-form-group>

          <div
            v-if="downstream[downstream.active].module"
            class="mb-2"
          >
            {{ $t('edit.federationSettings.downstream.description') }}
            <b-form-checkbox
              :checked="downstream[downstream.active].allFields[downstream[downstream.active].module]"
              @change="selectAllFields($event, 'downstream')"
            >
              <strong>{{ $t('edit.federationSettings.downstream.allFields') }}</strong>
            </b-form-checkbox>
          </div>
          <div
            v-if="downstream[downstream.active].module"
            class="overflow-auto"
          >
            <!-- list of fields per shared module -->
            <div
              v-for="sharedModuleFields in activeSharedModules"
              :key="`${downstream.active}_${sharedModuleFields.name}`"
              class="d-flex align-items-center justify-content-between"
            >
              <b-form-checkbox
                v-model="sharedModuleFields.map"
                class="my-2"
                @change="checkChange($event, 'downstream')"
              >
                {{ sharedModuleFields.label }}
              </b-form-checkbox>

              <!-- dropdown with a list of compose module fields -->
              <b-form-select
                v-show="sharedModuleFields.map"
                :key="`${downstream.active}_${sharedModuleFields.name}`"
                v-model="sharedModuleFields.mapped"
                :options="transformedModuleFields"
                value-field="name"
                text-field="label"
                class="w-50"
                @change="setUpdated('downstream')"
              />
            </div>
          </div>
        </div>

        <div
          v-else
          class="d-flex flex-grow-1 align-items-center justify-content-center"
        >
          {{ $t('edit.federationSettings.noNodes') }}
        </div>
      </b-tab>
    </b-tabs>
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
  },

  data () {
    return {
      showModal: false,

      servers: [],

      moduleFields: [],

      sharedModule: null,

      sharedModules: {},
      sharedModulesMapped: {},
      exposedModules: {},
      moduleMappings: {},

      downstream: {
        active: undefined,
        processing: false,
        enabled: false,
      },

      upstream: {
        active: undefined,
        processing: false,
        enabled: false,
      },
    }
  },

  computed: {

    //
    // shared modules
    //
    activeSharedModules () {
      if (!this.downstream[this.downstream.active].module) return []
      return (this.sharedModulesMapped[this.downstream.active] || {})[this.downstream[this.downstream.active].module] || []
    },

    transformedModuleMappings () {
      // get the transformed module fields
      const tf = this.transformFields(this.moduleFields)

      // get the module mappings and convert it to the appropriate structure
      const mm = ((this.sharedModules[this.downstream.active] || {})[this.sharedModule] || {}).fields || []

      return tf.map((el) => {
        el.origin.value = false

        if (mm.find((e) => e.origin.name === el.origin.name)) {
          el.destination.name = el.origin.name
          el.origin.value = true
        }
        return el
      })
    },

    // used on module field dropdown on field mapping screen
    transformedModuleFields () {
      return [
        { name: null, label: this.$t('edit.federationSettings.pickModuleField') },
        ...this.transformedModuleMappings.map((el) => ({
          name: el.origin.name,
          label: el.origin.label,
        })),
      ]
    },

    federationModalTitle () {
      const { handle } = this.module
      return handle ? this.$t('edit.federationSettings.specificTitle', { handle }) : this.$t('edit.federationSettings.title')
    },
  },

  watch: {
    modal: {
      immediate: true,
      handler (show = false) {
        this.showModal = show
      },
    },

    'module.fields': {
      immediate: true,
      handler (fields) {
        this.moduleFields = fields
          .map(f => {
            return {
              kind: f.kind,
              name: f.name,
              label: f.label,
              isMulti: f.isMulti,
              value: false,
              map: null,
            }
          })
          .sort((a, b) => a.label.localeCompare(b.label))
      },
    },

    'upstream.active': {
      handler (nodeID) {
        this.getNodeUpstream(nodeID)
      },
    },

    'downstream.active': {
      handler (nodeID) {
        this.getNodeDownstream(nodeID)
      },
    },
  },

  async mounted () {
    this.preload()
  },

  methods: {
    async preload () {
      await this.$FederationAPI.nodeSearch({ status: 'paired' })
        .then(({ set = [] }) => {
          this.servers = set.filter(({ canManageNode }) => canManageNode)
        })
        .catch(this.toastErrorHandler(this.$t('edit.federationSettings.error.fetch.node')))

      for (const node of this.servers) {
        await this.loadExposedModules(node.nodeID).catch(this.toastErrorHandler(this.$t('edit.federationSettings.error.fetch.exposed')))
        await this.loadSharedModules(node.nodeID).catch(this.toastErrorHandler(this.$t('edit.federationSettings.error.fetch.shared')))
        await this.loadModuleMappings(node.nodeID).catch(this.toastErrorHandler(this.$t('edit.federationSettings.error.fetch.mmap')))
      }

      this.sharedModulesMapped = this.getSharedModulesMapped()

      if (this.servers.length) {
        this.downstream.active = this.servers[0].nodeID
        this.upstream.active = this.servers[0].nodeID
      }
    },

    // shared module fields get prepopulated here
    // the module mappings also get applied here
    // on top of the fields
    getSharedModulesMapped () {
      const list = {}

      // first, prefill the shared module fields
      for (const nodeID in this.sharedModules) {
        list[nodeID] = {}

        for (const sm of this.sharedModules[nodeID]) {
          let f = sm.fields.sort((a, b) => a.label.localeCompare(b.label))

          // is there any mappings for this shared module?
          const mappedFields = ((this.moduleMappings[nodeID] || {})[sm.moduleID] || {}).fields || []

          // fetch the shared module fields and slap the
          // module mappings on top of them
          f = f.map((el) => {
            let found = false
            let mapped = (this.moduleFields.find(({ name }) => name === el.name) || {}).name || null

            if (mappedFields) {
              const m = mappedFields.find((mf) => el.name === mf.origin.name)

              mapped = ((m || {}).destination || {}).name || null
              found = !!mapped
            }

            return {
              ...el,
              map: found,
              mapped,
            }
          })

          list[nodeID][sm.moduleID] = f
        }
      }

      return list
    },

    async handleFederationSettingsSave () {
      // module mappings (downstream)
      for (const nodeID in this.sharedModulesMapped) {
        for (const moduleID in this.sharedModulesMapped[nodeID]) {
          const crtModule = this.sharedModules[nodeID].find(m => m.moduleID === moduleID)
          // Check if node module downstream settings were updated
          if (!crtModule || !(crtModule || {}).updated) {
            continue
          }

          const fields = this.toModuleMappingFormat(this.sharedModulesMapped[nodeID][moduleID])

          const payload = {
            nodeID,
            moduleID,
            composeModuleID: this.module.moduleID,
            composeNamespaceID: this.module.namespaceID,
            fields,
          }

          await this.persistModuleMappings(payload)
            .then(() => {
              // Reset update flag
              crtModule.updated = false
            })
            .catch(this.toastErrorHandler(this.$t('edit.federationSettings.persist.mmap')))
        }
      }

      const nodes = this.servers.map(s => s.nodeID)

      // upstream
      for (const nodeID of nodes) {
        // Check if node upstream settings were updated
        if (!this.upstream[nodeID] || !(this.upstream[nodeID] || {}).updated) {
          continue
        }

        const fields = ((this.upstream[nodeID] || {}).fields || []).filter((el) => el.value)

        const payload = {
          nodeID,
          moduleID: (this.exposedModules[nodeID] || {}).moduleID,
          composeModuleID: this.module.moduleID,
          composeNamespaceID: this.module.namespaceID,
          name: this.module.name,
          handle: this.module.handle,
          fields,
        }

        const response = await this.persistExposedModule(payload)
          .then(() => {
            // Reset update flag
            (this.upstream[nodeID] || {}).updated = false
          })
          .catch(this.toastErrorHandler(this.$t('edit.federationSettings.persist.exposed')))

        if (!response && !response.moduleID) {
          return
        }

        this.exposedModules[nodeID] = response
      }
    },

    // transform internal module mappings to
    // server api format
    // [{ name, kind, ...}] => [{origin: { name, kind }, destination: { name, kind }}]
    toModuleMappingFormat (fields) {
      return fields
        .filter((el) => el.map)
        .filter((el) => !!el.mapped)
        .map((el) => ({
          origin: {
            kind: el.kind,
            name: el.name,
            label: el.label,
            isMulti: el.isMulti,
          },
          destination: {
            kind: el.kind,
            name: el.mapped,
            label: el.label,
            isMulti: el.isMulti,
          },
        }))
    },

    transformFields (fields) {
      return fields.map((el) => ({
        origin: {
          kind: el.kind,
          name: el.name,
          label: el.label || 'N/A',
          isMulti: false,
        },
        destination: {
          kind: el.kind,
          name: '',
          label: '',
          isMulti: false,
        },
      }))
    },

    getNodeUpstream (nodeID) {
      if (this.upstream[nodeID]) {
        return
      }

      this.upstream.processing = true

      const exposedModule = this.exposedModules[nodeID] || {}

      const fields = (this.moduleFields || []).map(f => ({ ...f, value: false }))

      const exposedFields = exposedModule.fields || []

      exposedFields.forEach(({ name }) => {
        fields.find(f => f.name === name).value = true
      })

      const upstream = {
        options: [
          { moduleID: null, name: this.$t('edit.federationSettings.pickServer') },
          ...this.servers.filter(s => s.nodeID !== nodeID),
        ],
        copy: null,
        allFields: false,
        fields,
        updated: false,
        canManageModule: !!exposedModule.canManageModule || !!(this.servers.find(s => s.nodeID === nodeID) || {}).canCreateModule,
      }

      upstream.allFields = upstream.fields.filter(f => f.value).length === upstream.fields.length

      this.$set(this.upstream, nodeID, upstream)
      this.upstream.processing = false
    },

    async getNodeDownstream (nodeID) {
      if (this.downstream[nodeID]) {
        return
      }

      this.downstream.processing = true

      const fields = (this.moduleFields || []).map(f => ({ ...f, value: false }))
      const downstream = {
        options: [
          { moduleID: null, name: this.$t('edit.federationSettings.pickModule') },
          ...Object.values(this.sharedModules[nodeID] || {})
            .filter(({ canMapModule }) => canMapModule)
            .map(m => ({ moduleID: m.moduleID, name: m.name })),
        ],
        module: ((this.sharedModules[nodeID] || []).find(({ handle }) => handle === this.module.handle) || {}).moduleID || null,
        allFields: {},
        fields,
      }

      Object.entries(this.sharedModulesMapped[nodeID] || {}).forEach(([key, value]) => {
        if (value.length) {
          downstream.allFields[key] = (value || []).filter(f => f.map).length === value.length
        } else {
          downstream.allFields[key] = false
        }
      })

      this.$set(this.downstream, nodeID, downstream)
      this.downstream.processing = false
    },

    selectAllFields (value, target) {
      const active = this[target].active

      if (target === 'upstream') {
        this.upstream[active].fields = this.upstream[active].fields.map(f => ({ ...f, value }))
        this.upstream[active].allFields = value
        this.upstream[active].updated = true
      } else if (target === 'downstream') {
        this.sharedModulesMapped[active][this.downstream[active].module] = this.sharedModulesMapped[active][this.downstream[active].module].map(f => ({ ...f, map: value }))
        this.downstream[active].allFields[this.downstream[active].module] = value
        this.sharedModules[active].find(m => m.moduleID === this.downstream[active].module).updated = true
      }
    },

    setUpdated (target) {
      const active = this[target].active

      if (target === 'upstream') {
        this.upstream[active].updated = true
      } else if (target === 'downstream') {
        this.sharedModules[active].find(m => m.moduleID === this.downstream[active].module).updated = true
      }
    },

    copyUpstreamFrom (nodeID) {
      this.upstream[this.upstream.active].fields = this.upstream[this.upstream.active].fields.map(f => {
        let value = false
        if (this.upstream[nodeID]) {
          value = !!this.upstream[nodeID].fields.find(({ name, value }) => name === f.name && value)
        } else {
          value = !!this.exposedModules[nodeID].fields.find(({ name }) => name === f.name)
        }

        return {
          ...f,
          value,
        }
      })
    },

    // When field checkbox changes, check and update allFields checkbox value
    checkChange (value, target) {
      const active = this[target].active

      if (target === 'upstream') {
        this.upstream[active].allFields = value ? !this.upstream[active].fields.find(f => f.value === !value) : false
      } else if (target === 'downstream') {
        const allSameValue = !this.sharedModulesMapped[active][this.downstream[active].module].find(({ map }) => map === !value)
        this.downstream[active].allFields[this.downstream[active].module] = value ? allSameValue : false
      }

      this.setUpdated(target)
    },

    async persistExposedModule (payload) {
      if (payload.moduleID) {
        return this.$FederationAPI.manageStructureUpdateExposed(payload)
      }

      return this.$FederationAPI.manageStructureCreateExposed(payload)
    },

    async persistModuleMappings (payload) {
      return this.$FederationAPI.manageStructureCreateMappings(payload)
    },

    //
    // preloaders
    //
    async loadSharedModules (nodeID) {
      if (this.sharedModules[nodeID]) {
        return
      }

      await this.$FederationAPI.manageStructureListAll({ nodeID, shared: 1 })
        .then((data = []) => {
          this.sharedModules[nodeID] = data.map(d => ({ ...d, updated: false }))
        })
        .catch(this.toastErrorHandler(this.$t('edit.federationSettings.error.fetch.shared')))
    },

    async loadExposedModules (nodeID) {
      if (this.exposedModules[nodeID]) {
        return
      }

      await this.$FederationAPI.manageStructureListAll({ nodeID, exposed: 1 })
        .then((data = []) => {
          const exposedModule = data.find(({ composeModuleID }) => composeModuleID === this.module.moduleID)
          if (exposedModule) {
            this.exposedModules[nodeID] = exposedModule
          }
        })
        .catch(this.toastErrorHandler(this.$t('edit.federationSettings.error.fetch.exposed')))
    },

    async loadModuleMappings (nodeID) {
      if (this.moduleMappings[nodeID] || !this.sharedModules[nodeID]) {
        return
      }

      const mm = {}
      for (const { moduleID } of this.sharedModules[nodeID]) {
        mm[moduleID] = []
        await this.$FederationAPI.manageStructureReadMappings({ nodeID, moduleID, composeModuleID: this.module.moduleID })
          .then((data) => {
            mm[moduleID] = data
          })
          .catch(() => {})
      }

      this.moduleMappings[nodeID] = mm
    },
  },
}
</script>

<style lang="scss" scoped>
.tab-content {
  min-height: 0;
  max-height: 70vh;
  display: flex;
}

.server-list {
  max-width: 35%;
}
</style>
