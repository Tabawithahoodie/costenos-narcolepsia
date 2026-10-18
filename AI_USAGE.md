# Uso de IA

## Herramientas y para qué
- **Claude (Anthropic):** diseño del modelo, escribir el código en Go, los tests y
  el README. Lo corrí con `gofmt`, `go vet` y `go test` para comprobarlo.

## Prompts decisivos
_(Pega aquí 2 o 3 prompts reales que usaste. Por ejemplo: el enunciado completo
del trabajo con la instrucción de resolverlo parte por parte, y los prompts donde
pediste cambios o explicaciones.)_

## Casos donde el código generado estaba mal o no era idiomático
1. **Errores sin usar y mal formateados:** la primera versión de `errors.go` tenía
   errores centinela que nunca se usaban (`ErrNilStaffOrRoom`, `ErrNotAsleepToWake`)
   y `gofmt -l` marcó `errors.go` y `patient.go`. Se limpió el archivo, se dejó solo
   lo que se usa y se aplicó `gofmt -w`.
2. **Escenario que se salía del enunciado:** en la primera versión de `main.go`, la
   consulta 5.3 hacía dormir a un quinto paciente, algo que el escenario obligatorio
   no pide. Se cambió para que 5.3 imprima los resultados de las asignaciones del
   escenario (incluida la reasignación tras la liberación).

_(Si en tu proceso pasaron otros casos, agrégalos. Esta sección debe coincidir con
lo que cuentes en la sustentación.)_

## Qué aprendí
_(Escríbelo con tus palabras: embedding vs. herencia, interfaces implícitas,
campos en minúscula y paquetes, receptor por puntero vs. por valor, errores como
valores con `errors.Is` y `%w`.)_
