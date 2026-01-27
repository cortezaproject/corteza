export const NODE_DIMENSIONS = {
  // Default node dimensions
  WIDTH: 260,
  HEIGHT: 100,

  // Spacing
  HORIZONTAL_SEP: 200, // Between sibling branches
  VERTICAL_SEP: 100, // Between parent-child ranks
}

// Helper to get center offset for viewport centering
export function getNodeCenterOffset() {
  const width = NODE_DIMENSIONS.WIDTH
  const height = NODE_DIMENSIONS.HEIGHT

  return {
    x: width / 2,
    y: height / 2,
  }
}
