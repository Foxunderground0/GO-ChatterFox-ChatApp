package main

import (
	"context"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"time"
	"strconv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

func new_message(user string, message string, collection *mongo.Collection) {
    // Create a context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Use the created context for operations
    _, err := collection.InsertOne(ctx, bson.M{
		"time": time.Now().UnixNano(),
        "usr": user,
        "msg": message,
    })
    if err != nil {
        log.Println("Error inserting message:", err)
    }
}


func main(){
	ctx := context.TODO()
	file_name := "X509-cert-5136463749316099850.pem"
  	uri := "mongodb+srv://chatterfox-texts-db.hwyr17c.mongodb.net/?authSource=%24external&authMechanism=MONGODB-X509&retryWrites=true&w=majority&tlsCertificateKeyFile=" + file_name
  	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
  	clientOptions := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPIOptions)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil { log.Fatal(err) }
	defer client.Disconnect(ctx) // Run at the end to ensure safe disconnect
	
	collection := client.Database("Application-Data").Collection("Chats")
	docCount, err := collection.CountDocuments(ctx, bson.D{})

	fmt.Println(docCount)

	http.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Set the 200 status code
			w.WriteHeader(http.StatusOK)
			
			fmt.Fprintf(w, "OK \n")
		} else {
			http.Error(w, "Only Get requests are supported", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// Read the request body
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading text", http.StatusInternalServerError)
				return
			}
			defer r.Body.Close()

			// Get the IP address from the RemoteAddr field
			ip := r.RemoteAddr

			new_message(string(ip), string(body), collection)

			// Assuming the request body contains text data
			fmt.Fprintf(w, "OK \n")
		} else {
			http.Error(w, "Only POST requests are supported", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/read", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Define the options for sorting and limiting
			options := options.Find().SetSort(bson.D{{"time", -1}}).SetLimit(50)
	
			// Create a cursor with the find options
			cursor, err := collection.Find(ctx, bson.D{}, options)
			if err != nil {
				log.Fatal(err)
			}
			defer cursor.Close(ctx)
	
			// Create a buffer to accumulate the results
			var results []bson.M
	
			// Iterate through the cursor and process each document
			for cursor.Next(ctx) {
				var result bson.M
				if err := cursor.Decode(&result); err != nil {
					log.Fatal(err)
				}
				results = append(results, result)
			}

				// Reverse the order of the results
			for i, j := 0, len(results)-1; i < j; i, j = i+1, j-1 {
				results[i], results[j] = results[j], results[i]
			}
	
			if err := cursor.Err(); err != nil {
				log.Fatal(err)
			}
	        // Print the results to the response writer
			for _, result := range results {
				user, __ := result["usr"].(string)
				time, __ := result["time"].(int64)
				message, __ := result["msg"].(string)
				if __ {
					line := user + " at " + strconv.FormatInt(time, 10) +": " + message
					fmt.Fprintln(w, line)
				}
			}
	
			} else {
			http.Error(w, "Only GET requests are supported", http.StatusMethodNotAllowed)
		}
	})
	

	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)

}
