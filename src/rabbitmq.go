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
var pinCode int

func init() {
	log.Trace("Initialised rabbitmq package")
	conn, init_err = amqp.Dial("amqp://guest:password@localhost:5672/")
	failOnError(init_err, "Failed to connect to RabbitMQ")

	ch, init_err = conn.Channel()
	failOnError(init_err, "Failed to open a channel")
}

func SetCodes(pass string, pin int) {
	password = pass
	pinCode = pin 
}

func failOnError(err error, msg string) {
	if err != nil {
		log.WithFields(log.Fields{
			"Message": msg, "Error": err,
		}).Trace("Rabbitmq error")
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
		if key_id >= 0 {
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
}

func Subscribe() {
	conn, init_err = amqp.Dial("amqp://guest:" + password + "@localhost:5672/")
	failOnError(init_err, "Failed to connect to RabbitMQ")

	ch, init_err = conn.Channel()
	failOnError(init_err, "Failed to open a channel")

	log.Trace("Beginning rabbitmq initialisation")
	failOnError(init_err, "Rabbitmq error")
	if init_err == nil {
		var topics = [2]string{
			DATAINFO,
			REQUESTACCESS,
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
				s := string(d.Body[:])
				messages(d.RoutingKey, s)
				log.Debug("Checking states of received messages")
				checkState()
			}
			//This function is checked after to see if multiple errors occur then to
			//through an event message
		}()

		log.Trace(" [*] Waiting for logs. To exit press CTRL+C")
		<-forever
	}
}

func PublishEventNAC(component string, message string, time string) string {
	failure := ""

	eventNAC, err := json.Marshal(&EventNAC{
		Component:    component,
		Message:      message,
		Time:         time})
	if err != nil {
		failure = "Failed to convert EventNAC"
		log.Warn(failure)
	} else {
		if init_err == nil {
			log.Debug(string(eventNAC))
			err = ch.Publish(
				EXCHANGENAME, // exchange
				EVENTNAC,      // routing key
				false,        // mandatory
				false,        // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        []byte(eventNAC),
				})
			if err != nil {
				log.Fatal(err)
				failure = FAILUREPUBLISH
			}
		}
	}
	return failure
}

func PublishFailureNetwork(time string, reason string) string {
	failure := ""

	failureNetwork, err := json.Marshal(&FailureNetwork{
		Time:         time,
		Failure_type: reason})
	if err != nil {
		failure = "Failed to convert FailureNetwork"
		log.Warn(failure)
	} else {
		if init_err == nil {
			log.Debug(string(failureNetwork))
			err = ch.Publish(
				EXCHANGENAME, // exchange
				FAILURENETWORK,  // routing key
				false,        // mandatory
				false,        // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        []byte(failureNetwork),
				})
			if err != nil {
				log.Fatal(err)
				failure = FAILUREPUBLISH
			}
		}
	}
	return failure
}

func PublishRequestDatabase(id int, time_from string, time_to string, message string) string {
	failure := ""

	request, err := json.Marshal(&RequestDatabase{
		Request_id: id,
		Time_from:  time_from,
		Time_to:    time_to,
		Type:       message})
	if err != nil {
		failure = "Failed to convert RequestDatabase"
		log.Warn(failure)
	} else {
		if init_err == nil {
			log.Debug(string(request))
			err = ch.Publish(
				EXCHANGENAME, // exchange
				REQUESTDATABASE,  // routing key
				false,        // mandatory
				false,        // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        []byte(request),
				})
			if err != nil {
				log.Fatal(err)
				failure = FAILUREPUBLISH
			}
		}
	}
	return failure
}

func PublishDeviceFound(name string, address string, mac string) string {
	failure := ""

	device, err := json.Marshal(&DeviceFound{
		Device_name: name,
		Ip_address:  address,
		Mac: mac})
	if err != nil {
		failure = "Failed to convert DeviceFound"
		log.Warn(failure)
	} else {
		if init_err == nil {
			log.Debug(string(device))
			err = ch.Publish(
				EXCHANGENAME, // exchange
				DEVICEFOUND,  // routing key
				false,        // mandatory
				false,        // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        []byte(device),
				})
			if err != nil {
				log.Fatal(err)
				failure = FAILUREPUBLISH
			}
		}
	}
	return failure
}

func PublishAccessResponse(id int, result string) string {
	failure := ""

	access, err := json.Marshal(&AccessResponse{
		Id:     id,
		Result: result})
	if err != nil {
		failure = "Failed to convert AccessResponse"
		log.Warn(failure)
	} else {
		if init_err == nil {
			err = ch.Publish(
				EXCHANGENAME, // exchange
				ACCESSRESPONSE,  // routing key
				false,        // mandatory
				false,        // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        []byte(access),
				})
			if err != nil {
				log.Fatal(err)
				failure = FAILUREPUBLISH
			}
		}
	}
	return failure
}

func PublishUnauthorisedConnection(mac string, time string, alive bool) string {
	failure := ""

	connection, err := json.Marshal(&UnauthorisedConnection{
		Mac:  mac,
		Time: time,
		Alive: alive})
	if err != nil {
		failure = "Failed to convert UnauthorisedConnection"
		log.Warn(failure)
	} else {
		if init_err == nil {
			err = ch.Publish(
				EXCHANGENAME, // exchange
				UNAUTHORISEDCONNECTION,  // routing key
				false,        // mandatory
				false,        // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        []byte(connection),
				})
			if err != nil {
				log.Fatal(err)
				failure = FAILUREPUBLISH
			}
		}
	}
	return failure
}