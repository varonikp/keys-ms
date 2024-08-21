package httpserver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/varonikp/keys-ms/internal/common/server"
)

func getIDFromRequest(idKey string, w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars[idKey])
	if err != nil {
		server.BadRequest(fmt.Sprintf("missing %s", idKey), err, w)
		return 0, err
	}

	return id, nil
}
