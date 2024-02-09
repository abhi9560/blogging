package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql" // MySQL driver
)

// Post represents a blog post
type Post struct {
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Content   string `json:"content"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}

var db *sql.DB

func main() {
    // Initialize database connection
    initDB()

    // Create a new router
    router := mux.NewRouter()

    // Define API endpoints
    router.HandleFunc("/posts", getAllPosts).Methods("GET")
    router.HandleFunc("/posts/{id}", getPostByID).Methods("GET")
    router.HandleFunc("/posts", createPost).Methods("POST")
    router.HandleFunc("/posts/{id}", updatePost).Methods("PUT")
    router.HandleFunc("/posts/{id}", deletePost).Methods("DELETE")

    // Start the server
    http.ListenAndServe(":8080", router)
}

func initDB() {
    // Open database connection
    var err error
    db, err = sql.Open("mysql", "username:password@tcp(localhost:3306)/dbname")
    if err != nil {
        panic(err.Error())
    }

    // Check if the connection is successful
    err = db.Ping()
    if err != nil {
        panic(err.Error())
    }

    fmt.Println("Connected to the database")
}

func getAllPosts(w http.ResponseWriter, r *http.Request) {
    // Retrieve all posts from the database
    rows, err := db.Query("SELECT id, title, content, created_at, updated_at FROM posts")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // Create a slice to hold the posts
    var posts []Post

    // Iterate over the rows
    for rows.Next() {
        var post Post
        // Scan the values from the row into the post struct
        err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        // Append the post to the slice
        posts = append(posts, post)
    }
    if err := rows.Err(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Convert the posts slice to JSON
    jsonPosts, err := json.Marshal(posts)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Set the content type header and write the JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonPosts)
}

func getPostByID(w http.ResponseWriter, r *http.Request) {
    // Extract the post ID from the request URL
    vars := mux.Vars(r)
    postID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    // Retrieve the post from the database
    var post Post
    err = db.QueryRow("SELECT id, title, content, created_at, updated_at FROM posts WHERE id = ?", postID).
        Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Post not found", http.StatusNotFound)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    // Convert the post struct to JSON
    jsonPost, err := json.Marshal(post)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Set the content type header and write the JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonPost)
}

func createPost(w http.ResponseWriter, r *http.Request) {
    // Parse the request body to get the new post data
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer r.Body.Close()

    var newPost Post
    if err := json.Unmarshal(body, &newPost); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Execute an INSERT query to add the new post to the database
    result, err := db.Exec("INSERT INTO posts (title, content) VALUES (?, ?)", newPost.Title, newPost.Content)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Get the ID of the newly created post
    postID, err := result.LastInsertId()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Return the ID of the newly created post
    response := map[string]int{"id": int(postID)}
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Set the content type header and write the JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(jsonResponse)
}

func updatePost(w http.ResponseWriter, r *http.Request) {
    // Extract the post ID from the request URL
    vars := mux.Vars(r)
    postID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    // Parse the request body to get the updated post data
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer r.Body.Close()

    var updatedPost Post
    if err := json.Unmarshal(body, &updatedPost); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Execute an UPDATE query to modify the post in the database
    _, err = db.Exec("UPDATE posts SET title = ?, content = ? WHERE id = ?", updatedPost.Title, updatedPost.Content, postID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Return a success message
    successMessage := map[string]string{"message": "Post updated successfully"}
    jsonResponse, err := json.Marshal(successMessage)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Set the content type header and write the JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
    // Extract the post ID from the request URL
    vars := mux.Vars(r)
    postID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    // Execute a DELETE query to remove the post from the database
    _, err = db.Exec("DELETE FROM posts WHERE id = ?", postID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Return a success message
    successMessage := map[string]string{"message": "Post deleted successfully"}
    jsonResponse, err := json.Marshal(successMessage)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Set the content type header and write the JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}
