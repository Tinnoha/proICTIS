// ==========================================
// API CLIENT для Equipment Booking System
// ==========================================

const API_CONFIG = {
    BASE_URL: 'http://localhost:8080',
    TIMEOUT: 10000,
};

// ==========================================
// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
// ==========================================

function getAuthToken() {
    return localStorage.getItem('auth_token');
}

function setAuthToken(token) {
    localStorage.setItem('auth_token', token);
}

function getCurrentUser() {
    const user = localStorage.getItem('current_user');
    return user ? JSON.parse(user) : null;
}

function setCurrentUser(user) {
    localStorage.setItem('current_user', JSON.stringify(user));
}

function getHeaders(includeAuth = true) {
    const headers = { 'Content-Type': 'application/json' };
    if (includeAuth) {
        const token = getAuthToken();
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }
    }
    return headers;
}

async function handleResponse(response) {
    if (response.status === 204) return null;
    const data = await response.json();
    if (!response.ok) {
        const error = new Error(data.Error || 'Ошибка сервера');
        error.status = response.status;
        error.data = data;
        throw error;
    }
    return data;
}

async function fetchWithTimeout(url, options = {}, timeout = API_CONFIG.TIMEOUT) {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), timeout);
    try {
        const response = await fetch(url, { ...options, signal: controller.signal });
        clearTimeout(timeoutId);
        return response;
    } catch (error) {
        clearTimeout(timeoutId);
        if (error.name === 'AbortError') {
            throw new Error('Превышено время ожидания ответа от сервера');
        }
        throw error;
    }
}

// ==========================================
// AUTH API
// ==========================================

function loginWithYandex() {
    window.location.href = `${API_CONFIG.BASE_URL}/Regist`;
}

async function handleOAuthCallback() {
    try {
        const urlParams = new URLSearchParams(window.location.search);
        const token = urlParams.get('token');
        const email = urlParams.get('email');
        const firstName = urlParams.get('firstName');
        const lastName = urlParams.get('lastName');
        const avatar = urlParams.get('avatar');
        const role = urlParams.get('role');
        
        if (token) {
            setAuthToken(token);
            
            // 🔧 ОБРАБОТКА АВАТАРА ЯНДЕКСА
            let avatarUrl = decodeURIComponent(avatar || '');
            
            if (avatarUrl && avatarUrl !== 'null' && avatarUrl !== '') {
                // Если аватар есть, но это не полный URL - добавляем базовый URL Яндекса
                if (!avatarUrl.startsWith('http')) {
                    // Яндекс возвращает относительный путь, добавляем базовый URL
                    avatarUrl = 'https://avatars.yandex.net/get-yapic/' + avatarUrl + '/islands-200';
                }
            } else {
                // Если аватара нет - генерируем заглушку с инициалами
                avatarUrl = 'https://ui-avatars.com/api/?name=' + 
                           encodeURIComponent(firstName || 'Д') + 
                           '+' + encodeURIComponent(lastName || 'Н') + 
                           '&background=9d00d6&color=fff&size=100';
            }
            
            const user = {
                Id: token,
                Email: decodeURIComponent(email || ''),
                FirstName: decodeURIComponent(firstName || ''),
                SecondName: decodeURIComponent(lastName || ''),
                AvatarURL: avatarUrl,  // ✅ Теперь это полный URL
                Role: decodeURIComponent(role || 'student')
            };
            setCurrentUser(user);
            console.log('✅ Успешная авторизация:', user);
            return user;
        }
        
        const storedUser = localStorage.getItem('current_user');
        if (storedUser) {
            return JSON.parse(storedUser);
        }
        
        return null;
    } catch (error) {
        console.error('Ошибка обработки callback:', error);
        return null;
    }
}

function logout() {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('current_user');
    window.location.reload();
}

async function getCurrentUserInfo() {
    try {
        const token = getAuthToken();
        if (!token) return null;
        
        // Пробуем получить данные с сервера (если есть эндпоинт)
        // Или берём из localStorage
        const storedUser = localStorage.getItem('current_user');
        return storedUser ? JSON.parse(storedUser) : null;
    } catch (error) {
        console.error('Ошибка получения информации о пользователе:', error);
        return null;
    }
}

function isAuthenticated() {
    return !!getAuthToken();
}

function isAdmin() {
    const user = getCurrentUser();
    return user && (user.Role === 'admin' || user.Role === 'super_admin');
}

// ==========================================
// EQUIPMENT API
// ==========================================

async function getEquipment() {
    try {
        const response = await fetchWithTimeout(
            `${API_CONFIG.BASE_URL}/Equipment`,
            { headers: getHeaders(false) }
        );
        const data = await handleResponse(response);
        return data.map(item => ({
            id: item.Id,
            title: item.Name,
            category: item.TypeOfEquipment,
            description: item.Description || '',
            image: item.PhotoURL || '',
            auditory: item.Auditory || '',
            isActive: item.IsActive,
        }));
    } catch (error) {
        console.error('Ошибка получения оборудования:', error);
        alert('Не удалось загрузить оборудование. Проверьте подключение к серверу.');
        return [];
    }
}

async function getEquipmentById(id) {
    try {
        const response = await fetchWithTimeout(
            `${API_CONFIG.BASE_URL}/Equipment/id/${id}`,
            { headers: getHeaders(false) }
        );
        const data = await handleResponse(response);
        return {
            id: data.Id,
            title: data.Name,
            category: data.TypeOfEquipment,
            description: data.Description || '',
            image: data.PhotoURL || '',
            auditory: data.Auditory || '',
            isActive: data.IsActive,
        };
    } catch (error) {
        console.error('Ошибка получения оборудования:', error);
        throw error;
    }
}

async function createEquipment(equipmentData) {
    try {
        const user = getCurrentUser();
        if (!user || (user.Role !== 'admin' && user.Role !== 'super_admin')) {
            throw new Error('Только администраторы могут добавлять оборудование');
        }
        
        const requestBody = {
            admin_id: user.Id,
            tovars: [{
                name: equipmentData.title,
                type: equipmentData.category,
                auditory: equipmentData.auditory || '',
                description: equipmentData.description || '',
                photo_url: equipmentData.image || '',
                is_active: true,
            }]
        };
        
        const response = await fetchWithTimeout(
            `${API_CONFIG.BASE_URL}/Equipment`,
            {
                method: 'POST',
                headers: getHeaders(true),
                body: JSON.stringify(requestBody),
            }
        );
        
        await handleResponse(response);
        return true;
    } catch (error) {
        console.error('Ошибка создания оборудования:', error);
        alert(`Не удалось добавить оборудование: ${error.message}`);
        throw error;
    }
}

async function deleteEquipmentById(id) {
    try {
        const user = getCurrentUser();
        if (!user || (user.Role !== 'admin' && user.Role !== 'super_admin')) {
            throw new Error('Только администраторы могут удалять оборудование');
        }
        
        const requestBody = { admin_id: user.Id };
        
        const response = await fetchWithTimeout(
            `${API_CONFIG.BASE_URL}/Equipment/${id}`,
            {
                method: 'DELETE',
                headers: getHeaders(true),
                body: JSON.stringify(requestBody),
            }
        );
        
        await handleResponse(response);
        return true;
    } catch (error) {
        console.error('Ошибка удаления оборудования:', error);
        alert(`Не удалось удалить оборудование: ${error.message}`);
        throw error;
    }
}

async function updateEquipment(id, equipmentData) {
    try {
        const user = getCurrentUser();
        if (!user || (user.Role !== 'admin' && user.Role !== 'super_admin')) {
            throw new Error('Только администраторы могут редактировать оборудование');
        }
        
        const requestBody = {
            admin_id: user.Id,
            Equipment: {
                Id: id,
                Name: equipmentData.title,
                TypeOfEquipment: equipmentData.category,
                Auditory: equipmentData.auditory || '',
                Description: equipmentData.description || '',
                PhotoURL: equipmentData.image || '',
                IsActive: equipmentData.isActive !== undefined ? equipmentData.isActive : true,
            }
        };
        
        const response = await fetchWithTimeout(
            `${API_CONFIG.BASE_URL}/Equipment/${id}`,
            {
                method: 'PATCH',
                headers: getHeaders(true),
                body: JSON.stringify(requestBody),
            }
        );
        
        await handleResponse(response);
        return true;
    } catch (error) {
        console.error('Ошибка обновления оборудования:', error);
        alert(`Не удалось обновить оборудование: ${error.message}`);
        throw error;
    }
}

// ==========================================
// TYPES API
// ==========================================

async function getEquipmentTypes() {
    try {
        const response = await fetchWithTimeout(
            `${API_CONFIG.BASE_URL}/Types`,
            { headers: getHeaders(false) }
        );
        return await handleResponse(response);
    } catch (error) {
        console.error('Ошибка получения типов:', error);
        return [];
    }
}

// ==========================================
// ИНИЦИАЛИЗАЦИЯ
// ==========================================

function isOAuthCallback() {
    return window.location.pathname === '/callback' || 
           window.location.search.includes('token');
}

document.addEventListener('DOMContentLoaded', async () => {
    if (isOAuthCallback()) {
        await handleOAuthCallback();
        if (window.history.replaceState) {
            window.history.replaceState({}, document.title, window.location.pathname);
        }
    }
    console.log('API Client initialized');
    console.log('Backend URL:', API_CONFIG.BASE_URL);
});