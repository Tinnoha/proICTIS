// Данные оборудования
const equipmentData = {
  1: {
    img: "picture/kk.jpg",
    title: "Станок ЧПУ X-200",
    desc: "Высокоточный станок для обработки металла. Мощность 5 кВт, точность ±0.01 мм."
  },
  2: {
    img: "picture/kk.jpg",
    title: "Лазерный гравер ProLaser 3000",
    desc: "Профессиональный гравер для дерева, пластика и металла. Рабочая зона 600×400 мм."
  },
  3: {
    img: "picture/kk.jpg",
    title: "Пресс гидравлический HP-50",
    desc: "Гидравлический пресс с усилием 50 тонн. Подходит для штамповки и формовки."
  },
  4: {
    img: "picture/kk.jpg",
    title: "3D-принтер UltiMaker S5",
    desc: "Профессиональный 3D-принтер с двумя экструдерами. Область печати 330×240×300 мм."
  },
  5: {
    img: "picture/kk.jpg",
    title: "Фрезерный станок WoodMaster",
    desc: "Станок для обработки дерева и композитных материалов. Мощность двигателя 2.2 кВт."
  },

  6: {
    img: "picture/kk.jpg",
    title: "Плазменный резак PlasmaCut 100",
    desc: "Оборудование для плазменной резки металлов толщиной до 50 мм."
  },
  7: {
    img: "picture/kk.jpg",
    title: "Плазменный резак PlasmaCut 100",
    desc: "Оборудование для плазменной резки металлов толщиной до 50 мм."
  },
  8: {
    img: "picture/kk.jpg",
    title: "Плазменный резак PlasmaCut 100",
    desc: "Оборудование для плазменной резки металлов толщиной до 50 мм."
  },
  9: {
    img: "picture/kk.jpg",
    title: "Плазменный резак PlasmaCut 100",
    desc: "Оборудование для плазменной резки металлов толщиной до 50 мм."
  },
  10: {
    img: "picture/kk.jpg",
    title: "Плазменный резак PlasmaCut 100",
    desc: "Оборудование для плазменной резки металлов толщиной до 50 мм."
  },
  11: {
    img: "picture/kk.jpg",
    title: "Плазменный резак PlasmaCut 100",
    desc: "Оборудование для плазменной резки металлов толщиной до 50 мм."
  },
  12: {
    img: "picture/kk.jpg",
    title: "Плазменный резак PlasmaCut 100",
    desc: "Оборудование для плазменной резки металлов толщиной до 50 мм."
  },
  13: {
    img: "picture/kk.jpg",
    title: "Плазменный резак PlasmaCut 100",
    desc: "Оборудование для плазменной резки металлов толщиной до 50 мм."
  },
  14: {
    img: "picture/kk.jpg",
    title: "Плазменный резак PlasmaCut 100",
    desc: "Оборудование для плазменной резки металлов толщиной до 50 мм."
  },
  15: {
    img: "picture/kk.jpg",
    title: "Плазменный резак PlasmaCut 100",
    desc: "Оборудование для плазменной резки металлов толщиной до 50 мм."
  },
  16: {
    img: "picture/kk.jpg",
    title: "Плазменный резак PlasmaCut 100",
    desc: "Оборудование для плазменной резки металлов толщиной до 50 мм."
  },
  17: {
    img: "picture/kk.jpg",
    title: "Плазменный резак PlasmaCut 100",
    desc: "Оборудование для плазменной резки металлов толщиной до 50 мм."
  },
  18: {
    img: "picture/kk.jpg",
    title: "Плазменный резак PlasmaCut 100",
    desc: "Оборудование для плазменной резки металлов толщиной до 50 мм."
  },
  19: {
    img: "picture/kk.jpg",
    title: "Плазменный резак PlasmaCut 100",
    desc: "Оборудование для плазменной резки металлов толщиной до 50 мм."
  },
  20: {
    img: "picture/kk.jpg",
    title: "Плазменный резак PlasmaCut 100",
    desc: "Оборудование для плазменной резки металлов толщиной до 50 мм."
  }
};

// Элементы DOM
const cards = document.querySelectorAll('.card');
const modal = document.getElementById('modal');
const modalImg = document.getElementById('modal-img');
const modalTitle = document.getElementById('modal-title');
const modalDesc = document.getElementById('modal-desc');
const closeBtn = document.querySelector('.close');
const searchInput = document.querySelector('.search-input');
const searchBtn = document.querySelector('.search-btn');

// Открытие модального окна
cards.forEach(card => {
  card.addEventListener('click', () => {
    const id = card.getAttribute('data-id');
    const data = equipmentData[id];
    
    if (data) {
      modalImg.src = data.img;
      modalImg.alt = data.title;
      modalTitle.textContent = data.title;
      modalDesc.textContent = data.desc;
      
      modal.style.display = 'block';
      document.body.style.overflow = 'hidden';
    }
  });
});

// Закрытие модального окна
function closeModal() {
  modal.style.display = 'none';
  document.body.style.overflow = 'auto';
}

closeBtn.addEventListener('click', closeModal);

window.addEventListener('click', (e) => {
  if (e.target === modal) {
    closeModal();
  }
});

// Закрытие по Escape
document.addEventListener('keydown', (e) => {
  if (e.key === 'Escape' && modal.style.display === 'block') {
    closeModal();
  }
});

// Поиск (базовая реализация)
searchBtn.addEventListener('click', performSearch);
searchInput.addEventListener('keypress', (e) => {
  if (e.key === 'Enter') {
    performSearch();
  }
});

function performSearch() {
  const searchTerm = searchInput.value.toLowerCase().trim();
  
  if (searchTerm) {
    // Здесь можно реализовать фильтрацию карточек
    // Пока просто показываем сообщение
    // alert(`Поиск: "${searchTerm}"\n\n.`);
    
    // Пример фильтрации (раскомментировать когда нужно):
    /*
    cards.forEach(card => {
      const id = card.getAttribute('data-id');
      const data = equipmentData[id];
      const cardText = (data.title + ' ' + data.desc).toLowerCase();
      
      if (cardText.includes(searchTerm)) {
        card.style.display = 'block';
      } else {
        card.style.display = 'none';
      }
    });
    */
  }
}

// Адаптация высоты карточек
function adjustCardHeights() {
  if (window.innerWidth >= 992) {
    const cards = document.querySelectorAll('.card');
    let maxHeight = 0;
    
    // Сбрасываем высоту
    cards.forEach(card => {
      card.style.height = 'auto';
    });
    
    // Находим максимальную высоту в ряду
    cards.forEach(card => {
      const cardHeight = card.offsetHeight;
      if (cardHeight > maxHeight) {
        maxHeight = cardHeight;
      }
    });
    
    // Устанавливаем одинаковую высоту
    cards.forEach(card => {
      card.style.height = maxHeight + 'px';
    });
  }
}

// Вызываем при загрузке и изменении размера окна
window.addEventListener('load', adjustCardHeights);
window.addEventListener('resize', adjustCardHeights);