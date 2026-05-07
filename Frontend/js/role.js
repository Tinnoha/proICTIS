// Файл: role.js - Управление ролями пользователей
document.addEventListener('DOMContentLoaded', () => {
    // Проверка прав (только Super_Admin)
    const user = (typeof getCurrentUser === 'function') ? getCurrentUser() : null;
    if (!user || user.Role !== 'Super_Admin') return;

    const API_BASE_URL = 'http://localhost:8080';
    
    // Элементы
    const modal = document.getElementById('grantRoleModal');
    const tableBody = document.getElementById('usersTableBody');
    const selectAllCheckbox = document.getElementById('selectAllUsers');
    const roleSelect = document.getElementById('newRoleSelect');
    const confirmBtn = document.getElementById('confirmGrantRoleBtn');
    
    const selectedUsers = new Set();

    // Добавляем кнопку в выпадающее меню
    const userDropdown = document.querySelector('.user-dropdown-menu');
    if (userDropdown) {
        const divider = document.createElement('div');
        divider.className = 'dropdown-divider';
        
        const roleBtn = document.createElement('button');
        roleBtn.className = 'dropdown-role-btn';
        roleBtn.id = 'grantRoleBtn';
        roleBtn.innerHTML = `
            <i class="fas fa-user-plus"></i>
            <span>Выдать роль</span>
        `;
        
        // Вставляем перед кнопкой "Выйти"
        const logoutBtn = userDropdown.querySelector('.dropdown-logout');
        if (logoutBtn) {
            userDropdown.insertBefore(divider, logoutBtn);
            userDropdown.insertBefore(roleBtn, logoutBtn);
        }

        // Стили для кнопки
        const btnStyle = document.createElement('style');
        btnStyle.textContent = `
            .dropdown-role-btn {
                width: 100%;
                padding: 15px 20px;
                background: none;
                border: none;
                display: flex;
                align-items: center;
                gap: 12px;
                cursor: pointer;
                transition: all 0.3s ease;
                font-size: 0.95rem;
                color: #333;
                font-family: 'Inter', sans-serif;
                font-weight: 500;
            }
            .dropdown-role-btn:hover {
                background: rgba(157, 0, 214, 0.1);
                color: #9d00d6;
            }
            .dropdown-role-btn i {
                font-size: 1.1rem;
                color: #9d00d6;
            }
        `;
        document.head.appendChild(btnStyle);

        // Открытие модалки
        roleBtn.addEventListener('click', () => {
            modal.classList.add('active');
            loadUsers();
        });
    }

    // Закрытие модалки
    modal.querySelector('.modal-close').addEventListener('click', () => {
        modal.classList.remove('active');
        selectedUsers.clear();
        updateSelectAll();
    });
    
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            modal.classList.remove('active');
            selectedUsers.clear();
            updateSelectAll();
        }
    });

    // Загрузка пользователей
    async function loadUsers() {
        tableBody.innerHTML = '<tr><td colspan="4" style="text-align: center; padding: 30px; color: #999;">Загрузка...</td></tr>';
        
        try {
            const response = await fetch(`${API_BASE_URL}/User`, {
                method: 'POST',  // ← POST
                headers: { 
                    'Content-Type': 'application/json' 
                },
                body: JSON.stringify({
                    admin_id: user.Id  // ← UUID супер-админа в теле
                })
            });
            
            if (!response.ok) throw new Error('Ошибка загрузки');
            
            const users = await response.json();
            tableBody.innerHTML = '';
            
            if (users.length === 0) {
                tableBody.innerHTML = '<tr><td colspan="4" style="text-align: center; padding: 30px; color: #999;">Нет пользователей</td></tr>';
                return;
            }
            
            users.forEach(u => {
                const tr = document.createElement('tr');
                const fullName = `${u.first_name || ''} ${u.second_name || ''}`.trim() || 'Без имени';
                const roleName = getRoleName(u.role);
                const roleClass = u.role.toLowerCase();
                
                tr.innerHTML = `
                    <td><input type="checkbox" class="user-checkbox" data-id="${u.id}"></td>
                    <td>${fullName}</td>
                    <td>${u.email || '—'}</td>
                    <td><span class="role-badge ${roleClass}">${roleName}</span></td>
                `;
                
                const checkbox = tr.querySelector('.user-checkbox');
                checkbox.addEventListener('change', (e) => {
                    if (e.target.checked) {
                        selectedUsers.add(u.id);
                        tr.classList.add('selected');
                    } else {
                        selectedUsers.delete(u.id);
                        tr.classList.remove('selected');
                    }
                    updateSelectAll();
                });
                
                tableBody.appendChild(tr);
            });
            
        } catch (error) {
            console.error('Ошибка:', error);
            tableBody.innerHTML = '<tr><td colspan="4" style="text-align: center; padding: 30px; color: #ff4757;">Ошибка загрузки</td></tr>';
        }
    }

    // Выбор всех
    selectAllCheckbox.addEventListener('change', (e) => {
        const checkboxes = tableBody.querySelectorAll('.user-checkbox');
        checkboxes.forEach(cb => {
            cb.checked = e.target.checked;
            const userId = cb.dataset.id;
            if (e.target.checked) {
                selectedUsers.add(userId);
                cb.closest('tr').classList.add('selected');
            } else {
                selectedUsers.delete(userId);
                cb.closest('tr').classList.remove('selected');
            }
        });
    });

    function updateSelectAll() {
        const checkboxes = tableBody.querySelectorAll('.user-checkbox');
        const checkedCount = tableBody.querySelectorAll('.user-checkbox:checked').length;
        selectAllCheckbox.checked = checkedCount > 0 && checkedCount === checkboxes.length;
        selectAllCheckbox.indeterminate = checkedCount > 0 && checkedCount < checkboxes.length;
    }

    // Подтверждение
    confirmBtn.addEventListener('click', async () => {
        if (selectedUsers.size === 0) {
            alert('Выберите хотя бы одного пользователя!');
            return;
        }
        
        const newRole = roleSelect.value;
        if (!newRole) {
            alert('Выберите роль!');
            return;
        }
        
        const roleText = {
            'Super_Admin': 'Главный администратор',
            'Admin': 'Администратор',
            'Student': 'Студент'
        }[newRole];
        
        if (!confirm(`Выдать роль "${roleText}" ${selectedUsers.size} пользователю(ям)?`)) {
            return;
        }
        
        confirmBtn.disabled = true;
        confirmBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> <span>Обработка...</span>';
        
        let successCount = 0;
        let errorCount = 0;
        
        // Отправляем запросы параллельно
        const promises = Array.from(selectedUsers).map(userId => 
            fetch(`${API_BASE_URL}/User/admin`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    admin_id: user.id,
                    target_user_id: userId,
                    new_role: newRole
                })
            }).then(res => res.ok ? 'success' : 'error')
              .catch(() => 'error')
        );
        
        const results = await Promise.allSettled(promises);
        results.forEach(r => r.value === 'success' ? successCount++ : errorCount++);
        
        alert(`✅ Успешно: ${successCount}\n❌ Ошибок: ${errorCount}`);
        
        modal.classList.remove('active');
        selectedUsers.clear();
        updateSelectAll();
        confirmBtn.disabled = false;
        confirmBtn.innerHTML = '<i class="fas fa-user-shield"></i> <span>Выдать роль</span>';
        roleSelect.value = '';
    });

    function getRoleName(role) {
        const roles = {
            'Student': 'Студент',
            'Admin': 'Администратор',
            'Super_Admin': 'Главный администратор'
        };
        return roles[role] || role;
    }
});