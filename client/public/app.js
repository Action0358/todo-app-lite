// app.js
document.addEventListener('DOMContentLoaded', () => {
    const todoForm = document.getElementById('todo-form');
    const todoInput = document.getElementById('todo-input');
    const todoList = document.getElementById('todo-list');

    // Todoアイテムの状態を管理する配列
    let todos = [];

    // Todoアイテムを追加する関数
    function addTodo(text) {
        const todo = {
            id: Date.now(), // ユニークな識別子
            text: text,
            completed: false
        };
        todos.push(todo);
        renderTodos();
    }

    // Todoリストをレンダリングする関数
    function renderTodos() {
        todoList.innerHTML = ''; // リストをクリア
        todos.forEach(todo => {
            const li = document.createElement('li');
            li.innerHTML = `
                <input type="checkbox" class="todo-checkbox" ${todo.completed ? 'checked' : ''}>
                <span class="todo-text ${todo.completed ? 'completed' : ''}">${todo.text}</span>
                <div>
                    <span class="material-icons edit-icon">edit</span>
                    <span class="material-icons delete-icon">delete</span>
                </div>
            `;

            // チェックボックスのイベントリスナー
            const checkbox = li.querySelector('.todo-checkbox');
            checkbox.addEventListener('change', () => toggleTodo(todo.id));

            // 編集アイコンのイベントリスナー
            const editIcon = li.querySelector('.edit-icon');
            editIcon.addEventListener('click', () => editTodo(todo.id));

            // 削除アイコンのイベントリスナー
            const deleteIcon = li.querySelector('.delete-icon');
            deleteIcon.addEventListener('click', () => deleteTodo(todo.id));

            todoList.appendChild(li);
        });
    }

    // Todoのトグル（完了/未完了）
    function toggleTodo(id) {
        todos = todos.map(todo => 
            todo.id === id ? { ...todo, completed: !todo.completed } : todo
        );
        renderTodos();
    }

    // Todoの編集
    function editTodo(id) {
        const todo = todos.find(todo => todo.id === id);
        const newText = prompt('Todoを編集', todo.text);
        if (newText !== null) {
            todos = todos.map(todo => 
                todo.id === id ? { ...todo, text: newText } : todo
            );
            renderTodos();
        }
    }

    // Todoの削除
    function deleteTodo(id) {
        todos = todos.filter(todo => todo.id !== id);
        renderTodos();
    }

    // フォーム送信イベント
    todoForm.addEventListener('submit', (e) => {
        e.preventDefault(); // デフォルトの送信動作を防止
        const todoText = todoInput.value.trim();
        if (todoText) {
            addTodo(todoText);
            todoInput.value = ''; // 入力フィールドをクリア
        }
    });

    // 初期レンダリング
    renderTodos();
});