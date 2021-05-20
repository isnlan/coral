package migrate

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

type Migrate interface {
	Update() error
	Version() *Version
	Describe() string
}

type Version struct {
	ID    primitive.ObjectID `json:"id" bson:"_id"`      // id
	Major int                `json:"major" bson:"major"` // 主版本号
	Sub   int                `json:"sub" bson:"sub"`     // 子版本号
	Stage int                `json:"stage" bson:"stage"` // 阶段版本号
}

func LoadVersion(db *mongo.Database) (*Version, error) {
	coll := db.Collection("version")
	var v Version
	err := coll.FindOne(context.Background(), bson.D{}).Decode(&v)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			v := &Version{
				ID:    primitive.NewObjectID(),
				Major: 0,
				Sub:   0,
				Stage: 0,
			}
			_, err := coll.InsertOne(context.Background(), v)
			if err != nil {
				return nil, err
			}
			return v, nil
		}
		return nil, err
	}
	return &v, nil
}

func (v *Version) Value() int {
	return v.Stage + 100*v.Sub + 10000*v.Major
}

func (v *Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Sub, v.Stage)
}

func (v *Version) Update(version *Version) {
	v.Major = version.Major
	v.Sub = version.Sub
	v.Stage = version.Stage
}

func (v *Version) Save(db *mongo.Database) error {
	coll := db.Collection("version")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": v.ID}
	update := bson.D{
		{Key: "$set", Value: v},
	}
	_, err := coll.UpdateOne(ctx, filter, update)
	return err
}
