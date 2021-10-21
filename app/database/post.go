package database

import (
	"fiber-blog/config"
	"log"

	"gorm.io/gorm"
)

type Post struct {
	BaseModel
	Title   string `gorm:"title" json:"title"`
	Content string `gorm:"content" json:"content"`
}

const PostTable = "post"

func GetPost(id int64) (*Post, error) {
	post := Post{}
	if err := db().Table(PostTable).Select("*").Where("id=?", id).First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, config.ErrFoo
		}
		log.Println(err)
		return nil, err
	}
	return &post, nil
}

func SavePost(title, content string) error {
	post := Post{
		Title:   title,
		Content: content,
	}
	if err := db().Table(PostTable).Save(&post).Error; err != nil {
		log.Println(err)
		return err
	}
	log.Println(post)
	return nil
}
