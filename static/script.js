document.addEventListener('DOMContentLoaded', function() {
    loadTasks();
});

function loadTasks() {
    fetch('/api/tasks')
        .then(response => response.json())
        .then(tasks => {
            renderTasks(tasks);
        })
        .catch(error => {
            console.error('Error loading tasks:', error);
            alert('タスクの読み込みに失敗しました');
        });
}

function renderTasks(tasks) {
    const taskList = document.getElementById('taskList');
    const emptyState = document.getElementById('emptyState');
    
    taskList.innerHTML = '';
    
    if (tasks.length === 0) {
        emptyState.style.display = 'block';
        return;
    }
    
    emptyState.style.display = 'none';
    
    tasks.forEach(task => {
        const li = document.createElement('li');
        li.className = `task-item ${task.completed ? 'completed' : ''}`;
        
        li.innerHTML = `
            <input type="checkbox" class="task-checkbox" ${task.completed ? 'checked' : ''} 
                   onchange="toggleTask(${task.id})">
            <span class="task-title">${escapeHtml(task.title)}</span>
            <button class="delete-btn" onclick="deleteTask(${task.id})">削除</button>
        `;
        
        taskList.appendChild(li);
    });
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function addTask() {
    const input = document.getElementById('taskInput');
    const title = input.value.trim();
    
    if (!title) {
        alert('タスクの内容を入力してください');
        return;
    }
    
    fetch('/api/tasks', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ title: title })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            input.value = '';
            loadTasks(); // 画面を更新して最新の一覧を表示
        } else {
            alert('タスクの追加に失敗しました');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('エラーが発生しました');
    });
}

function toggleTask(id) {
    fetch('/api/tasks/' + id + '/toggle', {
        method: 'PUT'
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            loadTasks();
        } else {
            alert('タスクの更新に失敗しました');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('エラーが発生しました');
    });
}

function deleteTask(id) {
    if (confirm('このタスクを削除しますか？')) {
        fetch('/api/tasks/' + id, {
            method: 'DELETE'
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                loadTasks();
            } else {
                alert('タスクの削除に失敗しました');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('エラーが発生しました');
        });
    }
}

document.getElementById('taskInput').addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        addTask();
    }
});
