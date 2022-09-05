package dao

import (
	"auth-poc/svc/auth/adapter/grpc/pb"
	"auth-poc/svc/auth/application/dto"
	"auth-poc/svc/auth/application/entity"
	"auth-poc/svc/auth/constants"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (dao *Auth) GetStoredHash(phoneNumber string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DEFAULT_TIMEOUT)
	defer cancel()

	var result dto.TokenCollection

	coll := dao.DataStore.GetColl()
	err := coll.FindOne(ctx, bson.D{{Key: "phone_number", Value: phoneNumber}}).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil
		}
		return "", err
	}
	return result.Password, nil
}

func (dao *Auth) GetUserByPhoneNumber(phoneNumber string) (*entity.UserData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DEFAULT_TIMEOUT)
	defer cancel()
	userData, err := dao.Rpc.User.GetUserByPhoneNumber(ctx, &pb.GetUserByPhoneNumberArgs{
		PhoneNumber: phoneNumber,
	})
	if err != nil {
		return nil, err
	}
	return &entity.UserData{
		PhoneNumber: userData.PhoneNumber,
		UserID:      userData.UserID,
		UserType:    uint(userData.UserType),
	}, nil
}

func (dao *Auth) UpsertCredential(byteArr []byte, phoneNumber string) error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DEFAULT_TIMEOUT)
	defer cancel()
	data := &dto.TokenCollection{
		PhoneNumber: phoneNumber,
		Password:    string(byteArr),
	}
	filter := bson.D{{Key: "phone_number", Value: phoneNumber}}
	opts := options.Update().SetUpsert(true)
	coll := dao.DataStore.GetColl()

	update := bson.D{{Key: "$set", Value: data}}

	_, err := coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func (dao *Auth) RemoveCredential(phoneNumber string) error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DEFAULT_TIMEOUT)
	defer cancel()
	filter := bson.D{{Key: "phone_number", Value: phoneNumber}}
	coll := dao.DataStore.GetColl()

	_, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
