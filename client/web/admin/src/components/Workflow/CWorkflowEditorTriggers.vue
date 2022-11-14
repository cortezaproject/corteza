<template>
  <b-card
    class="shadow-sm mt-3"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-table
      id="trigger-list"
      hover
      responsive
      head-variant="light"
      class="mb-0"
      :items="triggers"
      :fields="triggerFields"
    >
      <template #cell(constraints)="trigger">
        <samp
          v-for="(c, index) in trigger.item.constraints"
          :key="index"
        >
          {{ c.name[0].toUpperCase() + c.name.slice(1).toLowerCase() }} {{ c.op }} "{{ c.values.join(' or ') }}"
          <code
            v-if="index < trigger.item.constraints.length - 1"
          >
            {{ $t('and') }}
          </code>
        </samp>
      </template>
    </b-table>

    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>
  </b-card>
</template>

<script>
export default {
  name: 'CWorkflowEditorTriggers',

  i18nOptions: {
    namespaces: 'automation.workflows',
    keyPrefix: 'editor.triggers',
  },

  props: {
    triggers: {
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
  },

  computed: {
    triggerFields () {
      return [
        {
          key: 'resourceType',
          formatter: (rt) => {
            return rt.split(':').map(s => {
              return s[0].toUpperCase() + s.slice(1).toLowerCase()
            }).join(' ')
          },
        },
        {
          key: 'eventType',
        },
        {
          key: 'constraints',
        },
      ]
    },
  },
}
</script>
