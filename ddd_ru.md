# Проектирование Магазина мебели

## 0. Event Storming
![ES](https://github.com/p12s/furniture-store/blob/master/images/ES.png?raw=true)  
[Обновляемая схема event storming](https://lucid.app/lucidchart/1482e706-4b6d-49f8-adce-e0b7932d8bbe/edit?viewport_loc=-128%2C-54%2C2307%2C1397%2C0_0&invitationId=inv_dd15d087-fe4e-4cea-b2f5-ce0f5ad99f35)  
![Domain model](https://github.com/p12s/furniture-store/blob/master/images/domain-model.png?raw=true)  
[Обновляемая схема доменной модели](https://www.xmind.net/m/EVD7bc)  
  
## 1. Описываем query/comands (DDD)  
#### Общий сервис аутентификации 
Зарегистрироваться могут все, в одном месте **(command)** создание и возврат "ключа" - это изменение стейта  
```
Actor   Account  
Command Sign up  
Data    Account (login, password, name, email, address)  
Event   Account.SignedUp  
```

#### Заказ    
Видеть товары магазина **(query, read model)**  
Видеть список своих заказанных товаров (корзину, что в ней) **(query, read model)**  
Положить товар в корзину **(command)**  
```
Actor   Account  
Command Add product to cart 
Data    Account (public_id), Product (public_id)     
Event   Order.ProductAdded  
```
Ввести купон на скидку **(command)**  
```
Actor   Account  
Command Add discount coupon 
Data    Account (public_id), Order (public_id)     
Event   Order.DiscountCouponAdded  
```
Оплатить заказ **(command)**  
```
Actor   Account  
Command Pay order
Data    Account (public_id), Product (public_id)     
Event   Order.Payed  
```
Увидеть что заказ передан в доставку **(query, read model)**  
Видеть изменение статуса заказа вплоть до доставки **(query, read model)**  

#### Доставка
Зайти в дашборд с созданными для него админом доступами **(query, read model)**  
Видеть список готовых к доставке товаров (оплаченных заказов) **(query, read model)**  
Взять заказ (видит имя клиента, адрес доставки - для принятия решения) **(command)**  
```
Actor   Account (with delivery role)  
Command Take order to deliver
Data    Account (public_id), Order (public_id) 
Event   Order.TakedToDeliver  
```
Отметить заказ доставленным **(command)**  
```
Actor   Account (with delivery role)  
Command Deliver order
Data    Account (public_id), Order (public_id)  
Event   Order.Delivered  
```

#### Диллер
Видеть дашборд склада **(query, read model)**  
Cоздать товар **(command)**  
```
Actor   Account (with dealer role)  
Command Create product
Data    Account (public_id), Product (data)  
Event   Product.Created  
```
Обновить товар **(command)**  
```
Actor   Account (with dealer role)  
Command Update product
Data    Account (public_id), Product (data)  
Event   Product.Updated  
```
Ввести скидку на товар **(command)**  
```
Actor   Account  
Command Add discount 
Data    Account (public_id), Product (public_id)     
Event   Product.DiscountAdded  
```
Удалить товар **(command)**  
```
Actor   Account (with dealer role)  
Command Delete product
Data    Account (public_id), Product (public_id)   
Event   Product.Deleted  
```

#### Админка
Видеть данные пользователей **(query, read model)**  
Создать пользователя (логин/пароль, имя/фам, мейл, роль) **(command)**  
```
Actor   Account (with admin role)  
Command Create account
Data    Account (data, role)
Event   Account.Created  
```
Поменять пользователю роль **(command)**  
```
Actor   Account (with admin role)  
Command Change account role
Data    Account (data, role)
Event   Account.RoleChanged  
```
Сбросить пароль **(command)**  
```
Actor   Account (with admin role)  
Command Reset password
Data    Account (public_id)
Event   Account.PassowrdReseted  
```
Заблокировать пользователя **(command)**  
```
Actor   Account (with admin role)  
Command Disable account
Data    Account (public_id)   
Event   Account.Disabled  
```

#### Бухгалтерия
Видеть сколько товаров получено от поставщиков **(query, read model)**  
Видеть сколько товаров куплено **(query, read model)**  
Видеть сколько товаров доставлено или ожидает доставки **(query, read model)**  

#### Аналитика
Видеть все по-максимуму, аналитики хотят все знать **(query, read model)**  

#### Cервис уведомлений (не разбираем)
Уведомление об оплате - успешно оплачен, ожидайте доставки  
Уведомление о доставке - курьер взялся доставить (его контакт?)  

## 2. Набрасываем модель данных  
- Account: login, password, name, email, address, role
- Order: account_id, status
- Product: dealer_id, name, price, discount(%) 
  
## 3. Выделяем домены (по акторам/контексту)  
[Обновляемая схема доменной модели](https://www.xmind.net/m/EVD7bc)  
- Auth domain  
    - Account  
        - id  
        - public_id (uuid)  
        - login  
        - pass  
        - role  
        - email
        - address  
- Product domain  
    - [Account - копия, "урезанная версия" домена Auth]  
        - public_id (uuid)
        - role  
    - Product  
        - id  
        - public_id (uuid)  
        - name  
        - price  
        - discount(%)  
- Ordering domain  
    - [Account - копия, "урезанная версия" домена Auth]  
        - public_id (uuid)
        - role  
    - [Product - копия, "урезанная версия" домена Product]  
        - public_id (uuid)
        - name  
        - price ?  
        - discount(%) ?  
    - Order  
        - id  
        - public_id (uuid)  
        - account public_id  (uuid)
        - status  
- Billing domain  
    - [Account - копия, "урезанная версия" домена Auth]  
        - public_id (uuid)
        - role  
    - Billing (auditlog)  
        - account public_id
        - order public_id  
        - status
        - price (from Money credited event?)  

## 4. Примеряемся как разделить домены на сервисы  
- Auth
- Product
- Ordering
- Billing

## 5. Определяем коммуникации между сервисами  
Один домен - один сервис.  
Выбираем **асинхронный** подход с отправкой событий в очередь - сервисы сами будут что-то делать, когда появится событие.    
Обновление данных учетных записей в в сервисах тоже **асинхронное**, считаем что нам не важна небольшая возможная задержка.    
Данные для коммуникаци между сервисами описаны в п.3 "Выделяем домены".  
