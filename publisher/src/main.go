package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var rabbit_host = os.Getenv("RABBIT_HOST")
var rabbit_port = os.Getenv("RABBIT_PORT")
var rabbit_user = os.Getenv("RABBIT_USERNAME")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")

func main() {

	router := httprouter.New()

	router.POST("/:queue/:message", func(_ http.ResponseWriter, req *http.Request, params httprouter.Params) {
		submit(req, params)
	})

	fmt.Println("Running!...")
	log.Fatal(http.ListenAndServe(":4000", router))
}

func submit(request *http.Request, params httprouter.Params) {
	message := params.ByName("message")
	queueName := params.ByName("queue")

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

	fmt.Println("publish to queue successfully!")
}
