# 🧪 Руководство по тестированию API UNO Spicchio

## 📋 Список исправлений

### ✅ Исправленные проблемы:

1. **500 ошибка на PUT /api/orders/{id}/status**
   - Добавлена валидация статусов
   - Улучшена обработка ошибок
   - Добавлены информативные сообщения об ошибках

2. **401 ошибка (Unauthorized)**
   - Улучшена проверка токена в middleware
   - Добавлена нормализация Bearer токена
   - Улучшены сообщения об ошибках

3. **405 ошибка (Method Not Allowed)**
   - Исправлены маршруты в роутере
   - Добавлен Recoverer middleware для отлова паник

---

## 🚀 Быстрый старт

### 1. Замените файлы в проекте:

```bash
# Создайте бэкап текущих файлов
cp internal/controller/http/handlers/order_handler.go internal/controller/http/handlers/order_handler.go.backup
cp internal/usecase/orders.go internal/usecase/orders.go.backup
cp internal/controller/http/middleware/auth.go internal/controller/http/middleware/auth.go.backup
cp internal/controller/http/router.go internal/controller/http/router.go.backup
```

Замените содержимое файлов на исправленные версии из артефактов выше.

### 2. Перезапустите сервер:

```bash
# Остановите текущий сервер (Ctrl+C)
# Запустите снова
make run
```

---

## 🔧 Тестирование в Postman

### Шаг 1: Импортируйте коллекцию

1. Откройте Postman
2. Нажмите **Import**
3. Скопируйте JSON из артефакта "UNO Spicchio - Postman Collection"
4. Вставьте и нажмите **Import**

### Шаг 2: Настройте переменные окружения

Коллекция уже содержит переменные:
- `baseUrl`: `http://localhost:8080/api`
- `token`: автоматически сохраняется после логина

### Шаг 3: Базовый сценарий тестирования

#### 📝 Полный цикл работы с заказом:

**1. Логин как официант:**
```
POST http://localhost:8080/api/auth/login

Body:
{
  "username": "waiter",
  "password": "waiter123"
}
```
✅ Токен автоматически сохранится

---

**2. Посмотреть доступные столы:**
```
GET http://localhost:8080/api/tables
Authorization: Bearer {{token}}
```

---

**3. Посмотреть меню:**
```
GET http://localhost:8080/api/dishes?active=true
Authorization: Bearer {{token}}
```

---

**4. Создать заказ:**
```
POST http://localhost:8080/api/orders
Authorization: Bearer {{token}}

Body:
{
  "table_number": 2,
  "notes": "Customer allergic to nuts",
  "items": [
    {
      "dish_id": 1,
      "qty": 2,
      "notes": "Extra cheese"
    },
    {
      "dish_id": 4,
      "qty": 1,
      "notes": "Medium rare"
    }
  ]
}
```

**Ожидаемый ответ:**
```json
{
  "success": true,
  "data": {
    "id": 6,
    "waiter_id": 3,
    "table_number": 2,
    "status": "new",
    "total": 12500,
    "notes": "Customer allergic to nuts",
    "created_at": "2025-11-01T10:30:00Z",
    "updated_at": "2025-11-01T10:30:00Z"
  }
}
```

---

**5. Логин как повар:**
```
POST http://localhost:8080/api/auth/login

Body:
{
  "username": "cook",
  "password": "cook123"
}
```

---

**6. Посмотреть новые заказы:**
```
GET http://localhost:8080/api/orders?status=new
Authorization: Bearer {{token}}
```

---

**7. Взять заказ в работу:**
```
PUT http://localhost:8080/api/orders/6/status
Authorization: Bearer {{token}}

Body:
{
  "status": "in_progress"
}
```

**Ожидаемый ответ:**
```json
{
  "success": true,
  "data": {
    "message": "order status updated successfully"
  }
}
```

---

**8. Отметить как готовый:**
```
PUT http://localhost:8080/api/orders/6/status
Authorization: Bearer {{token}}

Body:
{
  "status": "ready"
}
```

---

**9. Логин обратно как официант:**
```
POST http://localhost:8080/api/auth/login

Body:
{
  "username": "waiter",
  "password": "waiter123"
}
```

---

**10. Закрыть заказ (оплата):**

**Вариант A - Через /close (рекомендуется):**
```
PUT http://localhost:8080/api/orders/6/close
Authorization: Bearer {{token}}
```

**Вариант B - Через /status:**
```
PUT http://localhost:8080/api/orders/6/status
Authorization: Bearer {{token}}

Body:
{
  "status": "paid"
}
```

---

## ⚠️ Возможные ошибки и решения

### Ошибка: "invalid status change"

**Причина:** Попытка пропустить статус в цепочке

**Правильная последовательность:**
```
new → in_progress → ready → paid
```

**Пример ошибки:**
```
new → ready ❌ (нельзя пропустить in_progress)
```

**Решение:** Меняйте статусы последовательно

---

### Ошибка: "insufficient permissions"

**Причина:** Пользователь не имеет прав на это действие

**Права по ролям:**

| Действие | Admin | Manager | Waiter | Cook |
|----------|-------|---------|--------|------|
| Создать заказ | ✅ | ❌ | ✅ | ❌ |
| Изменить статус | ✅ | ❌ | ❌ | ✅ |
| Закрыть заказ | ✅ | ❌ | ✅ | ❌ |
| Аналитика | ✅ | ✅ | ❌ | ❌ |

**Решение:** Используйте правильный аккаунт для действия

---

### Ошибка: "insufficient stock for order"

**Причина:** На складе недостаточно ингредиентов

**Решение:** Сделайте поставку (требуется admin):
```
POST http://localhost:8080/api/supplies
Authorization: Bearer {{token}} (admin)

Body:
{
  "ingredient_id": 1,
  "qty": 50.0,
  "supplier_name": "Fresh Foods Inc"
}
```

---

### Ошибка: "table not found"

**Причина:** Указан несуществующий стол

**Решение:** Проверьте список столов:
```
GET http://localhost:8080/api/tables
```

Используйте `table.id` (не `table.name`)

---

### Ошибка: "dish not found"

**Причина:** Указан несуществующий `dish_id`

**Решение:** Проверьте список блюд:
```
GET http://localhost:8080/api/dishes
```

---

## 📊 Тестирование аналитики

**Требуется:** Admin или Manager

```
# Логин как admin
POST http://localhost:8080/api/auth/login
Body: {"username": "admin", "password": "admin123"}

# Дашборд за сегодня
GET http://localhost:8080/api/analytics/dashboard?period=today

# Дашборд за текущий месяц
GET http://localhost:8080/api/analytics/dashboard?period=current_month

# Популярные блюда
GET http://localhost:8080/api/analytics/dishes/popular?from=2025-10-01&to=2025-10-31&limit=5

# Продажи по категориям
GET http://localhost:8080/api/analytics/sales/by-category?from=2025-10-01&to=2025-10-31

# Производительность официантов
GET http://localhost:8080/api/analytics/waiters/performance?from=2025-10-01&to=2025-10-31
```

---

## 🧪 Проверка всех эндпоинтов

### Health Check (без авторизации):
```bash
curl http://localhost:8080/health
# Ожидается: {"status":"ok"}
```

### Auth:
- ✅ POST /api/auth/login
- ✅ GET /api/auth/me

### Orders:
- ✅ GET /api/orders
- ✅ GET /api/orders?status=new
- ✅ GET /api/orders/{id}
- ✅ POST /api/orders (waiter/admin)
- ✅ PUT /api/orders/{id}/status (cook/admin)
- ✅ PUT /api/orders/{id}/close (waiter/admin)

### Dishes:
- ✅ GET /api/dishes
- ✅ GET /api/dishes?active=true
- ✅ GET /api/dishes/{id}
- ✅ GET /api/dishes/{id}/ingredients

### Tables:
- ✅ GET /api/tables
- ✅ GET /api/tables/{id}
- ✅ PUT /api/tables/{id}/status (waiter/admin)

### Categories:
- ✅ GET /api/categories
- ✅ GET /api/categories/{id}

### Analytics:
- ✅ GET /api/analytics/dashboard
- ✅ GET /api/analytics/sales/summary
- ✅ GET /api/analytics/sales/by-category
- ✅ GET /api/analytics/dishes/popular
- ✅ GET /api/analytics/waiters/performance

---

## 📝 Логи для проверки

После успешных запросов вы должны видеть в логах:

```
2025/11/01 10:10:38 POST /api/auth/login 200 508ms
2025/11/01 10:11:08 GET /api/orders 200 1.5ms
2025/11/01 10:14:48 PUT /api/orders/1/status 200 4.3ms
2025/11/01 10:15:10 PUT /api/orders/1/close 200 2.1ms
```

**Не должно быть:**
- ❌ 401 (если токен правильный)
- ❌ 405 (неправильный HTTP метод)
- ❌ 500 (внутренняя ошибка сервера)

---

## 🎯 Checklist для проверки

- [ ] Сервер запущен без ошибок подключения к БД
- [ ] Логин работает для всех 4 ролей
- [ ] Создание заказа работает (waiter)
- [ ] Изменение статуса работает (cook)
- [ ] Закрытие заказа работает (waiter)
- [ ] Аналитика работает (admin/manager)
- [ ] Нет 401/405/500 ошибок в логах

---

## 💡 Советы

1. **Используйте Environment Variables в Postman** для быстрого переключения между пользователями
2. **Сохраняйте order_id** после создания для последующих тестов
3. **Проверяйте таблицу orders в БД** для отладки:
   ```sql
   SELECT id, status, table_number, total FROM orders ORDER BY created_at DESC LIMIT 5;
   ```

---

## 🔍 Дополнительная отладка

Если что-то не работает, проверьте:

1. **Токен сохранился?**
   ```
   Postman → Collection Variables → token (должен быть длинный JWT)
   ```

2. **Правильный формат Authorization header?**
   ```
   Authorization: Bearer eyJhbGc...
   ```

3. **БД доступна?**
   ```bash
   psql -h localhost -U restaurant_user -d restaurant_crm
   ```

4. **Есть тестовые данные?**
   ```sql
   SELECT COUNT(*) FROM dishes;
   SELECT COUNT(*) FROM tables;
   SELECT COUNT(*) FROM ingredients;
   ```

---

## 📞 Поддержка

Если проблема остается:
1. Проверьте логи сервера
2. Убедитесь, что используете правильную роль
3. Проверьте формат JSON в Body
4. Убедитесь, что порт 8080 свободен