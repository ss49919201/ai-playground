package mcp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/ss49919201/ai-kata/go/mcp-password-generator/internal/config"
)

type Server struct {
	config                   *config.Config
	tools                    map[string]Tool
	toolsMutex               sync.RWMutex
	handlers                 map[string]MethodHandler
	consentMap               map[string]bool
	consentLock              sync.RWMutex
	passwordGeneratorHandler func(input json.RawMessage) (interface{}, *Error)
}

type MethodHandler func(params json.RawMessage) (interface{}, *Error)

func NewServer(cfg *config.Config) *Server {
	s := &Server{
		config:     cfg,
		tools:      make(map[string]Tool),
		handlers:   make(map[string]MethodHandler),
		consentMap: make(map[string]bool),
	}

	s.handlers["tools/list"] = s.handleToolsList
	s.handlers["tools/call"] = s.handleToolsCall

	return s
}

func (s *Server) RegisterTool(tool Tool) {
	s.toolsMutex.Lock()
	defer s.toolsMutex.Unlock()
	s.tools[tool.Name] = tool
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp, _ := NewResponse(nil, nil, NewError(ParseError, "Parse error", err.Error()))
		writeResponse(w, resp)
		return
	}

	if req.JSONRPC != JSONRPCVersion {
		resp, _ := NewResponse(req.ID, nil, NewError(InvalidRequest, "Invalid JSON-RPC version", nil))
		writeResponse(w, resp)
		return
	}

	handler, ok := s.handlers[req.Method]
	if !ok {
		resp, _ := NewResponse(req.ID, nil, NewError(MethodNotFound, "Method not found", nil))
		writeResponse(w, resp)
		return
	}

	result, err := handler(req.Params)
	if err != nil {
		resp, _ := NewResponse(req.ID, nil, err)
		writeResponse(w, resp)
		return
	}

	resp, _ := NewResponse(req.ID, result, nil)
	writeResponse(w, resp)
}

func (s *Server) handleToolsList(params json.RawMessage) (interface{}, *Error) {
	var p ToolsListParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, NewError(InvalidParams, "Invalid parameters", err.Error())
	}

	s.toolsMutex.RLock()
	defer s.toolsMutex.RUnlock()

	tools := make([]Tool, 0, len(s.tools))
	for _, tool := range s.tools {
		tools = append(tools, tool)
	}

	return ToolsListResult{Tools: tools}, nil
}

func (s *Server) handleToolsCall(params json.RawMessage) (interface{}, *Error) {
	var p ToolsCallParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, NewError(InvalidParams, "Invalid parameters", err.Error())
	}

	s.toolsMutex.RLock()
	_, ok := s.tools[p.Name]
	s.toolsMutex.RUnlock()
	if !ok {
		return nil, NewError(ToolNotFound, fmt.Sprintf("Tool '%s' not found", p.Name), nil)
	}

	if s.config.MCP.RequireConsent {
		if p.Consent == nil || !*p.Consent {
			return nil, NewError(ConsentRequired, "Consent is required to execute this tool", nil)
		}
	}

	switch p.Name {
	case "generate_password":
		return s.passwordGeneratorHandler(p.Input)
	default:
		return nil, NewError(ToolExecutionFail, fmt.Sprintf("Tool '%s' execution not implemented", p.Name), nil)
	}
}

func writeResponse(w http.ResponseWriter, resp *Response) {
	w.Header().Set("Content-Type", "application/json")
	if resp.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
