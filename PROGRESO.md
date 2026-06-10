# 📅 Diario de Progreso - 10 Días

Este documento sirve para anotar el estado del proyecto y asegurar que llegamos a la entrega a tiempo.

## 🔄 Tareas de Migración (De JS a Go)
- [x] Migrar la lógica de Autenticación (JWT/Supabase) de JS a Go.
- [x] Migrar las rutas (endpoints) de usuarios e inventarios usando Gin. *(Perfil migrado)*
- [x] Traducir los servicios (lógica de negocio) de JS a Go.
- [x] Adaptar los modelos de GORM para que coincidan con las tablas creadas originalmente por Prisma/JS.

## Día 1-2: Setup y Fundamentos
- [x] Instalación de Go y configuración del proyecto (`go mod init`).
- [x] Creación del proyecto en Supabase y obtención de credenciales de PostgreSQL.
- [x] Conexión a la base de datos desde Go usando GORM o Ent.
- [ ] *Notas de bloqueo/aprendizaje:* 

## Día 3-4: Entidades e Infraestructura
- [x] Creación de los modelos (Users, Characters, Items, Game_Runs, Run_Inventory).
- [x] Ejecutar las migraciones automáticas hacia Supabase.
- [x] Levantar el servidor HTTP básico con el framework Gin.
- [ ] *Notas de bloqueo/aprendizaje:* 

## Día 5-6: Autenticación y Usuarios
- [x] Implementar registro y login de usuarios (Endpoints y lógica).
- [x] Generación y validación de tokens JWT.
- [ ] *Notas de bloqueo/aprendizaje:* 

## Día 7-8: Core del Roguelike (Personajes y Runs)
- [x] Endpoints de creación y selección de personajes.
- [x] Endpoints para Iniciar una Partida (Run) y guardar el estado actual del calabozo.
- [x] Guardado de puntuación al morir. (Integrado dinámicamente en SaveRun)
- [x] Mostrar ranking de puntuación con otros jugadores.
- [ ] *Notas de bloqueo/aprendizaje:* 

## Día 9-10: Orquestación y Microservicios (IA)
- [x] Configurar las variables de entorno `AI_SERVICE_URL` y `AI_INTERNAL_TOKEN` en `main.go`.
- [x] **Base de Conocimiento (RAG)**: Crear modelo `KnowledgeChunk` y endpoints CRUD protegidos para administradores.
- [x] **Cliente IA (`pkg/aiclient`)**: Crear cliente HTTP interno asíncrono que inyecte `X-Internal-Token`.
- [x] **Proxy FAQ**: Implementar `POST /api/faq` para rebotar las preguntas al microservicio.
- [x] **Triggers Automáticos (n8n)**:
  - [x] Enviar evento `user_registered` al finalizar `POST /auth/register`.
  - [x] Enviar evento `game_run_ended` al morir en `PUT /api/runs/save`.
- [x] **Discord Changelogs**: Implementar proxy `POST /api/admin/discord/changelog`.
- [ ] *Notas de bloqueo/aprendizaje:*

## Día 11-12: Pulido, Testeo y Entrega
- [ ] Limpieza de código y manejo de errores HTTP (400, 401, 404, 500).
- [ ] Refactorización mediante principios de Clean Code:
  - [x] KISS (Rutas centralizadas en un solo archivo, Repository Pattern evitado)
  - [x] DRY (Helper de Errores y Helper de GetUserID)
  - [ ] YAGNI (You Ain't Gonna Need It)
  - [x] SoC (Separation of Concerns: main.go vs routes.go)
  - [x] LoD (Law of Demeter)
  - [x] SRP (Single Responsibility Principle: lógica de daño en modelo)
  - [ ] OCP (Open/Closed Principle) // Aún no aplicado
  - [ ] LSP (Liskov Substitution Principle) // No aplica directamente
  - [x] ISP (Interface Segregation Principle) // Interfaces pequeñas y específicas
  - [x] DIP (Dependency Inversion Principle) // Inyección de dependencias
- [ ] Implementar medidas de seguridad básicas (Rate limiting, delays artificiales en Auth).
- [ ] Repaso para el correcto funcionamiento en un servidor linux / Render (para el uso correcto del Deploy).
- [ ] Implementación de tests mediante metodología de TDD (Crear tests aislados, explotación de endpoints, Test completo, etc.) 
- [ ] Documentación en Postman para probar la API rápidamente.
- [ ] Preparar presentación/defensa de la introducción a Go.