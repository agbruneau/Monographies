package mcpserver

import "github.com/modelcontextprotocol/go-sdk/mcp"

// toolError construit un résultat d'outil en **erreur métier** (IsError=true).
// Conforme à ADR-0001 : les erreurs métier (NotFound, règle violée) passent par
// IsError et le contenu texte, jamais par une erreur protocole -32602 (celle-ci
// est réservée au protocole ; cf. M1.7/FR-15).
func toolError(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
	}
}
