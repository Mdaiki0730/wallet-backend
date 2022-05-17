package main

import (
  "fmt"
  "log"
  "context"

  "gariwallet/internal/application/usecase"
  "gariwallet/internal/application/command"
)

func main() {
  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()
  fmt.Println("hello")
  wa := usecase.NewWalletApp()
  result, err := wa.Create(ctx, command.WalletCreate{})
  if err != nil {
    log.Fatal("can't create")
  }
  fmt.Println(result)
}
