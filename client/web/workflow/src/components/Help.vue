<template>
  <div class="flex flex-col gap-4 p-3">
    <div>
      <h5 class="flex items-center gap-2 font-semibold text-primary mb-2">
        {{ $t('help.basic-controls-label') }}
        <i class="pi pi-mouse" />
      </h5>
      <DataTable
        :value="basicControls"
        :show-headers="false"
        class="border border-surface rounded-border"
      >
        <Column field="action" />
        <Column field="shortcut" />
      </DataTable>
    </div>

    <div>
      <h5 class="flex items-center gap-2 font-semibold text-primary mb-1">
        {{ $t('help.keyboard-shortcuts-label') }}
        <i class="pi pi-th-large" />
      </h5>
      <small class="text-muted-color mb-2 block">
        {{ $t('help.mac-users') }}
      </small>
      <DataTable
        :value="keyboardShortcuts"
        :show-headers="false"
        class="border border-surface rounded-border"
      >
        <Column field="action" />
        <Column field="shortcut" />
      </DataTable>
    </div>
  </div>
</template>

<script>
export default {
  data () {
    return {
      basicControls: [
        { action: this.$t('help.basic-controls.action.select'), shortcut: this.$t('help.basic-controls.shortcut.left-click.label') },
        { action: this.$t('help.basic-controls.action.select-multiple'), shortcut: this.$t('help.basic-controls.shortcut.left-click.drag') },
        { action: this.$t('help.basic-controls.action.deselect'), shortcut: this.$t('help.basic-controls.shortcut.left-click.background') },
        { action: this.$t('help.basic-controls.action.pan'), shortcut: this.$t('help.basic-controls.shortcut.pan') },
        { action: this.$t('help.basic-controls.action.zoom-in/out'), shortcut: this.$t('help.basic-controls.shortcut.zoom') },
      ],

      keyboardShortcuts: [
        { action: this.$t('help.keyboard-shortcuts.action.select-all'), shortcut: 'Ctrl + A' },
        { action: this.$t('help.keyboard-shortcuts.action.add-selection'), shortcut: `Ctrl + ${this.$t('help.keyboard-shortcuts.shortcut.left-mouse-button')}` },
        { action: this.$t('help.keyboard-shortcuts.action.delete-selected-elements'), shortcut: 'Delete or Backspace' },
        { action: this.$t('help.keyboard-shortcuts.action.reset-view-to-start-point'), shortcut: 'Ctrl + Space' },
        { action: this.$t('help.keyboard-shortcuts.action.undo'), shortcut: 'Ctrl + Z' },
        { action: this.$t('help.keyboard-shortcuts.action.redo'), shortcut: 'Ctrl + Shift + Z' },
        { action: this.$t('help.keyboard-shortcuts.action.cut'), shortcut: 'Ctrl + X' },
        { action: this.$t('help.keyboard-shortcuts.action.copy'), shortcut: 'Ctrl + C' },
        { action: this.$t('help.keyboard-shortcuts.action.paste'), shortcut: 'Ctrl + V' },
        { action: this.$t('help.keyboard-shortcuts.action.save'), shortcut: 'Ctrl + S' },
        { action: this.$t('help.keyboard-shortcuts.action.nudge'), shortcut: `${this.$t('help.keyboard-shortcuts.shortcut.arrow-keys')} Shift ${this.$t('help.keyboard-shortcuts.shortcut.adjust-distance')}` },
        { action: this.$t('help.keyboard-shortcuts.action.hide-guides'), shortcut: `${this.$t('help.keyboard-shortcuts.shortcut.hold')} Alt` },
        { action: this.$t('help.keyboard-shortcuts.action.show-help'), shortcut: 'Shift + ?' },
      ],
    }
  },
}
</script>
