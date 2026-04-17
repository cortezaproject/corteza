import { vi } from 'vitest'
import type { MockFn } from './compose'

export type MockAutomationAPI = {
  workflowList: MockFn
  workflowRead: MockFn
  workflowCreate: MockFn
  workflowUpdate: MockFn
  workflowDelete: MockFn
  sessionList: MockFn
  sessionRead: MockFn
  triggerList: MockFn
  triggerRead: MockFn
  scriptList: MockFn
  [key: string]: MockFn
}

const listResult = () => ({ set: [], filter: { total: 0 } })

export function createMockAutomationAPI(overrides: Partial<MockAutomationAPI> = {}): MockAutomationAPI {
  return {
    workflowList: vi.fn().mockResolvedValue(listResult()),
    workflowRead: vi.fn().mockResolvedValue({}),
    workflowCreate: vi.fn().mockResolvedValue({}),
    workflowUpdate: vi.fn().mockResolvedValue({}),
    workflowDelete: vi.fn().mockResolvedValue({}),
    sessionList: vi.fn().mockResolvedValue(listResult()),
    sessionRead: vi.fn().mockResolvedValue({}),
    triggerList: vi.fn().mockResolvedValue(listResult()),
    triggerRead: vi.fn().mockResolvedValue({}),
    scriptList: vi.fn().mockResolvedValue(listResult()),
    ...overrides,
  }
}
