package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var rabbit_host = os.Getenv("RABBIT_HOST")
var rabbit_port = os.Getenv("RABBIT_PORT")
var rabbit_user = os.Getenv("RABBIT_USERNAME")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")
var publisher_port = os.Getenv("PUBLISHER_PORT")

type chatMessage struct{
	Id string `json:"id"`
	ChatId string `json:"chatId"`
	Content	string `json:"content"`
	From string `json:"from"`
	To string `json:"to"`
	CreatedAt time.Time `json:"createdAt"`
}

func main() {

	router := gin.Default()

	router.POST("/messages", submit)

	fmt.Println("Welcome aboard, Captain! All systems online")
	fmt.Println("Running on 0.0.0.0:"+publisher_port)
	router.Run("0.0.0.0:"+publisher_port)
}

func submit(c *gin.Context) {
	var newMessage chatMessage
	if err := c.BindJSON(&newMessage); err != nil {
		return
	}
	fmt.Println("POST Message data: ", newMessage)

	message := newMessage.Content
	queueName := newMessage.ChatId

	fmt.Println("Received message: " + message)

	rabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbit_user, rabbit_password, rabbit_host, rabbit_port)

	conn, dialErr := amqp.Dial(rabbitUrl)

	if dialErr != nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", dialErr)
	}

	defer conn.Close()

	channel, channelErr := conn.Channel()

	if channelErr != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", channelErr)
	}

	defer channel.Close()

	queue, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	publishErr := channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

	if publishErr != nil {
		log.Fatalf("%s: %s", "Failed to publish a message", publishErr)
	}
	c.IndentedJSON(http.StatusCreated, newMessage)

	fmt.Println("publish to queue successfully!")
}
