import { vi } from 'vitest'
import type { MockFn } from './compose'

export type MockSystemAPI = {
  userList: MockFn
  userRead: MockFn
  userCreate: MockFn
  userUpdate: MockFn
  userDelete: MockFn
  roleList: MockFn
  roleRead: MockFn
  roleCreate: MockFn
  roleUpdate: MockFn
  roleDelete: MockFn
  roleMemberList: MockFn
  roleMemberAdd: MockFn
  roleMemberRemove: MockFn
  permissionsEffective: MockFn
  applicationList: MockFn
  applicationRead: MockFn
  authClientList: MockFn
  authClientRead: MockFn
  actionlogList: MockFn
  settingsList: MockFn
  settingsUpdate: MockFn
  projectAiSystemList: MockFn
  projectAiSystemRead: MockFn
  projectAiSystemCreate: MockFn
  projectAiSystemUpdate: MockFn
  projectAiSystemDelete: MockFn
  projectAiSystemEntryAdd: MockFn
  projectAiSystemEntryRemove: MockFn
  projectFriaScenarioList: MockFn
  projectFriaScenarioRead: MockFn
  projectFriaScenarioCreate: MockFn
  projectFriaScenarioUpdate: MockFn
  projectFriaScenarioDelete: MockFn
  [key: string]: MockFn
}

const listResult = () => ({ set: [], filter: { total: 0 } })

export function createMockSystemAPI(overrides: Partial<MockSystemAPI> = {}): MockSystemAPI {
  return {
    userList: vi.fn().mockResolvedValue(listResult()),
    userRead: vi.fn().mockResolvedValue({}),
    userCreate: vi.fn().mockResolvedValue({}),
    userUpdate: vi.fn().mockResolvedValue({}),
    userDelete: vi.fn().mockResolvedValue({}),
    roleList: vi.fn().mockResolvedValue(listResult()),
    roleRead: vi.fn().mockResolvedValue({}),
    roleCreate: vi.fn().mockResolvedValue({}),
    roleUpdate: vi.fn().mockResolvedValue({}),
    roleDelete: vi.fn().mockResolvedValue({}),
    roleMemberList: vi.fn().mockResolvedValue(listResult()),
    roleMemberAdd: vi.fn().mockResolvedValue({}),
    roleMemberRemove: vi.fn().mockResolvedValue({}),
    permissionsEffective: vi.fn().mockResolvedValue([]),
    applicationList: vi.fn().mockResolvedValue(listResult()),
    applicationRead: vi.fn().mockResolvedValue({}),
    authClientList: vi.fn().mockResolvedValue(listResult()),
    authClientRead: vi.fn().mockResolvedValue({}),
    actionlogList: vi.fn().mockResolvedValue(listResult()),
    settingsList: vi.fn().mockResolvedValue([]),
    settingsUpdate: vi.fn().mockResolvedValue({}),

    // AI systems + FRIA risk scenarios (persisted 2026-07-30). The create/
    // update mocks echo back an id so the store's mappers, which key their
    // caches off it, produce usable rows rather than `undefined` ids.
    projectAiSystemList: vi.fn().mockResolvedValue(listResult()),
    projectAiSystemRead: vi.fn().mockResolvedValue({ projectAiSystemID: '1', entries: [] }),
    projectAiSystemCreate: vi.fn().mockResolvedValue({ projectAiSystemID: '1', entries: [] }),
    projectAiSystemUpdate: vi.fn().mockResolvedValue({ projectAiSystemID: '1', entries: [] }),
    projectAiSystemDelete: vi.fn().mockResolvedValue({}),
    projectAiSystemEntryAdd: vi.fn().mockResolvedValue({}),
    projectAiSystemEntryRemove: vi.fn().mockResolvedValue({}),
    projectFriaScenarioList: vi.fn().mockResolvedValue(listResult()),
    projectFriaScenarioRead: vi.fn().mockResolvedValue({ projectFriaScenarioID: '1' }),
    projectFriaScenarioCreate: vi.fn().mockResolvedValue({ projectFriaScenarioID: '1' }),
    projectFriaScenarioUpdate: vi.fn().mockResolvedValue({ projectFriaScenarioID: '1' }),
    projectFriaScenarioDelete: vi.fn().mockResolvedValue({}),

    ...overrides,
  }
}
