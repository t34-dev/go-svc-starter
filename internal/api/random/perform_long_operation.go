package random

import (
	"github.com/t34-dev/go-svc-starter/pkg/api/random_v1"
	"log"
	"time"
)

func (s *ImplementedRandom) PerformLongOperation(_ *random_v1.LongOperationRequest, stream random_v1.RandomService_PerformLongOperationServer) error {
	log.Println("Starting long operation")

	for i := 0; i < 3; i++ {
		select {
		case <-stream.Context().Done():
			log.Println("Client disconnected, stopping operation")
			return stream.Context().Err()
		default:
			log.Printf("Processing step %d/10\n", i+1)
			time.Sleep(time.Second)
			if err := stream.Send(&random_v1.LongOperationResponse{
				Status:   random_v1.LongOperationResponse_IN_PROGRESS,
				Message:  "Processing...",
				Progress: int32((i + 1) * 10),
			}); err != nil {
				log.Printf("Error sending update to client: %v\n", err)
				return err
			}
			log.Printf("Sent progress update: %d%%\n", (i+1)*10)
		}
	}

	log.Println("Long operation completed successfully")
	return stream.Send(&random_v1.LongOperationResponse{
		Status:  random_v1.LongOperationResponse_COMPLETED,
		Message: "Operation completed",
		Result:  "Long operation result",
	})
}
