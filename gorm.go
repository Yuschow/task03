package main

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string
	Posts []Post
}

type Post struct {
	gorm.Model
	UserID        uint
	Content       string
	Words         int
	CommentStatus string
	User          *User
	Comments      []Comment
}

type Comment struct {
	gorm.Model
	PostID  uint
	Content string
	Post    *Post
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	p.Words = len(p.Content)
	fmt.Println("BeforeCreate: 设置 Words字段值")
	return
}

func (c *Comment) BeforeDelete(tx *gorm.DB) (err error) {
	var count int64
	if err := tx.Model(&Comment{}).
		Where("post_id = ? AND deleted_at IS NULL", c.PostID).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 1 {
		return tx.Model(&Post{}).
			Where("id = ?", c.PostID).
			Update("comment_status", "无评论").Error
	}
	return nil
}

func taskGORM() {
	db, err := gorm.Open(sqlite.Open("task03.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		log.Fatal("migration failed:", err)
	}

	log.Println("Migration successful!")

	// 创建用户、帖子和评论
	// for i := 1; i <= 10; i++ {
	// 	user := User{
	// 		Name: fmt.Sprintf("User%d", i),
	// 		Posts: []Post{
	// 			{
	// 				Content: fmt.Sprintf("Post content %d", i),
	// 				Comments: []Comment{
	// 					{Content: fmt.Sprintf("Comment for post %d", i)},
	// 				},
	// 			},
	// 		},
	// 	}
	// 	// 一次性创建用户、其帖子和评论（利用 GORM 的关联功能）
	// 	db.Create(&user)
	// }
	var post Post
	db.Preload("Comments.Post.User").Where("id = ?", 14).Find(&post)
	for _, c := range post.Comments {
		fmt.Printf("user: %s, post: %s, words: %d, comment: %s", c.Post.User.Name, c.Post.Content, c.Post.Words, c.Content)
	}
	var comment Comment
	if err := db.First(&comment, 14).Error; err != nil {
		return
	}
	db.Delete(&comment)
}
