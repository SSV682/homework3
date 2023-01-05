# OTUS Course homework  "Microservice architecture" 

## Зависимости

Список зависимостей:

- [Minikube 1.27.1] (https://github.com/kubernetes/minikube/releases/tag/v1.27.1)
- [Kubectl 1.26.0] (https://github.com/kubernetes/kubectl/tree/release-1.26)
- [Istioctl 1.16.0] (https://github.com/istio/istio/releases/tag/1.16.0)
- [Metallb 0.13.7] (https://raw.githubusercontent.com/metallb/metallb/v0.13.7/config/manifests/metallb-native.yaml)

Некоторые операции будут совершаться с помощью утилиты `kubectl`

```shell
minikube start driver=docker
```

Чтобы подключить дополнения, выполните команду::
```shell
minikube addons enable ingress
minikube addons enable ingress-dns
```
Cоздать сетевой маршрут на хосте:
```shell
minikube tunnel
```

## Содержание

* [Описание стенда](#Описание-стенда)
* [Установка Istio](#Установка-Istio)
* [Установка Gateway](#Установка-Gateway)
* [Разворачиваем приложения](#Разворачиваем-приложения)
* [Установка Metallb](#Установка-Metallb)
* [Аутентификация и Авторизация](#Аутентификация-и-авторизация)

## Описание стенда

В кластере развернуто два пользовательских приложения: `auth-service`, `user-service`. А так же,
реализовано API Gateway с помощью Istio


### Регистрация

Для регистрации используется HTTP метод POST. Запросы к `/signup` попадают на `user-service` и не проходят проверку авторизации.

```http request
POST http://{address}/signup HTTP/1.1
```

В качестве входных параметров должно быть указано тело в JSON формате:

| Имя поля    | Тип      | Описание                    | Обязательный | Уникальный |
|-------------|----------|-----------------------------|--------------|------------|
| `username`  | `String` | Уникальное имя пользователя | Да           | Да         |
| `firstname` | `String` | Имя пользователя            | Нет          | Нет        |
| `lastname`  | `String` | Фамилия пользователя        | Нет          | Нет        |
| `email`     | `String` | Почта пользователя          | Нет          | Нет        |
| `phone`     | `String` | Телефон пользователя        | Нет          | Нет        |
| `password`  | `String` | Пароль пользователя         | Да           | Нет        |

**<u>Пример тела запроса</u>**: `json`

```json
{
    "Username": "LyricTurner74",
    "Firstname": "Lyric",
    "Lastname": "Turner",
    "Email": "Lyric_Turner74@yahoo.com",
    "Phone": "567-461-7480",
    "Password": "4LexKr4eV"
}
```

**<u>Пример тела ответа</u>**: `json`
```json
{
    "id": "a61a8a13-bc6c-4feb-bb77-ade9a31e5b63"
}
```

### Авторизация

Для авторизации используется HTTP метод POST. В результате получаем access-токен. Запросы к `/login` попадают на `auth-service` и не проходят проверку авторизации

```http request
POST http://{address}/login HTTP/1.1
```

Параметры запроса:

| Имя        | Положение | Опциональный | Тип      | Описание            |
|------------| --------- |--------------|----------|---------------------|
| `username` | `Query`   | Нет          | `String` | Имя пользователя    |
| `password` | `Query`   | Нет          | `String` | Пароль пользователя |

**<u>Пример тела ответа</u>**: `json`
```json
{
    "accessToken": "eyJhbGciOiJSUzI1NiIsImtpZCI6IlF5aF9OQjAxbklHOEVGUmxNaXdoZGtWMHhBST0iLCJ0eXAiOiJKV1QifQ.ewogImV4cCI6IDE2NzMwMTQ2ODUsCiAiaWF0IjogMTY3MjkyODI4NSwKICJpZF91c2VyIjogImE2MWE4YTEzLWJjNmMtNGZlYi1iYjc3LWFkZTlhMzFlNWI2MyIsCiAiaXNzIjogImh0dHA6Ly91c2Vyc2VydmljZS1hdXRoc2VydmljZS51c2Vyc2VydmljZS5zdmMuY2x1c3Rlci5sb2NhbCIsCiAianRpIjogIjlkODJhMWJlLWRhZDgtNDY2Yy1iNGU4LTg0MTlhY2NjNDJlMyIsCiAibmJmIjogMTY3MjkyODI4NQp9.QFe3zOuMIahirhJFBi5cqdiXKl0JzmusoCDoge5VXMaIDq7G6fSrCKJON64XkHlxa2IVbrblKoo0DcoPfApu41AHhRmOYBPSAxIvckc8ipRYPMOQo6HEbheoJ4FsMwrGJFNmjfK6VUUjzrYN0xClOZjohTNYPnzh_Hq2oczAOXr8VJGudJVW3x7luOIWN5e3aNQNuMWBSgsJM74KvMjrtO4SV3oCRQSCxRcedmXm8s5EACfo7Ucz78oxeYYwcUNuD3hgApx46NRSjhyvc2TKaJfK35gGS1U_AEJvExhJ3X1Ag9wJrNS9jS7jALj8C6I3JWqcHcIiVyxAqJ7Esqnp0A"
}
```

### Получение пользователем информации о себе

Для получения пользователем информации о себе используется HTTP метод GET. Запросы попадают на `user-service` и проходят проверку авторизации с помощью access-токена.

```http request
GET http://{address}/user HTTP/1.1
```
**<u>Пример тела ответа</u>**: `json`
```json
{
    "id": "a61a8a13-bc6c-4feb-bb77-ade9a31e5b63",
    "username": "LyricTurner74",
    "firstname": "Lyric",
    "lastname": "Turner",
    "email": "Lyric_Turner74@yahoo.com",
    "phone": "567-461-7480",
    "password": "4LexKr4eV"
}
```

### Обновление пользователем информации о себе

Для обновления пользователем информации о себе используется HTTP метод PATCH. Запросы попадают на `user-service` и проходят проверку авторизации с помощью access-токена.

```http request
PATCH http://{address}/user HTTP/1.1
```
**<u>Пример тела запроса</u>**: `json`
```json
{
    "id": "a61a8a13-bc6c-4feb-bb77-ade9a31e5b63",
    "username": "LyricTurner74",
    "firstname": "Lyric",
    "lastname": "Turner",
    "email": "Lyric_Turner74@yahoo.com",
    "phone": "567-461-7480",
    "password": "4LexKr4eV"
}
```

**<u>Пример тела ответа</u>**: `json`
```json
{
    "id": "a61a8a13-bc6c-4feb-bb77-ade9a31e5b63",
    "username": "LyricTurner74",
    "firstname": "Lyric",
    "lastname": "Turner",
    "email": "Lyric_Turner74@yahoo.com",
    "phone": "567-461-7480",
    "password": "4LexKr4eV"
}
```

### Удаление пользователем информации о себе

Для удаления пользователем информации о себе используется HTTP метод DELETE. Запросы попадают на `user-service` и проходят проверку авторизации с помощью access-токена.

```http request
DELETE http://{address}/user HTTP/1.1
```

## Установка Istio

Требуется установить istio:

```shell
istioctl install
```

```shell
istioctl operator init --watchedNamespaces istio-system --operatorNamespace istio-operator
```

Конфигурирование Istio с помощью файла манифеста:
```shell
kubectl apply -f ./k8s/apigateway/istio.yaml
```

## Установка Gateway

Установить apigateway с помощью файла манифеста:
```shell
kubectl apply -f ./k8s/apigateway/gateway.yaml
```

## Разворачиваем приложения

Создаем namespace приложения: 
```shell
kubectl create namespace userservice
```
Настроим его для использования с istio:
```shell
kubectl label namespace userservice istio-injection=enabled --overwrite
```
И установим приложение с помощью helm:
```shell
helm install --set name=userservice userservice --namespace userservice ./k8s/userservice
```

Так же требуется посмотреть <EXTERNAL-IP>:
```shell
kubectl get svc istio-ingressgateway -n istio-system
```
И прописать его в /etc/hosts
```shell
sudo nano /etc/hosts
```

## Установка Metallb

Чтобы установить MetalLB, поскольку он состоит из двух частей, нужно развернуть этот ресурс в minikube командой:
```shell
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.7/config/manifests/metallb-native.yaml
```
А затем мы должны выполнить конфигурацию, без которой он не будет работать, поэтому сначала просто примените этот манифест:
```shell
kubectl apply -f ./k8s/apigateway/metallb.yaml
```


## Аутентификация и Авторизация
На схеме представлен процесс аутентификации и авторизации пользователей для получения данных от сервиса пользователей 

<img width="695" alt="schema1" src="https://user-images.githubusercontent.com/16625234/210849244-cd803a43-6b19-44bd-8d66-8aad2a04db0c.png">


