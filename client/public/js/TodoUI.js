class TodoUI {
    constructor(elements) {
        this.elements = elements;
    }

    // Todoリストをレンダリングする関数
    renderTodos(todos, handlers) {
        if (!this.elements.list) {
            console.error('Todo list要素が見つかりません');
            return;
        }

        // リストをクリア
        this.elements.list.innerHTML = '';

        if (!Array.isArray(todos) || todos.length === 0) {
            this.elements.list.innerHTML = '<li class="empty-state">タスクがありません</li>';
            return;
        }

        todos.forEach(todo => {
            const li = document.createElement('li');
            li.innerHTML = `
                <input type="checkbox" class="todo-checkbox" ${todo.completed ? 'checked' : ''}>
                <span class="todo-text ${todo.completed ? 'completed' : ''}">${todo.title}</span>
                <div>
                    <span class="material-icons edit-icon">edit</span>
                    <span class="material-icons delete-icon">delete</span>
                </div>
            `;

            // イベントリスナーの設定
            this.addTodoItemListeners(li, todo.id, handlers);
            this.elements.list.appendChild(li);
        });
    }

    // Todoアイテムにイベントリスナーを追加
    addTodoItemListeners(li, todoId, handlers) {
        const checkbox = li.querySelector('.todo-checkbox');
        checkbox.addEventListener('change', () => handlers.onToggle(todoId));

        const editIcon = li.querySelector('.edit-icon');
        editIcon.addEventListener('click', () => handlers.onEdit(todoId));

        const deleteIcon = li.querySelector('.delete-icon');
        deleteIcon.addEventListener('click', () => handlers.onDelete(todoId));

        const todoText = li.querySelector('.todo-text');
        todoText.addEventListener("click", () => handlers.onShowModal(todoId));
    }

    showModal(todo, handlers) {
        const modal = document.getElementById('modal');
        const descriptionInput = document.getElementById('description');

        descriptionInput.value = todo.description === 'null' ? '' : todo.description;
        modal.style.display = "block";

        const closeModal = () => {
            modal.style.display = "none";
        };

        document.getElementById('saveDescription').onclick = () => {
            const description = descriptionInput.value.trim();
            if (description) {
                handlers.onSaveDescription(todo.id, description);
                closeModal();
            } else {
                alert('概要を入力してください');
            }
        };

        document.getElementsByClassName('close')[0].onclick = closeModal;
    }      
}