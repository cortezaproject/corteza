import path from 'path'
import * as grpc from '@grpc/grpc-js'
import * as protoLoader from '@grpc/proto-loader'

export function LoadDefinitions(protobuf: string): Promise<grpc.GrpcObject> {
  return protoLoader
    .load(path.join(protobuf), {})
    .then((def: protoLoader.PackageDefinition) => grpc.loadPackageDefinition(def))
}

export function VerifyDefinitions({ corredor }: grpc.GrpcObject): grpc.GrpcObject {
  const has = (svc: string): boolean => Object.hasOwnProperty.call(corredor, svc)

  if (!corredor || !has('ServerScripts') || !has('ClientScripts')) {
    throw new Error('invalid or incompatible protobuf files')
  }

  return corredor as grpc.GrpcObject
}
