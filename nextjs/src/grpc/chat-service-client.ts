import { chatClient } from './client'
import { ChatServiceClient as GrpcChatServiceClient } from './rpc/pb/ChatService'
export class ChatServiceClient {
  constructor(private chatClient: GrpcChatServiceClient) {}

  chatStream(data: { chat_id?: string; user_id: string; message: string }) {
    const stream = this.chatClient.chatStream({
      chatId: data.chat_id,
      userId: data.user_id,
      userMessage: data.message,
    })
    stream.on('data', (data) => {
      console.log(data)
    })
    stream.on('error', (err) => {
      console.log(err)
    })
    stream.on('end', () => {
      console.log('end')
    })
    return stream
  }
}

export class ChatServiceClientFactory {
  static create() {
    return new ChatServiceClient(chatClient)
  }
}
