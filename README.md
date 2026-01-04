# Jorge Díaz Alcojor

# SistemasDistribuidos-p4 – Taller distribuido y concurrente en Go

En esta práctica se modela un **taller de reparación de coches distribuido**, cuyo comportamiento se adapta dinámicamente en función del **estado recibido desde un servidor TCP**.

El sistema atiende coches de **tres categorías**:

* **A – Mecánica**
* **B – Eléctrica**
* **C – Carrocería**

Cada coche atraviesa **cuatro fases secuenciales** del taller (plaza, mecánico, limpieza y entrega), respetando la capacidad máxima de cada fase, las restricciones por categoría y las prioridades dinámicas indicadas por el estado remoto.

La práctica parte estrictamente del **código proporcionado por el profesor** (`servidor`, `mutua` y `taller`) y amplía únicamente el archivo permitido, organizando la lógica interna del taller en varios ficheros Go para mejorar la claridad y mantenibilidad del sistema.

---

## Estructura del proyecto

El proyecto se divide en **tres ejecutables independientes**:

### `servidor`

Servidor TCP que publica el estado del taller y acepta conexiones de clientes.

### `mutua`

Cliente que genera y envía códigos (0..9) al servidor, simulando cambios en el estado del sistema.

### `taller`

Cliente que se conecta al servidor y ejecuta la simulación concurrente del taller.
Incluye además los **tests automáticos** de la práctica.

---

## Funcionamiento general

El taller se conecta al servidor TCP y recibe **códigos numéricos entre 0 y 9** que representan el estado del sistema. Entre los estados posibles se encuentran:

* Taller **inactivo** o **cerrado**.
* Restricción de atención a una única categoría (`SOLO A`, `SOLO B` o `SOLO C`).
* Priorización dinámica de una categoría concreta (`PRIORIDAD A`, `PRIORIDAD B` o `PRIORIDAD C`).

El estado afecta tanto a **qué coches pueden avanzar por las fases** como al **orden de atención en las colas** del taller.

---

## Fases del taller

Cada coche atraviesa las siguientes fases:

* **Fase 0 – Plaza**
* **Fase 1 – Mecánico**
* **Fase 2 – Limpieza**
* **Fase 3 – Entrega**

Cada fase se implementa mediante **goroutines** que actúan como *workers* y utilizan **channels como semáforos** para representar los recursos físicos.
Los coches se gestionan mediante **colas con prioridad** y, antes de iniciar cualquier trabajo, se vuelve a comprobar el estado del taller para asegurar que se respetan las restricciones actuales.

---

## Formato de logs

Las trazas de ejecución siguen el formato exigido en el enunciado:

```
Tiempo {t} Coche {id} Incidencia {tipo} Fase {fase} Estado {Entra|Sale}
```

La impresión de logs se centraliza en una única goroutine para evitar interferencias entre goroutines concurrentes.

---

## Archivos principales del taller

### `runtime.go`

Inicializa toda la infraestructura concurrente del taller: canales de entrada de mensajes TCP, parser de mensajes, controlador del estado del taller, logger central y arranque de la simulación.
Es el punto de unión entre el código proporcionado por el profesor y la lógica interna del taller.

### `controller.go`

Implementa el **gestor del estado del taller**. Mantiene el estado actual del sistema y procesa los códigos recibidos del servidor.
Permite que las fases consulten el estado de forma segura sin necesidad de utilizar mutexes explícitos.

### `queues.go`

Implementa la estructura **`PhaseQueue`**, que representa las colas de cada fase con soporte de **prioridad por categoría** y prioridad dinámica según el estado del taller.
Cada cola se gestiona internamente mediante una goroutine.

### `phases.go`

Contiene la implementación de las **cuatro fases del taller**: plaza, mecánico, limpieza y entrega.
Cada fase respeta el estado del taller, bloquea en el recurso correspondiente, simula el tiempo de trabajo, genera logs de entrada y salida y pasa el coche a la siguiente fase.

### `sim.go`

Contiene la lógica de simulación de alto nivel.
Genera los coches por categoría, crea las colas y recursos, lanza los workers por fase y arranca el pipeline completo del taller.

### `logger.go`

Goroutine dedicada a la impresión de logs con formato consistente, evitando *interleaving* entre goroutines.

### `sim_test.go`

Archivo de **tests automáticos** usando el paquete `testing` de Go.

Se realizan las **seis comparativas exigidas** en el enunciado:

#### Distribuciones de coches

* A=10, B=10, C=10
* A=20, B=5, C=5
* A=5, B=5, C=20

#### Configuraciones de recursos

* NumPlazas=6, NumMecánicos=3
* NumPlazas=4, NumMecánicos=4

En total se ejecutan **seis tests**.
Cada test mide la **duración total** de la simulación y el **throughput** en coches por segundo, verificando que todos los coches alcanzan la fase de entrega.

Para evitar tiempos de ejecución excesivos, se utiliza un **factor de escala temporal**, que reduce proporcionalmente las esperas manteniendo las relaciones entre fases y categorías.

---

## Cómo ejecutar la práctica

Desde la raíz del proyecto:

### Ejecutar el servidor

```
go run ./servidor
```

### Ejecutar el taller

```
go run ./taller
```

### Ejecutar la mutua

```
go run ./mutua
```

El taller reaccionará en tiempo real a los estados enviados por la mutua a través del servidor.

---

## Cómo ejecutar los tests

```
go test ./taller -v
```

Se mostrarán por pantalla las métricas de **duración** y **throughput** para cada uno de los seis escenarios.

---

## Conclusión

La práctica implementa un **sistema distribuido y concurrente completo** que combina:

* comunicación TCP,
* concurrencia con goroutines y channels,
* control centralizado del estado,
* colas con prioridad,
* y pruebas comparativas automatizadas,

respetando estrictamente el código base proporcionado y utilizando únicamente mecanismos vistos en la asignatura.

---