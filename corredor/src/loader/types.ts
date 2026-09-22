export interface File {
  // Script location on filesystem
  src: string

  // Script reference
  ref: string

  // File modification time
  updatedAt: Date
}

export type FileType = 'client-scripts' | 'server-scripts' // | 'vue-components'
