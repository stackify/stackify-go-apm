package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful/v3"

	apm "github.com/stackify/stackify-go-apm"
	"github.com/stackify/stackify-go-apm/config"
	"github.com/stackify/stackify-go-apm/instrumentation/github.com/emicklei/go-restful/stackifyrestful"
)

func initStackifyTrace() (*apm.StackifyAPM, error) {
	return apm.NewStackifyAPM(
		config.WithApplicationName("Go Application"),
		config.WithEnvironmentName("Test"),
		config.WithDebug(true),
	)
}

type User struct {
	ID int `json:"id" description:"identifier of the user"`
}

type UserResource struct{}

func (u UserResource) findUser(req *restful.Request, resp *restful.Response) {
	uid := req.PathParameter("user-id")
	id, err := strconv.Atoi(uid)
	if err == nil && id >= 100 {
		_ = resp.WriteEntity(User{id})
		return
	}
	_ = resp.WriteErrorString(http.StatusNotFound, "User could not be found.")
}

func (u UserResource) WebService() *restful.WebService {
	ws := &restful.WebService{}

	ws.Path("/users").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{user-id}").To(u.findUser).
		Doc("get a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("integer").DefaultValue("1")).
		Writes(User{}).
		Returns(200, "OK", User{}).
		Returns(404, "Not Found", nil))
	return ws
}

func main() {
	stackifyAPM, err := initStackifyTrace()
	if err != nil {
		log.Fatalf("failed to initialize stackifyapm: %v", err)
	}
	defer stackifyAPM.Shutdown()

	u := UserResource{}
	filter := stackifyrestful.StackifyFilter()
	restful.DefaultContainer.Filter(filter)
	restful.DefaultContainer.Add(u.WebService())

	_ = http.ListenAndServe(":8000", nil)
}
