FROM python:3.11-slim

WORKDIR /app

# Зависимости
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Код бота (ВАЖНО: имя файла точно как на диске)
COPY VKbot.py .

# Запуск
CMD ["python", "VKbot.py"]