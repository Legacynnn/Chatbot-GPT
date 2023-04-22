'use client'

import useSWR from 'swr'
import { fetcher } from '@/http/http'
import { Chat, Message } from '@prisma/client'
import { useRouter } from 'next/navigation'

type ChatWithFirstMessage = Chat & {
  message: [Message]
}

export default function Home() {
  const router = useRouter()
  const { data: chats } = useSWR<ChatWithFirstMessage[]>('chats', fetcher, {
    fallbackData: [],
  })

  const { data: messages } = useSWR<ChatWithFirstMessage[]>(
    `chats//messages`,
    fetcher,
    {
      fallbackData: [],
    },
  )

  return (
    <div className="flex gap-5">
      <div className="flex flex-col">
        Barra lateral
        <button type="button">New Chat</button>
        <ul>
          {chats!.map((chat, key) => (
            <li key={key} onClick={() => router.push(`/?id=${chat.id}`)}>
              {chat.message[0].content}
            </li>
          ))}
        </ul>
        chats
      </div>
      <div>
        Centro
        <ul>
          <li>Mensagens</li>
        </ul>
        <form action="">
          <textarea placeholder="Digite sua pergunta"></textarea>
        </form>
      </div>
    </div>
  )
}
