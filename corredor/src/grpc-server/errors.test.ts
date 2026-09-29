import { describe, it } from 'mocha'
import { expect } from 'chai'
import * as grpc from '@grpc/grpc-js'
import pino from 'pino'
import { HandleException } from './errors'

const codec = {
  serialize: (v: unknown): Buffer => Buffer.from(JSON.stringify(v ?? {})),
  deserialize: (b: Buffer): unknown => JSON.parse(b.toString() || '{}'),
}

const service: grpc.ServiceDefinition = {
  Fail: {
    path: '/test.Test/Fail',
    requestStream: false,
    responseStream: false,
    requestSerialize: codec.serialize,
    requestDeserialize: codec.deserialize,
    responseSerialize: codec.serialize,
    responseDeserialize: codec.deserialize,
  },
}

// Throws the given value from a gRPC handler and returns the status the client receives
async function statusOf(thrown: unknown): Promise<grpc.ServiceError> {
  const server = new grpc.Server()
  server.addService(service, {
    Fail: (_: unknown, done: grpc.sendUnaryData<null>) => {
      HandleException(pino({ level: 'silent' }), thrown, done, grpc.status.UNKNOWN)
    },
  })

  const port = await new Promise<number>((resolve, reject) => {
    server.bindAsync('127.0.0.1:0', grpc.ServerCredentials.createInsecure(), (err, port) => {
      if (err) {
        reject(err)
      } else {
        resolve(port)
      }
    })
  })

  const Client = grpc.makeGenericClientConstructor(service, 'Test')
  const client = new Client(`127.0.0.1:${port}`, grpc.credentials.createInsecure()) as unknown as {
    Fail: (req: object, cb: (err: grpc.ServiceError) => void) => void
    close: () => void
  }

  try {
    return await new Promise(resolve => client.Fail({}, resolve))
  } finally {
    client.close()
    server.forceShutdown()
  }
}

describe('grpc exception handling', () => {
  it('sends the error message as the status details', async () => {
    const err = await statusOf(new Error('forced failure'))
    expect(err.code).to.equal(grpc.status.UNKNOWN)
    expect(err.details).to.equal('forced failure')
  })

  it('sends a thrown string as the status details', async () => {
    const err = await statusOf('bad input')
    expect(err.details).to.equal('bad input')
  })
})
