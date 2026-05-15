package apirest

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Muxcore-Media/core/pkg/contracts"
)

func init() {
	contracts.Register(func(deps contracts.ModuleDeps) contracts.Module {
		return NewModule(deps.Registry, deps.Routes)
	})
}

type Module struct {
	reg    contracts.ServiceRegistry
	routes contracts.RouteRegistrar
}

func NewModule(reg contracts.ServiceRegistry, routes contracts.RouteRegistrar) *Module {
	return &Module{reg: reg, routes: routes}
}

func (m *Module) Info() contracts.ModuleInfo {
	return contracts.ModuleInfo{
		ID:           "api-rest",
		Name:         "REST API",
		Version:      "1.0.0",
		Kind:         contracts.ModuleKindAPI,
		Description:  "REST API endpoints for module and system management",
		Author:       "MuxCore",
		Capabilities: []string{"api.rest", "api.v1"},
	}
}

func (m *Module) Init(ctx context.Context) error { return nil }

func (m *Module) Start(ctx context.Context) error {
	m.routes.HandleFunc("/api/v1/modules", m.handleModules)
	m.routes.HandleFunc("/api/v1/modules/", m.handleModuleByID)
	slog.Info("REST API routes registered")
	return nil
}

func (m *Module) Stop(ctx context.Context) error  { return nil }
func (m *Module) Health(ctx context.Context) error { return nil }

func (m *Module) handleModules(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var entries []contracts.ModuleEntry
	if kind := r.URL.Query().Get("kind"); kind != "" {
		entries = m.reg.FindByKind(contracts.ModuleKind(kind))
	} else {
		entries = m.reg.ListAll()
	}

	type moduleJSON struct {
		ID           string   `json:"id"`
		Name         string   `json:"name"`
		Version      string   `json:"version"`
		Kind         string   `json:"kind"`
		Description  string   `json:"description"`
		Capabilities []string `json:"capabilities"`
		// State is not exposed via ServiceRegistry by design — it's internal to core
	}

	modules := make([]moduleJSON, 0, len(entries))
	for _, e := range entries {
		modules = append(modules, moduleJSON{
			ID:           e.Info.ID,
			Name:         e.Info.Name,
			Version:      e.Info.Version,
			Kind:         string(e.Info.Kind),
			Description:  e.Info.Description,
			Capabilities: e.Info.Capabilities,
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{"modules": modules})
}

func (m *Module) handleModuleByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/api/v1/modules/"):]
	if id == "" {
		http.Error(w, "module id required", http.StatusBadRequest)
		return
	}
	// Strip any trailing path segments
	if idx := strings.Index(id, "/"); idx >= 0 {
		id = id[:idx]
	}

	entry, err := m.reg.Resolve(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":           entry.Info.ID,
		"name":         entry.Info.Name,
		"version":      entry.Info.Version,
		"kind":         string(entry.Info.Kind),
		"description":  entry.Info.Description,
		"author":       entry.Info.Author,
		"capabilities": entry.Info.Capabilities,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
