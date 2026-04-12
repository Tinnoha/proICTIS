import asyncio
import json
import logging
import os
from datetime import datetime, timedelta
from typing import Dict, List, Optional
import aiohttp
from aiogram import Bot, Dispatcher, types
from aiogram.filters import Command
from aiogram.types import Message, CallbackQuery
from aiogram.fsm.context import FSMContext
from aiogram.fsm.state import State, StatesGroup
from aiogram.fsm.storage.memory import MemoryStorage
from aiogram.utils.keyboard import InlineKeyboardBuilder
from dotenv import load_dotenv

load_dotenv()

# Настройка логирования
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

#Конфигурация
BOT_TOKEN = os.getenv("BOT_TOKEN")
API_BASE_URL = os.getenv("API_BASE_URL", "http://localhost:8080")

admin_ids_str = os.getenv("ADMIN_IDS", "")
ADMIN_IDS = []
if admin_ids_str:
    try:
        ADMIN_IDS = [int(x.strip()) for x in admin_ids_str.split(",") if x.strip()]
    except ValueError:
        logger.warning(f"Invalid ADMIN_IDS format: {admin_ids_str}")

# Проверка наличия токена
if not BOT_TOKEN:
    raise ValueError("BOT_TOKEN not found in .env file!")

# Инициализация бота и диспетчера
bot = Bot(token=BOT_TOKEN)
dp = Dispatcher(storage=MemoryStorage())

# Работа с API 
class BackendAPI:
    def __init__(self, base_url: str):
        self.base_url = base_url
        self.session: Optional[aiohttp.ClientSession] = None
    
    async def __aenter__(self):
        self.session = aiohttp.ClientSession()
        return self
    
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        if self.session:
            await self.session.close()
    
    async def _request(self, method: str, endpoint: str, data: dict = None, params: dict = None):
        """Универсальный метод для запросов к API"""
        url = f"{self.base_url}{endpoint}"
        headers = {"Content-Type": "application/json"}
        
        try:
            async with self.session.request(method, url, json=data, params=params, headers=headers) as response:
                if response.status == 204:
                    return None, response.status
                
                response_data = await response.json() if response.status != 204 else None
                return response_data, response.status
        except aiohttp.ClientError as e:
            logger.error(f"API request error: {e}")
            return {"Error": str(e)}, 500
    
    # Методы для работы с оборудованием
    async def get_equipment(self) -> tuple:
        """GET /Equipment - получить все оборудование"""
        return await self._request("GET", "/Equipment")
    
    async def get_equipment_by_type(self, equipment_type: str) -> tuple:
        """GET /Equipment/type/{type} - получить оборудование по типу"""
        return await self._request("GET", f"/Equipment/type/{equipment_type}")
    
    async def get_equipment_by_id(self, equipment_id: str) -> tuple:
        """GET /Equipment/id/{id} - получить оборудование по ID"""
        return await self._request("GET", f"/Equipment/id/{equipment_id}")
    
    # Методы для работы с типами оборудования
    async def get_types(self) -> tuple:
        """GET /Types - получить все типы оборудования"""
        return await self._request("GET", "/Types")
    
    # Методы для работы с бронированиями
    async def create_booking(self, user_id: str, equipment_id: str, start: str, end: str) -> tuple:
        """POST /Booking - создать бронирование"""
        data = {
            "user_id": user_id,
            "enviromt_id": equipment_id,
            "start": start,
            "end": end
        }
        return await self._request("POST", "/Booking", data=data)
    
    async def get_user_bookings(self, user_id: str) -> tuple:
        """GET /Booking/user/{id} - получить бронирования пользователя"""
        return await self._request("GET", f"/Booking/user/{user_id}")
    
    async def get_equipment_bookings(self, equipment_id: str) -> tuple:
        """GET /Booking/equipment/{id} - получить бронирования оборудования"""
        return await self._request("GET", f"/Booking/equipment/{equipment_id}")
    
    async def get_all_bookings(self) -> tuple:
        """GET /Booking - получить все бронирования (для админов)"""
        return await self._request("GET", "/Booking")
    
    async def update_booking_status(self, booking_id: str, admin_id: str, status: str) -> tuple:
        """PUT /Booking/{id} - обновить статус бронирования"""
        data = {
            "admin_id": admin_id,
            "status": status
        }
        return await self._request("PUT", f"/Booking/{booking_id}", data=data)
    
    async def return_equipment(self, booking_id: str, admin_id: str) -> tuple:
        """PUT /Booking/return/{id} - вернуть оборудование"""
        data = {"admin_id": admin_id}
        return await self._request("PUT", f"/Booking/return/{booking_id}", data=data)
    
    # Методы для работы с пользователями
    async def get_user_by_email(self, admin_id: str, email: str) -> tuple:
        """GET /User/email - получить пользователя по email"""
        data = {"AdminId": admin_id, "Email": email}
        return await self._request("GET", "/User/email", data=data)
    
    async def get_user_by_id(self, user_id: str) -> tuple:
        """GET /User/{id} - получить пользователя по ID"""
        return await self._request("GET", f"/User/{user_id}")
    
    async def get_all_users(self) -> tuple:
        """GET /User - получить всех пользователей"""
        return await self._request("GET", "/User")

# Инициализация API клиента
api = BackendAPI(API_BASE_URL)

# Состояния для FSM
class BookingStates(StatesGroup):
    choosing_equipment = State()
    choosing_start_date = State()
    choosing_end_date = State()
    confirming = State()

class AdminStates(StatesGroup):
    choosing_booking_for_action = State()
    choosing_status = State()
    input_email = State()

# Клавиатура
def get_main_keyboard(user_id: int) -> InlineKeyboardBuilder:
    """Создает основную клавиатуру для всех пользователей"""
    builder = InlineKeyboardBuilder()
    builder.button(text="🔧 Забронировать оборудование", callback_data="start_booking")
    builder.button(text="📋 Мои бронирования", callback_data="my_bookings")
    
    if user_id in ADMIN_IDS:
        builder.button(text="👑 Админ-панель", callback_data="admin_panel")
    
    builder.adjust(1)
    return builder

def get_equipment_keyboard(equipment_list: list) -> InlineKeyboardBuilder:
    """Создает клавиатуру со списком оборудования"""
    builder = InlineKeyboardBuilder()
    for eq in equipment_list:
        name = eq.get('Name', 'Без названия')
        eq_id = eq.get('Id', '')
        builder.button(text=name, callback_data=f"equipment:{eq_id}")
    builder.button(text="🔙 В главное меню", callback_data="back_to_main")
    builder.adjust(1)
    return builder

def get_date_keyboard(prefix: str = "start", days: int = 7) -> InlineKeyboardBuilder:
    """Создает клавиатуру с датами"""
    builder = InlineKeyboardBuilder()
    current_date = datetime.now()
    
    for i in range(days):
        date = current_date + timedelta(days=i)
        date_str = date.strftime("%Y-%m-%dT%H:%M:%S")
        display_date = date.strftime("%d.%m.%Y")
        
        if i == 0:
            display_date = f"Сегодня ({display_date})"
        elif i == 1:
            display_date = f"Завтра ({display_date})"
            
        builder.button(
            text=display_date,
            callback_data=f"{prefix}_date:{date_str}"
        )
    
    builder.button(text="❌ Отмена", callback_data="cancel")
    builder.adjust(1)
    return builder

def get_admin_keyboard() -> InlineKeyboardBuilder:
    """Создает клавиатуру для админ-панели"""
    builder = InlineKeyboardBuilder()
    builder.button(text="📊 Все бронирования", callback_data="admin_all_bookings")
    builder.button(text="✅ Подтвердить бронь", callback_data="admin_approve_booking")
    builder.button(text="❌ Отменить бронь", callback_data="admin_cancel_booking")
    builder.button(text="🔄 Вернуть оборудование", callback_data="admin_return_equipment")
    builder.button(text="👥 Найти пользователя", callback_data="admin_find_user")
    builder.button(text="🔙 В главное меню", callback_data="back_to_main")
    builder.adjust(1)
    return builder

def get_status_keyboard(booking_id: str) -> InlineKeyboardBuilder:
    """Клавиатура для выбора статуса бронирования"""
    builder = InlineKeyboardBuilder()
    builder.button(text="✅ Active", callback_data=f"status:{booking_id}:Active")
    builder.button(text="❌ Cancel", callback_data=f"status:{booking_id}:Cancel")
    builder.button(text="🔄 Returned", callback_data=f"status:{booking_id}:Returned")
    builder.button(text="⏳ Waiting answer", callback_data=f"status:{booking_id}:Waiting answer")
    builder.button(text="🔙 Назад", callback_data="admin_back")
    builder.adjust(1)
    return builder

def format_bookings_for_display(bookings: list) -> str:
    """Форматирует список бронирований для отображения"""
    if not bookings:
        return "У вас пока нет бронирований."
    
    result = "📋 Ваши бронирования:\n\n"
    for booking in bookings:
        booking_id = booking.get('ID', '')
        equipment_id = booking.get('EquipmentId', '')
        start = booking.get('BookStart', '')
        end = booking.get('BookEnd', '')
        status = booking.get('Status', '')
        
        equipment_name = equipment_id
        
        result += f"🔧 ID брони: {booking_id}\n"
        result += f"📅 {start} - {end}\n"
        result += f"📊 Статус: {status}\n\n"
    
    return result

# Обработчики команд
@dp.message(Command("start"))
async def cmd_start(message: Message):
    user = message.from_user
    welcome_text = (
        f"👋 Добро пожаловать в бот бронирования оборудования, {user.first_name}!\n\n"
        "Я помогу вам забронировать оборудование через интеграцию с бекендом.\n\n"
        "Используйте кнопки ниже для навигации:"
    )
    
    keyboard = get_main_keyboard(user.id)
    await message.answer(welcome_text, reply_markup=keyboard.as_markup())

@dp.message(Command("help"))
async def cmd_help(message: Message):
    help_text = (
        "📚 Справка по использованию бота:\n\n"
        "🔹 Используйте кнопки меню для навигации\n"
        "🔹 'Забронировать оборудование' - начать бронирование\n"
        "🔹 'Мои бронирования' - просмотр своих броней\n\n"
        "Статусы бронирований:\n"
        "• Waiting answer - ожидает подтверждения\n"
        "• Active - активна\n"
        "• Cancel - отменена\n"
        "• Returned - возвращена\n\n"
        "Все данные синхронизируются с бекендом через API."
    )
    await message.answer(help_text)

@dp.callback_query(lambda c: c.data == "back_to_main")
async def back_to_main(callback: CallbackQuery, state: FSMContext):
    await state.clear()
    user = callback.from_user
    keyboard = get_main_keyboard(user.id)
    await callback.message.edit_text("🏠 Главное меню:", reply_markup=keyboard.as_markup())
    await callback.answer()

@dp.callback_query(lambda c: c.data == "start_booking")
async def start_booking(callback: CallbackQuery, state: FSMContext):
    """Начинает процесс бронирования - получает список оборудования"""
    await callback.message.edit_text("🔄 Загрузка списка оборудования...")
    
    equipment_data, status = await api.get_equipment()
    
    if status != 200 or not equipment_data:
        await callback.message.edit_text(
            "❌ Не удалось загрузить список оборудования. Попробуйте позже.\n"
            f"Ошибка: {equipment_data.get('Error', 'Неизвестная ошибка') if equipment_data else 'Нет данных'}"
        )
        await callback.answer()
        return
    
    # Получаем только активное оборудование
    active_equipment = [eq for eq in equipment_data if eq.get('IsActive', True)]
    
    if not active_equipment:
        await callback.message.edit_text("📭 Нет доступного оборудования для бронирования.")
        await callback.answer()
        return
    
    keyboard = get_equipment_keyboard(active_equipment)
    await callback.message.edit_text(
        "Выберите оборудование для бронирования:",
        reply_markup=keyboard.as_markup()
    )
    await state.set_state(BookingStates.choosing_equipment)
    await callback.answer()

@dp.callback_query(lambda c: c.data.startswith("equipment:"))
async def process_equipment_choice(callback: CallbackQuery, state: FSMContext):
    equipment_id = callback.data.split(":")[1]
    await state.update_data(equipment_id=equipment_id)
    
    keyboard = get_date_keyboard("start")
    await callback.message.edit_text(
        "Выберите дату начала бронирования:",
        reply_markup=keyboard.as_markup()
    )
    await state.set_state(BookingStates.choosing_start_date)
    await callback.answer()

@dp.callback_query(lambda c: c.data.startswith("start_date:"))
async def process_start_date(callback: CallbackQuery, state: FSMContext):
    start_date = callback.data.split(":")[1]
    await state.update_data(start_date=start_date)
    
    keyboard = get_date_keyboard("end")
    await callback.message.edit_text(
        f"Дата начала: {start_date}\n\n"
        "Теперь выберите дату окончания бронирования:",
        reply_markup=keyboard.as_markup()
    )
    await state.set_state(BookingStates.choosing_end_date)
    await callback.answer()

@dp.callback_query(lambda c: c.data.startswith("end_date:"))
async def process_end_date(callback: CallbackQuery, state: FSMContext):
    end_date = callback.data.split(":")[1]
    data = await state.get_data()
    start_date = data.get('start_date')
    
    if end_date < start_date:
        await callback.message.edit_text(
            "❌ Дата окончания не может быть раньше даты начала!\n"
            "Попробуйте снова.",
            reply_markup=get_date_keyboard("end").as_markup()
        )
        await callback.answer()
        return
    
    await state.update_data(end_date=end_date)
    
    equipment_data, status = await api.get_equipment_by_id(data['equipment_id'])
    equipment_name = equipment_data.get('Name', 'Неизвестно') if equipment_data else data['equipment_id']
    
    keyboard = InlineKeyboardBuilder()
    keyboard.button(text="✅ Подтвердить", callback_data="confirm_booking")
    keyboard.button(text="❌ Отмена", callback_data="cancel")
    keyboard.adjust(1)
    
    await callback.message.edit_text(
        f"📝 Подтвердите бронирование:\n\n"
        f"🔧 Оборудование: {equipment_name}\n"
        f"📅 Период: {start_date} - {end_date}\n\n"
        f"Всё верно?",
        reply_markup=keyboard.as_markup()
    )
    await state.set_state(BookingStates.confirming)
    await callback.answer()

@dp.callback_query(lambda c: c.data == "confirm_booking")
async def confirm_booking(callback: CallbackQuery, state: FSMContext):
    data = await state.get_data()
    user = callback.from_user
    
    user_uuid = f"user_{user.id}"
    
    await callback.message.edit_text("🔄 Отправка запроса на бронирование...")
    
    # Создание брони в БД
    booking_data, status = await api.create_booking(
        user_id=user_uuid,
        equipment_id=data['equipment_id'],
        start=data['start_date'],
        end=data['end_date']
    )
    
    if status == 201:
        await state.clear()
        keyboard = get_main_keyboard(user.id)
        await callback.message.edit_text(
            "✅ Бронирование успешно создано!\n\n"
            f"🔧 ID брони: {booking_data.get('ID', 'N/A')}\n"
            f"📅 Период: {data['start_date']} - {data['end_date']}\n"
            f"📊 Статус: {booking_data.get('Status', 'Waiting answer')}\n\n"
            "Ожидайте подтверждения от администратора.",
            reply_markup=keyboard.as_markup()
        )
    else:
        error_msg = booking_data.get('Error', 'Неизвестная ошибка') if booking_data else 'Ошибка соединения'
        await callback.message.edit_text(
            f"❌ Ошибка при создании бронирования:\n{error_msg}\n\n"
            "Возможно, оборудование уже забронировано на выбранные даты."
        )
    
    await callback.answer()

@dp.callback_query(lambda c: c.data == "my_bookings")
async def show_my_bookings(callback: CallbackQuery):
    user = callback.from_user
    user_uuid = f"user_{user.id}"
    
    await callback.message.edit_text("🔄 Загрузка ваших бронирований...")
    
    bookings_data, status = await api.get_user_bookings(user_uuid)
    
    keyboard = InlineKeyboardBuilder()
    keyboard.button(text="🔙 В главное меню", callback_data="back_to_main")
    
    if status == 200 and bookings_data:
        response = format_bookings_for_display(bookings_data)
    elif status == 404:
        response = "У вас пока нет бронирований."
    else:
        response = f"❌ Ошибка загрузки бронирований. Статус: {status}"
    
    await callback.message.edit_text(response, reply_markup=keyboard.as_markup())
    await callback.answer()

# Админ-панель
@dp.callback_query(lambda c: c.data == "admin_panel")
async def admin_panel(callback: CallbackQuery):
    if callback.from_user.id not in ADMIN_IDS:
        await callback.answer("❌ У вас нет доступа к админ-панели", show_alert=True)
        return
    
    keyboard = get_admin_keyboard()
    await callback.message.edit_text(
        "👑 Админ-панель\n\n"
        "Выберите действие:",
        reply_markup=keyboard.as_markup()
    )
    await callback.answer()

@dp.callback_query(lambda c: c.data == "admin_all_bookings")
async def admin_all_bookings(callback: CallbackQuery):
    if callback.from_user.id not in ADMIN_IDS:
        await callback.answer("❌ Доступ запрещен", show_alert=True)
        return
    
    await callback.message.edit_text("🔄 Загрузка всех бронирований...")
    
    bookings_data, status = await api.get_all_bookings()
    
    keyboard = InlineKeyboardBuilder()
    keyboard.button(text="🔙 В админ-панель", callback_data="admin_panel")
    keyboard.button(text="🏠 В главное меню", callback_data="back_to_main")
    
    if status == 200 and bookings_data:
        response = "📊 ВСЕ БРОНИРОВАНИЯ:\n\n"
        for booking in bookings_data:
            response += (
                f"🔧 ID: {booking.get('ID', 'N/A')}\n"
                f"👤 User: {booking.get('UserId', 'N/A')}\n"
                f"📅 {booking.get('BookStart', '')} - {booking.get('BookEnd', '')}\n"
                f"📊 Статус: {booking.get('Status', '')}\n"
                f"{'-' * 30}\n"
            )
    elif status == 404:
        response = "📭 Нет активных бронирований."
    else:
        response = f"❌ Ошибка загрузки. Статус: {status}"
    
    await callback.message.edit_text(response, reply_markup=keyboard.as_markup())
    await callback.answer()

@dp.callback_query(lambda c: c.data == "admin_approve_booking")
async def admin_approve_booking_start(callback: CallbackQuery, state: FSMContext):
    if callback.from_user.id not in ADMIN_IDS:
        await callback.answer("❌ Доступ запрещен", show_alert=True)
        return
    
    await callback.message.edit_text("🔄 Загрузка ожидающих бронирований...")
    
    bookings_data, status = await api.get_all_bookings()
    
    if status != 200 or not bookings_data:
        await callback.message.edit_text(
            "❌ Нет бронирований для обработки.",
            reply_markup=get_admin_keyboard().as_markup()
        )
        await callback.answer()
        return
    
    waiting_bookings = [b for b in bookings_data if b.get('Status') == 'Waiting answer']
    
    if not waiting_bookings:
        await callback.message.edit_text(
            "📭 Нет ожидающих подтверждения бронирований.",
            reply_markup=get_admin_keyboard().as_markup()
        )
        await callback.answer()
        return
    
    builder = InlineKeyboardBuilder()
    for booking in waiting_bookings:
        booking_id = booking.get('ID', '')
        start = booking.get('BookStart', '')[:10]
        builder.button(
            text=f"Бронь {booking_id[:8]}... - {start}",
            callback_data=f"select_booking:{booking_id}"
        )
    builder.button(text="🔙 Назад", callback_data="admin_back")
    builder.adjust(1)
    
    await callback.message.edit_text(
        "Выберите бронирование для подтверждения:",
        reply_markup=builder.as_markup()
    )
    await state.set_state(AdminStates.choosing_booking_for_action)
    await callback.answer()

@dp.callback_query(lambda c: c.data.startswith("select_booking:"))
async def admin_select_booking(callback: CallbackQuery, state: FSMContext):
    booking_id = callback.data.split(":")[1]
    await state.update_data(booking_id=booking_id)
    
    keyboard = get_status_keyboard(booking_id)
    await callback.message.edit_text(
        f"Выберите новый статус для бронирования {booking_id}:",
        reply_markup=keyboard.as_markup()
    )
    await callback.answer()

@dp.callback_query(lambda c: c.data.startswith("status:"))
async def admin_update_status(callback: CallbackQuery, state: FSMContext):
    _, booking_id, status = callback.data.split(":")
    admin_id = f"admin_{callback.from_user.id}"
    
    await callback.message.edit_text(f"🔄 Обновление статуса бронирования...")
    
    result_data, result_status = await api.update_booking_status(booking_id, admin_id, status)
    
    if result_status in [200, 204]:
        keyboard = get_admin_keyboard()
        await callback.message.edit_text(
            f"✅ Статус бронирования {booking_id} успешно изменен на '{status}'.",
            reply_markup=keyboard.as_markup()
        )
    else:
        error_msg = result_data.get('Error', 'Неизвестная ошибка') if result_data else 'Ошибка соединения'
        await callback.message.edit_text(
            f"❌ Ошибка при обновлении статуса:\n{error_msg}",
            reply_markup=get_admin_keyboard().as_markup()
        )
    
    await state.clear()
    await callback.answer()

@dp.callback_query(lambda c: c.data == "admin_return_equipment")
async def admin_return_equipment_start(callback: CallbackQuery, state: FSMContext):
    if callback.from_user.id not in ADMIN_IDS:
        await callback.answer("❌ Доступ запрещен", show_alert=True)
        return
    
    await callback.message.edit_text("🔄 Загрузка активных бронирований...")
    
    bookings_data, status = await api.get_all_bookings()
    
    if status != 200 or not bookings_data:
        await callback.message.edit_text(
            "❌ Нет активных бронирований.",
            reply_markup=get_admin_keyboard().as_markup()
        )
        await callback.answer()
        return
    
    active_bookings = [b for b in bookings_data if b.get('Status') == 'Active']
    
    if not active_bookings:
        await callback.message.edit_text(
            "📭 Нет активных бронирований для возврата.",
            reply_markup=get_admin_keyboard().as_markup()
        )
        await callback.answer()
        return
    
    builder = InlineKeyboardBuilder()
    for booking in active_bookings:
        booking_id = booking.get('ID', '')
        start = booking.get('BookStart', '')[:10]
        builder.button(
            text=f"Бронь {booking_id[:8]}... - {start}",
            callback_data=f"return_booking:{booking_id}"
        )
    builder.button(text="🔙 Назад", callback_data="admin_back")
    builder.adjust(1)
    
    await callback.message.edit_text(
        "Выберите бронирование для возврата оборудования:",
        reply_markup=builder.as_markup()
    )
    await callback.answer()

@dp.callback_query(lambda c: c.data.startswith("return_booking:"))
async def admin_return_equipment_confirm(callback: CallbackQuery, state: FSMContext):
    booking_id = callback.data.split(":")[1]
    admin_id = f"admin_{callback.from_user.id}"
    
    builder = InlineKeyboardBuilder()
    builder.button(text="✅ Да, вернуть", callback_data=f"confirm_return:{booking_id}")
    builder.button(text="❌ Нет, отмена", callback_data="admin_back")
    builder.adjust(1)
    
    await callback.message.edit_text(
        f"⚠️ Вы уверены, что хотите вернуть оборудование по бронированию {booking_id}?",
        reply_markup=builder.as_markup()
    )
    await callback.answer()

@dp.callback_query(lambda c: c.data.startswith("confirm_return:"))
async def admin_return_equipment_final(callback: CallbackQuery, state: FSMContext):
    booking_id = callback.data.split(":")[1]
    admin_id = f"admin_{callback.from_user.id}"
    
    await callback.message.edit_text("🔄 Возврат оборудования...")
    
    result_data, result_status = await api.return_equipment(booking_id, admin_id)
    
    if result_status == 204:
        keyboard = get_admin_keyboard()
        await callback.message.edit_text(
            f"✅ Оборудование по бронированию {booking_id} успешно возвращено.",
            reply_markup=keyboard.as_markup()
        )
    else:
        error_msg = result_data.get('Error', 'Неизвестная ошибка') if result_data else 'Ошибка соединения'
        await callback.message.edit_text(
            f"❌ Ошибка при возврате оборудования:\n{error_msg}",
            reply_markup=get_admin_keyboard().as_markup()
        )
    
    await state.clear()
    await callback.answer()

@dp.callback_query(lambda c: c.data == "admin_back")
async def admin_back(callback: CallbackQuery, state: FSMContext):
    await state.clear()
    keyboard = get_admin_keyboard()
    await callback.message.edit_text(
        "👑 Админ-панель\n\nВыберите действие:",
        reply_markup=keyboard.as_markup()
    )
    await callback.answer()

@dp.callback_query(lambda c: c.data == "cancel")
async def cancel_operation(callback: CallbackQuery, state: FSMContext):
    await state.clear()
    user = callback.from_user
    keyboard = get_main_keyboard(user.id)
    await callback.message.edit_text(
        "❌ Операция отменена.\n\n🏠 Главное меню:",
        reply_markup=keyboard.as_markup()
    )
    await callback.answer()

# Запуск бота
async def main():
    print("🚀 Запуск Telegram бота...")
    print(f"🤖 API бекенда: {API_BASE_URL}")
    print(f"👑 Администраторы: {ADMIN_IDS}")
    
    # Инициализация сессии API
    async with api:
        await dp.start_polling(bot)

if __name__ == "__main__":
    asyncio.run(main())