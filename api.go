package main

import (
	"encoding/json"
	"log"
	"net/http"

	"go.temporal.io/sdk/client"
)

type API struct {
	temporal client.Client
}

func NewAPI(c client.Client) *API {
	return &API{temporal: c}
}

// POST /start-call
func (a *API) startCallHandler(w http.ResponseWriter, r *http.Request) {
	var lead Lead
	if err := json.NewDecoder(r.Body).Decode(&lead); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:        "lead-" + lead.ID,
		TaskQueue: "lead-calling-task-queue",
	}

	we, err := a.temporal.ExecuteWorkflow(r.Context(), workflowOptions, LeadCallingWorkflow, lead)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Workflow started:", we.GetID())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"workflowId": we.GetID(),
		"runId":      we.GetRunID(),
	})
}

// POST /3cx/webhook
func (a *API) threeCXWebhookHandler(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		LeadID string `json:"leadId"`
		Result string `json:"result"`
	}

	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	workflowID := "lead-" + payload.LeadID

	err := a.temporal.SignalWorkflow(
		r.Context(),
		workflowID,
		"",
		"call-result",
		payload.Result,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Signal sent to workflow:", workflowID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("signal sent"))
}
