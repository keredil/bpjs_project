package controllers

import (
	"backend/config"
	"backend/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPrograms godoc
// @Summary Menampilkan semua program BPJS
// @Description Endpoint ini mengambil semua data program dari database bpjs_db
// @Tags Programs
// @Produce json
// @Success 200 {object} map[string]interface{} "Daftar program berhasil diambil"
// @Failure 500 {object} map[string]interface{} "Terjadi kesalahan pada server"
// @Router /programs [get]
func GetPrograms(c *gin.Context) {
	var programs []models.Program
	result := config.DB.Find(&programs)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": programs})
}

// CreateProgram godoc
// @Summary Menambahkan program BPJS baru
// @Description Endpoint ini digunakan untuk menambahkan program baru ke database
// @Tags Programs
// @Accept json
// @Produce json
// @Param program body models.Program true "Data program baru"
// @Success 201 {object} map[string]interface{} "Program berhasil ditambahkan"
// @Failure 400 {object} map[string]interface{} "Data tidak valid"
// @Failure 500 {object} map[string]interface{} "Kesalahan server"
// @Router /programs [post]
func CreateProgram(c *gin.Context) {
	var program models.Program
	if err := c.ShouldBindJSON(&program); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := config.DB.Create(&program)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Program berhasil ditambahkan",
		"data":    program,
	})
}

// UpdateProgram godoc
// @Summary Memperbarui program BPJS
// @Description Endpoint ini memperbarui data program berdasarkan ID
// @Tags Programs
// @Accept json
// @Produce json
// @Param id path int true "ID program"
// @Param program body models.Program true "Data program yang diperbarui"
// @Success 200 {object} map[string]interface{} "Program berhasil diperbarui"
// @Failure 400 {object} map[string]interface{} "Data tidak valid"
// @Failure 404 {object} map[string]interface{} "Program tidak ditemukan"
// @Router /programs/{id} [put]
func UpdateProgram(c *gin.Context) {
	id := c.Param("id")
	var program models.Program

	// Cek apakah program ada
	if err := config.DB.First(&program, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Program tidak ditemukan"})
		return
	}

	var updated models.Program
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	program.Name = updated.Name
	program.Deskripsi = updated.Deskripsi
	program.BentukManfaat = updated.BentukManfaat
	program.ManfaatLengkap = updated.ManfaatLengkap

	config.DB.Save(&program)

	c.JSON(http.StatusOK, gin.H{
		"message": "Program berhasil diperbarui",
		"data":    program,
	})
}

// DeleteProgram godoc
// @Summary Menghapus program BPJS
// @Description Endpoint ini digunakan untuk menghapus data program berdasarkan ID
// @Tags Programs
// @Param id path int true "ID program"
// @Success 200 {object} map[string]interface{} "Program berhasil dihapus"
// @Failure 404 {object} map[string]interface{} "Program tidak ditemukan"
// @Router /programs/{id} [delete]
func DeleteProgram(c *gin.Context) {
	id := c.Param("id")
	var program models.Program

	if err := config.DB.First(&program, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Program tidak ditemukan"})
		return
	}

	config.DB.Delete(&program)
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Program dengan ID %s berhasil dihapus", id),
	})
}
