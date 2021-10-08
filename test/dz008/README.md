#НЕ НАЧАТА: [ЧАТ](https://otus.ru/learning/61597/#/homework-chat/12344/) / [ОТЧЕТ](REPORT.md)

-----

# Разделение монолита на сервисы

Цель: В результате выполнения ДЗ вы перенесете бизнес-домен монолитного приложения в отдельный сервис.
###В данном задании тренируются навыки:
- декомпозиции предметной области;
- разделение монолитного приложения;
- работа с HTTP;
- работа с REST API и gRPC;

### План выполнения:
1) Вынести систему диалогов в отдельный сервис.
2) Взаимодействия монолитного сервиса и сервиса чатов реализовать на Rest API или gRPC.
3) Организовать сквозное логирование запросов.
4) Предусмотреть то, что не все клиенты обновляют приложение быстро и кто-то может ходить через старое API.

ДЗ сдается в виде исходного кода на github и отчета по устройству системы.
Критерии оценки: Оценка происходит по принципу зачет/незачет.

### Требования:
- Описан протокол взаимодействия.
- Поддержаны старые клиенты.
- Новые клиенты верно ходят через новый API.

### Рекомендуем сдать до: 10.01.2022