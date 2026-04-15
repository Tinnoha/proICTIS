// ==========================================
// EQUIPMENT DATA (API версия)
// ==========================================
const API_URL = 'http://localhost:8080';
let equipmentData = [];
let equipmentTypes = [];

// Вспомогательные функции для работы с авторизацией
function getAuthToken() {
    return localStorage.getItem('auth_token');
}

function getCurrentUser() {
    const user = localStorage.getItem('current_user');
    return user ? JSON.parse(user) : null;
}

// Загрузка оборудования из API
// Загрузка оборудования из API
async function loadEquipmentData() {
    try {
        const response = await fetch(`${API_URL}/Equipment`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        
        console.log('📡 Raw data from API:', data);
        
        // Преобразуем данные из формата API в формат фронтенда
        equipmentData = data.map(item => {
            let imageUrl = item.photo_url || item.PhotoURL || '';
            
            // 🔧 Если это локальный путь (/static/...), добавляем API_URL
            if (imageUrl && imageUrl.startsWith('/static/')) {
                imageUrl = `${API_URL}${imageUrl}`;
            }
            
            return {
                id: item.id || item.Id,
                title: item.name || item.Name,
                category: item.type || item.type_of_equipment || item.TypeOfEquipment,
                description: item.description || item.Description || '',
                image: imageUrl,  // ✅ Теперь это полный URL
                auditory: item.auditory || item.Auditory || '',
                isActive: item.is_active !== undefined ? item.is_active : item.IsActive
            };
        });
        
        console.log('✅ Загружено оборудования:', equipmentData.length);
        console.log('📋 Processed data:', equipmentData);
        return equipmentData;
    } catch (error) {
        console.error('❌ Ошибка загрузки оборудования:', error);
        equipmentData = [];
        return [];
    }
}

// Загрузка типов оборудования
async function loadEquipmentTypes() {
    try {
        const response = await fetch(`${API_URL}/Types`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        
        // Поддерживаем оба формата (lowercase и PascalCase)
        equipmentTypes = data.map(item => ({
            id: item.Id || item.id,
            name: item.Name || item.name
        }));
        
        console.log('✅ Загружено типов:', equipmentTypes.length);
        return equipmentTypes;
    } catch (error) {
        console.error('❌ Ошибка загрузки типов:', error);
        equipmentTypes = [];
        return [];
    }
}

// Получить все данные
function getEquipmentData() {
    return equipmentData;
}

// Получить типы
function getEquipmentTypesData() {
    return equipmentTypes;
}

// Добавить оборудование
async function addEquipment(item) {
    // Получаем текущего авторизованного пользователя
    const user = getCurrentUser();
    if (!user) {
        alert('Пользователь не авторизован');
        return false;
    }

    const payload = {
        admin_id: user.Id,  // ✅ Используем реальный ID пользователя
        tovars: [{
            name: item.title,
            type: item.category,
            auditory: item.auditory || '',
            description: item.description || '',
            photo_url: item.image || '',
            is_active: true
        }]
    };

    try {
        const response = await fetch(`${API_URL}/Equipment`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
                // 🔧 Убрали Cache-Control - вызывал CORS ошибку
            },
            body: JSON.stringify(payload)
        });
        
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Ошибка ${response.status}: ${errorText}`);
        }
        
        const result = await response.json();
        console.log('✅ Оборудование добавлено:', result);
        
        await loadEquipmentData();
        return true;
    } catch (error) {
        console.error('❌ Ошибка добавления оборудования:', error);
        alert(`Не удалось добавить оборудование: ${error.message}`);
        return false;
    }
}

// Удалить оборудование
async function deleteEquipment(id) {
    // Получаем текущего авторизованного пользователя
    const user = getCurrentUser();
    if (!user) {
        alert('Пользователь не авторизован');
        return false;
    }

    // Проверяем роль
    if (user.Role !== 'Admin' && user.Role !== 'Super_Admin') {
        alert('Только администраторы могут удалять оборудование');
        return false;
    }

    const payload = {
        admin_id: user.Id  // ✅ Используем реальный ID пользователя
    };

    try {
        const response = await fetch(`${API_URL}/Equipment/${id}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
                // 🔧 Убрали Cache-Control - вызывал CORS ошибку
            },
            body: JSON.stringify(payload)
        });
        
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Ошибка ${response.status}: ${errorText}`);
        }
        
        console.log('✅ Оборудование удалено');
        
        // Перезагружаем список
        await loadEquipmentData();
        return true;
    } catch (error) {
        console.error('❌ Ошибка удаления оборудования:', error);
        alert(`Не удалось удалить оборудование: ${error.message}`);
        return false;
    }
}

function clearEquipmentData() {
    equipmentData = [];
    equipmentTypes = [];
    return equipmentData;
}

document.addEventListener('DOMContentLoaded', async () => {
    console.log('🔄 Инициализация Equipment API...');
    await loadEquipmentData();
    await loadEquipmentTypes();
    console.log('✅ Equipment API готов к работе');
});