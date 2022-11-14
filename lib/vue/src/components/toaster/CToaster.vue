<template>
  <div>
    <b-toast
      v-for="t in toasts"
      :key="t.id"
      visible
      v-bind="t.options || {}"
      :title="t.payload.title"
      :no-close-button="!t.actions.hide"
      toast-class="overflow-unset"
      @hide="t.actions.hide ? t.actions.hide.cb(t) : evtSink"
    >
      <b-card
        header-bg-variant="transparent"
        bg-variant="transparent"
        body-class="p-0"
        class="border-0"
      >
        <b-card-text>
          {{ t.payload.notes }}
        </b-card-text>

        <b-card-text v-if="t.payload.link">
          <c-toaster-link :link="t.payload.link" />
        </b-card-text>

        <component
          :is="actComponent(act)"
          v-for="([name, act]) in extraActions(t)"
          :key="name"
          class="mr-1"
          v-bind="act"
          @action="act.cb ? act.cb(t, $event) : evtSink"
        />
      </b-card>
    </b-toast>
  </div>
</template>

<script lang="js">
import * as actions from './actions'
import CToasterLink from './CToasterLink.vue'

export default {
  components: {
    CToasterLink,
  },

  props: {
    toasts: {
      type: Array,
      required: false,
      default: () => [],
    },
  },

  methods: {
    extraActions ({ actions = {} }) {
      const { hide, ...act } = actions
      return Object.entries(act)
    },

    actComponent ({ kind }) {
      const act = actions[kind]
      if (!act) {
        throw new Error('toast.actionKind.unknown')
      }
      return act
    },

    evtSink () {},
  },
}
</script>

<style lang="scss">
// BS only provides overflow-[auto, hidden]
.toast.overflow-unset {
  overflow: unset!important;
}
</style>
