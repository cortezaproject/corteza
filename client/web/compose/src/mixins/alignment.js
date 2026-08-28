// vue mixin for alignment expression parsing
import { NoID } from '@cortezaproject/corteza-js'
import { evaluatePrefilter } from 'corteza-webapp-compose/src/lib/record-filter'

export default {
  methods: {
    /**
     * Parse alignment expression and return structured data
     * Supports: L[...] for left, C[...] for center, R[...] for right
     */
    parseAlignmentExpression (expression, record) {
      if (!expression) {
        return { hasValidMarkers: false, left: '', center: '', right: '' }
      }

      const ctx = {
        record: record || {},
        user: this.$auth.user || {},
        recordID: (record || {}).recordID || NoID,
        ownerID: (record || {}).ownedBy || NoID,
        userID: (this.$auth.user || {}).userID || NoID,
      }

      // Find all markers
      const leftMatch = expression.match(/^L\[([^\]]*)\]/)
      const centerMatch = expression.match(/C\[([^\]]*)\]/)
      const rightMatch = expression.match(/R\[([^\]]*)\]/)

      // Check if there are ANY markers
      const hasAnyMarkers = leftMatch || centerMatch || rightMatch
      if (!hasAnyMarkers) {
        return { hasValidMarkers: false, left: '', center: '', right: '' }
      }

      // Check for duplicate markers
      const allMatches = expression.match(/([LCR])\[/g) || []
      const uniqueMarkers = [...new Set(allMatches)]
      if (allMatches.length !== uniqueMarkers.length) {
        return { hasValidMarkers: false, left: '', center: '', right: '' }
      }

      // Validate order: L must come before C, C must come before R
      const lIndex = expression.indexOf('L[')
      const cIndex = expression.indexOf('C[')
      const rIndex = expression.indexOf('R[')

      const hasL = lIndex !== -1
      const hasC = cIndex !== -1
      const hasR = rIndex !== -1

      if (hasL && hasC && lIndex > cIndex) return { hasValidMarkers: false, left: '', center: '', right: '' }
      if (hasC && hasR && cIndex > rIndex) return { hasValidMarkers: false, left: '', center: '', right: '' }
      if (hasL && hasR && lIndex > rIndex) return { hasValidMarkers: false, left: '', center: '', right: '' }

      // Extract and evaluate content
      let left = ''
      let center = ''
      let right = ''

      if (leftMatch) {
        try {
          left = evaluatePrefilter(leftMatch[1], ctx) || ''
        } catch (e) {
          left = leftMatch[1]
        }
      }

      if (centerMatch) {
        try {
          center = evaluatePrefilter(centerMatch[1], ctx) || ''
        } catch (e) {
          center = centerMatch[1]
        }
      }

      if (rightMatch) {
        try {
          right = evaluatePrefilter(rightMatch[1], ctx) || ''
        } catch (e) {
          right = rightMatch[1]
        }
      }

      return {
        hasValidMarkers: true,
        left,
        center,
        right,
      }
    },
  },
}
