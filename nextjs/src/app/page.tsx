'use client'

import useSWR from 'swr'
import ClientHttp, { fetcher } from '@/http/http'
import { Chat, Message } from '@prisma/client'
import { useRouter, useSearchParams } from 'next/navigation'
import { FormEvent, useEffect, useState } from 'react'

type ChatWithFirstMessage = Chat & {
  message: [Message]
}

export default function Home() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const chatIdParams = searchParams.get('id')
  const [chatId, setChatId] = useState(chatIdParams)

  const { data: chats, mutate: mutateChats } = useSWR<ChatWithFirstMessage[]>(
    'chats',
    fetcher,
    {
      fallbackData: [],
      revalidateOnFocus: false,
    },
  )

  const { data: messages, mutate: mutateMessages } = useSWR<Message[]>(
    chatIdParams ? `chats/${chatIdParams}/messages` : null,
    fetcher,
    {
      fallbackData: [],
    },
  )

  useEffect(() => {
    setChatId(chatIdParams)
  }, [chatIdParams])

  useEffect(() => {
    const textArea = document.querySelector('#message') as HTMLTextAreaElement
    textArea.addEventListener('keydown', (event) => {
      if (event.key === 'Enter' && !event.shiftKey) {
        event.preventDefault()
      }
    })
    textArea.addEventListener('keyup', (event) => {
      if (event.key === 'Enter' && !event.shiftKey) {
        const form = document.querySelector('#form') as HTMLFormElement
        const submitButton = form.querySelector('button') as HTMLButtonElement
        form.requestSubmit(submitButton)
      }
    })
  }, [])

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const textArea = event.currentTarget.querySelector(
      'textarea',
    ) as HTMLTextAreaElement
    const message = textArea.value

    if (!chatId) {
      const newChat: ChatWithFirstMessage = await ClientHttp.post(
        'chats',
        message,
      )
      mutateChats([newChat, ...chats!])
      setChatId(newChat.id)
    } else {
      const newMessage: Message = await ClientHttp.post(
        `chats/${chatId}/messages`,
        { message },
      )
      mutateMessages([...messages!, newMessage], false)
    }
    textArea.value = ''
  }

  return (
    <div className="flex gap-5">
      <div className="flex flex-col">
        Barra lateral
        <button type="button">New Chat</button>
        <ul>
          {chats!.map((chat, key) => (
            <li key={key} onClick={() => router.push(`/?id=${chat.id}`)}>
              {chat!.message[0].content}
            </li>
          ))}
        </ul>
        chats
      </div>
      <div>
        Centro
        <ul>
          {messages?.map((message, key) => {
            return <li key={message.id}>{message.content}</li>
          })}
        </ul>
        <form id="form" onSubmit={onSubmit}>
          <textarea id="message" placeholder="Digite sua pergunta"></textarea>
        </form>
      </div>
    </div>
  )
}
