package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"
)

var Zone = "us-central1-c"
var JsonPath = "/Users/ganesan.duraisamy/Desktop/Study/Go-Learn/gcpsdk/credentials.json"
var ProjectID = "kubernetes-283301"

func gke_create(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var art container.CreateClusterRequest
	json.Unmarshal(reqBody, &art)

	ctx := context.Background()

	containerService, err := container.NewService(ctx, option.WithCredentialsFile(JsonPath))

	if err != nil {
		log.Fatal(err)
	}

	resp, err := containerService.Projects.Zones.Clusters.Create(ProjectID, Zone, &art).Context(ctx).Do()

	json.NewEncoder(w).Encode(resp)

}

func gke_list(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	containerService, err := container.NewService(ctx, option.WithCredentialsFile(JsonPath))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := containerService.Projects.Zones.Clusters.List(ProjectID, Zone).Context(ctx).Do()

	json.NewEncoder(w).Encode(resp)
}

func gke_get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterid := vars["clusterid"]
	ctx := context.Background()
	containerService, err := container.NewService(ctx, option.WithCredentialsFile(JsonPath))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := containerService.Projects.Zones.Clusters.Get(ProjectID, Zone, clusterid).Context(ctx).Do()
	json.NewEncoder(w).Encode(resp)
}

func gke_delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterid := vars["clusterid"]
	ctx := context.Background()
	containerService, err := container.NewService(ctx, option.WithCredentialsFile(JsonPath))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := containerService.Projects.Zones.Clusters.Delete(ProjectID, Zone, clusterid).Context(ctx).Do()
	json.NewEncoder(w).Encode(resp)
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/api/gke/cluster/create", gke_create).Methods("POST")
	router.HandleFunc("/api/gke/clusters/list", gke_list).Methods("GET")
	router.HandleFunc("/api/gke/cluster/get/{clusterid}", gke_get).Methods("GET")
	router.HandleFunc("/api/gke/cluster/delete/{clusterid}", gke_delete).Methods("DELETE")
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
