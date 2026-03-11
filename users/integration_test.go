//go:build integration

package users

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/gothinkster/golang-gin-realworld-example-app/common"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupIntegrationDB(t *testing.T) *gorm.DB {
	t.Helper()
	// Save the original global DB so unit tests' TestMain DB is not corrupted
	origDB := common.DB
	os.Setenv("TEST_DB_PATH", fmt.Sprintf("./data/integration_%s.db", t.Name()))
	db := common.TestDBInit()
	AutoMigrate()
	t.Cleanup(func() {
		common.TestDBFree(db)
		// Restore the original global DB for subsequent unit tests
		common.DB = origDB
	})
	return db
}

func createIntegrationTestUser(t *testing.T, db *gorm.DB, username, email string) UserModel {
	t.Helper()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to generate password hash: %v", err)
	}
	user := UserModel{
		Username:     username,
		Email:        email,
		Bio:          "test bio",
		PasswordHash: string(passwordHash),
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

func TestIntegration_UserCRUD(t *testing.T) {
	asserts := assert.New(t)
	db := setupIntegrationDB(t)

	// Create
	username := fmt.Sprintf("cruduser%d", common.RandInt())
	email := fmt.Sprintf("crud%d@example.com", common.RandInt())
	user := createIntegrationTestUser(t, db, username, email)
	asserts.NotZero(user.ID, "User should have an ID after creation")
	asserts.Equal(username, user.Username, "Username should match")
	asserts.Equal(email, user.Email, "Email should match")

	// Read
	foundUser, err := FindOneUser(&UserModel{Email: email})
	asserts.NoError(err, "Should find user by email")
	asserts.Equal(user.ID, foundUser.ID, "Found user ID should match")
	asserts.Equal(username, foundUser.Username, "Found username should match")
	asserts.Equal(email, foundUser.Email, "Found email should match")
	asserts.Equal("test bio", foundUser.Bio, "Found bio should match")

	// Read by username
	foundByUsername, err := FindOneUser(&UserModel{Username: username})
	asserts.NoError(err, "Should find user by username")
	asserts.Equal(user.ID, foundByUsername.ID, "Found user ID should match when searching by username")

	// Update
	newBio := "updated bio"
	newImage := "https://example.com/new-image.png"
	err = foundUser.Update(UserModel{Bio: newBio, Image: &newImage})
	asserts.NoError(err, "Update should not return error")

	// Verify update
	updatedUser, err := FindOneUser(&UserModel{ID: user.ID})
	asserts.NoError(err, "Should find updated user")
	asserts.Equal(newBio, updatedUser.Bio, "Bio should be updated")
	asserts.NotNil(updatedUser.Image, "Image should not be nil after update")
	asserts.Equal(newImage, *updatedUser.Image, "Image should be updated")

	// Verify password is still valid after update
	asserts.NoError(updatedUser.checkPassword("password123"), "Password should still be valid after update")
}

func TestIntegration_UniqueEmailConstraint(t *testing.T) {
	asserts := assert.New(t)
	db := setupIntegrationDB(t)

	email := fmt.Sprintf("unique%d@example.com", common.RandInt())
	createIntegrationTestUser(t, db, fmt.Sprintf("user1_%d", common.RandInt()), email)

	// Try to create another user with the same email
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	duplicateUser := UserModel{
		Username:     fmt.Sprintf("user2_%d", common.RandInt()),
		Email:        email,
		Bio:          "another bio",
		PasswordHash: string(passwordHash),
	}
	err := db.Create(&duplicateUser).Error
	asserts.Error(err, "Creating user with duplicate email should fail")
	asserts.Contains(err.Error(), "UNIQUE constraint failed", "Error should mention unique constraint")
}

func TestIntegration_FollowUnfollowCycle(t *testing.T) {
	asserts := assert.New(t)
	db := setupIntegrationDB(t)

	// Create three users
	userA := createIntegrationTestUser(t, db, fmt.Sprintf("followerA%d", common.RandInt()), fmt.Sprintf("followerA%d@example.com", common.RandInt()))
	userB := createIntegrationTestUser(t, db, fmt.Sprintf("followeeB%d", common.RandInt()), fmt.Sprintf("followeeB%d@example.com", common.RandInt()))
	userC := createIntegrationTestUser(t, db, fmt.Sprintf("followeeC%d", common.RandInt()), fmt.Sprintf("followeeC%d@example.com", common.RandInt()))

	// Initially no followings
	asserts.False(userA.isFollowing(userB), "A should not be following B initially")
	asserts.False(userA.isFollowing(userC), "A should not be following C initially")
	asserts.Equal(0, len(userA.GetFollowings()), "A should have no followings initially")

	// A follows B
	err := userA.following(userB)
	asserts.NoError(err, "A following B should not error")
	asserts.True(userA.isFollowing(userB), "A should be following B after follow")
	asserts.False(userA.isFollowing(userC), "A should not be following C yet")
	asserts.Equal(1, len(userA.GetFollowings()), "A should have 1 following")

	// A follows C
	err = userA.following(userC)
	asserts.NoError(err, "A following C should not error")
	asserts.True(userA.isFollowing(userC), "A should be following C after follow")
	asserts.Equal(2, len(userA.GetFollowings()), "A should have 2 followings")

	// B follows C (verify independent follow relationships)
	err = userB.following(userC)
	asserts.NoError(err, "B following C should not error")
	asserts.True(userB.isFollowing(userC), "B should be following C")
	asserts.False(userB.isFollowing(userA), "B should not be following A")

	// A unfollows B
	err = userA.unFollowing(userB)
	asserts.NoError(err, "A unfollowing B should not error")
	asserts.False(userA.isFollowing(userB), "A should not be following B after unfollow")
	asserts.True(userA.isFollowing(userC), "A should still be following C")
	asserts.Equal(1, len(userA.GetFollowings()), "A should have 1 following after unfollow")

	// Verify B's followings are unaffected
	asserts.True(userB.isFollowing(userC), "B should still be following C")

	// A unfollows C
	err = userA.unFollowing(userC)
	asserts.NoError(err, "A unfollowing C should not error")
	asserts.Equal(0, len(userA.GetFollowings()), "A should have 0 followings after unfollowing all")

	// Idempotent unfollow (unfollowing someone not followed)
	err = userA.unFollowing(userB)
	asserts.NoError(err, "Unfollowing a non-followed user should not error")
}

func TestIntegration_ConcurrentFollows(t *testing.T) {
	asserts := assert.New(t)
	db := setupIntegrationDB(t)

	// Create target user
	target := createIntegrationTestUser(t, db, fmt.Sprintf("target%d", common.RandInt()), fmt.Sprintf("target%d@example.com", common.RandInt()))

	// Create multiple followers
	const numFollowers = 10
	followers := make([]UserModel, numFollowers)
	for i := 0; i < numFollowers; i++ {
		followers[i] = createIntegrationTestUser(t, db,
			fmt.Sprintf("cfollower%d_%d", i, common.RandInt()),
			fmt.Sprintf("cfollower%d_%d@example.com", i, common.RandInt()),
		)
	}

	// All followers follow target concurrently
	var wg sync.WaitGroup
	errCh := make(chan error, numFollowers)
	for i := 0; i < numFollowers; i++ {
		wg.Add(1)
		go func(follower UserModel) {
			defer wg.Done()
			if err := follower.following(target); err != nil {
				errCh <- err
			}
		}(followers[i])
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		asserts.NoError(err, "Concurrent follows should not produce errors")
	}

	// Verify all follow relationships exist
	for i, follower := range followers {
		asserts.True(follower.isFollowing(target), "Follower %d should be following target", i)
	}

	// Idempotent: follow again should not error
	err := followers[0].following(target)
	asserts.NoError(err, "Following again should not error (idempotent)")
}

func TestIntegration_FindOneUserEdgeCases(t *testing.T) {
	asserts := assert.New(t)
	db := setupIntegrationDB(t)

	// Non-existent user
	_, err := FindOneUser(&UserModel{Email: "nonexistent@example.com"})
	asserts.Error(err, "Finding non-existent user should return error")
	asserts.ErrorIs(err, gorm.ErrRecordNotFound, "Error should be ErrRecordNotFound")

	// Find by ID = 0 (should not match any user)
	_, err = FindOneUser(&UserModel{ID: 0})
	asserts.Error(err, "Finding user with ID=0 should return error")

	// Find by non-existent username
	_, err = FindOneUser(&UserModel{Username: "absolutelynonexistentuser999999"})
	asserts.Error(err, "Finding non-existent username should return error")

	// Create a user, then find by multiple matching conditions
	user := createIntegrationTestUser(t, db,
		fmt.Sprintf("edgeuser%d", common.RandInt()),
		fmt.Sprintf("edge%d@example.com", common.RandInt()),
	)
	found, err := FindOneUser(&UserModel{Username: user.Username, Email: user.Email})
	asserts.NoError(err, "Finding user by multiple matching conditions should work")
	asserts.Equal(user.ID, found.ID, "Found user should match")
}
