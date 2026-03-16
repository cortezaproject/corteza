<template>
  <div
    v-if="items.length"
    class="emoji-dropdown"
  >
    <button
      v-for="(item, index) in items"
      :key="item.name"
      type="button"
      :class="[
        'emoji-option',
        { 'emoji-option--highlighted': index === selectedIndex }
      ]"
      @click="handleClick(index)"
      @mouseenter="selectedIndex = index"
    >
      <span class="emoji-option-emoji">{{ item.emoji }}</span>
      <span class="emoji-option-name">:{{ item.name }}:</span>
    </button>
  </div>
</template>

<script>
export default {
  name: 'EmojiList',

  props: {
    items: {
      type: Array,
      required: true,
    },

    command: {
      type: Function,
      required: true,
    },
  },

  data() {
    return {
      selectedIndex: 0,
    }
  },

  watch: {
    items() {
      this.selectedIndex = 0
    },
  },

  methods: {
    onKeyDown({ event }) {
      if (event.key === 'ArrowUp') {
        this.upHandler()
        return true
      }

      if (event.key === 'ArrowDown') {
        this.downHandler()
        return true
      }

      if (event.key === 'Enter') {
        this.enterHandler()
        return true
      }

      return false
    },

    upHandler() {
      this.selectedIndex = ((this.selectedIndex + this.items.length) - 1) % this.items.length
    },

    downHandler() {
      this.selectedIndex = (this.selectedIndex + 1) % this.items.length
    },

    enterHandler() {
      this.selectItem(this.selectedIndex)
    },

    handleClick(index) {
      this.selectItem(index)
    },

    selectItem(index) {
      const item = this.items[index]

      if (item) {
        this.command({ name: item.name })
      }
    },
  },
}
</script>

<style scoped>
.emoji-dropdown {
  background: var(--p-content-background);
  border: 1px solid var(--p-content-border-color);
  border-radius: 0.25rem;
  box-shadow: 0 0.125rem 0.25rem rgba(0, 0, 0, 0.075);
  max-height: 200px;
  overflow-y: auto;
  font-size: 0.9rem;
}

.emoji-option {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  width: 100%;
  background: var(--p-content-background);
  color: var(--p-text-color);
  padding: 0.35rem 0.75rem;
  border: none;
  text-align: left;
  cursor: pointer;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.emoji-option:hover,
.emoji-option.emoji-option--highlighted {
  background: var(--p-content-hover-background);
  color: var(--p-text-color);
}

.emoji-option:active,
.emoji-option:focus {
  color: var(--p-primary-contrast-color);
  background-color: var(--p-primary-color);
  outline: none;
}

.emoji-option-emoji {
  font-size: 1.2em;
  line-height: 1;
}

.emoji-option-name {
  color: var(--p-text-muted-color);
  font-size: 0.85em;
}
</style>
