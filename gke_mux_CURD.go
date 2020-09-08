package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"fmt"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"
	"github.com/spf13/viper"
)

type Configuration struct {
        Server Server
        App    App
}

type Server struct {
        Port string
}

type App struct {
        Zone string
        ProjectID string
        CredentialsPath string
}

var (
    global_config Configuration
)

func read_config_file(configFile string)(bool, error) {
	viper.SetConfigName(configFile)
	viper.AddConfigPath(".")
        viper.AddConfigPath("/etc/gke")
        viper.AutomaticEnv()

        if err := viper.ReadInConfig(); err != nil {
                return false, fmt.Errorf("Error reading config file, %s", err)
        }
        err := viper.Unmarshal(&global_config)
        if err != nil {
                return false, fmt.Errorf("Error unmarshal config file, %s", err)
        }
        return true, nil
}

func gke_create(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var art container.CreateClusterRequest
	json.Unmarshal(reqBody, &art)

	ctx := context.Background()

	containerService, err := container.NewService(ctx, option.WithCredentialsFile(global_config.App.CredentialsPath))

	if err != nil {
		log.Fatal(err)
	}

	resp, err := containerService.Projects.Zones.Clusters.Create(global_config.App.ProjectID, global_config.App.Zone, &art).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(resp)

}

func gke_list(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	containerService, err := container.NewService(ctx, option.WithCredentialsFile(global_config.App.CredentialsPath))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := containerService.Projects.Zones.Clusters.List(global_config.App.ProjectID, global_config.App.Zone).Context(ctx).Do()

	json.NewEncoder(w).Encode(resp)
}

func gke_get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterid := vars["clusterid"]
	ctx := context.Background()
	containerService, err := container.NewService(ctx, option.WithCredentialsFile(global_config.App.CredentialsPath))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := containerService.Projects.Zones.Clusters.Get(global_config.App.ProjectID, global_config.App.Zone, clusterid).Context(ctx).Do()
	json.NewEncoder(w).Encode(resp)
}

func gke_delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterid := vars["clusterid"]
	ctx := context.Background()
	containerService, err := container.NewService(ctx, option.WithCredentialsFile(global_config.App.CredentialsPath))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := containerService.Projects.Zones.Clusters.Delete(global_config.App.ProjectID, global_config.App.Zone, clusterid).Context(ctx).Do()
	json.NewEncoder(w).Encode(resp)
}

func main() {
	var configfilename = "app"
	var _, err = read_config_file(configfilename)
	if err != nil {
		log.Fatal(err)
	}

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
