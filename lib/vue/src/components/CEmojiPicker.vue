<template>
  <div class="c-emoji-picker">
    <div class="c-emoji-picker-search-wrap">
      <input
        ref="searchInput"
        type="text"
        class="c-emoji-picker-search-input"
        :placeholder="labels.search || 'Search'"
        @input="onSearch"
        @click.stop
        @keydown.stop
      >
      <svg
        class="c-emoji-picker-search-icon"
        viewBox="0 0 16 16"
        fill="currentColor"
      >
        <path d="M11.742 10.344a6.5 6.5 0 1 0-1.397 1.398h-.001l3.85 3.85a1 1 0 0 0 1.415-1.414l-3.85-3.85zm-5.44 1.156a5 5 0 1 1 0-10 5 5 0 0 1 0 10z" />
      </svg>
    </div>

    <div
      ref="viewport"
      class="c-emoji-picker-viewport"
      :style="{ height: viewportHeight + 'px', width: viewportWidth + 'px' }"
      @scroll="onScroll"
    >
      <div
        ref="scrollContent"
        class="c-emoji-picker-scroll-content"
      />
    </div>

    <div
      v-if="showQuickReactions"
      class="c-emoji-picker-quick"
    >
      <div class="c-emoji-picker-quick-label">
        {{ labels.quickReactions || 'Quick Reactions' }}
      </div>
      <div class="c-emoji-picker-quick-row">
        <span
          v-for="qr in quickReactionsList"
          :key="qr.emoji"
          class="epi"
          :title="qr.name"
          @click.stop="onQuickReactionClick(qr)"
        >
          {{ qr.emoji }}
        </span>
      </div>
    </div>
  </div>
</template>

<script>
const STORAGE_KEY = 'corteza:emoji:frequently-used'
const MAX_FREQUENT = 18

// Virtual scroll constants (px)
const ITEM_SIZE = 30
const LABEL_HEIGHT = 26
const ITEMS_PER_ROW = 8
const BUFFER_PX = 90

const EMOJI_BLACKLIST = new Set(['relaxed', 'frowning_face'])

const GROUP_CONFIG = [
  { name: 'smileys-people', icon: '😀', label: 'Smileys & People', sources: ['', 'people & body'] },
  { name: 'animals & nature', icon: '🐶', label: 'Animals & Nature', sources: ['animals & nature'] },
  { name: 'food & drink', icon: '🍎', label: 'Food & Drink', sources: ['food & drink'] },
  { name: 'travel & places', icon: '🚗', label: 'Travel & Places', sources: ['travel & places'] },
  { name: 'activities', icon: '⚽', label: 'Activities', sources: ['activities'] },
  { name: 'objects', icon: '💡', label: 'Objects', sources: ['objects'] },
  { name: 'symbols', icon: '❤️', label: 'Symbols', sources: ['symbols'] },
]

export default {
  name: 'CEmojiPicker',

  props: {
    /**
     * Override emoji list. If not provided, uses tiptap's built-in emoji data.
     */
    emojis: {
      type: Array,
      default: () => [],
    },

    /**
     * Height of the scrollable emoji viewport in px.
     */
    viewportHeight: {
      type: Number,
      default: 260,
    },

    /**
     * Width of the scrollable emoji viewport in px.
     */
    viewportWidth: {
      type: Number,
      default: 260,
    },

    /**
     * Whether to show the "Frequently Used" section.
     */
    showFrequent: {
      type: Boolean,
      default: true,
    },

    /**
     * Whether to show the "Quick Reactions" footer.
     */
    showQuickReactions: {
      type: Boolean,
      default: true,
    },

    /**
     * Labels for translatable strings.
     */
    labels: {
      type: Object,
      default: () => ({}),
    },
  },

  data () {
    return {
      search: '',
      frequentlyUsed: [],
      virtualRows: [],
      totalHeight: 0,
      renderedRange: { start: -1, end: -1 },
    }
  },

  computed: {
    allEmojis () {
      return this.emojis || []
    },

    emojiByGroup () {
      const raw = {}
      for (const emoji of this.allEmojis) {
        if (emoji.name.startsWith('regional_indicator')) continue
        if (EMOJI_BLACKLIST.has(emoji.name)) continue
        if (!raw[emoji.group]) raw[emoji.group] = []
        raw[emoji.group].push(emoji)
      }

      const map = {}
      for (const group of GROUP_CONFIG) {
        const merged = []
        for (const src of group.sources) {
          if (raw[src]) merged.push(...raw[src])
        }
        if (merged.length) map[group.name] = merged
      }

      return map
    },

    groups () {
      const byGroup = this.emojiByGroup
      return GROUP_CONFIG.filter(g => byGroup[g.name]?.length)
    },

    frequentEmojis () {
      if (!this.showFrequent) return []
      return this.frequentlyUsed
        .map(name => this.allEmojis.find(e => e.name === name))
        .filter(Boolean)
    },

    quickReactionsList () {
      const names = ['+1', '-1', 'smile', 'tada', 'blush', 'rocket', 'eyes']
      return names
        .map(name => this.allEmojis.find(e => e.name === name))
        .filter(Boolean)
    },

    filteredEmojis () {
      if (!this.search) return []

      const q = this.search.toLowerCase()
      return this.allEmojis.filter((emoji) => {
        if (emoji.name.startsWith('regional_indicator')) return false
        if (EMOJI_BLACKLIST.has(emoji.name)) return false
        if (emoji.name.toLowerCase().includes(q)) return true
        if (emoji.shortcodes && emoji.shortcodes.some(s => s.toLowerCase().includes(q))) return true
        if (emoji.tags && emoji.tags.some(t => t.toLowerCase().includes(q))) return true
        return false
      }).slice(0, 60)
    },
  },

  mounted () {
    this.loadFrequentlyUsed()
  },

  methods: {
    /**
     * Call this externally after the picker becomes visible (e.g. popover @shown).
     * Resets search, rebuilds layout, focuses input.
     */
    async reset () {
      this.search = ''
      this.loadFrequentlyUsed()

      this.$nextTick(() => {
        if (this.$refs.searchInput) {
          this.$refs.searchInput.value = ''
          this.$refs.searchInput.focus()
        }

        const content = this.$refs.scrollContent
        if (content && !content._bound) {
          content.addEventListener('click', this.onEmojiClick)
          content._bound = true
        }

        this.rebuildVirtualRows()
        this.renderedRange = { start: -1, end: -1 }
        if (this.$refs.viewport) {
          this.$refs.viewport.scrollTop = 0
        }
        this.renderVisible()
      })
    },

    loadFrequentlyUsed () {
      try {
        const stored = localStorage.getItem(STORAGE_KEY)
        this.frequentlyUsed = stored ? JSON.parse(stored) : []
      } catch {
        this.frequentlyUsed = []
      }
    },

    saveFrequentlyUsed (emojiName) {
      const list = [emojiName, ...this.frequentlyUsed.filter(n => n !== emojiName)].slice(0, MAX_FREQUENT)
      this.frequentlyUsed = list
      try {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(list))
      } catch {
        // ignore
      }
    },

    onSearch (e) {
      this.search = e.target.value
      this.rebuildVirtualRows()
      this.renderedRange = { start: -1, end: -1 }
      if (this.$refs.viewport) {
        this.$refs.viewport.scrollTop = 0
      }
      this.renderVisible()
    },

    onEmojiClick (e) {
      const el = e.target.closest('[data-emoji]')
      if (el) {
        e.stopPropagation()
        const name = el.dataset.emoji
        const emoji = this.allEmojis.find(em => em.name === name)

        this.saveFrequentlyUsed(name)

        /**
         * Emitted when an emoji is selected.
         * @event select
         * @property {Object} emojiItem - The full emoji object { emoji, name, shortcodes, tags, group }
         */
        this.$emit('select', emoji || { name, emoji: el.textContent })
      }
    },

    rebuildVirtualRows () {
      const rows = []
      let y = 0

      const addSection = (label, emojis) => {
        rows.push({ type: 'label', text: label, y, height: LABEL_HEIGHT })
        y += LABEL_HEIGHT

        for (let i = 0; i < emojis.length; i += ITEMS_PER_ROW) {
          const chunk = emojis.slice(i, i + ITEMS_PER_ROW)
          rows.push({ type: 'emojis', items: chunk, y, height: ITEM_SIZE })
          y += ITEM_SIZE
        }
      }

      if (this.search) {
        const results = this.filteredEmojis
        if (results.length) {
          addSection(this.labels.searchResults || 'Search Results', results)
        } else {
          rows.push({ type: 'empty', y, height: 60 })
          y += 60
        }
      } else {
        const freq = this.frequentEmojis
        if (freq.length) {
          addSection(this.labels.frequentlyUsed || 'Frequently Used', freq)
        }

        const byGroup = this.emojiByGroup
        for (const group of this.groups) {
          const emojis = byGroup[group.name]
          if (emojis) {
            addSection(group.label, emojis)
          }
        }
      }

      this.virtualRows = rows
      this.totalHeight = y
    },

    onScroll () {
      this.renderVisible()
    },

    renderVisible () {
      const viewport = this.$refs.viewport
      const content = this.$refs.scrollContent
      if (!viewport || !content) return

      const rows = this.virtualRows

      content.style.height = this.totalHeight + 'px'

      if (!rows.length) {
        content.innerHTML = ''
        return
      }

      const scrollTop = viewport.scrollTop
      const rangeStart = scrollTop - BUFFER_PX
      const rangeEnd = scrollTop + this.viewportHeight + BUFFER_PX

      let first = this.findFirstRow(rangeStart)
      let last = this.findLastRow(rangeEnd)

      if (first === this.renderedRange.start && last === this.renderedRange.end) {
        return
      }
      this.renderedRange = { start: first, end: last }

      const parts = []
      for (let i = first; i <= last; i++) {
        const row = rows[i]
        if (row.type === 'label') {
          parts.push(`<div class="c-emoji-picker-section-label" style="position:absolute;top:${row.y}px;left:0;right:0;">${this.esc(row.text)}</div>`)
        } else if (row.type === 'emojis') {
          parts.push(`<div class="c-emoji-picker-vrow" style="position:absolute;top:${row.y}px;left:0;right:0;height:${ITEM_SIZE}px;">`)
          for (const e of row.items) {
            parts.push(`<span class="epi" data-emoji="${e.name}" title=":${e.name}:">${e.emoji}</span>`)
          }
          parts.push('</div>')
        } else if (row.type === 'empty') {
          parts.push(`<div class="c-emoji-picker-empty" style="position:absolute;top:${row.y}px;left:0;right:0;">${this.esc(this.labels.noResults || 'No emojis found')}</div>`)
        }
      }

      content.innerHTML = parts.join('')
    },

    findFirstRow (y) {
      const rows = this.virtualRows
      let lo = 0
      let hi = rows.length - 1
      while (lo < hi) {
        const mid = (lo + hi) >> 1
        if (rows[mid].y + rows[mid].height <= y) {
          lo = mid + 1
        } else {
          hi = mid
        }
      }
      return lo
    },

    findLastRow (y) {
      const rows = this.virtualRows
      let lo = 0
      let hi = rows.length - 1
      while (lo < hi) {
        const mid = (lo + hi + 1) >> 1
        if (rows[mid].y >= y) {
          hi = mid - 1
        } else {
          lo = mid
        }
      }
      return Math.min(lo, rows.length - 1)
    },

    onQuickReactionClick (emoji) {
      this.saveFrequentlyUsed(emoji.name)
      this.$emit('select', emoji)
    },

    esc (str) {
      return str.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')
    },
  },
}
</script>

<style lang="scss">
.c-emoji-picker-search-wrap {
  position: relative;
  padding: 0.4rem 0.5rem;
}

.c-emoji-picker-search-input {
  width: 100%;
  padding: 0.3rem 0.5rem 0.3rem 1.75rem;
  border: 1px solid var(--extra-light, #ddd);
  border-radius: 0.35rem;
  font-size: 0.8rem;
  outline: none;
  background: var(--white, #fff);
  color: var(--dark, #333);

  &:focus {
    border-color: var(--primary, #4080ff);
    box-shadow: 0 0 0 2px rgba(64, 128, 255, 0.15);
  }

  &::placeholder {
    color: var(--secondary, #999);
  }
}

.c-emoji-picker-search-icon {
  position: absolute;
  left: 1rem;
  top: 50%;
  transform: translateY(-50%);
  width: 0.8rem;
  height: 0.8rem;
  color: var(--secondary, #999);
  pointer-events: none;
}

.c-emoji-picker-viewport {
  overflow-y: auto;
  position: relative;
  padding: 0 0.4rem;
}

.c-emoji-picker-scroll-content {
  position: relative;
  width: 100%;
}

.c-emoji-picker-section-label {
  font-size: 0.7rem;
  font-weight: 700;
  color: var(--secondary, #888);
  text-transform: uppercase;
  letter-spacing: 0.03em;
  padding: 0.3rem 0.15rem 0.15rem;
  background: var(--white, #fff);
  line-height: 1;
  height: 26px;
  display: flex;
  align-items: center;
}

.c-emoji-picker-vrow {
  display: flex;
}

.epi {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  padding: 0;
  margin: 0;
  border-radius: 0.25rem;
  cursor: pointer;
  font-size: 1.15rem;
  line-height: 1;

  &:hover {
    background-color: var(--light, #f0f0f0);
    transform: scale(1.15);
  }

  &:active {
    transform: scale(0.95);
  }
}

.c-emoji-picker-empty {
  text-align: center;
  color: var(--secondary, #999);
  padding: 1.5rem 0;
  font-size: 0.8rem;
}

.c-emoji-picker-quick {
  border-top: 1px solid var(--extra-light, #e0e0e0);
  padding: 0.25rem 0.4rem 0.3rem;
}

.c-emoji-picker-quick-label {
  font-size: 0.65rem;
  font-weight: 700;
  color: var(--secondary, #888);
  text-transform: uppercase;
  letter-spacing: 0.03em;
  padding: 0.1rem 0.15rem;
}

.c-emoji-picker-quick-row {
  display: flex;
}
</style>
