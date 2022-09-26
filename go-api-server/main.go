package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"log"
	"errors"
)

type book struct {
	// `` is used to tell that this data will be send as json in api
	// variables here are started from capital letter, so that outside modules
	// can also access it
	// in json we want the fields in lower case, so we added the `` block which
	// says that ID will be sent as id in json response
	
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int `json:"quantity"`
}

// slice of book struct
var books = []book{
	{ID: "1", Title: "A1", Author: "a1", Quantity: 1},
	{ID: "2", Title: "A2", Author: "a2", Quantity: 2},
	{ID: "3", Title: "A3", Author: "a3", Quantity: 3},
	{ID: "4", Title: "A4", Author: "a4", Quantity: 4},
}

func getBooks(c *gin.Context) {
	// *gin.Context -> Has all the information about the request
	c.IndentedJSON(http.StatusOK, books)
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return 
	}
	c.IndentedJSON(http.StatusOK, book)
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	log.Println("id=", id)
	log.Println("ok=", ok)

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter"})
		return
	}
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return 
	}
	
	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available."})
		return 	
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	} 
	return nil, errors.New("book not found")
}

func createBook(c *gin.Context)  {
	var newBook book
	log.Println(c.Request.Body)
	if err := c.BindJSON(&newBook); err != nil {
		// Below is the code for debugging was very helpful
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "json decoding : " + err.Error(),
			"status": http.StatusBadRequest,
		})
		return
	}

	log.Println(newBook)
	log.Println("3.............")

	books = append(books, newBook)
	log.Println("4.............")
	c.IndentedJSON(http.StatusCreated, newBook)
}

func main()  {
	router := gin.Default()
	// curl localhost:8081/books 
	router.GET("/books", getBooks)
	// curl  -d '{"ID": "5","Title": "H","Author": "A","Quantity": 4}' -X POST http://localhost:8081/books
	router.POST("/books", createBook)
	// curl localhost:8081/books/2
	router.GET("/books/:id", bookById)
	// curl --request "PATCH" localhost:8081/checkout?id=1
	router.PATCH("/checkout", checkoutBook)
	router.Run("localhost:8081")
}