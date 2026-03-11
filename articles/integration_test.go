//go:build integration

package articles

import (
	"fmt"
	"os"
	"testing"

	"github.com/gothinkster/golang-gin-realworld-example-app/common"
	"github.com/gothinkster/golang-gin-realworld-example-app/users"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupArticleIntegrationDB(t *testing.T) *gorm.DB {
	t.Helper()
	origDB := common.DB
	os.Setenv("TEST_DB_PATH", fmt.Sprintf("./data/integration_%s.db", t.Name()))
	db := common.TestDBInit()
	users.AutoMigrate()
	db.AutoMigrate(&ArticleModel{})
	db.AutoMigrate(&TagModel{})
	db.AutoMigrate(&FavoriteModel{})
	db.AutoMigrate(&ArticleUserModel{})
	db.AutoMigrate(&CommentModel{})
	t.Cleanup(func() {
		common.TestDBFree(db)
		common.DB = origDB
	})
	return db
}

func createIntTestUser(t *testing.T, db *gorm.DB) users.UserModel {
	t.Helper()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to generate password hash: %v", err)
	}
	user := users.UserModel{
		Username:     fmt.Sprintf("testuser%d", common.RandInt()),
		Email:        fmt.Sprintf("test%d@example.com", common.RandInt()),
		Bio:          "test bio",
		PasswordHash: string(passwordHash),
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

func createIntTestArticle(t *testing.T, author ArticleUserModel, title, slug string, tags []string) ArticleModel {
	t.Helper()
	article := ArticleModel{
		Slug:        slug,
		Title:       title,
		Description: "Test Description",
		Body:        "Test Body Content",
		AuthorID:    author.ID,
		Author:      author,
	}
	if len(tags) > 0 {
		if err := article.setTags(tags); err != nil {
			t.Fatalf("Failed to set tags: %v", err)
		}
	}
	if err := SaveOne(&article); err != nil {
		t.Fatalf("Failed to save article: %v", err)
	}
	return article
}

func TestIntegration_ArticleCRUD(t *testing.T) {
	asserts := assert.New(t)
	db := setupArticleIntegrationDB(t)
	_ = db

	user := createIntTestUser(t, db)
	articleUser := GetArticleUserModel(user)

	slug := fmt.Sprintf("test-article-%d", common.RandInt())
	title := "Test Article Title"

	// Create
	article := createIntTestArticle(t, articleUser, title, slug, nil)
	asserts.NotZero(article.ID, "Article should have an ID after creation")
	asserts.Equal(slug, article.Slug, "Slug should match")
	asserts.Equal(title, article.Title, "Title should match")

	// Read
	found, err := FindOneArticle(&ArticleModel{Slug: slug})
	asserts.NoError(err, "Should find article by slug")
	asserts.Equal(article.ID, found.ID, "Found article ID should match")
	asserts.Equal(title, found.Title, "Found title should match")
	asserts.Equal("Test Description", found.Description, "Description should match")
	asserts.Equal("Test Body Content", found.Body, "Body should match")
	asserts.Equal(articleUser.ID, found.AuthorID, "AuthorID should match")
	asserts.Equal(user.Username, found.Author.UserModel.Username, "Author username should be preloaded")

	// Update
	newTitle := "Updated Article Title"
	err = found.Update(ArticleModel{Title: newTitle})
	asserts.NoError(err, "Update should not return error")

	updated, err := FindOneArticle(&ArticleModel{Model: gorm.Model{ID: found.ID}})
	asserts.NoError(err, "Should find updated article")
	asserts.Equal(newTitle, updated.Title, "Title should be updated")

	// Delete
	err = DeleteArticleModel(&ArticleModel{Model: gorm.Model{ID: article.ID}})
	asserts.NoError(err, "Delete should not return error")

	_, err = FindOneArticle(&ArticleModel{Model: gorm.Model{ID: article.ID}})
	asserts.Error(err, "Should not find deleted article")
}

func TestIntegration_TagManagement(t *testing.T) {
	asserts := assert.New(t)
	db := setupArticleIntegrationDB(t)

	user := createIntTestUser(t, db)
	articleUser := GetArticleUserModel(user)

	sharedTag := fmt.Sprintf("shared-%d", common.RandInt())
	uniqueTag1 := fmt.Sprintf("unique1-%d", common.RandInt())
	uniqueTag2 := fmt.Sprintf("unique2-%d", common.RandInt())

	// Create first article with tags
	article1 := createIntTestArticle(t, articleUser,
		"Tag Article 1", fmt.Sprintf("tag-article-1-%d", common.RandInt()),
		[]string{sharedTag, uniqueTag1},
	)
	asserts.Equal(2, len(article1.Tags), "Article 1 should have 2 tags")

	// Create second article with overlapping tag
	article2 := createIntTestArticle(t, articleUser,
		"Tag Article 2", fmt.Sprintf("tag-article-2-%d", common.RandInt()),
		[]string{sharedTag, uniqueTag2},
	)
	asserts.Equal(2, len(article2.Tags), "Article 2 should have 2 tags")

	// Verify tag deduplication: the shared tag should have the same ID in both articles
	var sharedTagID1, sharedTagID2 uint
	for _, tag := range article1.Tags {
		if tag.Tag == sharedTag {
			sharedTagID1 = tag.ID
		}
	}
	for _, tag := range article2.Tags {
		if tag.Tag == sharedTag {
			sharedTagID2 = tag.ID
		}
	}
	asserts.NotZero(sharedTagID1, "Shared tag should exist in article 1")
	asserts.Equal(sharedTagID1, sharedTagID2, "Shared tag should be deduplicated (same ID)")

	// Verify getAllTags returns the tags
	allTags, err := getAllTags()
	asserts.NoError(err, "getAllTags should not error")
	tagNames := make(map[string]bool)
	for _, tag := range allTags {
		tagNames[tag.Tag] = true
	}
	asserts.True(tagNames[sharedTag], "getAllTags should contain the shared tag")
	asserts.True(tagNames[uniqueTag1], "getAllTags should contain unique tag 1")
	asserts.True(tagNames[uniqueTag2], "getAllTags should contain unique tag 2")

	// Verify setTags with empty list
	emptyArticle := ArticleModel{}
	err = emptyArticle.setTags([]string{})
	asserts.NoError(err, "setTags with empty list should not error")
	asserts.Equal(0, len(emptyArticle.Tags), "Tags should be empty")
}

func TestIntegration_FavoriteOperations(t *testing.T) {
	asserts := assert.New(t)
	db := setupArticleIntegrationDB(t)

	user1 := createIntTestUser(t, db)
	user2 := createIntTestUser(t, db)
	articleUser1 := GetArticleUserModel(user1)
	articleUser2 := GetArticleUserModel(user2)

	article := createIntTestArticle(t, articleUser1,
		"Favorite Test", fmt.Sprintf("fav-test-%d", common.RandInt()), nil)

	// Initially no favorites
	asserts.Equal(uint(0), article.favoritesCount(), "Should have 0 favorites initially")
	asserts.False(article.isFavoriteBy(articleUser1), "Should not be favorited by user1 initially")
	asserts.False(article.isFavoriteBy(articleUser2), "Should not be favorited by user2 initially")

	// User1 favorites the article
	err := article.favoriteBy(articleUser1)
	asserts.NoError(err, "favoriteBy should not error")
	asserts.Equal(uint(1), article.favoritesCount(), "Should have 1 favorite after user1 favorites")
	asserts.True(article.isFavoriteBy(articleUser1), "Should be favorited by user1")
	asserts.False(article.isFavoriteBy(articleUser2), "Should not be favorited by user2")

	// User2 also favorites
	err = article.favoriteBy(articleUser2)
	asserts.NoError(err, "favoriteBy should not error for user2")
	asserts.Equal(uint(2), article.favoritesCount(), "Should have 2 favorites")

	// Idempotent favorite
	err = article.favoriteBy(articleUser1)
	asserts.NoError(err, "Duplicate favoriteBy should not error")
	asserts.Equal(uint(2), article.favoritesCount(), "Favorites count should not change on duplicate")

	// Test BatchGetFavoriteCounts
	article2 := createIntTestArticle(t, articleUser1,
		"Batch Fav Test", fmt.Sprintf("batch-fav-%d", common.RandInt()), nil)
	err = article2.favoriteBy(articleUser2)
	asserts.NoError(err)

	counts := BatchGetFavoriteCounts([]uint{article.ID, article2.ID})
	asserts.Equal(uint(2), counts[article.ID], "BatchGetFavoriteCounts should return 2 for article")
	asserts.Equal(uint(1), counts[article2.ID], "BatchGetFavoriteCounts should return 1 for article2")

	// Test BatchGetFavoriteStatus
	status := BatchGetFavoriteStatus([]uint{article.ID, article2.ID}, articleUser1.ID)
	asserts.True(status[article.ID], "User1 should have favorited article")
	asserts.False(status[article2.ID], "User1 should not have favorited article2")

	status2 := BatchGetFavoriteStatus([]uint{article.ID, article2.ID}, articleUser2.ID)
	asserts.True(status2[article.ID], "User2 should have favorited article")
	asserts.True(status2[article2.ID], "User2 should have favorited article2")

	// Unfavorite
	err = article.unFavoriteBy(articleUser1)
	asserts.NoError(err, "unFavoriteBy should not error")
	asserts.Equal(uint(1), article.favoritesCount(), "Should have 1 favorite after unfavorite")
	asserts.False(article.isFavoriteBy(articleUser1), "Should not be favorited by user1 after unfavorite")

	// Idempotent unfavorite
	err = article.unFavoriteBy(articleUser1)
	asserts.NoError(err, "Duplicate unFavoriteBy should not error")
}

func TestIntegration_CommentOperations(t *testing.T) {
	asserts := assert.New(t)
	db := setupArticleIntegrationDB(t)

	user := createIntTestUser(t, db)
	articleUser := GetArticleUserModel(user)

	article := createIntTestArticle(t, articleUser,
		"Comment Test", fmt.Sprintf("comment-test-%d", common.RandInt()), nil)

	// Create comments
	comment1 := CommentModel{
		ArticleID: article.ID,
		AuthorID:  articleUser.ID,
		Body:      "First comment",
	}
	err := db.Create(&comment1).Error
	asserts.NoError(err, "Creating comment 1 should not error")
	asserts.NotZero(comment1.ID, "Comment should have an ID")

	comment2 := CommentModel{
		ArticleID: article.ID,
		AuthorID:  articleUser.ID,
		Body:      "Second comment",
	}
	err = db.Create(&comment2).Error
	asserts.NoError(err, "Creating comment 2 should not error")

	// FindOneComment
	found, err := FindOneComment(&CommentModel{Model: gorm.Model{ID: comment1.ID}})
	asserts.NoError(err, "FindOneComment should not error")
	asserts.Equal("First comment", found.Body, "Comment body should match")
	asserts.Equal(article.ID, found.ArticleID, "Comment article ID should match")
	asserts.Equal(user.Username, found.Author.UserModel.Username, "Comment author should be preloaded")

	// getComments
	err = article.getComments()
	asserts.NoError(err, "getComments should not error")
	asserts.Equal(2, len(article.Comments), "Article should have 2 comments")

	// DeleteCommentModel
	err = DeleteCommentModel(&CommentModel{Model: gorm.Model{ID: comment1.ID}})
	asserts.NoError(err, "DeleteCommentModel should not error")

	// Verify comment is deleted (soft delete)
	_, err = FindOneComment(&CommentModel{Model: gorm.Model{ID: comment1.ID}})
	asserts.Error(err, "Deleted comment should not be found")

	// Remaining comment should still exist
	found2, err := FindOneComment(&CommentModel{Model: gorm.Model{ID: comment2.ID}})
	asserts.NoError(err, "Second comment should still exist")
	asserts.Equal("Second comment", found2.Body, "Second comment body should match")
}

func TestIntegration_FindManyArticle(t *testing.T) {
	asserts := assert.New(t)
	db := setupArticleIntegrationDB(t)

	// Create users and articles with various tags
	user1 := createIntTestUser(t, db)
	user2 := createIntTestUser(t, db)
	articleUser1 := GetArticleUserModel(user1)
	articleUser2 := GetArticleUserModel(user2)

	tagName := fmt.Sprintf("findmany-%d", common.RandInt())

	article1 := createIntTestArticle(t, articleUser1,
		"FindMany 1", fmt.Sprintf("findmany-1-%d", common.RandInt()),
		[]string{tagName})
	_ = createIntTestArticle(t, articleUser2,
		"FindMany 2", fmt.Sprintf("findmany-2-%d", common.RandInt()),
		[]string{tagName})
	_ = createIntTestArticle(t, articleUser1,
		"FindMany 3", fmt.Sprintf("findmany-3-%d", common.RandInt()),
		nil)

	// Query by tag
	articles, count, err := FindManyArticle(tagName, "", "20", "0", "")
	asserts.NoError(err, "FindManyArticle by tag should not error")
	asserts.Equal(2, count, "Should find 2 articles with the tag")
	asserts.Equal(2, len(articles), "Should return 2 articles")

	// Query by author
	_, count, err = FindManyArticle("", user1.Username, "20", "0", "")
	asserts.NoError(err, "FindManyArticle by author should not error")
	asserts.GreaterOrEqual(count, 2, "User1 should have at least 2 articles")

	// Query by author - user2
	_, count, err = FindManyArticle("", user2.Username, "20", "0", "")
	asserts.NoError(err, "FindManyArticle by user2 should not error")
	asserts.GreaterOrEqual(count, 1, "User2 should have at least 1 article")

	// Query by favorited
	err = article1.favoriteBy(articleUser2)
	asserts.NoError(err)
	_, count, err = FindManyArticle("", "", "20", "0", user2.Username)
	asserts.NoError(err, "FindManyArticle by favorited should not error")
	asserts.GreaterOrEqual(count, 1, "User2 should have at least 1 favorited article")

	// Default query (no filters)
	_, count, err = FindManyArticle("", "", "20", "0", "")
	asserts.NoError(err, "FindManyArticle default should not error")
	asserts.GreaterOrEqual(count, 3, "Should find at least 3 articles total")

	// Test pagination with limit
	limitArticles, _, err := FindManyArticle("", "", "1", "0", "")
	asserts.NoError(err, "FindManyArticle with limit should not error")
	asserts.LessOrEqual(len(limitArticles), 1, "Should return at most 1 article with limit=1")

	// Test pagination with offset
	_, _, err = FindManyArticle("", "", "20", "1", "")
	asserts.NoError(err, "FindManyArticle with offset should not error")

	// Test with invalid limit/offset (should use defaults)
	_, _, err = FindManyArticle("", "", "invalid", "invalid", "")
	asserts.NoError(err, "FindManyArticle with invalid pagination should not error")

	// Test with non-existent tag
	_, count, err = FindManyArticle("nonexistent-tag-xyz", "", "20", "0", "")
	asserts.NoError(err, "FindManyArticle with non-existent tag should not error")
	asserts.Equal(0, count, "Should find 0 articles with non-existent tag")

	// Test with non-existent author
	_, count, err = FindManyArticle("", "nonexistent-author-xyz", "20", "0", "")
	asserts.NoError(err, "FindManyArticle with non-existent author should not error")
	asserts.Equal(0, count, "Should find 0 articles with non-existent author")
}

func TestIntegration_ArticleFeed(t *testing.T) {
	asserts := assert.New(t)
	db := setupArticleIntegrationDB(t)

	// Create users
	follower := createIntTestUser(t, db)
	followed1 := createIntTestUser(t, db)
	followed2 := createIntTestUser(t, db)
	notFollowed := createIntTestUser(t, db)

	followerArticleUser := GetArticleUserModel(follower)
	followed1ArticleUser := GetArticleUserModel(followed1)
	followed2ArticleUser := GetArticleUserModel(followed2)
	notFollowedArticleUser := GetArticleUserModel(notFollowed)

	// Create follow relationships directly via DB (following is unexported in users package)
	err := db.Create(&users.FollowModel{
		FollowingID:  followed1.ID,
		FollowedByID: follower.ID,
	}).Error
	asserts.NoError(err)
	err = db.Create(&users.FollowModel{
		FollowingID:  followed2.ID,
		FollowedByID: follower.ID,
	}).Error
	asserts.NoError(err)

	// Create articles
	_ = createIntTestArticle(t, followed1ArticleUser,
		"Followed1 Article", fmt.Sprintf("followed1-%d", common.RandInt()), nil)
	_ = createIntTestArticle(t, followed2ArticleUser,
		"Followed2 Article", fmt.Sprintf("followed2-%d", common.RandInt()), nil)
	_ = createIntTestArticle(t, notFollowedArticleUser,
		"NotFollowed Article", fmt.Sprintf("notfollowed-%d", common.RandInt()), nil)

	// Get feed for follower
	articles, count, err := followerArticleUser.GetArticleFeed("20", "0")
	asserts.NoError(err, "GetArticleFeed should not error")
	asserts.GreaterOrEqual(count, 2, "Feed should contain at least 2 articles from followed users")

	// Verify all feed articles are from followed users
	for _, article := range articles {
		authorUserModelID := article.Author.UserModelID
		isFromFollowed := authorUserModelID == followed1.ID || authorUserModelID == followed2.ID
		asserts.True(isFromFollowed, "Feed articles should only be from followed users, got author user_model_id: %d", authorUserModelID)
	}

	// Test feed pagination
	articles, _, err = followerArticleUser.GetArticleFeed("1", "0")
	asserts.NoError(err, "GetArticleFeed with limit should not error")
	asserts.LessOrEqual(len(articles), 1, "Feed with limit=1 should return at most 1 article")

	// Test feed for user with no followings
	notFollowedArticleUser2 := GetArticleUserModel(notFollowed)
	articles, count, err = notFollowedArticleUser2.GetArticleFeed("20", "0")
	asserts.NoError(err, "GetArticleFeed for user with no followings should not error")
	asserts.Equal(0, count, "Feed should be empty for user with no followings")
	asserts.Equal(0, len(articles), "Feed articles should be empty")
}

func TestIntegration_GetArticleUserModel(t *testing.T) {
	asserts := assert.New(t)
	db := setupArticleIntegrationDB(t)

	user := createIntTestUser(t, db)

	// First call should create
	articleUser1 := GetArticleUserModel(user)
	asserts.NotZero(articleUser1.ID, "ArticleUserModel should have an ID")
	asserts.Equal(user.ID, articleUser1.UserModelID, "UserModelID should match")
	asserts.Equal(user.Username, articleUser1.UserModel.Username, "UserModel should be populated")

	// Second call should return the same (idempotent via FirstOrCreate)
	articleUser2 := GetArticleUserModel(user)
	asserts.Equal(articleUser1.ID, articleUser2.ID, "GetArticleUserModel should be idempotent")

	// Zero user should return empty model
	emptyArticleUser := GetArticleUserModel(users.UserModel{})
	asserts.Zero(emptyArticleUser.ID, "Zero user should return empty ArticleUserModel")
}

func TestIntegration_BatchOperationsWithEmptyInput(t *testing.T) {
	asserts := assert.New(t)
	_ = setupArticleIntegrationDB(t)

	// BatchGetFavoriteCounts with empty input
	counts := BatchGetFavoriteCounts([]uint{})
	asserts.NotNil(counts, "BatchGetFavoriteCounts with empty input should return non-nil map")
	asserts.Equal(0, len(counts), "BatchGetFavoriteCounts with empty input should return empty map")

	// BatchGetFavoriteStatus with empty article IDs
	status := BatchGetFavoriteStatus([]uint{}, 1)
	asserts.NotNil(status, "BatchGetFavoriteStatus with empty article IDs should return non-nil map")
	asserts.Equal(0, len(status), "BatchGetFavoriteStatus with empty article IDs should return empty map")

	// BatchGetFavoriteStatus with zero user ID
	status = BatchGetFavoriteStatus([]uint{1, 2, 3}, 0)
	asserts.NotNil(status, "BatchGetFavoriteStatus with zero user ID should return non-nil map")
	asserts.Equal(0, len(status), "BatchGetFavoriteStatus with zero user ID should return empty map")

	// BatchGetFavoriteCounts with non-existent article IDs
	counts = BatchGetFavoriteCounts([]uint{999999, 999998})
	asserts.NotNil(counts, "BatchGetFavoriteCounts with non-existent IDs should return non-nil map")
	asserts.Equal(uint(0), counts[999999], "Non-existent article should have 0 favorites")

	// BatchGetFavoriteStatus with non-existent IDs
	status = BatchGetFavoriteStatus([]uint{999999}, 999999)
	asserts.NotNil(status, "BatchGetFavoriteStatus with non-existent IDs should return non-nil map")
	asserts.False(status[999999], "Non-existent article should not be favorited")
}
