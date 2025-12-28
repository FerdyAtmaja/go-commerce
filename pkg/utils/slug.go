package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// GenerateSlug creates a URL-friendly slug from a string
func GenerateSlug(text string) string {
	// Convert to lowercase
	slug := strings.ToLower(text)
	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")
	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")
	return slug
}

// EnsureUniqueSlug generates a unique slug by appending counter if duplicate exists
func EnsureUniqueSlug(baseSlug string, checkExists func(string) bool) string {
	slug := baseSlug
	counter := 1
	
	for checkExists(slug) {
		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}
	
	return slug
}