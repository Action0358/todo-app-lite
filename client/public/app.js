// app.js
document.addEventListener('DOMContentLoaded', () => {
    const todoApp = {
        elements: {
            form: document.getElementById('todo-form'),
            input: document.getElementById('todo-input'),
            list: document.getElementById('todo-list')
        },
        config: {
            API_BASE_URL: 'http://localhost:8080'
        },
        todos: [],

        init() {
            // DOM要素の存在を確認
            if (!this.elements.form || !this.elements.input || !this.elements.list) {
                console.error('必要なDOM要素が見つかりません');
                return;
            }
            
            this.bindEvents();
            this.fetchTodos();
        },

        bindEvents() {
            // フォーム送信イベント
            this.elements.form.addEventListener('submit', (e) => {
                e.preventDefault();
                const todoText = this.elements.input.value.trim();
                if (todoText) {
                    this.addTodo(todoText);
                    this.elements.input.value = '';
                }
            });
        },

        // Todoの取得
        async fetchTodos() {
            try {
                const response = await fetch(`${this.config.API_BASE_URL}/todos`, {
                    method: 'GET'
                });
                if (!response.ok) throw new Error('Todoの取得に失敗');
                
                const data = await response.json();
                // レスポンスが配列であることを確認
                this.todos = Array.isArray(data) ? data : [];
                
                this.renderTodos();
            } catch (error) {
                // エラー時に空の配列を設定してからレンダリング
                this.todos = [];
                this.renderTodos();
                this.handleError('Todoの取得中にエラーが発生しました', error);
            }
        },

        // Todoの追加
        async addTodo(text) {
            try {
                const response = await fetch(`${this.config.API_BASE_URL}/todos`, {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({
                        title: text,
                        description: 'No description',
                        completed: false
                    })
                });
                if (!response.ok) throw new Error('Todoの追加に失敗');
                
                const newTodo = await response.json();
                
                // todosが未定義または配列でない場合、初期化
                if (!Array.isArray(this.todos)) {
                    this.todos = [];
                }
                
                this.todos.push(newTodo);
                this.renderTodos();
            } catch (error) {
                this.handleError('Todoの追加中にエラーが発生しました', error);
            }
        },

        // Todoの更新(トグル)
        async toggleTodo(id) {
            // todosが未定義または配列でない場合、何もしない
            if (!Array.isArray(this.todos)) {
                this.handleError('Todosが正しく初期化されていません');
                return;
            }
            
            const todo = this.todos.find(t => t.id === id);
            if (!todo) return;

            // フロントエンドの状態を更新
            const updatedTodo = { ...todo, completed: !todo.completed };
            this.todos = this.todos.map(t => t.id === id ? updatedTodo : t);
            this.renderTodos();

            try {
                const response = await fetch(`${this.config.API_BASE_URL}/todos/${id}`, {
                    method: 'PUT',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify(updatedTodo)  
                });
                if (!response.ok) throw new Error('Todoの更新に失敗');

                const result = await response.json();
                const index = this.todos.findIndex(t => t.id === id);
                if (index !== -1) {
                    this.todos[index] = result;
                    this.renderTodos();
                }
            } catch (error) {
                // エラー時に最新のTodoを再取得
                this.fetchTodos();
                this.handleError('Todoの更新中にエラーが発生しました', error);
            }
        },

        // Todoの編集
        async editTodo(id) {
            // todosが未定義または配列でない場合、何もしない
            if (!Array.isArray(this.todos)) {
                this.handleError('Todosが正しく初期化されていません');
                return;
            }
            
            const todo = this.todos.find(todo => todo.id === id);
            if (!todo) {
                this.handleError('指定されたTodoが見つかりません');
                return;
            }

            const newText = prompt('Todoを編集', todo.title);
            if (newText === null || newText.trim() === "") {
                return; // 空のテキストまたはキャンセルの場合は何もしない
            }

            // フロントエンドの状態を更新
            this.todos = this.todos.map(todo => 
                todo.id === id ? { ...todo, title: newText } : todo
            );
            this.renderTodos();

            // サーバーに更新を送信
            try {
                const response = await fetch(`${this.config.API_BASE_URL}/todos/${id}`, {
                    method: 'PUT',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({ 
                        ...todo,
                        title: newText 
                    })
                });
                if (!response.ok) throw new Error('Todoの更新に失敗');

                const updatedTodo = await response.json();
                const index = this.todos.findIndex(t => t.id === id);
                if (index !== -1) {
                    this.todos[index] = updatedTodo; // サーバーからの更新を反映
                    this.renderTodos();
                }
            } catch (error) {
                // エラー時に最新のTodoを再取得
                this.fetchTodos();
                this.handleError('Todoの編集中にエラーが発生しました', error);
            }
        },

        // Todoの削除
        async deleteTodo(id) {
            try {
                const response = await fetch(`${this.config.API_BASE_URL}/todos/${id}`, {
                    method: 'DELETE'
                });
                if (!response.ok) throw new Error('Todoの削除に失敗');
                
                // todosが未定義または配列でない場合、初期化
                if (!Array.isArray(this.todos)) {
                    this.todos = [];
                } else {
                    this.todos = this.todos.filter(t => t.id !== id);
                }
                
                this.renderTodos();
                
                // 全てのTodoが削除された場合、明示的に空の状態を表示
                if (this.todos.length === 0) {
                    this.elements.list.innerHTML = '<li class="empty-state">タスクがありません。新しいタスクを追加してください。</li>';
                }
            } catch (error) {
                // エラー時に最新のTodoを再取得
                this.fetchTodos();
                this.handleError('Todoの削除中にエラーが発生しました', error);
            }
        },

        // Todoリストをレンダリングする関数
        renderTodos() {
            // list要素が存在するか確認
            if (!this.elements.list) {
                console.error('Todo list要素が見つかりません');
                return;
            }
            
            // todosが未定義または配列でない場合、初期化
            if (!Array.isArray(this.todos)) {
                this.todos = [];
            }
            
            this.elements.list.innerHTML = ''; // リストをクリア
            
            // Todoがない場合の表示
            if (this.todos.length === 0) {
                this.elements.list.innerHTML = '<li class="empty-state">タスクがありません。新しいタスクを追加してください。</li>';
                return;
            }
            
            this.todos.forEach(todo => {
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
                this.addTodoItemListeners(li, todo.id);
                this.elements.list.appendChild(li);
            });
        },

        // Todoアイテムにイベントリスナーを追加
        addTodoItemListeners(li, todoId) {
            // チェックボックスのイベントリスナー
            const checkbox = li.querySelector('.todo-checkbox');
            checkbox.addEventListener('change', () => this.toggleTodo(todoId));

            // 編集アイコンのイベントリスナー
            const editIcon = li.querySelector('.edit-icon');
            editIcon.addEventListener('click', () => this.editTodo(todoId));

            // 削除アイコンのイベントリスナー
            const deleteIcon = li.querySelector('.delete-icon');
            deleteIcon.addEventListener('click', () => this.deleteTodo(todoId));
        },

        // エラー処理の統一
        handleError(message, error = null) {
            console.error(message, error);
            // ここにユーザーへのエラー表示ロジックを追加することも可能
        }
    };

    // アプリを初期化
    todoApp.init();
});