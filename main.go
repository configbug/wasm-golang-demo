package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"syscall/js"
)

func httpRequest() (string, error) {
	log.Println("(Log Golang): Acá inicia con la solicitud HTTP")

	httpClient := &http.Client{}
	url := "https://reqres.in/api/users"
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", err
	}

	request.Header.Add("Content-Type", "application/json")
	response, err := httpClient.Do(request)
	log.Println("(Log Golang): Acá realizó una solicitud HTTP")
	log.Printf("(Log Golang): Acá verificamos el código de respuesta de la solicitud HTTP: %d", response.StatusCode)

	if err != nil {
		log.Fatalf("Error al invocar: %v", err)
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	log.Println("(Log Golang): Acá leyó la respuesta de la solicitud HTTP")

	if err != nil {
		log.Fatalf("Error al leer la respuesta: %v", err)
	}

	var obj map[string]interface{}
	err = json.Unmarshal(responseBody, &obj)

	if err != nil {
		log.Fatalf("Error al serializar a JSON: %v", err)
	}

	prettyJSON, err := json.MarshalIndent(obj, "", "  ")

	if err != nil {
		log.Fatalf("Error al identar el JSON: %v", err)
	}

	stringResponse := string(prettyJSON)
	log.Printf("(Log Golang): Acá termino de obtener la solictud en formato JSON\n\n'%s'", stringResponse)
	return stringResponse, nil
}

func main() {
	btnRequest := js.Global().Get("document").Call("querySelector", "#btnRequest")
	dom := js.Global().Get("document").Call("querySelector", "#response")

	channel := make(chan struct{})

	btnRequest.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		dom.Set("innerHTML", "Cargando...")
		log.Println("(Log Golang): Se realiza la petición invocada desde el browser")

		go func() {
			response, err := httpRequest()

			if err != nil {
				dom.Set("innerHTML", "Error al realizar la petición")
				log.Printf("Error al realizar la petición: '%v'", err)
			}

			dom.Set("innerHTML", response)
		}()
		return nil
	}))

	<-channel
}
