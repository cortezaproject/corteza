// Mock user directory used for picking project teams. Deterministic, no API.

export const MOCK_USERS = [
  { id: 'u-mara', name: 'Mara Novak', email: 'mara.novak@example.com' },
  { id: 'u-tomaz', name: 'Tomaž Kovač', email: 'tomaz.kovac@example.com' },
  { id: 'u-lena', name: 'Lena Horvat', email: 'lena.horvat@example.com' },
  { id: 'u-pavel', name: 'Pavel Zupan', email: 'pavel.zupan@example.com' },
  { id: 'u-ana', name: 'Ana Marić', email: 'ana.maric@example.com' },
  { id: 'u-david', name: 'David Krištof', email: 'david.kristof@example.com' },
]

// The "logged in" mock user. Creator of new projects, auto-added as developer.
export const CURRENT_USER_ID = 'u-you'

export const findUser = id => MOCK_USERS.find(u => u.id === id)

export const userName = id => findUser(id)?.name || id

export const userInitials = id => {
  const name = findUser(id)?.name || ''
  return (
    name
      .split(/\s+/)
      .filter(Boolean)
      .slice(0, 2)
      .map(p => p[0]?.toUpperCase())
      .join('') || '?'
  )
}
