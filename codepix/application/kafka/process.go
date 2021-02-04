package kafka

import (
	"fmt"
	"os"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"

	appmodel "github.com/mateusjbarbosa/imersao-fullstack-fullcycle/codepix/application/model"

	"github.com/mateusjbarbosa/imersao-fullstack-fullcycle/codepix/application/factory"

	"github.com/jinzhu/gorm"
	"github.com/mateusjbarbosa/imersao-fullstack-fullcycle/codepix/application/usecase"
	"github.com/mateusjbarbosa/imersao-fullstack-fullcycle/codepix/domain/model"
)

type KafkaProcessor struct {
	Database     *gorm.DB
	Producer     *ckafka.Producer
	DeliveryChan chan ckafka.Event
}

func NewKafkaProcessor(database *gorm.DB, producer *ckafka.Producer, deliveryChan chan ckafka.Event) *KafkaProcessor {
	return &KafkaProcessor{
		Database:     database,
		Producer:     producer,
		DeliveryChan: deliveryChan,
	}
}

func (k *KafkaProcessor) Consume() {
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": os.Getenv("kafkaBootstrapServers"),
		"group.id":          os.Getenv("kafkaConsumerGroupId"),
		"auto.offset.reset": "earliest",
	}
	c, err := ckafka.NewConsumer(configMap)

	if err != nil {
		panic(err)
	}

	topics := []string{os.Getenv("kafkaTransactionTopic"), os.Getenv("kafkaTransactionConfirmationTopic")}
	c.SubscribeTopics(topics, nil)

	fmt.Println("kafka consumer has been started")

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Println(string(msg.Value))
			k.processMessage(msg)
		}
	}
}

func (k *KafkaProcessor) processMessage(msg *ckafka.Message) {
	transactionsTopic := os.Getenv("kafkaTransactionTopic")
	transactionConfirmationTopic := os.Getenv("kafkaTransactionConfirmationTopic")

	switch topic := *msg.TopicPartition.Topic; topic {
	case transactionsTopic:
		k.processTransaction(msg)
	case transactionConfirmationTopic:
		k.processTransactionConfirmation(msg)
	default:
		fmt.Println("not a valid topic", string(msg.Value))
	}
}

func (k *KafkaProcessor) processTransaction(msg *ckafka.Message) error {
	transaction := appmodel.NewTransaction()
	err := transaction.ParseJson(msg.Value)
	if err != nil {
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(k.Database)

	createdTransaction, err := transactionUseCase.Register(
		transaction.AccountID,
		transaction.Amount,
		transaction.PixKeyTo,
		transaction.PixKeyKindTo,
		transaction.Description,
		transaction.ID,
	)

	if err != nil {
		fmt.Println("error registering transaction", err)
		return err
	}

	topic := "bank" + createdTransaction.PixKeyTo.Account.Bank.Code
	transaction.ID = createdTransaction.ID
	transaction.Status = model.TransactionPending
	transactionJSON, err := transaction.ToJson()

	if err != nil {
		return err
	}

	err = Publish(string(transactionJSON), topic, k.Producer, k.DeliveryChan)

	if err != nil {
		return err
	}

	return nil
}

func (k *KafkaProcessor) processTransactionConfirmation(msg *ckafka.Message) error {
	transaction := appmodel.NewTransaction()
	err := transaction.ParseJson(msg.Value)

	if err != nil {
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(k.Database)

	if transaction.Status == model.TransactionConfirmed {
		err = k.confirmTransaction(transaction, transactionUseCase)

		if err != nil {
			return err
		}
	} else if transaction.Status == model.TransactionCompleted {
		_, err := transactionUseCase.Complete(transaction.ID)

		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (k *KafkaProcessor) confirmTransaction(transaction *appmodel.Transaction, transactionUseCase usecase.TransactionUseCase) error {
	confirmedTransaction, err := transactionUseCase.Confirm(transaction.ID)

	if err != nil {
		return err
	}

	topic := "bank" + confirmedTransaction.AccountFrom.Bank.Code
	transactionJSON, err := transaction.ToJson()

	if err != nil {
		return err
	}

	err = Publish(string(transactionJSON), topic, k.Producer, k.DeliveryChan)

	if err != nil {
		return err
	}

	return nil
}
