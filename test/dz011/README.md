#НЕ НАЧАТА: [ЧАТ](https://otus.ru/learning/61597/#/homework-chat/12347/) / [ОТЧЕТ](REPORT.md)

-----

# Внедрение docker и consul

Цель: В результате выполнения ДЗ вы интегрируете в ваш проект социальной сети docker и auto discovery сервисов с помощью consul
### В данном задании тренируются навыки:
- использование docker;
- использование consul;
- построение auto discovery;

### План выполнения:
1) Обернуть сервис диалогов в docker
2) Развернуть consul в вашей системе
3) Интегрировать auto discovery в систему диалогов
4) Научить монолитное приложение находить и равномерно нагружать все поднятые узлы сервиса диалогов
5) Опционально можно использовать nomad

6) ДЗ сдается в виде репозитория с исходными кодами на github и отчетом о выполненных шагах.

Критерии оценки: Оценка происходит по принципу зачет/незачет.

### Требования:
- Верно настроен docker.
- Обеспечено распределение нагрузки по экземплярам сервиса.
- Описан процесс развертки новых экземпляров.

### Рекомендуем сдать до: 31.01.2022