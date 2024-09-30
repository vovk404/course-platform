package entity

type Course struct {
	Id             string  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name           string  `json:"name" gorm:"index"`
	TeacherId      string  `json:"teacher_id" gorm:"index"`
	Author         string  `json:"author" gorm:"index"`
	Description    string  `json:"description"`
	Price          float32 `json:"price"`
	CourseLanguage string  `json:"course_language"`
}
