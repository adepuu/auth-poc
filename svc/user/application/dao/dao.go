package dao

import (
	"auth-poc/svc/user/adapter/grpc/pb"
	"auth-poc/svc/user/application/dto"
	"auth-poc/svc/user/application/entity"
	"auth-poc/svc/user/constants"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (dao *User) GetAllUser(limit, page int64) ([]*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DEFAULT_TIMEOUT)
	defer cancel()

	// Default pagination not present
	if limit == 0 {
		limit = 10
	}
	if page == 0 {
		page = 1
	}

	l := int64(limit)
	skip := int64(page*limit - limit)
	fOpt := options.FindOptions{Limit: &l, Skip: &skip}
	coll := dao.DataStore.GetColl()

	curr, err := coll.Find(ctx, bson.D{{}}, &fOpt)
	if err != nil {
		return nil, err
	}

	users := make([]*entity.User, 0)

	for curr.Next(ctx) {
		var usr entity.User
		if err := curr.Decode(&usr); err != nil {
			return nil, err
		}
		users = append(users, &usr)
	}

	return users, nil
}

func (dao *User) GetUserByKey(userID, phoneNumber string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DEFAULT_TIMEOUT)
	defer cancel()

	result := dto.UserMongoDoc{}

	filter := bson.D{{Key: "phone_number", Value: phoneNumber}}

	if len(userID) > 0 {
		objectId, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return nil, err
		}
		filter = bson.D{{Key: "_id", Value: objectId}}
	}

	err := dao.DataStore.GetColl().FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		Email:       result.Email,
		PhoneNumber: result.PhoneNumber,
		FullName:    result.FullName,
		UserID:      result.DocumentID.Hex(),
		UserType:    result.UserType,
	}, nil
}

func (dao *User) InsertUniqueUser(usr *entity.User) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DEFAULT_TIMEOUT)
	defer cancel()

	coll := dao.DataStore.GetColl()
	result, err := coll.InsertOne(ctx, usr)

	if err != nil {
		mwe := err.(mongo.WriteException)
		// check if _id or phone_number is unique/no duplicate
		// if there's duplicate value, will return nil for both value
		if mwe.HasErrorCode(11000) || mwe.HasErrorCode(11001) || mwe.HasErrorCode(12582) || mwe.HasErrorCodeWithMessage(16460, " E11000 ") {
			return nil, nil
		}
		return nil, err
	}

	usr.UserID = result.InsertedID.(primitive.ObjectID).Hex()
	return usr, nil
}

func (dao *User) UpsertUser(usr *entity.User) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DEFAULT_TIMEOUT)
	defer cancel()

	filter := bson.D{{Key: "phone_number", Value: usr.PhoneNumber}}

	if len(usr.UserID) > 0 {
		objectId, err := primitive.ObjectIDFromHex(usr.UserID)
		if err != nil {
			return nil, err
		}
		filter = bson.D{{Key: "_id", Value: objectId}}
	}

	opts := options.Update().SetUpsert(true)
	coll := dao.DataStore.GetColl()

	update := bson.D{{Key: "$set", Value: usr}}

	result, err := coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}

	if result.UpsertedID != nil {
		usr.UserID = result.UpsertedID.(primitive.ObjectID).Hex()
	}
	return usr, nil
}

func (dao *User) DeleteUser(id, phoneNumber string) error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DEFAULT_TIMEOUT)
	defer cancel()

	result, err := dao.Rpc.Auth.RemovePassword(ctx, &pb.RemovePasswordRequest{
		PhoneNumber: phoneNumber,
	})
	if err != nil {
		return err
	}

	if !result.Success {
		return errors.New("failed to store password")
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: objectId}}
	coll := dao.DataStore.GetColl()

	_, err = coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (dao *User) StorePassword(rawPassword string, phoneNumber string) error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DEFAULT_TIMEOUT)
	defer cancel()

	result, err := dao.Rpc.Auth.StorePassword(ctx, &pb.StorePasswordRequest{
		RawPassword: rawPassword,
		PhoneNumber: phoneNumber,
	})
	if err != nil {
		return err
	}

	if !result.Success {
		return errors.New("failed to store password")
	}

	return nil
}
