document.addEventListener('DOMContentLoaded', function() {
    const authToken = localStorage.getItem('auth_token');
    
    // Если пользователь не авторизован - скрываем секцию бронирования и ссылку в меню
    if (!authToken) {
        const bookingSection = document.getElementById('booking');
        if (bookingSection) {
            bookingSection.style.display = 'none';
        }
        
        const bookingNavLink = document.querySelector('a[href="#booking"]');
        if (bookingNavLink) {
            bookingNavLink.style.display = 'none';
        }
        
        // Завершаем выполнение, так как дальше ничего не нужно
        return;
    }
    
    const baseUrl = 'https://vk.com/write-237660555?ref=';
    
    // Функция получения финального URL
    function getFinalUrl() {
        const token = localStorage.getItem('auth_token');
        return token ? baseUrl + encodeURIComponent(token) : baseUrl;
    }
    
    // Функция открытия ссылки в новой вкладке
    function openInNewTab(url) {
        window.open(url, '_blank', 'noopener,noreferrer');
    }
    
    // --- Обработчик для текстовой ссылки ---
    const botLink = document.querySelector('.dynamic-link');
    if (botLink) {
        botLink.href = getFinalUrl(); // для отображения URL при наведении
        
        botLink.addEventListener('click', function(e) {
            e.preventDefault();
            const token = localStorage.getItem('auth_token');
            if (!token) {
                alert('Для бронирования необходимо авторизоваться');
                return;
            }
            const finalUrl = baseUrl + encodeURIComponent(token);
            openInNewTab(finalUrl);
        });
    }
    
    // --- Генерация QR-кода через библиотеку QRCode ---
    const qrImage = document.querySelector('[data-qr-image]');
    if (qrImage) {
        const token = localStorage.getItem('auth_token');
        
        // Функция для установки QR-кода в img
        function setQRCode(url) {
            if (typeof QRCode === 'undefined') {
                console.error('Библиотека QRCode не загружена');
                qrImage.alt = 'QR-код недоступен';
                return;
            }
            
            QRCode.toDataURL(url, {
                width: 200,
                margin: 2,
                color: {
                    dark: '#000000',
                    light: '#ffffff'
                }
            }, function(err, dataUrl) {
                if (err) {
                    console.error('Ошибка генерации QR-кода:', err);
                    qrImage.alt = 'Ошибка генерации QR';
                    return;
                }
                qrImage.src = dataUrl;
                qrImage.alt = 'QR-код для бронирования';
            });
        }
        
        if (token) {
            const finalUrl = baseUrl + encodeURIComponent(token);
            setQRCode(finalUrl);
        } else {
            setQRCode(baseUrl);
            qrImage.alt = 'QR-код (требуется авторизация)';
        }
        
        // Клик по QR-изображению
        qrImage.addEventListener('click', function() {
            const token = localStorage.getItem('auth_token');
            if (!token) {
                alert('Для бронирования необходимо авторизоваться');
                return;
            }
            const url = baseUrl + encodeURIComponent(token);
            openInNewTab(url);
        });
    }
});