# Сервис товаров

## Получение товара

Запрос:
```
GET localhost:8080/products/1
```

Ответ:
```
{
    "product": {
        "sku": 1,
        "price": 100,
        "name": "Крем для лица"
    }
}
```

# Полезные ссылки

+ Effective Go - https://go.dev/doc/effective_go
+ Go Wiki: Common Mistakes https://go.dev/wiki/CommonMistakes
+ Go Wiki: Go Code Review Comments https://go.dev/wiki/CodeReviewComments
+ Organizing a Go module - https://go.dev/doc/modules/layout
+ Project structure - https://github.com/golang-standards/project-layout/blob/master/README_ru.md
+ Uber style guide - https://github.com/uber-go/guide/blob/master/style.md
