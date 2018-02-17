package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

type GitlabPush struct {
	// push
	ObjectKind string `json:"object_kind"`
	// refs/heads/master
	Ref string `json:"ref"`
}

func gitHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("X-Gitlab-Token")
	if token != "mi-token-secreto" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"message\":\"El token de Gitlab no coincide con mi X-Gitlab-Token\"}"))
		return
	}

	log.Printf("X-Gitlab-Token recibido: %s\n", token)

	// push de gitlab
	pgit := GitlabPush{}
	err := json.NewDecoder(r.Body).Decode(&pgit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"message\":\"No fue posible decodificar el json\"}"))
		return
	}

	ref := "refs/heads/master"
	if pgit.Ref != ref {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"message\":\"La rama actualizada en gitlab no es master, por lo tanto no se ha llevado a cabo ninguna acción\"}"))
		return
	}

	out, err := exec.Command("git", "pull", "origin", "master").Output()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"message\":\"Error al ejecutar el comando git en el servidor\"}"))
		return
	}
	log.Printf("%s\n", out)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message\": \"ok\"}"))
}

func main() {
	port := fmt.Sprintf(":%d", 1030)

	http.HandleFunc("/", gitHandler)

	fmt.Printf("Ejecutando el servidor en %s\n", port)
	err := http.ListenAndServe(port, nil)
	log.Fatal(err)
}
