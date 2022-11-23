package messageQ

import (
	"context"
	"encoding/json"
	"github.com/ComfortDelgro/models"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// InitMessageQ 入参为只读string类型channel, int类型client ID号，string类型exchange Name
func InitMessageQ(c <-chan models.SMS, id int, exName string) {
	conn, err := amqp.Dial("amqp://MjphbXFwLWNuLTJyNDJ3d3VxMDAwMzpMVEFJNXRRcTVjWG1YSEg0OG1YUjcxeWk=:MDJDM0YwQjg2RUJGNkFBMzc2OTM0RDEzQjYyREJGNTlFRTE5OTQ5NzoxNjY1NzIyMTYyOTM3@amqp-cn-2r42wwuq0003.ap-southeast-1.amqp-0.net.mq.amqp.aliyuncs.com:5672/Vhost-CSD")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exName,   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		body := <-c
		m, err := json.Marshal(body)
		err = ch.PublishWithContext(ctx,
			exName, // exchange
			"",     // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				//Body:        []byte(body.Messages),
				Body: m,
			})
		failOnError(err, "Failed to publish message: "+body.Payload)

		log.Printf(" [x] Client%d push MQ %s", id, body.Payload)
	}
}

func InitMbQ(c <-chan models.MbRc, id int, exName string) {
	conn, err := amqp.Dial("amqp://MjphbXFwLWNuLTJyNDJ3d3VxMDAwMzpMVEFJNXRRcTVjWG1YSEg0OG1YUjcxeWk=:MDJDM0YwQjg2RUJGNkFBMzc2OTM0RDEzQjYyREJGNTlFRTE5OTQ5NzoxNjY1NzIyMTYyOTM3@amqp-cn-2r42wwuq0003.ap-southeast-1.amqp-0.net.mq.amqp.aliyuncs.com:5672/Vhost-CSD")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exName,   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		body := <-c
		m, err := json.Marshal(body)
		err = ch.PublishWithContext(ctx,
			exName, // exchange
			"",     // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        m,
			})
		failOnError(err, "Failed to publish message: "+body.Payload)

		log.Printf(" [x] Client%d push MQ %s", id, body.Payload)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
