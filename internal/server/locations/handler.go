package locations

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/alexPavlikov/go-wiki/internal/server/service"
	"golang.org/x/net/html"
)

type Handler struct {
	service service.Service
}

func New(service service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) WikiHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	start := ctx.Value("START_LINK").(string)
	end := ctx.Value("END_LINK").(string)

	urls, err := h.get(start)
	if err != nil {
		http.NotFound(w, r)
	}

	var checkURL = make(map[string]string, 0)

	for v := range urls {

		_, ok := checkURL[v]
		if !ok {
			checkURL[v] = v

			links, _ := h.get(v)
			_, ok := links[end]
			if !ok {
				for v := range links {
					_, ok := checkURL[v]
					if !ok {
						checkURL[v] = v
						links, _ := h.get(v)
						val, ok := links[end]
						if ok {
							var b bytes.Buffer
							html.Render(&b, val.Parent.Parent)
							fmt.Println("FINISH2", b.String())
							return
						}
					}
				}
			}
		}

	}

}

func (h *Handler) get(url string) (map[string]*html.Node, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		return nil, err
	}

	urls, err := h.service.ParseString(doc)
	if err != nil {
		return nil, err
	}

	return urls, err
}
