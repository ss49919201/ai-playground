let isEditing = false;

document.addEventListener('DOMContentLoaded', function() {
    loadMemos();
});

async function loadMemos() {
    try {
        const response = await fetch('/api/memos');
        const memos = await response.json();
        displayMemos(memos);
    } catch (error) {
        console.error('メモの読み込みに失敗しました:', error);
    }
}

function displayMemos(memos) {
    const container = document.getElementById('memos-container');
    
    if (!memos || memos.length === 0) {
        container.innerHTML = '<div class="empty-state">メモがありません</div>';
        return;
    }
    
    container.innerHTML = memos.map(memo => `
        <div class="memo-item">
            <div class="memo-title">${escapeHtml(memo.title)}</div>
            <div class="memo-content">${escapeHtml(memo.content)}</div>
            <div class="memo-meta">
                作成日: ${formatDate(memo.created_at)} | 更新日: ${formatDate(memo.updated_at)}
            </div>
            <div class="memo-actions">
                <button class="edit-btn" onclick="editMemo(${memo.id})">編集</button>
                <button class="delete-btn" onclick="deleteMemo(${memo.id})">削除</button>
            </div>
        </div>
    `).join('');
}

async function saveMemo() {
    const title = document.getElementById('memo-title').value.trim();
    const content = document.getElementById('memo-content').value.trim();
    const memoId = document.getElementById('memo-id').value;
    
    if (!title) {
        alert('タイトルを入力してください');
        return;
    }
    
    try {
        let response;
        if (isEditing && memoId) {
            response = await fetch(`/api/memos/${memoId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ title, content })
            });
        } else {
            response = await fetch('/api/memos', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ title, content })
            });
        }
        
        if (response.ok) {
            clearForm();
            loadMemos();
        } else {
            alert('保存に失敗しました');
        }
    } catch (error) {
        console.error('保存に失敗しました:', error);
        alert('保存に失敗しました');
    }
}

async function editMemo(id) {
    try {
        const response = await fetch(`/api/memos/${id}`);
        const memo = await response.json();
        
        document.getElementById('memo-id').value = memo.id;
        document.getElementById('memo-title').value = memo.title;
        document.getElementById('memo-content').value = memo.content;
        document.getElementById('form-title').textContent = 'メモを編集';
        document.getElementById('cancel-btn').style.display = 'inline-block';
        
        isEditing = true;
        
        document.getElementById('memo-title').focus();
    } catch (error) {
        console.error('メモの取得に失敗しました:', error);
        alert('メモの取得に失敗しました');
    }
}

async function deleteMemo(id) {
    if (!confirm('このメモを削除しますか？')) {
        return;
    }
    
    try {
        const response = await fetch(`/api/memos/${id}`, {
            method: 'DELETE'
        });
        
        if (response.ok) {
            loadMemos();
        } else {
            alert('削除に失敗しました');
        }
    } catch (error) {
        console.error('削除に失敗しました:', error);
        alert('削除に失敗しました');
    }
}

function cancelEdit() {
    clearForm();
}

function clearForm() {
    document.getElementById('memo-id').value = '';
    document.getElementById('memo-title').value = '';
    document.getElementById('memo-content').value = '';
    document.getElementById('form-title').textContent = '新しいメモ';
    document.getElementById('cancel-btn').style.display = 'none';
    isEditing = false;
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString('ja-JP', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
}