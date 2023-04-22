export default function Home() {
  return (
    <div className="flex gap-5">
      <div className="flex flex-col">
        Barra lateral
        <button type="button">New Chat</button>
        <ul>
          <li>Chat 1</li>
          <li>Chat 1</li>
          <li>Chat 1</li>
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
