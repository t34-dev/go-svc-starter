package common_imp

import (
	"github.com/t34-dev/go-svc-starter/pkg/api/common_v1"
	"log"
	"time"
)

func (s *ImplementedCommon) LongOperation(_ *common_v1.LongOperationRequest, stream common_v1.CommonV1_LongOperationServer) error {
	log.Println("Starting long operation")

	for i := 0; i < 3; i++ {
		select {
		case <-stream.Context().Done():
			log.Println("Client disconnected, stopping operation")
			return stream.Context().Err()
		default:
			log.Printf("Processing step %d/10\n", i+1)
			time.Sleep(time.Second)
			if err := stream.Send(&common_v1.LongOperationResponse{
				Status:   common_v1.LongOperationResponse_IN_PROGRESS,
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
	return stream.Send(&common_v1.LongOperationResponse{
		Status:  common_v1.LongOperationResponse_COMPLETED,
		Message: "Operation completed",
		Result:  "Long operation result",
	})
}
