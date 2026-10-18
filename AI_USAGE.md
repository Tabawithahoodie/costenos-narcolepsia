# Uso de IA

## Herramientas y para qué
- Durante la elaboración del proyecto utilicé ChatGPT y Claude como apoyo para entender los conceptos de Go y organizar el desarrollo del sistema del Hospital de los Costeños con Narcolepsia.

Claude me ayudo a resolver errores y plantear codigo, mientras que ChatGPT lo utilicé para entendermas a fondo como funcionaba cada cosa y a la vez resolver dudas.

La herramienta me ayudó a comprender cómo funcionan las estructuras (struct), los métodos, las interfaces, la composición y el manejo de errores. También la utilicé como guía para instalar Go, configurar Visual Studio Code y resolver dudas durante la programación.

El código se fue revisando y probando para entender qué hace cada parte y cómo se relaciona con los requisitos del ejercicio. 
El README. Lo corrí con `gofmt`, `go vet` y `go test` para comprobarlo.

## Prompts decisivos
**1.** Le pedi una guia basica de go para entender el trabajo con el siguiente prompt "Dame una guia basica de go para entender el proyecto"
**2.** Mande el archivo del trabajo y le pedi que me lo desglosara para entender todo como si fuera un diagrama de flujo.
**3.** Tuve problemas para descargar la extension oficial desde Go y habia instalado una que no era, asi que tuve que mandarle ss de lo que habia descargado y buscar el link oficial de implementacion de la extension a Visual Studio Code.

## Casos donde el código generado estaba mal o no era idiomático
1. **Errores sin usar y mal formateados:** la primera versión de `errors.go` tenía
   errores centinela que nunca se usaban (`ErrNilStaffOrRoom`, `ErrNotAsleepToWake`)
   y `gofmt -l` marcó `errors.go` y `patient.go`. Se limpió el archivo, se dejó solo
   lo que se usa y se aplicó `gofmt -w`.
2. **Escenario que se salía del enunciado:** en la primera versión de `main.go`, la
   consulta 5.3 hacía dormir a un quinto paciente, algo que el escenario obligatorio
   no pide. Se cambió para que 5.3 imprima los resultados de las asignaciones del
   escenario (incluida la reasignación tras la liberación).

## Qué aprendí
Aprendí a utilizar Go y a organizar un proyecto en diferentes archivos. También entendí cómo funcionan las estructuras, los métodos, las interfaces y el manejo de errores. Además, aprendí la importancia de realizar pruebas y comprender el código para poder explicar cómo funciona.
Durante este proyecto intente ser muy honesto y buscar entender con pruebas y errores y no simplemente dejar que la IA se encargue de todo
