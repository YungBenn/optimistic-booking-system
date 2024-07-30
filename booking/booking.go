package booking

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Booking struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     string             `bson:"user_id"`
	ResourceID string             `bson:"resource_id"`
	StartTime  time.Time          `bson:"start_time"`
	EndTime    time.Time          `bson:"end_time"`
	Version    int                `bson:"version"`
}

func CreateBooking(ctx context.Context, collection *mongo.Collection, booking Booking) (primitive.ObjectID, error) {
	// Check for existing booking with the same user, resource, and overlapping time
	filter := bson.M{
		"user_id":     booking.UserID,
		"resource_id": booking.ResourceID,
		"$or": []bson.M{
			{"start_time": bson.M{"$lt": booking.EndTime}, "end_time": bson.M{"$gt": booking.StartTime}},
		},
	}
	var existingBooking Booking
	err := collection.FindOne(ctx, filter).Decode(&existingBooking)
	if err != mongo.ErrNoDocuments {
		return primitive.NilObjectID, fmt.Errorf("conflict: overlapping booking exists")
	}

	// Set initial version to 1 for new bookings
	booking.Version = 1

	result, err := collection.InsertOne(ctx, booking)
	if err != nil {
		return primitive.NilObjectID, err
	}

	// Return the ID of the newly created booking
	return result.InsertedID.(primitive.ObjectID), nil
}

func UpdateBookingWithOptimisticLock(ctx context.Context, collection *mongo.Collection, id string, newStartTime, newEndTime time.Time, currentVersion int) error {
    // Define the filter to include the current version
	objID, _ := primitive.ObjectIDFromHex(id)
    filter := bson.M{"_id": objID, "version": currentVersion}

    // Define the update to set the new times and increment the version
    update := bson.M{
        "$set": bson.M{
            "start_time": newStartTime,
            "end_time":   newEndTime,
        },
        "$inc": bson.M{"version": 1},
    }

    // Perform the update
    result, err := collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }

    // Check if a document was modified (indicating the update was successful)
    if result.ModifiedCount == 0 {
        return fmt.Errorf("booking update failed due to version conflict")
    }

    return nil
}
