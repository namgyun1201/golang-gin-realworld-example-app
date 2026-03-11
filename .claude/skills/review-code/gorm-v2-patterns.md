# GORM v2 Pattern Reference

Anti-patterns to detect and their correct alternatives.

## Eager Loading

**Wrong** (GORM v1 pattern):
```go
db.Model(&user).Related(&articles)
```

**Correct** (GORM v2):
```go
db.Preload("Articles").Find(&user)
```

## Many-to-Many Relationships

**Wrong**:
```go
db.Model(&article).Related(&tags, "Tags")
```

**Correct**:
```go
db.Model(&article).Association("Tags").Find(&tags)
```

## Struct/Map Updates

**Wrong** (only updates first field):
```go
db.Model(&user).Update("username", "newname")
// or single-field Update for multiple fields
```

**Correct** (updates all non-zero fields):
```go
db.Model(&user).Updates(UserModel{Username: "newname", Bio: "new bio"})
// or with map for zero-value fields
db.Model(&user).Updates(map[string]interface{}{"username": "newname", "bio": ""})
```

## Delete with Pointer

**Wrong**:
```go
db.Delete(UserModel{}, id)
```

**Correct**:
```go
db.Delete(&UserModel{}, id)
```

## Count Return Type

**Wrong**:
```go
var count int
db.Model(&ArticleModel{}).Count(&count)
```

**Correct**:
```go
var count int64
db.Model(&ArticleModel{}).Count(&count)
// Convert safely when needed for uint
if count > 0 {
    uintCount := uint(count)
}
```

## Association Mode

**Wrong** (direct field manipulation):
```go
article.Tags = append(article.Tags, newTag)
db.Save(&article)
```

**Correct** (use Association):
```go
db.Model(&article).Association("Tags").Append(&newTag)
db.Model(&article).Association("Tags").Delete(&oldTag)
db.Model(&article).Association("Tags").Replace(&newTags)
db.Model(&article).Association("Tags").Clear()
```

## Transaction Handling

**Correct pattern**:
```go
err := db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&article).Error; err != nil {
        return err
    }
    if err := tx.Model(&article).Association("Tags").Replace(&tags); err != nil {
        return err
    }
    return nil
})
```

## AutoMigrate

**Correct** (use pointer to model):
```go
db.AutoMigrate(&UserModel{})
db.AutoMigrate(&ArticleModel{})
```
