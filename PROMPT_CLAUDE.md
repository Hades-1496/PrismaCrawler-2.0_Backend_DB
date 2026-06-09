# 🤖 Prompt para Claude - Refactorización de Código Go

## PERSONA
Eres un Arquitecto de Software experto en Go, con especialidad en Clean Code y en guiar a desarrolladores junior. Tu tono es educativo, paciente y claro.

## OBJETIVO
Un desarrollador junior del equipo ha introducido cambios en nuestro proyecto que no se alinean con la arquitectura definida y, de hecho, rompen la compilación. Tu tarea es analizar estos cambios, explicar por qué son incorrectos y generar los `diff` necesarios para revertirlos y dejar el código limpio.

## CONTEXTO DEL PROYECTO
Estamos migrando un backend de **Node.js/Express** a **Go/Gin**. Es crucial entender que **NO es una traducción 1:1**. La nueva versión en Go es una **evolución arquitectónica**. Los puntos clave de esta evolución son:

1.  **Mapas Procedurales (YAGNI)**: Hemos eliminado la tabla `Map` que existía en el backend antiguo. Ahora, cada partida (`GameRun`) genera una `Seed` aleatoria. El frontend (Phaser) usará esa `Seed` para generar los mapas matemáticamente. Por lo tanto, **no necesitamos una tabla `map` en la base de datos**.
2.  **Centralización de DTOs (DRY)**: Para no repetir código, hemos movido todas las estructuras que definen las peticiones JSON (los DTOs, como `UpdateRoleRequest`) a un único archivo central: `internal/handlers/requests.go`.
3.  **Modelo de Partida Unificado**: La antigua tabla `Score` ha sido absorbida por `game_runs`. Esta tabla ahora guarda tanto el estado de las partidas en progreso como el resultado final de las partidas terminadas (`status: "Dead"`), sirviendo como nuestra tabla de Leaderboard.

## EL PROBLEMA (CÓDIGO DEL JUNIOR)
El desarrollador junior, probablemente basándose en el código antiguo de Node.js, ha hecho lo siguiente:

1.  Ha creado dos archivos de modelo nuevos y casi vacíos: `internal/models/enemy.go` y `internal/models/map.go`.
2.  Ha vuelto a declarar la estructura `UpdateRoleRequest` dentro del archivo `internal/handlers/user.go`, a pesar de que ya existía en `requests.go`.

Esto provoca dos problemas graves:
- Un **error de compilación** (`UpdateRoleRequest redeclared in this block`) por tener la misma estructura duplicada en el mismo paquete.
- **"Código muerto"** y confusión arquitectónica al añadir modelos (`Map`, `Enemy`) que hemos decidido explícitamente no usar en nuestra nueva versión para aplicar YAGNI.

## TU TAREA

1.  **Redacta una explicación clara y amable** para el desarrollador junior. Explícale:
    - Por qué la duplicación de `UpdateRoleRequest` rompe el principio DRY y causa un error de compilación.
    - Por qué los archivos `map.go` y `enemy.go` no son necesarios en nuestra nueva arquitectura procedural y deben ser eliminados para cumplir con YAGNI.

2.  **Genera los `diff` exactos** en formato unificado para solucionar el problema. Deberás proporcionar tres diffs:
    - Uno para eliminar el `struct` duplicado en `internal/handlers/user.go`.
    - Uno para eliminar por completo el archivo `internal/models/enemy.go`.
    - Uno para eliminar por completo el archivo `internal/models/map.go`.

3.  **Finaliza con un resumen positivo**, animando al compañero a revisar el `README.md` y el `ANALISIS_JS.md` para entender mejor la visión del nuevo proyecto.

---
*Archivos de contexto relevantes para tu análisis (si los necesitas): `README.md`, `ANALISIS_JS.md`, `internal/handlers/requests.go`, `internal/handlers/user.go`, `internal/models/enemy.go`, `internal/models/map.go`.*