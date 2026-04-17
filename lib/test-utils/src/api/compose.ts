import { vi } from 'vitest'

export type MockFn = ReturnType<typeof vi.fn>

export type MockComposeAPI = {
  namespaceList: MockFn
  namespaceRead: MockFn
  namespaceCreate: MockFn
  namespaceUpdate: MockFn
  namespaceDelete: MockFn
  namespaceClone: MockFn
  moduleList: MockFn
  moduleRead: MockFn
  moduleCreate: MockFn
  moduleUpdate: MockFn
  moduleDelete: MockFn
  recordList: MockFn
  recordRead: MockFn
  recordCreate: MockFn
  recordUpdate: MockFn
  recordDelete: MockFn
  recordBulk: MockFn
  pageList: MockFn
  pageRead: MockFn
  pageCreate: MockFn
  pageUpdate: MockFn
  pageDelete: MockFn
  pageLayoutList: MockFn
  pageLayoutRead: MockFn
  chartList: MockFn
  chartRead: MockFn
  permissionsEffective: MockFn
  attachmentUpload: MockFn
  attachmentRead: MockFn
  [key: string]: MockFn
}

const listResult = () => ({ set: [], filter: { total: 0 } })

export function createMockComposeAPI(overrides: Partial<MockComposeAPI> = {}): MockComposeAPI {
  return {
    namespaceList: vi.fn().mockResolvedValue(listResult()),
    namespaceRead: vi.fn().mockResolvedValue({}),
    namespaceCreate: vi.fn().mockResolvedValue({}),
    namespaceUpdate: vi.fn().mockResolvedValue({}),
    namespaceDelete: vi.fn().mockResolvedValue({}),
    namespaceClone: vi.fn().mockResolvedValue({}),
    moduleList: vi.fn().mockResolvedValue(listResult()),
    moduleRead: vi.fn().mockResolvedValue({}),
    moduleCreate: vi.fn().mockResolvedValue({}),
    moduleUpdate: vi.fn().mockResolvedValue({}),
    moduleDelete: vi.fn().mockResolvedValue({}),
    recordList: vi.fn().mockResolvedValue(listResult()),
    recordRead: vi.fn().mockResolvedValue({}),
    recordCreate: vi.fn().mockResolvedValue({}),
    recordUpdate: vi.fn().mockResolvedValue({}),
    recordDelete: vi.fn().mockResolvedValue({}),
    recordBulk: vi.fn().mockResolvedValue(listResult()),
    pageList: vi.fn().mockResolvedValue(listResult()),
    pageRead: vi.fn().mockResolvedValue({}),
    pageCreate: vi.fn().mockResolvedValue({}),
    pageUpdate: vi.fn().mockResolvedValue({}),
    pageDelete: vi.fn().mockResolvedValue({}),
    pageLayoutList: vi.fn().mockResolvedValue(listResult()),
    pageLayoutRead: vi.fn().mockResolvedValue({}),
    chartList: vi.fn().mockResolvedValue(listResult()),
    chartRead: vi.fn().mockResolvedValue({}),
    permissionsEffective: vi.fn().mockResolvedValue([]),
    attachmentUpload: vi.fn().mockResolvedValue({}),
    attachmentRead: vi.fn().mockResolvedValue({}),
    ...overrides,
  }
}
