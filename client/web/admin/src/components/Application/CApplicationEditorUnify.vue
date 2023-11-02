<template>
  <b-card
    data-test-id="card-application-selector"
    header-bg-variant="white"
    footer-bg-variant="white"
    footer-class="d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <b-form
      @submit.prevent="$emit('submit', unify)"
    >
      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('name.label')"
            :description="$t('name.description')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="unify.name"
              data-test-id="input-name"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('logo.label')"
            :description="$t('logo.description')"
            label-class="text-primary"
          >
            <template #label>
              <div
                class="d-flex align-items-center"
              >
                {{ $t('logo.label') }}

                <b-button
                  v-if="showLogoPreview"
                  v-b-modal.logo
                  data-test-id="button-logo-show"
                  variant="link"
                  class="d-flex align-items-center border-0 p-0 ml-2"
                >
                  <font-awesome-icon
                    :icon="['fas', 'eye']"
                  />
                </b-button>

                <b-button
                  v-if="showLogoPreview"
                  data-test-id="button-logo-reset"
                  variant="light"
                  size="sm"
                  class="py-0 ml-2"
                  @click="resetLogo()"
                >
                  {{ $t('logo.reset') }}
                </b-button>
              </div>
            </template>
            <b-form-file
              v-model="unifyAssets.logo"
              data-test-id="file-logo-upload"
              accept="image/*"
              :placeholder="$t('logo.placeholder')"
              @change="$emit('change-detected')"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('url.label')"
            :description="$t('url.description')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="unify.url"
              data-test-id="input-url"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-row>
            <b-col
              cols="12"
              sm="6"
            >
              <b-form-group
                :label="$t('listed')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="unify.listed"
                  data-test-id="checkbox-listed"
                  switch
                  :labels="checkboxLabel"
                />
              </b-form-group>
            </b-col>

            <b-col
              cols="12"
              sm="6"
            >
              <b-form-group
                :label="$t('pinned')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="unify.pinned"
                  data-test-id="checkbox-pinned"
                  switch
                  :labels="checkboxLabel"
                  :disabled="!canPin"
                />
              </b-form-group>
            </b-col>
          </b-row>
        </b-col>

        <b-col
          cols="12"
        >
          <b-form-group
            :label="$t('config.label')"
            :description="$t('config.description')"
            label-class="text-primary"
          >
            <b-form-textarea
              v-model="unify.config"
              data-test-id="textarea-config"
              :state="configState"
              rows="10"
            />
          </b-form-group>
        </b-col>
      </b-row>
    </b-form>

    <template #header>
      <h3
        data-test-id="card-title"
        class="m-0"
      >
        {{ $t('title') }}
      </h3>
    </template>

    <template #footer>
      <c-button-submit
        :disabled="disabled"
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="$emit('submit', { unify, unifyAssets })"
      />
    </template>

    <b-modal
      id="logo"
      hide-header
      hide-footer
      centered
      body-class="p-1"
    >
      <b-img
        data-test-id="img-logo-preview"
        :src="unify.logo"
        fluid-grow
      />
    </b-modal>
  </b-card>
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'

export default {
  name: 'CApplicationEditorUnify',

  i18nOptions: {
    namespaces: 'system.applications',
    keyPrefix: 'editor.unify',
  },

  props: {
    unify: {
      type: Object,
      required: true,
    },

    application: {
      type: Object,
      required: true,
    },

    canPin: {
      type: Boolean,
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
  },

  data () {
    return {
      unifyAssets: {
        icon: undefined,
        logo: undefined,
      },

      checkboxLabel: {
        on: this.$t('general:label.general.yes'),
        off: this.$t('general:label.general.no'),
      },
    }
  },

  computed: {
    disabled () {
      return this.validConfig === false
    },

    validConfig () {
      if (!this.unify) {
        return null
      }

      try {
        if ((this.unify.config || '').trim() !== '') {
          JSON.parse(this.unify.config)
        }
        return null
      } catch (e) {
        return false
      }
    },

    configState () {
      if (((this.unify || {}).config || '').trim() === '') {
        return null
      } else {
        return this.validConfig
      }
    },

    showLogoPreview () {
      return this.unify.logoID !== NoID
    },
  },

  created () {
    this.unify.name = this.unify.name ? this.unify.name : this.application.name
  },

  methods: {
    resetLogo () {
      this.unify.logo = undefined
      this.unify.logoID = NoID
    },
  },
}
</script>
