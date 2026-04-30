// ==========================================
// МОБИЛЬНОЕ МЕНЮ
// ==========================================
const hamburger = document.querySelector('.hamburger');
const navMenu = document.querySelector('.nav-menu');
if (hamburger && navMenu) {
    hamburger.addEventListener('click', () => {
        hamburger.classList.toggle('active');
        navMenu.classList.toggle('active');
    });
}

// ==========================================
// ПЛАВНАЯ ПРОКРУТКА
// ==========================================
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault();
        const targetId = this.getAttribute('href');
        if (targetId === '#') return;
        
        const targetElement = document.querySelector(targetId);
        if (targetElement) {
            const navHeight = document.querySelector('.navbar').offsetHeight;
            const targetPosition = targetElement.offsetTop - navHeight - 30;
            
            window.scrollTo({
                top: targetPosition,
                behavior: 'smooth'
            });
            
            if (hamburger && navMenu) {
                hamburger.classList.remove('active');
                navMenu.classList.remove('active');
            }
        }
    });
});

// ==========================================
// АКТИВНАЯ ССЫЛКА ПРИ СКРОЛЛЕ
// ==========================================
const sections = document.querySelectorAll('header, section, footer');
const navLinks = document.querySelectorAll('.nav-link');
const footer = document.querySelector('footer');
window.addEventListener('scroll', () => {
    let current = '';
    const scrollPosition = window.scrollY + window.innerHeight / 2;
    const footerTop = footer.offsetTop;
    
    if (window.scrollY + window.innerHeight >= footerTop - 100) {
        current = 'contacts';
    } else {
        sections.forEach(section => {
            const sectionTop = section.offsetTop;
            const sectionHeight = section.clientHeight;
            const sectionId = section.getAttribute('id');
            
            if (scrollPosition >= sectionTop && scrollPosition < sectionTop + sectionHeight) {
                current = sectionId;
            }
        });
    }
    
    navLinks.forEach(link => {
        link.classList.remove('active');
        if (link.getAttribute('href').slice(1) === current) {
            link.classList.add('active');
        }
    });
});

// ==========================================
// РЕНДЕРИНГ ОБОРУДОВАНИЯ
// ==========================================
const equipmentGrid = document.querySelector('.equipment-grid');
const addCard = document.querySelector('.add-card');

async function renderEquipment() {
    if (!equipmentGrid) return;
    if (equipmentData.length === 0) {
        await loadEquipmentData();
    }

    // Очищаем только карточки оборудования, оставляя add-card
    const existingCards = equipmentGrid.querySelectorAll('.equip-card:not(.add-card)');
    existingCards.forEach(card => card.remove());

    // 🔧 1. Гарантируем, что карточка "Добавить" всегда стоит первой
    if (addCard && addCard.parentNode === equipmentGrid) {
        equipmentGrid.insertBefore(addCard, equipmentGrid.firstChild);
    }

    const equipment = getEquipmentData();

    equipment.forEach(item => {
        const card = createEquipmentCard(item);
        // 🔧 2. Добавляем карточки оборудования в конец сетки (строго после add-card)
        equipmentGrid.appendChild(card);
    });

    addViewModalListeners();
}

function createEquipmentCard(equipment) {
    const card = document.createElement('div');
    card.className = 'equip-card';
    card.dataset.id = equipment.id;
    card.dataset.category = equipment.category;
    card.style.cursor = 'pointer';
    
    const imageUrl = equipment.image && equipment.image !== '' 
        ? equipment.image 
        : 'https://placehold.co/300x200/9d00d6/ffffff?text=No+Image';
    
    card.innerHTML = `
        <div class="equip-img">
            <img src="${imageUrl}" alt="${equipment.title || 'Equipment'}" onerror="this.src='https://placehold.co/300x200/9d00d6/ffffff?text=No+Image'">
        </div>
        <div class="equip-info">
            <h3>${equipment.title || 'Без названия'}</h3>
        </div>
    `;
    return card;
}

function addViewModalListeners() {
    const cards = document.querySelectorAll('.equip-card:not(.add-card)');
    cards.forEach(card => {
        // 🔧 Удаляем старые обработчики (клонированием)
        const newCard = card.cloneNode(true);
        card.parentNode.replaceChild(newCard, card);
        
        newCard.addEventListener('click', () => {
            // 🔧 Если в режиме удаления - только выделяем, НЕ открываем модальное окно
            if (deleteMode) {
                handleDeleteModeClick({ 
                    currentTarget: newCard, 
                    preventDefault: () => {}, 
                    stopPropagation: () => {} 
                });
                return;
            }
            
            // Обычное поведение (не в режиме удаления)
            const id = newCard.dataset.id;
            const equipment = getEquipmentData().find(e => e.id === id);
            if (equipment) {
                openViewModal(equipment);
            }
        });
    });
}

// ==========================================
// МОДАЛЬНОЕ ОКНО ПРОСМОТРА
// ==========================================
let viewModal = document.querySelector('.view-modal');
let currentEquipmentId = null;

if (!viewModal) {
    const modalHTML = `<div class="view-modal">
        <div class="view-modal-content">
            <button class="view-modal-close">
                <i class="fas fa-times"></i>
            </button>
            <div class="view-modal-body">
                <div class="view-modal-image">
                    <img src="" alt="">
                </div>
                <div class="view-modal-info">
                    <h2 class="view-modal-title"></h2>
                    <div class="view-modal-category"></div>
                    <div class="view-modal-description"></div>
                    <button class="delete-btn">
                        <i class="fas fa-trash"></i> Удалить
                    </button>
                </div>
            </div>
        </div>
    </div>`;
    document.body.insertAdjacentHTML('beforeend', modalHTML);
    viewModal = document.querySelector('.view-modal');
}

function openViewModal(equipment) {
    if (!viewModal) return;
    currentEquipmentId = equipment.id;
    
    const img = viewModal.querySelector('.view-modal-image img');
    const title = viewModal.querySelector('.view-modal-title');
    const category = viewModal.querySelector('.view-modal-category');
    const description = viewModal.querySelector('.view-modal-description');
    
    let deleteBtn = viewModal.querySelector('.delete-btn');
    if (deleteBtn) {
        const newDeleteBtn = deleteBtn.cloneNode(true);
        deleteBtn.parentNode.replaceChild(newDeleteBtn, deleteBtn);
        deleteBtn = newDeleteBtn;
        
        const user = getCurrentUser();
        const userRole = (user?.Role || '').toLowerCase();
        if (user && (userRole === 'admin' || userRole === 'super_admin')) {
            deleteBtn.style.display = 'flex';
            deleteBtn.addEventListener('click', () => {
                openConfirmModal(currentEquipmentId);
            });
        } else {
            deleteBtn.style.display = 'none';
        }
    }
    
    const imageUrl = equipment.image && equipment.image !== '' 
        ? equipment.image 
        : 'https://placehold.co/400x400/9d00d6/ffffff?text=No+Image';
    
    img.src = imageUrl;
    img.onerror = function() {
        this.src = 'https://placehold.co/400x400/9d00d6/ffffff?text=No+Image';
    };
    
    img.alt = equipment.title || 'Equipment';
    title.textContent = equipment.title || 'Без названия';
    category.textContent = equipment.category || 'Без категории';
    
    const formattedDesc = equipment.description ? equipment.description.replace(/\n/g, '<br>') : 'Описание отсутствует';
    description.innerHTML = formattedDesc;
    
    viewModal.classList.add('active');
    document.body.style.overflow = 'hidden';
}

const closeViewModal = () => {
    if (!viewModal) return;
    viewModal.classList.remove('active');
    document.body.style.overflow = '';
};

if (viewModal) {
    const closeBtn = viewModal.querySelector('.view-modal-close');
    if (closeBtn) {
        closeBtn.addEventListener('click', closeViewModal);
    }
    
    viewModal.addEventListener('click', (e) => {
        if (e.target === viewModal) {
            closeViewModal();
        }
    });
}

// ==========================================
// МОДАЛЬНОЕ ОКНО ПОДТВЕРЖДЕНИЯ УДАЛЕНИЯ (одиночное)
// ==========================================
let confirmModal = document.querySelector('.confirm-modal');
if (!confirmModal) {
    const confirmHTML = `<div class="confirm-modal">
        <div class="confirm-modal-content">
            <div class="confirm-modal-icon">
                <i class="fas fa-exclamation-triangle"></i>
            </div>
            <h3 class="confirm-modal-title">Подтверждение удаления</h3>
            <p class="confirm-modal-text">Вы действительно хотите удалить эту карточку оборудования? Это действие нельзя отменить.</p>
            <div class="confirm-modal-buttons">
                <button class="confirm-btn cancel">Отмена</button>
                <button class="confirm-btn delete">Удалить</button>
            </div>
        </div>
    </div>`;
    document.body.insertAdjacentHTML('beforeend', confirmHTML);
    confirmModal = document.querySelector('.confirm-modal');
}

function openConfirmModal(equipmentId) {
    if (!confirmModal) return;
    currentEquipmentId = equipmentId;
    confirmModal.classList.add('active');
}

const closeConfirmModal = () => {
    if (!confirmModal) return;
    confirmModal.classList.remove('active');
};

if (confirmModal) {
    const cancelBtn = confirmModal.querySelector('.confirm-btn.cancel');
    const deleteBtn = confirmModal.querySelector('.confirm-btn.delete');
    
    if (cancelBtn) {
        cancelBtn.addEventListener('click', closeConfirmModal);
    }
    
    if (deleteBtn) {
        deleteBtn.addEventListener('click', async () => {
            if (currentEquipmentId !== null) {
                try {
                    const success = await deleteEquipment(currentEquipmentId);
                    
                    if (success) {
                        await new Promise(resolve => setTimeout(resolve, 500));
                        await renderEquipment();
                        
                        closeConfirmModal();
                        closeViewModal();
                        // alert('✅ Оборудование успешно удалено');
                        console.log('✅ Оборудование успешно удалено');
                    }
                } catch (error) {
                    // alert(`❌ Ошибка при удалении: ${error.message}`);
                    console.error('❌ Ошибка при удалении:', error);
                }
            }
        });
    }
    
    confirmModal.addEventListener('click', (e) => {
        if (e.target === confirmModal) {
            closeConfirmModal();
        }
    });
}

// ==========================================
// ЗАГРУЗКА КНОПОК КАТЕГОРИЙ ИЗ БД
// ==========================================
async function loadCategoryButtons() {
    try {
        const types = await getEquipmentTypesData();
        const categoriesContainer = document.querySelector('.categories');
        if (!categoriesContainer) return;
        
        categoriesContainer.innerHTML = '';
        
        // Добавляем кнопку "Все"
        const allButton = document.createElement('button');
        allButton.className = 'cat-btn active';
        allButton.textContent = 'Все';
        allButton.dataset.category = 'all';
        categoriesContainer.appendChild(allButton);
        
        // Добавляем кнопки из БД
        types.forEach(type => {
            const button = document.createElement('button');
            button.className = 'cat-btn';
            button.textContent = type.name || type.Name;
            button.dataset.category = type.name || type.Name;
            button.dataset.typeId = type.id || type.Id;
            categoriesContainer.appendChild(button);
        });
        
        // 🔧 ПРОВЕРКА ПРАВ: показываем кнопки только админам
        const user = getCurrentUser();
        const userRole = (user?.Role || '').toLowerCase();
        const isAdmin = userRole === 'admin' || userRole === 'super_admin';
         
        if (isAdmin) {
            // Кнопка добавления типа
            const addButton = document.createElement('button');
            addButton.className = 'cat-btn add-type-btn';
            addButton.innerHTML = '<i class="fas fa-plus"></i>';
            addButton.title = 'Добавить тип оборудования';
            categoriesContainer.appendChild(addButton);

            // Кнопка удаления типов (всегда видна админам)
            const deleteTypesBtn = document.createElement('button');
            deleteTypesBtn.className = 'cat-btn delete-types-btn';
            deleteTypesBtn.innerHTML = '<i class="fas fa-trash"></i>';
            deleteTypesBtn.title = 'Удалить типы/оборудование';
            categoriesContainer.appendChild(deleteTypesBtn);
            
            // Обработчик для кнопки удаления
            deleteTypesBtn.addEventListener('click', () => {
                if (deleteMode && selectedForDelete.length > 0) {
                    // 🔧 В режиме удаления + есть выбранные → подтверждаем удаление
                    openConfirmDeleteModal();
                } else if (deleteMode && selectedForDelete.length === 0) {
                    // 🔧 В режиме удаления, но ничего не выбрано
                    // alert('Пожалуйста, выберите объекты для удаления');
                    console.log('Пожалуйста, выберите объекты для удаления');
                } else {
                    // 🔧 Не в режиме удаления → открываем выбор объектов
                    openDeleteSelectionModal();
                }
            });
            
            // Обработчик для кнопки добавления типа
            addButton.addEventListener('click', () => {
                openAddTypeModal();
            });
        }
        
        // 🔧 Добавляем обработчики фильтрации ПРЯМО СЕЙЧАС (без клонирования!)
        addCategoryFilterListeners();
        console.log('✅ Загружено кнопок категорий:', types.length);
    } catch (error) {
        console.error('❌ Ошибка загрузки кнопок категорий:', error);
    }
}

function addCategoryFilterListeners() {
    // 🔧 Получаем актуальные кнопки КАЖДЫЙ РАЗ
    const catBtns = document.querySelectorAll('.cat-btn:not(.delete-types-btn):not(.add-type-btn):not(.cancel-delete-btn)');
    
    catBtns.forEach(btn => {
        // 🔧 Удаляем ВСЕ старые обработчики (клонированием)
        const newBtn = btn.cloneNode(true);
        btn.parentNode.replaceChild(newBtn, btn);
        
        // 🔧 Добавляем новый обработчик
        newBtn.addEventListener('click', () => {
            // 🔧 Если в режиме удаления - только выделяем, НЕ фильтруем
            if (deleteMode) {
                handleDeleteModeClick({ 
                    currentTarget: newBtn, 
                    preventDefault: () => {}, 
                    stopPropagation: () => {} 
                });
                return;
            }
            
            // 🔧 Обычная фильтрация (не в режиме удаления)
            // 🔧 Получаем АКТУАЛЬНЫЕ кнопки и снимаем выделение со ВСЕХ
            document.querySelectorAll('.cat-btn:not(.delete-types-btn):not(.add-type-btn):not(.cancel-delete-btn)')
                .forEach(b => b.classList.remove('active'));
            
            // 🔧 Добавляем выделение на нажатую кнопку
            newBtn.classList.add('active');
            
            const category = newBtn.dataset.category;
            const equipCards = document.querySelectorAll('.equip-card:not(.add-card)');
            
            equipCards.forEach(card => {
                const cardCategory = card.dataset.category;
                if (category === 'all' || cardCategory === category) {
                    card.style.display = 'flex';
                    card.style.animation = 'fadeIn 0.5s ease';
                } else {
                    card.style.display = 'none';
                }
            });
        });
    });
}

// ==========================================
// РЕЖИМ УДАЛЕНИЯ (МАССОВОЕ УДАЛЕНИЕ)
// ==========================================
let deleteMode = false;
let selectedForDelete = [];

function openDeleteSelectionModal() {
    const modal = document.querySelector('.delete-selection-modal');
    if (modal) {
        modal.classList.add('active');
        document.body.style.overflow = 'hidden';
    }
}

const closeDeleteSelectionModal = () => {
    const modal = document.querySelector('.delete-selection-modal');
    if (modal) {
        modal.classList.remove('active');
        document.body.style.overflow = '';
    }
};

function enterDeleteMode() {
    deleteMode = true;
    selectedForDelete = [];
    
    const categoriesContainer = document.querySelector('.categories');
    
    // 🔧 Добавляем ОТДЕЛЬНУЮ кнопку отмены (рядом с кнопкой удаления)
    const cancelDeleteBtn = document.createElement('button');
    cancelDeleteBtn.className = 'cat-btn cancel-delete-btn';
    cancelDeleteBtn.innerHTML = '<i class="fas fa-times"></i>';
    cancelDeleteBtn.title = 'Отменить удаление';
    cancelDeleteBtn.id = 'cancelDeleteModeBtn';
    
    // Вставляем после кнопки удаления
    const deleteTypesBtn = categoriesContainer.querySelector('.delete-types-btn');
    if (deleteTypesBtn && deleteTypesBtn.nextSibling) {
        categoriesContainer.insertBefore(cancelDeleteBtn, deleteTypesBtn.nextSibling);
    } else if (deleteTypesBtn) {
        categoriesContainer.appendChild(cancelDeleteBtn);
    }
    
    // Обработчик для кнопки отмены
    const cancelBtn = document.querySelector('#cancelDeleteModeBtn');
    if (cancelBtn) {
        cancelBtn.addEventListener('click', () => {
            openCancelDeleteConfirmModal();
        });
    }
    
    // 🔧 Пересоздаем карточки для удаления старых обработчиков
    const cards = document.querySelectorAll('.equip-card:not(.add-card)');
    cards.forEach(card => {
        const newCard = card.cloneNode(true);
        card.parentNode.replaceChild(newCard, card);
    });
    
    // 🔧 Пересоздаем кнопки категорий для удаления старых обработчиков
    const catBtns = document.querySelectorAll('.cat-btn:not(.delete-types-btn):not(.add-type-btn):not(.cancel-delete-btn)');
    catBtns.forEach(btn => {
        const newBtn = btn.cloneNode(true);
        btn.parentNode.replaceChild(newBtn, btn);
    });
    
    // 🔧 Добавляем обработчики для режима удаления
    addDeleteModeListeners();
    console.log('🗑️ Режим удаления активирован');
}

function exitDeleteMode() {
    deleteMode = false;
    
    // 🔧 Снимаем выделение со всех объектов
    document.querySelectorAll('.selected-for-delete').forEach(el => {
        el.classList.remove('selected-for-delete');
    });
    selectedForDelete = [];
    
    // 🔧 Удаляем кнопку отмены
    const cancelBtn = document.querySelector('#cancelDeleteModeBtn');
    if (cancelBtn) {
        cancelBtn.remove();
    }
    
    // 🔧 Пересоздаем карточки для удаления обработчиков режима удаления
    const cards = document.querySelectorAll('.equip-card:not(.add-card)');
    cards.forEach(card => {
        const newCard = card.cloneNode(true);
        card.parentNode.replaceChild(newCard, card);
    });
    
    // 🔧 Пересоздаем кнопки категорий для удаления обработчиков режима удаления
    const catBtns = document.querySelectorAll('.cat-btn:not(.delete-types-btn):not(.add-type-btn):not(.cancel-delete-btn)');
    catBtns.forEach(btn => {
        const newBtn = btn.cloneNode(true);
        btn.parentNode.replaceChild(newBtn, btn);
    });
    
    // 🔧 Восстанавливаем обычные обработчики
    addCategoryFilterListeners();
    addViewModalListeners();
    
    console.log('✅ Режим удаления деактивирован');
}

function addDeleteModeListeners() {
    // 🔧 Карточки оборудования - только выделение, без открытия модального окна
    document.querySelectorAll('.equip-card:not(.add-card)').forEach(card => {
        card.style.cursor = 'pointer';
        card.addEventListener('click', handleDeleteModeClick);
    });
    
    // 🔧 Кнопки категорий - только выделение, без фильтрации
    document.querySelectorAll('.cat-btn:not(.delete-types-btn):not(.add-type-btn):not(.cancel-delete-btn)').forEach(btn => {
        btn.style.cursor = 'pointer';
        btn.addEventListener('click', handleDeleteModeClick);
    });
}

function handleDeleteModeClick(e) {
    e.preventDefault();
    e.stopPropagation();
    
    const target = e.currentTarget;
    
    if (target.classList.contains('selected-for-delete')) {
        target.classList.remove('selected-for-delete');
        selectedForDelete = selectedForDelete.filter(item => item.element !== target);
    } else {
        target.classList.add('selected-for-delete');
        selectedForDelete.push({
            element: target,
            type: target.classList.contains('equip-card') ? 'equipment' : 'type',
            id: target.dataset.id || target.dataset.typeId,
            name: target.querySelector('h3')?.textContent || target.textContent
        });
    }
    
    console.log('📋 Выбрано объектов:', selectedForDelete.length);
    console.log('📋 Выбранные элементы:', selectedForDelete);
}

function openCancelDeleteConfirmModal() {
    const modal = document.createElement('div');
    modal.className = 'delete-selection-modal';
    modal.innerHTML = `<div class="delete-selection-modal-content">
        <div class="delete-selection-modal-icon">
            <i class="fas fa-question-circle"></i>
        </div>
        <h3 class="delete-selection-modal-title">Вы уверены, что хотите прекратить удаление?</h3>
        <p class="delete-selection-modal-text">Все выбранные объекты будут сняты с выделения.</p>
        <div class="delete-selection-modal-buttons">
            <button class="delete-selection-btn cancel">Отмена</button>
            <button class="delete-selection-btn ok">Да</button>
        </div>
    </div>`;
    document.body.appendChild(modal);
    
    setTimeout(() => modal.classList.add('active'), 10);
    
    const cancelBtn = modal.querySelector('.delete-selection-btn.cancel');
    const okBtn = modal.querySelector('.delete-selection-btn.ok');
    
    cancelBtn.addEventListener('click', () => {
        modal.classList.remove('active');
        setTimeout(() => modal.remove(), 300);
    });
    
    okBtn.addEventListener('click', () => {
        exitDeleteMode();
        modal.classList.remove('active');
        setTimeout(() => modal.remove(), 300);
    });
    
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            modal.classList.remove('active');
            setTimeout(() => modal.remove(), 300);
        }
    });
}

function openConfirmDeleteModal() {
    if (selectedForDelete.length === 0) {
        // alert('Пожалуйста, выберите объекты для удаления');
        console.log('Пожалуйста, выберите объекты для удаления');
        return;
    }
    
    const modal = document.createElement('div');
    modal.className = 'delete-selection-modal';
    modal.innerHTML = `
        <div class="delete-selection-modal-content">
            <div class="delete-selection-modal-icon">
                <i class="fas fa-exclamation-triangle"></i>
            </div>
            <h3 class="delete-selection-modal-title">Вы точно хотите удалить выбранные объекты?</h3>
            <p class="delete-selection-modal-text">
                Будет удалено объектов: ${selectedForDelete.length}<br>
                Это действие нельзя отменить.
            </p>
            <div class="delete-selection-modal-buttons">
                <button class="delete-selection-btn cancel">Отмена</button>
                <button class="delete-selection-btn ok">Удалить</button>
            </div>
        </div>
    `;
    document.body.appendChild(modal);
    
    setTimeout(() => modal.classList.add('active'), 10);
    
    const cancelBtn = modal.querySelector('.delete-selection-btn.cancel');
    const okBtn = modal.querySelector('.delete-selection-btn.ok');
    
    cancelBtn.addEventListener('click', () => {
        modal.classList.remove('active');
        setTimeout(() => modal.remove(), 300);
    });
    
    okBtn.addEventListener('click', async () => {
        await performDelete();
        // 🔧 НЕ выходим из режима удаления после удаления!
        modal.classList.remove('active');
        setTimeout(() => modal.remove(), 300);
    });
    
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            modal.classList.remove('active');
            setTimeout(() => modal.remove(), 300);
        }
    });
}

async function performDelete() {
    try {
        let successCount = 0;
        
        console.log('🗑️ [Mass] Начинаем удаление');
        console.log('📋 [Mass] Выбрано:', selectedForDelete.map(i => `${i.type}:${i.name || i.id}`));
        
        // 🔧 1. Сначала удаляем ВСЕ карточки оборудования
        const equipmentToDelete = selectedForDelete.filter(item => item.type === 'equipment');
        console.log('🔧 [Mass] Удаляем карточки:', equipmentToDelete.length);
        
        for (const item of equipmentToDelete) {
            console.log(`  → [Equip] Удаляем карточку: ${item.id}`);
            const success = await deleteEquipment(item.id);
            if (success) {
                successCount++;
                console.log(`  ✅ [Equip] Карточка ${item.id} удалена`);
            } else {
                console.log(`  ❌ [Equip] Карточка ${item.id} НЕ удалена`);
            }
        }
        
        // 🔧 2. ПРОВЕРЯЕМ через API, что карточки действительно удалились
        console.log('🔍 [Mass] Проверяем состояние БД...');
        const updatedEquipment = await loadEquipmentData();
        console.log('📊 [Mass] В БД осталось оборудования:', updatedEquipment.length);
        
        // 🔧 3. Потом удаляем типы (только если нет карточек с этим типом)
        const typesToDelete = selectedForDelete.filter(item => item.type === 'type');
        console.log('🔧 [Mass] Отправляем типы на удаление:', typesToDelete.length);

        for (const item of typesToDelete) {
            console.log(`  → [Type] Запрос на удаление типа: ${item.name}`);
            
            // УБРАЛИ ручную проверку. Бэкенд сам вернёт ошибку, если карточки есть.
            const success = await deleteEquipmentType(item.id);
            
            if (success) {
                successCount++;
                console.log(`  ✅ [Type] Тип "${item.name}" удалён`);
            } else {
                console.log(`  ⚠️ [Type] Тип "${item.name}" не удалён (модальное окно показано)`);
            }
        }
        
        console.log(`✅ [Mass] Итого удалено: ${successCount} из ${selectedForDelete.length}`);
        // alert(`✅ Удалено объектов: ${successCount} из ${selectedForDelete.length}`);
        
        exitDeleteMode();
        
        await loadEquipmentData();
        await loadEquipmentTypes();
        await loadCategoryButtons();
        await loadEquipmentTypesToDropdown();
        await renderEquipment();
        
    } catch (error) {
        console.error('❌ [Mass] Ошибка:', error);
        // alert('❌ Ошибка при удалении: ' + error.message);
    }
}

async function deleteEquipmentType(typeId) {
    const user = getCurrentUser();
    if (!user) {
        // alert('Пользователь не авторизован');
        console.log('Пользователь не авторизован');
        return false;
    }
    const payload = { admin_id: user.Id };

    try {
        const response = await fetch(`${API_URL}/Types/${typeId}`, {
            method: 'DELETE',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (!response.ok) {
            const errorText = await response.text();
            // Проверяем статус и ключевые слова ошибки из репозитория
            if (response.status >= 400 && (errorText.includes('существуют карточки') || errorText.includes('нельзя удалить тип'))) {
                showTypeDeleteErrorModal();
                return false; // Прерываем, не выбрасываем ошибку дальше
            }
            throw new Error(`Ошибка ${response.status}: ${errorText}`);
        }

        console.log('✅ Тип оборудования удалён');
        return true;
    } catch (error) {
        console.error('❌ Ошибка удаления типа:', error);
        // alert(`Не удалось удалить тип: ${error.message}`);
        return false;
    }
}

// 🔧 Красивое модальное окно для ошибки удаления типа
function showTypeDeleteErrorModal() {
    const modal = document.createElement('div');
    modal.className = 'delete-selection-modal';
    modal.innerHTML = `
        <div class="delete-selection-modal-content">
            <div class="delete-selection-modal-icon" style="color: #ff4757;">
                <i class="fas fa-exclamation-triangle" style="font-size: 4rem;"></i>
            </div>
            <h3 class="delete-selection-modal-title">Нельзя удалить выбранный тип оборудования</h3>
            <p class="delete-selection-modal-text">
                Существуют карточки с заданным типом!<br>
                Сначала удалите или измените всё оборудование этого типа.
            </p>
            <div class="delete-selection-modal-buttons">
                <button class="delete-selection-btn ok" style="background: linear-gradient(135deg, #9d00d6 0%, #048af1 100%);">
                    ОК
                </button>
            </div>
        </div>
    `;
    
    document.body.appendChild(modal);
    
    setTimeout(() => modal.classList.add('active'), 10);
    
    const okBtn = modal.querySelector('.delete-selection-btn.ok');
    okBtn.addEventListener('click', () => {
        modal.classList.remove('active');
        setTimeout(() => modal.remove(), 300);
    });
    
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            modal.classList.remove('active');
            setTimeout(() => modal.remove(), 300);
        }
    });
    
    // Закрытие по Escape
    const handleEscape = (e) => {
        if (e.key === 'Escape') {
            modal.classList.remove('active');
            setTimeout(() => modal.remove(), 300);
            document.removeEventListener('keydown', handleEscape);
        }
    };
    document.addEventListener('keydown', handleEscape);
}

// ==========================================
// МОДАЛЬНОЕ ОКНО ВЫБОРА УДАЛЕНИЯ (HTML)
// ==========================================
const deleteSelectionModalHTML = `<div class="delete-selection-modal">
    <div class="delete-selection-modal-content">
        <div class="delete-selection-modal-icon">
            <i class="fas fa-trash-alt"></i>
        </div>
        <h3 class="delete-selection-modal-title">Выберите объекты для удаления</h3>
        <p class="delete-selection-modal-text">
            Нажмите на карточки оборудования или кнопки типов, которые хотите удалить.<br>
            Выбранные объекты будут выделены.
        </p>
        <div class="delete-selection-modal-buttons">
            <button class="delete-selection-btn cancel">Отмена</button>
            <button class="delete-selection-btn ok">ОК</button>
        </div>
    </div>
</div>`;

// ==========================================
// МОДАЛЬНОЕ ОКНО ДОБАВЛЕНИЯ ТИПА ОБОРУДОВАНИЯ
// ==========================================
let addTypeModal = document.querySelector('.add-type-modal');
if (!addTypeModal) {
    const modalHTML = `<div class="add-type-modal">
        <div class="add-type-modal-content">
            <button class="add-type-modal-close">
                <i class="fas fa-times"></i>
            </button>
            <h2 class="add-type-modal-title">Добавить тип оборудования</h2>
            <div class="add-type-modal-body">
                <div class="form-group">
                    <label for="newEquipmentType">Название типа</label>
                    <input type="text" id="newEquipmentType" placeholder="Например: Лазерная резка">
                </div>
                <button class="add-type-submit-btn">
                    <i class="fas fa-plus"></i> Добавить тип
                </button>
            </div>
        </div>
    </div>`;
    document.body.insertAdjacentHTML('beforeend', modalHTML);
    addTypeModal = document.querySelector('.add-type-modal');
}

function openAddTypeModal() {
    if (!addTypeModal) return;
    addTypeModal.classList.add('active');
    document.body.style.overflow = 'hidden';
}

const closeAddTypeModal = () => {
    if (!addTypeModal) return;
    addTypeModal.classList.remove('active');
    document.body.style.overflow = '';
    const input = document.querySelector('#newEquipmentType');
    if (input) input.value = '';
};

if (addTypeModal) {
    const closeBtn = addTypeModal.querySelector('.add-type-modal-close');
    const submitBtn = addTypeModal.querySelector('.add-type-submit-btn');
    
    if (closeBtn) {
        closeBtn.addEventListener('click', closeAddTypeModal);
    }
    
    if (submitBtn) {
        submitBtn.addEventListener('click', async () => {
            const input = document.querySelector('#newEquipmentType');
            const typeName = input.value.trim();
            
            if (!typeName) {
                alert('Пожалуйста, введите название типа оборудования');
                // console.log('Пожалуйста, введите название типа оборудования');
                return;
            }
            
            const user = getCurrentUser();
            const userRole = (user?.Role || '').toLowerCase();
            if (userRole !== 'admin' && userRole !== 'super_admin') {
                // alert('Только администраторы могут добавлять типы оборудования');
                console.log('Только администраторы могут добавлять типы оборудования');
                closeAddTypeModal();
                return;
            }
            
            submitBtn.disabled = true;
            submitBtn.textContent = 'Добавление...';
            
            try {
                const success = await addEquipmentType(typeName);
                
                if (success) {
                    await loadEquipmentTypes();
                    await loadCategoryButtons();
                    await loadEquipmentTypesToDropdown();
                    closeAddTypeModal();
                    // alert('✅ Тип оборудования успешно добавлен!');
                    console.log('✅ Тип оборудования успешно добавлен!');
                }
            } catch (error) {
                // alert('❌ Ошибка: ' + error.message);
                console.error('❌ Ошибка: ', error);
            } finally {
                submitBtn.disabled = false;
                submitBtn.innerHTML = '<i class="fas fa-plus"></i> Добавить тип';
            }
        });
    }
    
    addTypeModal.addEventListener('click', (e) => {
        if (e.target === addTypeModal) {
            closeAddTypeModal();
        }
    });
}

document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && addTypeModal && addTypeModal.classList.contains('active')) {
        closeAddTypeModal();
    }
});

// ==========================================
// МОДАЛЬНОЕ ОКНО ДОБАВЛЕНИЯ ОБОРУДОВАНИЯ
// ==========================================
const modal = document.querySelector('.modal');
const closeModal = document.querySelector('.modal-close');
const uploadArea = document.querySelector('.upload-area');
const fileInput = document.querySelector('#fileInput');
const equipmentName = document.querySelector('#equipmentName');
const equipmentDescription = document.querySelector('#equipmentDescription');
const addBtn = document.querySelector('.add-equipment-btn');

if (addCard && modal) {
    addCard.addEventListener('click', (e) => {
        e.stopPropagation();
        modal.classList.add('active');
        document.body.style.overflow = 'hidden';
    });
}

if (closeModal && modal) {
    closeModal.addEventListener('click', () => {
        modal.classList.remove('active');
        document.body.style.overflow = '';
        resetModal();
    });
}

if (modal) {
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            modal.classList.remove('active');
            document.body.style.overflow = '';
            resetModal();
        }
    });
}

document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && modal && modal.classList.contains('active')) {
        modal.classList.remove('active');
        document.body.style.overflow = '';
        resetModal();
    }
});

if (uploadArea && fileInput) {
    uploadArea.addEventListener('click', () => {
        fileInput.click();
    });
    
    uploadArea.addEventListener('dragover', (e) => {
        e.preventDefault();
        uploadArea.classList.add('dragover');
    });
    
    uploadArea.addEventListener('dragleave', () => {
        uploadArea.classList.remove('dragover');
    });
    
    uploadArea.addEventListener('drop', (e) => {
        e.preventDefault();
        uploadArea.classList.remove('dragover');
        
        const files = e.dataTransfer.files;
        if (files.length > 0) {
            handleFile(files[0]);
        }
    });
    
    fileInput.addEventListener('change', (e) => {
        if (e.target.files.length > 0) {
            handleFile(e.target.files[0]);
        }
    });
}

let selectedFile = null;
let selectedFileName = '';

function handleFile(file) {
    if (!file.type.startsWith('image/')) {
        alert('Пожалуйста, выберите файл изображения (JPG, PNG, GIF и т.д.)');
        return;
    }
    
    selectedFile = file;
    selectedFileName = file.name;
    
    const reader = new FileReader();
    reader.onload = (e) => {
        uploadArea.innerHTML = `
            <div class="file-preview">
                <img src="${e.target.result}" alt="Preview">
                <div class="file-name">${selectedFileName}</div>
                <div class="change-file">Изменить</div>
            </div>
        `;
        
        uploadArea.addEventListener('click', () => {
            document.querySelector('#fileInput').click();
        });
    };
    reader.readAsDataURL(file);
}

if (addBtn) {
    addBtn.addEventListener('click', async () => {
        const name = equipmentName.value.trim();
        const category = document.querySelector('#equipmentCategory').value.trim();
        const description = equipmentDescription.value.trim();
        
        if (!selectedFile) {
            alert('Пожалуйста, загрузите изображение оборудования');
            return;
        }
        if (!name) {
            alert('Пожалуйста, введите название оборудования');
            equipmentName.focus();
            return;
        }
        if (!category) {
            alert('Пожалуйста, выберите категорию оборудования');
            document.querySelector('#equipmentCategory').focus();
            return;
        }
        if (!description) {
            alert('Пожалуйста, введите описание оборудования');
            equipmentDescription.focus();
            return;
        }

        addBtn.disabled = true;
        addBtn.textContent = 'Загрузка...';

        try {
            // 🔧 1. Загружаем изображение НА ЛОКАЛЬНЫЙ СЕРВЕР
            const formData = new FormData();
            formData.append('image', selectedFile);
            
            const uploadResponse = await fetch('http://localhost:8080/api/upload', {
                method: 'POST',
                body: formData
            });
            
            if (!uploadResponse.ok) {
                throw new Error('Не удалось загрузить изображение');
            }
            
            const uploadData = await uploadResponse.json();
            const imageUrl = uploadData.url; // ✅ Получаем URL от сервера (например: /static/equipment/abc123.png)
            
            console.log('✅ Изображение загружено:', imageUrl);
            
            // 🔧 2. Добавляем оборудование с URL
            const success = await addEquipment({
                image: imageUrl,  // URL с локального сервера
                title: name,
                category: category,
                description: description,
                auditory: ''
            });
            
            if (success) {
                await renderEquipment();
                modal.classList.remove('active');
                document.body.style.overflow = '';
                resetModal();
                // alert('✅ Оборудование успешно добавлено!');
                console.log('✅ Оборудование успешно добавлено!');
            }
        } catch (error) {
            console.error('Ошибка:', error);
            // alert('Ошибка: ' + error.message);
        } finally {
            addBtn.disabled = false;
            addBtn.textContent = 'Добавить';
        }
    });
}

function resetModal() {
    selectedFile = null;
    selectedFileName = '';
    
    if (uploadArea) {
        uploadArea.innerHTML = `<i class="fas fa-plus"></i><span>Нажмите для загрузки фото</span><small>или перетащите файл сюда</small>`;
    }
    if (equipmentName) equipmentName.value = '';
    const categoryInput = document.querySelector('#equipmentCategory');
    if (categoryInput) categoryInput.value = '';
    if (equipmentDescription) equipmentDescription.value = '';
    if (fileInput) fileInput.value = '';
}

// ==========================================
// ЗАГРУЗКА ТИПОВ ОБОРУДОВАНИЯ В DROPDOWN
// ==========================================
async function loadEquipmentTypesToDropdown() {
    try {
        const types = await getEquipmentTypesData();
        populateDropdown(types);
        console.log('✅ Загружено типов оборудования в dropdown:', types.length);
    } catch (error) {
        console.error('❌ Ошибка загрузки типов:', error);
        dropdownSelected.textContent = 'Ошибка загрузки';
        dropdownSelected.style.color = '#ff4757';
    }
}

// ==========================================
// ФИЛЬТРАЦИЯ И ПОИСК
// ==========================================
const searchInput = document.querySelector('.search-input');
const searchBtn = document.querySelector('.search-btn');

const performSearch = () => {
    const searchTerm = searchInput.value.toLowerCase().trim();
    if (!searchTerm) {
        document.querySelectorAll('.equip-card:not(.add-card)').forEach(card => {
            card.style.display = 'flex';
        });
        return;
    }

    const equipCards = document.querySelectorAll('.equip-card:not(.add-card)');
    equipCards.forEach(card => {
        // 🔧 Ищем ТОЛЬКО по названию (title), НЕ по категории
        const title = card.querySelector('h3').textContent.toLowerCase();
        
        if (title.includes(searchTerm)) {
            card.style.display = 'flex';
            card.style.animation = 'fadeIn 0.5s ease';
        } else {
            card.style.display = 'none';
        }
    });
};

if (searchBtn) {
    searchBtn.addEventListener('click', performSearch);
}

if (searchInput) {
    searchInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') performSearch();
    });
    
    searchInput.addEventListener('input', (e) => {
        if (e.target.value.trim() === '') {
            document.querySelectorAll('.equip-card:not(.add-card)').forEach(card => {
                card.style.display = 'flex';
            });
        }
    });
}

// ==========================================
// ДОБАВЛЕНИЕ ТИПА ОБОРУДОВАНИЯ
// ==========================================
async function addEquipmentType(typeName) {
    const user = getCurrentUser();
    if (!user) {
        // alert('Пользователь не авторизован');
        console.log('Пользователь не авторизован');
        return false;
    }

    const userRole = (user.Role || '').toLowerCase();
    if (userRole !== 'admin' && userRole !== 'super_admin') {
        // alert('Только администраторы могут добавлять типы оборудования');
        console.log('Только администраторы могут добавлять типы оборудования');
        return false;
    }

    const payload = {
        admin_id: user.Id,
        types: [{
            name: typeName
        }]
    };

    try {
        const response = await fetch(`${API_URL}/Types`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(payload)
        });
        
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Ошибка ${response.status}: ${errorText}`);
        }
        
        const result = await response.json();
        console.log('✅ Тип оборудования добавлен:', result);
        return true;
    } catch (error) {
        console.error('❌ Ошибка добавления типа:', error);
        throw error;
    }
}

// ==========================================
// ИНИЦИАЛИЗАЦИЯ
// ==========================================
document.addEventListener('DOMContentLoaded', async () => {
    // Добавляем модальное окно выбора удаления в DOM
    document.body.insertAdjacentHTML('beforeend', deleteSelectionModalHTML);
    
    // Инициализация обработчиков модального окна удаления
    const deleteModal = document.querySelector('.delete-selection-modal');
    if (deleteModal) {
        const cancelBtn = deleteModal.querySelector('.delete-selection-btn.cancel');
        const okBtn = deleteModal.querySelector('.delete-selection-btn.ok');
        
        if (cancelBtn) {
            cancelBtn.addEventListener('click', closeDeleteSelectionModal);
        }
        
        if (okBtn) {
            okBtn.addEventListener('click', () => {
                closeDeleteSelectionModal();
                enterDeleteMode();
            });
        }
        
        deleteModal.addEventListener('click', (e) => {
            if (e.target === deleteModal) {
                closeDeleteSelectionModal();
            }
        });
        
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape' && deleteModal && deleteModal.classList.contains('active')) {
                closeDeleteSelectionModal();
            }
        });
    }
    
    await renderEquipment();
    await loadEquipmentTypes();
    await loadCategoryButtons();
    await loadEquipmentTypesToDropdown();
    console.log('Project Space initialized');
});

// ==========================================
// АВТОРИЗАЦИЯ
// ==========================================
const authBtn = document.getElementById('authBtn');
const userMenu = document.getElementById('userMenu');
const userAvatar = document.getElementById('userAvatar');
const userDropdown = document.getElementById('userDropdown');
const logoutBtn = document.getElementById('logoutBtn');
const dropdownAvatar = document.getElementById('dropdownAvatar');
const dropdownName = document.getElementById('dropdownName');
const dropdownRole = document.getElementById('dropdownRole');
const dropdownEmail = document.getElementById('dropdownEmail');

function updateAuthButton() {
    const user = getCurrentUser();
    const token = getAuthToken();
    
    if (token && user && user.Email) {
        if (authBtn) authBtn.style.display = 'none';
        if (userMenu) userMenu.style.display = 'block';
        
        let avatarUrl = user.AvatarURL;
        if (!avatarUrl || avatarUrl === 'null' || avatarUrl === '') {
            avatarUrl = 'https://ui-avatars.com/api/?name=' + 
                       encodeURIComponent(user.FirstName || 'U') + 
                       '+' + encodeURIComponent(user.SecondName || 'N') + 
                       '&background=9d00d6&color=fff&size=100';
        }
        
        if (userAvatar) {
            userAvatar.src = avatarUrl;
            userAvatar.onerror = function() {
                this.src = 'https://ui-avatars.com/api/?name=' + 
                          encodeURIComponent(user.FirstName || 'U') + 
                          '&background=9d00d6&color=fff';
            };
        }
        if (dropdownAvatar) {
            dropdownAvatar.src = avatarUrl;
            dropdownAvatar.onerror = function() {
                this.src = 'https://ui-avatars.com/api/?name=' + 
                          encodeURIComponent(user.FirstName || 'U') + 
                          '&background=9d00d6&color=fff';
            };
        }
        
        const fullName = `${user.FirstName || ''} ${user.SecondName || ''}`.trim() || 'Имя Фамилия';
        if (dropdownName) dropdownName.textContent = fullName;
        if (dropdownRole) {
            const roleMap = {
                'super_admin': 'Главный администратор',
                'admin': 'Администратор', 
                'student': 'Студент'
            };
            const roleKey = (user.Role || '').toLowerCase();
            dropdownRole.textContent = roleMap[roleKey] || 'Пользователь';
            // Добавляем класс для стилизации
            dropdownRole.className = 'dropdown-role role-' + roleKey;
        }
        if (dropdownEmail) dropdownEmail.textContent = user.Email;
        
        const addCard = document.querySelector('.add-card');
        if (addCard) {
            const userRole = (user.Role || '').toLowerCase();
            const isAdmin = userRole === 'admin' || userRole === 'super_admin';
            addCard.style.display = isAdmin ? 'flex' : 'none';
            console.log('👤 Роль:', user.Role, '| Карточка:', isAdmin ? 'показана' : 'скрыта');
        }
    } else {
        if (authBtn) {
            authBtn.style.display = 'flex';
            authBtn.onclick = loginWithYandex;
        }
        if (userMenu) userMenu.style.display = 'none';
        if (userDropdown) userDropdown.classList.remove('active');
        
        const addCard = document.querySelector('.add-card');
        if (addCard) {
            addCard.style.display = 'none';
        }
    }
}

if (userAvatar && userAvatar.parentElement) {
    userAvatar.parentElement.addEventListener('click', (e) => {
        e.stopPropagation();
        if (userDropdown) {
            userDropdown.classList.toggle('active');
        }
    });
}

document.addEventListener('click', (e) => {
    if (userDropdown && userMenu && !userMenu.contains(e.target)) {
        userDropdown.classList.remove('active');
    }
});

if (logoutBtn) {
    logoutBtn.addEventListener('click', () => {
        logout();
        if (userDropdown) {
            userDropdown.classList.remove('active');
        }
    });
}

// ========== КАСТОМНЫЙ DROPDOWN ==========
const dropdownTrigger = document.getElementById('dropdownTrigger');
const dropdownMenu = document.getElementById('dropdownMenu');
const dropdownSelected = document.querySelector('.dropdown-selected');
const hiddenInput = document.getElementById('equipmentCategory');

// Открыть/закрыть меню
dropdownTrigger.addEventListener('click', () => {
    dropdownMenu.classList.toggle('open');
    dropdownTrigger.classList.toggle('active');
});

// Закрыть при клике вне
document.addEventListener('click', (e) => {
    if (!e.target.closest('.custom-dropdown')) {
        dropdownMenu.classList.remove('open');
        dropdownTrigger.classList.remove('active');
    }
});

// Выбрать опцию
function selectOption(option) {
    const value = option.dataset.value;
    const text = option.textContent.trim();
    
    // Обновляем видимый текст
    dropdownSelected.textContent = text;
    
    // Обновляем скрытый инпут
    hiddenInput.value = value;
    
    // Подсветка выбранной опции
    document.querySelectorAll('.dropdown-option').forEach(opt => {
        opt.classList.remove('selected');
    });
    option.classList.add('selected');
    
    // Закрываем меню
    dropdownMenu.classList.remove('open');
    dropdownTrigger.classList.remove('active');
    
    console.log('Выбрана категория:', text, '(value:', value, ')');
}

// Клик по опциям
dropdownMenu.addEventListener('click', (e) => {
    const option = e.target.closest('.dropdown-option');
    if (option) {
        selectOption(option);
    }
});

// Функция для динамической загрузки категорий из БД
function populateDropdown(types) {
    // Очищаем кроме плейсхолдера
    dropdownMenu.innerHTML = `
        <div class="dropdown-option" data-value="">Выберите категорию</div>
    `;
    
    types.forEach(type => {
        const option = document.createElement('div');
        option.className = 'dropdown-option';
        option.dataset.value = type.name || type.Name;
        option.textContent = type.name || type.Name;
        dropdownMenu.appendChild(option);
    });
}

document.addEventListener('DOMContentLoaded', () => {
    updateAuthButton();
    
    if (window.location.search.includes('token')) {
        handleOAuthCallback().then(() => {
            updateAuthButton();
            if (window.history.replaceState) {
                window.history.replaceState({}, document.title, window.location.pathname);
            }
        });
    }
});