import type { ErrorCode, ValidationCode } from './enums'

export interface C311FieldError {
  field: string
  code: ValidationCode
}

export interface C311ErrorPayload {
  error: ErrorCode
  message: string
  retryable: boolean
  errors?: Array<C311FieldError>
  current_version?: number
  failing_request_id?: string
  operation_id?: string
}

export class C311ApiError extends Error {
  readonly error: ErrorCode
  readonly code: ErrorCode
  readonly retryable: boolean
  readonly errors: Array<C311FieldError>
  readonly fieldErrors: Array<C311FieldError>
  readonly current_version?: number
  readonly currentVersion?: number
  readonly failing_request_id?: string
  readonly failingRequestID?: string
  readonly operation_id?: string
  readonly operationID?: string
  readonly status?: number
  readonly headers: Record<string, string>
  readonly retryAfter?: string

  constructor (payload: C311ErrorPayload, status?: number, headers: Record<string, string> = {}) {
    super(payload.message)
    this.name = 'C311ApiError'
    this.error = payload.error
    this.code = payload.error
    this.retryable = payload.retryable
    this.errors = payload.errors || []
    this.fieldErrors = this.errors
    this.current_version = payload.current_version
    this.currentVersion = payload.current_version
    this.failing_request_id = payload.failing_request_id
    this.failingRequestID = payload.failing_request_id
    this.operation_id = payload.operation_id
    this.operationID = payload.operation_id
    this.status = status
    this.headers = { ...headers }
    this.retryAfter = this.headers['retry-after'] || this.headers['Retry-After']
  }

  toJSON (): C311ErrorPayload {
    return {
      error: this.code,
      message: this.message,
      retryable: this.retryable,
      ...(this.fieldErrors.length ? { errors: this.fieldErrors } : {}),
      ...(this.currentVersion === undefined ? {} : { current_version: this.currentVersion }),
      ...(this.failingRequestID === undefined ? {} : { failing_request_id: this.failingRequestID }),
      ...(this.operationID === undefined ? {} : { operation_id: this.operationID }),
    }
  }
}

export function isC311ApiError (err: unknown): err is C311ApiError {
  return err instanceof C311ApiError
}

export function c311Error (payload: C311ErrorPayload, status?: number, headers?: Record<string, string>): C311ApiError {
  return new C311ApiError(payload, status, headers)
}
