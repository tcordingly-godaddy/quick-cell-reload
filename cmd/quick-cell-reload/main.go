package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/tcordingly-godaddy/quick-cell-reload/pkg/jobmeta"
)

const (
	defaultNamespace = "sites"
	defaultPattern   = "app-"
)

func main() {
	var (
		jobID      = flag.String("job", "", "Job ID to update")
		namespace  = flag.String("namespace", defaultNamespace, "Nomad namespace")
		listMeta   = flag.Bool("list", false, "List current meta tags for the job")
		jobPattern = flag.String("pattern", "", "Job pattern for updating multiple jobs")
		timeout    = flag.Duration("timeout", 30*time.Second, "Timeout for operations")
	)
	flag.Parse()

	if err := validateFlags(jobID, *listMeta, jobPattern); err != nil {
		log.Fatal(err)
	}

	updater, err := jobmeta.NewUpdater()
	if err != nil {
		log.Fatalf("Failed to create job meta updater: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	if err := runCommand(ctx, updater, jobID, namespace, *listMeta, jobPattern); err != nil {
		log.Fatalf("Command failed: %v", err)
	}

	log.Println("Job meta update completed successfully")
}

func validateFlags(jobID *string, listMeta bool, jobPattern *string) error {
	if listMeta && *jobID == "" {
		return fmt.Errorf("-list requires -job to be specified")
	}

	return nil
}

func runCommand(ctx context.Context, updater *jobmeta.Updater, jobID, namespace *string, listMeta bool, jobPattern *string) error {
	if listMeta {
		return updater.ListJobMeta(*jobID, *namespace)
	}

	// Generate a new reload hash
	reloadHash, err := generateReloadHash()
	if err != nil {
		return fmt.Errorf("failed to generate reload hash: %w", err)
	}

	metaUpdates := map[string]string{
		"reload-hash": reloadHash,
	}

	// If no job ID is provided, use the pattern (defaulting to "app-" if no pattern specified)
	if *jobID == "" {
		pattern := *jobPattern
		if pattern == "" {
			pattern = defaultPattern
		}
		return updater.UpdateMultipleJobs(ctx, *namespace, pattern, metaUpdates)
	}

	return updater.UpdateJobMeta(ctx, *jobID, *namespace, metaUpdates)
}

func generateReloadHash() (string, error) {
	// Generate 16 random bytes
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Convert to hex string
	return hex.EncodeToString(bytes), nil
}
