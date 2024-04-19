<template>
  <b-card
    body-class="p-0"
    header-class="border-bottom"
    footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <template #header>
      <h4 class="m-0">
        {{ $t('title') }}
      </h4>
    </template>

    <div class="d-flex align-items-center flex-grow-1 flex-wrap flex-fill-child gap-1 p-3">
      <b-button
        variant="primary"
        size="lg"
        @click="newOIDC()"
      >
        {{ $t('oidc.add') }}
      </b-button>

      <b-form-group
        :label="$t('enabled')"
        label-class="text-primary"
        class="mb-0 ml-auto"
      >
        <c-input-checkbox
          v-model="external.enabled"
          :value="true"
          :unchecked-value="false"
          :labels="checkboxLabel"
          switch
        />
      </b-form-group>
    </div>

    <b-table
      :items="providerItems"
      :fields="providerFields"
      :tbody-tr-class="(i) => i.rowBackground"
      head-variant="light"
      hover
      class="mb-0"
    >
      <template #cell(enabled)="{ item }">
        <b-checkbox
          :checked="item.enabled"
          @change="item.enable($event)"
        />
      </template>

      <template #cell(provider)="{ item }">
        {{ item.provider || item.tag }}
      </template>

      <template #cell(editor)="{ item }">
        <c-input-confirm
          v-if="item.delete"
          :icon="item.deleted ? ['fas', 'trash-restore'] : undefined"
          :variant="item.deleted ? 'outline-warning' : 'outline-danger'"
          :variant-ok="item.deleted ? 'warning' : 'danger'"
          @confirmed="item.delete()"
        />

        <b-button
          variant="link"
          @click="openEditor(item.editor)"
        >
          <font-awesome-icon
            :icon="['fas', 'wrench']"
          />
        </b-button>
      </template>
    </b-table>

    <b-modal
      v-model="modal.open"
      :title="modal.title"
      scrollable
      size="lg"
      title-class="text-capitalize"
      @ok="modal.updater(modal.data)"
    >
      <component
        :is="modal.component"
        v-model="modal.data"
      />
    </b-modal>

    <template #footer>
      <c-button-submit
        :disabled="!dirty || !canManage"
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="$emit('submit', changes)"
      />
    </template>
  </b-card>
</template>

<script>
import _ from 'lodash'
import OidcExternal from 'corteza-webapp-admin/src/components/Settings/System/Auth/ExternalOIDC'
import StandardExternal from 'corteza-webapp-admin/src/components/Settings/System/Auth/ExternalStd'
import SamlExternal from 'corteza-webapp-admin/src/components/Settings/System/Auth/ExternalSAML'
import ConfirmationToggle from 'corteza-webapp-admin/src/components/ConfirmationToggle'

const idpStandard = [
  'google',
  'github',
  'facebook',
  'linkedin',
  'nylas',
]

const idpSecurity = {
  permittedRoles: [],
  prohibitedRoles: [],
  forcedRoles: [],
}

/**
 * Give some structure to key-val settings
 */
function prepareExternal (external) {
  const extractKey = (name, t = 'string') => {
    const v = external.find(v => v.name === `auth.external.${name}`)

    switch (t) {
      case 'string':
        return (v || { value: null }).value || ''
      case 'boolean':
        return !!(v || { value: null }).value
      case 'array':
        return (v || { value: [] }).value || []
      case undefined:
        // use t=undefined to get raw value
        return v.value
    }
  }

  const extractKeys = (provider, base = {}) => {
    let out = { ...base }

    for (let k in base) {
      out[k] = extractKey(`providers.${provider}.${k}`, Array.isArray(out[k]) ? 'array' : typeof out[k])
    }

    return out
  }

  // careful with prefix!
  // ....providers.openid-connect....
  // ....providers....
  // ....saml
  const extractSec = (prefix) => {
    return { ...idpSecurity, ...(extractKey(`${prefix}.security`, undefined) || {}) }
  }

  const data = {
    enabled: !!(external.find(v => v.name === 'auth.external.enabled') || {}).value,

    oidc: [],
    standard: [],

    saml: {
      enabled: extractKey('saml.enabled'),
      cert: extractKey('saml.cert'),
      name: extractKey('saml.name'),
      key: extractKey('saml.key'),
      'sign-method': extractKey('saml.sign-method'),
      'sign-requests': extractKey('saml.sign-requests', 'boolean'),
      binding: extractKey('saml.binding'),
      idp: {
        url: extractKey('saml.idp.url'),
        'ident-name': extractKey('saml.idp.ident-name'),
        'ident-handle': extractKey('saml.idp.ident-handle'),
        'ident-identifier': extractKey('saml.idp.ident-identifier'),
      },
      security: extractSec('saml'),
    },
  }

  data.standard = idpStandard.map(handle => ({
    handle,
    ...extractKeys(handle, {
      enabled: false,
      secret: '',
      key: '',
      security: {},
      usage: [],
    }),
    security: extractSec(`providers.${handle}`),
  }))

  const prefix = `auth.external.providers.openid-connect.`

  data.oidc =
    [...new Set(external
      // Filter out all keys that start with prefix
      .filter(v => v.name.indexOf(prefix) === 0)
      // trim off the prefix and everything after next.
      .map(({ name }) => name.substring(prefix.length).split('.', 2)[0]))]
      .map(handle => ({
        ...extractKeys('openid-connect.' + handle, {
          enabled: false,
          issuer: '',
          key: '',
          secret: '',
          scope: '',
          security: {},
        }),
        handle,
        security: extractSec('providers.openid-connect.' + handle),
        deleted: false,
      }))

  return data
}

export default {
  name: 'CSystemEditorExternal',

  i18nOptions: {
    namespaces: 'system.settings',
    keyPrefix: 'editor.external',
  },

  components: {
    OidcExternal,
    StandardExternal,
    SamlExternal,
    ConfirmationToggle,
  },

  props: {
    value: {
      type: Array,
      required: true,
    },

    processing: {
      type: Boolean,
      value: false,
    },

    success: {
      type: Boolean,
      value: false,
    },

    canManage: {
      type: Boolean,
      required: true,
    },
  },

  data () {
    return {
      // current open modal
      modal: {
        open: false,
        editor: null,
        title: null,
        data: null,
      },

      // working copy
      external: prepareExternal(this.value),

      roles: [],
      checkboxLabel: {
        on: this.$t('general:label.general.yes'),
        off: this.$t('general:label.general.no'),
      },
    }
  },

  computed: {
    dirty () {
      return this.changes.length > 0
    },

    // copy of all values for comparison
    original () {
      return Object.freeze(prepareExternal(this.value))
    },

    /**
     * Configuration for table that shows external providers
     *
     * @returns object
     */
    providerFields () {
      return [
        { key: 'enabled', label: this.$t('table.header.enabled'), thStyle: { width: '50px' } },
        { key: 'provider', label: this.$t('table.header.provider'), thStyle: { width: '200px' }, tdClass: 'text-capitalize' },
        { key: 'info', label: this.$t('table.header.info') },
        { key: 'editor', label: '', thStyle: { width: '200px' }, tdClass: 'text-right' },
      ]
    },

    providerItems () {
      return [
        {
          rowBackground: _.isEqual(this.original.saml, this.external.saml) ? '' : 'bg-extra-light',
          provider: this.external.saml.name,
          info: this.external.saml.idp.url,
          tag: 'SAML',

          enabled: this.external.saml.enabled,
          enable: (val) => this.$set(this.external.saml, 'enabled', val),

          editor: {
            component: 'saml-external',
            data: this.external.saml,
            title: this.$t('saml.title'),
            updater: (changed) => this.updater('saml', changed),
          },
        },
        ...this.external.oidc
          .map((p, i) => ({
            rowBackground: (() => {
              if (_.isEqual(this.original.oidc[i], p)) {
                return ''
              }

              if (p.deleted) {
                return 'text-extra-light deleted'
              }

              return 'bg-extra-light'
            })(),
            provider: p.handle,
            tag: 'OIDC',
            info: p.issuer,

            enabled: p.enabled,
            deleted: p.deleted,
            enable: (val) => this.$set(this.external.oidc[i], 'enabled', val),
            delete: () => {
              this.$set(this.external.oidc[i], 'deleted', !p.deleted)
            },

            editor: {
              component: 'oidc-external',
              data: p,
              title: p.handle,
              updater: (changed) => this.updater('oidc', changed, i),
            },
          })),
        ...this.external.standard.map((p, i) => ({
          rowBackground: _.isEqual(this.original.standard[i], p) ? '' : 'bg-extra-light',
          provider: p.handle,
          info: p.key,

          enabled: p.enabled,
          enable: (val) => this.$set(this.external.standard[i], 'enabled', val),

          editor: {
            component: 'standard-external',
            data: p,
            title: p.handle,
            updater: (changed) => this.updater('standard', changed, i),
          },
        })),
      ]
    },

    /**
     * Converts changed settings back to values
     */
    changes () {
      let name, value
      let c = []

      const prefix = 'auth.external.providers'
      const o = this.original
      const e = this.external

      if (!_.isEqual(o.enabled, e.enabled)) {
        c.push({ name: 'auth.external.enabled', value: e.enabled })
      }

      /**
       * General purpose key mapper
       */
      const mapKeys = (prefix, wc, org, keys) => {
        for (const k of keys) {
          if (wc[k] === undefined) {
            console.error(`potential issue - unknown key "${prefix}.${k}" used`, wc)
          }

          if (_.isEqual(wc[k], org[k])) {
            continue
          }

          name = `${prefix}.${k}`
          value = wc[k]
          c.push({ name, value })
        }
      }

      e.standard.forEach((p, i) => {
        mapKeys(
          `${prefix}.${p.handle}`,
          p,
          o.standard[i],
          ['key', 'secret', 'enabled', 'security', 'usage'],
        )
      })

      // @todo how do we remove OIDC?
      const oidcKeys = ['key', 'secret', 'enabled', 'issuer', 'scope', 'security']
      e.oidc.forEach((p, i) => {
        if (p.deleted) {
          // use base set of keys used for OIDC provider
          // plus add a few extra ones that are usually set by the backend
          //
          // @todo it would probably smarter to trigger the backend that all settings
          //       with a certain key prefix should be deleted.
          //       Unfortunately, backend does not support that for now + could be a bit dangerous
          [...oidcKeys, 'weight', 'redirect', 'label']
            .forEach(name => c.push({ name: `${prefix}.openid-connect.${p.handle}.${name}`, value: null }))
        } else {
          mapKeys(
            `${prefix}.openid-connect.${p.handle}`,
            p,
            o.oidc[i] || {},
            oidcKeys
          )
        }
      })

      mapKeys(
        `auth.external.saml`,
        e.saml,
        o.saml,
        ['enabled', 'name', 'key', 'cert', 'sign-method', 'sign-requests', 'binding', 'security']
      )

      mapKeys(
        `auth.external.saml.idp`,
        e.saml.idp,
        o.saml.idp,
        ['url', 'ident-name', 'ident-handle', 'ident-identifier']
      )

      return c
    },
  },

  watch: {
    value: {
      immediate: true,
      handler () {
        this.external = prepareExternal(this.value)
      },
    },
  },

  methods: {
    openEditor ({ component, title, data, updater }) {
      this.modal.open = true
      this.modal.component = component
      this.modal.title = title
      this.modal.updater = updater

      // deref
      this.modal.data = JSON.parse(JSON.stringify(data))
    },

    newOIDC () {
      this.openEditor({
        component: 'oidc-external',
        title: this.$t('oidc.add'),
        data: {
          handle: '',
          enabled: true,
          issuer: '',
          key: '',
          secret: '',
          scope: '',
          fresh: true,
          security: { ...idpSecurity },
        },
        updater: (changed) => {
          this.updater('oidc', changed, -1)
        },
      })
    },

    /**
     * Please note that we're assuming that array items will not get changed!
     *
     * @param key
     * @param i
     * @param val
     */
    updater (key, val, i = undefined) {
      if (i === undefined) {
        this.$set(this.external, key, val)
      } else if (i < 0) {
        this.external[key].push(val)
      } else {
        this.$set(this.external[key], i, val)
      }
    },
  },
}
</script>
<style lang="scss">
.deleted {
  text-decoration: line-through;
}
</style>
