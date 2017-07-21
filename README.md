[![Go Report Card](https://goreportcard.com/badge/github.com/juusechec/goresource-proxy)](https://goreportcard.com/report/github.com/juusechec/goresource-proxy)

# goresource-proxy
A file proxy like https://github.com/Esri/resource-proxy Go Proxy Files

## Algunos ejemplos:
Se realiza la petición a http://localhost:12345/?form
con la URL https://www.google.com.co/ que resulta en
http://localhost:12345/?url=https%3A%2F%2Fwww.google.com.co%2F
la cuál funciona correctamente.

![screenshot](./images/Screenshot_from_2017-07-20_18-47-42.png)
![screenshot](./images/Screenshot_from_2017-07-20_18-48-22.png)


Se realiza la petición a http://localhost:12345/?form
con la URL https://www.google.com.co/ que resulta en
http://localhost:12345/?url=https%3A%2F%2Fwww.google.com.co%2F
la cuál es rechazada debido a que no está en el archivo [whitelist.srt](./whitelist.srt)

![screenshot](./images/Screenshot_from_2017-07-20_18-48-54.png)
![screenshot](./images/Screenshot_from_2017-07-20_18-49-07.png)

Puede cambiar los HEADERS de la petición con el parámetro ***headers*** de la
petición GET, separando cada cabecera con el string "\r\n" (sin comillas) y
separando clave y valor con ":" (dos puntos), un ejemplo puede ser
"Cookie:newcookie\r\nGoogle:GENIO". La construcción en lenguaje Javascript sería.
URL sería.
```js
var headers = 'Cookie:newcookie\\r\\nGoogle:GENIO';
var url = 'https://www.google.com.co';
console.log(headers, url);
var completeURL = 'http://localhost:12345/proxy?headers=' + escape(headers) + '&url=' + escape(url);
console.log(completeURL);
```
