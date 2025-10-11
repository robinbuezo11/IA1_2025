package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/otiai10/gosseract/v2"
)

// Función para extraer texto con gosseract
func extractTextFromImage(imagePath string) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(imagePath)
	text, err := client.Text()
	if err != nil {
		return "", err
	}
	return text, nil
}

// Enviar texto al modelo local de Ollama
func expandWithAI(text string) (string, error) {
	payload := map[string]string{"model": "phi3", "prompt": "Amplía el siguiente texto:\n" + text}
	data, _ := json.Marshal(payload)

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result bytes.Buffer
	decoder := json.NewDecoder(resp.Body)
	for decoder.More() {
		var msg map[string]interface{}
		if err := decoder.Decode(&msg); err != nil {
			break
		}
		if val, ok := msg["response"]; ok {
			result.WriteString(val.(string))
		}
	}

	return result.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Use: go run main.go <imagen>")
		return
	}

	imagePath := os.Args[1]

	fmt.Println("Extrayendo texto de:", imagePath)
	text, err := extractTextFromImage(imagePath)
	if err != nil {
		fmt.Println("Error al extraer texto:", err)
		return
	}

	fmt.Println("\nTexto detectado:")
	fmt.Println(text)

	fmt.Println("\nEnviando texto al modelo local...")
	expanded, err := expandWithAI(text)
	if err != nil {
		fmt.Println("Error al ampliar texto:", err)
		return
	}

	fmt.Println("\nTexto ampliado por IA:")
	fmt.Println(expanded)
}
