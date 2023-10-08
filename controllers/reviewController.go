package controllers

import (
	"context"
	"log"
	"net/http"
	"review-app/database"
	helper "review-app/helpers"
	"review-app/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var reviewCollection *mongo.Collection = database.OpenCollection(database.Client, "review")

// Add  new review
func AddReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Logged in account must be of the type USER
		if err := helper.VerifyUserType(c, "USER"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var review models.Reviews
		defer cancel()

		if err := c.BindJSON(&review); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationError := validate.Struct(&review); validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": validationError.Error()}})
			return
		}

		newReview := models.Reviews{
			Id:          primitive.NewObjectID(),
			Movie_id:    review.Movie_id,
			Reviewer_id: review.Reviewer_id,
			Review:      review.Review,
		}

		result, err := reviewCollection.InsertOne(ctx, newReview)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"Status":  http.StatusCreated,
			"Message": "success",
			"Data":    map[string]interface{}{"data": result}})
	}
}

// Filter reviews by movie_id
func ViewAMovieReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchreviews []models.Reviews
		queryParam := c.Query("movie_id")
		if queryParam == "" {
			log.Println("No movie id passed")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		searchquerydb, err := reviewCollection.Find(ctx, bson.M{"movie_id": bson.M{"$regex": queryParam}})
		if err != nil {
			c.IndentedJSON(404, "something went wrong in fetching the dbquery")
			return
		}
		err = searchquerydb.All(ctx, &searchreviews)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer searchquerydb.Close(ctx)
		if err := searchquerydb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}
		defer cancel()
		c.IndentedJSON(200, searchreviews)
	}
}

// Delete a review
func DeleteReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		reviewId := c.Param("_id")
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(reviewId)

		result, err := reviewCollection.DeleteOne(ctx, bson.M{"_id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				gin.H{
					" Status":  http.StatusNotFound,
					" Message": "error",
					" Data":    map[string]interface{}{"data": "Review with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			gin.H{
				"Status":  http.StatusOK,
				"Message": "success",
				"Data":    map[string]interface{}{"data": "Your review was successfully deleted!"}},
		)
	}
}

// Allow a user view all their Reviews
func AllUserReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchreviews []models.Reviews
		reviewId := c.Param("reviewer_id")
		if reviewId == "" {
			log.Println("No reviewer id passed")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		searchquerydb, err := reviewCollection.Find(ctx, bson.M{"reviewer_id": reviewId})
		if err != nil {
			c.IndentedJSON(404, "something went wrong in fetching the dbquery")
			return
		}
		err = searchquerydb.All(ctx, &searchreviews)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer searchquerydb.Close(ctx)
		if err := searchquerydb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}
		defer cancel()
		c.IndentedJSON(200, searchreviews)
	}
}

// Update a review
func EditReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var review models.Reviews
		defer cancel()

		if err := c.BindJSON(&review); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationError := validate.Struct(&review); validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": validationError.Error()}})
			return
		}

		// Write your code here
		update := bson.M{
			"review": "Good"}
		filterByID := bson.M{"review_id": 2}
		result, err := reviewCollection.UpdateOne(ctx, filterByID, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedReview models.Reviews
		if result.MatchedCount == 1 {
			err := reviewCollection.FindOne(ctx, filterByID).Decode(&updatedReview)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"Status":  http.StatusInternalServerError,
					"Message": "error",
					"Data":    map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"Status":  http.StatusOK,
			"Message": "Review updated successfully!",
			"Data":    updatedReview})
	}
}
