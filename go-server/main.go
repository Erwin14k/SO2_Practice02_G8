package main

import (
	"github.com/gorilla/mux" // Package for creating HTTP routers
	"github.com/rs/cors" // Package for Cross-Origin Resource Sharing (CORS) support
	"encoding/json" // Package for JSON encoding and decoding
	"io/ioutil" // Package for reading and writing files and data streams
	"net/http" // Package for HTTP client and server implementations
	"os/exec" // Package for executing external commands
	"strconv" // Package for string conversions
	"strings" // Package for manipulating strings
	"log" // Package for logging
	"fmt" // Package for formatted I/O
)

// Process represents a process with its properties
type Process struct {
	Pid     int    `json:"pid"`
	Nombre  string `json:"nombre"`
	Usuario string `json:"usuario"`
	Estado  string `json:"estado"`
	Ram     int    `json:"ram"`
	Padre   int    `json:"padre"`
}

// CPUInfo represents CPU information and process tasks
type CPUInfo struct {
	TotalCPU int       `json:"totalcpu"`
	Running  int       `json:"running"`
	Sleeping int       `json:"sleeping"`
	Stopped  int       `json:"stopped"`
	Zombie   int       `json:"zombie"`
	Total    int       `json:"total"`
	Tasks    []Process `json:"tasks"`
}

// RAMInfo represents RAM information
type RAMInfo struct {
	TotalRAM    int `json:"totalram"`
	RAMLibre    int `json:"ramlibre"`
	RAMOcupada  int `json:"ramocupada"`
}

// General represents general system information
type general struct {
	TotalRAM    int `json:"totalram"`
	RAMLibre    int `json:"ramlibre"`
	RAMOcupada  int `json:"ramocupada"`
	TotalCPU    int `json:"totalcpu"`
}

// Counters represents process counters
type counters struct {
	Running  int       `json:"running"`
	Sleeping int       `json:"sleeping"`
	Stopped  int       `json:"stopped"`
	Zombie   int       `json:"zombie"`
	Total    int       `json:"total"`
}

// AllData represents all system data
type AllData struct {
	AllGenerales    []general    `json:"AllGenerales"`
	AllTipoProcesos []Process  `json:"AllTipoProcesos"`
	AllProcesos     []counters   `json:"AllProcesos"`
}

// MemoryRow represents a row of memory information
type MemoryRow struct {
	InitialAddress string   `json:"initial_address"`
	FinalAddress   string   `json:"final_address"`
	SizeKB         int      `json:"size_kb"`
	Permissions    []string `json:"permissions"`
	Device         string   `json:"device"`
	File           string   `json:"file"`
}

// createData creates the system data
func createData() (string, error) {
	// read /proc/mem_grupo8 file
	outRAM, err := ioutil.ReadFile("/proc/mem_grupo8")
	if err != nil {
		fmt.Println(err)
	}
	// read /proc/cpu_grupo8 file
	outCPU, err := ioutil.ReadFile("/proc/cpu_grupo8")
	if err != nil {
		fmt.Println(err)
	}
	// --------- PROCESS ---------
	var cpuInfo CPUInfo
	err = json.Unmarshal(outCPU, &cpuInfo)
	if err != nil {
		fmt.Println("Error: Cpu json unmarshal failed", err)
		return "", err
	}

	for i, task := range cpuInfo.Tasks {
		uid, err := strconv.Atoi(task.Usuario)
		if err != nil {
			fmt.Println("Error: Failed to convert UID to int", err)
			return "", err
		}

		// Sh -> Interpreter
		// -c -> Read the command from the argument string
		// grep -m 1 -> Show the first match
		// cut -d: -f1 -> Cuts the string at the first delimiter and displays the first field
		cmdUsr := exec.Command("sh", "-c", "grep -m 1 '"+strconv.Itoa(uid)+":' /etc/passwd | cut -d: -f1")
		

		outUsr, err := cmdUsr.Output()
		if err != nil {
			fmt.Println("Error: Failed to get username for UID ", task.Usuario, err)
			return "", err
		}
		username := strings.TrimSpace(string(outUsr))
		cpuInfo.Tasks[i].Usuario = username
	}

	// --------- RAM ---------
	var ramInfo RAMInfo
	err = json.Unmarshal(outRAM, &ramInfo)
	if err != nil {
		fmt.Println("Error: Ram json unmarshal failed", err)
		return "", err
	}
	
	// allData struct variable contains the data to send to frontend
	allData := AllData{
		AllGenerales: []general{
			{
				TotalRAM:     ramInfo.TotalRAM,
				RAMLibre:     ramInfo.RAMLibre,
				RAMOcupada:   ramInfo.RAMOcupada,
				TotalCPU:     cpuInfo.TotalCPU,
			},
		},
		AllTipoProcesos: cpuInfo.Tasks,
		AllProcesos: []counters{
			{
				Running: cpuInfo.Running,
				Sleeping: cpuInfo.Sleeping,
				Stopped: cpuInfo.Stopped,
				Zombie: cpuInfo.Zombie,
				Total: cpuInfo.Total,
			},
		},
	}

	// Marshal allData variable to json
	allDataJSON, err := json.Marshal(allData)
	if err != nil {
		fmt.Println("Error: AllData json marshal failed", err)
		return "", err
	}

	return string(allDataJSON), nil
}

// handleGet handles the GET request
func handleGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET")

	allData, err := createData() // Get all process data in JSON format
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Return HTTP 500 Internal Server Error if data is empty
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, allData)
}

// handlePost handles the POST request
func handlePost(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body) // Read the request body
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Return HTTP 500 Internal Server Error if there's an error reading the body
		return
	}

	pid, err := strconv.Atoi(string(body)) // Convert the body to an integer (PID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // Return HTTP 400 Bad Request if the body is not a valid PID
		fmt.Fprintln(w, "Information: Invalid PID")
		return
	}

	cmd := exec.Command("sudo", "kill", strconv.Itoa(pid)) // Create a command to kill the process
	err = cmd.Run()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Return HTTP 500 Internal Server Error if there's an error killing the process
		fmt.Fprintln(w, "Error killing process")
		return
	}

	fmt.Println("Information: Process with PID", pid, "has been deleted") // Print the information about the deleted process
	w.WriteHeader(http.StatusOK)              // Set HTTP 200 OK status code
	fmt.Fprintln(w, "Process deleted")        // Write the response message to the response writer
}

// handleRoute handles the route "/"
func handleRoute(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to my API :D")
}

// handleMemory handles the route "/memory"
func handleMemory(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body) // Leer el cuerpo de la solicitud
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Devolver HTTP 500 Internal Server Error si hay un error al leer el cuerpo
		return
	}

	pid, err := strconv.Atoi(string(body)) // Convertir el cuerpo a un entero (PID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // Devolver HTTP 400 Bad Request si el cuerpo no es un PID válido
		fmt.Fprintln(w, "Invalid PID")
		return
	}

	// Crear un comando para ejecutar el cat en /proc/pid/maps
	cmd := exec.Command("sudo", "cat", fmt.Sprintf("/proc/%d/maps", pid))
	output, err := cmd.Output()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Devolver HTTP 500 Internal Server Error si hay un error al ejecutar el comando
		fmt.Fprintln(w, "Error reading process memory")
		return
	}

	rows := parseMemoryRows(string(output))
	jsonData, err := json.Marshal(rows)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Devolver HTTP 500 Internal Server Error si hay un error al generar el JSON
		fmt.Fprintln(w, "Error generating JSON")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Establecer el código de estado HTTP 200 OK

	// Escribir los datos del proceso en la respuesta
	fmt.Fprintln(w, string(jsonData))
}

// parseMemoryRows parses the memory rows from the output of `cat /proc/pid/maps`
func parseMemoryRows(output string) []MemoryRow {
	rows := []MemoryRow{}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		initialAddress := strings.Split(fields[0], "-")[0]
		finalAddress := strings.Split(fields[0], "-")[1]
		permissions := fields[1]
		device := fields[3]
		file := fields[5]

		sizeKB := calculateSizeKB(initialAddress, finalAddress)

		permissionsList := parsePermissions(permissions)

		row := MemoryRow{
			InitialAddress: initialAddress,
			FinalAddress:   finalAddress,
			SizeKB:         sizeKB,
			Permissions:    permissionsList,
			Device:         device,
			File:           file,
		}

		rows = append(rows, row)
	}

	return rows
}

// calculateSizeKB calculates the size in kilobytes from the initial and final addresses
func calculateSizeKB(initialAddress, finalAddress string) int {
	initial, err := strconv.ParseUint(initialAddress, 16, 64)
	if err != nil {
		return 0
	}

	final, err := strconv.ParseUint(finalAddress, 16, 64)
	if err != nil {
		return 0
	}

	size := final - initial
	sizeKB := int(size) / 1024

	return sizeKB
}

// parsePermissions parses the permissions string and returns a list of permissions
func parsePermissions(permissions string) []string {
	perms := []string{}

	if strings.Contains(permissions, "r") {
		perms = append(perms, "Lectura")
	}
	if strings.Contains(permissions, "w") {
		perms = append(perms, "Escritura")
	}
	if strings.Contains(permissions, "x") {
		perms = append(perms, "Ejecucion")
	}

	return perms
}

func main() {
	fmt.Println("************************************************************")
	fmt.Println("*                 SO2 Practica 2 - Grupo 8                 *")
	fmt.Println("************************************************************")

	router := mux.NewRouter().StrictSlash(true) // Create a new router instance
	router.HandleFunc("/", handleRoute) // Set the handler function for the root route ("/")
	router.HandleFunc("/tasks", handlePost).Methods("POST") // Set the handler function for the "/tasks" route with POST method
	router.HandleFunc("/tasks", handleGet).Methods("GET") // Set the handler function for the "/tasks" route with GET method
	router.HandleFunc("/memory", handleMemory).Methods("POST")

	handler := cors.Default().Handler(router) // Create a new CORS handler with default settings
	log.Fatal(http.ListenAndServe(":8080", handler)) // Start the HTTP server and listen on port 8080

	fmt.Println("Server on port 8080")
}
