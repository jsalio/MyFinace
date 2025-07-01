.PHONY: test test-coverage run

test:
	@echo "Ejecutando tests..."
	cd test && go test -v ./... -coverpkg=./...

test-coverage:
	@echo "Ejecutando tests con cobertura..."
	mkdir -p test/coverage
	cd test && go test -coverprofile=coverage/coverage.out -coverpkg=./...,Financial/Core/... ./...
	@echo "\nResumen de cobertura:"
	@cd test && go tool cover -func=coverage/coverage.out
	@echo "\nGenerando reporte HTML..."
	@cd test && go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@echo "Reporte de cobertura generado en: coverage/coverage.html"
	@xdg-open test/coverage/coverage.html 2>/dev/null || open test/coverage/coverage.html 2>/dev/null || echo "No se pudo abrir el navegador automáticamente. Abre test/coverage/coverage.html manualmente."

run:
	@echo "Regenerando documentación Swagger..."
	swag init -g intefaces/server.go
	@echo "Iniciando la aplicación..."
	go run main.go