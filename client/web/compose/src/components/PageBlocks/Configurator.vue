<template>
  <b-tabs
    data-test-id="page-block-configurator"
    nav-wrapper-class="bg-white border-bottom"
    card
    lazy
  >
    <b-tab
      data-test-id="general-tab"
      active
      :title="$t('general.label.general')"
      title-item-class="order-first"
    >
      <b-row>
        <b-col
          cols="12"
        >
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
              <code>${userID}</code>
            </i18next>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
        >
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
              <code>${userID}</code>
            </i18next>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          sm="6"
          class="mb-2"
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
              :placeholder="$t('general.label.none')"
              label="text"
            />
          </b-form-group>

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

          <b-form-checkbox
            v-if="block.kind !== 'Tabs'"
            v-model="block.meta.hidden"
            switch
          >
            {{ $t('general.hidden.label') }}
          </b-form-checkbox>
        </b-col>

        <b-col
          v-if="block.options.showRefresh !== undefined"
          cols="12"
          sm="6"
        >
          <b-form-group
            :label="$t('general.refresh.auto')"
            :description="$t('general.refresh.description')"
            label-class="text-primary"
            class="mb-1"
          >
            <b-input-group append="s">
              <b-form-input
                v-model="block.options.refreshRate"
                type="number"
                number
                min="0"
                @blur="updateRefresh"
              />
            </b-input-group>
          </b-form-group>
          <b-form-checkbox
            v-model="block.options.showRefresh"
            switch
            class="mb-2"
          >
            {{ $t('general.refresh.show') }}
          </b-form-checkbox>
        </b-col>

        <b-col
          v-if="block.options.magnifyOption !== undefined"
          cols="12"
          sm="6"
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
  },

  methods: {
    updateRefresh (e) {
      // If value is less than 5 but greater than 0 make it 5. Otherwise value stays the same.
      this.block.options.refreshRate = e.target.value < 5 && e.target.value > 0 ? 5 : e.target.value
    },
  },
}
</script>
<style scoped>
.mh-tab {
  max-height: calc(100vh - 16rem);
}
</style>
