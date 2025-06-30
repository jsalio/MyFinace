.PHONY: test test-coverage run

test:
	@echo "Ejecutando tests..."
	go test -v ./...

test-coverage:
	@echo "Ejecutando tests con cobertura..."
	go test -coverprofile=coverage.out ./...
	@echo "\nResumen de cobertura:"
	@go tool cover -func=coverage.out
	@echo "\nGenerando reporte HTML..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Reporte de cobertura generado en: coverage.html"
	@xdg-open coverage.html 2>/dev/null || open coverage.html 2>/dev/null || echo "No se pudo abrir el navegador automáticamente. Abre coverage.html manualmente."

run:
	@echo "Regenerando documentación Swagger..."
	swag init -g intefaces/server.go
	@echo "Iniciando la aplicación..."
	go run main.go
