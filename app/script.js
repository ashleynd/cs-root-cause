function addTodo(text) {
  const todoList = document.getElementById('todoList')
  const entry = document.createElement('li')

  entry.appendChild(document.createTextNode(text))
  todoList.appendChild(entry)
}

async function updateTodo() {
  const textBox = document.getElementById('description')

  const text = textBox.value
  if (text === '') {
    return
  }

  addTodo(text)
  await postData(text)

  textBox.value = ''
}

async function postData(description) {
  const response = await fetch('/update', {
    method: 'POST',
    body: description,
  })
  return response
}

document.getElementById('newTodoForm').addEventListener('submit', async e => {
  e.preventDefault()
  await updateTodo()
})
