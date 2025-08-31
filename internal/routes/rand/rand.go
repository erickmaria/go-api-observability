package routerand

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func Random(w http.ResponseWriter, r *http.Request) {
	// defer logger.NewLogger(r.Context())()

	slog.InfoContext(r.Context(), "random number generated")

	number := randomNumber(r.Context())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{\"number\": \"%d\"}", number)))
	json.NewEncoder(w)
}

func randomNumber(ctx context.Context) int {

	file, _ := os.Create("test.db")
	defer file.Close()

	database, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}

	database.Conn(ctx)
	conn, err := database.Conn(ctx)
	if err != nil {
		panic(err)
	}
	conn.Close()

	return rand.IntN(1000)
}
