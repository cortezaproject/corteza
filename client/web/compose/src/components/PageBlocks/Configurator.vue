<template>
  <b-tabs
    :value="activeTab"
    data-test-id="page-block-configurator"
    card
    lazy
  >
    <b-tab
      data-test-id="general-tab"
      :title="$t('general.label.general')"
      title-item-class="order-first"
    >
      <b-row>
        <b-col cols="12">
          <b-form-group
            :label="$t('general.titleLabel')"
            label-class="text-primary"
          >
            <b-input-group>
              <b-form-input
                id="title"
                v-model="block.title"
                :placeholder="$t('general.titlePlaceholder')"
              />

              <b-input-group-append>
                <page-translator
                  v-if="page"
                  :page="page"
                  :block.sync="block"
                  :disabled="isNew"
                  :highlight-key="`pageBlock.${block.blockID}.title`"
                />
              </b-input-group-append>
            </b-input-group>

            <i18next
              path="interpolationFootnote"
              tag="small"
              class="text-muted"
            >
              <code>${record.values.fieldName}</code>
              <code>${recordID}</code>
              <code>${ownerID}</code>
              <span><code>${userID}</code>, <code>${user.name}</code></span>
            </i18next>
          </b-form-group>
        </b-col>

        <b-col cols="12">
          <b-form-group
            :label="$t('general.descriptionLabel')"
            label-class="text-primary"
          >
            <b-input-group>
              <b-form-textarea
                id="description"
                v-model="block.description"
                :placeholder="$t('general.descriptionPlaceholder')"
              />
              <b-input-group-append>
                <page-translator
                  v-if="page"
                  :page="page"
                  :block.sync="block"
                  :disabled="isNew"
                  :highlight-key="`pageBlock.${block.blockID}.description`"
                />
              </b-input-group-append>
            </b-input-group>

            <i18next
              path="interpolationFootnote"
              tag="small"
              class="text-muted"
            >
              <code>${record.values.fieldName}</code>
              <code>${recordID}</code>
              <code>${ownerID}</code>
              <span><code>${userID}</code>, <code>${user.name}</code></span>
            </i18next>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('general.customID.label')"
            label-class="text-primary"
          >
            <b-form-input
              id="customID"
              v-model="block.meta.customID"
              :state="customIDState"
              :placeholder="$t('general.customID.placeholder')"
            />

            <b-form-invalid-feedback
              v-if="customIDState === false"
              :state="customIDState"
            >
              {{ $t('general.customID.invalid-state') }}
            </b-form-invalid-feedback>

            <b-form-text
              v-else
              class="text-muted"
            >
              {{ $t('general.customID.description') }}
            </b-form-text>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('general.customCSSClass.label')"
            label-class="text-primary"
          >
            <b-form-input
              id="customCSSClass"
              v-model="block.meta.customCSSClass"
              :state="customCSSClassState"
              :placeholder="$t('general.customCSSClass.placeholder')"
            />

            <b-form-invalid-feedback
              v-if="customCSSClassState === false"
              :state="customCSSClassState"
            >
              {{ $t('general.customCSSClass.invalid-state') }}
            </b-form-invalid-feedback>

            <b-form-text
              v-else
              class="text-muted"
            >
              {{ $t('general.customCSSClass.description') }}
            </b-form-text>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('general.headerStyle')"
            label-class="text-primary"
          >
            <c-input-select
              id="color"
              v-model="block.style.variants.headerText"
              :options="textVariants"
              :reduce="o => o.value"
              :clearable="false"
              :placeholder="$t('general.label.none')"
              label="text"
              class="mb-1"
            />

            <b-form-checkbox
              v-model="block.style.wrap.kind"
              value="card"
              unchecked-value="plain"
              switch
            >
              {{ $t('general.wrap') }}
            </b-form-checkbox>

            <b-form-checkbox
              v-model="block.style.border.enabled"
              switch
            >
              {{ $t('general.border.show') }}
            </b-form-checkbox>
          </b-form-group>
        </b-col>

        <b-col
          v-if="block.options.magnifyOption !== undefined"
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('general.magnifyLabel')"
            label-class="text-primary"
          >
            <b-form-select
              v-model="block.options.magnifyOption"
              :options="magnifyOptions"
            />
          </b-form-group>
        </b-col>

        <b-col
          v-if="block.options.showRefresh !== undefined"
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('general.refresh.auto')"
            :description="$t('general.refresh.description')"
            label-class="text-primary"
          >
            <b-input-group
              append="s"
              class="mb-1"
            >
              <b-form-input
                v-model="block.options.refreshRate"
                type="number"
                number
                min="0"
                @blur="updateRefresh"
              />
            </b-input-group>

            <b-form-checkbox
              v-model="block.options.showRefresh"
              switch
            >
              {{ $t('general.refresh.show') }}
            </b-form-checkbox>
          </b-form-group>
        </b-col>

        <b-col
          v-if="block.options.magnifyOption !== undefined"
          cols="12"
          lg="6"
          :offset-lg="block.options.showRefresh !== undefined ? 6 : 0"
        >
          <b-form-group
            :label="$t('general.magnifyLabel')"
            label-class="text-primary"
          >
            <b-form-select
              v-model="block.options.magnifyOption"
              :options="magnifyOptions"
            />
          </b-form-group>
        </b-col>
        <b-col
          cols="12"
          sm="12"
        >
          <hr>

          <h5 class="mb-3">
            {{ $t('general.visibility.label') }}
          </h5>

          <b-form-group
            label-class="d-flex align-items-center text-primary mb-0"
          >
            <template #label>
              {{ $t('general.visibility.condition.label') }}
              <c-hint
                :tooltip="$t('general.visibility.tooltip.performance.condition')"
                icon-class="text-warning"
              />
            </template>
            <b-input-group>
              <b-input-group-prepend>
                <b-button variant="extra-light">
                  ƒ
                </b-button>
              </b-input-group-prepend>
              <b-form-input
                v-model="block.meta.visibility.expression"
                :placeholder="$t('general.visibility.condition.placeholder')"
              />
              <b-input-group-append>
                <b-button
                  variant="outline-secondary"
                  :href="visibilityDocumentationURL"
                  class="d-flex justify-content-center align-items-center"
                  target="_blank"
                >
                  ?
                </b-button>
              </b-input-group-append>
            </b-input-group>

            <i18next
              v-if="isRecordPage"
              path="general.visibility.condition.description.record-page"
              tag="small"
              class="text-muted"
            >
              <code>record.values.fieldName</code>
              <code>user.(userID/email...)</code>
              <code>screen.(width/height)</code>
              <code>isView/isCreate/isEdit</code>
              <code>user.userID == record.values.createdBy</code>
              <code>screen.width &lt; 1024</code>
            </i18next>

            <i18next
              v-else
              path="general.visibility.condition.description.non-record-page"
              tag="small"
              class="text-muted"
            >
              <code>user.(userID/email...)</code>
              <code>screen.(width/height)</code>
              <code>user.email == "test@mail.com"</code>
              <code>screen.width &lt; 1024</code>
            </i18next>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          sm="12"
          class="pt-2"
        >
          <b-form-group
            :label="$t('general.visibility.roles.label')"
            label-class="text-primary"
          >
            <c-input-select
              v-model="currentRoles"
              :options="roles.options"
              :loading="roles.processing"
              :placeholder="$t('general.visibility.roles.placeholder')"
              :get-option-label="role => role.name"
              :reduce="role => role.roleID"
              :selectable="role => !currentRoles.includes(role.roleID)"
              multiple
            />
          </b-form-group>
        </b-col>
      </b-row>
    </b-tab>

    <page-block
      v-bind="{ ...$attrs, ...$props }"
      mode="configurator"
      class="mh-tab overflow-auto"
      v-on="$listeners"
    />

    <template #tabs-end>
      <page-translator
        v-if="page"
        :page="page"
        :block.sync="block"
        :disabled="isNew"
        button-variant="link"
      />
    </template>
  </b-tabs>
</template>
<script>
import { compose, NoID } from '@cortezaproject/corteza-js'
import { handle } from '@cortezaproject/corteza-vue'
import PageTranslator from 'corteza-webapp-compose/src/components/Admin/Page/PageTranslator'
import PageBlock from './index'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    PageBlock,
    PageTranslator,
  },

  props: {
    block: {
      type: compose.PageBlock,
      required: true,
    },

    page: {
      type: compose.Page,
      required: true,
    },
  },

  data () {
    return {
      roles: {
        processing: false,
        options: [],
      },
      abortableRequests: [],
    }
  },

  computed: {
    textVariants () {
      return [
        { value: 'dark', text: this.$t('general.style.default') },
        { value: 'primary', text: this.$t('general.style.primary') },
        { value: 'secondary', text: this.$t('general.style.secondary') },
        { value: 'success', text: this.$t('general.style.success') },
        { value: 'warning', text: this.$t('general.style.warning') },
        { value: 'danger', text: this.$t('general.style.danger') },
      ]
    },

    blockClass () {
      return [
        'text-' + this.block.style.variants.headerText,
      ]
    },

    isNew () {
      return this.block.blockID === NoID
    },

    magnifyOptions () {
      return [
        { value: '', text: this.$t('general.magnifyOptions.disabled') },
        { value: 'modal', text: this.$t('general.magnifyOptions.modal') },
        { value: 'fullscreen', text: this.$t('general.magnifyOptions.fullscreen') },
      ]
    },

    customIDState () {
      return handle.handleState(this.block.meta.customID)
    },

    customCSSClassState () {
      return handle.classState(this.block.meta.customCSSClass)
    },

    activeTab () {
      return this.isNew ? 0 : 1
    },

    isRecordPage () {
      return this.page && this.page.moduleID !== NoID
    },

    visibilityDocumentationURL () {
      // eslint-disable-next-line no-undef
      const [year, month] = VERSION.split('.')
      return `https://docs.cortezaproject.org/corteza-docs/${year}.${month}/integrator-guide/compose-configuration/page-layouts.html#visibility-condition`
    },

    currentRoles: {
      get () {
        if (!this.block.meta.visibility.roles) {
          return []
        }

        return this.block.meta.visibility.roles
      },

      set (roles) {
        this.$set(this.block.meta.visibility, 'roles', roles)
      },
    },
  },

  created () {
    this.fetchRoles()
  },

  beforeDestroy () {
    this.abortRequests()
    this.setDefaultValues()
  },

  methods: {
    updateRefresh (e) {
      // If value is less than 5 but greater than 0 make it 5. Otherwise value stays the same.
      this.block.options.refreshRate = e.target.value < 5 && e.target.value > 0 ? 5 : e.target.value
    },

    fetchRoles () {
      this.roles.processing = true

      const { response, cancel } = this.$SystemAPI
        .roleListCancellable({})

      this.abortableRequests.push(cancel)

      response()
        .then(({ set: roles = [] }) => {
          this.roles.options = roles.filter(({ meta }) => !(meta.context && meta.context.resourceTypes))
        }).finally(() => {
          this.roles.processing = false
        })
    },

    abortRequests () {
      this.abortableRequests.forEach((cancel) => {
        cancel()
      })
    },

    setDefaultValues () {
      this.roles = {
        processing: false,
        options: [],
      }
      this.abortableRequests = []
    },
  },
}
</script>
<style scoped>
.mh-tab {
  max-height: calc(100vh - 16rem);
}
</style>
