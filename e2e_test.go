//go:build e2e

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gothinkster/golang-gin-realworld-example-app/articles"
	"github.com/gothinkster/golang-gin-realworld-example-app/common"
	"github.com/gothinkster/golang-gin-realworld-example-app/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupE2EDB creates an isolated test DB for each E2E test to avoid cross-package conflicts.
func setupE2EDB(t *testing.T) *common.Database {
	t.Helper()
	os.Setenv("TEST_DB_PATH", fmt.Sprintf("./data/e2e_%s.db", t.Name()))
	db := common.TestDBInit()
	Migrate(db)
	t.Cleanup(func() {
		common.TestDBFree(db)
	})
	return &common.Database{DB: db}
}

// setupE2ERouter creates a full application router identical to the one in main().
func setupE2ERouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.RedirectTrailingSlash = false
	v1 := r.Group("/api")
	users.UsersRegister(v1.Group("/users"))
	v1.Use(users.AuthMiddleware(false))
	articles.ArticlesAnonymousRegister(v1.Group("/articles"))
	articles.TagsAnonymousRegister(v1.Group("/tags"))
	users.ProfileRetrieveRegister(v1.Group("/profiles"))
	v1.Use(users.AuthMiddleware(true))
	users.UserRegister(v1.Group("/user"))
	users.ProfileRegister(v1.Group("/profiles"))
	articles.ArticlesRegister(v1.Group("/articles"))
	return r
}

// e2eRequest performs an HTTP request against the given router and returns the recorder.
func e2eRequest(r *gin.Engine, method, url, body, token string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req, _ = http.NewRequest(method, url, bytes.NewBufferString(body))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Token "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// parseJSON unmarshals the response body into a generic map.
func parseJSON(body *bytes.Buffer) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(body.Bytes(), &result)
	return result
}

// extractUserToken extracts the token string from a user response.
func extractUserToken(data map[string]interface{}) string {
	u, ok := data["user"].(map[string]interface{})
	if !ok {
		return ""
	}
	token, _ := u["token"].(string)
	return token
}

// extractArticleSlug extracts the slug from an article response.
func extractArticleSlug(data map[string]interface{}) string {
	a, ok := data["article"].(map[string]interface{})
	if !ok {
		return ""
	}
	slug, _ := a["slug"].(string)
	return slug
}

// registerUser registers a new user and returns the parsed response and token.
func registerUser(r *gin.Engine, username, email, password string) (map[string]interface{}, string) {
	body := fmt.Sprintf(`{"user":{"username":"%s","email":"%s","password":"%s"}}`, username, email, password)
	w := e2eRequest(r, "POST", "/api/users", body, "")
	data := parseJSON(w.Body)
	return data, extractUserToken(data)
}

// TestE2E_UserJourney tests the complete user lifecycle:
// Register -> Login -> Get user -> Update user -> Get profile
func TestE2E_UserJourney(t *testing.T) {
	setupE2EDB(t)

	r := setupE2ERouter()
	assert := assert.New(t)
	require := require.New(t)

	// --- Register ---
	regBody := `{"user":{"username":"e2euser","email":"e2euser@test.com","password":"password123"}}`
	w := e2eRequest(r, "POST", "/api/users", regBody, "")
	assert.Equal(http.StatusCreated, w.Code, "register should return 201")

	regData := parseJSON(w.Body)
	regUser, ok := regData["user"].(map[string]interface{})
	require.True(ok, "response should contain user object")
	assert.Equal("e2euser", regUser["username"])
	assert.Equal("e2euser@test.com", regUser["email"])
	regToken, _ := regUser["token"].(string)
	require.NotEmpty(regToken, "register should return a token")

	// --- Login ---
	loginBody := `{"user":{"email":"e2euser@test.com","password":"password123"}}`
	w = e2eRequest(r, "POST", "/api/users/login", loginBody, "")
	assert.Equal(http.StatusOK, w.Code, "login should return 200")

	loginData := parseJSON(w.Body)
	loginUserObj, ok := loginData["user"].(map[string]interface{})
	require.True(ok, "login response should contain user object")
	assert.Equal("e2euser", loginUserObj["username"])
	loginToken, _ := loginUserObj["token"].(string)
	require.NotEmpty(loginToken, "login should return a token")

	// --- Get current user ---
	w = e2eRequest(r, "GET", "/api/user", "", loginToken)
	assert.Equal(http.StatusOK, w.Code, "get user should return 200")

	getUserData := parseJSON(w.Body)
	currentUser, ok := getUserData["user"].(map[string]interface{})
	require.True(ok, "get user response should contain user object")
	assert.Equal("e2euser", currentUser["username"])
	assert.Equal("e2euser@test.com", currentUser["email"])

	// --- Update user ---
	updateBody := `{"user":{"username":"e2eupdated","bio":"Updated bio","image":"http://example.com/img.png"}}`
	w = e2eRequest(r, "PUT", "/api/user", updateBody, loginToken)
	assert.Equal(http.StatusOK, w.Code, "update user should return 200")

	updateData := parseJSON(w.Body)
	updatedUser, ok := updateData["user"].(map[string]interface{})
	require.True(ok, "update response should contain user object")
	assert.Equal("e2eupdated", updatedUser["username"])
	assert.Equal("Updated bio", updatedUser["bio"])
	assert.Equal("http://example.com/img.png", updatedUser["image"])

	// --- Get profile of updated user ---
	w = e2eRequest(r, "GET", "/api/profiles/e2eupdated", "", loginToken)
	assert.Equal(http.StatusOK, w.Code, "get profile should return 200")

	profileData := parseJSON(w.Body)
	profile, ok := profileData["profile"].(map[string]interface{})
	require.True(ok, "profile response should contain profile object")
	assert.Equal("e2eupdated", profile["username"])
	assert.Equal("Updated bio", profile["bio"])
	assert.Equal(false, profile["following"], "should not be following self")
}

// TestE2E_ArticleLifecycle tests the full article CRUD:
// Register -> Create article with tags -> Get article -> Update -> List by author -> List by tag -> Delete
func TestE2E_ArticleLifecycle(t *testing.T) {
	setupE2EDB(t)

	r := setupE2ERouter()
	assert := assert.New(t)
	require := require.New(t)

	// --- Register author ---
	_, token := registerUser(r, "author1", "author1@test.com", "password123")
	require.NotEmpty(token)

	// --- Create article with tags ---
	createBody := `{"article":{"title":"Test Article Title","description":"Test description","body":"Test article body content","tagList":["golang","testing","e2e"]}}`
	w := e2eRequest(r, "POST", "/api/articles", createBody, token)
	assert.Equal(http.StatusCreated, w.Code, "create article should return 201")

	createData := parseJSON(w.Body)
	article, ok := createData["article"].(map[string]interface{})
	require.True(ok, "response should contain article object")
	assert.Equal("Test Article Title", article["title"])
	assert.Equal("Test description", article["description"])
	assert.Equal("Test article body content", article["body"])
	slug, _ := article["slug"].(string)
	require.NotEmpty(slug, "article should have a slug")

	tagList, ok := article["tagList"].([]interface{})
	require.True(ok, "article should have tagList")
	assert.Len(tagList, 3, "should have 3 tags")

	// Check author info
	authorObj, ok := article["author"].(map[string]interface{})
	require.True(ok, "article should have author")
	assert.Equal("author1", authorObj["username"])

	// --- Get single article ---
	w = e2eRequest(r, "GET", "/api/articles/"+slug, "", "")
	assert.Equal(http.StatusOK, w.Code, "get article should return 200")

	getArticleData := parseJSON(w.Body)
	gotArticle, ok := getArticleData["article"].(map[string]interface{})
	require.True(ok)
	assert.Equal("Test Article Title", gotArticle["title"])

	// --- Update article ---
	updateBody := `{"article":{"title":"Updated Title","body":"Updated body content"}}`
	w = e2eRequest(r, "PUT", "/api/articles/"+slug, updateBody, token)
	assert.Equal(http.StatusOK, w.Code, "update article should return 200")

	updateData := parseJSON(w.Body)
	updatedArticle, ok := updateData["article"].(map[string]interface{})
	require.True(ok)
	assert.Equal("Updated Title", updatedArticle["title"])
	assert.Equal("Updated body content", updatedArticle["body"])
	// Slug may change after title update
	updatedSlug, _ := updatedArticle["slug"].(string)
	require.NotEmpty(updatedSlug)

	// --- List articles by author ---
	w = e2eRequest(r, "GET", "/api/articles?author=author1", "", "")
	assert.Equal(http.StatusOK, w.Code, "list articles should return 200")

	listData := parseJSON(w.Body)
	articlesList, ok := listData["articles"].([]interface{})
	require.True(ok, "response should contain articles array")
	assert.GreaterOrEqual(len(articlesList), 1, "should have at least 1 article by author")

	articlesCount, ok := listData["articlesCount"].(float64)
	require.True(ok, "response should contain articlesCount")
	assert.GreaterOrEqual(articlesCount, float64(1))

	// --- List articles by tag ---
	w = e2eRequest(r, "GET", "/api/articles?tag=golang", "", "")
	assert.Equal(http.StatusOK, w.Code, "list by tag should return 200")

	tagListData := parseJSON(w.Body)
	tagArticles, ok := tagListData["articles"].([]interface{})
	require.True(ok)
	assert.GreaterOrEqual(len(tagArticles), 1, "should have at least 1 article with tag")

	// --- Delete article ---
	w = e2eRequest(r, "DELETE", "/api/articles/"+updatedSlug, "", token)
	assert.Equal(http.StatusOK, w.Code, "delete article should return 200")

	// Verify article is gone
	w = e2eRequest(r, "GET", "/api/articles/"+updatedSlug, "", "")
	assert.Equal(http.StatusNotFound, w.Code, "deleted article should return 404")
}

// TestE2E_SocialInteractions tests following/unfollowing and feed:
// Register 2 users -> Follow -> Get profile (following=true) -> Create article -> Feed -> Unfollow -> Feed empty
func TestE2E_SocialInteractions(t *testing.T) {
	setupE2EDB(t)

	r := setupE2ERouter()
	assert := assert.New(t)
	require := require.New(t)

	// --- Register two users ---
	_, token1 := registerUser(r, "socialuser1", "social1@test.com", "password123")
	require.NotEmpty(token1)
	_, token2 := registerUser(r, "socialuser2", "social2@test.com", "password123")
	require.NotEmpty(token2)

	// --- User1 follows User2 ---
	w := e2eRequest(r, "POST", "/api/profiles/socialuser2/follow", "", token1)
	assert.Equal(http.StatusOK, w.Code, "follow should return 200")

	followData := parseJSON(w.Body)
	followProfile, ok := followData["profile"].(map[string]interface{})
	require.True(ok)
	assert.Equal("socialuser2", followProfile["username"])
	assert.Equal(true, followProfile["following"], "following should be true after follow")

	// --- Get profile should show following=true ---
	w = e2eRequest(r, "GET", "/api/profiles/socialuser2", "", token1)
	assert.Equal(http.StatusOK, w.Code)

	profileData := parseJSON(w.Body)
	profile, ok := profileData["profile"].(map[string]interface{})
	require.True(ok)
	assert.Equal(true, profile["following"], "profile should show following=true")

	// --- User2 creates an article (should appear in User1's feed) ---
	articleBody := `{"article":{"title":"Feed Article","description":"For feed test","body":"Feed content","tagList":["feed"]}}`
	w = e2eRequest(r, "POST", "/api/articles", articleBody, token2)
	assert.Equal(http.StatusCreated, w.Code)

	// --- User1 checks feed ---
	w = e2eRequest(r, "GET", "/api/articles/feed", "", token1)
	assert.Equal(http.StatusOK, w.Code, "feed should return 200")

	feedData := parseJSON(w.Body)
	feedArticles, ok := feedData["articles"].([]interface{})
	require.True(ok)
	assert.GreaterOrEqual(len(feedArticles), 1, "feed should contain articles from followed user")

	// --- User1 unfollows User2 ---
	w = e2eRequest(r, "DELETE", "/api/profiles/socialuser2/follow", "", token1)
	assert.Equal(http.StatusOK, w.Code, "unfollow should return 200")

	unfollowData := parseJSON(w.Body)
	unfollowProfile, ok := unfollowData["profile"].(map[string]interface{})
	require.True(ok)
	assert.Equal(false, unfollowProfile["following"], "following should be false after unfollow")

	// --- Feed should now be empty ---
	w = e2eRequest(r, "GET", "/api/articles/feed", "", token1)
	assert.Equal(http.StatusOK, w.Code)

	emptyFeedData := parseJSON(w.Body)
	emptyFeedArticles, ok := emptyFeedData["articles"].([]interface{})
	require.True(ok)
	assert.Equal(0, len(emptyFeedArticles), "feed should be empty after unfollowing")
}

// TestE2E_CommentFlow tests the comment lifecycle:
// Register -> Create article -> Add comment -> List comments -> Delete comment
func TestE2E_CommentFlow(t *testing.T) {
	setupE2EDB(t)

	r := setupE2ERouter()
	assert := assert.New(t)
	require := require.New(t)

	// --- Register user ---
	_, token := registerUser(r, "commenter", "commenter@test.com", "password123")
	require.NotEmpty(token)

	// --- Create article ---
	articleBody := `{"article":{"title":"Comment Test Article","description":"For comments","body":"Article body","tagList":[]}}`
	w := e2eRequest(r, "POST", "/api/articles", articleBody, token)
	assert.Equal(http.StatusCreated, w.Code)

	slug := extractArticleSlug(parseJSON(w.Body))
	require.NotEmpty(slug)

	// --- Add comment ---
	commentBody := `{"comment":{"body":"This is a great article!"}}`
	w = e2eRequest(r, "POST", "/api/articles/"+slug+"/comments", commentBody, token)
	assert.Equal(http.StatusCreated, w.Code, "add comment should return 201")

	commentData := parseJSON(w.Body)
	comment, ok := commentData["comment"].(map[string]interface{})
	require.True(ok, "response should contain comment object")
	assert.Equal("This is a great article!", comment["body"])

	commentAuthor, ok := comment["author"].(map[string]interface{})
	require.True(ok, "comment should have author")
	assert.Equal("commenter", commentAuthor["username"])

	commentID := comment["id"].(float64)
	require.NotZero(commentID, "comment should have an id")

	// --- Add second comment ---
	commentBody2 := `{"comment":{"body":"Second comment here"}}`
	w = e2eRequest(r, "POST", "/api/articles/"+slug+"/comments", commentBody2, token)
	assert.Equal(http.StatusCreated, w.Code)

	// --- List comments ---
	w = e2eRequest(r, "GET", "/api/articles/"+slug+"/comments", "", "")
	assert.Equal(http.StatusOK, w.Code, "list comments should return 200")

	commentsData := parseJSON(w.Body)
	comments, ok := commentsData["comments"].([]interface{})
	require.True(ok, "response should contain comments array")
	assert.Equal(2, len(comments), "should have 2 comments")

	// --- Delete first comment ---
	deleteURL := fmt.Sprintf("/api/articles/%s/comments/%d", slug, int(commentID))
	w = e2eRequest(r, "DELETE", deleteURL, "", token)
	assert.Equal(http.StatusOK, w.Code, "delete comment should return 200")

	// --- Verify only one comment remains ---
	w = e2eRequest(r, "GET", "/api/articles/"+slug+"/comments", "", "")
	assert.Equal(http.StatusOK, w.Code)

	remainingData := parseJSON(w.Body)
	remainingComments, ok := remainingData["comments"].([]interface{})
	require.True(ok)
	assert.Equal(1, len(remainingComments), "should have 1 comment after deletion")
}

// TestE2E_FavoriteFlow tests the favorite lifecycle:
// Register 2 users -> Create article -> Favorite -> List favorited -> Unfavorite
func TestE2E_FavoriteFlow(t *testing.T) {
	setupE2EDB(t)

	r := setupE2ERouter()
	assert := assert.New(t)
	require := require.New(t)

	// --- Register two users ---
	_, authorToken := registerUser(r, "favauthor", "favauthor@test.com", "password123")
	require.NotEmpty(authorToken)
	_, readerToken := registerUser(r, "favreader", "favreader@test.com", "password123")
	require.NotEmpty(readerToken)

	// --- Author creates article ---
	articleBody := `{"article":{"title":"Favorite Test Article","description":"Fav test","body":"Favorite body","tagList":["favorite"]}}`
	w := e2eRequest(r, "POST", "/api/articles", articleBody, authorToken)
	assert.Equal(http.StatusCreated, w.Code)

	slug := extractArticleSlug(parseJSON(w.Body))
	require.NotEmpty(slug)

	// --- Reader favorites article ---
	w = e2eRequest(r, "POST", "/api/articles/"+slug+"/favorite", "", readerToken)
	assert.Equal(http.StatusOK, w.Code, "favorite should return 200")

	favData := parseJSON(w.Body)
	favArticle, ok := favData["article"].(map[string]interface{})
	require.True(ok)
	assert.Equal(true, favArticle["favorited"], "article should be favorited")
	assert.Equal(float64(1), favArticle["favoritesCount"], "favoritesCount should be 1")

	// --- List articles favorited by reader ---
	w = e2eRequest(r, "GET", "/api/articles?favorited=favreader", "", readerToken)
	assert.Equal(http.StatusOK, w.Code, "list favorited should return 200")

	favListData := parseJSON(w.Body)
	favArticles, ok := favListData["articles"].([]interface{})
	require.True(ok)
	assert.GreaterOrEqual(len(favArticles), 1, "should have at least 1 favorited article")

	// Verify the first article is the one we favorited
	firstFav, ok := favArticles[0].(map[string]interface{})
	require.True(ok)
	assert.Equal(true, firstFav["favorited"])

	// --- Reader unfavorites article ---
	w = e2eRequest(r, "DELETE", "/api/articles/"+slug+"/favorite", "", readerToken)
	assert.Equal(http.StatusOK, w.Code, "unfavorite should return 200")

	unfavData := parseJSON(w.Body)
	unfavArticle, ok := unfavData["article"].(map[string]interface{})
	require.True(ok)
	assert.Equal(false, unfavArticle["favorited"], "article should not be favorited")
	assert.Equal(float64(0), unfavArticle["favoritesCount"], "favoritesCount should be 0")

	// --- Verify favorited list is empty ---
	w = e2eRequest(r, "GET", "/api/articles?favorited=favreader", "", readerToken)
	assert.Equal(http.StatusOK, w.Code)

	emptyFavData := parseJSON(w.Body)
	emptyFavArticles, ok := emptyFavData["articles"].([]interface{})
	require.True(ok)
	assert.Equal(0, len(emptyFavArticles), "favorited list should be empty after unfavorite")
}

// TestE2E_TagsEndpoint tests that tags are returned correctly:
// Create articles with different tags -> Get tags list
func TestE2E_TagsEndpoint(t *testing.T) {
	setupE2EDB(t)

	r := setupE2ERouter()
	assert := assert.New(t)
	require := require.New(t)

	// --- Register user ---
	_, token := registerUser(r, "taguser", "taguser@test.com", "password123")
	require.NotEmpty(token)

	// --- Create articles with different tags ---
	article1 := `{"article":{"title":"Tag Article One","description":"Desc1","body":"Body1","tagList":["golang","backend"]}}`
	w := e2eRequest(r, "POST", "/api/articles", article1, token)
	assert.Equal(http.StatusCreated, w.Code)

	article2 := `{"article":{"title":"Tag Article Two","description":"Desc2","body":"Body2","tagList":["python","backend"]}}`
	w = e2eRequest(r, "POST", "/api/articles", article2, token)
	assert.Equal(http.StatusCreated, w.Code)

	article3 := `{"article":{"title":"Tag Article Three","description":"Desc3","body":"Body3","tagList":["javascript","frontend"]}}`
	w = e2eRequest(r, "POST", "/api/articles", article3, token)
	assert.Equal(http.StatusCreated, w.Code)

	// --- Get tags ---
	w = e2eRequest(r, "GET", "/api/tags", "", "")
	assert.Equal(http.StatusOK, w.Code, "get tags should return 200")

	tagsData := parseJSON(w.Body)
	tags, ok := tagsData["tags"].([]interface{})
	require.True(ok, "response should contain tags array")
	assert.GreaterOrEqual(len(tags), 5, "should have at least 5 unique tags")

	// Convert to string slice for easier checking
	tagStrings := make([]string, len(tags))
	for i, t := range tags {
		tagStrings[i], _ = t.(string)
	}
	assert.Contains(tagStrings, "golang")
	assert.Contains(tagStrings, "backend")
	assert.Contains(tagStrings, "python")
	assert.Contains(tagStrings, "javascript")
	assert.Contains(tagStrings, "frontend")
}

// TestE2E_ErrorCases tests various error scenarios:
// Unauthorized access, invalid data, not found
func TestE2E_ErrorCases(t *testing.T) {
	setupE2EDB(t)

	r := setupE2ERouter()
	assert := assert.New(t)
	require := require.New(t)

	// --- Unauthorized access to protected endpoints ---
	t.Run("unauthorized_get_user", func(t *testing.T) {
		w := e2eRequest(r, "GET", "/api/user", "", "")
		assert.Equal(http.StatusUnauthorized, w.Code, "get user without token should return 401")
	})

	t.Run("unauthorized_create_article", func(t *testing.T) {
		body := `{"article":{"title":"Unauth Article","description":"desc","body":"body","tagList":[]}}`
		w := e2eRequest(r, "POST", "/api/articles", body, "")
		assert.Equal(http.StatusUnauthorized, w.Code, "create article without token should return 401")
	})

	t.Run("unauthorized_update_user", func(t *testing.T) {
		body := `{"user":{"username":"hacker"}}`
		w := e2eRequest(r, "PUT", "/api/user", body, "")
		assert.Equal(http.StatusUnauthorized, w.Code, "update user without token should return 401")
	})

	t.Run("unauthorized_follow", func(t *testing.T) {
		w := e2eRequest(r, "POST", "/api/profiles/someone/follow", "", "")
		assert.Equal(http.StatusUnauthorized, w.Code, "follow without token should return 401")
	})

	t.Run("unauthorized_favorite", func(t *testing.T) {
		w := e2eRequest(r, "POST", "/api/articles/some-slug/favorite", "", "")
		assert.Equal(http.StatusUnauthorized, w.Code, "favorite without token should return 401")
	})

	t.Run("unauthorized_feed", func(t *testing.T) {
		w := e2eRequest(r, "GET", "/api/articles/feed", "", "")
		assert.Equal(http.StatusUnauthorized, w.Code, "feed without token should return 401")
	})

	t.Run("invalid_token", func(t *testing.T) {
		w := e2eRequest(r, "GET", "/api/user", "", "invalid.jwt.token")
		assert.Equal(http.StatusUnauthorized, w.Code, "invalid token should return 401")
	})

	// --- Invalid registration data ---
	t.Run("register_short_username", func(t *testing.T) {
		body := `{"user":{"username":"ab","email":"short@test.com","password":"password123"}}`
		w := e2eRequest(r, "POST", "/api/users", body, "")
		assert.Equal(http.StatusUnprocessableEntity, w.Code, "short username should return 422")
	})

	t.Run("register_short_password", func(t *testing.T) {
		body := `{"user":{"username":"validuser","email":"valid@test.com","password":"short"}}`
		w := e2eRequest(r, "POST", "/api/users", body, "")
		assert.Equal(http.StatusUnprocessableEntity, w.Code, "short password should return 422")
	})

	t.Run("register_invalid_email", func(t *testing.T) {
		body := `{"user":{"username":"validuser","email":"notanemail","password":"password123"}}`
		w := e2eRequest(r, "POST", "/api/users", body, "")
		assert.Equal(http.StatusUnprocessableEntity, w.Code, "invalid email should return 422")
	})

	t.Run("register_duplicate_email", func(t *testing.T) {
		body := `{"user":{"username":"dupuser1","email":"dup@test.com","password":"password123"}}`
		w := e2eRequest(r, "POST", "/api/users", body, "")
		assert.Equal(http.StatusCreated, w.Code, "first registration should succeed")

		body = `{"user":{"username":"dupuser2","email":"dup@test.com","password":"password123"}}`
		w = e2eRequest(r, "POST", "/api/users", body, "")
		assert.Equal(http.StatusUnprocessableEntity, w.Code, "duplicate email should return 422")
	})

	// --- Login with wrong credentials ---
	t.Run("login_wrong_password", func(t *testing.T) {
		// Register first
		registerUser(r, "logintest", "logintest@test.com", "password123")

		body := `{"user":{"email":"logintest@test.com","password":"wrongpassword"}}`
		w := e2eRequest(r, "POST", "/api/users/login", body, "")
		assert.Equal(http.StatusUnauthorized, w.Code, "wrong password should return 401")
	})

	t.Run("login_nonexistent_email", func(t *testing.T) {
		body := `{"user":{"email":"nonexistent@test.com","password":"password123"}}`
		w := e2eRequest(r, "POST", "/api/users/login", body, "")
		assert.Equal(http.StatusUnauthorized, w.Code, "nonexistent email should return 401")
	})

	// --- Not found ---
	t.Run("article_not_found", func(t *testing.T) {
		w := e2eRequest(r, "GET", "/api/articles/nonexistent-slug", "", "")
		assert.Equal(http.StatusNotFound, w.Code, "nonexistent article should return 404")
	})

	t.Run("profile_not_found", func(t *testing.T) {
		w := e2eRequest(r, "GET", "/api/profiles/nonexistentuser", "", "")
		assert.Equal(http.StatusNotFound, w.Code, "nonexistent profile should return 404")
	})

	// --- Authorization: delete another user's article ---
	t.Run("delete_other_users_article", func(t *testing.T) {
		_, ownerToken := registerUser(r, "artowner", "artowner@test.com", "password123")
		require.NotEmpty(ownerToken)
		_, otherToken := registerUser(r, "artother", "artother@test.com", "password123")
		require.NotEmpty(otherToken)

		createBody := `{"article":{"title":"Owner Article","description":"desc","body":"body","tagList":[]}}`
		w := e2eRequest(r, "POST", "/api/articles", createBody, ownerToken)
		assert.Equal(http.StatusCreated, w.Code)

		slug := extractArticleSlug(parseJSON(w.Body))
		require.NotEmpty(slug)

		// Other user tries to delete
		w = e2eRequest(r, "DELETE", "/api/articles/"+slug, "", otherToken)
		assert.Equal(http.StatusForbidden, w.Code, "deleting another user's article should return 403")
	})

	// --- Authorization: update another user's article ---
	t.Run("update_other_users_article", func(t *testing.T) {
		_, ownerToken := registerUser(r, "updowner", "updowner@test.com", "password123")
		require.NotEmpty(ownerToken)
		_, otherToken := registerUser(r, "updother", "updother@test.com", "password123")
		require.NotEmpty(otherToken)

		createBody := `{"article":{"title":"Owner Only Article","description":"desc","body":"body","tagList":[]}}`
		w := e2eRequest(r, "POST", "/api/articles", createBody, ownerToken)
		assert.Equal(http.StatusCreated, w.Code)

		slug := extractArticleSlug(parseJSON(w.Body))
		require.NotEmpty(slug)

		// Other user tries to update
		updateBody := `{"article":{"title":"Hacked Title"}}`
		w = e2eRequest(r, "PUT", "/api/articles/"+slug, updateBody, otherToken)
		assert.Equal(http.StatusForbidden, w.Code, "updating another user's article should return 403")
	})

	// --- Authorization: delete another user's comment ---
	t.Run("delete_other_users_comment", func(t *testing.T) {
		_, commOwnerToken := registerUser(r, "commowner", "commowner@test.com", "password123")
		require.NotEmpty(commOwnerToken)
		_, commOtherToken := registerUser(r, "commother", "commother@test.com", "password123")
		require.NotEmpty(commOtherToken)

		// Create article
		createBody := `{"article":{"title":"Comment Auth Article","description":"desc","body":"body","tagList":[]}}`
		w := e2eRequest(r, "POST", "/api/articles", createBody, commOwnerToken)
		assert.Equal(http.StatusCreated, w.Code)

		slug := extractArticleSlug(parseJSON(w.Body))
		require.NotEmpty(slug)

		// Owner adds comment
		commentBody := `{"comment":{"body":"Owner's comment"}}`
		w = e2eRequest(r, "POST", "/api/articles/"+slug+"/comments", commentBody, commOwnerToken)
		assert.Equal(http.StatusCreated, w.Code)

		commentData := parseJSON(w.Body)
		comment, ok := commentData["comment"].(map[string]interface{})
		require.True(ok)
		commentID := int(comment["id"].(float64))

		// Other user tries to delete
		deleteURL := fmt.Sprintf("/api/articles/%s/comments/%d", slug, commentID)
		w = e2eRequest(r, "DELETE", deleteURL, "", commOtherToken)
		assert.Equal(http.StatusForbidden, w.Code, "deleting another user's comment should return 403")
	})
}

// TestE2E_Pagination tests limit and offset query parameters:
// Create multiple articles -> Test limit/offset
func TestE2E_Pagination(t *testing.T) {
	setupE2EDB(t)

	r := setupE2ERouter()
	assert := assert.New(t)
	require := require.New(t)

	// --- Register user ---
	_, token := registerUser(r, "paguser", "paguser@test.com", "password123")
	require.NotEmpty(token)

	// --- Create 10 articles ---
	for i := 1; i <= 10; i++ {
		body := fmt.Sprintf(`{"article":{"title":"Pagination Article %d","description":"Desc %d","body":"Body %d","tagList":["pagination"]}}`, i, i, i)
		w := e2eRequest(r, "POST", "/api/articles", body, token)
		assert.Equal(http.StatusCreated, w.Code, fmt.Sprintf("creating article %d should return 201", i))
	}

	// --- Test default listing (should return articles) ---
	w := e2eRequest(r, "GET", "/api/articles", "", "")
	assert.Equal(http.StatusOK, w.Code)

	allData := parseJSON(w.Body)
	allArticles, ok := allData["articles"].([]interface{})
	require.True(ok)
	totalCount, _ := allData["articlesCount"].(float64)
	assert.Equal(float64(10), totalCount, "total count should be 10")

	// --- Test with limit=3 ---
	w = e2eRequest(r, "GET", "/api/articles?limit=3", "", "")
	assert.Equal(http.StatusOK, w.Code)

	limitData := parseJSON(w.Body)
	limitArticles, ok := limitData["articles"].([]interface{})
	require.True(ok)
	assert.Equal(3, len(limitArticles), "should return 3 articles with limit=3")
	limitCount, _ := limitData["articlesCount"].(float64)
	assert.Equal(float64(10), limitCount, "total count should still be 10 with limit")

	// --- Test with limit=3&offset=3 ---
	w = e2eRequest(r, "GET", "/api/articles?limit=3&offset=3", "", "")
	assert.Equal(http.StatusOK, w.Code)

	offsetData := parseJSON(w.Body)
	offsetArticles, ok := offsetData["articles"].([]interface{})
	require.True(ok)
	assert.Equal(3, len(offsetArticles), "should return 3 articles with offset=3")

	// Verify offset returns different articles than first page
	if len(allArticles) >= 6 && len(limitArticles) >= 1 && len(offsetArticles) >= 1 {
		firstPageSlug := limitArticles[0].(map[string]interface{})["slug"].(string)
		secondPageSlug := offsetArticles[0].(map[string]interface{})["slug"].(string)
		assert.NotEqual(firstPageSlug, secondPageSlug, "offset should return different articles")
	}

	// --- Test with large offset (beyond available articles) ---
	w = e2eRequest(r, "GET", "/api/articles?limit=5&offset=100", "", "")
	assert.Equal(http.StatusOK, w.Code)

	emptyData := parseJSON(w.Body)
	emptyArticles, ok := emptyData["articles"].([]interface{})
	require.True(ok)
	assert.Equal(0, len(emptyArticles), "offset beyond total should return empty list")

	// --- Test with limit=1 to verify exact behavior ---
	w = e2eRequest(r, "GET", "/api/articles?limit=1", "", "")
	assert.Equal(http.StatusOK, w.Code)

	oneData := parseJSON(w.Body)
	oneArticles, ok := oneData["articles"].([]interface{})
	require.True(ok)
	assert.Equal(1, len(oneArticles), "limit=1 should return exactly 1 article")
}
