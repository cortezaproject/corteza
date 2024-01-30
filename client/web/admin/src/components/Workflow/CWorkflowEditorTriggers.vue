<template>
  <b-card
    body-class="p-0"
    class="shadow-sm mt-3 overflow-hidden"
  >
    <b-table
      id="trigger-list"
      :items="triggers"
      :fields="triggerFields"
      hover
      responsive
      head-variant="light"
      class="mb-0"
    >
      <template #cell(constraints)="trigger">
        <samp
          v-for="(c, index) in trigger.item.constraints"
          :key="index"
        >
          <template
            v-if="c.name"
          >
            {{ c.name[0].toUpperCase() + c.name.slice(1).toLowerCase() }} {{ c.op }} "{{ c.values.join(' or ') }}"
          </template>
          <code
            v-if="index < trigger.item.constraints.length - 1"
          >
            {{ $t('and') }}
          </code>
        </samp>
      </template>
    </b-table>

    <template #header>
      <h4 class="m-0">
        {{ $t('title') }}
      </h4>
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
