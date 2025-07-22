package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github"
)

func main() {
	// Create a GitHub client
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	// Set up webhook handler
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		// Parse the webhook payload
		hook, err := client.Webhooks.Parse(r, func(webhook scm.Webhook) (string, error) {
			// Return your webhook secret here
			return "your-webhook-secret", nil
		})
		if err != nil {
			log.Printf("Error parsing webhook: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Handle pull request review events
		if reviewHook, ok := hook.(*scm.PullRequestReviewHook); ok {
			fmt.Printf("Pull request review event: %s\n", reviewHook.Action)
			fmt.Printf("Repository: %s/%s\n", reviewHook.Repo.Namespace, reviewHook.Repo.Name)
			fmt.Printf("Pull Request #%d: %s\n", reviewHook.PullRequest.Number, reviewHook.PullRequest.Title)
			fmt.Printf("Review by %s: %s\n", reviewHook.Review.Author.Login, reviewHook.Review.Body)
			
			// This is where Drone CI would trigger builds based on review events
			switch reviewHook.Action {
			case scm.ActionCreate: // Review submitted
				fmt.Println("Review submitted - triggering build")
			case scm.ActionUpdate: // Review edited
				fmt.Println("Review edited - updating build status")
			case scm.ActionDelete: // Review dismissed
				fmt.Println("Review dismissed - cancelling build")
			}
		}

		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Webhook server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}