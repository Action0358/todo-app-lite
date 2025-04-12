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
            this.addTodoItemListeners(li, todo.id, todo.title, handlers);
            this.elements.list.appendChild(li);
        });
    }

    // モーダルを表示
    showModal(todo, handlers) {
        const modal = document.getElementById('modal');
        const descriptionInput = document.getElementById('description');
        const saveButton = document.getElementById('saveDescription');
        const closeButton = document.getElementById('save-cancel-button')

        descriptionInput.value = todo.description === 'null' ? '' : todo.description;

         // モーダルウィンドウを表示
        modal.style.display = "block";

         // モーダルウィンドウを非表示
        const closeModal = () => {
            modal.style.display = "none";
        };

        // 保存ボタンのイベント設定
        saveButton.onclick = () => {
            const description = descriptionInput.value.trim();
            if (description) {
                handlers.onSaveDescription(todo.id, description);
                closeModal();
            } else {
                alert('概要を入力してください');
            }
        };

        // クリックしたときに閉じる処理
        const handleCloseClick = (event) => {
            // モーダルの背景部分（モーダル自体）がクリックされた場合のみ閉じる
            if (event.target === modal || event.target === closeButton) {
                closeModal();
            }
        };
        // 共通のクリックイベント登録
        modal.onclick = handleCloseClick;
        closeButton.onclick = handleCloseClick;
    }
    
    // 削除確認ポップアップを表示
    showDeleteConfirmation(todoId, todoTitle, handlers) {
        const modal = document.getElementById('delete-confirmation-modal');
        const modalMassage = document.getElementById('delete-confirmation-message');
        const deleteButton = document.getElementById('delete-confirm-button');
        const cancelButton = document.getElementById('delete-cancel-button');
        
        // メッセージを設定
        modalMassage.textContent = 
            `「${todoTitle}」を削除してもよろしいですか？`;
        
        // モーダルを表示
        modal.style.display = "block";
        
        // 確認ボタンのイベント設定（シンプル化）
        deleteButton.onclick = () => {
            modal.style.display = "none";
            handlers.onDelete(todoId);
        };
        
        // キャンセルボタンのイベント設定（シンプル化）
        cancelButton.onclick = () => {
            modal.style.display = "none";
        };
        
        // モーダル自体をクリックしたときに閉じる処理（シンプル化）
        modal.onclick = (event) => {
            // モーダルの背景部分がクリックされた場合のみ閉じる
            if (event.target === modal) {
                modal.style.display = "none";
            }
        };
    }

    // Todoアイテムにイベントリスナーを追加
    addTodoItemListeners(li, todoId, todoTitle, handlers) {
        const checkbox = li.querySelector('.todo-checkbox');
        checkbox.addEventListener('change', () => handlers.onToggle(todoId));

        const editIcon = li.querySelector('.edit-icon');
        editIcon.addEventListener('click', () => handlers.onEdit(todoId));

        const deleteIcon = li.querySelector('.delete-icon');
        deleteIcon.addEventListener('click', () => {
            this.showDeleteConfirmation(todoId, todoTitle, handlers);
        });

        const todoText = li.querySelector('.todo-text');
        todoText.addEventListener("click", () => handlers.onShowModal(todoId));
    }
}