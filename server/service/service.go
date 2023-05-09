package service

import (
	"log"
	"net/http"
	"path"
	"stacksviz/datasource"

	"github.com/google/traceviz/server/go/handlers"
	querydispatcher "github.com/google/traceviz/server/go/query_dispatcher"
)

type Service struct {
	queryHandler handlers.QueryHandler
	assetHandler *handlers.AssetHandler
}

func New(assetRoot, collectionRoot string) (*Service, error) {
	cf := datasource.NewStacksFetcher(collectionRoot)
	ds := datasource.New(cf)
	qd, err := querydispatcher.New(ds)
	if err != nil {
		return nil, err
	}
	assetHandler := handlers.NewHandler()
	addFileAsset := func(resourceName, resourceType, filename string) {
		log.Printf("Serving asset '%s' at '%s'",
			path.Join(assetRoot, filename),
			resourceName)
		assetHandler.With(
			resourceName,
			handlers.NewFileAsset(
				path.Join(assetRoot, filename),
				resourceType,
			),
		)
	}
	addFileAsset("/logviz-theme.css", "text/css", "logviz-theme.css")
	addFileAsset("/index.html", "text/html", "index.html")
	addFileAsset("main.js", "application/javascript", "main.js")
	addFileAsset("polyfills.js", "application/javascript", "polyfills.js")
	addFileAsset("runtime.js", "application/javascript", "runtime.js")
	return &Service{
		queryHandler: handlers.NewQueryHandler(qd),
		assetHandler: assetHandler,
	}, nil
}

func (s *Service) RegisterHandlers(mux *http.ServeMux) {
	for path, handler := range s.queryHandler.HandlersByPath() {
		mux.HandleFunc(path, handler)
	}
}
