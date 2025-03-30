class TodoAPI {
    constructor(baseURL) {
        this.baseURL = baseURL;
    }

    async fetchTodos() {
        const response = await fetch(`${this.baseURL}/todos`, {
            method: 'GET'
        });
        if (!response.ok) throw new Error('Todoの取得に失敗しました');
        return await response.json();
    }

    async addTodo(text) {
        const response = await fetch(`${this.baseURL}/todos`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                title: text,
                description: 'null',
                completed: false
            })
        });
        if (!response.ok) throw new Error('Todoの追加に失敗');
        return await response.json();
    }

    async updateTodo(id, todoData) {
        const response = await fetch(`${this.baseURL}/todos/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(todoData)
        });
        if (!response.ok) throw new Error('Todoの更新に失敗');
        return await response.json();
    }

    async deleteTodo(id) {
        const response = await fetch(`${this.baseURL}/todos/${id}`, {
            method: 'DELETE'
        });
        if (!response.ok) throw new Error('Todoの削除に失敗');
    }
}