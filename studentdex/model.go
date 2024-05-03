package studentdex

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"grpc/student-api/pb"
)

type Student struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	Grade     string    `json:"grade"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func studentToResponse(p Student) pb.Student {
	/*types := make([]pb.Type, len(p.Types))
	for i, t := range p.Types {
		types[i] = pb.Type(pb.Type_value[t])
	}*/

	return pb.Student{
		Id:        p.ID,
		Name:      p.Name,
		Grade:     p.Grade,
		CreatedAt: timestamppb.New(p.CreatedAt),
		UpdatedAt: timestamppb.New(p.UpdatedAt),
	}
}
