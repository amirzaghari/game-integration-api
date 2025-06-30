package http

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

const swaggerHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Game Integration API - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui.css" />
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *, *:before, *:after {
            box-sizing: inherit;
        }
        body {
            margin:0;
            background: #fafafa;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: '/swagger/doc.json',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                onComplete: function() {
                    // Pre-fill the Bearer token field
                    setTimeout(function() {
                        const authInput = document.querySelector('input[placeholder*="Authorization"]');
                        if (authInput && !authInput.value) {
                            authInput.value = 'Bearer ';
                            authInput.focus();
                            // Position cursor after "Bearer "
                            authInput.setSelectionRange(7, 7);
                        }
                    }, 1000);
                }
            });
        };
    </script>
</body>
</html>
`

// SwaggerHandler serves the custom Swagger UI with pre-filled Bearer token
func SwaggerHandler() gin.HandlerFunc {
	tmpl := template.Must(template.New("swagger").Parse(swaggerHTML))

	return func(c *gin.Context) {
		if c.Request.URL.Path == "/swagger/index.html" {
			c.Header("Content-Type", "text/html")
			tmpl.Execute(c.Writer, nil)
			return
		}

		// For other swagger paths, serve the default swagger files
		c.Next()
	}
}
