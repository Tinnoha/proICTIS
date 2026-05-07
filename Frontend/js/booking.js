// ==========================================
// BOOKING & QR-CODE (исправленная версия)
// ==========================================

document.addEventListener('DOMContentLoaded', function() {
    // --- 1. Проверка авторизации и скрытие секции для неавторизованных ---
    const authToken = localStorage.getItem('auth_token');
    if (!authToken) {
        const bookingSection = document.getElementById('booking');
        if (bookingSection) bookingSection.style.display = 'none';
        const bookingNavLink = document.querySelector('a[href="#booking"]');
        if (bookingNavLink) bookingNavLink.style.display = 'none';
        return;
    }

    // Базовый URL для API (из конфигурации, доступной после api.js)
    const API_BASE = (typeof API_CONFIG !== 'undefined' && API_CONFIG.BASE_URL) 
                     ? API_CONFIG.BASE_URL 
                     : 'http://localhost:8080';

    // --- 2. Функция получения VK-ссылки с бэкенда ---
    async function fetchVkLink() {
        const currentUser = getCurrentUser();
        if (!currentUser || !currentUser.Id) {
            console.error('❌ Не удалось получить ID пользователя');
            return null;
        }

        try {
            console.log('🔍 Запрос VK-ссылки для user_id:', currentUser.Id);
            const response = await fetch(`${API_BASE}/User/vk`, {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    // Если бэкенд требует авторизацию, можно добавить токен (но в текущей реализации не требуется)
                },
                body: JSON.stringify({ user_id: currentUser.Id })
            });

            if (!response.ok) {
                const errText = await response.text();
                throw new Error(`Ошибка ${response.status}: ${errText}`);
            }

            const data = await response.json();
            console.log('✅ Получена VK-ссылка:', data.link);
            return data.link || null;
        } catch (error) {
            console.error('❌ Ошибка при получении VK-ссылки:', error);
            // alert('Не удалось получить ссылку для бронирования. Попробуйте позже.');
            return null;
        }
    }

    // --- 3. Вспомогательные функции ---
    function openInNewTab(url) {
        if (url) window.open(url, '_blank', 'noopener,noreferrer');
    }

    // --- 4. Обработчик текстовой ссылки «бота» ---
    const botLink = document.querySelector('.dynamic-link');
    if (botLink) {
        botLink.href = '#'; // убираем старую прямую ссылку

        botLink.addEventListener('click', async function(e) {
            e.preventDefault();
            const link = await fetchVkLink();
            if (link) {
                openInNewTab(link);
            } else {
                alert('Не удалось получить ссылку. Попробуйте позже.');
            }
        });
    }

    // --- 5. QR-код ---
    const qrImage = document.querySelector('[data-qr-image]');
    if (qrImage) {
        // Функция для установки QR-кода в img
        function setQRCode(url) {
            if (typeof QRCode === 'undefined') {
                console.error('Библиотека QRCode не загружена');
                qrImage.alt = 'QR-код недоступен';
                return;
            }

            if (!url) {
                console.error('Нет URL для генерации QR-кода');
                qrImage.alt = 'QR-код (нет ссылки)';
                return;
            }

            QRCode.toDataURL(url, {
                width: 200,
                margin: 2,
                color: { dark: '#000000', light: '#ffffff' }
            }, function(err, dataUrl) {
                if (err) {
                    console.error('Ошибка генерации QR-кода:', err);
                    qrImage.alt = 'Ошибка генерации QR';
                    return;
                }
                qrImage.src = dataUrl;
                qrImage.alt = 'QR-код для бронирования';
                console.log('✅ QR-код обновлён');
            });
        }

        // Асинхронная инициализация: получаем ссылку и генерируем QR
        (async function initQR() {
            const link = await fetchVkLink();
            if (link) {
                setQRCode(link);
                // Клик по QR-изображению тоже идёт по актуальной ссылке
                qrImage.addEventListener('click', () => openInNewTab(link));
            } else {
                qrImage.alt = 'QR-код (не удалось загрузить ссылку)';
                // Показываем заглушку (пустой QR)
                setQRCode('about:blank');
            }
        })();
    }
});