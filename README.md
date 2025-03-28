# Накопительная система лояльности «Гофермарт»

## Общие требования

Система представляет собой HTTP API со следующими требованиями к бизнес-логике:

* регистрация, аутентификация и авторизация пользователей;
* приём номеров заказов от зарегистрированных пользователей;
* учёт и ведение списка переданных номеров заказов зарегистрированного пользователя;
* учёт и ведение накопительного счёта зарегистрированного пользователя;
* проверка принятых номеров заказов через систему расчёта баллов лояльности;
* начисление за каждый подходящий номер заказа положенного вознаграждения на счёт лояльности пользователя.

![image](https://pictures.s3.yandex.net:443/resources/gophermart2x_1634502166.png)

### Абстрактная схема взаимодействия с системой

Ниже представлена абстрактная бизнес-логика взаимодействия пользователя с системой:

1. Пользователь регистрируется в системе лояльности «Гофермарт».
2. Пользователь совершает покупку в интернет-магазине «Гофермарт».
3. Заказ попадает в систему расчёта баллов лояльности.
4. Пользователь передаёт номер совершённого заказа в систему лояльности.
5. Система связывает номер заказа с пользователем и сверяет номер с системой расчёта баллов лояльности.
6. При наличии положительного расчёта баллов лояльности производится начисление баллов лояльности на счёт пользователя.
7. Пользователь списывает доступные баллы лояльности для частичной или полной оплаты последующих заказов в интернет-магазине «Гофермарт».

Примечания:

- пункт 2 представлен как гипотетический и не требует реализации в данной работе;
- пункт 3 реализован в системе расчёта баллов лояльности и не требует реализации в данной работе.

### Система расчета баллов лояльности

Система расчета баллов лояльности является внешним сервисом в доверенном контуре. Он работает по принципу чёрного ящика и недоступен для инспекции внешними клиентами. Система рассчитывает положенные баллы лояльности за совершённый заказ по сложным алгоритмам, которые могут меняться в любой момент времени.

Внешнему потребителю доступна только информация о количестве положенных за конкретный заказ баллов лояльности. Причины наличия или отсутствия начислений внешнему потребителю неизвестны.

Протокол взаимодействия с сервисом базы будет предоставлен в конце.

## Сводное HTTP API

Накопительная система лояльности «Гофермарт» должна предоставлять следующие HTTP-хендлеры:

* `POST /api/user/register` — регистрация пользователя;
* `POST /api/user/login` — аутентификация пользователя;
* `POST /api/user/orders` — загрузка пользователем номера заказа для расчёта;
* `GET /api/user/orders` — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
* `GET /api/user/balance` — получение текущего баланса счёта баллов лояльности пользователя;
* `POST /api/user/balance/withdraw` — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
* `GET /api/user/withdrawals` — получение информации о выводе средств с накопительного счёта пользователем.

### Общие ограничения и требования

* хранилище данных — PostgreSQL;
* структура таблиц остаётся на усмотрение студента;
* типы и формат хранения данных (в том числе паролей и прочей чувствительной информации) остаётся на усмотрение студента;
* клиент может поддерживать HTTP-запросы/ответы со сжатием данных;
* клиент не обязан делать запросы соответственно нижеизложенной спецификации API, любая проверка запроса остаётся на усмотрение студента;
* формат и алгоритм проверки аутентификации и авторизации пользователя остаётся на усмотрение студента;
* номера заказов уникальны и никогда не повторяются;
* номер заказа может быть принят в обработку только один раз от одного пользователя;
* номер заказа может не иметь никакого начисления;
* вознаграждение начисляется и тратится в виртуальных баллах из расчёта 1 балл = 1 рубль.

### **Регистрация пользователя**

Хендлер: `POST /api/user/register`.

Регистрация производится по паре логин/пароль. Каждый логин должен быть уникальным.
После успешной регистрации должна происходить автоматическая аутентификация пользователя.

Формат запроса:

```
POST /api/user/register HTTP/1.1
Content-Type: application/json
...

{
	"login": "<login>",
	"password": "<password>"
}
```

Возможные коды ответа:

- `200` — пользователь успешно зарегистрирован и аутентифицирован;
- `400` — неверный формат запроса;
- `409` — логин уже занят;
- `500` — внутренняя ошибка сервера.

### **Аутентификация пользователя**

Хендлер: `POST /api/user/login`.

Аутентификация производится по паре логин/пароль.

Формат запроса:

```
POST /api/user/login HTTP/1.1
Content-Type: application/json
...

{
	"login": "<login>",
	"password": "<password>"
}
```

Возможные коды ответа:

- `200` — пользователь успешно аутентифицирован;
- `400` — неверный формат запроса;
- `401` — неверная пара логин/пароль;
- `500` — внутренняя ошибка сервера.

### **Загрузка номера заказа**

Хендлер: `POST /api/user/orders`.

Хендлер доступен только аутентифицированным пользователям. Номером заказа является последовательность цифр произвольной длины.

Номер заказа может быть проверен на корректность ввода с помощью [алгоритма Луна](https://ru.wikipedia.org/wiki/Алгоритм_Луна).

Формат запроса:

```
POST /api/user/orders HTTP/1.1
Content-Type: text/plain
...

12345678903
```

Возможные коды ответа:

- `200` — номер заказа уже был загружен этим пользователем;
- `202` — новый номер заказа принят в обработку;
- `400` — неверный формат запроса;
- `401` — пользователь не аутентифицирован;
- `409` — номер заказа уже был загружен другим пользователем;
- `422` — неверный формат номера заказа;
- `500` — внутренняя ошибка сервера.

### **Получение списка загруженных номеров заказов**

Хендлер: `GET /api/user/orders`.

Хендлер доступен только авторизованному пользователю. Номера заказа в выдаче должны быть отсортированы по времени загрузки от самых новых к самым старым. Формат даты — RFC3339.

Доступные статусы обработки расчётов:

- `NEW` — заказ загружен в систему, но не попал в обработку;
- `PROCESSING` — вознаграждение за заказ рассчитывается;
- `INVALID` — система расчёта вознаграждений отказала в расчёте;
- `PROCESSED` — данные по заказу проверены и информация о расчёте успешно получена.

Формат запроса:

```
GET /api/user/orders HTTP/1.1
Content-Length: 0
```

Возможные коды ответа:

- `200` — успешная обработка запроса.

  Формат ответа:

    ```
    200 OK HTTP/1.1
    Content-Type: application/json
    ...
    
    [
    	{
            "number": "9278923470",
            "status": "PROCESSED",
            "accrual": 500,
            "uploaded_at": "2020-12-10T15:15:45+03:00"
        },
        {
            "number": "12345678903",
            "status": "PROCESSING",
            "uploaded_at": "2020-12-10T15:12:01+03:00"
        },
        {
            "number": "346436439",
            "status": "INVALID",
            "uploaded_at": "2020-12-09T16:09:53+03:00"
        }
    ]
    ```

- `204` — нет данных для ответа.
- `401` — пользователь не авторизован.
- `500` — внутренняя ошибка сервера.

### **Получение текущего баланса пользователя**

Хендлер: `GET /api/user/balance`.

Хендлер доступен только авторизованному пользователю. В ответе должны содержаться данные о текущей сумме баллов лояльности, а также сумме использованных за весь период регистрации баллов.

Формат запроса:

```
GET /api/user/balance HTTP/1.1
Content-Length: 0
```

Возможные коды ответа:

- `200` — успешная обработка запроса.

  Формат ответа:

    ```
    200 OK HTTP/1.1
    Content-Type: application/json
    ...
    
    {
    	"current": 500.5,
    	"withdrawn": 42
    }
    ```

- `401` — пользователь не авторизован.
- `500` — внутренняя ошибка сервера.

### **Запрос на списание средств**

Хендлер: `POST /api/user/balance/withdraw`

Хендлер доступен только авторизованному пользователю. Номер заказа представляет собой гипотетический номер нового заказа пользователя в счет оплаты которого списываются баллы.

Примечание: для успешного списания достаточно успешной регистрации запроса, никаких внешних систем начисления не предусмотрено и не требуется реализовывать.

Формат запроса:

```
POST /api/user/balance/withdraw HTTP/1.1
Content-Type: application/json

{
	"order": "2377225624",
    "sum": 751
}
```

Здесь `order` — номер заказа, а `sum` — сумма баллов к списанию в счёт оплаты.

Возможные коды ответа:

- `200` — успешная обработка запроса;
- `401` — пользователь не авторизован;
- `402` — на счету недостаточно средств;
- `422` — неверный номер заказа;
- `500` — внутренняя ошибка сервера.

### **Получение информации о выводе средств**

Хендлер: `GET /api/user/withdrawals`.

Хендлер доступен только авторизованному пользователю. Факты выводов в выдаче должны быть отсортированы по времени вывода от самых новых к самым старым. Формат даты — RFC3339.

Формат запроса:

```
GET /api/user/withdrawals HTTP/1.1
Content-Length: 0
```

Возможные коды ответа:

- `200` — успешная обработка запроса.

  Формат ответа:

    ```
    200 OK HTTP/1.1
    Content-Type: application/json
    ...
    
    [
        {
            "order": "2377225624",
            "sum": 500,
            "processed_at": "2020-12-09T16:09:57+03:00"
        }
    ]
    ```

- `204` - нет ни одного списания.
- `401` — пользователь не авторизован.
- `500` — внутренняя ошибка сервера.

## Взаимодействие с системой расчёта начислений баллов лояльности

Для взаимодействия с системой доступен один хендлер:

- `GET /api/orders/{number}` — получение информации о расчёте начислений баллов лояльности.

Формат запроса:

```
GET /api/orders/{number} HTTP/1.1
Content-Length: 0
```

Возможные коды ответа:

- `200` — успешная обработка запроса.

  Формат ответа:

    ```
    200 OK HTTP/1.1
    Content-Type: application/json
    ...
    
    {
        "order": "<number>",
        "status": "PROCESSED",
        "accrual": 500
    }
    ```

  Поля объекта ответа:

    - `order` — номер заказа;
    - `status` — статус расчёта начисления:

        - `REGISTERED` — заказ зарегистрирован, но не начисление не рассчитано;
        - `INVALID` — заказ не принят к расчёту, и вознаграждение не будет начислено;
        - `PROCESSING` — расчёт начисления в процессе;
        - `PROCESSED` — расчёт начисления окончен;

    - `accrual` — рассчитанные баллы к начислению, при отсутствии начисления — поле отсутствует в ответе.

- `204` - заказ не зарегистрирован в системе расчета.

- `429` — превышено количество запросов к сервису.

  Формат ответа:

    ```
    429 Too Many Requests HTTP/1.1
    Content-Type: text/plain
    Retry-After: 60
    
    No more than N requests per minute allowed
    ```

- `500` — внутренняя ошибка сервера.

Заказ может быть взят в расчёт в любой момент после его совершения. Время выполнения расчёта системой не регламентировано. Статусы `INVALID` и `PROCESSED` являются окончательными.

Общее количество запросов информации о начислении не ограничено.

## Конфигурирование сервиса накопительной системы лояльности

Сервис должн поддерживать конфигурирование следующими методами:

- адрес и порт запуска сервиса: переменная окружения ОС `RUN_ADDRESS` или флаг `-a`
- адрес подключения к базе данных: переменная окружения ОС `DATABASE_URI` или флаг `-d`
- адрес системы расчёта начислений: переменная окружения ОС `ACCRUAL_SYSTEM_ADDRESS` или флаг `-r`
