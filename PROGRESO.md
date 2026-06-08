# 📅 Diario de Progreso - 10 Días

Este documento sirve para anotar el estado del proyecto y asegurar que llegamos a la entrega a tiempo.

## 🔄 Tareas de Migración (De JS a Go)
- [ ] Migrar la lógica de Autenticación (JWT/Supabase) de JS a Go.
- [ ] Migrar las rutas (endpoints) de usuarios e inventarios usando Gin.
- [ ] Traducir los servicios (lógica de negocio) de JS a Go.
- [ ] Adaptar los modelos de GORM para que coincidan con las tablas creadas originalmente por Prisma/JS.

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
- [ ] Implementar registro y login de usuarios (Endpoints y lógica).
- [ ] Generación y validación de tokens JWT.
- [ ] *Notas de bloqueo/aprendizaje:* 

## Día 7-8: Core del Roguelike (Personajes y Runs)
- [ ] Endpoints de creación y selección de personajes.
- [ ] Endpoints para Iniciar una Partida (Run) y guardar el estado actual del calabozo.
- [ ] *Notas de bloqueo/aprendizaje:* 

## Día 9-10: Pulido, Testeo y Entrega
- [ ] Documentación en Postman para probar la API rápidamente.
- [ ] Limpieza de código y manejo de errores HTTP (400, 401, 404, 500).
- [ ] Preparar presentación/defensa de la introducción a Go.