let _seq = 9000
const nextID = () => String(_seq++)

export interface UserOverrides {
  userID?: string
  name?: string
  handle?: string
  email?: string
}

export function makeUser(overrides: UserOverrides = {}) {
  return {
    userID: overrides.userID ?? nextID(),
    name: overrides.name ?? 'Test User',
    handle: overrides.handle ?? 'test-user',
    email: overrides.email ?? 'test@example.com',
  }
}
