package jobmeta

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/nomad/api"
)

// Updater handles updating meta tags in Nomad jobs
type Updater struct {
	client *api.Client
}

// NewUpdater creates a new JobMetaUpdater instance
func NewUpdater() (*Updater, error) {
	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Nomad client: %w", err)
	}

	return &Updater{
		client: client,
	}, nil
}

// UpdateJobMeta updates the meta tags of a specific job and submits the update
func (u *Updater) UpdateJobMeta(ctx context.Context, jobID, namespace string, metaUpdates map[string]string) error {
	// Get the current job
	job, _, err := u.client.Jobs().Info(jobID, &api.QueryOptions{
		Namespace: namespace,
	})
	if err != nil {
		return fmt.Errorf("failed to get job %s: %w", jobID, err)
	}

	// Initialize meta map if it doesn't exist
	if job.Meta == nil {
		job.Meta = make(map[string]string)
	}

	// Update the meta tags
	for key, value := range metaUpdates {
		job.Meta[key] = value
		log.Printf("Setting meta tag %s = %s", key, value)
	}

	// Submit the job update
	writeOpts := &api.WriteOptions{
		Namespace: namespace,
	}

	response, _, err := u.client.Jobs().Register(job, writeOpts)
	if err != nil {
		return fmt.Errorf("failed to register job update: %w", err)
	}

	log.Printf("Job update submitted successfully. Evaluation ID: %s", response.EvalID)
	return nil
}

// ListJobMeta displays the current meta tags for a job
func (u *Updater) ListJobMeta(jobID, namespace string) error {
	job, _, err := u.client.Jobs().Info(jobID, &api.QueryOptions{
		Namespace: namespace,
	})
	if err != nil {
		return fmt.Errorf("failed to get job %s: %w", jobID, err)
	}

	if job.Meta == nil || len(job.Meta) == 0 {
		fmt.Printf("No meta tags found for job %s\n", jobID)
		return nil
	}

	fmt.Printf("Meta tags for job %s:\n", jobID)
	for key, value := range job.Meta {
		fmt.Printf("  %s = %s\n", key, value)
	}

	return nil
}

// UpdateMultipleJobs updates meta tags for multiple jobs based on a pattern or list
func (u *Updater) UpdateMultipleJobs(ctx context.Context, namespace string, jobPattern string, metaUpdates map[string]string) error {
	// List jobs in the namespace
	jobs, _, err := u.client.Jobs().List(&api.QueryOptions{
		Namespace: namespace,
	})
	if err != nil {
		return fmt.Errorf("failed to list jobs: %w", err)
	}

	updatedCount := 0
	for _, jobStub := range jobs {
		// Check if job matches pattern (simple prefix matching for now)
		if strings.HasPrefix(jobStub.ID, jobPattern) {
			log.Printf("Updating job: %s", jobStub.ID)
			if err := u.UpdateJobMeta(ctx, jobStub.ID, namespace, metaUpdates); err != nil {
				log.Printf("Failed to update job %s: %v", jobStub.ID, err)
				continue
			}
			updatedCount++
		}
	}

	log.Printf("Updated %d jobs", updatedCount)
	return nil
}

// WaitForEvaluation waits for a Nomad evaluation to complete
func (u *Updater) WaitForEvaluation(evalID string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for evaluation %s", evalID)
		case <-ticker.C:
			eval, _, err := u.client.Evaluations().Info(evalID, nil)
			if err != nil {
				log.Printf("Error checking evaluation status: %v", err)
				continue
			}

			switch eval.Status {
			case "complete":
				log.Printf("Evaluation %s completed successfully", evalID)
				return nil
			case "failed":
				return fmt.Errorf("evaluation %s failed", evalID)
			default:
				log.Printf("Evaluation %s status: %s", evalID, eval.Status)
			}
		}
	}
}
