package main

import (
  "io"
  "log"
  "net/http"
  "os/exec"
  "sync"
)

func main() {
  cmd := exec.Command("./github-mcp-server", "stdio")
  stdin, err := cmd.StdinPipe()
  if err != nil {
    log.Fatal(err)
  }
  stdout, err := cmd.StdoutPipe()
  if err != nil {
    log.Fatal(err)
  }

  if err := cmd.Start(); err != nil {
    log.Fatal(err)
  }

  var mu sync.Mutex

  http.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    defer mu.Unlock()

    // Forward request body to MCP stdin
    _, err := io.Copy(stdin, r.Body)
    if err != nil {
      http.Error(w, "Failed to write to MCP server", 500)
      return
    }
    r.Body.Close()

    // Read response from MCP stdout and copy to response
    _, err = io.Copy(w, stdout)
    if err != nil {
      http.Error(w, "Failed to read from MCP server", 500)
      return
    }
  })

  log.Println("HTTP MCP proxy listening on :8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
