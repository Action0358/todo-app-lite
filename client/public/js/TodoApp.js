class TodoApp {
    constructor() {
        this.elements = {
            form: document.querySelector('#todo-form'),
            input: document.querySelector('#todo-input'),
            list: document.querySelector('#todo-list'),
        };

        this.api = new TodoAPI('http://localhost:8080');
        this.ui = new TodoUI(this.elements);
        this.todos = [];

        this.handlers = {
            onToggle: (id) => this.toggleTodo(id),
            onEdit: (id) => this.editTodo(id),
            onDelete: (id) => this.deleteTodo(id),
            onShowModal: (id) => this.showModal(id),
            onSaveDescription: (id, description) => this.saveDescription(id, description)
        };
    }

    init() {
        // DOM要素の存在を確認
        if (!this.elements.form || !this.elements.input || !this.elements.list) {
            console.error('必要なDOM要素が見つかりません');
            return;
        }

        this.bindEvents();
        this.fetchTodos();
    }

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
    }

    // Todoの取得
    async fetchTodos() {
        try {
            const data = await this.api.fetchTodos();
            // レスポンスが配列であることを確認
            this.todos = Array.isArray(data) ? data.filter(todo => !todo.completed) : [];
            this.renderTodos();
        } catch (error) {
            // エラー時に空の配列を設定してからレンダリング
            this.todos = [];
            this.renderTodos();
            this.handleError(error);
        }
    }

    // Todoの追加
    async addTodo(text) {
        try {
            const newTodo = await this.api.addTodo(text);
            // todosが未定義または配列でない場合、初期化
            if (!Array.isArray(this.todos)) {
                this.todos = [];
            }
            this.todos.push(newTodo);
            this.renderTodos();
        } catch (error) {
            this.handleError(error);
        }
    }

    // Todoの更新
    async toggleTodo(id) {
        const todo = this.todos.find(todo => todo.id === id);
        if (!todo) {
            this.handleError('指定されたTodoが見つかりません');
            return;
        }

        try {
            // 完了状態を反転する
            const updatedTodo = { ...todo, completed: true };
            // データベースでは完了状態を更新するだけ（削除しない）
            await this.api.updateTodo(id, updatedTodo);
            // フロント側では完了状態のTodoをリストから除外する
            this.todos = this.todos.filter(todo => todo.id !== id);
        
            this.renderTodos();
        } catch (error) {
            this.fetchTodos();
            this.handleError(error);
        }
    }

    // Todoの編集
    async editTodo(id) {
        const todo = this.todos.find(todo => todo.id === id);
        if (!todo) {
            this.handleError('指定されたTodoが見つかりません');
            return;
        }

        const newText = prompt('Todoを編集', todo.title);
        if (newText === null || newText.trim() === "") return;
        try {
            const updatedTodo = await this.api.updateTodo(id, { ...todo, title: newText });
            this.todos = this.todos.map(todo => todo.id === id ? updatedTodo : todo);
            this.renderTodos();
        } catch (error) {
            this.fetchTodos(); // エラー時に最新のTodoを再取得
            this.handleError(error);
        }
    }

    // Todoの削除
    async deleteTodo(id) {
        try {
            await this.api.deleteTodo(id);
            if (!Array.isArray(this.todos)) {
                this.todos = [];
            } else {
                this.todos = this.todos.filter(todo => todo.id !== id);
            }
            this.renderTodos();
        } catch (error) {
            this.fetchTodos(); // エラー時に最新のTodoを再取得
            this.handleError(error);
        }
    }

    // 説明の保存
    async saveDescription(id, description) {
        try {
            const todo = this.todos.find(todo => todo.id === id);
            if (!todo) {
                this.handleError('指定されたTodoが見つかりません');
                return;
            }
            
            const updatedTodo = { ...todo, description };
            const result = await this.api.updateTodo(id, updatedTodo);
            this.todos = this.todos.map(todo => todo.id === id ? result : todo);
            this.renderTodos();
        } catch (error) {
            this.fetchTodos(); // エラー時に最新のTodoを再取得
            this.handleError(error);
        }
    }

    showModal(id) {
        const todo = this.todos.find(todo => todo.id === id);
        if(!todo) return;

        this.ui.showModal(todo, this.handlers);
    }
    
    // Todoリストをレンダリング
    renderTodos() {
        this.ui.renderTodos(this.todos, this.handlers);
    }

    handleError(error) {
        console.error(error);
    }
}