package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Kitap struct {
	ID     uint   `gorm:"primaryKey"`
	Baslik string `gorm:"not null"`
	ISBN   string `gorm:"unique"`
	Yazar  string `gorm:"not null"`
}

var db *gorm.DB

// Kitapları listeleme fonksiyonu

func listele(c *gin.Context) {
	var kitaplar []Kitap
	err := db.Find(&kitaplar).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kitaplar)
}

// ID'ye göre Kitap listeleme fonksiyonu

func listeleAranan(c *gin.Context) {
	id := c.Param("id")
	var arananKitap Kitap
	err := db.Where("id = ?", id).First(&arananKitap).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, arananKitap)
}

// Kitap ekleme fonksiyonu

func ekle(c *gin.Context) {
	var kitap Kitap
	err := c.BindJSON(&kitap)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = db.Create(&kitap).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, kitap)
}

// Kitap güncelleme fonksiyonu

func guncelle(c *gin.Context) {
	id := c.Param("id")
	var kitap Kitap
	err := db.First(&kitap, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitap bulunamadı"})
		return
	}
	var updatedData Kitap
	err = c.BindJSON(&updatedData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	kitap.Baslik = updatedData.Baslik
	err = db.Save(&kitap).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kitap)
}

// Kitap silme fonksiyonu

func sil(c *gin.Context) {
	id := c.Param("id")
	var kitap Kitap
	err := db.First(&kitap, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitap bulunamadı"})
		return
	}
	err = db.Delete(&kitap).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kitap silindi"})
}

func main() {
	// Gorm veritabanı bağlantısı
	var err error
	db, err = gorm.Open(sqlite.Open("Kitaplar.db"), &gorm.Config{}) // db değişkenine değer ataması yapılıyor
	if err != nil {
		log.Fatal(err)
	}
	// kitaplar tablosunu oluşturma
	err = db.AutoMigrate(&Kitap{})
	if err != nil {
		log.Fatal(err)
	}

	// Gin router oluşturma
	r := gin.Default()

	// Kitaplar listeleme endpoint'i
	r.GET("/kitaplar", listele)
	// Kitap listeleme örneği: curl -X GET http://localhost:8080/kitaplar

	// Aranan kitabı listeleme endpoint'i
	r.GET("/kitaplar/:id", listeleAranan)
	// Aranan kitap listeleme örneği: curl -X GET http://localhost:8080/kitaplar/1

	// Kitap ekleme endpoint'i
	r.POST("/kitaplar", ekle)
	// Kitap ekleme örneği: curl -X POST -H "Content-Type: application/json" -d '{"baslik":"Suç ve Ceza","isbn":"9789750729869","yazar":"Fyodor Dostoyevski"}' http://localhost:8080/kitaplar

	// Kitap güncelleme endpoint'i
	r.PUT("/kitaplar/:id", guncelle)
	// Kitap güncelleme örneği: curl -X PUT -H "Content-Type: application/json" -d '{"baslik":"Savaş ve Barış","isbn":"9789750729869","yazar":"Lev Tolstoy"}' http://localhost:8080/kitaplar/1

	// Kitap silme endpoint'i
	r.DELETE("/kitaplar/:id", sil)
	// Kitap silme örneği: curl -X DELETE http://localhost:8080/kitaplar/1

	r.Run("localhost:8080")
}
