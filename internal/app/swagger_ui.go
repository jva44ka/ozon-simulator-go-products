package app

const swaggerUiHtml = `
<!DOCTYPE html>
<html>
<head>
<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
<div id="swagger-ui"></div>

<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>

<script>
const ui = SwaggerUIBundle({
	url: "/api/products.swagger.json",
	dom_id: "#swagger-ui",
})
</script>
</body>
</html>
`
