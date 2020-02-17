package broker

import (
	"encoding/json"
	"fmt"
	"g-tech.com/dto"
	"g-tech.com/infrastructure"
	"g-tech.com/infrastructure/logger"
	"github.com/streadway/amqp"
)

/**
 * Connects to RabbitMQ server
 */
func Connect(host string, port int, userName string, password string) (conn *amqp.Connection) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", userName, password, host, port)

	// Make a infrastructure
	conn, err := amqp.Dial(url)

	if err != nil {
		logger.Error("Failed to connect to RabbitMQ %s", err.Error())
		return nil
	}

	return conn
}

/**
 * Creates a RabbitMQ queue
 */
func CreateQueue(ch *amqp.Channel, name string) amqp.Queue {
	queue, err := ch.QueueDeclare(
		name,  				// name of the queue
		true, 		// should the message be persistent? also queue will survive if the cluster gets reset
		false, 	// auto delete if there's no consumers (like queues that have anonymous names, often used with fanout exchange)
		false, 	// exclusive means I should get an error if any other consumer subscribes to this queue
		false, 		// no-wait means I don't want RabbitMQ to wait if there's a queue successfully setup
		nil,   		// arguments for more advanced configuration
	)

	if err != nil {
		logger.Error("Failed to declare a queue %s: %s", name, err.Error())
		//return nil
	}

	return queue
}

/**
 * Pushes a message to RabbitMQ queue
 * @Return {bool}
 */
func PushMessage(rbChannel *amqp.Channel, rbExchange string, rbRouteKey string, messageType string, obj interface{}) error {
	objTask := dto.Task {
		MessageType: messageType,
		Data:     obj,
	}
	//if (messageType == "ResultPost"){
	//	fmt.Println(util.ToJSON(objTask))
	//}
	task, err := json.Marshal(objTask)

	if err != nil {
		logger.Error("Failed to parse message %s", err.Error())
		return err
	}

	err = rbChannel.Publish (
		rbExchange,    			// exchange
		rbRouteKey, 			// routing key
		false,       	// mandatory
		false,       	// immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: infrastructure.DefaultContentType,
			Body:        task,
		})

	if err != nil {
		logger.Error("Failed to publish a message %s", err.Error())
		return err
	}

	return nil
}

/**
 * Handles error
 */
func HandleError(err error)  {
	logger.Error("[RabbitMQ]", err.Error())
}