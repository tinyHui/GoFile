package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tinyhui/GoFile/fileop"
	"github.com/tinyhui/GoFile/router"
	"github.com/tinyhui/GoFile/utils"
	"github.com/tinyhui/GoFile/utils/log"
)

func main() {
	var logger = log.GetLogger()

	parameters := utils.LoadParameters()

	handler := router.NewHandler(
		parameters.StorageRoot,
		fileop.NewFileOp(),
		fileop.NewFileStatic(),
	)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", parameters.Port),
		Handler:      router.InitRouter(handler),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		logger.Errorln(err)
	}
}
