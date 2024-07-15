package main

import (
    "database/sql"
    "fmt"
    "html/template"
    "net/http"
    "strings"
    "log"

    "github.com/gorilla/mux"
    _ "github.com/mattn/go-sqlite3"
)

// Data structure for the greeting template
//type HostData struct {
//    Name string
//}
/// Database initialization function

type AddHostData struct {
    Message string
}

type hostData struct {
    id int
    hostname string
    mac string
    ipv4 string
    ipv6 string
    domain string
    status string
    vlan string
    cnames string
    notes string
}

func initDB() (*sql.DB, error) {
    db, err := sql.Open("sqlite3", "./webapp.db")
    if err != nil {
        return nil, err
    }

    // Create the users table if it doesn't exist
        _, err = db.Exec(`CREATE TABLE IF NOT EXISTS hosts (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            hostname TEXT NOT NULL UNIQUE,
            mac TEXT,
            ipv4 TEXT,
            ipv6 TEXT,
            domain TEXT,
            status TEXT,
            vlan TEXT,
            cnames TEXT,
            notes TEXT
        )
    `)
    if err != nil {
        return nil, err
    }

    return db, nil
}

// Function to insert a user into the database
func insertHost(db *sql.DB, hostname, mac, ipv4, ipv6, domain, status, vlan, cnames, notes string) error {
    _, err := db.Exec("INSERT INTO hosts (hostname, mac, ipv4, ipv6, domain, status, vlan, cnames, notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", hostname, mac, ipv4, ipv6, domain, status, vlan, cnames, notes)
    return err
}

// Function to insert a user into the database
func queryAllHosts(db *sql.DB) []hostData {
    var hostList []hostData

    rows, err := db.Query("SELECT * FROM hosts")
    if err != nil {
        log.Println(err)
    }
    defer rows.Close()

    for rows.Next() {
        var host hostData
        err := rows.Scan(&host.id, &host.hostname, &host.mac, &host.ipv4, &host.ipv6, &host.domain, &host.status, &host.vlan, &host.cnames, &host.notes)
        if err != nil {
            log.Println(err)
        }
    hostList = append(hostList, host)
    }
    return hostList
}


func main() {
    // Initialize the database
    db, err := initDB()
    if err != nil {
        fmt.Println("Error initializing database:", err)
        return
    }
    defer db.Close()

    // Create a new router
    router := mux.NewRouter()

    // Define a handler function for the home page
    homeHandler := func(w http.ResponseWriter, r *http.Request) {
        // Parse the HTML template
        tmpl, err := template.ParseFiles("./web/templates/index.html")
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        // Execute the template
        err = tmpl.Execute(w, nil)
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
    }

    // Define a handler function for the addhost page
    addHostHandler := func(w http.ResponseWriter, r *http.Request) {
        // Check if the request is a POST request
        if r.Method == http.MethodPost {
            // Parse the form data
            err := r.ParseForm()
            if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }

            // Get the username and password from the form
            hostname := strings.TrimSpace(r.Form.Get("hostname"))
            mac := strings.TrimSpace(r.Form.Get("mac"))
            ipv4 := strings.TrimSpace(r.Form.Get("ipv4"))
            ipv6 := strings.TrimSpace(r.Form.Get("ipv6"))
            domain := strings.TrimSpace(r.Form.Get("domain"))
            status := strings.TrimSpace(r.Form.Get("status"))
            vlan := strings.TrimSpace(r.Form.Get("vlan"))
            cnames := strings.TrimSpace(r.Form.Get("cnames"))
            notes := strings.TrimSpace(r.Form.Get("notes"))

            insertHost(db, hostname, mac, ipv4, ipv6, domain, status, vlan, cnames, notes)
            http.Redirect(w, r, "/hostlist", http.StatusFound)

        }

        // Parse the HTML template
        tmpl, err := template.ParseFiles("./web/templates/addhost.html")
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        // Execute the template
        err = tmpl.Execute(w, nil)
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
    }

    // Define a handler function for the host list page
    hostListHandler := func(w http.ResponseWriter, r *http.Request) {
        list := queryAllHosts(db) 
        fmt.Println(list)

        // Parse the HTML template
        tmpl, err := template.ParseFiles("./web/templates/hostlist.html")
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        // Execute the template
        err = tmpl.Execute(w, nil)
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
    }
    // Register the handlers for the home, greet, and login pages
    router.HandleFunc("/", homeHandler).Methods(http.MethodGet)
    router.HandleFunc("/hostlist", hostListHandler).Methods(http.MethodGet)
    router.HandleFunc("/addhost", addHostHandler).Methods(http.MethodGet, http.MethodPost)

    // Serve static files from the "static" directory
    router.PathPrefix("./web/static/").Handler(http.StripPrefix("./web/static/", http.FileServer(http.Dir("static"))))

    // Start the HTTP server on port 8080 using the router
    fmt.Println("Server is listening on :8080...")
    http.ListenAndServe(":8080", router)
}
