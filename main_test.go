package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cobot00/actor-hit/models"
)

func TestRandomChoice(t *testing.T) {
	var images []models.Image
	images = append(images, models.Image{Name: "1", Type: "", Path: "", VoiceActor: ""})
	images = append(images, models.Image{Name: "2", Type: "", Path: "", VoiceActor: ""})
	images = append(images, models.Image{Name: "3", Type: "", Path: "", VoiceActor: ""})
	images = append(images, models.Image{Name: "4", Type: "", Path: "", VoiceActor: ""})
	images = append(images, models.Image{Name: "5", Type: "", Path: "", VoiceActor: ""})

	actual := randomChoice(images, 3)
	assert.Equal(t, 3, len(actual))
}
