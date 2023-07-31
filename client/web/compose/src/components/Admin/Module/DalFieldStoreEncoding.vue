<template>
  <b-row
    cols="12"
    class="mx-1 mb-2"
  >
    <b-col
      cols="3"
      align-self="center"
    >
      <b-form-checkbox
        v-if="allowOmitStrategy"
        v-model="use"
        :disabled="disabled"
      >
        {{ label }}

        <b-badge
          v-if="!use"
          variant="info"
          class="ml-2 align-middle"
        >
          {{ $t('unavailable') }}
        </b-badge>
      </b-form-checkbox>
      <div
        v-else
        class="font-weight-bold"
      >
        {{ label }}
      </div>
    </b-col>

    <b-col
      cols="3"
    >
      <c-input-select
        v-show="strategy !== 'omit'"
        v-model="strategy"
        :options="strategies"
        :disabled="!use"
        label="text"
        :reduce="strategy => strategy.value"
      />
    </b-col>

    <b-col
      v-if="strategy === ''"
      cols="6"
    >
      <b-form-input
        :value="storeIdent"
        :placeholder="$t('ident.placeholder')"
        size="sm"
        readonly
      />
    </b-col>

    <b-col
      v-else-if="showIdentInput"
      cols="6"
    >
      <b-form-input
        v-model="draft.ident"
        :placeholder="$t('ident.placeholder')"
        :disabled="disableIdentInput"
        size="sm"
      />
    </b-col>
  </b-row>
</template>
<script>
import { defaultConfigDraft, types } from './encoding-strategy'

export default {
  i18nOptions: {
    namespaces: 'module',
    keyPrefix: 'edit.config.dal.encoding-strategy',
  },

  props: {
    config: {
      type: Object,
      required: true,
    },

    field: {
      type: String,
      required: true,
    },

    label: {
      type: String,
      required: true,
    },

    isMulti: {
      type: Boolean,
      default: false,
    },

    // default store-ident
    storeIdent: {
      type: String,
      required: true,
    },

    defaultStrategy: {
      type: String,
      default: types.Plain,
    },

    allowOmitStrategy: {
      type: Boolean,
      default: true,
    },

    disabled: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      // holds working copy of strategy config
      draft: defaultConfigDraft(this.config, this.storeIdent),

      // strategy before omit
      undoOmit: this.defaultStrategy,
    }
  },

  computed: {
    strategies () {
      return [
        { value: types.Plain, text: this.$t('strategies.plain.label'), disabled: this.isMulti },
        { value: types.Alias, text: this.$t('strategies.alias.label'), disabled: this.isMulti },
        { value: types.JSON, text: this.$t('strategies.json.label') },
      ].filter(({ disabled }) => !disabled)
    },

    showIdentInput () {
      return [types.JSON, types.Alias, types.Plain].includes(this.strategy)
    },

    disableIdentInput () {
      return [types.Plain].includes(this.strategy)
    },

    // current strategy
    strategy: {
      get () {
        // iterate over all types and return the first one that matches
        for (const t of Object.values(types)) {
          if (this.config[t] === undefined) {
            continue
          }

          return t
        }

        return this.defaultStrategy
      },

      set (strategy) {
        this.$emit('change', { strategy, config: this.draft })
      },
    },

    // use => when field is not used it is omitted
    use: {
      get () {
        return this.strategy !== types.Omit
      },

      set (use) {
        if (this.strategy !== types.Omit) {
          this.undoOmit = this.strategy
        }

        this.strategy = use ? this.undoOmit : types.Omit
      },
    },
  },

  watch: {
    draft: {
      deep: true,
      handler (config) {
        this.$emit('change', { strategy: this.strategy, config })
      },
    },
  },
}
</script>
