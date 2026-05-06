import os
import asyncio
import aiohttp
import json
import re
import time
from dotenv import load_dotenv
from datetime import datetime, timedelta
from typing import Dict, Optional, List
import vk_api
from vk_api.bot_longpoll import VkBotLongPoll, VkBotEventType
from vk_api.keyboard import VkKeyboard, VkKeyboardColor
import logging

load_dotenv()

VK_TOKEN = os.getenv('VK_TOKEN')
GROUP_ID = int(os.getenv('GROUP_ID')) if os.getenv('GROUP_ID') else None
API_BASE_URL = os.getenv('API_BASE_URL', 'http://localhost:8080')
API_TIMEOUT = int(os.getenv('API_TIMEOUT', 30))

if not VK_TOKEN:
    raise ValueError("VK_TOKEN не найден в .env файле")
if not GROUP_ID:
    raise ValueError("GROUP_ID не найден в .env файле")

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


class EquipmentAPIClient:
    
    def __init__(self, base_url: str, timeout: int = 30):
        self.base_url = base_url
        self.timeout = timeout
        self.users_cache = {}
    
    async def _request(self, method: str, endpoint: str, data: Dict = None, params: Dict = None):
        url = f"{self.base_url}{endpoint}"
        async with aiohttp.ClientSession() as session:
            try:
                async with session.request(
                    method=method,
                    url=url,
                    json=data,
                    params=params,
                    timeout=self.timeout,
                    headers={"Content-Type": "application/json"}
                ) as response:
                    if response.status == 204:
                        return {"success": True}
                    if response.status in [200, 201]:
                        text = await response.text()
                        try:
                            return json.loads(text)
                        except:
                            return {"success": True, "data": text}
                    return None
            except Exception as e:
                logger.error(f"Ошибка API запроса: {e}")
                return None
    
    async def get_user_by_vk_id(self, vk_id: int, force_refresh: bool = False) -> Optional[Dict]:
        """Получить пользователя по VK ID"""
        vk_id_int = int(vk_id)
        
        if not force_refresh and vk_id_int in self.users_cache:
            return self.users_cache[vk_id_int]
        
        result = await self._request("GET", "/User/by-vk", params={"vkid": vk_id_int})
        
        if result and isinstance(result, dict) and 'id' in result:
            self.users_cache[vk_id_int] = result
            return result
        
        return None
    
    async def get_user_by_id(self, user_id: str) -> Optional[Dict]:
        """Получить пользователя по ID"""
        result = await self._request("GET", f"/User/{user_id}")
        if result and isinstance(result, dict) and 'id' in result:
            return result
        return None
    
    async def get_all_equipment(self) -> List[Dict]:
        """Получить всё оборудование"""
        result = await self._request("GET", "/Equipment")
        if isinstance(result, list):
            return result
        return []
    
    async def get_equipment_by_id(self, equipment_id: str) -> Optional[Dict]:
        """Получить оборудование по ID"""
        result = await self._request("GET", f"/Equipment/id/{equipment_id}")
        if result and isinstance(result, dict):
            return result
        return None
    
    async def create_booking(self, user_id: str, equipment_id: str, start: str, end: str) -> Optional[Dict]:
        """Создать бронирование"""
        data = {
            "user_id": user_id,
            "enviromt_id": equipment_id,
            "start": start,
            "end": end
        }
        result = await self._request("POST", "/Booking", data=data)
        return result if isinstance(result, dict) else None
    
    async def get_user_bookings(self, user_id: str) -> List[Dict]:
        """Получить бронирования пользователя"""
        result = await self._request("GET", f"/Booking/user/{user_id}")
        if isinstance(result, list):
            return result
        return []
    
    async def cancel_user_booking(self, booking_id: str) -> bool:
        """Отменить бронирование (для обычного пользователя)"""
        data = {"status": "Canceled"}
        result = await self._request("PUT", f"/Booking/{booking_id}", data=data)
        return result is not None
    
    async def confirm_booking(self, booking_id: str, admin_id: str) -> bool:
        """Подтвердить бронирование (только для админа)"""
        data = {
            "admin_id": admin_id,
            "status": "Confirmed"
        }
        result = await self._request("PUT", f"/Booking/{booking_id}", data=data)
        return result is not None
    
    async def reject_booking(self, booking_id: str, admin_id: str) -> bool:
        """Отклонить бронирование (только для админа)"""
        data = {
            "admin_id": admin_id,
            "status": "Canceled"
        }
        result = await self._request("PUT", f"/Booking/{booking_id}", data=data)
        return result is not None
    
    async def get_all_bookings(self) -> List[Dict]:
        """Получить все бронирования (для админов)"""
        result = await self._request("GET", "/Booking")
        if isinstance(result, list):
            return result
        return []
    
    async def link_vk_user(self, vk_id: int, token: str) -> bool:
        """Привязать VK ID к пользователю"""
        data = {"vk_id": int(vk_id), "token": token}
        result = await self._request("PATCH", "/User/vk", data=data)
        return result is not None


class EquipmentKeyboardPaginator:
    
    def __init__(self, items_per_page: int = 8):
        self.items_per_page = items_per_page
        self.user_pages = {}
        self.user_equipment_cache = {}
    
    def create_keyboard(self, equipment_list: List[Dict], page: int = 0) -> VkKeyboard:
        keyboard = VkKeyboard(one_time=False)
        
        start_idx = page * self.items_per_page
        end_idx = min(start_idx + self.items_per_page, len(equipment_list))
        
        for eq in equipment_list[start_idx:end_idx]:
            status_icon = "✅" if eq.get('is_active', True) else "❌"
            button_text = f"{status_icon} {eq.get('name', 'Без названия')[:35]}"
            keyboard.add_button(button_text, color=VkKeyboardColor.SECONDARY)
            keyboard.add_line()
        
        has_prev = page > 0
        has_next = end_idx < len(equipment_list)
        
        if has_prev or has_next:
            if has_prev:
                keyboard.add_button("◀ Назад", color=VkKeyboardColor.PRIMARY)
            if has_next:
                keyboard.add_button("Вперед ▶", color=VkKeyboardColor.PRIMARY)
            keyboard.add_line()
        
        keyboard.add_button("🔙 Главное меню", color=VkKeyboardColor.NEGATIVE)
        
        return keyboard
    
    def get_page_info(self, vk_id: int, total_items: int) -> tuple:
        current_page = self.user_pages.get(vk_id, 0)
        total_pages = (total_items + self.items_per_page - 1) // self.items_per_page
        return current_page, total_pages
    
    def next_page(self, vk_id: int):
        current = self.user_pages.get(vk_id, 0)
        self.user_pages[vk_id] = current + 1
    
    def prev_page(self, vk_id: int):
        current = self.user_pages.get(vk_id, 0)
        if current > 0:
            self.user_pages[vk_id] = current - 1
    
    def reset_page(self, vk_id: int):
        self.user_pages[vk_id] = 0
    
    def get_current_page(self, vk_id: int) -> int:
        return self.user_pages.get(vk_id, 0)
    
    def cache_equipment(self, vk_id: int, equipment_list: List[Dict]):
        self.user_equipment_cache[vk_id] = equipment_list
    
    def get_cached_equipment(self, vk_id: int) -> Optional[List[Dict]]:
        return self.user_equipment_cache.get(vk_id)


class VKBookingBot:
    def __init__(self):
        self.vk_session = vk_api.VkApi(token=VK_TOKEN)
        self.vk = self.vk_session.get_api()
        self.longpoll = VkBotLongPoll(self.vk_session, GROUP_ID)
        self.api = EquipmentAPIClient(API_BASE_URL, API_TIMEOUT)
        self.paginator = EquipmentKeyboardPaginator(items_per_page=8)
        
        self.user_states = {}
        self.user_data_cache = {}
        self.pending_bookings = []
        self.selected_booking = None
    
    def get_user_id(self, vk_id: int) -> Optional[str]:
        user = self.user_data_cache.get(vk_id)
        if not user:
            return None
        return user.get('id')
    
    def is_admin(self, vk_id: int) -> bool:
        user = self.user_data_cache.get(vk_id)
        if not user:
            return False
        role = user.get('role', 'student')
        return role in ['Admin', 'Super_Admin']
    
    def create_main_keyboard(self, vk_id: int) -> VkKeyboard:
        keyboard = VkKeyboard(one_time=False)
        keyboard.add_button("📋 Список оборудования", color=VkKeyboardColor.PRIMARY)
        keyboard.add_line()
        keyboard.add_button("📅 Мои бронирования", color=VkKeyboardColor.SECONDARY)
        
        if self.is_admin(vk_id):
            keyboard.add_line()
            keyboard.add_button("👑 Все бронирования", color=VkKeyboardColor.PRIMARY)
            keyboard.add_line()
            keyboard.add_button("⏳ Управление заявками", color=VkKeyboardColor.POSITIVE)
        
        return keyboard
    
    async def send_message(self, peer_id: int, message: str, keyboard: VkKeyboard = None):
        try:
            params = {
                "peer_id": peer_id,
                "message": message,
                "random_id": 0
            }
            if keyboard:
                params["keyboard"] = keyboard.get_keyboard()
            
            self.vk.messages.send(**params)
        except Exception as e:
            logger.error(f"Не удалось отправить сообщение: {e}")
    
    async def link_user_and_finish_auth(self, peer_id: int, vk_id: int, ref_token: str):
        """Привязка пользователя по реферальной ссылке"""
        if vk_id in self.user_data_cache:
            user = self.user_data_cache[vk_id]
            await self.send_message(
                peer_id,
                f"✅ Вы уже авторизованы!\n\nДобро пожаловать, {user.get('first_name', 'Пользователь')}!",
                keyboard=self.create_main_keyboard(vk_id)
            )
            return
        
        user = await self.api.get_user_by_vk_id(vk_id)
        if user:
            self.user_data_cache[vk_id] = user
            await self.send_message(
                peer_id,
                f"✅ Вы уже авторизованы!\n\nДобро пожаловать, {user.get('first_name', 'Пользователь')}!",
                keyboard=self.create_main_keyboard(vk_id)
            )
            return
        
        success = await self.api.link_vk_user(vk_id, ref_token)
        
        if success:
            await asyncio.sleep(2)
            for attempt in range(5):
                user = await self.api.get_user_by_vk_id(vk_id)
                if user:
                    self.user_data_cache[vk_id] = user
                    role_emoji = "👨‍🎓" if user.get('role') == 'student' else "👑"
                    await self.send_message(
                        peer_id,
                        f"✅ **Авторизация успешна!** {role_emoji}\n\n"
                        f"Добро пожаловать, {user.get('first_name', '')} {user.get('second_name', '')}!\n\n"
                        f"Теперь вы можете бронировать оборудование.",
                        keyboard=self.create_main_keyboard(vk_id)
                    )
                    return
                await asyncio.sleep(1)
        
        await self.send_message(
            peer_id,
            "❌ **Ошибка авторизации!**\n\n"
            "Пожалуйста, попробуйте снова:\n"
            "1. Перейдите в личный кабинет на сайте\n"
            "2. В меню выберите пункт «Бронирование» \n"
            "3. Перейдите по QR-коду\n"
            "4. Отправьте любое сообщение"
        )
    
    async def process_start(self, peer_id: int, vk_id: int, ref_token: str = None):
        """Обработка команды Начать"""
        if vk_id in self.user_data_cache:
            user = self.user_data_cache[vk_id]
            await self.send_message(
                peer_id,
                f"✅ Вы уже авторизованы!\n\nДобро пожаловать, {user.get('first_name', 'Пользователь')}!",
                keyboard=self.create_main_keyboard(vk_id)
            )
            return
        
        user = await self.api.get_user_by_vk_id(vk_id)
        if user:
            self.user_data_cache[vk_id] = user
            await self.send_message(
                peer_id,
                f"✅ Вы уже авторизованы!\n\nДобро пожаловать, {user.get('first_name', 'Пользователь')}!",
                keyboard=self.create_main_keyboard(vk_id)
            )
            return
        
        if ref_token:
            await self.link_user_and_finish_auth(peer_id, vk_id, ref_token)
            return
        
        await self.send_message(
            peer_id,
            "👋 Привет! Я бот для бронирования оборудования.\n\n"
            "Для начала работы необходимо привязать аккаунт.\n\n"
            "1. Перейдите в личный кабинет на сайте\n"
            "2. В меню выберите пункт «Бронирование» \n"
            "3. Перейдите по QR-коду и отправьте любое сообщение\n\n"
            "Если вы изначально перешли по ссылке из QR-кода, то отправьте любое сообщение"
        )
    
    async def show_equipment_list(self, peer_id: int, vk_id: int):
        equipment = await self.api.get_all_equipment()
        
        if not equipment:
            await self.send_message(peer_id, "❌ Оборудование не найдено")
            return
        
        active_equipment = [eq for eq in equipment if eq.get('is_active', True)]
        
        if not active_equipment:
            await self.send_message(peer_id, "❌ Нет доступного оборудования")
            return
        
        self.paginator.cache_equipment(vk_id, active_equipment)
        
        current_page = self.paginator.get_current_page(vk_id)
        total_pages = (len(active_equipment) + self.paginator.items_per_page - 1) // self.paginator.items_per_page
        
        if current_page >= total_pages:
            current_page = total_pages - 1
            self.paginator.user_pages[vk_id] = current_page
        if current_page < 0:
            current_page = 0
            self.paginator.user_pages[vk_id] = 0
        
        keyboard = self.paginator.create_keyboard(active_equipment, current_page)
        
        start_idx = current_page * self.paginator.items_per_page
        end_idx = min(start_idx + self.paginator.items_per_page, len(active_equipment))
        
        message = f"📋 **Оборудование** (стр. {current_page + 1} из {total_pages}):\n\n"
        for eq in active_equipment[start_idx:end_idx]:
            message += f"**{eq.get('name', 'Без названия')}**\n"
            message += f"Тип: {eq.get('type', 'Не указан')}\n"
            message += f"Аудитория: {eq.get('auditory', 'Не указана')}\n\n"
        
        await self.send_message(peer_id, message, keyboard=keyboard)
    
    async def show_my_bookings(self, peer_id: int, vk_id: int):
        if vk_id not in self.user_data_cache:
            await self.send_message(peer_id, "❌ Сначала авторизуйтесь")
            return
    
        user_id = self.get_user_id(vk_id)
        if not user_id:
            await self.send_message(peer_id, "❌ Ошибка: ID пользователя не найден")
            return
    
        bookings = await self.api.get_user_bookings(user_id)
    
        if not bookings:
            await self.send_message(peer_id, "📭 У вас нет бронирований")
            return
    
        message = "📅 **Ваши бронирования:**\n\n"
        for booking in bookings:
            equipment_id = booking.get('EquipmentId')
            equipment = await self.api.get_equipment_by_id(equipment_id)
        
            if equipment:
                equipment_name = equipment.get('name', 'Неизвестно')
                auditory = equipment.get('auditory', 'Не указана')
            else:
                equipment_name = "Неизвестное оборудование"
                auditory = "Не указана"
        
            status_map = {
                'Waiting answer': '⏳',
                'Confirmed': '✅',
                'Active': '🔵',
                'Completed': '✔️',
                'Canceled': '❌'
            }
            status_icon = status_map.get(booking.get('Status'), '❓')
        
            start_str = booking.get('BookStart', '')
            end_str = booking.get('BookEnd', '')
        
            try:
                if start_str:
                    start_display = datetime.fromisoformat(start_str.replace('Z', '+00:00')).strftime('%d.%m.%Y %H:%M')
                else:
                    start_display = 'Не указано'
            except:
                start_display = start_str[:16] if start_str else 'Не указано'
        
            try:
                if end_str:
                    end_display = datetime.fromisoformat(end_str.replace('Z', '+00:00')).strftime('%d.%m.%Y %H:%M')
                else:
                    end_display = 'Не указано'
            except:
                end_display = end_str[:16] if end_str else 'Не указано'
        
            message += f"{status_icon} **{equipment_name}**\n"
            message += f"📅 {start_display} - {end_display}\n"
            message += f"📍 {auditory}\n"
            message += f"🏷 Статус: {booking.get('Status')}\n"
            message += "-" * 30 + "\n"
    
        await self.send_message(peer_id, message)
    
    async def cancel_booking_flow(self, peer_id: int, vk_id: int):
        if vk_id not in self.user_data_cache:
            await self.send_message(peer_id, "❌ Сначала авторизуйтесь")
            return
    
        user_id = self.get_user_id(vk_id)
        if not user_id:
            await self.send_message(peer_id, "❌ Ошибка: ID пользователя не найден")
            return
    
        bookings = await self.api.get_user_bookings(user_id)
    
        if not bookings:
            await self.send_message(peer_id, "📭 Нет бронирований")
            return
    
        cancellable_bookings = [b for b in bookings if b.get('Status') in ['Waiting answer', 'Confirmed', 'Active']]
    
        if not cancellable_bookings:
            await self.send_message(peer_id, "📭 Нет бронирований, которые можно отменить")
            return
    
        keyboard = VkKeyboard(one_time=True)
        for booking in cancellable_bookings[:10]:
            equipment_id = booking.get('EquipmentId')
            equipment = await self.api.get_equipment_by_id(equipment_id)
            if equipment:
                status_emoji = "⏳" if booking.get('Status') == 'Waiting answer' else "✅"
                button_text = f"{status_emoji} {equipment.get('name')[:30]}"
                keyboard.add_button(button_text, color=VkKeyboardColor.NEGATIVE)
                keyboard.add_line()
    
        keyboard.add_button("🔙 Назад", color=VkKeyboardColor.PRIMARY)
    
        self.user_states[vk_id] = {
            'state': 'cancelling_booking',
            'bookings': cancellable_bookings
        }
    
        await self.send_message(peer_id, "❌ **Выберите бронирование для отмены:**", keyboard=keyboard)
    
    async def process_booking_cancellation(self, peer_id: int, vk_id: int, text: str):
        if vk_id not in self.user_states:
            return
    
        state = self.user_states[vk_id]
        if state.get('state') != 'cancelling_booking':
            return
    
        if text == "🔙 Назад":
            del self.user_states[vk_id]
            await self.send_message(peer_id, "❌ Отмена бронирования отменена")
            return
    
        clean_text = text
        if clean_text.startswith("⏳ ") or clean_text.startswith("✅ ") or clean_text.startswith("❌ "):
            clean_text = clean_text[2:]
    
        for booking in state['bookings']:
            equipment_id = booking.get('EquipmentId')
            equipment = await self.api.get_equipment_by_id(equipment_id)
            if equipment and equipment.get('name') == clean_text:
                booking_id = booking.get('ID')
                success = await self.api.cancel_user_booking(booking_id)
                if success:
                    await self.send_message(peer_id, f"✅ Бронирование '{equipment.get('name')}' успешно отменено")
                    self.user_data_cache.pop(vk_id, None)
                    user = await self.api.get_user_by_vk_id(vk_id)
                    if user:
                        self.user_data_cache[vk_id] = user
                else:
                    await self.send_message(peer_id, f"❌ Ошибка при отмене бронирования '{equipment.get('name')}'")
            
                del self.user_states[vk_id]
                return
    
        await self.send_message(peer_id, "❌ Бронирование не найдено")
    
    async def handle_equipment_selection(self, peer_id: int, vk_id: int, text: str):
        equipment = self.paginator.get_cached_equipment(vk_id)
        if not equipment:
            await self.send_message(peer_id, "❌ Список оборудования не найден. Запросите список заново.")
            return
        
        if text == "Вперед ▶":
            self.paginator.next_page(vk_id)
            await self.show_equipment_list(peer_id, vk_id)
            return
        
        elif text == "◀ Назад":
            self.paginator.prev_page(vk_id)
            await self.show_equipment_list(peer_id, vk_id)
            return
        
        clean_text = text
        if clean_text.startswith("✅ ") or clean_text.startswith("❌ "):
            clean_text = clean_text[2:]
        
        selected = next((eq for eq in equipment if eq.get('name') == clean_text), None)
        if selected:
            await self.request_booking_dates(peer_id, vk_id, selected)
        else:
            await self.send_message(peer_id, "❌ Оборудование не найдено")
    
    async def request_booking_dates(self, peer_id: int, vk_id: int, equipment: Dict):
        self.user_states[vk_id] = {
            'state': 'booking_dates',
            'equipment': equipment
        }
        
        keyboard = VkKeyboard(one_time=False)
        
        today = datetime.now().date()
        for i in range(11):
            date = today + timedelta(days=i)
            display_date = date.strftime("%d.%m.%Y")
            
            if i == 0:
                display_date = f"📅 Сегодня ({display_date})"
            elif i == 1:
                display_date = f"📅 Завтра ({display_date})"
            else:
                display_date = f"📅 {display_date}"
            
            keyboard.add_button(display_date, color=VkKeyboardColor.SECONDARY)
            
            if i % 2 == 1:
                keyboard.add_line()
        
        keyboard.add_line()
        keyboard.add_button("🔙 Отмена", color=VkKeyboardColor.NEGATIVE)
        
        await self.send_message(
            peer_id,
            f"📅 **Бронирование: {equipment.get('name', 'Оборудование')}**\n\n"
            f"Выберите дату начала (время будет 12:00):",
            keyboard=keyboard
        )
    
    async def process_booking_dates(self, peer_id: int, vk_id: int, text: str):
        if vk_id not in self.user_states:
            return
        
        state = self.user_states[vk_id]
        if state.get('state') != 'booking_dates':
            return
        
        if text == "🔙 Отмена":
            del self.user_states[vk_id]
            await self.send_message(peer_id, "❌ Бронирование отменено")
            return
        
        date_match = re.search(r'(\d{2}\.\d{2}\.\d{4})', text)
        if not date_match:
            await self.send_message(peer_id, "❌ Пожалуйста, выберите дату из кнопок")
            return
        
        date_str = date_match.group(1)
        
        try:
            selected_date = datetime.strptime(date_str, "%d.%m.%Y").date()
            start_time = datetime.combine(selected_date, datetime.strptime("12:00", "%H:%M").time())
            
            if start_time < datetime.now():
                await self.send_message(peer_id, "❌ Нельзя бронировать в прошлом. Выберите будущую дату.")
                return
            
            state['start_time'] = start_time
            state['state'] = 'booking_duration'
            
            keyboard = VkKeyboard(one_time=False)
            for days in [1, 2, 3, 4, 5, 6, 7]:
                keyboard.add_button(f"{days} день", color=VkKeyboardColor.SECONDARY)
                if days == 4:
                    keyboard.add_line()
            keyboard.add_line()
            keyboard.add_button("🔙 Отмена", color=VkKeyboardColor.NEGATIVE)
            
            await self.send_message(
                peer_id,
                f"⏰ **Выберите продолжительность бронирования:**\n\n"
                f"📅 Начало: {start_time.strftime('%d.%m.%Y в 12:00')}\n\n"
                f"⚠️ Максимум 7 дней",
                keyboard=keyboard
            )
        except Exception as e:
            logger.error(f"Ошибка парсинга даты: {e}")
            await self.send_message(peer_id, "❌ Неверная дата")
    
    async def process_booking_duration(self, peer_id: int, vk_id: int, text: str):
        if vk_id not in self.user_states:
            return
    
        state = self.user_states[vk_id]
        if state.get('state') != 'booking_duration':
            return
    
        if text == "🔙 Отмена":
            del self.user_states[vk_id]
            await self.send_message(peer_id, "❌ Бронирование отменено")
            return
    
        try:
            days = int(text.split()[0])
        
            if 1 <= days <= 7:
                end_time = state['start_time'] + timedelta(days=days)
            
                user_id = self.get_user_id(vk_id)
                if not user_id:
                    await self.send_message(peer_id, "❌ Ошибка: ID пользователя не найден")
                    del self.user_states[vk_id]
                    return
            
                equipment = state['equipment']
            
                start_with_tz = state['start_time'].isoformat() + '+03:00'
                end_with_tz = end_time.isoformat() + '+03:00'
            
                booking = await self.api.create_booking(
                    user_id,
                    equipment.get('id'),
                    start_with_tz,
                    end_with_tz
                )
            
                if booking:
                    await self.send_message(
                        peer_id,
                        f"✅ **Бронирование создано!**\n\n"
                        f"Оборудование: {equipment.get('name')}\n"
                        f"Начало: {state['start_time'].strftime('%d.%m.%Y в 12:00')}\n"
                        f"Конец: {end_time.strftime('%d.%m.%Y в 12:00')}\n"
                        f"Статус: Ожидает подтверждения администратора",
                        keyboard=self.create_main_keyboard(vk_id)
                    )
                else:
                    await self.send_message(
                        peer_id,
                        "❌ Ошибка создания бронирования.\n"
                        "Возможно, оборудование уже забронировано на это время."
                    )
            
                del self.user_states[vk_id]
            else:
                await self.send_message(peer_id, "❌ Количество дней должно быть от 1 до 7")
        except (ValueError, IndexError):
            await self.send_message(peer_id, "❌ Пожалуйста, выберите количество дней из меню")
    
    async def admin_show_all_bookings(self, peer_id: int, vk_id: int):
        if not self.is_admin(vk_id):
            await self.send_message(peer_id, "❌ Доступно только администраторам")
            return
    
        bookings = await self.api.get_all_bookings()
    
        if not bookings:
            await self.send_message(peer_id, "📭 Нет бронирований")
            return
    
        users_cache = {}
    
        message = "👑 **Все бронирования:**\n\n"
        for booking in bookings[:20]:
            equipment_id = booking.get('EquipmentId')
            equipment = await self.api.get_equipment_by_id(equipment_id)
            equipment_name = equipment.get('name') if equipment else "Неизвестно"
        
            user_id_str = booking.get('UserId')
            user_name = "Неизвестный пользователь"
        
            if user_id_str in users_cache:
                user_name = users_cache[user_id_str]
            else:
                user_info = await self.api.get_user_by_id(user_id_str)
                if user_info:
                    first_name = user_info.get('first_name', '')
                    second_name = user_info.get('second_name', '')
                    if first_name or second_name:
                        user_name = f"{first_name} {second_name}".strip()
                    else:
                        user_name = user_id_str[:8]
                else:
                    user_name = user_id_str[:8]
                users_cache[user_id_str] = user_name
        
            start_str = booking.get('BookStart', '')
            end_str = booking.get('BookEnd', '')
        
            status_map = {
                'Waiting answer': '⏳ Ожидает',
                'Confirmed': '✅ Подтверждено',
                'Active': '🔵 Активно',
                'Completed': '✔️ Завершено',
                'Canceled': '❌ Отменено'
            }
            status_text = status_map.get(booking.get('Status'), booking.get('Status'))
        
            message += f"**{equipment_name}**\n"
            message += f"👤 Пользователь: {user_name}\n"
            message += f"📅 {start_str[:16] if start_str else 'Не указано'} - {end_str[:16] if end_str else 'Не указано'}\n"
            message += f"🏷 Статус: {status_text}\n"
            message += "-" * 30 + "\n"
    
        keyboard = VkKeyboard(one_time=False)
        keyboard.add_button("🔙 Главное меню", color=VkKeyboardColor.NEGATIVE)
    
        await self.send_message(peer_id, message, keyboard=keyboard)
    
    async def admin_manage_bookings(self, peer_id: int, vk_id: int):
        """Показать список ожидающих бронирований для управления"""
        if not self.is_admin(vk_id):
            await self.send_message(peer_id, "❌ Доступно только администраторам")
            return

        bookings = await self.api.get_all_bookings()

        if not bookings:
            await self.send_message(peer_id, "📭 Нет бронирований")
            return

        waiting_bookings = [b for b in bookings if b.get('Status') == 'Waiting answer']

        if not waiting_bookings:
            await self.send_message(peer_id, "📭 Нет заявок, ожидающих подтверждения")
            return

        users_cache = {}

        message = "⏳ **Заявки на бронирование:**\n\n"

        for idx, booking in enumerate(waiting_bookings[:15], 1):
            equipment_id = booking.get('EquipmentId')
            equipment = await self.api.get_equipment_by_id(equipment_id)
            equipment_name = equipment.get('name') if equipment else "Неизвестно"
    
            user_id_str = booking.get('UserId')
            user_name = "Неизвестный пользователь"
    
            if user_id_str in users_cache:
                user_name = users_cache[user_id_str]
            else:
                user_info = await self.api.get_user_by_id(user_id_str)
                if user_info:
                    first_name = user_info.get('first_name', '')
                    second_name = user_info.get('second_name', '')
                    if first_name or second_name:
                        user_name = f"{first_name} {second_name}".strip()
                    else:
                        user_name = user_id_str[:8]
                else:
                    user_name = user_id_str[:8]
                users_cache[user_id_str] = user_name

            start_str = booking.get('BookStart', '')
            end_str = booking.get('BookEnd', '')

            message += f"{idx}. **{equipment_name}**\n"
            message += f"   👤 Пользователь: {user_name}\n"
            message += f"   📅 {start_str[:16] if start_str else 'Не указано'} - {end_str[:16] if end_str else 'Не указано'}\n\n"

        keyboard = VkKeyboard(one_time=True)
    
        buttons_in_row = 0
        max_buttons_in_row = 4
    
        for i in range(1, min(len(waiting_bookings), 13) + 1):
            keyboard.add_button(str(i), color=VkKeyboardColor.PRIMARY)
            buttons_in_row += 1
        
            if buttons_in_row >= max_buttons_in_row:
                keyboard.add_line()
                buttons_in_row = 0
    
        if buttons_in_row > 0:
            keyboard.add_line()
    
        keyboard.add_button("🔙 Главное меню", color=VkKeyboardColor.NEGATIVE)

        self.pending_bookings = waiting_bookings
        self.user_states[vk_id] = {'state': 'admin_selecting_booking'}

        await self.send_message(peer_id, message, keyboard=keyboard)
    
    async def admin_booking_action(self, peer_id: int, vk_id: int, text: str):
        """Обработка выбора заявки для подтверждения/отклонения"""
        if not self.is_admin(vk_id):
            await self.send_message(peer_id, "❌ Доступно только администраторам")
            return

        if text == "🔙 Главное меню":
            if vk_id in self.user_states:
                del self.user_states[vk_id]
            await self.send_message(peer_id, "Главное меню", keyboard=self.create_main_keyboard(vk_id))
            return

        try:
            idx = int(text) - 1
            if 0 <= idx < len(self.pending_bookings):
                selected_booking = self.pending_bookings[idx]
            else:
                await self.send_message(peer_id, "❌ Неверный номер заявки")
                return
        except ValueError:
            await self.send_message(peer_id, "❌ Пожалуйста, выберите номер заявки из меню")
            return

        self.selected_booking = selected_booking
        
        equipment_id = selected_booking.get('EquipmentId')
        equipment = await self.api.get_equipment_by_id(equipment_id)
        equipment_name = equipment.get('name') if equipment else "Неизвестно"
        
        user_id_str = selected_booking.get('UserId')
        user_info = await self.api.get_user_by_id(user_id_str)
        user_name = "Неизвестный"
        if user_info:
            first_name = user_info.get('first_name', '')
            second_name = user_info.get('second_name', '')
            user_name = f"{first_name} {second_name}".strip()
        
        self.user_states[vk_id] = {'state': 'admin_choosing_action'}
        
        keyboard = VkKeyboard(one_time=True)
        keyboard.add_button("✅ Подтвердить", color=VkKeyboardColor.POSITIVE)
        keyboard.add_button("❌ Отклонить", color=VkKeyboardColor.NEGATIVE)
        keyboard.add_line()
        keyboard.add_button("🔙 Назад", color=VkKeyboardColor.PRIMARY)
    
        await self.send_message(
            peer_id,
            f"**Вы выбрали заявку:**\n\n"
            f"🖥 Оборудование: {equipment_name}\n"
            f"👤 Пользователь: {user_name}\n"
            f"📅 Начало: {selected_booking.get('BookStart', '')[:16]}\n"
            f"📅 Конец: {selected_booking.get('BookEnd', '')[:16]}\n\n"
            f"Выберите действие:",
            keyboard=keyboard
        )
    
    async def admin_process_action(self, peer_id: int, vk_id: int, text: str):
        """Подтверждение или отклонение заявки"""
        if not self.is_admin(vk_id):
            await self.send_message(peer_id, "❌ Доступно только администраторам")
            return

        if not self.selected_booking:
            await self.send_message(peer_id, "❌ Заявка не выбрана")
            return

        booking = self.selected_booking
        booking_id = booking.get('ID')
        equipment_id = booking.get('EquipmentId')
        equipment = await self.api.get_equipment_by_id(equipment_id)
        equipment_name = equipment.get('name') if equipment else "Неизвестно"
        
        admin_user = self.user_data_cache.get(vk_id)
        admin_id = admin_user.get('id') if admin_user else None
        
        if not admin_id:
            await self.send_message(peer_id, "❌ Ошибка: ID администратора не найден")
            return

        if text == "✅ Подтвердить":
            success = await self.api.confirm_booking(booking_id, admin_id)
            if success:
                await self.send_message(peer_id, f"✅ Бронирование '{equipment_name}' подтверждено!")
            else:
                await self.send_message(peer_id, f"❌ Ошибка при подтверждении бронирования")

        elif text == "❌ Отклонить":
            success = await self.api.reject_booking(booking_id, admin_id)
            if success:
                await self.send_message(peer_id, f"❌ Бронирование '{equipment_name}' отклонено")
            else:
                await self.send_message(peer_id, f"❌ Ошибка при отклонении бронирования")

        elif text == "🔙 Назад":
            self.selected_booking = None
            if vk_id in self.user_states:
                del self.user_states[vk_id]
            await self.admin_manage_bookings(peer_id, vk_id)
            return

        self.selected_booking = None
        if vk_id in self.user_states:
            del self.user_states[vk_id]

        await self.admin_manage_bookings(peer_id, vk_id)
    
    async def handle_message(self, event):
        peer_id = event.obj.message['peer_id']
        vk_id = event.obj.message['from_id']
        text = event.obj.message['text'].strip()

        last_check = getattr(self, '_last_role_check', {}).get(vk_id, 0)

        if time.time() - last_check > 300:
            fresh_user = await self.api.get_user_by_vk_id(vk_id, force_refresh=True)
            if fresh_user:
                self.user_data_cache[vk_id] = fresh_user
            self._last_role_check = getattr(self, '_last_role_check', {})
            self._last_role_check[vk_id] = time.time()

        ref_token = event.obj.message.get('payload') or event.obj.message.get('ref')

        if vk_id in self.user_states:
            state = self.user_states[vk_id].get('state')
            if state == 'admin_selecting_booking':
                await self.admin_booking_action(peer_id, vk_id, text)
                return
            elif state == 'admin_choosing_action':
                await self.admin_process_action(peer_id, vk_id, text)
                return
    
        if text.startswith("✅ ") or text.startswith("❌ "):
            if self.paginator.get_cached_equipment(vk_id):
                await self.handle_equipment_selection(peer_id, vk_id, text)
            else:
                await self.show_equipment_list(peer_id, vk_id)
            return
    
        if text == "Вперед ▶" or text == "◀ Назад":
            if self.paginator.get_cached_equipment(vk_id):
                await self.handle_equipment_selection(peer_id, vk_id, text)
            else:
                await self.show_equipment_list(peer_id, vk_id)
            return
    
        if vk_id in self.user_states:
            state = self.user_states[vk_id]['state']
        
            if state == 'selecting_equipment':
                await self.handle_equipment_selection(peer_id, vk_id, text)
                return
            elif state == 'booking_dates':
                await self.process_booking_dates(peer_id, vk_id, text)
                return
            elif state == 'booking_duration':
                await self.process_booking_duration(peer_id, vk_id, text)
                return
            elif state == 'cancelling_booking':
                await self.process_booking_cancellation(peer_id, vk_id, text)
                return
    
        if text == "🔙 Главное меню":
            await self.send_message(peer_id, "Главное меню", keyboard=self.create_main_keyboard(vk_id))
            return
    
        if vk_id in self.user_data_cache:
            if text in ["📋 Список оборудования", "список оборудования", "список"]:
                self.paginator.reset_page(vk_id)
                await self.show_equipment_list(peer_id, vk_id)
            elif text in ["📅 Мои бронирования", "мои бронирования"]:
                await self.show_my_bookings(peer_id, vk_id)
            elif text in ["👑 Все бронирования", "все бронирования"] and self.is_admin(vk_id):
                await self.admin_show_all_bookings(peer_id, vk_id)
            elif text in ["⏳ Управление заявками", "управление заявками"] and self.is_admin(vk_id):
                await self.admin_manage_bookings(peer_id, vk_id)
            else:
                await self.send_message(peer_id, "❓ Неизвестная команда. Используйте кнопки меню.", keyboard=self.create_main_keyboard(vk_id))
            return
    
        # Авторизация
        if ref_token:
            await self.link_user_and_finish_auth(peer_id, vk_id, ref_token)
            return
    
        if text.lower() in ["start", "/start", "начать"]:
            await self.process_start(peer_id, vk_id, None)
            return
    
        # Попытка получить пользователя по VK ID
        user = await self.api.get_user_by_vk_id(vk_id)
        if user:
            self.user_data_cache[vk_id] = user
            await self.send_message(
                peer_id,
                f"✅ Вы уже авторизованы!\n\nДобро пожаловать, {user.get('first_name', 'Пользователь')}!",
                keyboard=self.create_main_keyboard(vk_id)
            )
            return
    
        await self.send_message(
            peer_id,
            "🔐 **Для работы с ботом необходима авторизация на сайте**\n\n"
            "1. Перейдите в личный кабинет на сайте\n"
            "2. В меню выберите пункт «Бронирование» \n"
            "3. Перейдите по QR-коду\n\n"
            "После перехода отправьте любое сообщение."
        )
    
    async def run(self):
        print("🚀 Бот запущен. Ожидание сообщений...")
        
        for event in self.longpoll.listen():
            if event.type == VkBotEventType.MESSAGE_NEW:
                try:
                    await self.handle_message(event)
                except Exception as e:
                    print(f"❌ Ошибка: {e}")
                    import traceback
                    traceback.print_exc()


def main():
    bot = VKBookingBot()
    asyncio.run(bot.run())


if __name__ == "__main__":
    main()