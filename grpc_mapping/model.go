package grpc_mapping

const (
	InfoKey = "metadata_info_key"
)

type Info struct {
	ID            string `bson:"_id,omitempty"`
	Fullname      string `bson:"fullname,omitempty"`
	Password      string `bson:"password,omitempty"`
	Phone         string `bson:"phone,omitempty"`
	Email         string `bson:"email,omiempty"`
	Age           int32  `bson:"age,omitempty"`
	Role          int    `bson:"role,omitempty"`
	Sex           int    `bson:"sex,omitempty"`
	ProvinceID    int32  `bson:"province_id,omitempty"`
	DistrictID    int32  `bson:"district_id,omitempty"`
	WardID        int32  `bson:"ward_id,omitempty"`
	LocationScore int64  `bson:"location_score,omitempty"`
	IdentityCard  string `bson:"identity_card,omitempty"`
	Job           string `bson:"job,omitempty"`
}
