import { prisma } from '@/app/prisma/prisma'
import { NextRequest, NextResponse } from 'next/server'

export async function POST(request: NextRequest) {
  const chatCreated = prisma.chat.create({
    data: {},
  })

  return NextResponse.json(chatCreated)
}
