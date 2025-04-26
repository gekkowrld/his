// tools/run.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	buildDeps()

	fmt.Println("Starting API...")
	api := exec.CommandContext(ctx, "./his")
	api.Stdout = os.Stdout
	api.Stderr = os.Stderr
	if err := api.Start(); err != nil {
		log.Fatalf("Failed to start API: %v\n", err)
	}

	fmt.Println("Starting client...")
	client := exec.CommandContext(ctx, "./client/client")
	client.Stdout = os.Stdout
	client.Stderr = os.Stderr
	if err := client.Start(); err != nil {
		log.Printf("Failed to start client: %v\n", err)
		_ = api.Process.Kill()
		os.Exit(1)
	}

	<-ctx.Done()
	fmt.Println("\nReceived signal, shutting down...")

	_ = api.Process.Signal(syscall.SIGTERM)
	_ = client.Process.Signal(syscall.SIGTERM)

	_ = api.Wait()
	_ = client.Wait()

	fmt.Println("Both API and client exited.")
}

func buildDeps() {
	fmt.Println("Building API...")
	his := exec.Command("go", "build", "-tags", "fts5", "-ldflags=-w -s", "-o", "his", ".")
	his.Stdout = os.Stdout
	his.Stderr = os.Stderr
	if err := his.Run(); err != nil {
		log.Fatalf("API build failed: %v\n", err)
	}

	clientDir := "client"
	origDir, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(origDir)
	}()

	fmt.Println("Building client...")
	err := os.Chdir(clientDir)
	if err != nil {
		log.Fatalf("Failed to change to client dir: %v\n", err)
	}

	client := exec.Command("go", "build", "-ldflags=-w -s", "-o", filepath.Join("..", "client", "client"), ".")
	client.Stdout = os.Stdout
	client.Stderr = os.Stderr
	if err := client.Run(); err != nil {
		log.Fatalf("Client build failed: %v\n", err)
	}
}
