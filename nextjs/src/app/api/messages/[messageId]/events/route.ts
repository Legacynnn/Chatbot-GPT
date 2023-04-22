import { prisma } from '@/app/prisma/prisma'
import { NextRequest } from 'next/server'
export async function GET(
  request: NextRequest,
  { params }: { params: { chatId: string } },
) {
  const message = prisma.message.findUniqueOrThrow({
    where: {
      id: params.chatId,
    },
  })

  const transformStream = new TransformStream()
  const writter = transformStream.writable.getWriter()
  const encoder = new TextEncoder()

  return new Response(transformStream.readable, {
    headers: {
      'Content-Type': 'text/event-stream',
      Connection: 'keep-alive',
      'Cache-Control': 'no-cache, no-transform',
    },
    status: 200,
  })
}
