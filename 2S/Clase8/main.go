package main

import (
	"fmt"

	"gocv.io/x/gocv"
)

func main() {
	// Abrir cámara predeterminada de la computadora
	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		fmt.Printf("Error al abrir la cámara: %v\n", err)
		return
	}
	defer webcam.Close()

	// Crear ventana para mostrar video
	window := gocv.NewWindow("Detección de cambios de diapositiva")
	defer window.Close()

	frame := gocv.NewMat()
	defer frame.Close()
	prevFrame := gocv.NewMat()
	defer prevFrame.Close()
	gray := gocv.NewMat()
	defer gray.Close()

	threshold := 20.0 // Umbral para detectar cambio
	fmt.Println("Presiona ESC para salir")

	for {
		if ok := webcam.Read(&frame); !ok || frame.Empty() {
			continue
		}

		// Convertir a escala de grises
		gocv.CvtColor(frame, &gray, gocv.ColorBGRToGray)

		if !prevFrame.Empty() {
			diff := gocv.NewMat()
			gocv.AbsDiff(gray, prevFrame, &diff)

			mean := gocv.Mean(diff)
			if mean[0] > threshold {
				fmt.Println("¡Cambio de diapositiva detectado!")
			}

			diff.Close()
		}

		// Guardar frame actual como previo
		gray.CopyTo(&prevFrame)

		// Mostrar video en ventana
		window.IMShow(frame)
		if window.WaitKey(1) == 27 { // ESC para salir
			break
		}
	}
}
