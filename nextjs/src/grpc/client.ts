import * as protoloader from '@grpc/proto-loader'
import * as grpc from '@grpc/grpc-js'
import path from 'path'
import { ProtoGrpcType } from './rpc/chat'

const packageDefinition = protoloader.loadSync(
  path.resolve(process.cwd(), 'proto', 'proto.', 'chat.proto'),
)

const proto = grpc.loadPackageDefinition(
  packageDefinition,
) as unknown as ProtoGrpcType

export const chatClient = new proto.pb.ChatService(
  'localhost:50051',
  grpc.credentials.createInsecure(),
)
