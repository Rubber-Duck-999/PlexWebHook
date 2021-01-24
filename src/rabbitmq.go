package main

import (
	"time"

	"github.com/clarketm/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var ch *amqp.Channel
var init_err error
var password string

//

func init() {
	log.Trace("Initialised rabbitmq package")
	password = ""
	_, _, day := time.Now().Date()
	log.Debug("Day currently: ", day)
}

func SetCodes(pass string) {
	password = pass
}

func failOnError(err error, msg string) {
	if err != nil {
		log.WithFields(log.Fields{
			"Message": msg, "Error": err,
		}).Error("Rabbitmq error")
	}
}

func getTime() string {
	t := time.Now()
	log.Trace(t.Format(TIMEFORMAT))
	return t.Format(TIMEFORMAT)
}

func messages(routing_key string, value string) {
	log.Warn("Adding messages to map")
	if SubscribedMessagesMap == nil {
		log.Warn("Creation of messages map")
		SubscribedMessagesMap = make(map[uint32]*MapMessage)
		messages(routing_key, value)
	} else {
		_, valid := SubscribedMessagesMap[key_id]
		if valid {
			log.Debug("Key already exists, checking next field: ", key_id)
			key_id++
			messages(routing_key, value)
		} else {
			log.Debug("Key does not exists, adding new field: ", key_id)
			entry := MapMessage{value, routing_key, getTime(), true}
			SubscribedMessagesMap[key_id] = &entry
			key_id++
		}
	}
}

func SetConnection() error {
	conn, init_err = amqp.Dial("amqp://guest:" + password + "@localhost:5672/")
	failOnError(init_err, "Failed to connect to RabbitMQ")

	ch, init_err = conn.Channel()
	failOnError(init_err, "Failed to open a channel")
	return init_err
}

func Subscribe() {
	init := SetConnection()

	DevicesList = make(map[uint32]*Device)
	if init == nil {
		var topics = [1]string{
			DEVICERESPONSE,
		}

		err := ch.ExchangeDeclare(
			EXCHANGENAME, // name
			EXCHANGETYPE, // type
			true,         // durable
			false,        // auto-deleted
			false,        // internal
			false,        // no-wait
			nil,          // arguments
		)
		failOnError(err, "FH - Failed to declare an exchange")

		q, err := ch.QueueDeclare(
			"",    // name
			false, // durable
			false, // delete when usused
			true,  // exclusive
			false, // no-wait
			nil,   // arguments
		)
		failOnError(err, "Failed to declare a queue")

		for _, s := range topics {
			log.Printf("Binding queue %s to exchange %s with routing key %s",
				q.Name, EXCHANGENAME, s)
			err = ch.QueueBind(
				q.Name,       // queue name
				s,            // routing key
				EXCHANGENAME, // exchange
				false,
				nil)
			failOnError(err, "Failed to bind a queue")
		}

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto ack
			false,  // exclusive
			false,  // no local
			false,  // no wait
			nil,    // args
		)
		failOnError(err, "Failed to register a consumer")

		forever := make(chan bool)

		go func() {
			for d := range msgs {
				log.Trace("Sending message to callback")
				log.Trace(d.RoutingKey)
				s := string(d.Body)
				messages(d.RoutingKey, s)
				log.Debug("Checking states of received messages")
				checkState()
			}
			//This function is checked after to see if multiple errors occur then to
			//through an event message
		}()

		go checkDevices()

		log.Trace(" [*] Waiting for logs. To exit press CTRL+C")
		<-forever
	}
}

func Publish(message []byte, routingKey string) string {
	if init_err == nil {
		err := ch.Publish(
			EXCHANGENAME, // exchange
			routingKey,   // routing key
			false,        // mandatory
			false,        // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte(message),
			})
		if err != nil {
			log.Fatal(err)
			return FAILUREPUBLISH
		}
	}
	return ""
}

func PublishFailureNetwork(time string, reason string) string {
	failureNetwork, err := json.Marshal(&FailureNetwork{
		Time:         time,
		Failure_type: reason})
	if err != nil {
		return "Failed to convert FailureNetwork"
	} else {
		return Publish(failureNetwork, FAILURENETWORK)
	}
}

func PublishDeviceFound(name string, address string, status int) string {
	device, err := json.Marshal(&DeviceFound{
		Device_name: name,
		Ip_address:  address,
		Status:      status})
	if err != nil {
		return "Failed to convert DeviceFound"
	} else {
		return Publish(device, DEVICEFOUND)
	}
}

func PublishDeviceRequest(id uint32, name string, mac string) string {
	device, err := json.Marshal(&DeviceRequest{
		Request_id: id,
		Name:       name,
		Mac:        mac})
	if err != nil {
		return "Failed to convert DeviceRequest"
	} else {
		return Publish(device, DEVICEREQUEST)
	}
}
