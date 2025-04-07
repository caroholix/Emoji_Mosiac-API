package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// Function to process uploaded images
func processImage(w http.ResponseWriter, r *http.Request) {
	// Ensure required folders exist
	os.MkdirAll("src/main/input_images", os.ModePerm)
	os.MkdirAll("output/image_output", os.ModePerm)

	// Get the uploaded file
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Image file missing", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Generate a unique filename using timestamp + random number
	rand.Seed(time.Now().UnixNano())
	uniqueID := strconv.Itoa(int(time.Now().UnixNano())) + strconv.Itoa(rand.Intn(10000))
	inputPath := "src/main/input_images/" + uniqueID + ".png"
	outputPath := "output/image_output/" + uniqueID + ".png"

	// Save the uploaded file
	outFile, err := os.Create(inputPath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()
	io.Copy(outFile, file)

	cmd := exec.Command("./main", "video", "4", "input_images/")
cmd.Dir = "src/main" // This sets the working directory


	output, err := cmd.CombinedOutput()

	// If there's an error, return it
	if err != nil {
		fmt.Println("❌ Spritefire Error:", err)
		fmt.Println("📄 Spritefire Output:", string(output))
		http.Error(w, "Error processing image:\n"+string(output), http.StatusInternalServerError)
		return
	}

	// Log success
	fmt.Println("✅ Image processed successfully:", inputPath)

	// Send the processed image back
	http.ServeFile(w, r, outputPath)

	// Clean up: Delete files after 10 seconds (optional)
	go func() {
		time.Sleep(10 * time.Second)
		os.Remove(inputPath)
		os.Remove(outputPath)
	}()
}

func main() {
	http.HandleFunc("/process-image", processImage) // API route
	port := "8080"
	fmt.Println("🚀 Server running on 0.0.0.0:" + port + "...")
	err := http.ListenAndServe("0.0.0.0:"+port, nil) // Allow external access
	if err != nil {
		fmt.Println("❌ Error starting server:", err)
	}
}

