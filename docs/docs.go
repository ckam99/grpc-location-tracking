package docs

import (
	"fmt"
	"net/http"
)

var SwaggerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	// Write the HTML content to the response
	fmt.Fprintln(w, swaggerUI)
})

var swaggerUI = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Swagger UI</title>
    <!-- Swagger UI CSS CDN -->
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.14.1/swagger-ui.min.css">
    <style>
        body { margin: 0; }
        #swagger-ui { height: 100vh; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <!-- Swagger UI JS CDN -->
    <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.14.1/swagger-ui-bundle.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.14.1/swagger-ui-standalone-preset.min.js"></script>
    <script>
        const ui = SwaggerUIBundle({
            url: "/swagger/apidocs.swagger.json", // Replace with your Swagger API URL
            dom_id: '#swagger-ui',
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIStandalonePreset
            ],
            layout: "StandaloneLayout"
        });
    </script>
</body>
</html>
`
