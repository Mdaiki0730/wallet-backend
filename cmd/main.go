package main

import (
  "fmt"
  "log"
  "context"

  "gariwallet/internal/application/usecase"
  "gariwallet/internal/infrastructure"
  "gariwallet/pkg/database"
)

func init() {
  log.SetPrefix("Wallet Server: ")
}

func main() {
  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()

  mongoClient := database.ConnectMongoDB(ctx)
  defer database.DisconnectDB(ctx, mongoClient)

  wr := infrastructure.NewWalletRepository(mongoClient)
  wa := usecase.NewWalletApp(wr)
  result, err := wa.Create(ctx)
  if err != nil {
    log.Fatal("can't create")
  }
  fmt.Println(result)
}
