/**
 * Execution trace types for TAQ (Next-Gen Automation)
 *
 * Mirrors the backend types from pkg/automation_exec/types.
 */

/**
 * A single frame in the execution trace.
 * Each step execution produces one StackFrame.
 * Triggers do NOT produce StackFrames — only steps do.
 */
export interface StackFrame {
  id: string
  stepID: string
  parentID?: string
  handle?: string
  kind?: string
  input?: Record<string, unknown>
  args?: Record<string, unknown>
  output?: Record<string, unknown>
  startedAt: string
  endedAt?: string
  error?: string
}

/**
 * Result of a TAQ execution.
 * Returned by the exec endpoint after the automation completes.
 */
export interface ExecutionResult {
  executionID: string
  executableID: string
  revision: number
  status: ExecutionStatus
  error?: string
  startedAt: string
  endedAt?: string
  duration: string
}

export type ExecutionStatus =
  | 'created'
  | 'running'
  | 'paused'
  | 'completed'
  | 'failed'
  | 'cancelled'

/**
 * UI-level trace status for the builder overlay.
 */
export type TraceStatus = 'idle' | 'running' | 'completed' | 'failed'
